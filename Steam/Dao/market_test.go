package Dao

import (
	"testing"

	"github.com/JovanniChen/SteamDB/Steam/Model"
)

func TestExtractReceiverName(t *testing.T) {
	ownerDescriptions := []Model.OwnerDescription{
		{
			Type:  "html",
			Value: `已发送给 <a href="https://steamcommunity.com/profiles/76561199533675524" data-miniprofile="1573409796">4wzwg</a>`,
		},
	}

	if got := extractReceiverName(ownerDescriptions); got != "4wzwg" {
		t.Fatalf("expected receiver name 4wzwg, got %q", got)
	}
}

func TestExtractReceiverNameWithEstimatedTime(t *testing.T) {
	ownerDescriptions := []Model.OwnerDescription{
		{
			Type:  "html",
			Value: `预计在9 月 1 日 下午 10:25发送给 <a href="https://steamcommunity.com/profiles/76561198699275203" data-miniprofile="739009475">Why always me</a>`,
		},
	}

	if got := extractReceiverName(ownerDescriptions); got != "Why always me" {
		t.Fatalf("expected receiver name Why always me, got %q", got)
	}
}

func TestExtractDescriptionName(t *testing.T) {
	descriptions := []Model.DescriptionText{
		{Value: "这是礼物的描述内容"},
	}

	if got := extractDescriptionName(descriptions); got != "这是礼物的描述内容" {
		t.Fatalf("expected description name, got %q", got)
	}
}

func TestExtractDescriptionNameSkipsEmptyValues(t *testing.T) {
	descriptions := []Model.DescriptionText{
		{Value: "  "},
		{Value: "第二条描述"},
	}

	if got := extractDescriptionName(descriptions); got != "第二条描述" {
		t.Fatalf("expected second description value, got %q", got)
	}
}

func TestBuildSteamGiftItemsIncludesEnglishMarketNames(t *testing.T) {
	chineseResponse := Model.InventoryResponse{
		Assets: []Model.Asset{
			{AssetID: "asset-1", ClassID: "class-1", InstanceID: "instance-1"},
		},
		Descriptions: []Model.Description{
			{
				ClassID:        "class-1",
				InstanceID:     "instance-1",
				Name:           "中文名称",
				MarketName:     "中文市场名称",
				MarketHashName: "中文 Hash 名称",
			},
		},
	}
	englishDescriptionMap := map[string]Model.Description{
		"class-1_instance-1": {
			ClassID:        "class-1",
			InstanceID:     "instance-1",
			MarketName:     "English Market Name",
			MarketHashName: "English Market Hash Name",
		},
	}

	items := buildSteamGiftItems(chineseResponse, englishDescriptionMap)
	if len(items) != 1 {
		t.Fatalf("expected one item, got %d", len(items))
	}
	if items[0].MarketName != "中文市场名称" || items[0].MarketHashName != "中文 Hash 名称" {
		t.Fatalf("expected Chinese market fields to be preserved, got %+v", items[0])
	}
	if items[0].EnglishMarketName != "English Market Name" {
		t.Fatalf("expected English market name, got %q", items[0].EnglishMarketName)
	}
	if items[0].EnglishMarketHashName != "English Market Hash Name" {
		t.Fatalf("expected English market hash name, got %q", items[0].EnglishMarketHashName)
	}
}
