package rest

import (
	"github.com/gin-gonic/gin"
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

func EditTrojanUsers(c *gin.Context, newUsers []rq.GlobalModel, deletee bool) {
	utils.CurrentInboundName = "Trojan"
	if len(inbound.TrojanPtr) == 0 {
		utils.ApiLogInfo("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.TrojanUser{
			Name:     user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.TrojanUser](convertedUser)
		for i := range inbound.TrojanPtr {
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.TrojanPtr[i].Tag() == model.Tag {
						if len(model.Name) > 0 {
							convertedUser.Name = model.Name
						}
						if len(model.Password) > 0 {
							convertedUser.Password = model.Password
						}
						break
					}
				}
			}
			if !deletee {
				if len(convertedUser.Password) == 0 || len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.TrojanPtr[i].Tag() + "] User failed to add name or password invalid: " + dbUser.UserJson)
					continue
				}
				founded := false
				for _, inboundUsers := range inbound.TrojanPtr[i].Users {
					if inboundUsers.Name == convertedUser.Name {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.TrojanUser](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.TrojanPtr[i].Tag() + "] User Added: " + dbUser.UserJson)
					_ = inbound.TrojanPtr[i].Service.AddUser(
						common.MapIndexedString([]option.TrojanUser{convertedUser}, func(index any, it option.TrojanUser) string {
							return it.Name
						}), common.Map([]option.TrojanUser{convertedUser}, func(it option.TrojanUser) string {
							return it.Password
						}))
					inbound.TrojanPtr[i].Users[convertedUser.Name] = convertedUser
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.TrojanPtr[i].Tag() + "]   User already exist: " + dbUser.UserJson)
				}
			} else {
				if len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.TrojanPtr[i].Tag() + "]   User failed to deletee name invalid")
					continue
				}
				_ = inbound.TrojanPtr[i].Service.DeleteUser(
					common.MapIndexedString([]option.TrojanUser{convertedUser}, func(index any, it option.TrojanUser) string {
						return it.Name
					}), common.Map([]option.TrojanUser{convertedUser}, func(it option.TrojanUser) string {
						return it.Password
					}))
				delete(inbound.TrojanPtr[i].Users, convertedUser.Name)
			}

			box.EditUserInV2rayApi(convertedUser.Name, deletee)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeTrojan, deletee)
		}
	}
}
