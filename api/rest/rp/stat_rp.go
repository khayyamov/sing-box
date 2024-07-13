package rp

type StatRpList struct {
	ListUserStats []StatRp `json:"list_user_stats"`
}

type StatRp struct {
	Id       string `json:"id"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}
