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
	"github.com/sagernet/sing-box/option"
)

func EditHysteriaUsers(c *gin.Context, newUsers []rq.GlobalModel, deletee bool) {
	utils.CurrentInboundName = "Hysteria"
	if len(inbound.HysteriaPtr) == 0 {
		utils.ApiLogInfo("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for _, user := range newUsers {
		convertedUser := option.HysteriaUser{
			Name:       user.Name,
			AuthString: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.HysteriaUser](convertedUser)
		for i := range inbound.HysteriaPtr {
			userList := make([]string, 0, len(newUsers))
			userPasswordList := make([]string, 0, len(newUsers))
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.HysteriaPtr[i].Tag() == model.Tag {
						if len(model.Name) > 0 {
							convertedUser.Name = model.Name
						}
						if len(model.Password) > 0 {
							convertedUser.AuthString = model.Password
						}
						break
					}
				}
			}
			if !deletee {

				if len(convertedUser.Auth) == 0 && len(convertedUser.AuthString) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.HysteriaPtr[i].Tag() + "] User failed to add Auth or AuthString invalid")
					continue
				}
				founded := false
				for _, inboundUsers := range inbound.TrojanPtr[i].Users {
					if inboundUsers.Name == convertedUser.Name {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.HysteriaUser](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.HysteriaPtr[i].Tag() + "] User Added: " + dbUser.UserJson)
					userList = append(userList, convertedUser.Name)
					userPasswordList = append(userPasswordList, convertedUser.AuthString)
					inbound.HysteriaPtr[i].Service.AddUser(userList, userPasswordList)
					inbound.HysteriaPtr[i].Users[convertedUser.Name] = convertedUser
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.HysteriaPtr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {
				userList = append(userList, convertedUser.Name)
				userPasswordList = append(userPasswordList, convertedUser.AuthString)
				inbound.HysteriaPtr[i].Service.DeleteUser(userList, userPasswordList)
				delete(inbound.HysteriaPtr[i].Users, convertedUser.Name)
			}
			box.EditUserInV2rayApi(convertedUser.Name, deletee)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeHysteria, deletee)
		}
	}
}
