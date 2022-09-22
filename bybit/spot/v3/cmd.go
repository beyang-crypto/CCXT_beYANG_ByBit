package v3

//	Необходим для удобного создания подписок
type Cmd struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}
