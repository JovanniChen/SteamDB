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

	"example.com/m/v2/Steam"
	"example.com/m/v2/Steam/Errors"
	"example.com/m/v2/Steam/Model"
	"example.com/m/v2/Steam/Param"
	"example.com/m/v2/Steam/Protoc"
	"google.golang.org/protobuf/proto"
)

// GetReacionts 获取指定目标的积分反应记录
// 查询某个用户或内容的积分反应历史记录
// 参数：
//   targetId - 目标ID(用户SteamID或内容ID)
//   targetType - 目标类型(1=用户档案, 2=用户生成内容等)
// 返回值：操作成功返回nil，失败返回错误
func (d *Dao) GetReacionts(targetId uint64, targetType int32) error {
	// 构建请求数据
	reactionsSend := &Protoc.ReactionsSend{
		Targetid:   targetId,
		TargetType: targetType,
	}

	// 序列化为protobuf格式
	data, err := proto.Marshal(reactionsSend)
	if err != nil {
		return err
	}

	// 获取访问令牌
	accessToken, _ := d.AccessToken()
	
	// 构建请求参数
	params := Param.Params{}
	params.SetString("access_token", accessToken)
	params.SetString("origin", Steam.CommunityOrigin)
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Steam.GetReactions+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}
	
	// 发送请求
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	// 解析protobuf响应
	reactionsReceive := &Protoc.ReactionsReceive{}
	err = proto.Unmarshal(buf.Bytes(), reactionsReceive)
	if err != nil {
		return err
	}

	// 输出反应数据(调试用)
	fmt.Printf("%v\n", reactionsReceive)
	return nil
}

// GetSummary 获取指定用户的积分摘要信息
// 查询用户的积分余额、等级等摘要数据
// 参数：steamId - Steam用户ID
// 返回值：操作成功返回nil，失败返回错误
func (d *Dao) GetSummary(steamId uint64) error {
	// 构建请求数据
	summarySend := &Protoc.SummarySend{
		Steamid: steamId,
	}

	// 序列化为protobuf格式
	data, err := proto.Marshal(summarySend)
	if err != nil {
		return err
	}

	// 获取访问令牌
	accessToken, _ := d.AccessToken()
	
	// 构建请求参数
	params := Param.Params{}
	params.SetString("access_token", accessToken)
	params.SetString("origin", Steam.CommunityOrigin)
	params.SetString("spoof_steamid", "")  // 伪装SteamID(通常为空)
	params.SetString("input_protobuf_encoded", base64.StdEncoding.EncodeToString(data))

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Steam.GetSummary+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Add("origin", Steam.CommunityOrigin)
	req.Header.Set("referer", Steam.CommunityOrigin+"/")

	// 发送请求
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	// 解析protobuf响应
	summaryReceive := &Protoc.SummaryReceive{}
	if err := proto.Unmarshal(buf.Bytes(), summaryReceive); err != nil {
		return err
	}

	// 输出摘要数据(调试用)
	fmt.Printf("%+v\n", summaryReceive)

	return nil
}

// GetReactionConfig 获取积分反应配置信息
// 查询当前可用的积分反应类型和消耗配置
// 返回值：操作成功返回nil，失败返回错误
func (d *Dao) GetReactionConfig() error {
	// 构建请求参数
	params := Param.Params{}
	params.SetString("origin", Steam.CommunityOrigin)
	params.SetString("input_protobuf_encoded", "") // 配置查询无需额外数据

	// 输出请求URL(调试用)
	fmt.Println(Steam.GetReactionConfig + "?" + params.ToUrl())

	// 创建GET请求
	req, err := d.NewRequest(http.MethodGet, Steam.GetReactionConfig+"?"+params.ToUrl(), nil)
	if err != nil {
		return err
	}
	
	// 设置请求头
	req.Header.Add("origin", Steam.CommunityOrigin)
	req.Header.Set("referer", Steam.CommunityOrigin+"/")

	// 发送请求
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return Errors.ResponseError(resp.StatusCode)
	}

	// 读取响应数据
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	// 输出原始响应(调试用)
	fmt.Println(buf.String())

	// 解析JSON响应
	var data Model.JsonData
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return err
	}

	// 处理反应配置数据(此处已注释)
	// for _, v := range data.Response.Reactions {
	// 	fmt.Printf("+%v\n", v)
	// }

	return nil
}

// AddReaction 为指定目标添加积分反应
// 向目标用户或内容添加反应(如点赞、获奖等)，消耗对应积分
// 参数：
//   targetId - 目标ID(用户SteamID或内容ID)
//   targetType - 目标类型(1=用户档案, 2=用户生成内容等)
//   reactionId - 反应类型ID(参考GetReactionConfig获取)
// 返回值：操作成功返回nil，失败返回错误
func (d *Dao) AddReaction(targetId uint64, targetType int32, reactionId uint32) error {
	// 构建请求数据
	addReactionSend := &Protoc.AddReactionSend{
		Targetid:   targetId,
		TargetType: targetType,
		Reactionid: reactionId,
	}

	// 序列化为protobuf格式
	data, err := proto.Marshal(addReactionSend)
	if err != nil {
		return err
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
	req, err := d.NewRequest(http.MethodPost, Steam.AddReaction+"?"+params.ToUrl(), strings.NewReader(params1.Encode()))
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := d.RetryRequest(Steam.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode != 200 {
		return Errors.ResponseError(resp.StatusCode)
	}

	return nil
}
