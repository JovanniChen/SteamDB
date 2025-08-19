// point.go - Steam积分系统相关功能
// 实现Steam积分(Points)的获取、消费和反应管理功能
package Dao

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Protoc"
	"google.golang.org/protobuf/proto"
)

// GetReacionts 获取指定目标的积分反应记录
// 查询某个用户或内容的积分反应历史记录
// 参数：
//
//	targetId - 目标ID(用户SteamID或内容ID)
//	targetType - 目标类型(1=用户档案, 2=用户生成内容等)
//
// 返回值：反应记录数据和可能的错误
func (d *Dao) GetReacionts(targetId uint64, targetType int32) (*Protoc.ReactionsReceive, error) {
	// 构建请求数据
	reactionsSend := &Protoc.ReactionsSend{
		Targetid:   targetId,
		TargetType: targetType,
	}

	// 序列化为protobuf格式
	data, err := proto.Marshal(reactionsSend)
	if err != nil {
		return nil, err
	}

	// 获取访问令牌
	accessToken, _ := d.AccessToken()

	// 构建请求参数
	params := Param.Params{}
	params.SetString("access_token", accessToken)
	params.SetString("origin", Constants.CommunityOrigin)
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Constants.GetReactions+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	// 解析protobuf响应
	reactionsReceive := &Protoc.ReactionsReceive{}

	// 使用重试机制解析protobuf
	if err = protoUnmarshalWithRetry(buf.Bytes(), reactionsReceive, "GetReacionts", 3); err != nil {
		return nil, err
	}

	// 输出反应数据(调试用)
	return reactionsReceive, nil
}

// GetSummary 获取指定用户的积分摘要信息
// 查询用户的积分余额、等级等摘要数据
// 参数：steamId - Steam用户ID
// 返回值：积分摘要数据和可能的错误
func (d *Dao) GetSummary(steamId uint64) (*Protoc.SummaryReceive, error) {
	// 构建请求数据
	summarySend := &Protoc.SummarySend{
		Steamid: steamId,
	}

	// 序列化为protobuf格式
	data, err := proto.Marshal(summarySend)
	if err != nil {
		return nil, err
	}

	// 获取访问令牌
	accessToken, _ := d.AccessToken()

	// 构建请求参数
	params := Param.Params{}
	params.SetString("access_token", accessToken)
	params.SetString("origin", Constants.CommunityOrigin)
	params.SetString("spoof_steamid", "") // 伪装SteamID(通常为空)
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Constants.GetSummary+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", Constants.CommunityOrigin+"/")

	// 发送请求
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	// 解析protobuf响应
	summaryReceive := &Protoc.SummaryReceive{}

	// 使用重试机制解析protobuf
	if err := protoUnmarshalWithRetry(buf.Bytes(), summaryReceive, "GetSummary", 3); err != nil {
		return nil, err
	}

	// 输出摘要数据(调试用)
	fmt.Printf("%+v\n", summaryReceive)

	return summaryReceive, nil
}

// GetReactionConfig 获取积分反应配置信息
// 查询当前可用的积分反应类型和消耗配置
// 返回值：反应配置数据和可能的错误
func (d *Dao) GetReactionConfig() (*Model.JsonData, error) {
	// 构建请求参数
	params := Param.Params{}
	params.SetString("origin", Constants.CommunityOrigin)
	params.SetString("input_protobuf_encoded", "") // 配置查询无需额外数据

	// 输出请求URL(调试用)
	fmt.Println(Constants.GetReactionConfig + "?" + params.ToUrl())

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Constants.GetReactionConfig+"?"+params.ToUrl(), nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Add("origin", Constants.CommunityOrigin)
	req.Header.Set("referer", Constants.CommunityOrigin+"/")

	// 发送请求
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	// 输出原始响应(调试用)
	fmt.Println(buf.String())

	// 解析JSON响应
	var data Model.JsonData
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return nil, err
	}

	// 处理反应配置数据(此处已注释)
	// for _, v := range data.Response.Reactions {
	// 	fmt.Printf("+%v\n", v)
	// }

	return &data, nil
}

// AddReaction 为指定目标添加积分反应
// 向目标用户或内容添加反应(如点赞、获奖等)，消耗对应积分
// 参数：
//
//	targetId - 目标ID(用户SteamID或内容ID)
//	targetType - 目标类型(1=用户档案, 2=用户生成内容等)
//	reactionId - 反应类型ID(参考GetReactionConfig获取)
//
// 返回值：反应结果数据和可能的错误
func (d *Dao) AddReaction(targetId uint64, targetType int32, reactionId uint32) (*Protoc.AddReactionReceive, error) {
	// 构建请求数据
	addReactionSend := &Protoc.AddReactionSend{
		Targetid:   targetId,
		TargetType: targetType,
		Reactionid: reactionId,
	}

	// 序列化为protobuf格式
	data, err := proto.Marshal(addReactionSend)
	if err != nil {
		return nil, err
	}

	// 获取访问令牌
	accessToken, _ := d.AccessToken()

	// 构建URL查询参数(包含访问令牌)
	params := Param.Params{}
	params.SetString("access_token", accessToken)

	// 构建POST请求体参数(包含protobuf数据)
	params1 := Param.Params{}
	params1.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	// 创建POST请求
	req, err := d.NewRequest(http.MethodPost, Constants.AddReaction+"?"+params.ToUrl(), strings.NewReader(params1.Encode()))
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != 200 {
		return nil, Errors.ResponseError(resp.StatusCode)
	}

	// // 读取响应数据
	// buf := new(bytes.Buffer)
	// if _, err := buf.ReadFrom(resp.Body); err != nil {
	// 	return nil, err
	// }

	// // 解析protobuf响应
	// addReactionReceive := &Protoc.AddReactionReceive{}

	// // 使用重试机制解析protobuf
	// if err := protoUnmarshalWithRetry(buf.Bytes(), addReactionReceive, "AddReaction", 3); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
