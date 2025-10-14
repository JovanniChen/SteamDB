// SteamDB Go客户端库 - 主要客户端接口
// 提供简单易用的Steam平台交互功能，包括登录、积分系统、Steam Guard等
package Steam

import (
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Model"
	"github.com/JovanniChen/SteamDB/Steam/Protoc"
)

// Client Steam客户端结构体
// 封装了Steam平台的所有交互功能，提供统一的API接口
type Client struct {
	dao *Dao.Dao // 底层数据访问对象
}

// Config 客户端配置选项
// 用于初始化Steam客户端时的配置设置
type Config struct {
	Proxy   string        // 代理服务器地址，格式: "host:port"，空字符串表示不使用代理
	Timeout time.Duration // 请求超时时间，0表示使用默认值
}

// DefaultConfig 返回默认配置
// 使用推荐的默认设置初始化客户端配置
func DefaultConfig() *Config {
	return &Config{
		Proxy:   "",               // 不使用代理
		Timeout: 30 * time.Second, // 30秒超时
	}
}

func NewConfig(proxy string) *Config {
	return &Config{
		Proxy:   proxy,
		Timeout: 30 * time.Second,
	}
}

// NewClient 创建新的Steam客户端实例
// 参数:
//
//	config - 客户端配置，可以为nil使用默认配置
//
// 返回值:
//
//	*Client - Steam客户端实例
//	error - 初始化错误
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 创建底层DAO对象
	dao := Dao.New(config.Proxy)

	return &Client{
		dao: dao,
	}, nil
}

// LoginCredentials 登录凭据结构体
// 包含用户的基本登录信息
type LoginCredentials struct {
	Username     string // Steam用户名
	Password     string // Steam密码
	SharedSecret string // Steam Guard共享密钥(base64编码)，如果没有2FA可以为空
	MaFile       string
}

// UserInfo 用户信息结构体
// 包含登录成功后的用户详细信息
type UserInfo struct {
	SteamID      uint64 // Steam ID
	Username     string // 用户名
	Nickname     string // 昵称
	AccessToken  string // 访问令牌
	RefreshToken string // 刷新令牌
	CountryCode  string // 国家代码
	MaFile       string
}

// Login 执行Steam登录
// 使用用户凭据登录Steam平台，支持Steam Guard双因素认证
// 参数:
//
//	credentials - 登录凭据信息
//
// 返回值:
//
//	*UserInfo - 用户信息
//	error - 登录错误
func (c *Client) Login(credentials *LoginCredentials) (*UserInfo, error) {
	// 执行登录
	err := c.dao.Login(credentials.Username, credentials.Password, credentials.SharedSecret)
	if err != nil {
		return nil, err
	}

	if err := c.dao.UserInfo(); err != nil {
		return nil, err
	}

	// 在登录完成后立即获取所有用户信息，确保一致性
	steamID := c.GetSteamID()
	nickname := c.GetNickname()
	refreshToken := c.GetRefreshToken()
	countryCode := c.GetCountryCode()

	// 获取访问令牌
	accessToken, err := c.dao.AccessToken()
	if err != nil {
		return nil, err
	}

	// 使用同一时间点获取的一致数据创建UserInfo
	userInfo := &UserInfo{
		SteamID:      steamID,
		Username:     credentials.Username,
		Nickname:     nickname,
		RefreshToken: refreshToken,
		CountryCode:  countryCode,
		AccessToken:  accessToken,
		MaFile:       credentials.MaFile,
	}

	return userInfo, nil
}

// GetTokenCode 获取Steam Guard令牌代码
// 基于共享密钥生成当前时间的6位数字验证码
// 参数:
//
//	sharedSecret - Steam Guard共享密钥(base64编码)
//
// 返回值:
//
//	string - 6位数字验证码
//	error - 生成错误
func (c *Client) GetTokenCode(sharedSecret string) (string, error) {
	return c.dao.GetTokenCode(sharedSecret)
}

// PointsSummary 积分摘要信息
// 包含用户的Steam积分系统详细信息
type PointsSummary struct {
	SteamID uint64 // Steam ID
	Points  int64  // 当前积分数量
	Level   int32  // 用户等级
}

// GetPointsSummary 获取用户积分摘要
// 查询指定用户的Steam积分系统信息
// 参数:
//
//	steamID - 目标用户的Steam ID
//
// 返回值:
//
//	*PointsSummary - 积分摘要信息
//	error - 查询错误
func (c *Client) GetPointsSummary(steamID uint64) (*PointsSummary, error) {
	summaryData, err := c.dao.GetSummary(steamID)
	if err != nil {
		return nil, err
	}

	return &PointsSummary{
		SteamID: steamID,
		Points:  int64(summaryData.GetSummary().GetPoints()),
		Level:   0, // SummaryPoint结构中没有Level字段，设为0
	}, nil
}

// ReactionConfig 反应配置信息
// 描述可用的反应类型和消耗积分
type ReactionConfig struct {
	ReactionID       uint32   // 反应ID
	PointsCost       int64    // 消耗积分数量
	ValidTargetTypes []uint32 // 可用的目标类型列表
}

// GetReactionConfig 获取反应配置
// 查询当前可用的所有反应类型及其配置信息
// 返回值:
//
//	[]ReactionConfig - 反应配置列表
//	error - 查询错误
func (c *Client) GetReactionConfig() ([]ReactionConfig, error) {
	configData, err := c.dao.GetReactionConfig()
	if err != nil {
		return nil, err
	}

	var reactions []ReactionConfig
	for _, reaction := range configData.Response.Reactions {
		// 转换数据类型
		var validTargetTypes []uint32
		for _, t := range reaction.ValidTargetTypes {
			validTargetTypes = append(validTargetTypes, uint32(t))
		}

		reactions = append(reactions, ReactionConfig{
			ReactionID:       uint32(reaction.ReactionID),
			PointsCost:       int64(reaction.PointsCost),
			ValidTargetTypes: validTargetTypes,
		})
	}

	return reactions, nil
}

// AddReactionResult 添加反应结果
// 包含添加反应操作的执行结果
type AddReactionResult struct {
	Success        bool  // 操作是否成功
	PointsConsumed int64 // 消耗的积分数量
}

// AddReaction 为指定用户添加反应
// 使用积分为目标用户添加指定类型的反应
// 参数:
//
//	targetSteamID - 目标用户的Steam ID
//	reactionType - 反应类型ID
//	reactionID - 具体反应ID
//
// 返回值:
//
//	*AddReactionResult - 添加结果
//	error - 操作错误
func (c *Client) AddReaction(targetSteamID uint64, reactionType uint32, reactionID uint32, pointsCost int64) (*AddReactionResult, error) {
	pointsRemainning, err := c.GetPointsSummary(c.GetSteamID())
	if err != nil {
		return nil, err
	}

	_, err = c.dao.AddReaction(targetSteamID, int32(reactionType), reactionID)
	if err != nil {
		return nil, err
	}

	pointsNow, err := c.GetPointsSummary(c.GetSteamID())
	if err != nil {
		return nil, err
	}

	return &AddReactionResult{
		Success:        pointsRemainning.Points-pointsNow.Points == pointsCost, // 简单判断：非nil表示成功
		PointsConsumed: pointsRemainning.Points - pointsNow.Points,             // 暂时无法从结果中获取具体消耗积分数
	}, nil
}

// GetReactions 获取用户的反应记录
// 查询指定用户收到的所有反应记录
// 参数:
//
//	steamID - 目标用户的Steam ID
//	reactionType - 反应类型过滤器，0表示获取所有类型
//
// 返回值:
//
//	*Protoc.ReactionsReceive - 反应记录数据
//	error - 查询错误
func (c *Client) GetReactions(steamID uint64, reactionType uint32) (*Protoc.ReactionsReceive, error) {
	return c.dao.GetReacionts(steamID, int32(reactionType))
}

func (c *Client) GetSteamIDByFriendLink(friendLink string) (uint64, error) {
	return c.dao.GetSteamIDByFriendLink(friendLink)
}

func (c *Client) GetInventory(gameID int, categoryId int) ([]Model.Item, error) {
	return c.dao.GetInventory(gameID, categoryId)
}

func (c *Client) PutList(assetID string, price float64, currency int, maFileContent string) error {
	return c.dao.PutList(assetID, price, currency, maFileContent)
}

func (c *Client) BuyListing(creatorId string, name string, buyerPrice float64, sellerReceivePrice float64, maFileContent string) error {
	return c.dao.BuyListing(creatorId, name, buyerPrice, sellerReceivePrice, "0", maFileContent)
}

func (c *Client) CreateOrder(marketHashName string, price float64, quantity int64, maFileContent string) error {
	return c.dao.CreateOrder(marketHashName, price, quantity, maFileContent)
}

func (c *Client) RemoveMyListings(creatorId string) error {
	return c.dao.RemoveMyListings(creatorId)
}

func (c *Client) RemoveAllMyListings() error {
	return c.dao.RemoveAllMyListings()
}

func (c *Client) GetMyListings() ([]Model.MyListingReponse, error) {
	return c.dao.GetMyListings()
}

func (c *Client) GetConfirmations(maFileContent string) error {
	return c.dao.GetAndOperateConfirmations("allow", maFileContent)
}

// CheckLoginStatus 检查登录状态
// 验证当前客户端是否已成功登录Steam
// 参数:
//
//	url - 要检查的Steam URL（用于确定检查哪个域名的登录状态）
//
// 返回值:
//
//	bool - true表示已登录，false表示未登录
func (c *Client) CheckLoginStatus(url string) bool {
	return c.dao.CheckLogin(url)
}

// SetLanguage 设置Steam语言偏好
// 更改用户的Steam界面显示语言
// 参数:
//
//	language - 语言代码（如: "english", "schinese", "japanese"等）
//
// 返回值:
//
//	error - 设置错误
func (c *Client) SetLanguage(language string) error {
	return c.dao.SetLanguage(language)
}

// GetUserInfo 获取详细用户信息
// 查询当前登录用户的详细资料信息
// 返回值:
//
//	error - 查询错误
//
// 注意: 此方法会更新内部用户信息状态，具体数据需要通过其他getter方法获取
func (c *Client) GetUserInfo() error {
	return c.dao.UserInfo()
}

// 以下是一些便捷的getter方法，用于获取用户信息

// GetUsername 获取当前登录用户的用户名
func (c *Client) GetUsername() string {
	return c.dao.GetUsername()
}

// GetSteamID 获取当前登录用户的Steam ID
func (c *Client) GetSteamID() uint64 {
	return c.dao.GetSteamID()
}

// GetAccessToken 获取当前有效的访问令牌
func (c *Client) GetAccessToken() (string, error) {
	return c.dao.AccessToken()
}

// GetRefreshToken 获取刷新令牌
func (c *Client) GetRefreshToken() string {
	return c.dao.GetRefreshToken()
}

// GetNickname 获取用户昵称
func (c *Client) GetNickname() string {
	return c.dao.GetNickname()
}

// GetCountryCode 获取用户国家代码
func (c *Client) GetCountryCode() string {
	return c.dao.GetCountryCode()
}

// GetLoginCookies 获取登录Cookie信息
func (c *Client) GetLoginCookies() map[string]*Dao.LoginCookie {
	return c.dao.GetLoginCookies()
}

func (c *Client) GetBalance() int {
	return c.dao.GetBalance()
}

func (c *Client) GetWaitBalance() int {
	return c.dao.GetWaitBalance()
}

// SetLoginInfo 设置登录信息（用于恢复会话）
func (c *Client) SetLoginInfo(username string, steamID uint64, nickname string, countryCode string, accessToken string, refreshToken string, loginCookies map[string]*Dao.LoginCookie) {
	c.dao.SetLoginInfoDirect(username, steamID, nickname, countryCode, accessToken, refreshToken, loginCookies)
}
