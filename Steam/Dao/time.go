// time.go - Steam服务器时间同步功能
// 用于获取Steam服务器标准时间，确保时间敏感操作的准确性
package Dao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
)

// QueryTime Steam服务器时间查询响应结构体
// 包含Steam服务器返回的时间同步相关参数
type QueryTime struct {
	Response struct {
		ServerTime                        string `json:"server_time"`                           // 服务器当前时间(秒级时间戳)
		SkewToleranceSeconds              string `json:"skew_tolerance_seconds"`                // 允许的时间偏差秒数
		LargeTimeJink                     string `json:"large_time_jink"`                       // 大时间跳跃阈值
		ProbeFrequencySeconds             int    `json:"probe_frequency_seconds"`               // 探测频率(秒)
		AdjustedTimeProbeFrequencySeconds int    `json:"adjusted_time_probe_frequency_seconds"` // 调整后的探测频率
		HintProbeFrequencySeconds         int    `json:"hint_probe_frequency_seconds"`          // 提示探测频率
		SyncTimeout                       int    `json:"sync_timeout"`                          // 同步超时时间
		TryAgainSeconds                   int    `json:"try_again_seconds"`                     // 重试间隔秒数
		MaxAttempts                       int    `json:"max_attempts"`                          // 最大尝试次数
	} `json:"response,omitempty"`
}

// SteamTime 获取Steam服务器标准时间
// 返回经过时间偏差修正的Steam服务器时间
// 返回值：Unix时间戳(秒)和可能的错误
func (d *Dao) SteamTime() (int64, error) {
	// 获取本地时间与Steam服务器的时间偏差
	offset, err := d.timeOffset()
	if err != nil {
		return 0, err
	}

	fmt.Println("**************************")
	fmt.Println("offset:", offset)
	fmt.Println("**************************")

	// 返回修正后的Steam服务器时间
	i := time.Now().Unix() + offset
	return i, nil
}

// timeOffset 计算本地时间与Steam服务器时间的偏差
// 通过调用Steam时间查询API获取服务器时间，并计算偏差
// 返回值：时间偏差(秒)和可能的错误
func (d *Dao) timeOffset() (int64, error) {
	// 创建HTTP请求查询Steam服务器时间
	req, err := d.NewRequest("POST", Constants.QueryTime, nil)
	if err != nil {
		return 0, err
	}

	// 发送请求获取响应
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
	if resp.StatusCode != 200 {
		return 0, Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应体数据
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return 0, err
	}

	// 解析JSON响应
	body := &QueryTime{}
	err = json.Unmarshal(buf.Bytes(), body)
	if err != nil {
		return 0, err
	}

	// 检查服务器时间是否有效
	if body.Response.ServerTime == "" {
		return 0, Errors.Error("ServerTime is empty")
	}

	// 计算时间偏差：服务器时间 - 本地时间
	timeoffset, _ := strconv.ParseInt(body.Response.ServerTime, 10, 64)
	timeoffset -= time.Now().Unix()
	return timeoffset, nil
}
