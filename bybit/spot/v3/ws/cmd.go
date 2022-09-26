package ws

//	Необходим для удобного создания подписок
type Cmd struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}
