package Dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
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

func (d *Dao) InitTransaction() error {
	Logger.Infof("[%s]初始化交易", d.GetUsername())
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

	if resp.StatusCode != 200 {
		return Errors.GetCheckoutError(resp.StatusCode)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))

	return nil
}

func (d *Dao) CancelTransaction(transactionID string) error {
	Logger.Infof("[%s]取消交易", d.GetUsername())

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

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))

	return nil
}

func (d *Dao) GetFinalPrice(transactionID string) error {
	Logger.Infof("[%s]获取最终价格", d.GetUsername())

	params := Param.Params{}
	params.SetInt64("count", 5)
	params.SetString("transid", transactionID)
	params.SetString("purchasetype", "self")
	params.SetInt64("microtxnid", -1)
	params.SetInt64("cart", -1)
	params.SetInt64("gidReplayOfTransID", -1)

	req, err := d.Request(http.MethodGet, Constants.Getfinalprice+"?"+params.ToUrl(), nil)
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

func (d *Dao) UnsendGift(giftId string) error {
	Logger.Infof("撤回赠送礼物[%s][%s]", d.GetUsername(), giftId)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))

	return nil
}

//	{
//	    "success": 1,
//	    "purchaseresultdetail": 0,
//	    "purchasereceipt": {
//	        "paymentmethod": 11,
//	        "purchasestatus": 1,
//	        "resultdetail": 0,
//	        "baseprice": "3000",
//	        "totaldiscount": "2700",
//	        "tax": "0",
//	        "shipping": "0",
//	        "packageid": -1,
//	        "transactiontime": 1764225154,
//	        "transactionid": "193220467646477559",
//	        "currencycode": 23,
//	        "formattedTotal": "¥ 3.00",
//	        "rewardPointsBalance": "275"
//	    },
//	    "strReceiptPageHTML": "",
//	    "bShowBRSpecificCreditCardError": false
//	}
func (d *Dao) TransactionStatus(transId string) error {
	Logger.Infof("获取交易状态[%s][%s]", d.GetUsername(), transId)

	var result Model.TransactionStatusResponse

	for i := 0; i < 10; i++ {
		params := Param.Params{}
		params.SetInt64("count", int64(i+1))
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
		// fmt.Println(resp.StatusCode)
		// fmt.Println(string(body))

		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}

		fmt.Printf("result = %+v\n", result)

		// if result.Success == 1 {
		// 	return nil
		// }
		// time.Sleep(5 * time.Second)
	}
	return nil
}
