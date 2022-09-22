package v3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/goccy/go-json" // для создания собственных json файлов и преобразования json в структуру
)

func (ex *ByBitWS) GetBalance() WalletBalance {
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
				"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3",
				"File": "api.go",
				"Functions" : "(ex *ByBitWS) GetBalance() (WalletBalance)",
				"Function where err" : "json.Unmarshal",
				"Exchange" : "Bybit",
				"Comment" : %s to WalletBalance struct,
				"Error" : %s
			}`, string(data), err)
		log.Fatal()
	}

	return walletBalance

}
