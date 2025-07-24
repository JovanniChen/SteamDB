package Dao

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"example.com/m/v2/Steam"
	"example.com/m/v2/Steam/Errors"
	"example.com/m/v2/Steam/Param"
	"example.com/m/v2/Steam/Protoc"
	"google.golang.org/protobuf/proto"
)

func (d *Dao) GetReacionts(targetId uint64, targetType int32) error {
	reactionsSend := &Protoc.ReactionsSend{
		Targetid:   targetId,
		TargetType: targetType,
	}

	data, err := proto.Marshal(reactionsSend)
	if err != nil {
		return err
	}

	accessToken, _ := d.AccessToken()
	params := Param.Params{}
	params.SetString("access_token", accessToken)
	params.SetString("origin", Steam.CommunityOrigin)
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	req, err := d.NewRequest(http.MethodGet, Steam.GetReactions+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	reactionsReceive := &Protoc.ReactionsReceive{}
	err = proto.Unmarshal(buf.Bytes(), reactionsReceive)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) GetSummary(steamId uint64) error {
	summarySend := &Protoc.SummarySend{
		Steamid: steamId,
	}

	data, err := proto.Marshal(summarySend)
	if err != nil {
		return err
	}

	accessToken, _ := d.AccessToken()
	params := Param.Params{}
	params.SetString("access_token", accessToken)
	params.SetString("origin", Steam.CommunityOrigin)
	params.SetString("spoof_steamid", "")
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	req, err := d.NewRequest(http.MethodGet, Steam.GetSummary+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("origin", Steam.CommunityOrigin)
	req.Header.Set("referer", Steam.CommunityOrigin+"/")

	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	summaryReceive := &Protoc.SummaryReceive{}
	if err := proto.Unmarshal(buf.Bytes(), summaryReceive); err != nil {
		return err
	}

	fmt.Printf("%+v\n", summaryReceive)

	return nil
}
