package Dao

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	"golang.org/x/net/html"
)

type InitTransactionParams struct {
	GidShoppingCart          int64  `form:"gidShoppingCart"`
	GidReplayOfTransID       int64  `form:"gidReplayOfTransID"`
	UseAccountCart           int64  `form:"bUseAccountCart"`
	PaymentMethod            string `form:"PaymentMethod"`
	AbortPendingTransactions int64  `form:"abortPendingTransactions"`
	HasCardInfo              int64  `form:"bHasCardInfo"`
	CardNumber               string `form:"CardNumber"`
	CardExpirationYear       string `form:"CardExpirationYear"`
	CardExpirationMonth      string `form:"CardExpirationMonth"`
	FirstName                string `form:"FirstName"`
	LastName                 string `form:"LastName"`
	Address                  string `form:"Address"`
	AddressTwo               string `form:"AddressTwo"`
	Country                  string `form:"Country"`
	City                     string `form:"City"`
	State                    string `form:"State"`
	PostalCode               string `form:"PostalCode"`
	Phone                    string `form:"Phone"`
	ShippingFirstName        string `form:"ShippingFirstName"`
	ShippingLastName         string `form:"ShippingLastName"`
	ShippingAddress          string `form:"ShippingAddress"`
	ShippingAddressTwo       string `form:"ShippingAddressTwo"`
	ShippingCountry          string `form:"ShippingCountry"`
	ShippingCity             string `form:"ShippingCity"`
	ShippingState            string `form:"ShippingState"`
	ShippingPostalCode       string `form:"ShippingPostalCode"`
	ShippingPhone            string `form:"ShippingPhone"`
	IsGift                   int64  `form:"bIsGift"`
	GifteeAccountID          int64  `form:"GifteeAccountID"`
	GifteeEmail              string `form:"GifteeEmail"`
	GifteeName               string `form:"GifteeName"`
	GiftMessage              string `form:"GiftMessage"`
	Sentiment                string `form:"Sentiment"`
	Signature                string `form:"Signature"`
	ScheduledSendOnDate      int64  `form:"ScheduledSendOnDate"`
	BankAccount              string `form:"BankAccount"`
	BankCode                 string `form:"BankCode"`
	BankIBAN                 string `form:"BankIBAN"`
	BankBIC                  string `form:"BankBIC"`
	TPBankID                 string `form:"TPBankID"`
	BankAccountID            string `form:"BankAccountID"`
	SaveBillingAddress       int64  `form:"bSaveBillingAddress"`
	GidPaymentID             string `form:"gidPaymentID"`
	UseRemainingSteamAccount int64  `form:"bUseRemainingSteamAccount"`
	PreAuthOnly              int64  `form:"bPreAuthOnly"`
	SessionID                string `form:"sessionid"`
}

//	{
//	    "success": 2,
//	    "purchaseresultdetail": 53,
//	    "paymentmethod": 11,
//	    "transid": "18446744073709551615",
//	    "transactionprovider": 0,
//	    "paymentmethodcountrycode": "",
//	    "paypaltoken": "",
//	    "paypalacct": 0,
//	    "packagewitherror": -1,
//	    "appcausingerror": 0,
//	    "pendingpurchasepaymentmethod": 0,
//	    "authorizationurl": ""
//	}
type InitTransactionResponse struct {
	Success                      int    `json:"success"`
	PurchaseResultDetail         int    `json:"purchaseresultdetail"`
	PaymentMethod                int    `json:"paymentmethod"`
	TransID                      string `json:"transid"`
	TransactionProvider          int    `json:"transactionprovider"`
	PaymentMethodCountryCode     string `json:"paymentmethodcountrycode"`
	PaypalToken                  string `json:"paypaltoken"`
	PaypalAcct                   int    `json:"paypalacct"`
	PackageWithError             int    `json:"packagewitherror"`
	AppCausingError              int    `json:"appcausingerror"`
	PendingPurchasePaymentMethod int    `json:"pendingpurchasepaymentmethod"`
	AuthorizationURL             string `json:"authorizationurl"`
}

// gidShoppingCart: -1
// gidReplayOfTransID: -1
// bUseAccountCart: 1
// PaymentMethod: alipay
// abortPendingTransactions: 0
// bHasCardInfo: 0
// CardNumber:
// CardExpirationYear:
// CardExpirationMonth:
// FirstName:
// LastName:
// Address:
// AddressTwo:
// Country: CN
// City:
// State:
// PostalCode:
// Phone:
// ShippingFirstName:
// ShippingLastName:
// ShippingAddress:
// ShippingAddressTwo:
// ShippingCountry: CN
// ShippingCity:
// ShippingState:
// ShippingPostalCode:
// ShippingPhone:
// bIsGift: 0
// GifteeAccountID: 0
// GifteeEmail:
// GifteeName:
// GiftMessage:
// Sentiment:
// Signature:
// ScheduledSendOnDate: 0
// BankAccount:
// BankCode:
// BankIBAN:
// BankBIC:
// TPBankID:
// BankAccountID:
// bSaveBillingAddress: 1
// gidPaymentID:
// bUseRemainingSteamAccount: 0
// bPreAuthOnly: 0
// sessionid: 5bf319d1458a73814a5873be

func (d *Dao) InitTransaction() (string, error) {
	params := Param.Params{}
	params.SetInt64("gidShoppingCart", -1)
	params.SetInt64("gidReplayOfTransID", -1)
	params.SetInt64("bUseAccountCart", 1)
	params.SetString("PaymentMethod", "alipay")
	params.SetInt64("abortPendingTransactions", 0)
	params.SetInt64("bHasCardInfo", 0)
	params.SetString("CardNumber", "")
	params.SetString("CardExpirationYear", "")
	params.SetString("CardExpirationMonth", "")
	params.SetString("FirstName", "")
	params.SetString("LastName", "")
	params.SetString("Address", "")
	params.SetString("AddressTwo", "")
	params.SetString("Country", "CN")
	params.SetString("City", "")
	params.SetString("State", "")
	params.SetString("PostalCode", "")
	params.SetString("Phone", "")
	params.SetString("ShippingFirstName", "")
	params.SetString("ShippingLastName", "")
	params.SetString("ShippingAddress", "")
	params.SetString("ShippingAddressTwo", "")
	params.SetString("ShippingCountry", "CN")
	params.SetString("ShippingCity", "")
	params.SetString("ShippingState", "")
	params.SetString("ShippingPostalCode", "")
	params.SetString("ShippingPhone", "")
	params.SetInt64("bIsGift", 0)
	params.SetInt64("GifteeAccountID", 0)
	params.SetString("GifteeEmail", "")
	params.SetString("GifteeName", "")
	params.SetString("GiftMessage", "")
	params.SetString("Sentiment", "")
	params.SetString("Signature", "")
	params.SetInt64("ScheduledSendOnDate", 0)
	params.SetString("BankAccount", "")
	params.SetString("BankCode", "")
	params.SetString("BankIBAN", "")
	params.SetString("BankBIC", "")
	params.SetString("TPBankID", "")
	params.SetString("BankAccountID", "")
	params.SetInt64("bSaveBillingAddress", 1)
	params.SetString("gidPaymentID", "")
	params.SetInt64("bUseRemainingSteamAccount", 1)
	params.SetInt64("bPreAuthOnly", 0)

	if d.GetLoginCookies()["checkout.steampowered.com"] == nil {
		return "", errors.New("checkout.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["checkout.steampowered.com"].SessionId
	params.SetString("sessionid", sessionId)

	req, err := d.Request(http.MethodPost, Constants.InitTransaction, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("origin", "https://checkout.steampowered.com")
	req.Header.Set("referer", "https://checkout.steampowered.com/checkout/?accountcart=1")

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return "", fmt.Errorf("初始化交易失败,返回状态码: %d", resp.StatusCode)
	}

	var response InitTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析初始化交易响应失败: %w", err)
	}

	if response.Success != 1 {
		fmt.Println(string(body))
		return "", errors.New(Errors.GetCheckoutError(response.PurchaseResultDetail))
	}

	return response.TransID, nil
}

func (d *Dao) InitAddFundsTransaction(gidShoppingCart, microTxn string) (string, error) {
	params := Param.Params{}
	params.SetString("gidShoppingCart", gidShoppingCart)
	params.SetInt64("gidReplayOfTransID", -1)
	params.SetInt64("bUseAccountCart", 0)
	params.SetString("PaymentMethod", "alipay")
	params.SetInt64("abortPendingTransactions", 0)
	params.SetInt64("bHasCardInfo", 0)
	params.SetString("CardNumber", "")
	params.SetString("CardExpirationYear", "")
	params.SetString("CardExpirationMonth", "")
	params.SetString("FirstName", "")
	params.SetString("LastName", "")
	params.SetString("Address", "")
	params.SetString("AddressTwo", "")
	params.SetString("Country", "CN")
	params.SetString("City", "")
	params.SetString("State", "")
	params.SetString("PostalCode", "")
	params.SetString("Phone", "")
	params.SetString("ShippingFirstName", "")
	params.SetString("ShippingLastName", "")
	params.SetString("ShippingAddress", "")
	params.SetString("ShippingAddressTwo", "")
	params.SetString("ShippingCountry", "CN")
	params.SetString("ShippingCity", "")
	params.SetString("ShippingState", "")
	params.SetString("ShippingPostalCode", "")
	params.SetString("ShippingPhone", "")
	params.SetInt64("bIsGift", 0)
	params.SetInt64("GifteeAccountID", 0)
	params.SetString("GifteeEmail", "")
	params.SetString("GifteeName", "")
	params.SetString("GiftMessage", "")
	params.SetString("Sentiment", "")
	params.SetString("Signature", "")
	params.SetInt64("ScheduledSendOnDate", 0)
	params.SetString("BankAccount", "")
	params.SetString("BankCode", "")
	params.SetString("BankIBAN", "")
	params.SetString("BankBIC", "")
	params.SetString("TPBankID", "")
	params.SetString("BankAccountID", "")
	params.SetInt64("bSaveBillingAddress", 1)
	params.SetString("gidPaymentID", "")
	params.SetInt64("bUseRemainingSteamAccount", 0)
	params.SetInt64("bPreAuthOnly", 0)

	if d.GetLoginCookies()["checkout.steampowered.com"] == nil {
		return "", errors.New("checkout.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["checkout.steampowered.com"].SessionId
	params.SetString("sessionid", sessionId)

	req, err := d.Request(http.MethodPost, Constants.InitTransaction, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("origin", "https://checkout.steampowered.com")
	req.Header.Set("referer", fmt.Sprintf("https: //checkout.steampowered.com/checkout?cart=%s&microtxn=%s", gidShoppingCart, microTxn))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return "", fmt.Errorf("初始化交易失败,返回状态码: %d", resp.StatusCode)
	}

	var response InitTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析初始化交易响应失败: %w", err)
	}

	fmt.Println(response)

	if response.Success != 1 {
		fmt.Println(string(body))
		return "", errors.New(Errors.GetCheckoutError(response.PurchaseResultDetail))
	}

	return response.TransID, nil
}

func (d *Dao) InitAddFundsTransactionWithCountry(gidShoppingCart, microTxn, country string) (string, error) {
	params := Param.Params{}
	params.SetString("gidShoppingCart", gidShoppingCart)
	params.SetInt64("gidReplayOfTransID", -1)
	params.SetInt64("bUseAccountCart", 0)
	params.SetString("PaymentMethod", "alipay")
	params.SetInt64("abortPendingTransactions", 0)
	params.SetInt64("bHasCardInfo", 0)
	params.SetString("CardNumber", "")
	params.SetString("CardExpirationYear", "")
	params.SetString("CardExpirationMonth", "")
	params.SetString("FirstName", "")
	params.SetString("LastName", "")
	params.SetString("Address", "")
	params.SetString("AddressTwo", "")
	params.SetString("Country", country)
	params.SetString("City", "")
	params.SetString("State", "")
	params.SetString("PostalCode", "")
	params.SetString("Phone", "")
	params.SetString("ShippingFirstName", "")
	params.SetString("ShippingLastName", "")
	params.SetString("ShippingAddress", "")
	params.SetString("ShippingAddressTwo", "")
	params.SetString("ShippingCountry", country)
	params.SetString("ShippingCity", "")
	params.SetString("ShippingState", "")
	params.SetString("ShippingPostalCode", "")
	params.SetString("ShippingPhone", "")
	params.SetInt64("bIsGift", 0)
	params.SetInt64("GifteeAccountID", 0)
	params.SetString("GifteeEmail", "")
	params.SetString("GifteeName", "")
	params.SetString("GiftMessage", "")
	params.SetString("Sentiment", "")
	params.SetString("Signature", "")
	params.SetInt64("ScheduledSendOnDate", 0)
	params.SetString("BankAccount", "")
	params.SetString("BankCode", "")
	params.SetString("BankIBAN", "")
	params.SetString("BankBIC", "")
	params.SetString("TPBankID", "")
	params.SetString("BankAccountID", "")
	params.SetInt64("bSaveBillingAddress", 1)
	params.SetString("gidPaymentID", "")
	params.SetInt64("bUseRemainingSteamAccount", 0)
	params.SetInt64("bPreAuthOnly", 0)

	if d.GetLoginCookies()["checkout.steampowered.com"] == nil {
		return "", errors.New("checkout.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["checkout.steampowered.com"].SessionId
	params.SetString("sessionid", sessionId)

	req, err := d.Request(http.MethodPost, Constants.InitTransaction, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("origin", "https://checkout.steampowered.com")
	req.Header.Set("referer", fmt.Sprintf("https: //checkout.steampowered.com/checkout?cart=%s&microtxn=%s", gidShoppingCart, microTxn))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return "", fmt.Errorf("初始化交易失败,返回状态码: %d", resp.StatusCode)
	}

	var response InitTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析初始化交易响应失败: %w", err)
	}

	fmt.Println(response)

	if response.Success != 1 {
		fmt.Println(string(body))
		return "", errors.New(Errors.GetCheckoutError(response.PurchaseResultDetail))
	}

	return response.TransID, nil
}

func (d *Dao) InitConcurrentTransaction() (string, error) {
	params := Param.Params{}
	params.SetInt64("gidShoppingCart", -1)
	params.SetInt64("gidReplayOfTransID", -1)
	params.SetInt64("bUseAccountCart", 1)
	params.SetString("PaymentMethod", "steamaccount")
	params.SetInt64("abortPendingTransactions", 0)
	params.SetInt64("bHasCardInfo", 0)
	params.SetString("CardNumber", "")
	params.SetString("CardExpirationYear", "")
	params.SetString("CardExpirationMonth", "")
	params.SetString("FirstName", "")
	params.SetString("LastName", "")
	params.SetString("Address", "")
	params.SetString("AddressTwo", "")
	params.SetString("Country", "CN")
	params.SetString("City", "")
	params.SetString("State", "")
	params.SetString("PostalCode", "")
	params.SetString("Phone", "")
	params.SetString("ShippingFirstName", "")
	params.SetString("ShippingLastName", "")
	params.SetString("ShippingAddress", "")
	params.SetString("ShippingAddressTwo", "")
	params.SetString("ShippingCountry", "CN")
	params.SetString("ShippingCity", "")
	params.SetString("ShippingState", "")
	params.SetString("ShippingPostalCode", "")
	params.SetString("ShippingPhone", "")
	params.SetInt64("bIsGift", 0)
	params.SetInt64("GifteeAccountID", 0)
	params.SetString("GifteeEmail", "")
	params.SetString("GifteeName", "")
	params.SetString("GiftMessage", "")
	params.SetString("Sentiment", "")
	params.SetString("Signature", "")
	params.SetInt64("ScheduledSendOnDate", 0)
	params.SetString("BankAccount", "")
	params.SetString("BankCode", "")
	params.SetString("BankIBAN", "")
	params.SetString("BankBIC", "")
	params.SetString("TPBankID", "")
	params.SetString("BankAccountID", "")
	params.SetInt64("bSaveBillingAddress", 1)
	params.SetString("gidPaymentID", "")
	params.SetInt64("bUseRemainingSteamAccount", 1)
	params.SetInt64("bPreAuthOnly", 0)

	if d.GetLoginCookies()["checkout.steampowered.com"] == nil {
		return "", errors.New("checkout.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["checkout.steampowered.com"].SessionId
	params.SetString("sessionid", sessionId)

	req, err := d.Request(http.MethodPost, Constants.InitTransaction, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("origin", "https://checkout.steampowered.com")
	req.Header.Set("referer", "https://checkout.steampowered.com/checkout/?accountcart=1")

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return "", fmt.Errorf("初始化交易失败,返回状态码: %d", resp.StatusCode)
	}

	var response InitTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析初始化交易响应失败: %w", err)
	}

	if response.Success != 1 {
		fmt.Println(string(body))
		return "", errors.New(Errors.GetCheckoutError(response.PurchaseResultDetail))
	}

	return response.TransID, nil
}

func (d *Dao) FinalizeTransaction(transactionID string) error {
	params := Param.Params{}
	params.SetString("transid", transactionID)
	params.SetString("CardCVV2", "")
	params.SetString("browserInfo", `{"language":"zh-CN","javaEnabled":"false","colorDepth":24,"screenHeight":1890,"screenWidth":3360}"`)
	req, err := d.Request(http.MethodPost, Constants.FinalizeTransaction, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("origin", "https://checkout.steampowered.com")
	req.Header.Set("referer", "https://checkout.steampowered.com/checkout/?accountcart=1")

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

	if resp.StatusCode != 200 {
		return fmt.Errorf("完成交易失败,返回状态码: %d", resp.StatusCode)
	}

	// {"success":22,"purchaseresultdetail":0,"bShowBRSpecificCreditCardError":false}
	var response struct {
		Success                       int  `json:"success"`
		PurchaseResultDetail          int  `json:"purchaseresultdetail"`
		ShowBRSpecificCreditCardError bool `json:"bShowBRSpecificCreditCardError"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("解析完成交易响应失败: %w", err)
	}

	if response.Success != 22 {
		return fmt.Errorf("完成交易失败,返回Success: %d, PurchaseResultDetail: %d, ShowBRSpecificCreditCardError: %t", response.Success, response.PurchaseResultDetail, response.ShowBRSpecificCreditCardError)
	}

	return nil
}

type CancelTransactionResponse struct {
	Success int `json:"success"`
}

func (d *Dao) CancelTransaction(transactionID string) error {
	params := Param.Params{}
	params.SetString("transid", transactionID)
	params.SetString("lock_transaction", "false")
	params.SetString("cancel_reason", "老子高兴，你管得着吗？")
	req, err := d.Request(http.MethodPost, Constants.CancelCartTrans, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("origin", "https://checkout.steampowered.com")
	req.Header.Set("referer", "https://checkout.steampowered.com/checkout/?accountcart=1")

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response CancelTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Success != 2 {
		return fmt.Errorf("取消交易失败,返回Success: %d", response.Success)
	}

	return nil
}

type GetFinalPriceResponse struct {
	Success                    int    `json:"success"`
	Total                      int    `json:"total"`
	FormattedSteamAccountTotal string `json:"formattedSteamAccountTotal"`
	FormattedProviderTotal     string `json:"formattedProviderTotal"`
}

// FinalPriceDetails 包含解析后的价格详情
type FinalPriceDetails struct {
	SteamAccountTotal float64 // Steam账户总额（浮点数）
	ProviderTotal     float64 // 提供商总额（浮点数）
	Total             int     // 总金额（整数，单位：分）
}

// parsePriceString 解析价格字符串（如 "¥ 45.90"）并返回浮点数
func parsePriceString(priceStr string) (float64, error) {
	// 移除所有空格和货币符号
	priceStr = strings.ReplaceAll(priceStr, "¥", "")
	priceStr = strings.ReplaceAll(priceStr, " ", "")
	priceStr = strings.ReplaceAll(priceStr, ",", "")
	priceStr = strings.TrimSpace(priceStr)

	// 转换为浮点数
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("解析价格失败 '%s': %w", priceStr, err)
	}
	return price, nil
}

func (d *Dao) GetFinalPrice(transactionID string) (int, error) {
	Logger.Infof("[%s]获取最终价格", d.GetUsername())

	params := Param.Params{}
	params.SetInt64("count", 1)
	params.SetString("transid", transactionID)
	params.SetString("purchasetype", "self")
	params.SetInt64("microtxnid", -1)
	params.SetInt64("cart", -1)
	params.SetInt64("gidReplayOfTransID", -1)

	req, err := d.Request(http.MethodGet, Constants.Getfinalprice+"?"+params.ToUrl(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return 0, fmt.Errorf("获取最终价格失败,返回状态码: %d", resp.StatusCode)
	}

	var response GetFinalPriceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("解析最终价格响应失败: %w", err)
	}

	if response.Success != 1 {
		fmt.Println(string(body))
		return 0, fmt.Errorf("获取最终价格失败,返回Success: %d", response.Success)
	}

	return response.Total, nil
}

// GetFinalPriceWithDetails 获取最终价格并返回解析后的价格详情（包含浮点数）
func (d *Dao) GetFinalPriceWithDetails(transactionID string) (*FinalPriceDetails, error) {
	Logger.Infof("[%s]获取最终价格详情", d.GetUsername())

	params := Param.Params{}
	params.SetInt64("count", 1)
	params.SetString("transid", transactionID)
	params.SetString("purchasetype", "self")
	params.SetInt64("microtxnid", -1)
	params.SetInt64("cart", -1)
	params.SetInt64("gidReplayOfTransID", -1)

	req, err := d.Request(http.MethodGet, Constants.Getfinalprice+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
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

	fmt.Println(string(body))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取最终价格失败,返回状态码: %d", resp.StatusCode)
	}

	var response GetFinalPriceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("解析最终价格响应失败: %w", err)
	}

	// {1 9600 ¥ 45.90 ¥ 50.10}
	fmt.Println(response)

	if response.Success != 1 {
		fmt.Println(string(body))
		return nil, fmt.Errorf("获取最终价格失败,返回Success: %d", response.Success)
	}

	// 解析格式化的价格字符串为浮点数
	steamAccountTotal, err := parsePriceString(response.FormattedSteamAccountTotal)
	if err != nil {
		return nil, fmt.Errorf("解析 SteamAccountTotal 失败: %w", err)
	}

	providerTotal, err := parsePriceString(response.FormattedProviderTotal)
	if err != nil {
		return nil, fmt.Errorf("解析 ProviderTotal 失败: %w", err)
	}

	return &FinalPriceDetails{
		SteamAccountTotal: steamAccountTotal,
		ProviderTotal:     providerTotal,
		Total:             response.Total,
	}, nil
}

func (d *Dao) GetAddFundFinalPrice(transactionID, cart, microTxn string) (int, error) {
	Logger.Infof("[%s]获取最终价格", d.GetUsername())

	params := Param.Params{}
	params.SetInt64("count", 1)
	params.SetString("transid", transactionID)
	params.SetString("purchasetype", "self")
	params.SetString("microtxnid", microTxn)
	params.SetString("cart", cart)
	params.SetInt64("gidReplayOfTransID", -1)

	req, err := d.Request(http.MethodGet, Constants.Getfinalprice+"?"+params.ToUrl(), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("referer", fmt.Sprintf("https: //checkout.steampowered.com/checkout?cart=%s&microtxn=%s", cart, microTxn))

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("获取最终价格失败,返回状态码: %d", resp.StatusCode)
	}

	var response GetFinalPriceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("解析最终价格响应失败: %w", err)
	}

	if response.Success != 1 {
		fmt.Println(string(body))
		return 0, fmt.Errorf("获取最终价格失败,返回Success: %d", response.Success)
	}

	return response.Total, nil
}

func (d *Dao) TestGetPayLinkAgain() (string, error) {
	paramsForPayLink := Param.Params{}
	// 	  Amount: 141
	//   Currency: CNY
	//   MethodID: 24
	//   Description: Steam Purchase
	//   SkinID: 101
	//   MerchantID: 1102
	//   ReturnURL: https://checkout.steampowered.com/paypal/smart2pay/56990384110504959/
	//   Country: CN
	//   CustomerEmail: MauriceDonovan8084@outlook.com
	//   CustomerName: dih9u8nad
	//   SkipHPP: 1
	//   Articles: Name=Lossless Scaling, DYNASTY WARRIORS ORIGINS Visions of Four Heroes with Pre-purchase Bonus, PRAGMATA Deluxe Edition&Quantity=1
	//   Hash: 3a1aec77155182fbe2bcf021ec9af8d504d1c7c6f6ca4bc0841fe2c4dcd7419f
	//   MerchantTransactionID: 56990384110504959
	paramsForPayLink.SetString("Amount", "141")
	paramsForPayLink.SetString("Currency", "CNY")
	paramsForPayLink.SetString("MethodID", "24")
	paramsForPayLink.SetString("Description", "Steam Purchase")
	paramsForPayLink.SetString("SkinID", "101")
	paramsForPayLink.SetString("MerchantID", "1102")
	paramsForPayLink.SetString("ReturnURL", "https://checkout.steampowered.com/paypal/smart2pay/56990384110504959/")
	paramsForPayLink.SetString("Country", "CN")
	paramsForPayLink.SetString("CustomerEmail", "MauriceDonovan8084@outlook.com")
	paramsForPayLink.SetString("CustomerName", "dih9u8nad")
	paramsForPayLink.SetString("SkipHPP", "1")
	paramsForPayLink.SetString("Articles", "Name=Lossless Scaling, DYNASTY WARRIORS ORIGINS Visions of Four Heroes with Pre-purchase Bonus, PRAGMATA Deluxe Edition&Quantity=1")
	paramsForPayLink.SetString("Hash", "3a1aec77155182fbe2bcf021ec9af8d504d1c7c6f6ca4bc0841fe2c4dcd7419f")
	paramsForPayLink.SetString("MerchantTransactionID", "56990384110504959")

	paslice := []string{"MerchantID", "MerchantTransactionID", "Amount", "Currency", "ReturnURL", "MethodID", "Country", "CustomerEmail", "CustomerName", "SkipHPP", "Articles", "Description", "SkinID", "Hash"}
	reqForPayLink, err := d.NewRequest(http.MethodPost, "https://globalapi.smart2pay.com", strings.NewReader(paramsForPayLink.EncodeBy(paslice)))

	if err != nil {
		return "", err
	}

	respForPayLink, err := d.RetryRequest(Constants.Tries, reqForPayLink)
	if err != nil {
		return "", err
	}
	defer respForPayLink.Body.Close()

	// 检查响应是否被 gzip 压缩
	var reader io.Reader = respForPayLink.Body
	if respForPayLink.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(respForPayLink.Body)
		if err != nil {
			return "", err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	bodyForPayLink, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	redirectURL, err := ExtractRedirectURL(string(bodyForPayLink))
	if err != nil {
		return "", err
	}
	if redirectURL == "" {
		return "", errors.New("未找到重定向URL")
	}

	return "https://globalep1.smart2pay.com/" + redirectURL, nil
}

func (d *Dao) AccessCheckoutURL(transactionID string) (string, error) {
	params := Param.Params{}
	params.SetString("transid", transactionID)

	req, err := d.Request(http.MethodGet, Constants.ExternallLink+"?"+params.ToUrl(), nil)
	if err != nil {
		return "", err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析HTML表单
	formData, err := ParsePaymentForm(string(body))
	if err != nil {
		return "", err
	}

	paramsForPayLink := Param.Params{}

	for key, value := range formData.Fields {
		fmt.Printf("  %s: %s\n", key, value)
		paramsForPayLink.SetString(key, value)
	}

	paslice := []string{"MerchantID", "MerchantTransactionID", "Amount", "Currency", "ReturnURL", "MethodID", "Country", "CustomerEmail", "CustomerName", "SkipHPP", "Articles", "Description", "SkinID", "Hash"}
	reqForPayLink, err := d.NewRequest(http.MethodPost, formData.Action, strings.NewReader(paramsForPayLink.EncodeBy(paslice)))

	if err != nil {
		return "", err
	}

	respForPayLink, err := d.RetryRequest(Constants.Tries, reqForPayLink)
	if err != nil {
		return "", err
	}
	defer respForPayLink.Body.Close()

	// 检查响应是否被 gzip 压缩
	var reader io.Reader = respForPayLink.Body
	if respForPayLink.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(respForPayLink.Body)
		if err != nil {
			return "", err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	bodyForPayLink, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	redirectURL, err := ExtractRedirectURL(string(bodyForPayLink))
	if err != nil {
		return "", err
	}
	if redirectURL == "" {
		return "", errors.New("未找到重定向URL")
	}

	return "https://globalep1.smart2pay.com/" + redirectURL, nil
}

// ExtractRedirectURL 从HTML中提取重定向URL
// 处理像 "Object moved to <a href="...">here</a>" 这样的重定向页面
func ExtractRedirectURL(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %w", err)
	}

	var redirectURL string
	var findLink func(*html.Node)
	findLink = func(n *html.Node) {
		if redirectURL != "" {
			return
		}

		// 查找 <a> 标签
		if n.Type == html.ElementNode && n.Data == "form" {
			for _, attr := range n.Attr {
				if attr.Key == "action" {
					// 检查链接是否是有效的URL
					if strings.HasPrefix(attr.Val, "AlipayPlus/Payment") {
						redirectURL = attr.Val
						return
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			findLink(child)
		}
	}

	findLink(doc)

	if redirectURL == "" {
		return "", errors.New("未找到重定向URL")
	}

	// HTML实体已经被html.Parse自动解码，所以 &amp; 会变成 &
	return redirectURL, nil
}

// PaymentFormData 存储支付表单的数据
type PaymentFormData struct {
	Action string            // 表单的action URL
	Fields map[string]string // 表单中的所有隐藏字段
}

// ParsePaymentForm 解析HTML中的支付表单
func ParsePaymentForm(htmlContent string) (*PaymentFormData, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	formData := &PaymentFormData{
		Fields: make(map[string]string),
	}

	// 递归遍历HTML节点查找form
	var findForm func(*html.Node) bool
	findForm = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "form" {
			// 获取form的action属性
			for _, attr := range n.Attr {
				if attr.Key == "action" {
					formData.Action = attr.Val
					break
				}
			}

			// 查找所有input字段
			var findInputs func(*html.Node)
			findInputs = func(node *html.Node) {
				if node.Type == html.ElementNode && node.Data == "input" {
					var name, value string
					var isHidden bool

					for _, attr := range node.Attr {
						switch attr.Key {
						case "name":
							name = attr.Val
						case "value":
							value = attr.Val
						case "type":
							if attr.Val == "hidden" {
								isHidden = true
							}
						}
					}

					// 只添加隐藏字段
					if isHidden && name != "" {
						formData.Fields[name] = value
					}
				}

				for child := node.FirstChild; child != nil; child = child.NextSibling {
					findInputs(child)
				}
			}

			findInputs(n)
			return true
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if findForm(child) {
				return true
			}
		}
		return false
	}

	if !findForm(doc) {
		return nil, errors.New("未找到支付表单")
	}

	if formData.Action == "" {
		return nil, errors.New("表单没有action属性")
	}

	return formData, nil
}

// 在 checkout.go 中添加新函数
func (d *Dao) GetAlipayURL(transID string) (string, error) {
	externallinkURL := fmt.Sprintf("https://checkout.steampowered.com/checkout/externallink/?transid=%s", transID)

	// 创建一个不自动跟随重定向的客户端
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 返回错误以停止自动重定向，手动处理
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, externallinkURL, nil)
	if err != nil {
		return "", err
	}

	// 添加必要的 Cookie
	cookies := d.GetLoginCookies()["checkout.steampowered.com"]
	if cookies != nil {
		req.AddCookie(&http.Cookie{
			Name:  "sessionid",
			Value: cookies.SessionId,
		})
		req.AddCookie(&http.Cookie{
			Name:  "steamLoginSecure",
			Value: cookies.SteamLoginSecure,
		})
	}

	req.Header.Set("User-Agent", "Mozilla/5.0...")

	// 跟踪重定向链
	finalURL := externallinkURL
	for i := 0; i < 10; i++ { // 最多跟踪10次重定向
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			if location == "" {
				break
			}
			finalURL = location
			req, err = http.NewRequest(http.MethodGet, location, nil)
			if err != nil {
				return "", err
			}
			// 继续使用相同的 Cookie
			if cookies != nil {
				req.AddCookie(&http.Cookie{
					Name:  "sessionid",
					Value: cookies.SessionId,
				})
			}
		} else {
			break
		}
	}

	return finalURL, nil
}

type UnsendGiftResponse struct {
	Success int `json:"success"`
}

func (d *Dao) UnsendGift(giftId string) error {
	if d.GetLoginCookies()["checkout.steampowered.com"] == nil {
		return errors.New("checkout.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["checkout.steampowered.com"].SessionId

	params := Param.Params{}
	params.SetString("GiftGID", giftId)
	params.SetString("SessionID", sessionId)

	req, err := d.Request(http.MethodPost, Constants.UnsendGiftSubmit, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("撤回赠送礼物失败,返回状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response UnsendGiftResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Success != 1 {
		return fmt.Errorf("撤回赠送礼物失败,返回Success: %d", response.Success)
	}

	return nil
}

func (d *Dao) TransactionStatus(transId string, count int) error {
	var result Model.TransactionStatusResponse

	params := Param.Params{}
	params.SetInt64("count", int64(count))
	params.SetString("transid", transId)

	req, err := d.Request(http.MethodGet, Constants.TransactionStatus+"?"+params.ToUrl(), nil)
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

	if resp.StatusCode != 200 {
		return errors.New("获取交易状态失败,返回状态码: " + strconv.Itoa(resp.StatusCode))
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Success != 1 {
		return fmt.Errorf("订单未完成,返回Success: %d", result.Success)
	}

	return nil
}

func (d *Dao) SetCountry(country string) error {
	if d.GetLoginCookies()["store.steampowered.com"] == nil {
		return errors.New("store.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["store.steampowered.com"].SessionId

	params := Param.Params{}
	params.SetString("sessionid", sessionId)
	params.SetString("cc", country)

	req, err := d.Request(http.MethodPost, Constants.SetCountry, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Referer", "https://store.steampowered.com/steamaccount/addfunds")
	req.Header.Set("Origin", "https://store.steampowered.com")

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("SetCountry response:", string(body))

	if resp.StatusCode != 200 {
		return fmt.Errorf("SetCountry error,返回状态码: %d", resp.StatusCode)
	}

	return nil
}

// action add_to_cart
// currency CNY
// amount 3001
// sessionID 3e5bdd99e638b8ac492ed063
// mtreturnurl
func (d *Dao) AddFundsSubmit(amount int) (string, error) {
	if d.GetLoginCookies()["store.steampowered.com"] == nil {
		fmt.Println("store.steampowered.com cookie not found")
		return "", fmt.Errorf("store.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["store.steampowered.com"].SessionId

	params := Param.Params{}
	params.SetString("action", "add_to_cart")
	params.SetString("currency", "CNY")
	params.SetInt64("amount", int64(amount))
	params.SetString("sessionID", sessionId)
	params.SetString("mtreturnurl", "")

	req, err := d.Request(http.MethodPost, Constants.AddFundsSubmit, strings.NewReader(params.Encode()))
	if err != nil {
		fmt.Println("d.Request error:", err)
		return "", fmt.Errorf("d.Request error: %w", err)
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		fmt.Println("d.RetryRequest error:", err)
		return "", fmt.Errorf("d.RetryRequest error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AddFundsSubmit error,返回状态码: %d", resp.StatusCode)
	}

	finalURL := resp.Request.URL
	cart := finalURL.Query().Get("cart")
	microtxn := finalURL.Query().Get("microtxn")
	if cart == "" || microtxn == "" {
		redir := finalURL.Query().Get("redir")
		if redir != "" {
			if redirURL, parseErr := url.Parse(redir); parseErr == nil {
				if cart == "" {
					cart = redirURL.Query().Get("cart")
				}
				if microtxn == "" {
					microtxn = redirURL.Query().Get("microtxn")
				}
			}
		}
	}

	fmt.Println("最终状态码:", resp.StatusCode)
	fmt.Println("最终URL:", finalURL.String())
	fmt.Println("cart:", cart)
	fmt.Println("microtxn:", microtxn)

	transId, err := d.InitAddFundsTransaction(cart, microtxn)
	if err != nil {
		return "", fmt.Errorf("InitAddFundsTransaction error: %w", err)
	}

	rs, err := d.GetAddFundFinalPrice(transId, cart, microtxn)
	if err != nil {
		return "", fmt.Errorf("GetAddFundFinalPrice error: %w", err)
	}
	fmt.Println("rs:", rs)
	payLink, err := d.AccessCheckoutURL(transId)
	if err != nil {
		return "", fmt.Errorf("AccessCheckoutURL error: %w", err)
	}
	fmt.Println(payLink)

	return payLink, nil
}

func (d *Dao) AddFundsSubmitWithCountry(amount int, country string) (string, error) {
	if d.GetLoginCookies()["store.steampowered.com"] == nil {
		fmt.Println("store.steampowered.com cookie not found")
		return "", fmt.Errorf("store.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["store.steampowered.com"].SessionId

	params := Param.Params{}
	params.SetString("action", "add_to_cart")
	params.SetString("currency", Constants.Countries[country].CurrencyCode)
	params.SetInt64("amount", int64(amount))
	params.SetString("sessionID", sessionId)
	params.SetString("mtreturnurl", "")

	req, err := d.Request(http.MethodPost, Constants.AddFundsSubmit, strings.NewReader(params.Encode()))
	if err != nil {
		fmt.Println("d.Request error:", err)
		return "", fmt.Errorf("d.Request error: %w", err)
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		fmt.Println("d.RetryRequest error:", err)
		return "", fmt.Errorf("d.RetryRequest error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AddFundsSubmit error,返回状态码: %d", resp.StatusCode)
	}

	finalURL := resp.Request.URL
	cart := finalURL.Query().Get("cart")
	microtxn := finalURL.Query().Get("microtxn")
	if cart == "" || microtxn == "" {
		redir := finalURL.Query().Get("redir")
		if redir != "" {
			if redirURL, parseErr := url.Parse(redir); parseErr == nil {
				if cart == "" {
					cart = redirURL.Query().Get("cart")
				}
				if microtxn == "" {
					microtxn = redirURL.Query().Get("microtxn")
				}
			}
		}
	}

	fmt.Println("最终状态码:", resp.StatusCode)
	fmt.Println("最终URL:", finalURL.String())
	fmt.Println("cart:", cart)
	fmt.Println("microtxn:", microtxn)

	transId, err := d.InitAddFundsTransactionWithCountry(cart, microtxn, country)
	if err != nil {
		return "", fmt.Errorf("InitAddFundsTransactionWithCountry error: %w", err)
	}

	rs, err := d.GetAddFundFinalPrice(transId, cart, microtxn)
	if err != nil {
		return "", fmt.Errorf("GetAddFundFinalPrice error: %w", err)
	}
	fmt.Println("rs:", rs)
	payLink, err := d.AccessCheckoutURL(transId)
	if err != nil {
		return "", fmt.Errorf("AccessCheckoutURL error: %w", err)
	}
	fmt.Println(payLink)

	return payLink, nil
}

func (d *Dao) LogoutAll() error {
	if d.GetLoginCookies()["store.steampowered.com"] == nil {
		fmt.Println("store.steampowered.com cookie not found")
		return fmt.Errorf("store.steampowered.com cookie not found")
	}
	sessionId := d.GetLoginCookies()["store.steampowered.com"].SessionId

	params := Param.Params{}
	params.SetString("action", "deauthorize")
	params.SetString("sessionid", sessionId)

	fmt.Println(Constants.LogoutAll)
	fmt.Println(params.Encode())

	req, err := d.Request(http.MethodPost, Constants.LogoutAll, strings.NewReader(params.Encode()))
	if err != nil {
		fmt.Println("d.Request error:", err)
		return fmt.Errorf("d.Request error: %w", err)
	}

	req.Header.Set("origin", "https://store.steampowered.com")
	req.Header.Set("referer", "https://store.steampowered.com/account/authorizeddevices")

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		fmt.Println("d.RetryRequest error:", err)
		return fmt.Errorf("d.RetryRequest error: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("最终状态码:", resp.StatusCode)
	fmt.Println("最终URL:", resp.Request.URL.String())

	if resp.StatusCode != 200 {
		return fmt.Errorf("LogoutAll error,返回状态码: %d", resp.StatusCode)
	}

	return nil
}
