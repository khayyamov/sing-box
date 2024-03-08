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

func AddUserToShadowsocksMulti(c *gin.Context) {
	var rq option.ShadowsocksUser
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []option.ShadowsocksUser{rq}
	for i := range inbound.ShadowsocksMultiPtr {
		err := inbound.ShadowsocksMultiPtr[i].Service.AddUsersWithPasswords(common.MapIndexed(newUsers, func(index int, user option.ShadowsocksUser) int {
			return index
		}), common.Map(newUsers, func(user option.ShadowsocksUser) string {
			return user.Password
		}))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().AddUserToDb(users, C.TypeShadowsocksMulti)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
