package rest

import (
	"github.com/gin-gonic/gin"
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
	err := inbound.TrojanPtr.Service.AddUser(common.MapIndexed(newUsers, func(index int, it option.TrojanUser) int {
		return index
	}), common.Map(inbound.TrojanPtr.Users, func(it option.TrojanUser) string {
		return it.Password
	}))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {

	}
}
