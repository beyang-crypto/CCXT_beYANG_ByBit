package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

const (
	RestMainnetBybit  = "https://api.bybit.com"
	RestMainnetBytick = "https://api.bytick.com"

	RestTestnetBybit = "https://api-testnet.bybit.com"
)

type Configuration struct {
	Addr      string `json:"addr"`
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
	DebugMode bool   `json:"debug_mode"`
}

type ByBitRest struct {
	cfg *Configuration
}

func (b *ByBitRest) GetPair(args ...string) string {
	pair := args[0] + args[1]
	return strings.ToUpper(pair)
}

func New(config *Configuration) *ByBitRest {

	// 	потом тут добавятся различные другие настройки
	b := &ByBitRest{
		cfg: config,
	}
	return b
}

func (ex *ByBitRest) GetBalance() interface{} {
	//	https://bybit-exchange.github.io/docs/spot/?python--pybit#t-wallet
	//	получение времяни
	ts := time.Now().UTC().Unix() * 1000
	apiKey := ex.cfg.ApiKey
	secretKey := ex.cfg.SecretKey

	parms := fmt.Sprintf("api_key=%s&recv_Window=5000&timestamp=%d", apiKey, ts)
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(parms))
	parms += "&sign=" + hex.EncodeToString(mac.Sum(nil))
	//	реализация метода GET
	resp, err := http.Get("https://api.bybit.com/spot/v1/account?" + parms)

	//	код для вывода полученных данных
	if err != nil {
		log.Fatalln(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if ex.cfg.DebugMode {
		log.Printf("STATUS: DEBUG\tEXCHANGE: Bybit\tAPI: Rest\tBybitWalletBalance %v", string(data))
	}

	// {
	// 	"ret_code": 0,
	// 	"ret_msg": "",
	// 	"ext_code": null,
	// 	"ext_info": null,
	// 	"result": {
	// 		"balances": [
	// 			{
	// 				"coin": "USDT",
	// 				"coinId": "USDT",
	// 				"coinName": "USDT",
	// 				"total": "10",
	// 				"free": "10",
	// 				"locked": "0"
	// 			}
	// 		]
	// 	}
	// }

	var walletBalance WalletBalance
	err = json.Unmarshal(data, &walletBalance)
	if err != nil {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_beYANG_Bybit/spot/v3/rest",
				"File": "client.go",
				"Functions" : "(ex *ByBitRest) GetBalance() WalletBalance",
				"Function where err" : "json.Unmarshal",
				"Exchange" : "Bybit",
				"Comment" : %s to WalletBalance struct,
				"Error" : %s
			}`, string(data), err)
		log.Fatal()
	}

	return walletBalance

}
