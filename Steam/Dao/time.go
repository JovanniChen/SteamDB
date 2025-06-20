package Dao

import (
	"bytes"
	"encoding/json"
	"example.com/m/v2/Steam"
	"example.com/m/v2/Steam/Errors"
	"strconv"
	"time"
)

type QueryTime struct {
	Response struct {
		ServerTime                        string `json:"server_time"`
		SkewToleranceSeconds              string `json:"skew_tolerance_seconds"`
		LargeTimeJink                     string `json:"large_time_jink"`
		ProbeFrequencySeconds             int    `json:"probe_frequency_seconds"`
		AdjustedTimeProbeFrequencySeconds int    `json:"adjusted_time_probe_frequency_seconds"`
		HintProbeFrequencySeconds         int    `json:"hint_probe_frequency_seconds"`
		SyncTimeout                       int    `json:"sync_timeout"`
		TryAgainSeconds                   int    `json:"try_again_seconds"`
		MaxAttempts                       int    `json:"max_attempts"`
	} `json:"response,omitempty"`
}

func (d *Dao) SteamTime() (int64, error) {
	offset, err := d.timeOffset()
	if err != nil {
		return 0, err
	}
	i := time.Now().Unix() + offset
	return i, nil

}

func (d *Dao) timeOffset() (int64, error) {
	req, err := d.NewRequest("POST", Steam.QueryTime, nil)
	if err != nil {
		return 0, err
	}
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, Errors.ResponseError(resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return 0, err
	}
	body := &QueryTime{}
	err = json.Unmarshal(buf.Bytes(), body)
	if err != nil {
		return 0, err
	}
	if body.Response.ServerTime == "" {
		return 0, Errors.Error("ServerTime is empty")
	}
	timeoffset, _ := strconv.ParseInt(body.Response.ServerTime, 10, 64)
	timeoffset -= time.Now().Unix()
	return timeoffset, nil
}
