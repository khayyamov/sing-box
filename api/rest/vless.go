package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"net/http"
)

func AddUserToVless(c *gin.Context) {
	var rq option.VLESSUser
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []option.VLESSUser{rq}
	inbound.VLESSPtr.Service.AddUser(common.MapIndexed(newUsers, func(index int, _ option.VLESSUser) int {
		return index
	}), common.Map(newUsers, func(it option.VLESSUser) string {
		return it.UUID
	}), common.Map(newUsers, func(it option.VLESSUser) string {
		return it.Flow
	}))
	err := db.GetDb().AddVlessUser(newUsers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
