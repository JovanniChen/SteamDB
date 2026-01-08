package Model

type AddCartItem struct {
	PackageID       uint32
	BundleID        uint32
	AccountidGiftee uint32
	Message         string
}

type GamePurchaseAction struct {
	IsBundle        int    `json:"isBundle"`        // 0=标准版, 1=捆绑包
	BundleInfoTexts string `json:"bundleInfoTexts"` // 捆绑包信息文本
	GameName        string `json:"gameName"`        // 游戏名称
	FinalPrice      string `json:"finalPrice"`      // 最终价格（数字）
	FinalPriceText  string `json:"finalPriceText"`  // 最终价格文本（带货币符号）
	CountryCode     string `json:"countryCode"`     // 国家代码
	AddToCartIds    string `json:"addToCartIds"`    // 添加到购物车的ID
}
