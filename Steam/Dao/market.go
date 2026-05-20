package Dao

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

func (d *Dao) GetMarketListings(gameID int, gameName string, start, count int, country, language string, currency int) (*Model.GetMarketListingIntegrationResponse, error) {
	Logger.Infof("获取市场列表")

	// if start < 0 || count > 100 {
	// 	return nil, errors.New("start must be >= 0 and count must be <= 100")
	// }

	marketUrl := fmt.Sprintf(Constants.GetMarketListing+"%d/%s/render", gameID, gameName)

	params := Param.Params{}
	params.SetString("query", "")
	params.SetString("start", strconv.Itoa(start))
	params.SetString("count", strconv.Itoa(count))
	params.SetString("country", country)
	params.SetString("language", language)
	params.SetString("currency", strconv.Itoa(currency))

	req, err := d.Request(http.MethodGet, marketUrl+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "bMarketOptOut", Value: "1"})

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var response Model.GetMarketListingResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("获取市场列表失败")
	}

	marketListingItems := make([]Model.GetMarketListingItem, 0)
	for _, listing := range response.ListingInfo {
		marketListingItems = append(marketListingItems, Model.GetMarketListingItem{
			AssetID:               listing.AssetInfo.ID,
			ListingID:             listing.ListingID,
			Price:                 listing.Price,
			Fee:                   listing.Fee,
			ConvertedPrice:        listing.ConvertedPrice,
			ConvertedFee:          listing.ConvertedFee,
			ConvertedSteamFee:     listing.ConvertedSteamFee,
			ConvertedPublisherFee: listing.ConvertedPublisherFee,
			ConvertedPricePerUnit: listing.ConvertedPricePerUnit,
		})
	}

	return &Model.GetMarketListingIntegrationResponse{
		Start:      response.Start,
		PageSize:   response.PageSize,
		TotalCount: response.TotalCount,
		Items:      marketListingItems,
	}, nil
}

// GetMyListings 获取用户的上架列表
// 返回两个列表：已上架的物品和等待确认的物品
func (d *Dao) GetMyListings() (Model.GetMyListingResponse, error) {
	Logger.Infof("获取用户 %s 的上架列表", d.GetUsername())

	var response Model.GetMyListingResponse

	params := Param.Params{}
	params.SetString("count", "50")
	params.SetString("norender", "1")

	req, err := d.NewRequest(http.MethodGet, Constants.GetMyListings+"?"+params.ToUrl(), nil)
	if err != nil {
		return response, err
	}

	// 如果有会话信息，添加Cookie
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		req.AddCookie(&http.Cookie{Name: "sessionid", Value: d.GetLoginCookies()["steamcommunity.com"].SessionId})
		req.AddCookie(&http.Cookie{Name: "steamLoginSecure", Value: d.GetLoginCookies()["steamcommunity.com"].SteamLoginSecure})
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return response, Errors.ResponseError(resp.StatusCode)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Error("JSON解析错误:", err)
		return response, err
	}

	Logger.Infof("获取用户[%s]的上架列表: %+v", d.GetUsername(), response)

	// Logger.Infof("[特殊打印	]获取用户[%s]的上架列表的html，数量: %s", d.GetUsername(), response.ResultsHTML)

	// activeListings, pendingListings, err := parseSteamMarketHTML(response.ResultsHTML)
	// if err != nil {
	// 	Logger.Error("从html中解析上架物品失败:", err)
	// 	return nil, err
	// }

	// Logger.Infof("[特殊打印]获取用户[%s]的已上架列表，数量: %d", d.GetUsername(), len(activeListings))
	// for _, listing := range activeListings {
	// 	Logger.Infof("[特殊打印]已上架列表: %+v", listing)
	// }
	// Logger.Infof("[特殊打印]获取用户[%s]的待确认上架列表，数量: %d", d.GetUsername(), len(pendingListings))
	// for _, listing := range pendingListings {
	// 	Logger.Infof("[特殊打印]待确认上架列表: %+v", listing)
	// }

	// for _, listing := range pendingListings {
	// 	Logger.Infof("删除用户 [%s] 的等待确认物品，creatorId: %s", d.GetUsername(), listing.ListingID)
	// 	err := d.RemoveMyListings(listing.ListingID)
	// 	if err != nil {
	// 		Logger.Errorf("删除listing失败 [%s]: %v", listing.ListingID, err)
	// 		// 可以选择继续删除下一个，或者返回错误
	// 		continue
	// 	}
	// 	Logger.Infof("成功删除listing [%s]", listing.ListingID)
	// 	fmt.Printf("%+v\n", listing)
	// }

	return response, nil
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

func (d *Dao) buy(gameId int, currency string, creatorId string, name string, buyerPrice, sellerReceivePrice int, confirmation string) buyResult {
	Logger.Infof("[%s]购买[%s][%s][%d][%d][%s]", d.GetUsername(), creatorId, name, buyerPrice, sellerReceivePrice, confirmation)

	fee := buyerPrice - sellerReceivePrice

	params := Param.Params{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		params.SetString("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}
	params.SetString("currency", currency)
	params.SetInt64("subtotal", int64(sellerReceivePrice))
	params.SetInt64("fee", int64(fee))
	params.SetInt64("total", int64(buyerPrice))
	params.SetString("quantity", "1")
	params.SetString("billing_state", "")
	params.SetInt64("tradefee_tax", 0)
	params.SetInt64("save_my_address", 0)
	params.SetString("confirmation", confirmation)

	fmt.Println(params.Encode())

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
	// https: //steamcommunity.com/market/listings/440/Specialized%20Killstreak%20Maul%20Kit%20Fabricator
	req.Header.Set("referer", fmt.Sprintf("%s/market/listings/%d/%s", Constants.CommunityOrigin, gameId, name))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            Errors.ErrRetryRequest,
		}
	}
	if resp == nil || resp.Body == nil {
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("购买失败: 空响应"),
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

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		Logger.Warnf("用户 [%s] 购买物品遇到速率限制 (429)", d.GetUsername())
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("购买失败: %w", Errors.ErrRateLimited),
		}
	case http.StatusBadGateway:
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
				error:            fmt.Errorf("购买失败: %w", Errors.ErrAccountBan),
			}
		}
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("%w", Errors.ErrServerError),
		}
	case http.StatusNotAcceptable:
		var buyListingResp Model.BuyListingNeedConfirmationResponse
		if err := json.Unmarshal(body, &buyListingResp); err != nil {
			return buyResult{
				success:          false,
				needConfirmation: false,
				confirmationId:   "",
				error:            err,
			}
		}
		confirmationId := ""
		if buyListingResp.Confirmation != nil {
			confirmationId = buyListingResp.Confirmation["confirmation_id"]
		}
		if buyListingResp.NeedConfirmation && confirmationId == "" {
			return buyResult{
				success:          false,
				needConfirmation: true,
				confirmationId:   "",
				error:            fmt.Errorf("购买失败: 缺少 confirmation_id"),
			}
		}
		return buyResult{
			success:          buyListingResp.Success == 22,
			needConfirmation: buyListingResp.NeedConfirmation,
			confirmationId:   confirmationId,
			error:            nil,
		}
	case http.StatusOK:
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
		}
		Logger.Infof("用户[%s]购买物品[creatorId: %s][confirmation: '%s']失败", d.GetUsername(), creatorId, confirmation)
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("购买失败，错误码: %d", buyListingResp.WalletInfo.Success),
		}
	case http.StatusBadRequest:
		return buyResult{
			success:          false,
			needConfirmation: false,
			confirmationId:   "",
			error:            fmt.Errorf("购买失败，错误码: %d，错误信息: %s", resp.StatusCode, "返回为空"),
		}
	default:
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
func (d *Dao) BuyListing(gameId int, currency string, creatorId string, name string, buyerPrice, sellerReceivePrice int, confirmation string, maFileContent string) error {
	br := d.buy(gameId, currency, creatorId, name, buyerPrice, sellerReceivePrice, confirmation)
	if br.success && br.needConfirmation {
		if err := d.ConfirmationForBuyList("allow", maFileContent); err != nil {
			return err
		}
		brAgain := d.buy(gameId, currency, creatorId, name, buyerPrice, sellerReceivePrice, br.confirmationId)
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
func (d *Dao) GetSteamGift(gameId int, categoryId int) ([]Model.Item, error) {

	inventoryUrl := fmt.Sprintf("%s/%d/%d/%d", Constants.GetInventory, d.GetSteamID(), gameId, categoryId)
	req, err := d.Request(http.MethodGet, inventoryUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建库存请求失败: %w", err)
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, fmt.Errorf("执行库存请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取库存响应失败: %w", err)
	}

	// 检查是否为GZIP压缩数据
	if len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b {
		// 解压GZIP数据
		reader, _ := gzip.NewReader(bytes.NewReader(body))
		defer reader.Close()
		body, _ = io.ReadAll(reader)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取库存失败,返回状态码: %d", resp.StatusCode)
	}

	var inventoryResponse Model.InventoryResponse
	if err := json.Unmarshal(body, &inventoryResponse); err != nil {
		return nil, fmt.Errorf("解析库存响应失败: %w", err)
	}

	if inventoryResponse.Success != 1 {
		return nil, fmt.Errorf("库存API返回失败，success=%d", inventoryResponse.Success)
	}

	var steamGiftResponse []Model.Item
	for _, asset := range inventoryResponse.Assets {
		steamGiftResponse = append(steamGiftResponse, Model.Item{
			AssetID: asset.AssetID,
		})
	}

	return steamGiftResponse, nil
}

func (d *Dao) IsAccountBanned() bool {
	req, err := d.Request(http.MethodGet, Constants.Market, nil)

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// 您的帐户目前已被锁定
	// Your account is currently locked

	if strings.Contains(string(body), "您的帐户目前已被锁定") || strings.Contains(string(body), "Your account is currently locked") {
		return true
	}
	return false
}

type GetInventoryOption func(*inventoryOptions)

type inventoryOptions struct {
	tradable   *int
	marketable *int
	commodity  *int
}

func WithTradable(tradable int) GetInventoryOption {
	return func(o *inventoryOptions) {
		o.tradable = &tradable
	}
}

func WithMarketable(marketable int) GetInventoryOption {
	return func(o *inventoryOptions) {
		o.marketable = &marketable
	}
}

func WithCommodity(commodity int) GetInventoryOption {
	return func(o *inventoryOptions) {
		o.commodity = &commodity
	}
}

// GetInventory 获取用户库存
func (d *Dao) GetInventory(gameId, categoryId int, opts ...GetInventoryOption) ([]Model.Item, error) {
	username := d.GetUsername()
	Logger.Infof("开始获取用户 [%s] 的库存，游戏ID: %d, 分类ID: %d", username, gameId, categoryId)

	inventoryUrl := fmt.Sprintf("%s/%d/%d/%d", Constants.GetInventory, d.GetSteamID(), gameId, categoryId)
	req, err := d.Request(http.MethodGet, inventoryUrl, nil)
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

	switch resp.StatusCode {
	case 429:
		Logger.Warnf("用户 [%s] 获取库存遇到速率限制 (429)", username)
		Logger.Warnf("响应内容: %s", string(body))
		return nil, fmt.Errorf("获取库存失败: %w", Errors.ErrRateLimited)
	case 401, 403:
		Logger.Warnf("用户 [%s] 获取库存遇到授权失败 (401/403)", username)
		Logger.Warnf("响应内容: %s", string(body))
		return nil, fmt.Errorf("获取库存失败: %w", Errors.ErrAuthorizationFailed)
	}

	var inventoryResponse Model.InventoryResponse
	if err := json.Unmarshal(body, &inventoryResponse); err != nil {
		Logger.Errorf("解析库存响应失败，用户: [%s], 错误: %v", username, err)
		Logger.Warnf("响应内容: %s", string(body))
		return nil, fmt.Errorf("解析库存响应失败: %w", err)
	}

	if inventoryResponse.Success != 1 {
		Logger.Errorf("库存API返回失败，用户: [%s], success字段: %d", username, inventoryResponse.Success)
		return nil, fmt.Errorf("库存API返回失败，success=%d", inventoryResponse.Success)
	}

	// 转换为内部物品模型
	items := d.processInventoryData(&inventoryResponse, username, opts...)

	Logger.Infof("获取用户 [%s] 的库存完成，共找到 %d 个可交易物品", username, len(items))

	return items, nil
}

// processInventoryData 处理库存数据并返回可交易物品列表
func (d *Dao) processInventoryData(inventoryResponse *Model.InventoryResponse, username string, opts ...GetInventoryOption) []Model.Item {
	// 边界检查
	if inventoryResponse == nil {
		Logger.Warnf("用户 [%s] 库存响应为空", username)
		return []Model.Item{}
	}

	if len(inventoryResponse.Assets) == 0 {
		Logger.Infof("用户 [%s] 的库存为空", username)
		return []Model.Item{}
	}

	options := &inventoryOptions{}
	for _, opt := range opts {
		opt(options)
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
			if !exists {
				missingDescCount++
				continue
			}

			if options.tradable != nil && desc.Tradable != *options.tradable {
				continue
			}
			if options.marketable != nil && desc.Marketable != *options.marketable {
				continue
			}
			if options.commodity != nil && desc.Commodity != *options.commodity {
				continue
			}

			items = append(items, Model.Item{
				AssetID:    asset.AssetID,
				ClassID:    asset.ClassID,
				InstanceID: asset.InstanceID,
				Name:       desc.Name,
				MarketName: desc.MarketName,
				Tradable:   desc.Tradable == 1,
				Marketable: desc.Marketable == 1,
				Commodity:  desc.Commodity == 1,
			})
			filteredCount++
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
func (d *Dao) PutList(gameId int, contextId int, assetID string, price int, currency int, maFileContent string) (Model.MyListingReponse, error) {
	Logger.Infof("用户 [%s] 上架物品，AssetID: %s, 价格: %d", d.GetUsername(), assetID, price)

	data := url.Values{}
	if d.GetLoginCookies()["steamcommunity.com"] != nil {
		data.Set("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	}
	data.Set("appid", strconv.Itoa(gameId))
	data.Set("contextid", strconv.Itoa(contextId)) // 分类
	data.Set("assetid", assetID)
	data.Set("amount", "1")
	data.Set("price", strconv.Itoa(price))

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
		Logger.Infof("确认结果: %+v", result)
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
		Result: Model.MyListingReponse{
			ListingID: confirmResp.Confirmations[0].CreatorID,
		},
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

func (d *Dao) ConfirmationForSendGift(op string, maFileContent string) *Model.ConfirmationResult {
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
		Result: Model.MyListingReponse{
			ListingID: confirmResp.Confirmations[0].CreatorID,
		},
	}

	for i := len(confirmResp.Confirmations) - 1; i >= 0; i-- {
		Logger.Infof("处理第 %d 个确认项", i+1)
		conf := confirmResp.Confirmations[i]
		if conf.Type != 2 {
			Logger.Infof("非赠送饰品确认不予处理:%+v", conf)
			continue
		}

		if i != 0 {
			Logger.Infof("【允许】其他赠送饰品确认")
			sTime, _ := d.GetSteamTimeLocal()
			err := d.AllowSingleConfirmation(pt, conf, sTime)
			if err != nil {
				Logger.Errorf("【允许】其他赠送饰品确认失败，用户: [%s], 错误: %v", username, err)
			}
			Logger.Errorf("【允许】其他赠送饰品确认成功，用户: [%s]", username)
		} else {
			Logger.Infof("【允许】本次赠送饰品确认")
			for j := 0; j < Constants.Tries; j++ {
				sTime, _ := d.GetSteamTimeLocal()
				err = d.AllowSingleConfirmation(pt, conf, sTime)
				if err != nil {
					Logger.Errorf("第 %d 次【允许】待确认失败，用户: [%s], 错误: %v", j+1, username, err)
					time.Sleep(100 * time.Millisecond)
					continue
				} else {
					Logger.Infof("【允许】成功本次赠送饰品确认:%+v", conf)
					finalResult.Success = true
					return finalResult
				}

			}
		}
	}

	return finalResult
}

func (d *Dao) GetConfirmations(maFileContent string) error {
	Logger.Infof("开始获取用户 [%s] 的待确认请求", d.GetUsername())

	var response Model.ConfirmationsResponse

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

	// req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 9; Valve Steam App Version/3)")
	// req.Header.Set("mobileClient", "android")
	// req.Header.Set("mobileClientVersion", "777777 3.6.4")

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

	if err := json.Unmarshal(body, &response); err != nil {
		Logger.Errorf("解析待确认响应失败，用户: [%s], 错误: %v", username, err)
		return fmt.Errorf("解析待确认响应失败: %v", err)
	}

	if !response.Success {
		Logger.Errorf("待确认API返回失败，用户: [%s], success字段: %t, 返回码：%d", username, response.Success, resp.StatusCode)
		return fmt.Errorf("待确认API返回失败")
	}

	for _, conf := range response.Confirmations {
		Logger.Infof("待确认列表: %+v", conf)
	}

	for i := len(response.Confirmations) - 1; i >= 0; i-- {
		Logger.Infof("正在【允许】第 %d 个待确认", i+1)
		conf := response.Confirmations[i]
		err := d.AllowSingleConfirmation(pt, conf, steamTime)
		if err != nil {
			Logger.Errorf("【允许】待确认失败，用户: [%s], 错误: %v", username, err)
		}
		Logger.Errorf("【允许】待确认成功，用户: [%s]", username)
	}

	// for i := len(response.Confirmations) - 1; i >= 0; i-- {
	// 	Logger.Infof("正在【拒绝】第 %d 个待确认", i+1)
	// 	conf := response.Confirmations[i]
	// 	err := d.CancelSingleConfirmation(pt, conf, steamTime)
	// 	if err != nil {
	// 		Logger.Errorf("【拒绝】待确认失败，用户: [%s], 错误: %v", username, err)
	// 	}
	// 	Logger.Errorf("【拒绝】待确认成功，用户: [%s]", username)
	// }

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

func (d *Dao) GetPartnerInventory(partnerUrl string, gameId, contextId int) ([]Model.PartnerIntegrationItem, error) {
	partnerIntegrations := make([]Model.PartnerIntegrationItem, 0)

	u, err := url.Parse(partnerUrl)
	if err != nil {
		return nil, err
	}
	partner := u.Query().Get("partner")
	partnerID, err := strconv.ParseUint(partner, 10, 64)
	if err != nil {
		return nil, err
	}
	steamId := Utils.FriendCodeToSteamID64(uint32(partnerID))

	cookies, ok := d.GetLoginCookies()["steamcommunity.com"]
	if !ok {
		return nil, errors.New("sessionid not found")
	}
	sessionid := cookies.SessionId

	params := Param.Params{}
	params.SetString("sessionid", sessionid)
	params.SetString("partner", strconv.Itoa(int(steamId)))
	params.SetString("appid", strconv.Itoa(gameId))
	params.SetString("contextid", strconv.Itoa(contextId))

	fmt.Println(params.ToUrl())

	req, err := d.Request(http.MethodGet, Constants.GetPartnerInventory+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("referer", partnerUrl)

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解码 JSON
	var partnerInventoryResp Model.PartnerInventoryResponse
	err = json.Unmarshal(body, &partnerInventoryResp)
	if err != nil {
		log.Fatal("JSON 解析失败:", err)
	}

	// 输出基本信息
	fmt.Printf("成功: %v\n", partnerInventoryResp.Success)
	fmt.Printf("游戏: %s\n", partnerInventoryResp.RGAppInfo.Name)
	fmt.Println(len(partnerInventoryResp.RGInventory))
	fmt.Println(len(partnerInventoryResp.RGDescriptions))

	for _, item := range partnerInventoryResp.RGInventory {
		partnerIntegrations = append(partnerIntegrations, Model.PartnerIntegrationItem{
			ID:             item.ID,
			MarketName:     partnerInventoryResp.RGDescriptions[item.ClassID+"_"+item.InstanceID].MarketName,
			MarketHashName: partnerInventoryResp.RGDescriptions[item.ClassID+"_"+item.InstanceID].MarketHashName,
		})
	}

	return partnerIntegrations, nil
}

// https://steamcommunity.com/tradeoffer/new/?partner=352956450&token=U4SAf1wu

// sessionid 33f2afd25ecbf44da6443dfc
// serverid 1
// partner 76561199668111414
// tradeoffermessage
// json_tradeoffer
// {"newversion":true,"version":2,"me":{"assets":[{"appid":440,"contextid":"2","amount":1,"assetid":"16317805093"}],"currency":[],"ready":false},"them":{"assets":[],"currency":[],"ready":false}}
// captcha
// trade_offer_create_params
// {"trade_offer_access_token":"fGsR2IfZ"}
func (d *Dao) SendGift(partnerUrl, assetId, maFileContent string) error {
	if d.GetLoginCookies()["steamcommunity.com"] == nil {
		return errors.New("sessionid not found")
	}

	sessionid := d.GetLoginCookies()["steamcommunity.com"].SessionId

	parsedURL, err := url.Parse(partnerUrl)
	if err != nil {
		return err
	}
	partner := parsedURL.Query().Get("partner")
	partnerID, err := strconv.ParseUint(partner, 10, 64)
	if err != nil {
		return err
	}
	steamId := Utils.FriendCodeToSteamID64(uint32(partnerID))

	token := parsedURL.Query().Get("token")

	jsonTradeoffer := fmt.Sprintf("{\"newversion\":true,\"version\":2,\"me\":{\"assets\":[{\"appid\":440,\"contextid\":\"2\",\"amount\":1,\"assetid\":\"%s\"}],\"currency\":[],\"ready\":false},\"them\":{\"assets\":[],\"currency\":[],\"ready\":false}}", assetId)

	tradeOfferAccessToken := fmt.Sprintf("{\"trade_offer_access_token\":\"%s\"}", token)

	params := Param.Params{}
	params.SetString("sessionid", sessionid)
	params.SetString("serverid", "1")
	params.SetString("partner", strconv.Itoa(int(steamId)))
	params.SetString("tradeoffermessage", "")
	params.SetString("json_tradeoffer", jsonTradeoffer)
	params.SetString("captcha", "")
	params.SetString("trade_offer_create_params", tradeOfferAccessToken)

	fmt.Println(params.Encode())

	req, err := d.Request(http.MethodPost, Constants.SendTradeOffer, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("origin", "https://steamcommunity.com")
	req.Header.Set("referer", partnerUrl)

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("赠送饰品失败,返回状态码: %d", resp.StatusCode)
	}

	fmt.Println(string(body))

	var sendGiftResp Model.SendGiftResponse
	if err := json.Unmarshal(body, &sendGiftResp); err != nil {
		return fmt.Errorf("解析赠送饰品响应失败: %w", err)
	}

	if sendGiftResp.NeedsMobileConfirmation {
		result := d.ConfirmationForSendGift("allow", maFileContent)
		if !result.Success {
			return fmt.Errorf("赠送饰品失败")
		}
	}

	return nil
}
