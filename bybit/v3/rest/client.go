package rest

import (
	"log"
	"strings"

	binanceRest "github.com/TestingAccMar/CCXT_beYANG_Binance/binance/rest"
	bybitRest "github.com/TestingAccMar/CCXT_beYANG_ByBit/bybit/spot/v3/rest"
)

type ExchangeRest struct {
	Name           string
	Addr           string
	ApiKey         string
	SecretKey      string
	APIKeyPassword string //	необходим только для okx
	DebugMode      bool
}

func NewExchange(exWs ExchangeRest) EchangeInterface {
	var ex EchangeInterface
	switch strings.ToLower(exWs.Name) {
	case "binance":
		cfg := &binanceRest.Configuration{
			Addr:      exWs.Addr,
			ApiKey:    exWs.ApiKey,
			SecretKey: exWs.SecretKey,
			DebugMode: exWs.DebugMode,
		}
		ex = binanceRest.New(cfg)
	case "bybit":
		cfg := &bybitRest.Configuration{
			Addr:      exWs.Addr,
			ApiKey:    exWs.ApiKey,
			SecretKey: exWs.SecretKey,
			DebugMode: exWs.DebugMode,
		}
		ex = bybitRest.New(cfg)
	case "ftx":
	case "okx":
	default:
		log.Printf(`
		{
			"Status" : "INFO",
			"Comment" : "Данная биржа пока не поддерживается"
		}`)
		log.Fatal()
	}
	return ex
}

type EchangeInterface interface {
	GetBalance() binanceRest.WalletBalance
}
