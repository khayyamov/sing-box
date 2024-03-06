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
	inbound.TrojanPtr.Users = append(inbound.TrojanPtr.Users, rq)
	inbound.TrojanPtr.Service.UpdateUsers(common.MapIndexed(inbound.TrojanPtr.Users, func(index int, it option.TrojanUser) int {
		return index
	}), common.Map(inbound.TrojanPtr.Users, func(it option.TrojanUser) string {
		return it.Password
	}))
}
