package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"net/http"
)

func AddUserToVmess(c *gin.Context) {
	var rq option.VMessUser
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []option.VMessUser{rq}
	err := inbound.VMessPtr.Service.AddUser(common.MapIndexed(newUsers, func(index int, it option.VMessUser) int {
		return index
	}), common.Map(newUsers, func(it option.VMessUser) string {
		return it.UUID
	}), common.Map(newUsers, func(it option.VMessUser) int {
		return it.AlterId
	}))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
