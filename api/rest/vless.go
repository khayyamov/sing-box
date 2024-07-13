package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/api/db"
	"github.com/sagernet/sing-box/api/db/entity"
	"github.com/sagernet/sing-box/api/rest/rq"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"net/http"
)

func EditVlessUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.VLESSPtr) == 0 {
		log.Info("No Active Vless outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.VLESSUser{
			Name: user.UUID,
			UUID: user.UUID,
			Flow: user.Flow,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VLESSUser](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			log.Error("User already exist: " + dbUser.UserJson)
			continue
		}
		_, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		box.EditUserInV2rayApi(user.UUID, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeVLESS, delete)
		for i := range inbound.VLESSPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.VLESSPtr[i].Tag() == model.Tag {
							if len(model.Flow) > 0 {
								convertedUser.Flow = model.Flow
							}
							break
						}
					}
				}
				inbound.VLESSPtr[i].Service.AddUser(
					common.MapIndexed([]option.VLESSUser{convertedUser}, func(index int, it option.VLESSUser) int {
						return len(inbound.VLESSPtr[i].Users) + index
					}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
						return it.UUID
					}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
						return it.Flow
					}))
				inbound.VLESSPtr[i].Users = append(inbound.VLESSPtr[i].Users, convertedUser)
			} else {
				inbound.VLESSPtr[i].Service.DeleteUser(
					common.MapIndexed([]option.VLESSUser{convertedUser}, func(index int, it option.VLESSUser) int {
						return index
					}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
						return it.UUID
					}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
						return it.Flow
					}))
				for j := range newUsers {
					for k := range inbound.VLESSPtr[i].Users {
						if newUsers[j].UUID == inbound.VLESSPtr[i].Users[k].UUID {
							inbound.VLESSPtr[i].Users = append(
								inbound.VLESSPtr[i].Users[:k],
								inbound.VLESSPtr[i].Users[k+1:]...)
							break
						}
					}
				}
			}
		}
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
