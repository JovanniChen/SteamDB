package Dao

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/JovanniChen/SteamDB/Steam/Model"
)

type marketListingsRoundTripFunc func(*http.Request) (*http.Response, error)

func (f marketListingsRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetNewMarketListingsOptionalPriceFilter(t *testing.T) {
	min, max := 3440, 3450
	tests := []struct {
		name   string
		filter []Model.MarketListingPriceFilter
		want   map[string]int
	}{
		{name: "no price filter"},
		{name: "currency only", filter: []Model.MarketListingPriceFilter{{Currency: 23}}, want: map[string]int{"eCurrency": 23}},
		{name: "currency and range", filter: []Model.MarketListingPriceFilter{{Currency: 23, Min: &min, Max: &max}}, want: map[string]int{"eCurrency": 23, "unMin": 3440, "unMax": 3450}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rawQP string
			d := &Dao{
				credentials: &Credentials{LoginCookies: map[string]*LoginCookie{
					"steamcommunity.com": {SessionId: "session"},
				}},
				httpCli: &http.Client{Transport: marketListingsRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					rawQP = req.URL.Query().Get("qp")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"data":{"listings":[]}}`)),
						Header:     make(http.Header),
						Request:    req,
					}, nil
				})},
			}

			if _, err := d.GetNewMarketListings(440, "Item", 0, tt.filter...); err != nil {
				t.Fatalf("GetNewMarketListings returned error: %v", err)
			}

			var query []map[string]json.RawMessage
			if err := json.Unmarshal([]byte(rawQP), &query); err != nil {
				t.Fatalf("invalid qp JSON %q: %v", rawQP, err)
			}
			var price map[string]int
			rawPrice, hasPrice := query[0]["price"]
			if len(tt.want) == 0 {
				if hasPrice {
					t.Fatalf("expected no price filter, got %s", rawPrice)
				}
				return
			}
			if !hasPrice {
				t.Fatal("expected price filter")
			}
			if err := json.Unmarshal(rawPrice, &price); err != nil {
				t.Fatalf("invalid price filter %s: %v", rawPrice, err)
			}
			if len(price) != len(tt.want) {
				t.Fatalf("price filter = %#v, want %#v", price, tt.want)
			}
			for key, want := range tt.want {
				if price[key] != want {
					t.Fatalf("price filter %s = %d, want %d", key, price[key], want)
				}
			}
		})
	}
}
