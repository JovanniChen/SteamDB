package Constants

// Country 国家信息，包含名称、ISO国家代码、货币代码、Steam内部货币ID和货币符号
type Country struct {
	Name                string `json:"name"`
	Code                string `json:"code"`
	CurrencyCode        string `json:"currency_code"`
	SteamCurrencyID     int    `json:"steam_currency_id"`      // Steam 社区市场内部货币ID，如人民币=23
	SteamCurrencyFullID int    `json:"steam_currency_full_id"` // Steam 社区市场完整货币ID，如人民币=2023
	CurrencySymbol      string `json:"currency_symbol"`
}

// Countries Steam 社区市场支持的货币/国家，以 ISO 国家代码为键
// 1 - USD / 美元
// 2 - GBP / 英镑
// 3 - EUR / 欧元
// 4 - CHF / 瑞士法郎
// 5 - RUB / 俄罗斯卢布
// 6 - PLN / 波兰兹罗提
// 7 - BRL / 巴西雷亚尔
// 8 - JPY / 日元
// 9 - NOK / 挪威克朗
// 10 - IDR / 印度尼西亚卢比
// 11 - MYR / 马来西亚林吉特
// 12 - PHP / 菲律宾比索
// 13 - SGD / 新加坡元
// 14 - THB / 泰铢
// 15 - VND / 越南盾
// 16 - KRW / 韩元
// 18 - UAH / 乌克兰格里夫纳
// 19 - MXN / 墨西哥比索
// 20 - CAD / 加拿大元
// 21 - AUD / 澳大利亚元
// 22 - NZD / 新西兰元
// 23 - CNY / 人民币（元）
// 24 - INR / 印度卢比
// 25 - CLP / 智利比索
// 26 - PEN / 秘鲁索尔
// 27 - COP / 哥伦比亚比索
// 28 - ZAR / 南非兰特
// 29 - HKD / 港元
// 30 - TWD / 新台币
// 31 - SAR / 沙特里亚尔
// 32 - AED / 阿联酋迪拉姆
// 35 - ILS / 以色列新谢克尔
// 37 - KZT / 哈萨克斯坦坚戈
// 38 - KWD / 科威特第纳尔
// 39 - QAR / 卡塔尔里亚尔
// 40 - CRC / 哥斯达黎加科朗
// 41 - UYU / 乌拉圭比索
var Countries = map[string]Country{
	"US": {Name: "美元", SteamCurrencyID: 1, SteamCurrencyFullID: 2001, CurrencyCode: "USD", CurrencySymbol: "$"},
	"GB": {Name: "英镑", SteamCurrencyID: 2, SteamCurrencyFullID: 2002, CurrencyCode: "GBP", CurrencySymbol: "£"},
	"EU": {Name: "欧元", SteamCurrencyID: 3, SteamCurrencyFullID: 2003, CurrencyCode: "EUR", CurrencySymbol: "€"},
	"CH": {Name: "瑞士法郎", SteamCurrencyID: 4, SteamCurrencyFullID: 2004, CurrencyCode: "CHF", CurrencySymbol: "CHF"},
	"RU": {Name: "俄罗斯卢布", SteamCurrencyID: 5, SteamCurrencyFullID: 2005, CurrencyCode: "RUB", CurrencySymbol: "₽"},
	"PL": {Name: "波兰兹罗提", SteamCurrencyID: 6, SteamCurrencyFullID: 2006, CurrencyCode: "PLN", CurrencySymbol: "zł"},
	"BR": {Name: "巴西雷亚尔", SteamCurrencyID: 7, SteamCurrencyFullID: 2007, CurrencyCode: "BRL", CurrencySymbol: "R$"},
	"JP": {Name: "日元", SteamCurrencyID: 8, SteamCurrencyFullID: 2008, CurrencyCode: "JPY", CurrencySymbol: "¥"},
	"NO": {Name: "挪威克朗", SteamCurrencyID: 9, SteamCurrencyFullID: 2009, CurrencyCode: "NOK", CurrencySymbol: "kr"},
	"ID": {Name: "印度尼西亚卢比", SteamCurrencyID: 10, SteamCurrencyFullID: 2010, CurrencyCode: "IDR", CurrencySymbol: "Rp"},
	"MY": {Name: "马来西亚林吉特", SteamCurrencyID: 11, SteamCurrencyFullID: 2011, CurrencyCode: "MYR", CurrencySymbol: "RM"},
	"PH": {Name: "菲律宾比索", SteamCurrencyID: 12, SteamCurrencyFullID: 2012, CurrencyCode: "PHP", CurrencySymbol: "₱"},
	"SG": {Name: "新加坡元", SteamCurrencyID: 13, SteamCurrencyFullID: 2013, CurrencyCode: "SGD", CurrencySymbol: "S$"},
	"TH": {Name: "泰铢", SteamCurrencyID: 14, SteamCurrencyFullID: 2014, CurrencyCode: "THB", CurrencySymbol: "฿"},
	"VN": {Name: "越南盾", SteamCurrencyID: 15, SteamCurrencyFullID: 2015, CurrencyCode: "VND", CurrencySymbol: "₫"},
	"KR": {Name: "韩元", SteamCurrencyID: 16, SteamCurrencyFullID: 2016, CurrencyCode: "KRW", CurrencySymbol: "₩"},
	"UA": {Name: "乌克兰格里夫纳", SteamCurrencyID: 18, SteamCurrencyFullID: 2018, CurrencyCode: "UAH", CurrencySymbol: "₴"},
	"MX": {Name: "墨西哥比索", SteamCurrencyID: 19, SteamCurrencyFullID: 2019, CurrencyCode: "MXN", CurrencySymbol: "Mex$"},
	"CA": {Name: "加拿大元", SteamCurrencyID: 20, SteamCurrencyFullID: 2020, CurrencyCode: "CAD", CurrencySymbol: "C$"},
	"AU": {Name: "澳大利亚元", SteamCurrencyID: 21, SteamCurrencyFullID: 2021, CurrencyCode: "AUD", CurrencySymbol: "A$"},
	"NZ": {Name: "新西兰元", SteamCurrencyID: 22, SteamCurrencyFullID: 2022, CurrencyCode: "NZD", CurrencySymbol: "NZ$"},
	"CN": {Name: "人民币", SteamCurrencyID: 23, SteamCurrencyFullID: 2023, CurrencyCode: "CNY", CurrencySymbol: "￥"},
	"IN": {Name: "印度卢比", SteamCurrencyID: 24, SteamCurrencyFullID: 2024, CurrencyCode: "INR", CurrencySymbol: "₹"},
	"CL": {Name: "智利比索", SteamCurrencyID: 25, SteamCurrencyFullID: 2025, CurrencyCode: "CLP", CurrencySymbol: "CLP$"},
	"PE": {Name: "秘鲁索尔", SteamCurrencyID: 26, SteamCurrencyFullID: 2026, CurrencyCode: "PEN", CurrencySymbol: "S/."},
	"CO": {Name: "哥伦比亚比索", SteamCurrencyID: 27, SteamCurrencyFullID: 2027, CurrencyCode: "COP", CurrencySymbol: "COL$"},
	"ZA": {Name: "南非兰特", SteamCurrencyID: 28, SteamCurrencyFullID: 2028, CurrencyCode: "ZAR", CurrencySymbol: "R"},
	"HK": {Name: "港元", SteamCurrencyID: 29, SteamCurrencyFullID: 2029, CurrencyCode: "HKD", CurrencySymbol: "HK$"},
	"TW": {Name: "新台币", SteamCurrencyID: 30, SteamCurrencyFullID: 2030, CurrencyCode: "TWD", CurrencySymbol: "NT$"},
	"SA": {Name: "沙特里亚尔", SteamCurrencyID: 31, SteamCurrencyFullID: 2031, CurrencyCode: "SAR", CurrencySymbol: "﷼"},
	"AE": {Name: "阿联酋迪拉姆", SteamCurrencyID: 32, SteamCurrencyFullID: 2032, CurrencyCode: "AED", CurrencySymbol: "د.إ"},
	"IL": {Name: "以色列新谢克尔", SteamCurrencyID: 35, SteamCurrencyFullID: 2035, CurrencyCode: "ILS", CurrencySymbol: "₪"},
	"KZ": {Name: "哈萨克斯坦坚戈", SteamCurrencyID: 37, SteamCurrencyFullID: 2037, CurrencyCode: "KZT", CurrencySymbol: "₸"},
	"KW": {Name: "科威特第纳尔", SteamCurrencyID: 38, SteamCurrencyFullID: 2038, CurrencyCode: "KWD", CurrencySymbol: "د.ك"},
	"QA": {Name: "卡塔尔里亚尔", SteamCurrencyID: 39, SteamCurrencyFullID: 2039, CurrencyCode: "QAR", CurrencySymbol: "ر.ق"},
	"CR": {Name: "哥斯达黎加科朗", SteamCurrencyID: 40, SteamCurrencyFullID: 2040, CurrencyCode: "CRC", CurrencySymbol: "₡"},
	"UY": {Name: "乌拉圭比索", SteamCurrencyID: 41, SteamCurrencyFullID: 2041, CurrencyCode: "UYU", CurrencySymbol: "$U"},
}
