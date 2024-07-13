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

func EditTrojanUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.TrojanPtr) == 0 {
		log.Info("No Active Trojan outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.TrojanUser{
			Name:     user.UUID,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.TrojanUser](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			log.Error("User already exist: " + dbUser.UserJson)
			continue
		}
		_, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		box.EditUserInV2rayApi(user.UUID, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeTrojan, delete)
		for i := range inbound.TrojanPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.TrojanPtr[i].Tag() == model.Tag {
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
				_ = inbound.TrojanPtr[i].Service.AddUser(
					common.MapIndexed([]option.TrojanUser{convertedUser}, func(index int, it option.TrojanUser) int {
						return len(inbound.TrojanPtr[i].Users) + index
					}), common.Map([]option.TrojanUser{convertedUser}, func(it option.TrojanUser) string {
						return it.Password
					}))
				inbound.TrojanPtr[i].Users = append(inbound.TrojanPtr[i].Users, convertedUser)
			} else {
				_ = inbound.TrojanPtr[i].Service.DeleteUser(
					common.MapIndexed([]option.TrojanUser{convertedUser}, func(index int, it option.TrojanUser) int {
						return index
					}), common.Map([]option.TrojanUser{convertedUser}, func(it option.TrojanUser) string {
						return it.Password
					}))
				for j := range newUsers {
					for k := range inbound.TrojanPtr[i].Users {
						if newUsers[j].UUID == inbound.TrojanPtr[i].Users[k].Name {
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
