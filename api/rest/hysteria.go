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

func EditHysteriaUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.HysteriaPtr) == 0 {
		log.Info("No Active HysteriaPtr outbound found to add users to it")
		return
	}
	userList := make([]int, 0, len(newUsers))
	userNameList := make([]string, 0, len(newUsers))
	userPasswordList := make([]string, 0, len(newUsers))
	for index, user := range newUsers {
		convertedUser := option.HysteriaUser{
			Name:       user.UUID,
			AuthString: user.AuthString,
			Auth:       []byte(user.Auth),
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.HysteriaUser](convertedUser)
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
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeHysteria, delete)
		for i := range inbound.HysteriaPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.HysteriaPtr[i].Tag() == model.Tag {
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
					continue
				}
				inbound.HysteriaPtr[i].Service.AddUser(userList, userPasswordList)
				inbound.HysteriaPtr[i].UserNameList = append(inbound.HysteriaPtr[i].UserNameList, convertedUser.Name)
			} else {

				inbound.HysteriaPtr[i].Service.DeleteUser(userList, userPasswordList)
				for j := range newUsers {
					for k := range inbound.HysteriaPtr[i].UserNameList {
						if newUsers[j].UUID == inbound.HysteriaPtr[i].UserNameList[k] {
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
