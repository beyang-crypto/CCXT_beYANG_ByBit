package v3

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

// https://bybit-exchange.github.io/docs/spot/v3/#t-websocektbookticker

type dataBookTicker struct {
	Ap string `json:"ap"`
	Aq string `json:"aq"`
	Bp string `json:"bp"`
	Bq string `json:"bq"`
	S  string `json:"s"`
	T  int64  `json:"t"`
}

type BookTicker struct {
	Data  dataBookTicker `json:"data"`
	Type  string         `json:"type"`
	Topic string         `json:"topic"`
	//Ts    int64          `json:"ts"`
}
