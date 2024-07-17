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
	if len(inbound.ShadowsocksRelayPtr) == 0 {
		utils.ApiLogInfo("No Active ShadowsocksRelayPtr outbound found to add users to it")
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
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			utils.ApiLogError("User already exist: " + dbUser.UserJson)
			continue
		}
		box.EditUserInV2rayApi(user.Name, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeShadowsocksRelay, delete)
		for i := range inbound.ShadowsocksRelayPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.ShadowsocksRelayPtr[i].Tag() == model.Tag {
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
					continue
				}
				_ = inbound.ShadowsocksRelayPtr[i].Service.AddUsersWithPasswords(
					common.MapIndexed([]option.ShadowsocksDestination{convertedUser}, func(index int, user option.ShadowsocksDestination) int {
						return len(inbound.ShadowsocksRelayPtr[i].Destinations) + index
					}), common.Map([]option.ShadowsocksDestination{convertedUser}, func(user option.ShadowsocksDestination) string {
						return user.Password
					}), common.Map([]option.ShadowsocksDestination{convertedUser}, option.ShadowsocksDestination.Build))
				inbound.ShadowsocksRelayPtr[i].Destinations = append(inbound.ShadowsocksRelayPtr[i].Destinations, []option.ShadowsocksDestination{convertedUser}...)
			} else {
				_ = inbound.ShadowsocksRelayPtr[i].Service.DeleteUsersWithPasswords(
					common.MapIndexed([]option.ShadowsocksDestination{convertedUser}, func(index int, user option.ShadowsocksDestination) int {
						return len(inbound.ShadowsocksRelayPtr[i].Destinations) + index
					}), common.Map([]option.ShadowsocksDestination{convertedUser}, func(user option.ShadowsocksDestination) string {
						return user.Password
					}))

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
