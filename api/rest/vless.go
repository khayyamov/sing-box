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

func EditVlessUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	utils.CurrentInboundName = "Vless"
	if len(inbound.VLESSPtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.VLESSUser{
			Name: user.Name,
			UUID: user.UUID,
			Flow: user.Flow,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VLESSUser](convertedUser)
		for i := range inbound.VLESSPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.VLESSPtr[i].Tag() == model.Tag {
							if len(model.Name) > 0 {
								convertedUser.Name = model.Name
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
				_, err := uuid.FromString(user.UUID)
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
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.VLESSPtr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {
				_, err := uuid.FromString(user.UUID)
				if len(convertedUser.UUID) == 0 || err != nil {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.VLESSPtr[i].Tag() + "] User failed to delete uuid invalid")
					continue
				}
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

			box.EditUserInV2rayApi(convertedUser.Name, delete)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeVLESS, delete)
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
