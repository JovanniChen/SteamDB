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
