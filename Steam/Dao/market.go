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
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Utils"
)

// GetMyListings 获取用户的上架列表
// 返回两个列表：已上架的物品和等待确认的物品
func (d *Dao) GetMyListings() (activeListings []Model.MyListingReponse, err error) {
	Logger.Infof("获取用户 %s 的上架列表", d.GetUsername())
	params := Param.Params{}
	params.SetString("count", "50")

	req, err := d.NewRequest(http.MethodGet, Constants.GetMyListings+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, Errors.ResponseError(resp.StatusCode)
	}

	var response Model.GetMyListingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Error("JSON解析错误:", err)
		return nil, err
	}

	// Logger.Infof("[特殊打印	]获取用户[%s]的上架列表的html，数量: %s", d.GetUsername(), response.ResultsHTML)

	activeListings, pendingListings, err := parseSteamMarketHTML(response.ResultsHTML)
	if err != nil {
		Logger.Error("从html中解析上架物品失败:", err)
		return nil, err
	}

	Logger.Infof("[特殊打印]获取用户[%s]的已上架列表，数量: %d", d.GetUsername(), len(activeListings))
	for _, listing := range activeListings {
		Logger.Infof("[特殊打印]已上架列表: %+v", listing)
	}
	Logger.Infof("[特殊打印]获取用户[%s]的待确认上架列表，数量: %d", d.GetUsername(), len(pendingListings))
	for _, listing := range pendingListings {
		Logger.Infof("[特殊打印]待确认上架列表: %+v", listing)
	}

	for _, listing := range pendingListings {
		Logger.Infof("删除用户 [%s] 的等待确认物品，creatorId: %s", d.GetUsername(), listing.ListingID)
		err := d.RemoveMyListings(listing.ListingID)
		if err != nil {
			Logger.Errorf("删除listing失败 [%s]: %v", listing.ListingID, err)
			// 可以选择继续删除下一个，或者返回错误
			continue
		}
		Logger.Infof("成功删除listing [%s]", listing.ListingID)
		fmt.Printf("%+v\n", listing)
	}

	return activeListings, nil
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

func (d *Dao) buy(gameId string, creatorId string, name string, buyerPrice float64, sellerReceivePrice float64, confirmation string) buyResult {
	Logger.Infof("[%s]购买[%s][%s][%.02f][%.02f][%s]", d.GetUsername(), creatorId, name, buyerPrice, sellerReceivePrice, confirmation)

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
			error:            Errors.ErrNewRequest,
		}
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/market/listings/%s/%s", gameId, Constants.CommunityOrigin, name))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            Errors.ErrRetryRequest,
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
				error:            Errors.ErrGzipReader,
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
			error:            Errors.ErrIOReadAll,
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
	} else if resp.StatusCode == 502 {
		Logger.Warnf("用户 [%s] 购买物品遇到服务器错误 (502)", d.GetUsername())
		var buyListingFailedResp Model.BuyListingFailedResponse
		if err := json.Unmarshal(body, &buyListingFailedResp); err != nil {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            err,
			}
		}

		if buyListingFailedResp.Message == `Your account is currently unable to use the Community Market.` || buyListingFailedResp.Message == `您的帐户当前无法使用社区市场。` {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            Errors.ErrAccountBan,
			}
		} else {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            fmt.Errorf("%w", Errors.ErrServerError),
			}
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
func (d *Dao) BuyListing(gameId, creatorId string, name string, buyerPrice float64, sellerReceivePrice float64, confirmation string, maFileContent string) error {
	br := d.buy(gameId, creatorId, name, buyerPrice, sellerReceivePrice, confirmation)
	if br.success && br.needConfirmation {
		if err := d.ConfirmationForBuyList("allow", maFileContent); err != nil {
			return err
		}
		brAgain := d.buy(gameId, creatorId, name, buyerPrice, sellerReceivePrice, br.confirmationId)
		if brAgain.success {
			return nil
		} else {
			return brAgain.error
		}
	}
	return br.error
}

func (d *Dao) createOrder(gameId int, marketHashName string, price float64, quantity int64, confirmation string, maFileContent string) error {
	Logger.Infof("用户 [%s] 开始挂单，饰品名称: %s，数量：%d", d.GetUsername(), marketHashName, quantity)

	var createOrderResp Model.CreateOrderResponse

	priceTotalStr := strconv.FormatFloat(price*float64(quantity)*100, 'f', 0, 64)

	params := Param.Params{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		params.SetString("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}
	params.SetString("currency", "23")
	params.SetInt64("appid", int64(gameId))
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
	d.createOrder(570, marketHashName, price, quantity, "", maFileContent)
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

	// 检查是否为GZIP压缩数据
	if len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b {
		// 解压GZIP数据
		reader, _ := gzip.NewReader(bytes.NewReader(body))
		defer reader.Close()
		body, _ = io.ReadAll(reader)
	}

	Logger.Infof("库存响应状态码: %d", resp.StatusCode)
	Logger.Infof("库存响应: %s", string(body))

	switch resp.StatusCode {
	case 429:
		Logger.Warnf("用户 [%s] 获取库存遇到速率限制 (429)", username)
		return nil, fmt.Errorf("获取库存失败: %w", Errors.ErrRateLimited)
	case 401, 403:
		Logger.Warnf("用户 [%s] 获取库存遇到授权失败 (401/403)", username)
		return nil, fmt.Errorf("获取库存失败: %w", Errors.ErrAuthorizationFailed)
	}

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("获取库存失败，错误码: %d", resp.StatusCode)
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
func (d *Dao) PutList(gameId int, contextId int, assetID string, price float64, currency int, maFileContent string) (Model.MyListingReponse, error) {
	Logger.Infof("用户 [%d] 上架物品，AssetID: %s, 价格: %.2f", d.GetSteamID(), assetID, price)

	data := url.Values{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		data.Set("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}
	data.Set("appid", strconv.Itoa(gameId))
	data.Set("contextid", strconv.Itoa(contextId)) // 分类
	data.Set("assetid", assetID)
	data.Set("amount", "1")
	data.Set("price", strconv.FormatInt(int64(price*100), 10))

	req, err := d.Request(http.MethodPost, Constants.PutList, strings.NewReader(data.Encode()))
	if err != nil {
		return Model.MyListingReponse{}, err
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/profiles/%d/inventory", Constants.CommunityOrigin, d.GetSteamID()))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return Model.MyListingReponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Model.MyListingReponse{}, err
	}

	Logger.Debugf("[PutListing][%s]HTTP响应状态码: %d，响应内容: %s", assetID, resp.StatusCode, string(body))

	// 检测429状态码（访问频繁）
	if resp.StatusCode == http.StatusTooManyRequests {
		Logger.Warnf("用户 [%d] 上架物品遇到速率限制 (429)", d.GetSteamID())
		return Model.MyListingReponse{}, fmt.Errorf("上架失败: %w", Errors.ErrRateLimited)
	}

	// 先行处理返回状态码不为200的情况
	if resp.StatusCode != http.StatusOK {
		return Model.MyListingReponse{}, fmt.Errorf("上架失败: %v", string(body))
	}

	var sellResp Model.PutListResponse
	if err := json.Unmarshal(body, &sellResp); err != nil {
		return Model.MyListingReponse{}, fmt.Errorf("解析上架响应失败: %w", err)
	}

	// 再行处理返回数据不为成功的情况
	if !sellResp.Success {
		switch sellResp.Message {
		case "您的帐户当前无法使用社区市场。":
		case "Your account is currently unable to use the Community Market.":
			return Model.MyListingReponse{}, Errors.ErrAccountBan
		default:
			return Model.MyListingReponse{}, fmt.Errorf("%s", sellResp.Message)
		}

	}

	// 如果需要手机令牌确认
	if sellResp.RequiresConfirmation == 1 && sellResp.NeedsMobileConfirmation {
		Logger.Infof("物品上架需要手机令牌确认，assetID: %s", assetID)
		result := d.ConfirmationForPutList("allow", maFileContent)
		if !result.Success {
			return Model.MyListingReponse{}, fmt.Errorf("上架确认失败: %s", assetID)
		} else {
			return result.Result, nil
		}
	} else {
		Logger.Warnf("无法进行确认操作，assetID: %s, RequiresConfirmation: %d", assetID, sellResp.RequiresConfirmation)
	}
	return Model.MyListingReponse{}, nil
}

func (d *Dao) ConfirmationForPutList(op string, maFileContent string) *Model.ConfirmationResult {
	username := d.GetUsername()
	Logger.Infof("开始获取用户 [%s] 待确认请求", username)

	pt, err := Utils.LoadMaFile(maFileContent)
	if err != nil {
		Logger.Errorf("加载 [%s] 令牌文件失败，错误： %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}

	steamTime, err := d.GetSteamTimeLocal()
	if err != nil {
		Logger.Errorf("获取 Steam 服务器时间失败，错误： %v", err)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}

	queryParams, err := Utils.GenerateConfirmationQueryParams(pt.MaFile.DeviceID, pt.MaFile.IdentitySecret, strconv.Itoa(int(pt.MaFile.Session.SteamID)), steamTime, "conf")
	if err != nil {
		Logger.Errorf("构建获取待确认请求参数失败，错误： %v", err)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}

	req, err := d.Request(http.MethodGet, Constants.GetConfirmationList+"?"+queryParams.ToUrl(), nil)
	if err != nil {
		Logger.Errorf("创建待确认请求失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		Logger.Errorf("执行待确认请求失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger.Errorf("读取待确认响应失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}

	var confirmResp Model.ConfirmationsResponse
	if err := json.Unmarshal(body, &confirmResp); err != nil {
		Logger.Errorf("解析待确认响应失败，用户: [%s], 错误: %v", username, err)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}

	Logger.Debugf("获取用户 [%s] 的待确认列表响应:%+v,返回内容：%s", username, confirmResp, string(body))

	if !confirmResp.Success {
		Logger.Errorf("待确认API返回失败，用户: [%s], success字段: %t, 返回码：%d", username, confirmResp.Success, resp.StatusCode)
		return &Model.ConfirmationResult{
			Success: false,
		}
	}

	Logger.Infof("获取用户 [%s] 的待确认完成，共找到 %d 个待确认请求", username, len(confirmResp.Confirmations))

	// 初始化最终返回结果
	finalResult := &Model.ConfirmationResult{
		Success: false,
	}

	for i := len(confirmResp.Confirmations) - 1; i >= 0; i-- {
		Logger.Infof("处理第 %d 个确认项", i+1)
		conf := confirmResp.Confirmations[i]
		if conf.Type != 3 {
			Logger.Infof("非上架饰品确认不予处理:%+v", conf)
			continue
		}

		if i != 0 {
			Logger.Infof("处理其他购买饰品确认")
			sTime, _ := d.GetSteamTimeLocal()
			err := d.AllowSingleConfirmation(pt, conf, sTime)
			if err != nil {
				Logger.Errorf("处理其他购买饰品确认失败，用户: [%s], 错误: %v", username, err)
			}
			Logger.Errorf("处理其他购买饰品确认成功，用户: [%s]", username)
		} else {
			Logger.Infof("处理本次购买饰品确认")
			for j := 0; j < Constants.Tries; j++ {
				sTime, _ := d.GetSteamTimeLocal()
				err = d.AllowSingleConfirmation(pt, conf, sTime)
				if err != nil {
					Logger.Errorf("第 %d 次允许待确认失败，用户: [%s], 错误: %v", j+1, username, err)
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					Logger.Infof("处理成功本次上架确认:%+v", conf)
					finalResult.Success = true
					return finalResult
				}

			}
		}

		// price, ok := utils.ExtractPrice(conf.Headline)
		// if !ok {
		// 	Logger.Errorf("提取价格失败，用户: [%s], 错误: %v", username, err)
		// }
		// Logger.Infof("提取价格成功，用户: [%s], 价格: %.2f", username, price)

		// finalResult.Result = Model.MyListingReponse{
		// 	ListingID:          conf.CreatorID,
		// 	MarketHashName:     conf.Summary[0],
		// 	BuyerPrice:         0,
		// 	SellerReceivePrice: 0,
		// }

		// switch op {
		// case "allow":
		// 	err = d.AllowSingleConfirmation(pt, conf, steamTime)
		// 	if err != nil {
		// 		Logger.Errorf("允许待确认失败，用户: [%s], 错误: %v", username, err)
		// 		return &Model.ConfirmationResult{
		// 			Success: false,
		// 			Result:  []string{},
		// 		}
		// 	}
		// 	finalResult.Result = append(finalResult.Result, conf.CreatorID)
		// case "cancel":
		// 	err = d.CancelSingleConfirmation(pt, conf, steamTime)
		// 	if err != nil {
		// 		Logger.Errorf("拒绝待确认失败，用户: [%s], 错误: %v", username, err)
		// 		return &Model.ConfirmationResult{
		// 			Success: false,
		// 			Result:  []string{},
		// 		}
		// 	}
		// }

	}

	return finalResult
}

func (d *Dao) GetConfirmations(maFileContent string) error {
	Logger.Infof("开始获取用户 [%s] 的待确认请求", d.GetUsername())
	username := d.GetUsername()
	pt, err := Utils.LoadMaFile(maFileContent)
	if err != nil {
		Logger.Errorf("加载 [%s] 令牌文件失败，错误： %v", username, err)
		return err
	}

	steamTime, err := d.GetSteamTimeLocal()
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

	req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 9; Valve Steam App Version/3)")
	req.Header.Set("mobileClient", "android")
	req.Header.Set("mobileClientVersion", "777777 3.6.4")

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

	// d.ConfirmationForPutList("allow", maFileContent)
	// d.ConfirmationForBuyListAndOrder("allow", maFileContent)

	return nil
}

func (d *Dao) ConfirmationForBuyList(op string, maFileContent string) error {
	username := d.GetUsername()
	Logger.Infof("开始获取用户 [%s] 的购买饰品待确认请求", username)

	pt, err := Utils.LoadMaFile(maFileContent)
	if err != nil {
		Logger.Errorf("加载 [%s] 令牌文件失败，错误： %v", username, err)
		return err
	}

	steamTime, err := d.GetSteamTimeLocal()
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

	for _, conf := range confirmResp.Confirmations {
		if conf.Type != 12 {
			continue
		}

		for i := range Constants.Tries {
			err = d.AllowSingleConfirmation(pt, conf, steamTime)
			if err != nil {
				if i == Constants.Tries-1 {
					Logger.Errorf("最终购买饰品确认失败，用户: [%s], 错误: %v", username, err)
					return err
				}
				time.Sleep(100 * time.Millisecond)
				Logger.Errorf("第 %d 次购买饰品确认失败，用户: [%s], 错误: %v", i, username, err)
				continue
			} else {
				Logger.Infof("处理成功本次购买饰品确认:%+v", conf)
				break
			}
		}

		// switch op {
		// case "allow":
		// 	err = d.AllowSingleConfirmation(pt, conf, steamTime)
		// 	if err != nil {
		// 		Logger.Errorf("允许待确认失败，用户: [%s], 错误: %v", username, err)
		// 		return err
		// 	}
		// case "cancel":
		// 	err = d.CancelSingleConfirmation(pt, conf, steamTime)
		// 	if err != nil {
		// 		Logger.Errorf("拒绝待确认失败，用户: [%s], 错误: %v", username, err)
		// 		return err
		// 	}
		// }
	}

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

	steamTime, err := d.GetSteamTimeLocal()
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
		if conf.Type != 12 {
			continue
		}
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
	steamTime, err := d.GetSteamTimeLocal()
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

// 保留原有的正则表达式方法作为备用
// 返回两个列表：已上架的物品和等待确认的物品
func parseSteamMarketHTML(htmlContent string) (activeListings []Model.MyListingReponse, pendingListings []Model.MyListingReponse, err error) {
	// 首先尝试XPath方法
	// activeItems, pendingItems, err := parseSteamMarketHTMLWithXPath(htmlContent)
	// if err == nil && len(activeItems) > 0 {
	// 	fmt.Println("XPath方法成功，获取到", len(activeItems), "个已上架物品，", len(pendingItems), "个等待确认物品")
	// 	return activeItems, pendingItems, nil
	// }

	// XPath方法失败时使用正则表达式方法（支持中英文）
	return parseSteamMarketHTMLWithRegex(htmlContent)
}

// parseSteamMarketHTMLWithRegex 使用正则表达式解析（支持中英文）
// 返回两个列表：已上架的物品和等待确认的物品
func parseSteamMarketHTMLWithRegex(htmlContent string) (activeListings []Model.MyListingReponse, pendingListings []Model.MyListingReponse, err error) {
	// 检测是否为中文版本
	isChinese := strings.Contains(htmlContent, "我正在出售的物品") || strings.Contains(htmlContent, "这是买家所要支付")

	// 解析已上架物品
	activeListings = parseListingsFromSection(htmlContent, "tabContentsMyActiveMarketListingsRows", isChinese, true)
	Logger.Infof("正则表达式方法共解析到 %d 个已上架物品", len(activeListings))

	// 解析等待确认的物品
	// 等待确认的物品通常紧跟在已上架物品后面，在同一个 table 中，但不在 tabContentsMyActiveMarketListingsRows 内
	// 我们需要搜索包含 "My listings awaiting confirmation" 的区域
	pendingListings = parseListingsFromHTMLByKeyword(htmlContent, "My listings awaiting confirmation", "我的等待确认的上架物品", isChinese, false)
	Logger.Infof("正则表达式方法共解析到 %d 个等待确认物品", len(pendingListings))

	return activeListings, pendingListings, nil
}

// parseListingsFromSection 从指定的section中解析物品列表
func parseListingsFromSection(htmlContent string, sectionID string, isChinese bool, isActiveListing bool) []Model.MyListingReponse {
	var items []Model.MyListingReponse

	// 提取指定区域的HTML内容
	// 使用字符串查找而不是正则表达式，避免嵌套 div 的问题
	startMarker := `<div id="` + sectionID + `">`
	startIdx := strings.Index(htmlContent, startMarker)
	if startIdx == -1 {
		Logger.Debugf("未找到区域开始标记 (%s)", sectionID)
		return items
	}

	// 从开始标记后开始搜索
	contentStart := startIdx + len(startMarker)

	// 找到这个 div 的结束标签 - 搜索到下一个 section 的开始
	endMarkers := []string{
		`<div class="my_listing_section`, // 下一个 section
		`</div>`,                         // 如果是最后一个 section
	}

	endIdx := len(htmlContent)
	// 寻找最早出现的结束标记
	for _, marker := range endMarkers {
		if idx := strings.Index(htmlContent[contentStart:], marker); idx != -1 {
			potentialEnd := contentStart + idx
			if potentialEnd < endIdx {
				endIdx = potentialEnd
			}
			break // 找到第一个就停止
		}
	}

	targetHTML := htmlContent[contentStart:endIdx]
	Logger.Debugf("成功提取%s区域，长度: %d", sectionID, len(targetHTML))

	// 按行分割每个listing item（更可靠的方法）
	// 先找到所有 mylisting_xxx 的 id 和位置，然后按区域分割
	listingIDRegex := regexp.MustCompile(`id="mylisting_(\d+)"`)
	listingIDMatches := listingIDRegex.FindAllStringSubmatchIndex(targetHTML, -1)

	var rowMatches [][]string
	for i, match := range listingIDMatches {
		listingID := targetHTML[match[2]:match[3]]

		// 确定这个 listing 的起始位置（回溯找到 <div 开始）
		startPos := match[0]
		for startPos > 0 && targetHTML[startPos-1] != '<' {
			startPos--
		}
		if startPos > 0 {
			startPos-- // 包含 '<'
		}

		// 确定这个 listing 的结束位置（找到下一个 mylisting 或区域结束）
		endPos := len(targetHTML)
		if i+1 < len(listingIDMatches) {
			endPos = listingIDMatches[i+1][0]
		}

		rowHTML := targetHTML[startPos:endPos]
		rowMatches = append(rowMatches, []string{rowHTML, listingID})
	}

	for _, rowMatch := range rowMatches {
		if len(rowMatch) < 2 {
			continue
		}

		listingID := rowMatch[1]
		rowHTML := rowMatch[0] // 完整的匹配内容

		// 根据 isActiveListing 参数决定是否检查 RemoveMarketListing
		if isActiveListing {
			// 已上架物品：只处理包含 RemoveMarketListing 的物品
			if !strings.Contains(rowHTML, "RemoveMarketListing") {
				Logger.Debugf("跳过非已上架物品，Listing ID: %s", listingID)
				continue
			}
		} else {
			// 等待确认物品：只处理包含 CancelMarketListingConfirmation 的物品
			if !strings.Contains(rowHTML, "CancelMarketListingConfirmation") {
				Logger.Debugf("跳过非等待确认物品，Listing ID: %s", listingID)
				continue
			}
		}

		item := Model.MyListingReponse{
			ListingID: listingID,
		}

		// 提取Asset ID
		var assetIDRegex *regexp.Regexp
		if isActiveListing {
			assetIDRegex = regexp.MustCompile(`RemoveMarketListing\('mylisting',\s*'[^']+',\s*\d+,\s*'[^']+',\s*'([^']+)'\)`)
		} else {
			assetIDRegex = regexp.MustCompile(`CancelMarketListingConfirmation\('mylisting',\s*'[^']+',\s*\d+,\s*'[^']+',\s*'([^']+)'\)`)
		}
		if assetIDMatch := assetIDRegex.FindStringSubmatch(rowHTML); len(assetIDMatch) > 1 {
			item.AssetID = assetIDMatch[1]
		}

		// 提取物品名称（从market listings URL中获取）
		nameRegex := regexp.MustCompile(`href="https://steamcommunity\.com/market/listings/\d+/([^"]+)"`)
		if nameMatch := nameRegex.FindStringSubmatch(rowHTML); len(nameMatch) > 1 {
			if decodedName, err := url.QueryUnescape(nameMatch[1]); err == nil {
				item.MarketHashName = strings.TrimSpace(decodedName)
			} else {
				item.MarketHashName = strings.TrimSpace(nameMatch[1])
			}
		}

		// 提取买家价格
		var buyerPriceRegex *regexp.Regexp
		if isChinese {
			buyerPriceRegex = regexp.MustCompile(`这是买家所要支付[^>]*>\s*([^<]+)\s*<`)
		} else {
			buyerPriceRegex = regexp.MustCompile(`This is the price the buyer pays[^>]*>\s*([^<]+)\s*<`)
		}
		if priceMatch := buyerPriceRegex.FindStringSubmatch(rowHTML); len(priceMatch) > 1 {
			priceStr := strings.ReplaceAll(priceMatch[1], "\\n", "")
			priceStr = strings.ReplaceAll(priceStr, "\\t", "")
			priceStr = strings.ReplaceAll(priceStr, "\n", "")
			priceStr = strings.ReplaceAll(priceStr, "\t", "")
			item.BuyerPrice = parsePrice(strings.TrimSpace(priceStr))
		}

		// 提取卖家到账价格
		sellerPriceRegex := regexp.MustCompile(`\(¥\s*([^)]+)\)`)
		if sellerMatch := sellerPriceRegex.FindStringSubmatch(rowHTML); len(sellerMatch) > 1 {
			item.SellerReceivePrice = parsePrice(sellerMatch[1])
		}

		items = append(items, item)
	}

	return items
}

// parseListingsFromHTMLByKeyword 通过关键词在HTML中查找区域并解析物品列表
func parseListingsFromHTMLByKeyword(htmlContent string, englishKeyword string, chineseKeyword string, isChinese bool, isActiveListing bool) []Model.MyListingReponse {
	var items []Model.MyListingReponse

	// 根据语言选择关键词
	keyword := englishKeyword
	if isChinese && chineseKeyword != "" {
		keyword = chineseKeyword
	}

	// 找到关键词的位置
	keywordIdx := strings.Index(htmlContent, keyword)
	if keywordIdx == -1 {
		Logger.Debugf("未找到关键词: %s", keyword)
		return items
	}

	// 从关键词位置向后查找，找到这个区域的内容
	// 查找从关键词开始到下一个 my_listing_section 或文档结束的内容
	startIdx := keywordIdx
	endIdx := len(htmlContent)

	// 查找下一个section或文档结束
	nextSectionIdx := strings.Index(htmlContent[startIdx+len(keyword):], `<div class="my_listing_section`)
	if nextSectionIdx != -1 {
		endIdx = startIdx + len(keyword) + nextSectionIdx
	}

	targetHTML := htmlContent[startIdx:endIdx]
	Logger.Debugf("成功提取关键词区域 (%s)，长度: %d", keyword, len(targetHTML))

	// 按行分割每个listing item
	listingIDRegex := regexp.MustCompile(`id="mylisting_(\d+)"`)
	listingIDMatches := listingIDRegex.FindAllStringSubmatchIndex(targetHTML, -1)

	var rowMatches [][]string
	for i, match := range listingIDMatches {
		listingID := targetHTML[match[2]:match[3]]

		// 确定这个 listing 的起始位置
		startPos := match[0]
		for startPos > 0 && targetHTML[startPos-1] != '<' {
			startPos--
		}
		if startPos > 0 {
			startPos--
		}

		// 确定这个 listing 的结束位置
		endPos := len(targetHTML)
		if i+1 < len(listingIDMatches) {
			endPos = listingIDMatches[i+1][0]
		}

		rowHTML := targetHTML[startPos:endPos]
		rowMatches = append(rowMatches, []string{rowHTML, listingID})
	}

	for _, rowMatch := range rowMatches {
		if len(rowMatch) < 2 {
			continue
		}

		listingID := rowMatch[1]
		rowHTML := rowMatch[0]

		// 根据 isActiveListing 参数决定是否检查对应的JavaScript函数
		if isActiveListing {
			if !strings.Contains(rowHTML, "RemoveMarketListing") {
				Logger.Debugf("跳过非已上架物品，Listing ID: %s", listingID)
				continue
			}
		} else {
			if !strings.Contains(rowHTML, "CancelMarketListingConfirmation") {
				Logger.Debugf("跳过非等待确认物品，Listing ID: %s", listingID)
				continue
			}
		}

		item := Model.MyListingReponse{
			ListingID: listingID,
		}

		// 提取Asset ID
		var assetIDRegex *regexp.Regexp
		if isActiveListing {
			assetIDRegex = regexp.MustCompile(`RemoveMarketListing\('mylisting',\s*'[^']+',\s*\d+,\s*'[^']+',\s*'([^']+)'\)`)
		} else {
			assetIDRegex = regexp.MustCompile(`CancelMarketListingConfirmation\('mylisting',\s*'[^']+',\s*\d+,\s*'[^']+',\s*'([^']+)'\)`)
		}
		if assetIDMatch := assetIDRegex.FindStringSubmatch(rowHTML); len(assetIDMatch) > 1 {
			item.AssetID = assetIDMatch[1]
		}

		// 提取物品名称
		nameRegex := regexp.MustCompile(`href="https://steamcommunity\.com/market/listings/\d+/([^"]+)"`)
		if nameMatch := nameRegex.FindStringSubmatch(rowHTML); len(nameMatch) > 1 {
			if decodedName, err := url.QueryUnescape(nameMatch[1]); err == nil {
				item.MarketHashName = strings.TrimSpace(decodedName)
			} else {
				item.MarketHashName = strings.TrimSpace(nameMatch[1])
			}
		}

		// 提取买家价格
		var buyerPriceRegex *regexp.Regexp
		if isChinese {
			buyerPriceRegex = regexp.MustCompile(`这是买家所要支付[^>]*>\s*([^<]+)\s*<`)
		} else {
			buyerPriceRegex = regexp.MustCompile(`This is the price the buyer pays[^>]*>\s*([^<]+)\s*<`)
		}
		if priceMatch := buyerPriceRegex.FindStringSubmatch(rowHTML); len(priceMatch) > 1 {
			priceStr := strings.ReplaceAll(priceMatch[1], "\\n", "")
			priceStr = strings.ReplaceAll(priceStr, "\\t", "")
			priceStr = strings.ReplaceAll(priceStr, "\n", "")
			priceStr = strings.ReplaceAll(priceStr, "\t", "")
			item.BuyerPrice = parsePrice(strings.TrimSpace(priceStr))
		}

		// 提取卖家到账价格
		sellerPriceRegex := regexp.MustCompile(`\(¥\s*([^)]+)\)`)
		if sellerMatch := sellerPriceRegex.FindStringSubmatch(rowHTML); len(sellerMatch) > 1 {
			item.SellerReceivePrice = parsePrice(sellerMatch[1])
		}

		items = append(items, item)
	}

	return items
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
