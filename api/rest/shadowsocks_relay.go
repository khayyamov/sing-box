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

func EditShadowsocksRelayUsers(c *gin.Context, newUsers []rq.GlobalModel, deletee bool) {
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
			if !deletee {
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
						common.MapIndexedString([]option.ShadowsocksDestination{convertedUser}, func(index any, user option.ShadowsocksDestination) string {
							return user.Name
						}), common.Map([]option.ShadowsocksDestination{convertedUser}, func(user option.ShadowsocksDestination) string {
							return user.Password
						}), common.Map([]option.ShadowsocksDestination{convertedUser}, option.ShadowsocksDestination.Build))
					inbound.ShadowsocksRelayPtr[i].Destinations[convertedUser.Name] = convertedUser
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.ShadowsocksRelayPtr[i].Tag() + "]   User already exist: " + dbUser.UserJson)
				}
			} else {
				if len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.ShadowsocksRelayPtr[i].Tag() + "]   User failed to deletee name invalid")
					continue
				}
				_ = inbound.ShadowsocksRelayPtr[i].Service.DeleteUsersWithPasswords(
					common.MapIndexedString([]option.ShadowsocksDestination{convertedUser}, func(index any, user option.ShadowsocksDestination) string {
						return user.Name
					}), common.Map([]option.ShadowsocksDestination{convertedUser}, func(user option.ShadowsocksDestination) string {
						return user.Password
					}))

				delete(inbound.ShadowsocksRelayPtr[i].Destinations, convertedUser.Name)
			}

			box.EditUserInV2rayApi(convertedUser.Name, deletee)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeShadowsocksRelay, deletee)
		}
	}
}
