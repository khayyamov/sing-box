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
	"github.com/sagernet/sing/common/auth"
)

func EditNaiveUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.NaivePtr) == 0 {
		log.Info("No Active NaivePtr outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := auth.User{
			Username: user.UUID,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[auth.User](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			log.Error("User already exist: " + dbUser.UserJson)
			continue
		}
		_, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		box.EditUserInV2rayApi(user.UUID, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeNaive, delete)
		for i := range inbound.NaivePtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.NaivePtr[i].Tag() == model.Tag {
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
				inbound.NaivePtr[i].Authenticator.AddUserToAuthenticator([]auth.User{convertedUser})
			} else {
				inbound.NaivePtr[i].Authenticator.DeleteUserToAuthenticator([]auth.User{convertedUser})
			}
		}
	}
}
