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
)

func EditShadowsocksRelayUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.ShadowsocksRelayPtr) == 0 {
		log.Info("No Active ShadowsocksRelayPtr outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.ShadowsocksDestination{
			Name:     user.UUID,
			Password: user.Password,
			ServerOptions: option.ServerOptions{
				Server:     user.ServerAddress,
				ServerPort: user.ServerPort,
			},
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.ShadowsocksDestination](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			log.Error("User already exist: " + dbUser.UserJson)
			continue
		}
		_, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		box.EditUserInV2rayApi(user.UUID, delete)
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
						if newUsers[j].UUID == inbound.ShadowsocksRelayPtr[i].Destinations[k].Name {
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
