package ws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"      //  Для вытаскивания одного значения из файла json
	"github.com/chuckpreslar/emission" // Эмитер необходим для удобного выполнения функции в какой-то момент
	"github.com/goccy/go-json"         // для создания собственных json файлов и преобразования json в структуру
	"github.com/gorilla/websocket"
)

const (
	HostMainnetPublicTopics  = "wss://stream.bybit.com/spot/public/v3"
	HostMainnetPrivateTopics = "wss://stream.bybit.com/spot/private/v3"

	HostTestnetPublicTopics  = "wss://stream-testnet.bybit.com/spot/public/v3"
	HostTestnetPrivateTopics = "wss://stream-testnet.bybit.com/spot/private/v3"
)

const (
	ChannelBookTicker = "bookticker." //Topic: bookticker.{symbol}
)

type Configuration struct {
	Addr      string `json:"addr"`
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
	DebugMode bool   `json:"debug_mode"`
}

type ByBitWS struct {
	cfg  *Configuration
	conn *websocket.Conn

	mu            sync.RWMutex
	subscribeCmds []Cmd //	сохраняем все подписки у данной биржи, чтоб при переподключении можно было к ним повторно подключиться

	emitter *emission.Emitter
}

func (b *ByBitWS) GetPair(coin1 string, coin2 string) string {
	return coin1 + coin2
}

func New(config *Configuration) *ByBitWS {

	// 	потом тут добавятся различные другие настройки
	b := &ByBitWS{
		cfg:     config,
		emitter: emission.NewEmitter(),
	}
	return b
}

func (b *ByBitWS) Subscribe(args ...string) {
	switch len(args) {
	case 2:
		b.Subscribe2(args[0], args[1])
	default:
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
				"File": "client.go",
				"Functions" : "(b *ByBitWS) Subscribe(args ...string)",
				"Exchange" : "Bybit",
				"Data" : [%v],
				"Comment" : "Слишком много аргументов"
			}`, args)
		log.Fatal()
	}
}

func (b *ByBitWS) Subscribe2(channel string, coin string) {
	cmd := Cmd{
		Op:   "subscribe",
		Args: []interface{}{channel + coin},
	}
	b.subscribeCmds = append(b.subscribeCmds, cmd)
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку part 1")
	}
	b.SendCmd(cmd)
}

//	отправка команды на сервер в отдельной функции для того, чтобы при переподключении быстро подписаться на все предыдущие каналы
func (b *ByBitWS) SendCmd(cmd Cmd) {
	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws/ws",
				"File": "client.go",
				"Functions" : "(b *ByBitWS) sendCmd(cmd Cmd)",
				"Function where err" : "json.Marshal",
				"Exchange" : "Bybit",
				"Data" : [%s],
				"Error" : %s
			}`, cmd, err)
		log.Fatal()
	}
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку part 2")
	}
	b.Send(string(data))
}

func (b *ByBitWS) Send(msg string) (err error) {
	defer func() {
		// recover необходим для корректной обработки паники
		if r := recover(); r != nil {
			if err != nil {
				log.Printf(`
					{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws/ws",
						"File": "client.go",
						"Functions" : "(b *ByBitWS) Send(msg string) (err error)",
						"Function where err" : "b.conn.WriteMessage",
						"Exchange" : "Bybit",
						"Data" : [websocket.TextMessage, %s],
						"Error" : %s,
						"Recover" : %v
					}`, msg, err, r)
				log.Fatal()
			}
			err = errors.New(fmt.Sprintf("BybitWs send error: %v", r))
		}
	}()
	if b.cfg.DebugMode {
		log.Printf("Отправка сообщения на сервер. текст сообщения:%s", msg)
	}

	err = b.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	return
}

// подключение к серверу и постоянное чтение приходящих ответов
func (b *ByBitWS) Start() error {
	if b.cfg.DebugMode {
		log.Printf("Начало подключения к серверу")
	}
	b.connect()

	cancel := make(chan struct{})

	go func() {
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				b.ping()
			case <-cancel:
				return
			}
		}
	}()

	go func() {
		defer close(cancel)

		for {
			_, data, err := b.conn.ReadMessage()
			if err != nil {

				if websocket.IsCloseError(err, 1006) {
					b.closeAndReconnect()
					//Необходим вызв SubscribeToTicker в отдельной горутине, рекурсия, думаю, тут неуместна
					log.Printf("Status: INFO	ошибка 1006 начинается переподключение к серверу")

				} else {
					b.conn.Close()
					log.Printf(`
						{
							"Status" : "Error",
							"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
							"File": "client.go",
							"Functions" : "(b *ByBitWS) Start() error",
							"Function where err" : "b.conn.ReadMessage",
							"Exchange" : "Bybit",
							"Error" : %s
						}`, err)
					log.Fatal()
				}
			} else {
				b.messageHandler(data)
			}
		}
	}()

	return nil
}

//	Необходим для приватных каналов
func (b *ByBitWS) Auth() {
	expires := time.Now().Unix()*1000 + 10000
	req := fmt.Sprintf("GET/realtime%d", expires)
	sig := hmac.New(sha256.New, []byte(b.cfg.SecretKey))
	sig.Write([]byte(req))
	signature := hex.EncodeToString(sig.Sum(nil))

	cmd := Cmd{
		Op: "auth",
		Args: []interface{}{
			b.cfg.ApiKey,
			//fmt.Sprintf("%v", expires),
			expires,
			signature,
		},
	}
	b.SendCmd(cmd)
}

func (b *ByBitWS) connect() {

	c, _, err := websocket.DefaultDialer.Dial(b.cfg.Addr, nil)
	if err != nil {
		log.Printf(`{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
						"File": "client.go",
						"Functions" : "(b *ByBitWS) connect()",
						"Function where err" : "websocket.DefaultDialer.Dial",
						"Exchange" : "Bybit",
						"Data" : [%s, nil],
						"Error" : %s
					}`, b.cfg.Addr, err)
		log.Fatal()
	}
	b.conn = c
	for _, cmd := range b.subscribeCmds {
		b.SendCmd(cmd)
	}
}

func (b *ByBitWS) closeAndReconnect() {
	b.conn.Close()
	b.connect()
}

func (b *ByBitWS) ping() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("BybitWs ping error: %v", r)
		}
	}()

	//	https://bybit-exchange.github.io/docs/spot/v3/#t-websocketauthentication
	err := b.conn.WriteMessage(websocket.TextMessage, []byte(`{"op":"ping"}`))
	if err != nil {
		log.Printf("BybitWs ping error: %v", err)
	}
}

func (b *ByBitWS) messageHandler(data []byte) {

	if b.cfg.DebugMode {
		log.Printf("BybitWs %v", string(data))
	}

	//	в ошибке нет необходимости, т.к. она выходит каждый раз, когда не найдет элемент
	typeJSON, _ := jsonparser.GetString(data, "type")

	switch typeJSON {
	case "delta":
		topic, _ := jsonparser.GetString(data, "topic")

		// 	пример: bookticker.BTCUSDT
		//	topicArr = ["bookticker", "BTCUSDT"]
		topicArr := strings.Split(topic, ".")
		switch topicArr[0] {
		case "bookticker":
			var bookTicker BookTicker
			err := json.Unmarshal(data, &bookTicker)
			if err != nil {
				log.Printf(`
					{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
						"File": "client.go",
						"Functions" : "(b *ByBitWS) messageHandler(data []byte)",
						"Function where err" : "json.Unmarshal",
						"Exchange" : "Bybit",
						"Comment" : %s to BookTicker struct,
						"Error" : %s
					}`, string(data), err)
				log.Fatal()
			}
			b.processBookTicker(topicArr[1], bookTicker)
		default:
			log.Printf(`
				{
					"Status" : "INFO",
					"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
					"File": "client.go",
					"Functions" : "(b *ByBitWS) messageHandler(data []byte)",
					"Exchange" : "Bybit",
					"Comment" : "Ответ от неизвестного канала"
					"Message" : %s
				}`, string(data))
			log.Fatal()
		}
	case "error":
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
				"File": "client.go",
				"Functions" : "(b *ByBitWS) messageHandler(data []byte)",
				"Exchange" : "Bybit",
				"Message" : %s
			}`, string(data))
		log.Fatal()
	default:
		opJSON, _ := jsonparser.GetString(data, "op")
		switch opJSON {
		case "subscribe":
			successJSON, _ := jsonparser.GetBoolean(data, "success")
			if !successJSON {
				log.Printf(`
					{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
						"File": "client.go",
						"Functions" : "(b *ByBitWS) messageHandler(data []byte)",
						"Exchange" : "Bybit",
						"Comment" : "Проблема с подпиской"
						"Message" : %s
					}`, string(data))
				log.Fatal()
			}
		case "pong":
		default:
			log.Printf(`
				{
					"Status" : "INFO",
					"Path to file" : "CCXT_BEYANG_BYBIT/spot/v3/ws",
					"File": "client.go",
					"Functions" : "(b *ByBitWS) messageHandler(data []byte)",
					"Exchange" : "Bybit",
					"Comment" : "не известный ответ от сервера"
					"Message" : %s
				}`, string(data))
			log.Fatal()
		}
	}
}
