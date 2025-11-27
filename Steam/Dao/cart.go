package Dao

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
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

func (d *Dao) AddItemToCart() error {
	fmt.Println("d.GetCountryCode() =", d.GetCountryCode())

	item := &Protoc.Item{
		Packageid: 1074193,
		GiftInfo: &Protoc.GiftInfo{
			AccountidGiftee: 352956450,
			GiftMessage: &Protoc.GiftMessage{
				Gifteename: "",
				Message:    "test",
				Sentiment:  "",
				Signature:  "",
			},
		},
		Flag: &Protoc.Flag{
			IsGift:    true,
			IsPrivate: false,
		},
	}

	addCartSend := &Protoc.AddCartSend{
		UserCountry: d.GetCountryCode(),
		Items: []*Protoc.Item{
			item,
		},
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

	fmt.Println(resp.StatusCode)

	fmt.Println(resp.Header.Get("X-Eresult"))
	if resp.Header.Get("X-Eresult") != "1" {
		return errors.New("add item to cart failed")
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

	fmt.Println(addCartReceive)

	return nil
}
