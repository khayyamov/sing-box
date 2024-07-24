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

func EditHysteria2Users(c *gin.Context, newUsers []rq.GlobalModel, deletee bool) {
	utils.CurrentInboundName = "Hysteria2"
	for _, user := range newUsers {
		convertedUser := option.Hysteria2User{
			Name:     user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.Hysteria2User](convertedUser)
		for i := range inbound.Hysteria2Ptr {
			userList := make([]string, 0, len(newUsers))
			userPasswordList := make([]string, 0, len(newUsers))
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.Hysteria2Ptr[i].Tag() == model.Tag {
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
			if !deletee {
				if len(convertedUser.Password) == 0 || len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.Hysteria2Ptr[i].Tag() + "] User failed to add password invalid")
					continue
				}
				founded := false
				for _, inboundUsers := range inbound.Hysteria2Ptr[i].Users {
					if inboundUsers.Name == convertedUser.Name {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.Hysteria2User](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.Hysteria2Ptr[i].Tag() + "] User Added: " + dbUser.UserJson)
					userList = append(userList, convertedUser.Name)
					userPasswordList = append(userPasswordList, convertedUser.Password)
					inbound.Hysteria2Ptr[i].Service.AddUser(userList, userPasswordList)
					inbound.Hysteria2Ptr[i].Users[convertedUser.Name] = convertedUser
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.Hysteria2Ptr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {

				if len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.Hysteria2Ptr[i].Tag() + "] User failed to deletee Name invalid")
					continue
				}
				userList = append(userList, convertedUser.Name)
				userPasswordList = append(userPasswordList, convertedUser.Password)
				inbound.Hysteria2Ptr[i].Service.DeleteUser(userList, userPasswordList)
				delete(inbound.Hysteria2Ptr[i].Users, convertedUser.Name)
			}
			box.EditUserInV2rayApi(convertedUser.Name, deletee)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeHysteria2, deletee)
		}
	}
}
