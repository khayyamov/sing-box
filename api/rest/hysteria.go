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

func EditHysteriaUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	utils.CurrentInboundName = "Hysteria"
	if len(inbound.HysteriaPtr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for index, user := range newUsers {
		convertedUser := option.HysteriaUser{
			Name:       user.Name,
			AuthString: user.AuthString,
			Auth:       []byte(user.Auth),
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.HysteriaUser](convertedUser)
		for i := range inbound.HysteriaPtr {
			userList := make([]int, 0, len(newUsers))
			userNameList := make([]string, 0, len(newUsers))
			userPasswordList := make([]string, 0, len(newUsers))
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.HysteriaPtr[i].Tag() == model.Tag {
							if len(model.Name) > 0 {
								convertedUser.Name = model.Name
							}
							if len(model.Auth) > 0 {
								convertedUser.Auth = []byte(model.Auth)
							}
							if len(model.AuthString) > 0 {
								convertedUser.AuthString = model.AuthString
							}
							break
						}
					}
				}

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
					userList = append(userList, index)
					userNameList = append(userNameList, user.Name)
					userPasswordList = append(userPasswordList, user.Password)
					inbound.HysteriaPtr[i].Service.AddUser(userList, userPasswordList)
					inbound.HysteriaPtr[i].UserNameList = append(inbound.HysteriaPtr[i].UserNameList, convertedUser.Name)
					box.EditUserInV2rayApi(user.Name, delete)
					db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeHysteria, delete)
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.HysteriaPtr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {
				userList = append(userList, index)
				userNameList = append(userNameList, user.Name)
				userPasswordList = append(userPasswordList, user.Password)
				inbound.HysteriaPtr[i].Service.DeleteUser(userList, userPasswordList)
				box.EditUserInV2rayApi(user.Name, delete)
				for j := range newUsers {
					for k := range inbound.HysteriaPtr[i].UserNameList {
						if newUsers[j].Name == inbound.HysteriaPtr[i].UserNameList[k] {
							inbound.HysteriaPtr[i].UserNameList = append(
								inbound.HysteriaPtr[i].UserNameList[:k],
								inbound.HysteriaPtr[i].UserNameList[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}
