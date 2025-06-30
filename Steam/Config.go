package Steam

type domain struct {
	Store     string
	Community string
	Help      string
	CheckOut  string
	Tv        string
	Friend    string
	Api       string
	Login     string
}

const (
	Scheme    string = "https://"
	Tries     int    = 3
	UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.95 Safari/537.36"
)

var (
	Domain domain = domain{
		Store:     "store.steampowered.com",
		Community: "steamcommunity.com",
		Help:      "help.steampowered.com",
		CheckOut:  "checkout.steampowered.com",
		Tv:        "steam.tv",
		Friend:    "s.team",
		Api:       "api.steampowered.com",
		Login:     "login.steampowered.com",
	}

	Origin   string = Scheme + Domain.Store
	Account  string = Scheme + Domain.Store + "/account/"             // 获取账号信息
	Language string = Scheme + Domain.Store + "/account/setlanguage/" // 设置语言

	GetPasswordRSAPublicKey        string = Scheme + Domain.Api + "/IAuthenticationService/GetPasswordRSAPublicKey/v1"
	BeginAuthSessionViaCredentials string = Scheme + Domain.Api + "/IAuthenticationService/BeginAuthSessionViaCredentials/v1"
	PollAuthSessionStatus          string = Scheme + Domain.Api + "/IAuthenticationService/PollAuthSessionStatus/v1"
	UpdateCode                     string = Scheme + Domain.Api + "/IAuthenticationService/UpdateAuthSessionWithSteamGuardCode/v1"
	AjaxRefresh                    string = Scheme + Domain.Login + "/jwt/ajaxrefresh"
	FinalizeLogin                  string = Scheme + Domain.Login + "/jwt/finalizelogin"
	CheckEmailCode                 string = Scheme + Domain.Login + "/jwt/checkdevice/"

	// 手机令牌
	QueryTime         string = Scheme + Domain.Api + "/ITwoFactorService/QueryTime/v1/"              //steam服务器时间
	ConfirmationList  string = Scheme + Domain.Community + "/ITwoFactorService/ConfirmationList/v1/" //市场待确认列表
	Confirmation      string = Scheme + Domain.Community + "/mobileconf/ajaxop"                      //确认列表操作 GET
	MultiConfirmation string = Scheme + Domain.Community + "/mobileconf/multiajaxop"                 //确认列表批量操作 测试失败
)
