package Model

type StorePurchaseHistoryResult struct {
	TotalRecords                  int
	LatestUnrefundedGiftPurchases []StorePurchaseHistoryRecord
}

type StorePurchaseHistoryRecord struct {
	Index           int
	TransactionID   string
	Date            string
	Items           []string
	Receivers       []string
	TransactionType string
	Payment         string
	BasePrice       string
	Tax             string
	Shipping        string
	Total           string
	WalletChange    string
	WalletBalance   string
	Refunded        bool
}
