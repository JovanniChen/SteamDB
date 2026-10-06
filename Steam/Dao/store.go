package Dao

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

	ioutil.WriteFile("store_purchase_history.html", body, 0644)

	result, err := parseStorePurchaseHistory(body)
	if err != nil {
		return nil, err
	}
	if err := d.populateStorePurchaseInventoryAssets(result); err != nil {
		return nil, err
	}
	return result, nil
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
	items := nodeTexts(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_items ")]//div[contains(@style, "clear")]`)
	if len(items) == 0 {
		if item := nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_items ")]`); item != "" {
			items = []string{item}
		}
	}
	receivers := nodeTexts(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_items ")]//a[contains(@href, "steamcommunity.com/profiles/")]`)

	return Model.StorePurchaseHistoryRecord{
		Index:           index,
		TransactionID:   extractTransactionID(htmlquery.SelectAttr(row, "onclick")),
		TransactionURL:  extractTransactionURL(htmlquery.SelectAttr(row, "onclick")),
		Date:            nodeText(row, `.//td[contains(concat(" ", normalize-space(@class), " "), " wht_date ")]`),
		Items:           items,
		Receivers:       receivers,
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

// populateStorePurchaseInventoryAssets follows the detail pages only for the
// records selected by filterLatestUnrefundedGiftPurchases. A transaction can
// contain several line items, so all inventory keys are retained in order.
func (d *Dao) populateStorePurchaseInventoryAssets(result *Model.StorePurchaseHistoryResult) error {
	if result == nil {
		return nil
	}

	transactionAssets := make(map[string][]string)
	for i := range result.LatestUnrefundedGiftPurchases {
		record := &result.LatestUnrefundedGiftPurchases[i]
		transactionURL := record.TransactionURL
		if transactionURL == "" && record.TransactionID != "" {
			transactionURL = Constants.HelpWithTransaction + "?transid=" + url.QueryEscape(record.TransactionID)
		}
		if transactionURL == "" {
			continue
		}

		assets, ok := transactionAssets[transactionURL]
		if !ok {
			transactionBody, err := d.getStorePurchasePage(transactionURL)
			if err != nil {
				return fmt.Errorf("获取交易详情失败 (%s): %w", record.TransactionID, err)
			}
			purchaseURLs, err := parseStorePurchaseDetailURLs(transactionBody, transactionURL)
			if err != nil {
				return fmt.Errorf("解析交易详情失败 (%s): %w", record.TransactionID, err)
			}

			assets = make([]string, 0)
			for _, purchaseURL := range purchaseURLs {
				purchaseBody, err := d.getStorePurchasePage(purchaseURL)
				if err != nil {
					return fmt.Errorf("获取购买详情失败 (%s): %w", purchaseURL, err)
				}
				assets = appendUniqueStrings(assets, parseStorePurchaseInventoryAssetKeys(purchaseBody)...)
			}
			transactionAssets[transactionURL] = assets
		}
		record.InventoryAssetKeys = append([]string(nil), assets...)
	}
	return nil
}

func (d *Dao) getStorePurchasePage(pageURL string) ([]byte, error) {
	req, err := d.Request(http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status = %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func parseStorePurchaseDetailURLs(body []byte, baseURL string) ([]string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, node := range htmlquery.Find(doc, `//a[contains(@href, "HelpWithMyPurchase")]`) {
		href := strings.TrimSpace(htmlquery.SelectAttr(node, "href"))
		if href == "" {
			continue
		}
		resolved, err := resolveStorePurchaseURL(baseURL, href)
		if err != nil {
			continue
		}
		urls = appendUniqueStrings(urls, resolved)
	}
	return urls, nil
}

func parseStorePurchaseInventoryAssetKeys(body []byte) []string {
	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil
	}

	keys := make([]string, 0)
	for _, node := range htmlquery.Find(doc, `//div[contains(concat(" ", normalize-space(@class), " "), " purchase_line_items ")]//a[contains(@href, "/inventory/#")]`) {
		href := strings.TrimSpace(htmlquery.SelectAttr(node, "href"))
		parsed, err := url.Parse(href)
		if err != nil || parsed.Fragment == "" {
			continue
		}
		if regexp.MustCompile(`^\d+_\d+_\d+$`).MatchString(parsed.Fragment) {
			keys = appendUniqueStrings(keys, parsed.Fragment)
		}
	}
	return keys
}

func resolveStorePurchaseURL(baseURL, href string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	target, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(target).String(), nil
}

func appendUniqueStrings(values []string, additions ...string) []string {
	seen := make(map[string]struct{}, len(values)+len(additions))
	for _, value := range values {
		seen[value] = struct{}{}
	}
	for _, value := range additions {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		values = append(values, value)
	}
	return values
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

func extractTransactionURL(onclick string) string {
	matches := regexp.MustCompile(`["']([^"']*HelpWithTransaction\?transid=\d+)["']`).FindStringSubmatch(onclick)
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

func nodeTexts(root *html.Node, path string) []string {
	nodes := htmlquery.Find(root, path)
	texts := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if text := cleanText(htmlquery.InnerText(node)); text != "" {
			texts = append(texts, text)
		}
	}
	return texts
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
