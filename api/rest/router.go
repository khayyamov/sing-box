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

		r.GET(constant.GetStats, GetAllUsersStats)
		r.GET(constant.GetStat, GetAUserStat)

		r.POST(constant.AddToAll, AddToAll)
		r.POST(constant.DeleteToAll, DeleteToAll)

		err := r.Run(constant.ApiHost + ":" + constant.ApiPort)
		if err != nil {
			panic("Api exception:" + err.Error())
		}
	}()
}
