package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/constant"
)

func HandleApiRoutes() {
	go func() {
		r := gin.Default()
		r.POST(constant.RouteAddUserToVmess, AddUserToVmess)
		r.POST(constant.RouteAddUserToVless, AddUserToVless)
		r.POST(constant.RouteAddUserToHysteria2, AddUserToHysteria2)
		err := r.Run(":3000")
		if err != nil {
			panic("Api exception:" + err.Error())
		}
	}()
}
