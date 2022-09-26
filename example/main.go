package main

import (
	"ccxt_beyang/rest"
	restconsts "ccxt_beyang/rest/consts"
	"fmt"
	"log"
	"time"
)

func main() {

	exRestBincance := rest.ExchangeRest{
		Name:      "Binance",
		Addr:      restconsts.BybitRestBaseEndpoint,
		ApiKey:    "vKzCLmKCkZfIkWX8mkNCpJjVhQDs6hWbwvflRVnKbXh5C17npsDmSKjQpURwyanU",
		SecretKey: "qwyNSz4Gyk4DeEtBkjoSYuuol3DQe9n5nktt1gV757Ixda5FBiTmDjpfvQh9CUIA",
		DebugMode: true,
	}

	binanceRestEx := rest.NewExchange(exRestBincance)
	go func() {
		time.Sleep(1 * time.Second)
		balance := binanceRestEx.GetBalance()
		for _, coins := range balance.Balances {
			log.Printf("coin = %s, total = %s", coins.Asset, coins.Free)
		}
	}()
	// exBinance := ws.ExchangeWS{
	// 	Name:           "Binance",
	// 	Addr:           wsconsts.BinanceHostMainnetPublicTopics,
	// 	ApiKey:         "",
	// 	APIKeyPassword: "",
	// 	DebugMode:      true,
	// }

	// exBybit := ws.ExchangeWS{
	// 	Name:           "Bybit",
	// 	Addr:           wsconsts.BybitHostMainnetPublicTopics,
	// 	ApiKey:         "",
	// 	APIKeyPassword: "",
	// 	DebugMode:      true,
	// }

	// binanceEx := ws.NewExchange(exBinance)
	// bybitEx := ws.NewExchange(exBybit)
	// ws.Start(binanceEx)
	// ws.Start(bybitEx)

	// pairBinance := ws.GetPair(binanceEx, "btc", "usdt")
	// pairBybit := ws.GetPair(bybitEx, "BTC", "USDT")
	// ws.Subscribe(binanceEx, wsconsts.BinanceChannelTicker, pairBinance)
	// ws.Subscribe(bybitEx, wsconsts.BybitChannelBookTicker, pairBybit)

	// ws.On(binanceEx, wsconsts.BinanceChannelTicker, handleBookTickerBinance)
	// ws.On(bybitEx, wsconsts.BybitChannelBookTicker, handleBookTickerBybit)

	forever := make(chan struct{})
	<-forever
}

func handleBookTickerBybit(symbol string, data interface{}) {
	log.Printf("Bybit Ticker  %s: %s", symbol, fmt.Sprintf("%+v\n", data))
}

func handleBookTickerBinance(symbol string, data interface{}) {
	log.Printf("Binance Ticker  %s: %s", symbol, fmt.Sprintf("%+v\n", data))
}
