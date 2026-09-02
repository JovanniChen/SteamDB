package Dao

import (
	"testing"

	"github.com/JovanniChen/SteamDB/Steam/Model"
)

func TestExtractEstimatedSendTime(t *testing.T) {
	ownerDescriptions := []Model.OwnerDescription{
		{
			Type:  "bbcode",
			Value: "预计在[date]1788319301[/date]发送给 [persona]739009475[/persona]",
		},
	}

	if got := extractEstimatedSendTime(ownerDescriptions); got != 1788319301 {
		t.Fatalf("expected estimated send time 1788319301, got %d", got)
	}
}

func TestExtractEstimatedSendTimeWithoutDate(t *testing.T) {
	ownerDescriptions := []Model.OwnerDescription{
		{Type: "bbcode", Value: "礼物已发送"},
	}

	if got := extractEstimatedSendTime(ownerDescriptions); got != 0 {
		t.Fatalf("expected no estimated send time, got %d", got)
	}
}
