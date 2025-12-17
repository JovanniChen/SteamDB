package Model

type TransactionStatusResponse struct {
	Success              int `json:"success"`
	PurchaseResultDetail int `json:"purchaseresultdetail"`
	PurchaseReceipt      struct {
		PaymentMethod       int    `json:"paymentmethod"`
		PurchaseStatus      int    `json:"purchasestatus"`
		ResultDetail        int    `json:"resultdetail"`
		BasePrice           string `json:"baseprice"`
		TotalDiscount       string `json:"totaldiscount"`
		Tax                 string `json:"tax"`
		Shipping            string `json:"shipping"`
		PackageID           int    `json:"packageid"`
		TransactionTime     int64  `json:"transactiontime"`
		TransactionID       string `json:"transactionid"`
		CurrencyCode        int    `json:"currencycode"`
		FormattedTotal      string `json:"formattedTotal"`
		RewardPointsBalance string `json:"rewardPointsBalance"`
	} `json:"purchasereceipt"`
	StrReceiptPageHTML             string `json:"strReceiptPageHTML"`
	BShowBRSpecificCreditCardError bool   `json:"bShowBRSpecificCreditCardError"`
}
