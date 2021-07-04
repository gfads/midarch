package messages

type AOR struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	Id        int    `json:"id"`
	ProxyName string `json:"proxy"`
}
