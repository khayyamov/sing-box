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

func EditShadowsocksMultiUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	utils.CurrentInboundName = "ShadowsocksMulti"
	if len(inbound.ShadowsocksMultiPtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.ShadowsocksUser{
			Name:     user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.ShadowsocksUser](convertedUser)
		for i := range inbound.ShadowsocksMultiPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.ShadowsocksMultiPtr[i].Tag() == model.Tag {
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
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.ShadowsocksMultiPtr[i].Tag() + "]  User failed to add password invalid")
					continue
				}
				founded := false
				for _, inboundUsers := range inbound.ShadowsocksMultiPtr[i].Users {
					if inboundUsers.Name == convertedUser.Name {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.ShadowsocksUser](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.ShadowsocksMultiPtr[i].Tag() + "]  User Added: " + dbUser.UserJson)
					inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users, []option.ShadowsocksUser{convertedUser}...)
					_ = inbound.ShadowsocksMultiPtr[i].Service.AddUsersWithPasswords(
						common.MapIndexed([]option.ShadowsocksUser{convertedUser}, func(index int, user option.ShadowsocksUser) int {
							return len(inbound.ShadowsocksMultiPtr[i].Users) + index
						}), common.Map([]option.ShadowsocksUser{convertedUser}, func(user option.ShadowsocksUser) string {
							return user.Password
						}))
					inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users, convertedUser)
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.ShadowsocksMultiPtr[i].Tag() + "]  User already exist: " + dbUser.UserJson)
				}
			} else {
				_ = inbound.ShadowsocksMultiPtr[i].Service.DeleteUsersWithPasswords(
					common.MapIndexed([]option.ShadowsocksUser{convertedUser}, func(index int, user option.ShadowsocksUser) int {
						return index
					}), common.Map([]option.ShadowsocksUser{convertedUser}, func(user option.ShadowsocksUser) string {
						return user.Password
					}))
				for j := range newUsers {
					for k := range inbound.ShadowsocksMultiPtr[i].Users {
						if newUsers[j].Name == inbound.ShadowsocksMultiPtr[i].Users[k].Name {
							inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users[:k],
								inbound.ShadowsocksMultiPtr[i].Users[k+1:]...)
							break
						}
					}
				}
			}

			box.EditUserInV2rayApi(convertedUser.Name, delete)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeShadowsocksMulti, delete)
		}
	}
}
