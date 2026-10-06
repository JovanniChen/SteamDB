package Dao

import (
	"strings"
	"testing"

	"github.com/antchfx/htmlquery"
)

func TestParseStorePurchaseHistoryRowParsesMultipleGifts(t *testing.T) {
	doc, err := htmlquery.Parse(strings.NewReader(`
		<table><tbody>
			<tr class="wallet_table_row" onclick="location.href='?transid=123'">
				<td class="wht_date">2026 年 8 月 25 日</td>
				<td class="wht_items">
					<div style="clear: both">游戏 A</div>
					<div class="wth_payment">礼物已发送给 <a href="https://steamcommunity.com/profiles/1/">用户 A</a></div>
					<div style="clear: both">游戏 B</div>
					<div class="wth_payment">礼物已发送给 <a href="https://steamcommunity.com/profiles/2/">用户 B</a></div>
				</td>
				<td class="wht_type"><div>礼物购买</div><div class="wth_payment">钱包</div></td>
			</tr>
		</tbody></table>
	`))
	if err != nil {
		t.Fatal(err)
	}

	rows := htmlquery.Find(doc, `//tr[contains(concat(" ", normalize-space(@class), " "), " wallet_table_row ")]`)
	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d", len(rows))
	}

	record := parseStorePurchaseHistoryRow(rows[0], 1)
	if len(record.Items) != 2 || record.Items[0] != "游戏 A" || record.Items[1] != "游戏 B" {
		t.Fatalf("unexpected items: %#v", record.Items)
	}
	if len(record.Receivers) != 2 || record.Receivers[0] != "用户 A" || record.Receivers[1] != "用户 B" {
		t.Fatalf("unexpected receivers: %#v", record.Receivers)
	}
}

func TestParseStorePurchaseDetailURLs(t *testing.T) {
	urls, err := parseStorePurchaseDetailURLs([]byte(`
		<html><body>
			<a href="/zh-cn/wizard/HelpWithMyPurchase?line_item=274313161229376464&amp;transid=274313161229376463">项目 1</a>
			<a href="https://help.steampowered.com/zh-cn/wizard/HelpWithMyPurchase?line_item=274313161229376465&amp;transid=274313161229376463">项目 2</a>
			<a href="/zh-cn/wizard/HelpWithMyPurchase?line_item=274313161229376464&amp;transid=274313161229376463">重复项目</a>
		</body></html>
	`), "https://help.steampowered.com/zh-cn/wizard/HelpWithTransaction?transid=274313161229376463")
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 2 {
		t.Fatalf("expected two unique purchase URLs, got %d: %#v", len(urls), urls)
	}
	if urls[0] != "https://help.steampowered.com/zh-cn/wizard/HelpWithMyPurchase?line_item=274313161229376464&transid=274313161229376463" {
		t.Fatalf("unexpected first purchase URL: %s", urls[0])
	}
}

func TestParseStorePurchaseInventoryAssetKeys(t *testing.T) {
	keys := parseStorePurchaseInventoryAssetKeys([]byte(`
		<div class="purchase_line_items">
			<div>此份礼物已存放在您的 <a href="https://steamcommunity.com/my/inventory/#753_1_274313161229392851">Steam 库存</a>中。</div>
			<div>您将<a href="https://steamcommunity.com/my/inventory/#753_1_274313161229392852">这份礼物</a>发送给 <a href="https://steamcommunity.com//profiles/76561199669923296">ERT234</a>。</div>
		</div>
		<div class="purchase_line_items"><a href="https://steamcommunity.com/my/inventory/#753_1_274313161229392851">重复</a></div>
	`))
	if len(keys) != 2 {
		t.Fatalf("expected two unique inventory keys, got %d: %#v", len(keys), keys)
	}
	if keys[0] != "753_1_274313161229392851" || keys[1] != "753_1_274313161229392852" {
		t.Fatalf("unexpected inventory keys: %#v", keys)
	}
}
