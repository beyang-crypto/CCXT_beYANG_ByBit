package ws

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
