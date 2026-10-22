package Dao

import (
	"strings"
	"testing"

	"github.com/antchfx/htmlquery"
)

func TestExtractFinalPriceSupportsEuropeanDecimalComma(t *testing.T) {
	doc, err := htmlquery.Parse(strings.NewReader(`<div class="discount_final_price">12,54€</div>`))
	if err != nil {
		t.Fatal(err)
	}
	price := extractFinalPrice(htmlquery.FindOne(doc, "//div"))
	if price != "12.54" {
		t.Fatalf("expected 12.54, got %q", price)
	}
}

func TestExtractFinalPricePrefersCentsAttribute(t *testing.T) {
	doc, err := htmlquery.Parse(strings.NewReader(`<div class="discount_block" data-price-final="1254"><div class="discount_final_price">12,54€</div></div>`))
	if err != nil {
		t.Fatal(err)
	}
	price := extractFinalPrice(htmlquery.FindOne(doc, "//div"))
	if price != "12.54" {
		t.Fatalf("expected 12.54 from cents attribute, got %q", price)
	}
}
