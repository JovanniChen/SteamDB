package Dao

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Protoc"
	"google.golang.org/protobuf/proto"
)

func (d *Dao) ClearCart() error {
	accessToken, _ := d.AccessToken()
	params := Param.Params{}
	params.SetString("access_token", accessToken)

	req, err := d.NewRequest(http.MethodPost, Constants.ClearCart+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Errors.ErrClearCartFailed
	}

	return nil
}

func (d *Dao) GetCart() error {
	req, err := d.Request(http.MethodGet, Constants.CartIndex, nil)
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

func (d *Dao) AddItemToCart(addCartItems []Model.AddCartItem) error {
	items := make([]*Protoc.Item, 0)
	for _, addCartItem := range addCartItems {
		var item *Protoc.Item
		if addCartItem.BundleID != 0 {
			item = &Protoc.Item{
				Bundleid: addCartItem.BundleID,
			}
		} else {
			item = &Protoc.Item{
				Packageid: addCartItem.PackageID,
			}
		}
		item.GiftInfo = &Protoc.GiftInfo{
			AccountidGiftee: int32(addCartItem.AccountidGiftee),
			GiftMessage: &Protoc.GiftMessage{
				Gifteename: "",
				Message:    addCartItem.Message,
				Sentiment:  "",
				Signature:  "",
			},
		}
		item.Flag = &Protoc.Flag{
			IsGift:    true,
			IsPrivate: false,
		}

		items = append(items, item)
	}

	addCartSend := &Protoc.AddCartSend{
		UserCountry: d.GetCountryCode(),
		Items:       items,
	}

	// 序列化为protobuf格式
	data, err := proto.Marshal(addCartSend)
	if err != nil {
		return err
	}

	accessToken, _ := d.AccessToken()
	params := Param.Params{}
	params.SetString("access_token", accessToken)

	// 构建POST请求体参数(包含protobuf数据)
	params1 := Param.Params{}
	params1.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	req, err := d.NewRequest(http.MethodPost, Constants.AddItemsToCart+"?"+params.ToUrl(), strings.NewReader(params1.Encode()))
	if err != nil {
		return err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Header.Get("X-Eresult") != "1" {
		return fmt.Errorf("add item to cart failed: %s", resp.Header.Get("X-Eresult"))
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	// 解析protobuf响应
	addCartReceive := &Protoc.AddCartReceive{}

	// 使用重试机制解析protobuf
	if err := protoUnmarshalWithRetry(buf.Bytes(), addCartReceive, "AddItemToCart", 3); err != nil {
		return err
	}

	// fmt.Println("******************************************************")
	// fmt.Println(addCartReceive)
	// fmt.Println("******************************************************")
	// fmt.Println(addCartReceive.Cart.Subtotal)
	// fmt.Println("******************************************************")

	return nil
}

func (d *Dao) ValidateCart() error {
	accessToken, _ := d.AccessToken()
	params := Param.Params{}
	params.SetString("access_token", accessToken)

	req, err := d.NewRequest(http.MethodGet, Constants.ValidateCart+"?"+params.ToUrl(), nil)
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
	fmt.Println("=====", string(body))

	return nil
}

func (d *Dao) GetProductByAppUrl(url string) error {
	req, err := d.Request(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	cookies := d.GetLoginCookies()["store.steampowered.com"]
	if cookies != nil {
		req.Header.Add("cookie", fmt.Sprintf("sessionid=%s;steamLoginSecure=%s", cookies.SessionId, cookies.SteamLoginSecure))
	}

	req.Header.Add("cookie", "birthtime=325958401;lastagecheckage=1-May-1980;wants_mature_content=1;Steam_Language=schinese")
	req.Header.Add("cookie", "app_impressions=1222140@1_5_9__412|960990@1_5_9__412|960910@1_5_9__412|312840@1_5_9__412|1222140@1_5_9__412|960990@1_5_9__412|960910@1_5_9__412|312840@1_5_9__412|1222140@1_5_9__412|960990@1_5_9__412|960910@1_5_9__412|312840@1_5_9__412|252490@1_8_512_513_520|218620@1_8_512_513_520|252490@1_8_512_513_520|218620@1_8_512_513_520|1222140@1_5_9__412|960990@1_5_9__412|960910@1_5_9__412|312840@1_5_9__412|292030@1_5_9__300|960990@1_5_9__300|990080@1_5_9__300|960910@1_5_9__300|1222140@1_5_9__412|960990@1_5_9__412|960910@1_5_9__412|312840@1_5_9__412|292030@1_5_9__300|960990@1_5_9__300|1091500@1_5_9__300|960910@1_5_9__300|2947440@1_430_4__300|567640@1_430_4__300|2172010@1_430_4__300|3101040@1_430_4__300|1222140@1_5_9__412|960990@1_5_9__412|960910@1_5_9__412|312840@1_5_9__412|1222140@1_5_9__412|960990@1_5_9__412|960910@1_5_9__412|312840@1_5_9__412")
	req.Header.Add("cookie", "recentapps=%7B%221222140%22%3A1766996240%7D")

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
	// fmt.Println("=====", string(body))
	// 保存这个string(body)到项目根目录
	os.WriteFile("product.html", body, 0644)

	return nil
}
