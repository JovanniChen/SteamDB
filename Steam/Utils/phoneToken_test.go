package Utils

import (
	"math"
	"testing"
)

func TestLoadMaFileAcceptsSteamIDStringAndNumber(t *testing.T) {
	const expectedSteamID uint64 = 76561198000000000

	tests := []struct {
		name    string
		steamID string
	}{
		{name: "string", steamID: `"76561198000000000"`},
		{name: "number", steamID: `76561198000000000`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phoneToken, err := LoadMaFile(`{"Session":{"steamid":` + tt.steamID + `}}`)
			if err != nil {
				t.Fatalf("LoadMaFile() error = %v", err)
			}

			if got := uint64(phoneToken.MaFile.Session.SteamID); got != expectedSteamID {
				t.Fatalf("SteamID = %d, want %d", got, expectedSteamID)
			}
			if got := phoneToken.MaFile.Session.SteamID.String(); got != "76561198000000000" {
				t.Fatalf("SteamID.String() = %q, want %q", got, "76561198000000000")
			}
		})
	}
}

func TestLoadMaFileRejectsInvalidSteamID(t *testing.T) {
	if _, err := LoadMaFile(`{"Session":{"steamid":"invalid"}}`); err == nil {
		t.Fatal("LoadMaFile() expected an error for an invalid SteamID")
	}
}

func TestPhoneTokenGenerateConfirmationQueryParamsPreservesUint64SteamID(t *testing.T) {
	phoneToken := &PhoneToken{MaFile: &MaFile{DeviceID: "android:test"}}

	params, err := phoneToken.GenerateConfirmationQueryParams(1, math.MaxUint64, "conf")
	if err != nil {
		t.Fatalf("GenerateConfirmationQueryParams() error = %v", err)
	}
	if got := params.GetString("a"); got != "18446744073709551615" {
		t.Fatalf("SteamID query parameter = %q, want %q", got, "18446744073709551615")
	}
}
