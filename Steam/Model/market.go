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
}

// ConfirmationsResponse 确认列表响应
type ConfirmationsResponse struct {
	Success       bool           `json:"success"`
	Confirmations []Confirmation `json:"conf"`
}

// Confirmation 待确认项目
type Confirmation struct {
	ID        string `json:"id"`
	Nonce     string `json:"nonce"`
	CreatorID string `json:"creator_id"`
}

// ProcessConfirmationResponse 处理确认请求响应
type ProcessConfirmationResponse struct {
	Success bool `json:"success"`
}

type BuyListingResponse struct {
	NeedConfirmation bool              `json:"need_confirmation"`
	Confirmation     map[string]string `json:"confirmation"`
	Success          int               `json:"success"`
}
