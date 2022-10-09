package rest

// https://bybit-exchange.github.io/docs/spot/v3/#t-balance
type WalletBalance struct {
	RetCode    int                     `json:"retCode"`
	RetMsg     string                  `json:"retMsg"`
	Result     ResultWalletBalance     `json:"result"`
	RetExtMap  RetExtMapWalletBalance  `json:"retExtMap"`
	RetExtInfo RetExtInfoWalletBalance `json:"retExtInfo"`
	Time       int64                   `json:"time"`
}
type BalancesWalletBalance struct {
	Coin   string `json:"coin"`
	CoinID string `json:"coinId"`
	Total  string `json:"total"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}
type ResultWalletBalance struct {
	Balances []BalancesWalletBalance `json:"balances"`
}
type RetExtMapWalletBalance struct {
}
type RetExtInfoWalletBalance struct {
}

func BybitToWalletBalance(data interface{}) (WalletBalance, bool) {
	bt, ok := data.(WalletBalance)
	return bt, ok
}
