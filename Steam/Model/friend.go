package Model

// FriendLinkStatus 好友链接状态类型
type FriendLinkStatus int

const (
	// FriendLinkStatusFailed 链接解析失败
	FriendLinkStatusFailed FriendLinkStatus = iota
	// FriendLinkStatusOnlyInfo 只获取到信息（链接可能过期）
	FriendLinkStatusOnlyInfo
	// FriendLinkStatusIsFriend 已经是好友
	FriendLinkStatusIsFriend
	// FriendLinkStatusSuccess 链接有效，可以添加好友
	FriendLinkStatusSuccess
)

// FriendInfo 好友信息结构体
type FriendInfo struct {
	PersonName string `json:"person_name"` // 用户昵称
	FriendCode uint64 `json:"friend_code"` // 好友代码（miniprofile ID）
	SessionID  string `json:"session_id"`  // 会话ID
	AbuseID    string `json:"abuse_id"`    // 滥用举报ID
}

// FriendLinkParseResult 好友链接解析结果
type FriendLinkParseResult struct {
	Status FriendLinkStatus `json:"status"` // 链接状态
	Msg    string           `json:"msg"`    // 状态消息
	Data   *FriendInfo      `json:"data"`   // 好友信息数据
}

type AddFriendByLinkResult struct {
	Success int `json:"success"`
}

// AddFriendByCodeResult 通过好友代码添加好友的响应结构
type AddFriendByCodeResult struct {
	Success      bool   `json:"success"`       // 是否成功
	ErrorText    string `json:"error_text"`    // 错误信息
	Failed       int    `json:"failed"`        // 失败状态码
	PlayerName   string `json:"player_name"`   // 玩家名称
	Invited      []any  `json:"invited"`       // 邀请列表
	FailedInvite []any  `json:"failed_invite"` // 失败的邀请列表
}
