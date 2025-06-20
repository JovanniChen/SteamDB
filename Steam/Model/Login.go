package Model

import (
	"encoding/hex"
	"math/big"
	"strconv"
)

type SteamPublicKey struct {
	PublicKeyExp string `json:"publickey_exp,omitempty"`
	PublicKeyMod string `json:"publickey_mod,omitempty"`
	SteamID      uint64 `json:"steamid,string,omitempty"`
	Success      bool   `json:"success"`
	Timestamp    uint64 `json:"timestamp,string,omitempty"`
	TokenGID     string `json:"token_gid,omitempty"`
}

func (spk SteamPublicKey) Modulus() (*big.Int, error) {
	by, er := hex.DecodeString(spk.PublicKeyMod)
	if er != nil {
		return nil, er
	}
	bi := big.NewInt(0)
	return bi.SetBytes(by), nil
}

func (spk SteamPublicKey) Exponent() (int64, error) {
	return strconv.ParseInt(spk.PublicKeyExp, 16, 0)
}

type message struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type RefreshResponse struct {
	message
	Cookie map[string]string `json:"cookie,omitempty"`
}

type LoginResponse struct {
	message
	Data struct {
		ClientID  string
		RequestID string
		SteamID   string
	}
}

type FinalizeResponse struct {
	message
	SteamID      string `json:"steamID"`
	TransferInfo []struct {
		Url    string `json:"url"`
		Params struct {
			Nonce string `json:"nonce"`
			Auth  string `json:"auth"`
		} `json:"params"`
	} `json:"transfer_info"`
}

type CheckLoginResponse struct {
	message
	Url  string `json:"url"`
	Data struct {
		SteamLoginSecure string `json:"steamLoginSecure"`
		SessionId        string `json:"sessionid"`
	}
}
