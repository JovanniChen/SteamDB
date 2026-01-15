package Dao

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Protoc"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
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

func (d *Dao) GetProductByAppID(appID int) (map[string]Model.GamePurchaseAction, error) {
	url := fmt.Sprintf("https://store.steampowered.com/app/%d", appID)
	products, err := d.GetProductByAppUrl(url)
	if err != nil {
		return nil, err
	}
	// 转换成map[string]Model.GamePurchaseAction
	productMap := make(map[string]Model.GamePurchaseAction)
	for _, product := range products {
		productMap[product.AddToCartIds] = product
	}
	return productMap, nil
}

func (d *Dao) GetProductByAppUrl(url string) ([]Model.GamePurchaseAction, error) {
	// ?cc=cn&l=schinese
	req, err := d.Request(http.MethodGet, url+"?cc=cn&l=schinese", nil)
	if err != nil {
		return nil, err
	}

	cookies := d.GetLoginCookies()["store.steampowered.com"]
	if cookies != nil {
		req.Header.Add("cookie", fmt.Sprintf("sessionid=%s;steamLoginSecure=%s", cookies.SessionId, cookies.SteamLoginSecure))
	}

	req.AddCookie(&http.Cookie{Name: "birthtime", Value: "0"})
	req.AddCookie(&http.Cookie{Name: "lastagecheckage", Value: "1-January-1970"})
	req.AddCookie(&http.Cookie{Name: "mature_content", Value: "1"})
	req.AddCookie(&http.Cookie{Name: "wants_mature_content", Value: "1"})

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.StatusCode)
	// fmt.Println("=====", string(body))
	// 保存这个string(body)到项目根目录
	// os.WriteFile("product.html", body, 0644)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取商品信息失败,返回状态码: %d,url: %s", resp.StatusCode, url)
	}

	// 解析游戏购买信息
	results, err := ParseGamePurchaseActions(string(body), url)
	if err != nil {
		fmt.Printf("解析游戏购买信息失败,url: %s,err: %v", url, err)
		return nil, fmt.Errorf("解析游戏购买信息失败")
	}

	fmt.Printf("解析到 %d 个购买选项:\n", len(results))
	for i, result := range results {
		fmt.Printf("选项 %d:\n", i+1)
		fmt.Printf("  是否为捆绑包: %d\n", result.IsBundle)
		fmt.Printf("  游戏名称: %s\n", result.GameName)
		fmt.Printf("  价格: %s\n", result.FinalPrice)
		fmt.Printf("  类型: %s\n", result.BundleInfoTexts)
		fmt.Printf("  购物车ID: %s\n", result.AddToCartIds)
		fmt.Println()
	}

	return results, nil
}

// extractCountryCode 从HTML中提取国家代码
func extractCountryCode(doc *html.Node) string {
	// 查找 application_config 元素
	configNode := htmlquery.FindOne(doc, "//div[@id='application_config']")
	if configNode == nil {
		return "CN" // 默认返回CN
	}

	// 获取 data-config 属性
	configAttr := htmlquery.SelectAttr(configNode, "data-config")
	if configAttr == "" {
		return "CN"
	}

	// 解析JSON配置
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configAttr), &config); err != nil {
		return "CN"
	}

	// 提取国家代码
	if country, ok := config["COUNTRY"].(string); ok {
		return country
	}

	return "CN"
}

// extractFinalPrice 从HTML节点中提取最终价格
func extractFinalPrice(priceNode *html.Node) string {
	if priceNode == nil {
		return ""
	}

	// 优先从 data-price-final 属性提取（单位是分，需要除以100）
	priceFinal := htmlquery.SelectAttr(priceNode, "data-price-final")
	if priceFinal != "" {
		// 将分转换为元
		if price, err := strconv.ParseFloat(priceFinal, 64); err == nil {
			return fmt.Sprintf("%.2f", price/100.0)
		}
	}

	// 如果没有 data-price-final，从文本内容提取
	priceText := strings.TrimSpace(htmlquery.InnerText(priceNode))
	if priceText != "" {
		// 使用正则表达式提取价格数字
		priceRegex := regexp.MustCompile(`(?i)(?:[￥$]|CNY|HKD|USD)\s*(\d+\.?\d*)`)
		matches := priceRegex.FindStringSubmatch(priceText)
		if len(matches) > 1 {
			return matches[1]
		}

		// 如果没有找到带货币符号的价格，尝试只匹配数字
		numberRegex := regexp.MustCompile(`(\d+\.?\d*)`)
		matches = numberRegex.FindStringSubmatch(priceText)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// CartInfo 购物车信息结构
type CartInfo struct {
	ID   string
	Type string
}

// extractCartInfo 从wrapper节点中提取购物车信息
func extractCartInfo(wrapper *html.Node, isBundle bool) *CartInfo {
	cartInfo := &CartInfo{}

	// 查找添加到购物车的按钮
	// 尝试多种可能的选择器
	addToCartNodes := htmlquery.Find(wrapper, ".//a[contains(@class, 'btn_addtocart')]")
	if len(addToCartNodes) == 0 {
		addToCartNodes = htmlquery.Find(wrapper, ".//a[contains(@onclick, 'addtocart')]")
	}
	if len(addToCartNodes) == 0 {
		addToCartNodes = htmlquery.Find(wrapper, ".//div[contains(@class, 'game_purchase_action')]//a")
	}

	for _, node := range addToCartNodes {
		// 跳过"捆绑包信息"链接（btn_packageinfo 类）
		nodeClass := htmlquery.SelectAttr(node, "class")
		if strings.Contains(nodeClass, "btn_packageinfo") {
			continue
		}

		// 检查链接文本，跳过"捆绑包信息"链接
		linkText := strings.ToLower(strings.TrimSpace(htmlquery.InnerText(node)))
		if strings.Contains(linkText, "捆绑包信息") || strings.Contains(linkText, "bundle info") || strings.Contains(linkText, "package info") {
			continue
		}
		// 优先检查 href 属性（Steam 使用 javascript:addToCart()）
		href := htmlquery.SelectAttr(node, "href")
		if href != "" {
			// 提取ID，匹配类似 javascript:addToCart(123456) 或 javascript:addBundleToCart( 123456 )（允许空格）
			re := regexp.MustCompile(`javascript:(?:addToCart|addBundleToCart)\s*\(\s*(\d+)\s*\)`)
			matches := re.FindStringSubmatch(href)
			if len(matches) > 1 {
				cartInfo.ID = strings.TrimSpace(matches[1])
				if strings.Contains(href, "addBundleToCart") {
					cartInfo.Type = "addbundletocart"
				} else {
					cartInfo.Type = "addtocart"
				}
				return cartInfo
			}
		}

		// 获取onclick属性
		onclick := htmlquery.SelectAttr(node, "onclick")
		if onclick != "" {
			// 提取ID，匹配类似 addtocart(123456) 或 addbundletocart( 123456 )（允许空格）
			re := regexp.MustCompile(`(?:addtocart|addbundletocart)\s*\(\s*(\d+)\s*\)`)
			matches := re.FindStringSubmatch(onclick)
			if len(matches) > 1 {
				cartInfo.ID = strings.TrimSpace(matches[1])
				if strings.Contains(onclick, "addbundletocart") {
					cartInfo.Type = "addbundletocart"
				} else {
					cartInfo.Type = "addtocart"
				}
				return cartInfo
			}
		}

		// 尝试从data-ds-appid或data-ds-packageid获取
		appID := htmlquery.SelectAttr(node, "data-ds-appid")
		packageID := htmlquery.SelectAttr(node, "data-ds-packageid")
		if appID != "" {
			cartInfo.ID = appID
			cartInfo.Type = "addtocart"
			return cartInfo
		}
		if packageID != "" {
			cartInfo.ID = packageID
			if isBundle {
				cartInfo.Type = "addbundletocart"
			} else {
				cartInfo.Type = "addtocart"
			}
			return cartInfo
		}
	}

	// 如果找不到链接，尝试从表单中提取
	formNodes := htmlquery.Find(wrapper, ".//form")
	for _, formNode := range formNodes {
		formName := htmlquery.SelectAttr(formNode, "name")
		// 检查是否是添加到购物车的表单（包括 add_to_cart 和 add_bundle_to_cart）
		if strings.Contains(formName, "add_to_cart") || strings.Contains(formName, "add_bundle_to_cart") {
			// 优先查找 bundleid（捆绑包）
			bundleidNode := htmlquery.FindOne(formNode, ".//input[@name='bundleid']")
			if bundleidNode != nil {
				cartInfo.ID = strings.TrimSpace(htmlquery.SelectAttr(bundleidNode, "value"))
				cartInfo.Type = "addbundletocart"
				return cartInfo
			}
			// 然后查找 subid（标准版）
			subidNode := htmlquery.FindOne(formNode, ".//input[@name='subid']")
			if subidNode != nil {
				cartInfo.ID = strings.TrimSpace(htmlquery.SelectAttr(subidNode, "value"))
				cartInfo.Type = "addtocart"
				return cartInfo
			}
		}
	}

	return nil
}

// convertToBundleInfo 将类型转换为捆绑包信息文本
func convertToBundleInfo(cartType string) string {
	if cartType == "addbundletocart" || cartType == "bundle" {
		return "捆绑包"
	}
	return "标准版"
}

// ParseGamePurchaseActions 解析游戏购买操作信息
func ParseGamePurchaseActions(htmlContent, url string) ([]Model.GamePurchaseAction, error) {
	results := make([]Model.GamePurchaseAction, 0)

	// 解析HTML
	doc, err := htmlquery.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	// 提取国家代码
	countryCode := extractCountryCode(doc)
	if countryCode == "" {
		countryCode = "CN"
	}

	// 确定货币标志
	moneyFlag := "￥"
	if countryCode == "CN" {
		moneyFlag = "￥"
	} else if countryCode == "HK" {
		moneyFlag = "HK$"
	}

	// 判断URL类型
	isAppURL := regexp.MustCompile(`/app/`).MatchString(url)

	if isAppURL {
		// 处理 /app/ 类型的URL
		wrappers := htmlquery.Find(doc, "//div[contains(@class, 'game_area_purchase_game_wrapper')]")
		for _, wrapper := range wrappers {
			// 判断是否为捆绑包（检查 wrapper 本身是否有 dynamic_bundle_description 类）
			wrapperClass := htmlquery.SelectAttr(wrapper, "class")
			isBundle := strings.Contains(wrapperClass, "dynamic_bundle_description")

			// 提取游戏名称
			nameNodes := htmlquery.Find(wrapper, ".//*[self::h1 or self::h2][contains(@id, '_purchase_') or contains(@id, 'bundle_purchase_label_')]")
			gameName := ""
			if len(nameNodes) > 0 {
				gameName = strings.TrimSpace(htmlquery.InnerText(nameNodes[0]))
				// 清理游戏名称
				gameName = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(gameName, "")
				gameName = regexp.MustCompile(`^(购买|Buy|BUY)\s*`).ReplaceAllString(gameName, "")
				// 移除捆绑包标签
				gameName = regexp.MustCompile(`\s*捆绑包.*$`).ReplaceAllString(gameName, "")
				gameName = strings.TrimSpace(gameName)
			}

			// 提取价格节点
			priceNode := htmlquery.FindOne(wrapper, ".//div[contains(@class, 'game_purchase_price')]")
			if priceNode == nil {
				// 尝试从 discount_final_price 提取
				priceNode = htmlquery.FindOne(wrapper, ".//div[contains(@class, 'discount_final_price')]")
			}

			finalPrice := extractFinalPrice(priceNode)

			// 提取购物车信息
			cartInfo := extractCartInfo(wrapper, isBundle)

			if cartInfo != nil && cartInfo.ID != "" && finalPrice != "" {
				results = append(results, Model.GamePurchaseAction{
					IsBundle: func() int {
						if isBundle {
							return 1
						} else {
							return 0
						}
					}(),
					BundleInfoTexts: convertToBundleInfo(func() string {
						if isBundle {
							return "bundle"
						} else {
							return "standard"
						}
					}()),
					GameName:       gameName,
					FinalPrice:     finalPrice,
					FinalPriceText: moneyFlag + " " + finalPrice,
					CountryCode:    countryCode,
					AddToCartIds:   cartInfo.ID,
				})
			} else {
				fmt.Printf("解析失败 - gameName:%s, cartInfo:%+v, finalPrice:%s, url: %s\n", gameName, cartInfo, finalPrice, url)
			}
		}
	} else {
		// 处理其他类型的URL（如 /sub/ 或 /bundle/）
		wrappers := htmlquery.Find(doc, "//*[@id='game_area_purchase_top']/div")
		for _, wrapper := range wrappers {
			// 提取游戏名称（在wrapper内部查找h1）
			nameNode := htmlquery.FindOne(wrapper, ".//h1")
			if nameNode == nil {
				fmt.Printf("获取商品名称失败,url: %s 获取其中一个商品信息异常！\n", url)
				continue
			}

			gameName := strings.TrimSpace(htmlquery.InnerText(nameNode))
			// 清理游戏名称
			gameName = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(gameName, "")
			gameName = regexp.MustCompile(`^(购买|Buy|BUY)\s*`).ReplaceAllString(gameName, "")
			gameName = strings.TrimSpace(gameName)

			// 提取价格节点
			priceNode := htmlquery.FindOne(wrapper, ".//div[contains(@class, 'game_purchase_price')]")
			if priceNode == nil {
				// 尝试从 discount_final_price 提取
				priceNode = htmlquery.FindOne(wrapper, ".//div[contains(@class, 'discount_final_price')]")
			}

			finalPrice := extractFinalPrice(priceNode)

			// 提取购物车信息
			cartInfo := extractCartInfo(wrapper, false)

			if cartInfo != nil && cartInfo.ID != "" && finalPrice != "" {
				isBundle := cartInfo.Type == "addbundletocart"
				results = append(results, Model.GamePurchaseAction{
					IsBundle: func() int {
						if isBundle {
							return 1
						} else {
							return 0
						}
					}(),
					BundleInfoTexts: convertToBundleInfo(cartInfo.Type),
					GameName:        gameName,
					FinalPrice:      finalPrice,
					FinalPriceText:  moneyFlag + " " + finalPrice,
					CountryCode:     countryCode,
					AddToCartIds:    cartInfo.ID,
				})
			}
		}
	}

	return results, nil
}
