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
	"github.com/sagernet/sing-box/option"
)

func EditHysteria2Users(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.Hysteria2Ptr) == 0 {
		log.Info("No Active Hysteria2 outbound found to add users to it")
		return
	}
	userList := make([]int, 0, len(newUsers))
	userNameList := make([]string, 0, len(newUsers))
	userPasswordList := make([]string, 0, len(newUsers))
	for index, user := range newUsers {
		convertedUser := option.Hysteria2User{
			Name:     user.UUID,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.Hysteria2User](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			log.Error("User already exist: " + dbUser.UserJson)
			continue
		}
		_, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		userList = append(userList, index)
		userNameList = append(userNameList, user.UUID)
		userPasswordList = append(userPasswordList, user.Password)
		box.EditUserInV2rayApi(user.UUID, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeHysteria2, delete)
		for i := range inbound.Hysteria2Ptr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.Hysteria2Ptr[i].Tag() == model.Tag {
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
				inbound.Hysteria2Ptr[i].Service.AddUser(userList, userPasswordList)
				inbound.Hysteria2Ptr[i].Users = append(inbound.Hysteria2Ptr[i].Users, convertedUser)
			} else {

				inbound.Hysteria2Ptr[i].Service.DeleteUser(userList, userPasswordList)
				for j := range newUsers {
					for k := range inbound.Hysteria2Ptr[i].Users {
						if newUsers[j].UUID == inbound.Hysteria2Ptr[i].Users[k].Name {
							inbound.Hysteria2Ptr[i].Users = append(
								inbound.Hysteria2Ptr[i].Users[:k],
								inbound.Hysteria2Ptr[i].Users[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}
