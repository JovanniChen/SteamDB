package Dao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Utils"
)

// GetMyListings 获取用户的上架列表
func (d *Dao) GetMyListings() error {
	fmt.Printf("获取用户 %s 的上架列表", d.GetUsername())

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

	fmt.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	// 输出原始响应(调试用)
	fmt.Println(buf.String())

	return nil
}

// GetInventory 获取用户库存
func (d *Dao) GetInventory() error {
	fmt.Printf("获取用户[%d]的库存\n", d.GetSteamID())

	inventoryUrl := fmt.Sprintf("%s/%d/570/2", Constants.GetInventory, d.GetSteamID())
	req, err := d.NewRequest(http.MethodGet, inventoryUrl, nil)
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

	fmt.Println(string(body))

	return nil
}

// PutList 上架物品，需要二次手机令牌确认
func (d *Dao) PutList(assetID string, price float64, currency int) error {
	fmt.Printf("用户 [%d] 上架物品，AssetID: %s, 价格: %.2f\n", d.GetSteamID(), assetID, price)

	data := url.Values{}
	data.Set("sessionid", d.GetLoginCookies()["steamcommunity.com"].SessionId)
	data.Set("appid", "570")
	data.Set("contextid", "2")
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

	var sellResp Model.PutListResponse
	if err := json.Unmarshal(body, &sellResp); err != nil {
		return fmt.Errorf("解析上架响应失败: %w", err)
	}

	if !sellResp.Success {
		return fmt.Errorf("上架失败: %v", sellResp)
	}

	// 如果需要手机令牌确认
	if sellResp.RequiresConfirmation == 1 && sellResp.NeedsMobileConfirmation {
		fmt.Printf("物品上架需要手机令牌确认，assetID: %s\n", assetID)
		// 进行手机令牌操作
	} else {
		fmt.Printf("无法进行确认操作，assetID: %s, RequiresConfirmation: %d\n", assetID, sellResp.RequiresConfirmation)
	}

	return nil
}

func (d *Dao) AcceptConfirmations() {}

func (d *Dao) AcceptSingleConfirmation(phoneToken *Utils.PhoneToken, conf Model.Confirmation, timestamp int64) error {
	fmt.Printf("确认用户 [%d] 上架物品，confID: %s\n", d.GetSteamID(), conf.ID)

	params, err := phoneToken.GenerateConfirmationQueryParams(timestamp, "allow")
	if err != nil {
		return err
	}

	params.SetString("op", "allow")
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

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))

	return nil
}

func (d *Dao) CancelSingleConfirmation(phoneToken *Utils.PhoneToken, conf Model.Confirmation, timestamp int64) error {
	fmt.Printf("取消确认用户 [%d] 上架物品，confID: %s\n", d.GetSteamID(), conf.ID)

	params, err := phoneToken.GenerateConfirmationQueryParams(timestamp, "cancel")
	if err != nil {
		return err
	}

	params.SetString("op", "cancel")
	params.SetString("cid", conf.ID)

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

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))

	return nil
}

func (d *Dao) GetConfirmations() error {
	fmt.Printf("获取用户 [%d] 待确认上架物品\n", d.GetSteamID())

	steamTime, err := d.SteamTime()
	if err != nil {
		return err
	}

	pt, err := Utils.LoadMaFile(d.GetUsername() + ".maFile")
	if err != nil {
		return err
	}

	queryParams, err := Utils.GenerateConfirmationQueryParams(pt.MaFile.DeviceID, pt.MaFile.IdentitySecret, strconv.Itoa(int(pt.MaFile.Session.SteamID)), steamTime, "conf")
	if err != nil {
		return err
	}

	req, err := d.Request(http.MethodGet, Constants.GetConfirmationList+"?"+queryParams.ToUrl(), nil)
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

	fmt.Println(string(body))

	var confirmResp Model.ConfirmationsResponse
	if err := json.Unmarshal(body, &confirmResp); err != nil {
		return fmt.Errorf("解析确认列表失败: %v", err)
	}

	if !confirmResp.Success {
		return fmt.Errorf("获取确认列表失败")
	}

	for _, conf := range confirmResp.Confirmations {
		d.AcceptSingleConfirmation(pt, conf, steamTime)
	}

	return nil
}
