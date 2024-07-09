package rest

import (
	box "github.com/sagernet/sing-box"
)

func AddUserToV2rayApi(user string) {
	if box.BoxInstance.Router().V2RayServer() != nil {
		box.BoxInstance.Router().V2RayServer().StatsService().AddUser(user)
	}
}
