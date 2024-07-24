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
)

func EditVmessUsers(c *gin.Context, newUsers []rq.GlobalModel, deletee bool) {
	utils.CurrentInboundName = "Vmess"
	if len(inbound.VMessPtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.VMessUser{
			Name:    user.Name,
			UUID:    user.UUID,
			AlterId: user.AlterId,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VMessUser](convertedUser)
		for i := range inbound.VMessPtr {
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.VMessPtr[i].Tag() == model.Tag {
						if len(model.Name) > 0 {
							convertedUser.Name = model.Name
						}
						if len(model.UUID) > 0 {
							convertedUser.UUID = model.UUID
						}
						if len(model.Flow) > 0 {
							convertedUser.AlterId = model.AlterId
						}
						break
					}
				}
			}
			if !deletee {
				_, err := uuid.FromString(user.UUID)
				if len(convertedUser.UUID) == 0 || err != nil {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.VMessPtr[i].Tag() + "] User failed to add uuid invalid")
					continue
				}
				founded := false
				for _, inboundUsers := range inbound.VMessPtr[i].Users {
					if inboundUsers.UUID == convertedUser.UUID {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VMessUser](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.VMessPtr[i].Tag() + "] User Added: " + dbUser.UserJson)
					inbound.VMessPtr[i].Service.AddUser(
						common.MapIndexedString([]option.VMessUser{convertedUser}, func(index any, it option.VMessUser) string {
							return it.UUID
						}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) string {
							return it.UUID
						}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) int {
							return it.AlterId
						}))
					inbound.VMessPtr[i].Users[convertedUser.UUID] = convertedUser
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.VMessPtr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {
				_, err := uuid.FromString(user.UUID)
				if len(convertedUser.UUID) == 0 || err != nil {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.VMessPtr[i].Tag() + "] User failed to deletee uuid invalid")
					continue
				}
				inbound.VMessPtr[i].Service.DeleteUser(
					common.MapIndexedString([]option.VMessUser{convertedUser}, func(index any, it option.VMessUser) string {
						return it.UUID
					}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) string {
						return it.UUID
					}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) int {
						return it.AlterId
					}))
				delete(inbound.VMessPtr[i].Users, convertedUser.UUID)
			}

			box.EditUserInV2rayApi(convertedUser.Name, deletee)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeVMess, deletee)
		}
	}
}
