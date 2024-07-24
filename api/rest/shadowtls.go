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
	utils.CurrentInboundName = "ShadowTls"
	if len(inbound.ShadowTlsPtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := shadowtls.User{
			Name:     user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[shadowtls.User](convertedUser)
		for i := range inbound.ShadowTlsPtr {
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.ShadowTlsPtr[i].Tag() == model.Tag {
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

			if !delete {
				if len(convertedUser.Password) == 0 || len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.ShadowTlsPtr[i].Tag() + "]   User failed to add name or password invalid: " + dbUser.UserJson)
					continue
				}

				founded := false
				for _, inboundUsers := range inbound.ShadowTlsPtr[i].Service.Users {
					if inboundUsers.Name == convertedUser.Name {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[shadowtls.User](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.ShadowTlsPtr[i].Tag() + "]   User Added: " + dbUser.UserJson)
					inbound.ShadowTlsPtr[i].Service.AddUser([]shadowtls.User{convertedUser})
					inbound.ShadowTlsPtr[i].Service.Users = append(inbound.ShadowTlsPtr[i].Service.Users, convertedUser)
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.ShadowTlsPtr[i].Tag() + "]   User already exist: " + dbUser.UserJson)
				}
			} else {

				if len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.ShadowTlsPtr[i].Tag() + "]   User failed to delete name invalid")
					continue
				}
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

			box.EditUserInV2rayApi(convertedUser.Name, delete)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeShadowTLS, delete)
		}
	}
}
