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

func AddUserToShadowsocksRelay(c *gin.Context) {
	var rq option.ShadowsocksDestination
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []option.ShadowsocksDestination{rq}
	for i := range inbound.ShadowsocksRelayPtr {
		err := inbound.ShadowsocksRelayPtr[i].Service.AddUsersWithPasswords(common.MapIndexed(newUsers, func(index int, user option.ShadowsocksDestination) int {
			return index
		}), common.Map(newUsers, func(user option.ShadowsocksDestination) string {
			return user.Password
		}), common.Map(newUsers, option.ShadowsocksDestination.Build))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().AddUserToDb(users, C.TypeShadowsocksRelay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
