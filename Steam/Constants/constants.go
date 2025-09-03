// Constants包 - Steam平台相关的配置和常量定义
// 包含Steam各个服务的域名、API端点和系统配置
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
	GetMyListings       string = Scheme + Domain.Community + "/market/mylistings"  // 获取用户的上架列表
	GetInventory        string = Scheme + Domain.Community + "/inventory"          // 获取用户库存
	PutList             string = Scheme + Domain.Community + "/market/sellitem"    // 上架物品
	GetConfirmationList string = Scheme + Domain.Community + "/mobileconf/getlist" // 获取待确认列表
	Confirmation        string = Scheme + Domain.Community + "/mobileconf/ajaxop"  // 确认上架
)
