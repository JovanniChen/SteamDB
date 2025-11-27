package Model

type TransactionStatusResponse struct {
	// 	{
	//     "success": 22,
	//     "purchaseresultdetail": 0,
	//     "purchasereceipt": {
	//         "paymentmethod": 11,
	//         "purchasestatus": 0,
	//         "resultdetail": 0,
	//         "baseprice": "3000",
	//         "totaldiscount": "2700",
	//         "tax": "0",
	//         "shipping": "0",
	//         "packageid": -1,
	//         "transactiontime": 1764225154,
	//         "transactionid": "193220467646477559",
	//         "currencycode": 23,
	//         "formattedTotal": "¥ 3.00",
	//         "rewardPointsBalance": "233"
	//     },
	//     "strReceiptPageHTML": "",
	//     "bShowBRSpecificCreditCardError": false
	// }
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
