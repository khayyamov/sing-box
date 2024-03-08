package rest

import (
	"github.com/gin-gonic/gin"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
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

	for i := range inbound.VLESSPtr {
		inbound.VLESSPtr[i].Service.AddUser(common.MapIndexed(newUsers, func(index int, _ option.VLESSUser) int {
			return index
		}), common.Map(newUsers, func(it option.VLESSUser) string {
			return it.UUID
		}), common.Map(newUsers, func(it option.VLESSUser) string {
			return it.Flow
		}))
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().AddUserToDb(users, C.TypeVLESS)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
func AddVlessInbound(c *gin.Context) {
	var rq option.VLESSInboundOptions
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := box.AddInbound(box.Options{
		Options: option.Options{
			Inbounds: []option.Inbound{
				{
					Type:         C.TypeVLESS,
					VLESSOptions: rq,
				},
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
