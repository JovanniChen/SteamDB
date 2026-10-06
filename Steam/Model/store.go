package Model

type StorePurchaseHistoryResult struct {
	TotalRecords                  int
	LatestUnrefundedGiftPurchases []StorePurchaseHistoryRecord
}

type StorePurchaseHistoryRecord struct {
	Index           int
	TransactionID   string
	TransactionURL  string
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
	// InventoryAssetKeys contains the appid_contextid_assetid values found
	// on the purchase detail pages for this transaction.
	InventoryAssetKeys []string
}
