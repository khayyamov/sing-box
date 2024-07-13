package rest

import (
	"github.com/gin-gonic/gin"
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

func EditShadowsocksMultiUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.ShadowsocksMultiPtr) == 0 {
		log.Info("No Active Vless outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.ShadowsocksUser{
			Name:     user.UUID,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.ShadowsocksUser](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			log.Error("User already exist: " + dbUser.UserJson)
			continue
		}
		box.EditUserInV2rayApi(user.UUID, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeShadowsocksMulti, delete)
		for i := range inbound.ShadowsocksMultiPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.ShadowsocksMultiPtr[i].Tag() == model.Tag {
							if len(model.Password) > 0 {
								convertedUser.Password = model.Password
							}
							break
						}
					}
				}
				if len(convertedUser.Password) == 0 {
					continue
				}
				inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users, []option.ShadowsocksUser{convertedUser}...)
				_ = inbound.ShadowsocksMultiPtr[i].Service.AddUsersWithPasswords(
					common.MapIndexed([]option.ShadowsocksUser{convertedUser}, func(index int, user option.ShadowsocksUser) int {
						return len(inbound.ShadowsocksMultiPtr[i].Users) + index
					}), common.Map([]option.ShadowsocksUser{convertedUser}, func(user option.ShadowsocksUser) string {
						return user.Password
					}))
				inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users, convertedUser)
			} else {
				_ = inbound.ShadowsocksMultiPtr[i].Service.DeleteUsersWithPasswords(
					common.MapIndexed([]option.ShadowsocksUser{convertedUser}, func(index int, user option.ShadowsocksUser) int {
						return index
					}), common.Map([]option.ShadowsocksUser{convertedUser}, func(user option.ShadowsocksUser) string {
						return user.Password
					}))
				for j := range newUsers {
					for k := range inbound.ShadowsocksMultiPtr[i].Users {
						if newUsers[j].UUID == inbound.ShadowsocksMultiPtr[i].Users[k].Name {
							inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users[:k],
								inbound.ShadowsocksMultiPtr[i].Users[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}