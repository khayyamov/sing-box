package rq

type GlobalModel struct {
	Username      string `json:"username"`       // naive
	Password      string `json:"password"`       // shadowtls, naive, trojan, hysteria2, tuic, shadowsocks_multi, shadowsocks_relay
	Name          string `json:"name"`           // vless, vmess, tuic, trojan, shadowtls, hysteria, hysteria2, shadowsocks_multi, shadowsocks_relay
	Uuid          string `json:"uuid"`           // vless, vmess, tuic
	Flow          string `json:"flow"`           // vless, vmess
	Auth          string `json:"auth"`           // hysteria
	AuthString    string `json:"auth_string"`    // hysteria
	ServerAddress string `json:"server_address"` // shadowsocks_relay
	ServerPort    uint16 `json:"server_port"`    // shadowsocks_relay

	AddToAll          bool `json:"add_to_all"`
	Hysteria          bool `json:"hysteria"`
	Hysteria2         bool `json:"hysteria2"`
	Naive             bool `json:"naive"`
	Shadowsocks_multi bool `json:"shadowsocks_multi"`
	Shadowsocks_relay bool `json:"shadowsocks_relay"`
	Shadowtls         bool `json:"shadowtls"`
	Trojan            bool `json:"trojan"`
	Tuic              bool `json:"tuic"`
	Vless             bool `json:"vless"`
	Vmess             bool `json:"vmess"`
}
