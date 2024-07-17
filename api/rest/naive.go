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
	"github.com/sagernet/sing/common/auth"
)

func EditNaiveUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.NaivePtr) == 0 {
		utils.ApiLogInfo("No Active NaivePtr outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := auth.User{
			Username: user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[auth.User](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			utils.ApiLogError("User already exist: " + dbUser.UserJson)
			continue
		}
		box.EditUserInV2rayApi(user.Name, delete)
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
