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

func EditHysteria2Users(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	utils.CurrentInboundName = "Hysteria2"
	if len(inbound.Hysteria2Ptr) == 0 {
		utils.ApiLogError("No Active " + utils.CurrentInboundName + " outbound found to add users to it")
		return
	}
	for index, user := range newUsers {
		convertedUser := option.Hysteria2User{
			Name:     user.Name,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.Hysteria2User](convertedUser)
		for i := range inbound.Hysteria2Ptr {
			userList := make([]int, 0, len(newUsers))
			userPasswordList := make([]string, 0, len(newUsers))
			if !delete {
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
					userList = append(userList, len(inbound.Hysteria2Ptr[i].UserNameList)+index)
					userPasswordList = append(userPasswordList, user.Password)
					inbound.Hysteria2Ptr[i].Service.AddUser(userList, userPasswordList)
					inbound.Hysteria2Ptr[i].Users = append(inbound.Hysteria2Ptr[i].Users, convertedUser)
					inbound.Hysteria2Ptr[i].UserNameList = append(inbound.Hysteria2Ptr[i].UserNameList, convertedUser.Name)
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.Hysteria2Ptr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {

				if len(convertedUser.Name) == 0 {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.Hysteria2Ptr[i].Tag() + "] User failed to delete Name invalid")
					continue
				}
				userList = append(userList, len(inbound.Hysteria2Ptr[i].UserNameList)+index)
				userPasswordList = append(userPasswordList, user.Password)
				inbound.Hysteria2Ptr[i].Service.DeleteUser(userList, userPasswordList)
				for j := range newUsers {
					for k := range inbound.Hysteria2Ptr[i].Users {
						if newUsers[j].Name == inbound.Hysteria2Ptr[i].Users[k].Name {
							inbound.Hysteria2Ptr[i].Users = append(
								inbound.Hysteria2Ptr[i].Users[:k],
								inbound.Hysteria2Ptr[i].Users[k+1:]...)
							inbound.Hysteria2Ptr[i].UserNameList = append(
								inbound.Hysteria2Ptr[i].UserNameList[:k],
								inbound.Hysteria2Ptr[i].UserNameList[k+1:]...)
							break
						}
					}
				}
			}
			box.EditUserInV2rayApi(convertedUser.Name, delete)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeHysteria2, delete)
		}
	}
}
