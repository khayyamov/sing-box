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

		r.POST(constant.AddToAll, AddToAll)
		r.POST(constant.DeleteToAll, DeleteToAll)

		r.POST(constant.RouteAddUserToVmess, AddUserToVmess)
		r.POST(constant.RouteDeleteUserToVmess, DeleteUserToVmess)

		r.POST(constant.RouteAddUserToVless, AddUserToVless)
		r.POST(constant.RouteDeleteUserToVless, DeleteUserToVless)
		r.POST(constant.RouteAddVlessInbound, AddVlessInbound)

		r.POST(constant.RouteAddUserToHysteria2, AddUserToHysteria2)
		r.POST(constant.RouteDeleteUserToHysteria2, DeleteUserToHysteria2)

		r.POST(constant.RouteAddUserToHysteria, AddUserToHysteria)
		r.POST(constant.RouteDeleteUserToHysteria, DeleteUserToHysteria)

		r.POST(constant.RouteAddUserToShadowtls, AddUserToShadowtls)
		r.POST(constant.RouteDeleteUserToShadowtls, DeleteUserToShadowtls)

		r.POST(constant.RouteAddUserToNaive, AddUserTonNaive)
		r.POST(constant.RouteDeleteUserToNaive, DeleteUserTonNaive)

		r.POST(constant.RouteAddUserToTuic, AddUserToTuic)
		r.POST(constant.RouteDeleteUserToTuic, DeleteUserToTuic)

		r.POST(constant.RouteAddUserToTrojan, AddUserToTrojan)
		r.POST(constant.RouteDeleteUserToTrojan, DeleteUserToTrojan)

		r.POST(constant.RouteAddUserToShadowsocksMulti, AddUserToShadowsocksMulti)
		r.POST(constant.RouteDeleteUserToShadowsocksMulti, DeleteUserToShadowsocksMulti)

		err := r.Run(":3000")
		if err != nil {
			panic("Api exception:" + err.Error())
		}
	}()
}
