package ws

import (
	"log"
	"strings"

	binanceWs "github.com/TestingAccMar/CCXT_beYANG_Binance/binance/ws"
	bybit "github.com/TestingAccMar/CCXT_beYANG_ByBit/bybit/spot/v3"
	"github.com/chuckpreslar/emission"
)

//  Для вытаскивания одного значения из файла json
// Эмитер необходим для удобного выполнения функции в какой-то момент
// для создания собственных json файлов и преобразования json в структуру

// из нее уже будет создаваться все остальное
type ExchangeWS struct {
	Name           string
	Addr           string
	ApiKey         string
	SecretKey      string
	APIKeyPassword string //	необходим только для okx
	DebugMode      bool
}

func NewExchange(exWs ExchangeWS) EchangeInterface {
	var ex EchangeInterface
	switch strings.ToLower(exWs.Name) {
	case "binance":
		cfg := &binanceWs.Configuration{
			Addr:      exWs.Addr,
			ApiKey:    exWs.ApiKey,
			SecretKey: exWs.SecretKey,
			DebugMode: exWs.DebugMode,
		}
		ex = binanceWs.New(cfg)
	case "bybit":
		cfg := &bybit.Configuration{
			Addr:      exWs.Addr,
			ApiKey:    exWs.ApiKey,
			SecretKey: exWs.SecretKey,
			DebugMode: exWs.DebugMode,
		}
		ex = bybit.New(cfg)
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
	Subscribe(channel string, coin string)
	GetPair(coin1 string, coin2 string) string
	Start() error

	On(event interface{}, listener interface{}) *emission.Emitter
	Emit(event interface{}, arguments ...interface{}) *emission.Emitter
	Off(event interface{}, listener interface{}) *emission.Emitter
}

func GetPair(ex EchangeInterface, coin1 string, coin2 string) string {
	return ex.GetPair(coin1, coin2)
}

func Start(ex EchangeInterface) {
	ex.Start()
}

func Subscribe(ex EchangeInterface, channel string, coin string) {
	ex.Subscribe(channel, coin)
}

func On(ex EchangeInterface, event interface{}, listener interface{}) *emission.Emitter {
	return ex.On(event, listener)
}

func Emit(ex EchangeInterface, event interface{}, arguments ...interface{}) *emission.Emitter {
	return ex.Emit(event, arguments)
}

func Off(ex EchangeInterface, event interface{}, listener interface{}) *emission.Emitter {
	return ex.Off(event, listener)
}
