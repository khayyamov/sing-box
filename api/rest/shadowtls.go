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
	shadowtls "github.com/sagernet/sing-shadowtls"
)

func EditShadowtlsUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.ShadowTlsPtr) == 0 {
		log.Info("No Active ShadowTls outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := shadowtls.User{
			Name:     user.UUID,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[shadowtls.User](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			log.Error("User already exist: " + dbUser.UserJson)
			continue
		}
		_, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		box.EditUserInV2rayApi(user.UUID, delete)
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
						if newUsers[j].UUID == inbound.ShadowTlsPtr[i].Service.Users[k].Name {
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
