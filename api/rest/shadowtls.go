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
	shadowtls "github.com/sagernet/sing-shadowtls"
)

func EditShadowtlsUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.ShadowTlsPtr) == 0 {
		utils.ApiLogInfo("No Active ShadowTlsPtr outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := shadowtls.User{
			Name:     user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[shadowtls.User](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			utils.ApiLogError("User already exist: " + dbUser.UserJson)
			continue
		}
		box.EditUserInV2rayApi(user.Name, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeShadowTLS, delete)
		for i := range inbound.ShadowTlsPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.ShadowTlsPtr[i].Tag() == model.Tag {
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
				inbound.ShadowTlsPtr[i].Service.AddUser([]shadowtls.User{convertedUser})
				inbound.ShadowTlsPtr[i].Service.Users = append(inbound.ShadowTlsPtr[i].Service.Users, convertedUser)
			} else {

				_ = inbound.ShadowTlsPtr[i].Service.DeleteUser([]shadowtls.User{convertedUser})
				for j := range newUsers {
					for k := range inbound.ShadowTlsPtr[i].Service.Users {
						if newUsers[j].Name == inbound.ShadowTlsPtr[i].Service.Users[k].Name {
							inbound.ShadowTlsPtr[i].Service.Users = append(
								inbound.ShadowTlsPtr[i].Service.Users[:k],
								inbound.ShadowTlsPtr[i].Service.Users[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}
