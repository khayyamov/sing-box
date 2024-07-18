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

func EditShadowsocksRelayUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	utils.CurrentInboundName = "ShadowsocksRelay"
	if len(inbound.ShadowsocksRelayPtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.ShadowsocksDestination{
			Name:     user.Name,
			Password: user.Password,
			ServerOptions: option.ServerOptions{
				Server:     user.ServerAddress,
				ServerPort: user.ServerPort,
			},
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.ShadowsocksDestination](convertedUser)
		for i := range inbound.ShadowsocksRelayPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.ShadowsocksRelayPtr[i].Tag() == model.Tag {
							if len(model.Name) > 0 {
								convertedUser.Name = model.Name
							}
							if len(model.Password) > 0 {
								convertedUser.Password = model.Password
							}
							if model.ServerPort > 0 {
								convertedUser.ServerPort = model.ServerPort
							}
							if len(model.ServerAddress) > 0 {
								convertedUser.Server = model.ServerAddress
							}
							break
						}
					}
				}

				if len(convertedUser.Password) == 0 || len(convertedUser.Server) == 0 || convertedUser.ServerPort == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.ShadowsocksRelayPtr[i].Tag() + "]   User failed to add password or Server or ServerPort invalid")
					continue
				}

				founded := false
				for _, inboundUsers := range inbound.ShadowsocksRelayPtr[i].Destinations {
					if inboundUsers.Name == convertedUser.Name {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.ShadowsocksDestination](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.ShadowsocksRelayPtr[i].Tag() + "]   User Added: " + dbUser.UserJson)
					_ = inbound.ShadowsocksRelayPtr[i].Service.AddUsersWithPasswords(
						common.MapIndexed([]option.ShadowsocksDestination{convertedUser}, func(index int, user option.ShadowsocksDestination) int {
							return len(inbound.ShadowsocksRelayPtr[i].Destinations) + index
						}), common.Map([]option.ShadowsocksDestination{convertedUser}, func(user option.ShadowsocksDestination) string {
							return user.Password
						}), common.Map([]option.ShadowsocksDestination{convertedUser}, option.ShadowsocksDestination.Build))
					inbound.ShadowsocksRelayPtr[i].Destinations = append(inbound.ShadowsocksRelayPtr[i].Destinations, []option.ShadowsocksDestination{convertedUser}...)
					box.EditUserInV2rayApi(user.Name, delete)
					db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeShadowsocksRelay, delete)
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.ShadowsocksRelayPtr[i].Tag() + "]   User already exist: " + dbUser.UserJson)
				}
			} else {
				if len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.ShadowsocksRelayPtr[i].Tag() + "]   User failed to delete name invalid")
					continue
				}
				_ = inbound.ShadowsocksRelayPtr[i].Service.DeleteUsersWithPasswords(
					common.MapIndexed([]option.ShadowsocksDestination{convertedUser}, func(index int, user option.ShadowsocksDestination) int {
						return len(inbound.ShadowsocksRelayPtr[i].Destinations) + index
					}), common.Map([]option.ShadowsocksDestination{convertedUser}, func(user option.ShadowsocksDestination) string {
						return user.Password
					}))
				box.EditUserInV2rayApi(user.Name, delete)
				for j := range newUsers {
					for k := range inbound.ShadowsocksRelayPtr[i].Destinations {
						if newUsers[j].Name == inbound.ShadowsocksRelayPtr[i].Destinations[k].Name {
							inbound.ShadowsocksRelayPtr[i].Destinations = append(
								inbound.ShadowsocksRelayPtr[i].Destinations[:k],
								inbound.ShadowsocksRelayPtr[i].Destinations[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}
