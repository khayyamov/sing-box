package rq

type GlobalModel struct {
	Tag              string        `json:"tag"`
	ReplacementField []GlobalModel `json:"replacement_field"`
	UUID             string        `json:"uuid"`           // vless, vmess, tuic
	Name             string        `json:"name"`           // tuic, trojan, shadowtls, hysteria, hysteria2, shadowsocks_multi, shadowsocks_relay
	Password         string        `json:"password"`       // shadowtls, naive, trojan, hysteria2, tuic, shadowsocks_multi, shadowsocks_relay
	Flow             string        `json:"flow"`           // vless, vmess
	Auth             string        `json:"auth"`           // hysteria
	AuthString       string        `json:"auth_string"`    // hysteria
	ServerAddress    string        `json:"server_address"` // shadowsocks_relay
	ServerPort       uint16        `json:"server_port"`    // shadowsocks_relay
	AlterId          int           `json:"alter_id"`       // vmess
}
