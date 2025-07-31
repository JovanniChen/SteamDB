// Model包 - Steam积分系统相关的数据模型
// 定义了Steam积分反应系统的数据结构
package Model

// Reaction Steam积分反应结构体
// 代表一个可用的积分反应类型(如点赞、喜爱等)
type Reaction struct {
	ReactionID        int   `json:"reactionid"`         // 反应ID，唯一标识反应类型
	PointsCost        int   `json:"points_cost"`        // 消耗积分数量
	PointsTransferred int   `json:"points_transferred"` // 传输给目标用户的积分数
	ValidTargetTypes  []int `json:"valid_target_types"` // 有效的目标类型数组
	ValidUgcTypes     []int `json:"valid_ugc_types"`    // 有效的UGC(用户生成内容)类型
}

// Response Steam API响应的数据部分
// 包含所有可用的反应类型列表
type Response struct {
	Reactions []Reaction `json:"reactions"` // 可用反应类型数组
}

// JsonData 完整的JSON响应结构体
// Steam积分反应API返回的最外层数据结构
type JsonData struct {
	Response Response `json:"response"` // 包含实际数据的响应对象
}
