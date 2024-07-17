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

func EditVmessUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.VMessPtr) == 0 {
		utils.ApiLogInfo("No Active VMessPtr outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.VMessUser{
			Name:    user.Name,
			UUID:    user.UUID,
			AlterId: user.AlterId,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VMessUser](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			utils.ApiLogError("User already exist: " + dbUser.UserJson)
			continue
		}
		_, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		box.EditUserInV2rayApi(user.UUID, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeVMess, delete)
		for i := range inbound.VMessPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.VMessPtr[i].Tag() == model.Tag {
							if len(model.Flow) > 0 {
								convertedUser.AlterId = model.AlterId
							}
							break
						}
					}
				}
				inbound.VMessPtr[i].Service.AddUser(
					common.MapIndexed([]option.VMessUser{convertedUser}, func(index int, it option.VMessUser) int {
						return len(inbound.VMessPtr[i].Users) + index
					}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) string {
						return it.UUID
					}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) int {
						return it.AlterId
					}))
				inbound.VMessPtr[i].Users = append(inbound.VMessPtr[i].Users, convertedUser)
			} else {
				inbound.VMessPtr[i].Service.DeleteUser(
					common.MapIndexed([]option.VMessUser{convertedUser}, func(index int, it option.VMessUser) int {
						return index
					}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) string {
						return it.UUID
					}), common.Map([]option.VMessUser{convertedUser}, func(it option.VMessUser) int {
						return it.AlterId
					}))
				for j := range newUsers {
					for k := range inbound.VMessPtr[i].Users {
						if newUsers[j].UUID == inbound.VMessPtr[i].Users[k].UUID {
							inbound.VMessPtr[i].Users = append(
								inbound.VMessPtr[i].Users[:k],
								inbound.VMessPtr[i].Users[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}
