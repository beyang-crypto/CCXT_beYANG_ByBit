package main

import (
	"log"

	bybitWS "github.com/TestingAccMar/CCXT_beYANG_ByBit/bybit/spot/v3/ws"
)

func main() {
	cfg := &bybitWS.Configuration{
		Addr:      bybitWS.HostMainnetPublicTopics,
		ApiKey:    "",
		SecretKey: "",
		DebugMode: true,
	}
	b := bybitWS.New(cfg)
	b.Start()

	pair1 := b.GetPair("btc", "usdt")
	pair2 := b.GetPair("eth", "usdt")
	pair3 := b.GetPair("xrp", "usdt")
	pair4 := b.GetPair("ada", "usdt")
	pair5 := b.GetPair("sol", "usdt")
	pair6 := b.GetPair("doge", "usdt")
	pair7 := b.GetPair("matic", "usdt")
	pair8 := b.GetPair("shib", "usdt")
	pair9 := b.GetPair("trx", "usdt")
	pair10 := b.GetPair("uni", "usdt")
	pair11 := b.GetPair("avax", "usdt")
	pair12 := b.GetPair("ltc", "usdt")
	pair13 := b.GetPair("etc", "usdt")
	pair14 := b.GetPair("link", "usdt")
	pair15 := b.GetPair("atom", "usdt")

	b.Subscribe(bybitWS.ChannelBookTicker, []string{pair1, pair2, pair3, pair4, pair5, pair6, pair7, pair8, pair9, pair10})
	b.Subscribe(bybitWS.ChannelBookTicker, []string{pair11, pair12, pair13, pair14, pair15})
	b.On(bybitWS.ChannelBookTicker, handleBestBidPrice)
	// cfg := &bybitRest.Configuration{
	// 	Addr:      bybitRest.RestMainnetBybit,
	// 	ApiKey:    "",
	// 	SecretKey: "",
	// 	DebugMode: true,
	// }
	// b := bybitRest.New(cfg)
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	balance := bybitRest.BybitToWalletBalance(b.GetBalance())
	// 	for _, coins := range balance.Result.Balances {
	// 		log.Printf("coin = %s, total = %s", coins.Coin, coins.Total)
	// 	}
	// }()

	//	не дает прекратить работу программы
	forever := make(chan struct{})
	<-forever
}

func handleBookTicker(name string, symbol string, data bybitWS.BookTicker) {
	log.Printf("%s BookTicker  %s: %v", name, symbol, data)
}

func handleBestBidPrice(name string, symbol string, data bybitWS.BookTicker) {
	log.Printf("%s BookTicker  %s: BestBidPrice : %s", name, symbol, data.Data.Bp)
}
