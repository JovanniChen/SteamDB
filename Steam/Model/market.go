package Model

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

// AcceptConfirmationResponse 接受确认响应
type AcceptConfirmationResponse struct {
	Success bool `json:"success"`
}
