package Dao

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Utils"
	"github.com/antchfx/htmlquery"
)

// GetMyListings 获取用户的上架列表
func (d *Dao) GetMyListings() ([]Model.MyListingReponse, error) {
	Logger.Infof("获取用户 %s 的上架列表", d.GetUsername())

	var items []Model.MyListingReponse
	params := Param.Params{}
	params.SetString("count", "30")

	req, err := d.NewRequest(http.MethodGet, Constants.GetMyListings+"?"+params.ToUrl(), nil)
	if err != nil {
		return items, err
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return items, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return items, err
	}

	if resp.StatusCode != http.StatusOK {
		return items, Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return items, err
	}

	var response Model.GetMyListingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Error("JSON解析错误:", err)
		return items, err
	}

	items, err = parseSteamMarketHTML(response.ResultsHTML)
	if err != nil {
		Logger.Error("从html中解析上架物品失败:", err)
		return items, err
	}

	return items, nil
}

// Remove 删除上架物品
func (d *Dao) RemoveMyListings(creatorId string) error {
	Logger.Infof("删除用户 [%d] 的上架物品，creatorId: %s", d.GetSteamID(), creatorId)

	params := Param.Params{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		params.SetString("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}

	req, err := d.NewRequest(http.MethodPost, Constants.RemoveMyListings+"/"+creatorId, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/market", Constants.CommunityOrigin))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// Logger.Debug(string(body))

	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	return nil
}

func (d *Dao) RemoveAllMyListings() error {
	return nil
}

type buyResult struct {
	success          bool
	needConfirmation bool
	confirmationId   string
	error            error
}

func (d *Dao) buy(creatorId string, name string, buyerPrice float64, sellerReceivePrice float64, confirmation string) buyResult {
	Logger.Infof("用户[%s]购买物品[%s][%s][%.02f][%.02f][%s]: ", d.GetUsername(), creatorId, name, buyerPrice, sellerReceivePrice, confirmation)

	buyerPriceStr := strconv.FormatFloat(buyerPrice*100, 'f', 0, 64)
	sellerReceivePriceStr := strconv.FormatFloat(sellerReceivePrice*100, 'f', 0, 64)
	feeStr := strconv.FormatFloat(buyerPrice*100-sellerReceivePrice*100, 'f', 0, 64)

	params := Param.Params{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		params.SetString("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}
	params.SetString("currency", "23")
	params.SetString("subtotal", sellerReceivePriceStr)
	params.SetString("fee", feeStr)
	params.SetString("total", buyerPriceStr)
	params.SetString("quantity", "1")
	params.SetString("confirmation", confirmation)
	params.SetInt64("save_my_address", 0)

	req, err := d.NewRequest(http.MethodPost, Constants.BuyListing+"/"+creatorId, strings.NewReader(params.Encode()))
	if err != nil {
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            err,
		}
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/market/listings/440/%s", Constants.CommunityOrigin, name))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            err,
		}
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            err,
			}
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		reader = flate.NewReader(resp.Body)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            err,
		}
	}

	Logger.Debugf("[BuyListing][%s]HTTP响应状态码: %d，响应内容: %s", creatorId, resp.StatusCode, string(body))

	// 检测429状态码（访问频繁）
	if resp.StatusCode == http.StatusTooManyRequests {
		Logger.Warnf("用户 [%s] 购买物品遇到速率限制 (429)", d.GetUsername())
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("购买失败: %w", Errors.ErrRateLimited),
		}
	} else if resp.StatusCode == http.StatusNotAcceptable {
		var buyListingResp Model.BuyListingNeedConfirmationResponse
		if err := json.Unmarshal(body, &buyListingResp); err != nil {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            err,
			}
		}
		return buyResult{
			success:          buyListingResp.Success == 22,
			needConfirmation: buyListingResp.NeedConfirmation,
			confirmationId:   buyListingResp.Confirmation["confirmation_id"],
			error:            nil,
		}

	} else if resp.StatusCode == http.StatusOK {
		Logger.Infof("用户[%s]购买物品[creatorId: %s][confirmation: '%s']", d.GetUsername(), creatorId, confirmation)
		var buyListingResp Model.BuyListingResponse
		if err := json.Unmarshal(body, &buyListingResp); err != nil {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            err,
			}
		}

		if buyListingResp.WalletInfo.Success == 1 {
			Logger.Infof("用户[%s]购买物品[creatorId: %s][confirmation: '%s']成功", d.GetUsername(), creatorId, confirmation)
			return buyResult{
				success:          true,
				needConfirmation: false,
				confirmationId:   "",
				error:            nil,
			}
		} else {
			Logger.Infof("用户[%s]购买物品[creatorId: %s][confirmation: '%s']失败", d.GetUsername(), creatorId, confirmation)
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            fmt.Errorf("购买失败，错误码: %d", buyListingResp.WalletInfo.Success),
			}
		}
	} else if resp.StatusCode == http.StatusBadRequest {
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("购买失败，错误码: %d，错误信息: %s", resp.StatusCode, "返回为空"),
		}
	} else {
		var buyListingFailedResp Model.BuyListingFailedResponse
		if err := json.Unmarshal(body, &buyListingFailedResp); err != nil {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            err,
			}
		}
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("购买失败，错误码: %d，错误信息: %s", resp.StatusCode, buyListingFailedResp.Message),
		}
	}
}

// BuyListing 购买物品
func (d *Dao) BuyListing(creatorId string, name string, buyerPrice float64, sellerReceivePrice float64, confirmation string, maFileContent string) error {
	br := d.buy(creatorId, name, buyerPrice, sellerReceivePrice, confirmation)
	if br.success && br.needConfirmation {
		Logger.Infof("用户[%s]购买物品[%s][%s]需要手机令牌确认", d.GetUsername(), creatorId, name)
		for i := range Constants.Tries {
			if err := d.ConfirmationForBuyListAndOrder("allow", maFileContent); err != nil {
				if i == Constants.Tries-1 {
					return err
				}
			} else {
				break
			}
		}
		brAgain := d.buy(creatorId, name, buyerPrice, sellerReceivePrice, br.confirmationId)
		if brAgain.success {
			return nil
		} else {
			return brAgain.error
		}
	} else {
		Logger.Infof("用户[%s]购买物品[%s][%s][%+v]", d.GetUsername(), creatorId, name, br)
	}
	return br.error
}

func (d *Dao) createOrder(marketHashName string, price float64, quantity int64, confirmation string, maFileContent string) error {
	Logger.Infof("用户 [%s] 开始挂单，饰品名称: %s，数量：%d", d.GetUsername(), marketHashName, quantity)

	var createOrderResp Model.CreateOrderResponse

	priceTotalStr := strconv.FormatFloat(price*float64(quantity)*100, 'f', 0, 64)

	params := Param.Params{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		params.SetString("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}
	params.SetString("currency", "23")
	params.SetInt64("appid", int64(Constants.TeamFortress2))
	params.SetString("market_hash_name", marketHashName)
	params.SetString("price_total", priceTotalStr)
	params.SetInt64("tradefee_tax", 0)
	params.SetInt64("quantity", quantity)
	params.SetInt64("save_my_address", 0)
	params.SetString("confirmation", confirmation)

	req, err := d.NewRequest(http.MethodPost, Constants.CreateOrder, strings.NewReader(params.Encode()))
	if err != nil {

		return err
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/market/listings/440/%s", Constants.CommunityOrigin, marketHashName))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		Logger.Errorf("[CreateOrder] 创建订单请求时 RetryRequest 失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			Logger.Errorf("[CreateOrder] 创建订单请求时 NewReader 失败: %v", err)
			return err
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		reader = flate.NewReader(resp.Body)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		Logger.Errorf("[CreateOrder] 创建订单请求时 ReadAll 失败: %v", err)
		return err
	}

	// type CreateOrderResponse struct {
	// 	NeedConfirmation bool              `json:"need_confirmation"`
	// 	Confirmation     map[string]string `json:"confirmation"`
	// 	Success          int               `json:"success"`
	// }

	if resp.StatusCode != 429 {
		Logger.Debugf("[CreateOrder] HTTP响应状态码: %d, HTTP响应内容: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &createOrderResp); err != nil {
		Logger.Errorf("[CreateOrder] 创建订单请求时 Unmarshal 失败: %v", err)
		return err
	}

	Logger.Debugf("[CreateOrder] 创建订单响应: %+v", createOrderResp)

	if resp.StatusCode == http.StatusNotAcceptable && createOrderResp.Success == 22 {
		if createOrderResp.NeedConfirmation {
			for i := range Constants.Tries {
				Logger.Debugf("挂单需要手机令牌确认，第 %d 次尝试", i+1)
				if err := d.ConfirmationForBuyListAndOrder("allow", maFileContent); err != nil {
					if i == Constants.Tries-1 {
						Logger.Errorf("[CreateOrder] 创建订单请求时 ConfirmationForBuyListAndOrder 失败: %v", err)
						return err
					}
					continue
				}

				// if err := d.createOrder(marketHashName, price, quantity, createOrderResp.Confirmation["confirmation_id"], maFileContent); err != nil {
				// 	Logger.Errorf("[CreateOrder] 创建订单请求时 createOrder 失败: %v", err)
				// 	return err
				// }
			}

		}
	}

	return nil
}

func (d *Dao) CreateOrder(marketHashName string, price float64, quantity int64, maFileContent string) error {
	d.createOrder(marketHashName, price, quantity, "", maFileContent)
	return nil
}

// GetInventory 获取用户库存
func (d *Dao) GetInventory(gameId int, categoryId int) ([]Model.Item, error) {
	username := d.GetUsername()
	Logger.Infof("开始获取用户 [%s] 的库存，游戏ID: %d, 分类ID: %d", username, gameId, categoryId)

	inventoryUrl := fmt.Sprintf("%s/%d/%d/%d", Constants.GetInventory, d.GetSteamID(), gameId, categoryId)
	req, err := d.NewRequest(http.MethodGet, inventoryUrl, nil)
	if err != nil {
		Logger.Errorf("创建库存请求失败，用户: [%s], 错误: %v", username, err)
		return nil, fmt.Errorf("创建库存请求失败: %w", err)
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		Logger.Errorf("执行库存请求失败，用户: [%s], 错误: %v", username, err)
		return nil, fmt.Errorf("执行库存请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger.Errorf("读取库存响应失败，用户: [%s], 错误: %v", username, err)
		return nil, fmt.Errorf("读取库存响应失败: %w", err)
	}

	Logger.Debugf("库存响应: %s", string(body))

	// 检测429状态码（访问频繁）
	if resp.StatusCode == http.StatusTooManyRequests {
		Logger.Warnf("用户 [%s] 获取库存遇到速率限制 (429)", username)
		return nil, fmt.Errorf("获取库存失败: %w", Errors.ErrRateLimited)
	}

	// 检查是否为GZIP压缩数据
	if len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b {
		// 解压GZIP数据
		reader, _ := gzip.NewReader(bytes.NewReader(body))
		defer reader.Close()
		body, _ = io.ReadAll(reader)
	}

	var inventoryResponse Model.InventoryResponse
	if err := json.Unmarshal(body, &inventoryResponse); err != nil {
		Logger.Errorf("解析库存响应失败，用户: [%s], 错误: %v", username, err)
		return nil, fmt.Errorf("解析库存响应失败: %w", err)
	}

	if inventoryResponse.Success != 1 {
		Logger.Errorf("库存API返回失败，用户: [%s], success字段: %d", username, inventoryResponse.Success)
		return nil, fmt.Errorf("库存API返回失败，success=%d", inventoryResponse.Success)
	}

	// 转换为内部物品模型
	items := d.processInventoryData(&inventoryResponse, username)

	Logger.Infof("获取用户 [%s] 的库存完成，共找到 %d 个可交易物品", username, len(items))

	return items, nil
}

// processInventoryData 处理库存数据并返回可交易物品列表
func (d *Dao) processInventoryData(inventoryResponse *Model.InventoryResponse, username string) []Model.Item {
	// 边界检查
	if inventoryResponse == nil {
		Logger.Warnf("用户 [%s] 库存响应为空", username)
		return []Model.Item{}
	}

	if len(inventoryResponse.Assets) == 0 {
		Logger.Infof("用户 [%s] 的库存为空", username)
		return []Model.Item{}
	}

	Logger.Debugf("开始处理用户 [%s] 的库存数据，资产数量: %d, 描述数量: %d",
		username, len(inventoryResponse.Assets), len(inventoryResponse.Descriptions))

	// 预分配 map 容量
	descMap := make(map[string]Model.Description, len(inventoryResponse.Descriptions))

	// 构建描述映射
	for _, desc := range inventoryResponse.Descriptions {
		// 直接字符串拼接，避免 fmt.Sprintf 的开销
		key := desc.ClassID + "_" + desc.InstanceID
		descMap[key] = desc
	}

	// 预估可交易物品数量，预分配切片容量
	capacity := len(inventoryResponse.Assets)
	if capacity > 100 {
		capacity = capacity / 3 // 经验值：约1/3的物品可交易
	}
	items := make([]Model.Item, 0, capacity)

	// 统计计数器
	var tradableCount, marketableCount, filteredCount, missingDescCount int

	// 处理资产
	for _, asset := range inventoryResponse.Assets {
		key := asset.ClassID + "_" + asset.InstanceID
		if desc, exists := descMap[key]; exists {
			if desc.Tradable == 1 {
				tradableCount++
			}
			if desc.Marketable == 1 {
				marketableCount++
			}

			// 只包含可交易和可市场交易的物品
			if desc.Marketable == 1 {
				items = append(items, Model.Item{
					AssetID:    asset.AssetID,
					ClassID:    asset.ClassID,
					InstanceID: asset.InstanceID,
					Name:       desc.Name,
					MarketName: desc.MarketName,
					Tradable:   true,
					Marketable: true,
				})
				filteredCount++
			}
		} else {
			missingDescCount++
			Logger.Debugf("用户 [%s] 的资产 %s 缺少描述信息", username, asset.AssetID)
		}
	}

	// 如果有大量缺失描述的物品，记录警告
	if missingDescCount > len(inventoryResponse.Assets)/4 {
		Logger.Warnf("用户 [%s] 有 %d 个物品缺少描述信息，占总数的 %.1f%%",
			username, missingDescCount, float64(missingDescCount)*100/float64(len(inventoryResponse.Assets)))
	}

	Logger.Debugf("用户 [%s] 库存处理完成 - 总物品: %d, 可交易: %d, 可市场交易: %d, 筛选后: %d, 缺少描述: %d",
		username, len(inventoryResponse.Assets), tradableCount, marketableCount, filteredCount, missingDescCount)

	return items
}

// PutList 上架物品，需要二次手机令牌确认
func (d *Dao) PutList(assetID string, price float64, currency int, maFileContent string) ([]string, error) {
	Logger.Infof("用户 [%d] 上架物品，AssetID: %s, 价格: %.2f", d.GetSteamID(), assetID, price)

	data := url.Values{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		data.Set("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}
	data.Set("appid", strconv.Itoa(Constants.TeamFortress2))
	data.Set("contextid", "2") // 分类
	data.Set("assetid", assetID)
	data.Set("amount", "1")
	data.Set("price", strconv.FormatInt(int64(price*100), 10))

	req, err := d.Request(http.MethodPost, Constants.PutList, strings.NewReader(data.Encode()))
	if err != nil {
		return []string{}, err
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/profiles/%d/inventory", Constants.CommunityOrigin, d.GetSteamID()))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	Logger.Infof("用户 [%s] 上架物品，返回结果: %s，返回状态码: %d", d.GetUsername(), string(body), resp.StatusCode)

	// 检测429状态码（访问频繁）
	if resp.StatusCode == http.StatusTooManyRequests {
		Logger.Warnf("用户 [%d] 上架物品遇到速率限制 (429)", d.GetSteamID())
		return []string{}, fmt.Errorf("上架失败: %w", Errors.ErrRateLimited)
	}

	// 先行处理返回状态码不为200的情况
	if resp.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("上架失败: %v", string(body))
	}

	var sellResp Model.PutListResponse
	if err := json.Unmarshal(body, &sellResp); err != nil {
		return []string{}, fmt.Errorf("解析上架响应失败: %w", err)
	}

	// 再行处理返回数据不为成功的情况
	if !sellResp.Success {
		return []string{}, fmt.Errorf("上架失败: %v", sellResp)
	}

	// 如果需要手机令牌确认
	if sellResp.RequiresConfirmation == 1 && sellResp.NeedsMobileConfirmation {
		Logger.Infof("物品上架需要手机令牌确认，assetID: %s", assetID)
		// 进行手机令牌操作，多次尝试
		for i := range Constants.Tries {
			result := d.ConfirmationForPutList("allow", maFileContent)
			if !result.Success {
				if i == Constants.Tries-1 {
					return []string{}, fmt.Errorf("上架失败: %v", "尝试三次确认均未成功")
				}
				continue
			} else {
				return result.Result, nil
			}
		}
	} else {
		Logger.Warnf("无法进行确认操作，assetID: %s, RequiresConfirmation: %d", assetID, sellResp.RequiresConfirmation)
	}
	return []string{}, nil
}

func (d *Dao) ConfirmationForPutList(op string, maFileContent string) *Model.ConfirmationResult {
	username := d.GetUsername()
	Logger.Infof("开始获取用户 [%s] 待确认请求", username)

	pt, err := Utils.LoadMaFile(maFileContent)
	if err != nil {
		Logger.Errorf("加载 [%s] 令牌文件失败，错误： %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}

	steamTime, err := d.SteamTime()
	if err != nil {
		Logger.Errorf("获取 Steam 服务器时间失败，错误： %v", err)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}

	queryParams, err := Utils.GenerateConfirmationQueryParams(pt.MaFile.DeviceID, pt.MaFile.IdentitySecret, strconv.Itoa(int(pt.MaFile.Session.SteamID)), steamTime, "conf")
	if err != nil {
		Logger.Errorf("构建获取待确认请求参数失败，错误： %v", err)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}

	req, err := d.Request(http.MethodGet, Constants.GetConfirmationList+"?"+queryParams.ToUrl(), nil)
	if err != nil {
		Logger.Errorf("创建待确认请求失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		Logger.Errorf("执行待确认请求失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger.Errorf("读取待确认响应失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}

	var confirmResp Model.ConfirmationsResponse
	if err := json.Unmarshal(body, &confirmResp); err != nil {
		Logger.Errorf("解析待确认响应失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}

	Logger.Debugf("待确认响应:%+v", confirmResp)

	if !confirmResp.Success {
		Logger.Errorf("待确认API返回失败，用户: [%s], success字段: %t, 返回码：%d", username, confirmResp.Success, resp.StatusCode)
		return &Model.ConfirmationResult{
			Success: false,
			Result:  []string{},
		}
	}

	Logger.Infof("获取用户 [%s] 的待确认完成，共找到 %d 个待确认请求", username, len(confirmResp.Confirmations))

	// 初始化最终返回结果
	finalResult := &Model.ConfirmationResult{
		Success: true,
		Result:  []string{},
	}

	for _, conf := range confirmResp.Confirmations {
		if conf.Type != 3 {
			Logger.Infof("非上架确认不予处理:%+v", conf)
			continue
		}
		switch op {
		case "allow":
			err = d.AllowSingleConfirmation(pt, conf, steamTime)
			if err != nil {
				Logger.Errorf("允许待确认失败，用户: [%s], 错误: %v", username, err)
				return &Model.ConfirmationResult{
					Success: false,
					Result:  []string{},
				}
			}
			finalResult.Result = append(finalResult.Result, conf.CreatorID)
		case "cancel":
			err = d.CancelSingleConfirmation(pt, conf, steamTime)
			if err != nil {
				Logger.Errorf("拒绝待确认失败，用户: [%s], 错误: %v", username, err)
				return &Model.ConfirmationResult{
					Success: false,
					Result:  []string{},
				}
			}
		}
	}

	return finalResult
}

func (d *Dao) GetConfirmations(maFileContent string) error {
	username := d.GetUsername()
	pt, err := Utils.LoadMaFile(maFileContent)
	if err != nil {
		Logger.Errorf("加载 [%s] 令牌文件失败，错误： %v", username, err)
		return err
	}

	steamTime, err := d.SteamTime()
	if err != nil {
		Logger.Errorf("获取 Steam 服务器时间失败，错误： %v", err)
		return err
	}

	queryParams, err := Utils.GenerateConfirmationQueryParams(pt.MaFile.DeviceID, pt.MaFile.IdentitySecret, strconv.Itoa(int(pt.MaFile.Session.SteamID)), steamTime, "conf")
	if err != nil {
		Logger.Errorf("构建获取待确认请求参数失败，错误： %v", err)
		return err
	}

	req, err := d.Request(http.MethodGet, Constants.GetConfirmationList+"?"+queryParams.ToUrl(), nil)
	if err != nil {
		Logger.Errorf("创建待确认请求失败，用户: [%s], 错误: %v", username, err)
		return err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		Logger.Errorf("执行待确认请求失败，用户: [%s], 错误: %v", username, err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger.Errorf("读取待确认响应失败，用户: [%s], 错误: %v", username, err)
		return err
	}

	Logger.Infof("获取到的确认列表: %s", string(body))

	return nil
}

func (d *Dao) ConfirmationForBuyListAndOrder(op string, maFileContent string) error {
	username := d.GetUsername()
	Logger.Infof("开始获取用户 [%s] 待确认请求", username)

	pt, err := Utils.LoadMaFile(maFileContent)
	if err != nil {
		Logger.Errorf("加载 [%s] 令牌文件失败，错误： %v", username, err)
		return err
	}

	steamTime, err := d.SteamTime()
	if err != nil {
		Logger.Errorf("获取 Steam 服务器时间失败，错误： %v", err)
		return err
	}

	queryParams, err := Utils.GenerateConfirmationQueryParams(pt.MaFile.DeviceID, pt.MaFile.IdentitySecret, strconv.Itoa(int(pt.MaFile.Session.SteamID)), steamTime, "conf")
	if err != nil {
		Logger.Errorf("构建获取待确认请求参数失败，错误： %v", err)
		return err
	}

	req, err := d.Request(http.MethodGet, Constants.GetConfirmationList+"?"+queryParams.ToUrl(), nil)
	if err != nil {
		Logger.Errorf("创建待确认请求失败，用户: [%s], 错误: %v", username, err)
		return err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		Logger.Errorf("执行待确认请求失败，用户: [%s], 错误: %v", username, err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger.Errorf("读取待确认响应失败，用户: [%s], 错误: %v", username, err)
		return err
	}

	var confirmResp Model.ConfirmationsResponse
	if err := json.Unmarshal(body, &confirmResp); err != nil {
		Logger.Errorf("解析待确认响应失败，用户: [%s], 错误: %v", username, err)
		return fmt.Errorf("解析待确认响应失败: %v", err)
	}

	if !confirmResp.Success {
		Logger.Errorf("待确认API返回失败，用户: [%s], success字段: %t, 返回码：%d", username, confirmResp.Success, resp.StatusCode)
		return fmt.Errorf("待确认API返回失败")
	}

	Logger.Infof("获取用户 [%s] 的待确认完成，共找到 %d 个待确认请求", username, len(confirmResp.Confirmations))

	for i, conf := range confirmResp.Confirmations {
		Logger.Infof("confirmResp.Confirmations[%d] = %+v", i, conf)
	}

	for _, conf := range confirmResp.Confirmations {
		switch op {
		case "allow":
			err = d.AllowSingleConfirmation(pt, conf, steamTime)
			if err != nil {
				Logger.Errorf("允许待确认失败，用户: [%s], 错误: %v", username, err)
				return err
			}
		case "cancel":
			err = d.CancelSingleConfirmation(pt, conf, steamTime)
			if err != nil {
				Logger.Errorf("拒绝待确认失败，用户: [%s], 错误: %v", username, err)
				return err
			}
		}
	}

	return nil
}

func (d *Dao) AcceptConfirmations() {}

func (d *Dao) processSingleConfirmation(phoneToken *Utils.PhoneToken, conf Model.Confirmation, op string) error {
	Logger.Infof("处理用户 [%s] 确认请求，confID: %s，操作：%s", d.GetUsername(), conf.ID, op)

	steamTime, err := d.SteamTime()
	if err != nil {
		Logger.Errorf("获取 Steam 服务器时间失败，错误： %v", err)
		return err
	}

	params, err := phoneToken.GenerateConfirmationQueryParams(steamTime, op)
	if err != nil {
		return err
	}

	params.SetString("op", op)
	params.SetString("cid", conf.ID)
	params.SetString("ck", conf.Nonce)

	req, err := d.Request(http.MethodGet, Constants.Confirmation+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var acceptResp Model.ProcessConfirmationResponse
	if err := json.Unmarshal(body, &acceptResp); err != nil {
		return fmt.Errorf("解析接受确认响应失败: %w", err)
	}

	if !acceptResp.Success {
		return fmt.Errorf("接受确认失败")
	}

	return nil
}

func (d *Dao) AllowSingleConfirmation(phoneToken *Utils.PhoneToken, conf Model.Confirmation, timestamp int64) error {
	return d.processSingleConfirmation(phoneToken, conf, "allow")
}

func (d *Dao) CancelSingleConfirmation(phoneToken *Utils.PhoneToken, conf Model.Confirmation, timestamp int64) error {
	return d.processSingleConfirmation(phoneToken, conf, "cancel")
}

func parseSteamMarketHTMLWithXPath(htmlContent string) ([]Model.MyListingReponse, error) {
	var items []Model.MyListingReponse

	// 解析HTML文档
	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析HTML文档失败: %v", err)
	}

	// 使用XPath查找所有市场列表行
	listingRows := htmlquery.Find(doc, "//div[@class='market_listing_row market_recent_listing_row']")

	for _, row := range listingRows {
		item := Model.MyListingReponse{}

		// 提取Listing ID (从id属性中获取)
		if idAttr := htmlquery.SelectAttr(row, "id"); idAttr != "" {
			// 从 "mylisting_654831914932301672" 中提取数字部分
			if strings.HasPrefix(idAttr, "mylisting_") {
				item.ListingID = strings.TrimPrefix(idAttr, "mylisting_")
			}
		}

		// 提取物品名称 (使用正则表达式从HTML中提取URL中的英文名称)
		rowHTML := htmlquery.OutputHTML(row, false)
		nameRegex := regexp.MustCompile(`href="https://steamcommunity.com/market/listings/\d+/([^"]+)"`)
		nameMatch := nameRegex.FindStringSubmatch(rowHTML)
		if len(nameMatch) > 1 {
			// URL解码物品名称
			encodedName := nameMatch[1]
			if decodedName, err := url.QueryUnescape(encodedName); err == nil {
				item.MarketHashName = strings.TrimSpace(decodedName)
			} else {
				// 如果URL解码失败，使用简单的字符串替换
				decodedName := strings.ReplaceAll(encodedName, "%20", " ")
				decodedName = strings.ReplaceAll(decodedName, "%25", "%")
				decodedName = strings.ReplaceAll(decodedName, "%26", "&")
				decodedName = strings.ReplaceAll(decodedName, "%27", "'")
				item.MarketHashName = strings.TrimSpace(decodedName)
			}
		}

		// 提取买家价格 (第一个价格span)
		buyerPriceNode := htmlquery.FindOne(row, ".//span[@class='market_listing_price']/span/span[1]")
		if buyerPriceNode != nil {
			priceText := strings.TrimSpace(htmlquery.InnerText(buyerPriceNode))
			item.BuyerPrice = parsePrice(priceText)
		}

		// 提取卖家到账价格 (括号内的价格)
		sellerPriceNode := htmlquery.FindOne(row, ".//span[@class='market_listing_price']/span/span[2]")
		if sellerPriceNode != nil {
			priceText := strings.TrimSpace(htmlquery.InnerText(sellerPriceNode))
			// 移除括号
			priceText = strings.Trim(priceText, "()")
			item.SellerReceivePrice = parsePrice(priceText)
		}

		items = append(items, item)
	}

	return items, nil
}

// 保留原有的正则表达式方法作为备用
func parseSteamMarketHTML(htmlContent string) ([]Model.MyListingReponse, error) {
	// 首先尝试XPath方法
	items, err := parseSteamMarketHTMLWithXPath(htmlContent)
	if err == nil && len(items) > 0 {
		fmt.Println("XPath方法成功，获取到", len(items), "个物品")
		return items, nil
	}

	// XPath方法失败时使用正则表达式方法（支持中英文）
	return parseSteamMarketHTMLWithRegex(htmlContent)
}

// parseSteamMarketHTMLWithRegex 使用正则表达式解析（支持中英文）
func parseSteamMarketHTMLWithRegex(htmlContent string) ([]Model.MyListingReponse, error) {
	var items []Model.MyListingReponse

	// 检测是否为中文版本
	isChinese := strings.Contains(htmlContent, "我正在出售的物品") || strings.Contains(htmlContent, "这是买家所要支付")

	// 先找到所有唯一的listing ID
	listingIDRegex := regexp.MustCompile(`listing_(\d+)`)
	listingMatches := listingIDRegex.FindAllStringSubmatch(htmlContent, -1)

	// 使用map去重
	uniqueListingIDs := make(map[string]bool)
	for _, match := range listingMatches {
		if len(match) > 1 {
			uniqueListingIDs[match[1]] = true
		}
	}

	// 提取所有物品名称（从URL中获取英文名称）
	// 简单有效的方法：用简单的字符串操作
	nameMatches := extractItemNamesFromHTML(htmlContent)

	// 根据语言版本选择不同的正则表达式
	var priceMatches [][]string
	if isChinese {
		// 中文版价格提取
		buyerPriceRegex := regexp.MustCompile(`这是买家所要支付[^>]*>\s*([^<]+)\s*<`)
		priceMatches = buyerPriceRegex.FindAllStringSubmatch(htmlContent, -1)
	} else {
		// 英文版价格提取
		buyerPriceRegex := regexp.MustCompile(`This is the price the buyer pays[^>]*>\s*([^<]+)\s*<`)
		priceMatches = buyerPriceRegex.FindAllStringSubmatch(htmlContent, -1)
	}

	// 提取所有卖家到账价格
	sellerPriceRegex := regexp.MustCompile(`\(¥ ([^)]+)\)`)
	sellerMatches := sellerPriceRegex.FindAllStringSubmatch(htmlContent, -1)

	// 将listing ID转换为切片以便按顺序访问
	var listingIDs []string
	for listingID := range uniqueListingIDs {
		listingIDs = append(listingIDs, listingID)
	}

	// 为每个listing ID创建物品，按顺序分配信息
	for i, listingID := range listingIDs {
		item := Model.MyListingReponse{
			ListingID: listingID,
		}

		// 按顺序分配信息
		if i < len(nameMatches) {
			item.MarketHashName = strings.TrimSpace(nameMatches[i])
		}
		if i < len(priceMatches) {
			// 清理价格字符串
			price := strings.ReplaceAll(priceMatches[i][1], "\\n", "")
			price = strings.ReplaceAll(price, "\\t", "")
			price = strings.ReplaceAll(price, "\n", "")
			price = strings.ReplaceAll(price, "\t", "")
			item.BuyerPrice = parsePrice(price)
		}
		if i < len(sellerMatches) {
			item.SellerReceivePrice = parsePrice(sellerMatches[i][1])
		}

		items = append(items, item)
	}

	return items, nil
}

// extractItemNamesFromHTML 从HTML中提取物品的英文名称
func extractItemNamesFromHTML(htmlContent string) []string {
	var names []string

	// 直接在整个HTML内容中搜索，不按行分割
	searchPattern := "steamcommunity.com/market/listings/440/"
	content := htmlContent

	for {
		// 找到下一个市场链接
		pos := strings.Index(content, searchPattern)
		if pos == -1 {
			break
		}

		// 移动到物品名称开始位置
		startPos := pos + len(searchPattern)

		// 找到物品名称结束位置（引号）
		endPos := strings.Index(content[startPos:], "\"")
		if endPos != -1 {
			encodedName := content[startPos : startPos+endPos]
			if decodedName, err := url.QueryUnescape(encodedName); err == nil {
				// 清理物品名称，移除多余的字符
				cleanedName := strings.TrimSpace(decodedName)
				cleanedName = strings.TrimRight(cleanedName, "\\") // 移除末尾的反斜杠
				names = append(names, cleanedName)
			}
		}

		// 移动到下一个搜索位置
		content = content[startPos+1:]
	}

	return names
}

// parsePrice 从价格字符串中提取数字部分并转换为float64
func parsePrice(priceStr string) float64 {
	// 移除所有非数字和小数点的字符
	priceStr = strings.TrimSpace(priceStr)
	// 使用正则表达式提取数字部分（包括小数点）
	priceRegex := regexp.MustCompile(`(\d+\.?\d*)`)
	matches := priceRegex.FindStringSubmatch(priceStr)
	if len(matches) > 1 {
		if price, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return price
		}
	}
	return 0.0
}
