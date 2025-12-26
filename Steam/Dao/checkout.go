package Dao

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
		return "", fmt.Errorf("初始化交易失败,返回状态码: %d", resp.StatusCode)
	}

	fmt.Println(string(body))

	var response InitTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析初始化交易响应失败: %w", err)
	}

	if response.Success != 1 {
		return "", errors.New(Errors.GetCheckoutError(response.PurchaseResultDetail))
	}

	return response.TransID, nil
}

type CancelTransactionResponse struct {
	Success int `json:"success"`
}

func (d *Dao) CancelTransaction(transactionID string) error {
	params := Param.Params{}
	params.SetString("transid", transactionID)
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
	Success int `json:"success"`
	Total   int `json:"total"`
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
		return 0, fmt.Errorf("获取最终价格失败,返回状态码: %d", resp.StatusCode)
	}

	var response GetFinalPriceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("解析最终价格响应失败: %w", err)
	}

	if response.Success != 1 {
		return 0, fmt.Errorf("获取最终价格失败,返回Success: %d", response.Success)
	}

	return response.Total, nil
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

	// reqForPayLink.Header.Set("Upgrade-Insecure-Requests", "1")
	// reqForPayLink.Header.Set("Origin", "https://checkout.steampowered.com")
	// reqForPayLink.Header.Set("Referer", "https://checkout.steampowered.com/checkout/?accountcart=1")
	// reqForPayLink.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// reqForPayLink.Header.Set("Sec-Fetch-Site", "cross-site")
	// reqForPayLink.Header.Set("Sec-Fetch-Mode", "navigate")
	// reqForPayLink.Header.Set("Sec-Fetch-Dest", "document")
	// reqForPayLink.Header.Set("Accept-Encoding", "gzip, deflate")
	// reqForPayLink.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	// reqForPayLink.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	// reqForPayLink.Header.Set("sec-ch-ua", "Not=A?Brand\";v=\"24\", \"Chromium\";v=\"140\"")
	// reqForPayLink.Header.Set("sec-ch-ua-mobile", "?0")
	// reqForPayLink.Header.Set("sec-ch-ua-platform", "Windows")

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
