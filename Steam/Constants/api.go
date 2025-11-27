package Constants

// domain 结构体定义Steam平台的各个域名
// 用于构建不同服务的完整URL地址
type domain struct {
	Store     string // Steam商店域名
	Community string // Steam社区域名
	Help      string // Steam帮助中心域名
	CheckOut  string // Steam结账服务域名
	Tv        string // Steam TV域名
	Friend    string // Steam好友服务域名
	Api       string // Steam API服务域名
	Login     string // Steam登录服务域名
}

// 系统级常量定义
const (
	Scheme    string = "https://"                                                                                                            // HTTPS协议前缀
	Tries     int    = 3                                                                                                                     // 默认重试次数
	UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36" // 浏览器用户代理字符串，用于模拟真实浏览器请求
	// UserAgent string = "Dalvik/2.1.0 (Linux; U; Android 9; Valve Steam App Version/3)"
)

// 全局变量定义
var (
	// Domain Steam平台各服务域名配置
	Domain domain = domain{
		Store:     "store.steampowered.com",    // Steam商店
		Community: "steamcommunity.com",        // Steam社区
		Help:      "help.steampowered.com",     // Steam帮助
		CheckOut:  "checkout.steampowered.com", // Steam结账
		Tv:        "steam.tv",                  // Steam电视
		Friend:    "s.team",                    // Steam好友短链接
		Api:       "api.steampowered.com",      // Steam API
		Login:     "login.steampowered.com",    // Steam登录
	}

	// 基础服务端点URL
	Origin          string = Scheme + Domain.Store                           // Steam商店主域名
	Account         string = Scheme + Domain.Store + "/account/"             // 获取账号信息的端点
	Language        string = Scheme + Domain.Store + "/account/setlanguage/" // 设置语言偏好的端点
	CommunityOrigin string = Scheme + Domain.Community                       // Steam社区主域名

	// 身份认证相关API端点
	GetPasswordRSAPublicKey        string = Scheme + Domain.Api + "/IAuthenticationService/GetPasswordRSAPublicKey/v1"             // 获取密码RSA公钥，用于密码加密
	BeginAuthSessionViaCredentials string = Scheme + Domain.Api + "/IAuthenticationService/BeginAuthSessionViaCredentials/v1"      // 通过凭据开始认证会话
	PollAuthSessionStatus          string = Scheme + Domain.Api + "/IAuthenticationService/PollAuthSessionStatus/v1"               // 轮询认证会话状态
	UpdateCode                     string = Scheme + Domain.Api + "/IAuthenticationService/UpdateAuthSessionWithSteamGuardCode/v1" // 使用Steam Guard代码更新认证会话
	AjaxRefresh                    string = Scheme + Domain.Login + "/jwt/ajaxrefresh"                                             // JWT令牌刷新端点
	FinalizeLogin                  string = Scheme + Domain.Login + "/jwt/finalizelogin"                                           // 完成登录流程的端点
	CheckEmailCode                 string = Scheme + Domain.Login + "/jwt/checkdevice/"                                            // 检查邮箱验证码的端点

	// 手机令牌(Steam Guard移动认证器)相关API端点
	QueryTime         string = Scheme + Domain.Api + "/ITwoFactorService/QueryTime/v1/"              // 查询Steam服务器时间，用于时间同步
	ConfirmationList  string = Scheme + Domain.Community + "/ITwoFactorService/ConfirmationList/v1/" // 获取市场交易待确认列表
	MultiConfirmation string = Scheme + Domain.Community + "/mobileconf/multiajaxop"                 // 批量确认操作端点(目前测试失败)

	// Steam积分/点数系统相关API端点
	GetReactions      string = Scheme + Domain.Api + "/ILoyaltyRewardsService/GetReactions/v1"      // 获取用户的反应/表情记录
	GetSummary        string = Scheme + Domain.Api + "/ILoyaltyRewardsService/GetSummary/v1"        // 获取积分系统摘要信息
	GetReactionConfig string = Scheme + Domain.Api + "/ILoyaltyRewardsService/GetReactionConfig/v1" // 获取反应配置信息
	AddReaction       string = Scheme + Domain.Api + "/ILoyaltyRewardsService/AddReaction/v1"       // 添加反应/表情到指定内容

	// 市场交易相关API端点
	GetMyListings       string = Scheme + Domain.Community + "/market/mylistings"     // 获取用户的上架列表
	RemoveMyListings    string = Scheme + Domain.Community + "/market/removelisting"  // 删除用户的已上架或待确认的物品
	GetInventory        string = Scheme + Domain.Community + "/inventory"             // 获取用户库存
	PutList             string = Scheme + Domain.Community + "/market/sellitem"       // 上架物品
	GetConfirmationList string = Scheme + Domain.Community + "/mobileconf/getlist"    // 获取待确认列表
	Confirmation        string = Scheme + Domain.Community + "/mobileconf/ajaxop"     // 确认上架
	BuyListing          string = Scheme + Domain.Community + "/market/buylisting"     // 购买物品
	CreateOrder         string = Scheme + Domain.Community + "/market/createbuyorder" // 创建订单

	// 游戏更新
	GetGameUpdateInofs    string = Scheme + Domain.Store + "/news/app" // 获取游戏更新信息
	CheckAccountAvailable string = Scheme + Domain.Community + "/market/eligibilitycheck/"

	// 好友相关API端点
	AddFriendByLink  string = Scheme + Domain.Community + "/invites/ajaxredeem/"      // 通过链接添加好友
	AddFriendAjax    string = Scheme + Domain.Community + "/actions/AddFriendAjax"    // 通过好友码添加好友
	RemoveFriendAjax string = Scheme + Domain.Community + "/actions/RemoveFriendAjax" // 通过好友码删除好友
	Ajaxresolveusers string = Scheme + Domain.Community + "/actions/ajaxresolveusers" // 查看好友信息

	// 购物车相关API端点
	ClearCart      string = Scheme + Domain.Api + "/IAccountCartService/DeleteCart/v1" // 清空购物车
	CartIndex      string = Scheme + Domain.Store + "/cart/"
	AddItemsToCart string = Scheme + Domain.Api + "/IAccountCartService/AddItemsToCart/v1"

	// 订单相关 API 端点
	InitTransaction     string = Scheme + Domain.CheckOut + "/checkout/inittransaction/"
	CancelCartTrans     string = Scheme + Domain.CheckOut + "/checkout/canceltransaction/"
	Finalizetransaction string = Scheme + Domain.CheckOut + "/checkout/finalizetransaction/"
	Getfinalprice       string = Scheme + Domain.CheckOut + "/checkout/getfinalprice/"

	// 礼物相关 API
	UnsendGiftSubmit  string = Scheme + Domain.CheckOut + "/checkout/unsendgiftsubmit/"
	TransactionStatus string = Scheme + Domain.CheckOut + "/checkout/transactionstatus/"
)
