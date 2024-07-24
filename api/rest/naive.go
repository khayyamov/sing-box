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
	utils.CurrentInboundName = "Naive"
	if len(inbound.NaivePtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := auth.User{
			Username: user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[auth.User](convertedUser)
		for i := range inbound.NaivePtr {
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.NaivePtr[i].Tag() == model.Tag {
						if len(model.Name) > 0 {
							convertedUser.Username = model.Name
						}
						if len(model.Password) > 0 {
							convertedUser.Password = model.Password
						}
						break
					}
				}
			}
			if !delete {
				if len(convertedUser.Password) == 0 || len(convertedUser.Username) == 0 {
					continue
				}
				if !inbound.NaivePtr[i].Authenticator.Verify(convertedUser.Username, convertedUser.Password) {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[auth.User](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.NaivePtr[i].Tag() + "] User Added: " + dbUser.UserJson)
					inbound.NaivePtr[i].Authenticator.AddUserToAuthenticator([]auth.User{convertedUser})
				} else {
					utils.ApiLogInfo("NaivePtr: User already exist: " + dbUser.UserJson)
				}
			} else {
				if len(convertedUser.Username) == 0 {
					utils.ApiLogError("NaivePtr: User failed to delete Username invalid")
					continue
				}
				inbound.NaivePtr[i].Authenticator.DeleteUserToAuthenticator([]auth.User{convertedUser})
			}

			box.EditUserInV2rayApi(convertedUser.Username, delete)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeNaive, delete)
		}
	}
}
