package Dao

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func (d *Dao) GetPackageDetails(subID int) error {
	params := Param.Params{}
	params.SetString("packageids", strconv.Itoa(subID))
	params.SetString("cc", "cn")
	params.SetString("l", "schinese")

	// 创建请求
	req, err := d.NewRequest(http.MethodGet, Constants.PackageDetails+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	// 发送请求，重定向会自动处理，cookie 会从 jar 中自动获取
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查最终重定向后的状态码
	if resp.StatusCode != 200 {
		return errors.New("status != 200")
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

// 获取商店购物记录
func (d *Dao) GetStorePurchaseHistory() (*Model.StorePurchaseHistoryResult, error) {
	req, err := d.Request(http.MethodGet, Constants.StorePurchaseHistory, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status != 200")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseStorePurchaseHistory(body)
}

func parseStorePurchaseHistory(body []byte) (*Model.StorePurchaseHistoryResult, error) {
	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	rows := htmlquery.Find(doc, `//tr[contains(concat(" ", normalize-space(@class), " "), " wallet_table_row ")]`)
	if len(rows) == 0 && htmlquery.FindOne(doc, `//table[contains(concat(" ", normalize-space(@class), " "), " wallet_history_table ")]`) == nil {
		title := nodeText(doc, `//title`)
		if title == "" {
			title = "未知页面"
		}
		return nil, fmt.Errorf("未获取到消费历史表格，当前页面标题: %s，可能是登录状态失效", title)
	}

	records := make([]Model.StorePurchaseHistoryRecord, 0, len(rows))
	for i, row := range rows {
		records = append(records, parseStorePurchaseHistoryRow(row, i+1))
	}

	return &Model.StorePurchaseHistoryResult{
		TotalRecords:                  len(records),
		LatestUnrefundedGiftPurchases: filterLatestUnrefundedGiftPurchases(records),
	}, nil
}

func parseStorePurchaseHistoryRow(row *html.Node, index int) Model.StorePurchaseHistoryRecord {
	return Model.StorePurchaseHistoryRecord{
		Index:           index,
		TransactionID:   extractTransactionID(htmlquery.SelectAttr(row, "onclick")),
		Date:            nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_date ")]`),
		Item:            firstNonEmptyText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_items ")]//div[contains(@style, "clear")]`, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_items ")]`),
		Receiver:        nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_items ")]//a[contains(@href, "steamcommunity.com/profiles/")]`),
		TransactionType: firstNonEmptyText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_type ")]/div[not(contains(concat(" ", normalize-space(@class), " "), " wth_payment "))]`, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_type ")]`),
		Payment:         nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_type ")]//div[contains(concat(" ", normalize-space(@class), " "), " wth_payment ")]`),
		BasePrice:       nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_base_price ")]`),
		Tax:             nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_tax ")]`),
		Shipping:        nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_shipping ")]`),
		Total:           nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_total ")]`),
		WalletChange:    nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_wallet_change ")]`),
		WalletBalance:   nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_wallet_balance ")]`),
		Refunded:        strings.Contains(htmlquery.SelectAttr(row, "class"), "wht_item_refunded") || htmlquery.FindOne(row, `.//*[contains(concat(" ", normalize-space(@class), " "), " wht_refunded ")]`) != nil,
	}
}

func filterLatestUnrefundedGiftPurchases(records []Model.StorePurchaseHistoryRecord) []Model.StorePurchaseHistoryRecord {
	filtered := make([]Model.StorePurchaseHistoryRecord, 0)
	for _, record := range records {
		if record.TransactionType != "礼物购买" || record.Refunded {
			break
		}
		filtered = append(filtered, record)
	}
	return filtered
}

func extractTransactionID(onclick string) string {
	matches := regexp.MustCompile(`transid=(\d+)`).FindStringSubmatch(onclick)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func firstNonEmptyText(root *html.Node, paths ...string) string {
	for _, path := range paths {
		if text := nodeText(root, path); text != "" {
			return text
		}
	}
	return ""
}

func nodeText(root *html.Node, path string) string {
	node := htmlquery.FindOne(root, path)
	if node == nil {
		return ""
	}
	return cleanText(htmlquery.InnerText(node))
}

func cleanText(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
