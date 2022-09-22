package main

import (
	"log"
	"time"

	bybitSpotV3 "github.com/testingaccmar/ccxt_beyang_bybit/bybit/spot/v3"
)

func main() {
	cfg := &bybitSpotV3.Configuration{
		Addr:      bybitSpotV3.HostMainnetPublicTopics,
		ApiKey:    "",
		SecretKey: "",
		DebugMode: true,
	}
	b := bybitSpotV3.New(cfg)
	b.Start()

	pair := b.GetPair("BTC", "USDT")
	b.Subscribe(bybitSpotV3.ChannelBookTicker, pair)

	b.On(bybitSpotV3.ChannelBookTicker, handleBookTicker)
	b.On(bybitSpotV3.ChannelBookTicker, handleBestBidPrice)

	go func() {
		time.Sleep(5 * time.Second)
		balance := b.GetBalance()
		for _, coins := range balance.Result.Balances {
			log.Printf("coin = %s, total = %s", coins.Coin, coins.Total)
		}
	}()

	//	не дает прекратить работу программы
	forever := make(chan struct{})
	<-forever
}

func handleBookTicker(symbol string, data bybitSpotV3.BookTicker) {
	log.Printf("Bybit BookTicker  %s: %v", symbol, data)
}

func handleBestBidPrice(symbol string, data bybitSpotV3.BookTicker) {
	log.Printf("Bybit BookTicker  %s: BestBidPrice : %s", symbol, data.Data.Bp)
}
