package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"net/http"
)

func AddUserToTrojan(c *gin.Context) {
	var rq option.TrojanUser
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []option.TrojanUser{rq}
	for i := range inbound.TrojanPtr {
		err := inbound.TrojanPtr[i].Service.AddUser(common.MapIndexed(newUsers, func(index int, it option.TrojanUser) int {
			return index
		}), common.Map(newUsers, func(it option.TrojanUser) string {
			return it.Password
		}))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().AddUserToDb(users, C.TypeTrojan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
