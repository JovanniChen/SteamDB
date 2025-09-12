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
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Utils"
)

// GetMyListings 获取用户的上架列表
func (d *Dao) GetMyListings() error {
	Logger.Infof("获取用户 %s 的上架列表", d.GetUsername())

	req, err := d.NewRequest(http.MethodGet, Constants.GetMyListings, nil)
	if err != nil {
		return err
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
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

	Logger.Debug(string(body))

	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	// 输出原始响应(调试用)
	Logger.Debug(buf.String())

	return nil
}

// Remove 删除上架物品
func (d *Dao) RemoveMyListings(creatorId string) error {
	Logger.Infof("删除用户 [%d] 的上架物品，creatorId: %s", d.GetSteamID(), creatorId)

	params := Param.Params{}
	params.SetString("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)

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

	fmt.Println(req.URL.String())
	fmt.Println(req.Header)

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	Logger.Debug(string(body))

	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	return nil
}

func (d *Dao) RemoveAllMyListings() error {
	return nil
}

// BuyListing 购买物品
func (d *Dao) BuyListing(creatorId string, name string, confirmation string) error {
	Logger.Infof("购买用户 [%d] 的上架物品，creatorId: %s", d.GetSteamID(), creatorId)

	params := Param.Params{}
	params.SetString("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	params.SetString("currency", "23")
	params.SetString("subtotal", "10") // 转换为分
	params.SetString("fee", "2")
	params.SetString("total", "12")
	params.SetString("quantity", "1")
	params.SetString("confirmation", confirmation)

	req, err := d.NewRequest(http.MethodPost, Constants.BuyListing+"/"+creatorId, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/market/listings/570/%s", Constants.CommunityOrigin, name))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	Logger.Debug(resp.StatusCode)
	Logger.Debug(resp.Body)

	var reader io.Reader = resp.Body

	Logger.Debug(resp.Header.Get("Content-Encoding"))
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		reader = flate.NewReader(resp.Body)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	Logger.Debug(string(body))

	var buyResp Model.BuyListingResponse
	if err := json.Unmarshal(body, &buyResp); err != nil {
		return fmt.Errorf("解析购买响应失败: %w", err)
	}

	Logger.Debug(buyResp)

	if resp.StatusCode == http.StatusNotAcceptable {
		if buyResp.NeedConfirmation {
			Logger.Debug("购买需要手机令牌确认")
			d.GetConfirmations(creatorId, name, buyResp.Confirmation["confirmation_id"])
		}
	}

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
			if desc.Tradable == 1 && desc.Marketable == 1 {
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
func (d *Dao) PutList(assetID string, price float64, currency int) error {
	Logger.Infof("用户 [%d] 上架物品，AssetID: %s, 价格: %.2f", d.GetSteamID(), assetID, price)

	data := url.Values{}
	data.Set("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	data.Set("appid", "570")   // dota2
	data.Set("contextid", "2") // 分类
	data.Set("assetid", assetID)
	data.Set("amount", "1")
	data.Set("price", strconv.FormatInt(int64(price*100), 10))

	req, err := d.Request(http.MethodPost, Constants.PutList, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", fmt.Sprintf("%s/profiles/%d/inventory", Constants.CommunityOrigin, d.GetSteamID()))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	Logger.Infof("用户 [%d] 上架物品，返回结果: %s，返回状态码: %d", d.GetSteamID(), string(body), resp.StatusCode)

	// 先行处理返回状态码不为200的情况
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("上架失败: %v", string(body))
	}

	var sellResp Model.PutListResponse
	if err := json.Unmarshal(body, &sellResp); err != nil {
		return fmt.Errorf("解析上架响应失败: %w", err)
	}

	// 再行处理返回数据不为成功的情况
	if !sellResp.Success {
		return fmt.Errorf("上架失败: %v", sellResp)
	}

	// 如果需要手机令牌确认
	if sellResp.RequiresConfirmation == 1 && sellResp.NeedsMobileConfirmation {
		Logger.Infof("物品上架需要手机令牌确认，assetID: %s", assetID)
		// 进行手机令牌操作
	} else {
		Logger.Warnf("无法进行确认操作，assetID: %s, RequiresConfirmation: %d", assetID, sellResp.RequiresConfirmation)
	}

	return nil
}

func (d *Dao) GetConfirmations(creatorId string, name string, id string) error {
	username := d.GetUsername()
	Logger.Infof("开始获取用户 [%s] 待确认请求", username)

	pt, err := Utils.LoadMaFile(d.GetUsername() + ".maFile")
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
		Logger.Errorf("待确认API返回失败，用户: [%s], success字段: %t", username, confirmResp.Success)
		return fmt.Errorf("待确认API返回失败")
	}

	Logger.Infof("获取用户 [%s] 的待确认完成，共找到 %d 个待确认请求", username, len(confirmResp.Confirmations))

	for i, conf := range confirmResp.Confirmations {
		Logger.Infof("confirmResp.Confirmations[%d] = %+v", i, conf)
	}

	for _, conf := range confirmResp.Confirmations {
		err = d.AllowSingleConfirmation(pt, conf, steamTime)
		if err != nil {
			Logger.Errorf("允许待确认失败，用户: [%s], 错误: %v", username, err)
			return err
		}

		// err = d.CancelSingleConfirmation(pt, conf, steamTime)
		// if err != nil {
		// 	Logger.Errorf("拒绝待确认失败，用户: [%s], 错误: %v", username, err)
		// 	return err
		// }
	}

	// Logger.Debug("开始最后一次购买操作")
	// err = d.BuyListing(creatorId, name, id)
	// if err != nil {
	// 	Logger.Debug("最后一次购买操作失败")
	// 	Logger.Debug(err)
	// }
	// Logger.Debug("最后一次购买操作成功")

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

	Logger.Debug(resp.StatusCode)
	Logger.Debug(string(body))

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
