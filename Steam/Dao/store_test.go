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
