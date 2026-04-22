package Model

// 库存相关结构体
type InventoryResponse struct {
	Success      int8          `json:"success"`
	Assets       []Asset       `json:"assets"`
	Descriptions []Description `json:"descriptions"`
}

type Asset struct {
	AppID      int    `json:"appid"`
	ContextID  string `json:"contextid"`
	AssetID    string `json:"assetid"`
	ClassID    string `json:"classid"`
	InstanceID string `json:"instanceid"`
	Amount     string `json:"amount"`
}

type Description struct {
	ClassID    string `json:"classid"`
	InstanceID string `json:"instanceid"`
	Name       string `json:"name"`
	MarketName string `json:"market_name"`
	Tradable   int    `json:"tradable"`
	Marketable int    `json:"marketable"`
	Commodity  int    `json:"commodity"`
}

// 将 Asset 和 Description 结合起来的结构体
type Item struct {
	AssetID    string  `json:"asset_id"`    // 物品资产ID
	ClassID    string  `json:"class_id"`    // 物品类别ID
	InstanceID string  `json:"instance_id"` // 物品实例ID
	Name       string  `json:"name"`        // 物品名称
	MarketName string  `json:"market_name"` // 市场名称
	Price      float64 `json:"price"`       // 价格
	Currency   int     `json:"currency"`    // 货币类型
	Tradable   bool    `json:"tradable"`    // 是否可交易
	Marketable bool    `json:"marketable"`  // 是否可在市场交易
	ListingID  string  `json:"listing_id"`  // 上架ID(如果已上架)
}

// PutListResponse 上架物品响应
type PutListResponse struct {
	Success                 bool   `json:"success"`
	RequiresConfirmation    int    `json:"requires_confirmation"`
	NeedsMobileConfirmation bool   `json:"needs_mobile_confirmation"`
	NeedsEmailConfirmation  bool   `json:"needs_email_confirmation"`
	EmailDomain             string `json:"email_domain"`
	Message                 string `json:"message"`
}

// ConfirmationsResponse 确认列表响应
type ConfirmationsResponse struct {
	Success       bool           `json:"success"`
	Confirmations []Confirmation `json:"conf"`
}

// Confirmation 待确认项目
type Confirmation struct {
	Type         int      `json:"type"`
	TypeName     string   `json:"type_name"`
	ID           string   `json:"id"`
	CreatorID    string   `json:"creator_id"`
	Nonce        string   `json:"nonce"`
	CreationTime int64    `json:"creation_time"`
	Cancel       string   `json:"cancel"`
	Accept       string   `json:"accept"`
	Icon         string   `json:"icon"`
	Multi        bool     `json:"multi"`
	Headline     string   `json:"headline"`
	Summary      []string `json:"summary"`
}

// ProcessConfirmationResponse 处理确认请求响应
type ProcessConfirmationResponse struct {
	Success bool `json:"success"`
}

type ConfirmationResult struct {
	Success bool
	Result  MyListingReponse
}

type BuyListingNeedConfirmationResponse struct {
	NeedConfirmation bool              `json:"need_confirmation"`
	Confirmation     map[string]string `json:"confirmation"`
	Success          int               `json:"success"`
}

// {"wallet_info":{"wallet_currency":23,"wallet_country":"CN","wallet_state":"","wallet_fee":"1","wallet_fee_minimum":"1","wallet_fee_percent":"0.05","wallet_publisher_fee_percent_default":"0.10","wallet_fee_base":"0","wallet_balance":"1036","wallet_delayed_balance":"0","wallet_max_balance":"1400000","wallet_trade_max_balance":"1260000","success":1,"rwgrsn":-2}}
type BuyListingResponse struct {
	WalletInfo WalletInfo `json:"wallet_info"`
}

type WalletInfo struct {
	WalletCurrency                   int     `json:"wallet_currency"`           // 币种代码
	WalletCountry                    string  `json:"wallet_country"`            // 国家代码
	WalletState                      string  `json:"wallet_state"`              // 州/省（空串）
	WalletFee                        float64 `json:"wallet_fee,string"`         // 手续费（元）
	WalletFeeMinimum                 float64 `json:"wallet_fee_minimum,string"` // 最低手续费
	WalletFeePercent                 float64 `json:"wallet_fee_percent,string"` // 手续费百分比
	WalletPublisherFeePercentDefault float64 `json:"wallet_publisher_fee_percent_default,string"`
	WalletFeeBase                    float64 `json:"wallet_fee_base,string"` // 基础手续费
	WalletBalance                    int64   `json:"wallet_balance,string"`  // 当前余额（分）
	WalletDelayedBalance             int64   `json:"wallet_delayed_balance,string"`
	WalletMaxBalance                 int64   `json:"wallet_max_balance,string"`       // 钱包上限
	WalletTradeMaxBalance            int64   `json:"wallet_trade_max_balance,string"` // 单次交易上限
	Success                          int     `json:"success"`
	Rwgrsn                           int     `json:"rwgrsn"`
}

type BuyListingFailedResponse struct {
	Message string `json:"message"`
}

type CreateOrderResponse struct {
	NeedConfirmation bool              `json:"need_confirmation"`
	Confirmation     map[string]string `json:"confirmation"`
	Success          int               `json:"success"`
}

type GetMyListingResponse struct {
	Success           bool               `json:"success"`
	PageSize          int                `json:"pagesize"`
	TotalCount        int                `json:"total_count"`
	Start             int                `json:"start"`
	NumActiveListings int                `json:"num_active_listings"`
	Listings          []MyListingReponse `json:"listings"`
	ListingsToConfirm []MyListingReponse `json:"listings_to_confirm"`
}

type MyListingReponse struct {
	ListingID          string         `json:"listingid"`
	Asset              MyListingAsset `json:"asset"`
	SellerReceivePrice int            `json:"converted_price"` // 卖家到账价
	Fee                int            `json:"fee"`             // 手续费
}

type MyListingAsset struct {
	MarketHashName string `json:"market_hash_name"`
}

type MarketListingResponse struct {
	Success     bool                       `json:"success"`
	Start       int                        `json:"start"`
	PageSize    int                        `json:"pagesize"`
	TotalCount  int                        `json:"total_count"`
	ListingInfo map[string]ListingInfoItem `json:"listinginfo"`
}

type ListingInfoItem struct {
	ListingID             string `json:"listingid"`
	ConvertedSteamFee     int    `json:"converted_steam_fee"`
	ConvertedPublisherFee int    `json:"converted_publisher_fee"`
	ConvertedPricePerUnit int    `json:"converted_price_per_unit"`
}
