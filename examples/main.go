package main

import (
	"log"
	"time"

	bybitRest "github.com/TestingAccMar/CCXT_beYANG_ByBit/bybit/spot/v3/rest"
	bybitWS "github.com/TestingAccMar/CCXT_beYANG_ByBit/bybit/spot/v3/ws"
)

func main() {
	cfg := &bybitRest.Configuration{
		Addr:      bybitRest.RestMainnetBybit,
		ApiKey:    "",
		SecretKey: "",
		DebugMode: true,
	}
	b := bybitRest.New(cfg)
	go func() {
		time.Sleep(5 * time.Second)
		balance := bybitRest.BybitToWalletBalance(b.GetBalance())
		for _, coins := range balance.Result.Balances {
			log.Printf("coin = %s, total = %s", coins.Coin, coins.Total)
		}
	}()

	//	не дает прекратить работу программы
	forever := make(chan struct{})
	<-forever
}

func handleBookTicker(symbol string, data bybitWS.BookTicker) {
	log.Printf("Bybit BookTicker  %s: %v", symbol, data)
}

func handleBestBidPrice(symbol string, data bybitWS.BookTicker) {
	log.Printf("Bybit BookTicker  %s: BestBidPrice : %s", symbol, data.Data.Bp)
}
