package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/api/db"
	"github.com/sagernet/sing-box/api/db/entity"
	"github.com/sagernet/sing-box/api/rest/rq"
	"github.com/sagernet/sing-box/api/utils"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"net/http"
)

func EditVlessUsers(c *gin.Context, newUsers []rq.GlobalModel, deletee bool) {
	utils.CurrentInboundName = "Vless"
	for _, user := range newUsers {
		for i := range inbound.VLESSPtr {
			convertedUser := option.VLESSUser{
				Name: user.Name,
				UUID: user.UUID,
				Flow: user.Flow,
			}
			dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VLESSUser](convertedUser)
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.VLESSPtr[i].Tag() == model.Tag {
						if len(model.Name) > 0 {
							convertedUser.Name = model.Name
						} else {
							convertedUser.Name = convertedUser.Name
						}
						if len(model.UUID) > 0 {
							convertedUser.UUID = model.UUID
						}
						if len(model.Flow) > 0 {
							convertedUser.Flow = model.Flow
						}
						break
					}
				}
			}
			if !deletee {
				_, err := uuid.FromString(convertedUser.UUID)
				if len(convertedUser.UUID) == 0 || err != nil {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.VLESSPtr[i].Tag() + "] User failed to add name or password invalid: " + dbUser.UserJson)
					continue
				}
				founded := false
				for _, inboundUsers := range inbound.VLESSPtr[i].Users {
					if inboundUsers.UUID == convertedUser.UUID {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VLESSUser](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.VLESSPtr[i].Tag() + "] User Added: " + dbUser.UserJson)
					inbound.VLESSPtr[i].Service.AddUser(common.MapIndexedString([]option.VLESSUser{convertedUser}, func(index any, it option.VLESSUser) string {
						return it.UUID
					}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
						return it.UUID
					}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
						return it.Flow
					}))
					inbound.VLESSPtr[i].Users[convertedUser.UUID] = convertedUser
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.VLESSPtr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {
				_, err := uuid.FromString(convertedUser.UUID)
				if len(convertedUser.UUID) == 0 || err != nil {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.VLESSPtr[i].Tag() + "] User failed to delete uuid invalid")
					continue
				}
				inbound.VLESSPtr[i].Service.DeleteUser(common.MapIndexedString([]option.VLESSUser{convertedUser}, func(index any, it option.VLESSUser) string {
					return it.UUID
				}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
					return it.UUID
				}), common.Map([]option.VLESSUser{convertedUser}, func(it option.VLESSUser) string {
					return it.Flow
				}))
				delete(inbound.VLESSPtr[i].Users, convertedUser.UUID)
			}

			box.EditUserInV2rayApi(convertedUser.Name, deletee)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeVLESS, deletee)
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
