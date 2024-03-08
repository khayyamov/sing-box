package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/constant"
	"github.com/sagernet/sing-box/api/db/mysql_config"
)

func HandleApiRoutes() {
	go func() {
		r := gin.Default()
		mysql_config.MySqlInstance()
		r.POST(constant.RouteAddUserToVmess, AddUserToVmess)
		r.POST(constant.RouteAddUserToVless, AddUserToVless)
		r.POST(constant.RouteAddVlessInbound, AddVlessInbound)
		r.POST(constant.RouteAddUserToHysteria2, AddUserToHysteria2)
		r.POST(constant.RouteAddUserToHysteria, AddUserToHysteria)
		r.POST(constant.RouteAddUserToShadowtls, AddUserToShadowtls)
		r.POST(constant.RouteAddUserToNaive, AddUserTonNaive)
		r.POST(constant.RouteAddUserToTuic, AddUserToTuic)
		r.POST(constant.RouteAddUserToTrojan, AddUserToTrojan)
		r.POST(constant.RouteAddUserToShadowsocksMulti, AddUserToShadowsocksMulti)
		r.POST(constant.RouteAddUserToShadowsocksRelay, AddUserToShadowsocksRelay)
		err := r.Run(":3000")
		if err != nil {
			panic("Api exception:" + err.Error())
		}
	}()
}
