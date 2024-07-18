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

func EditTrojanUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	utils.CurrentInboundName = "Trojan"
	if len(inbound.TrojanPtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.TrojanUser{
			Name:     user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.TrojanUser](convertedUser)
		for i := range inbound.TrojanPtr {
			if !delete {
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
						common.MapIndexed([]option.TrojanUser{convertedUser}, func(index int, it option.TrojanUser) int {
							return len(inbound.TrojanPtr[i].Users) + index
						}), common.Map([]option.TrojanUser{convertedUser}, func(it option.TrojanUser) string {
							return it.Password
						}))
					inbound.TrojanPtr[i].Users = append(inbound.TrojanPtr[i].Users, convertedUser)
					box.EditUserInV2rayApi(user.Name, delete)
					db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeTrojan, delete)
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.TrojanPtr[i].Tag() + "]   User already exist: " + dbUser.UserJson)
				}
			} else {
				if len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.TrojanPtr[i].Tag() + "]   User failed to delete name invalid")
					continue
				}
				_ = inbound.TrojanPtr[i].Service.DeleteUser(
					common.MapIndexed([]option.TrojanUser{convertedUser}, func(index int, it option.TrojanUser) int {
						return index
					}), common.Map([]option.TrojanUser{convertedUser}, func(it option.TrojanUser) string {
						return it.Password
					}))
				box.EditUserInV2rayApi(user.Name, delete)
				for j := range newUsers {
					for k := range inbound.TrojanPtr[i].Users {
						if newUsers[j].Name == inbound.TrojanPtr[i].Users[k].Name {
							inbound.TrojanPtr[i].Users = append(
								inbound.TrojanPtr[i].Users[:k],
								inbound.TrojanPtr[i].Users[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}
