package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/api/db"
	"github.com/sagernet/sing-box/api/db/entity"
	"github.com/sagernet/sing-box/api/rest/rq"
	"github.com/sagernet/sing-box/api/utils"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
)

func EditTuicUsers(c *gin.Context, newUsers []rq.GlobalModel, delete bool) {
	if len(inbound.TUICPtr) == 0 {
		utils.ApiLogInfo("No Active TUICPtr outbound found to add users to it")
		return
	}
	var userList []int
	var userNameList []string
	var userUUIDList [][16]byte
	var userPasswordList []string
	for index, user := range newUsers {
		convertedUser := option.TUICUser{
			Name:     user.Name,
			UUID:     user.UUID,
			Password: user.Password,
		}
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.TUICUser](convertedUser)
		if db.GetDb().IsUserExistInRamUsers(dbUser) && !delete {
			utils.ApiLogError("User already exist: " + dbUser.UserJson)
			continue
		}
		userUUID, err := uuid.FromString(user.UUID)
		if err != nil {
			continue
		}
		userList = append(userList, index)
		userNameList = append(userNameList, convertedUser.Name)
		userUUIDList = append(userUUIDList, userUUID)
		userPasswordList = append(userPasswordList, user.Password)
		box.EditUserInV2rayApi(user.UUID, delete)
		db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeTUIC, delete)
		for i := range inbound.TUICPtr {
			if !delete {
				if len(user.ReplacementField) > 0 {
					for _, model := range user.ReplacementField {
						if inbound.TUICPtr[i].Tag() == model.Tag {
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
				inbound.TUICPtr[i].Service.AddUser(userList, userUUIDList, userPasswordList)
				inbound.TUICPtr[i].Users = append(inbound.TUICPtr[i].Users, convertedUser)
			} else {

				inbound.TUICPtr[i].Service.DeleteUser(userList, userUUIDList)
				for j := range newUsers {
					for k := range inbound.TUICPtr[i].Users {
						if newUsers[j].UUID == inbound.TUICPtr[i].Users[k].UUID {
							inbound.TUICPtr[i].Users = append(
								inbound.TUICPtr[i].Users[:k],
								inbound.TUICPtr[i].Users[k+1:]...)
							break
						}
					}
				}
			}
		}
	}
}

//func EditTuicUsers(c *gin.Context, newUsers []option.TUICUser, delete bool) {
//	for _, user := range newUsers {
//		box.EditUserInV2rayApi(user.Name, delete)
//	}
//	var userList []int
//	var userNameList []string
//	var userUUIDList [][16]byte
//	var userPasswordList []string
//	for index, user := range newUsers {
//		if user.UUID == "" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "missing uuid for user " + user.UUID})
//		}
//		userUUID, err := uuid.FromString(user.UUID)
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid for user " + user.UUID})
//		}
//		userList = append(userList, index)
//		userNameList = append(userNameList, user.Name)
//		userUUIDList = append(userUUIDList, userUUID)
//		userPasswordList = append(userPasswordList, user.Password)
//	}
//	for i := range inbound.TUICPtr {
//		if !delete {
//			inbound.TUICPtr[i].Service.AddUser(userList, userUUIDList, userPasswordList)
//		} else {
//			inbound.TUICPtr[i].Service.DeleteUser(userList, userUUIDList)
//		}
//	}
//	users, err := db.ConvertProtocolModelToDbUser(newUsers)
//	err = db.GetDb().EditDbUser(users, C.TypeTUIC, delete)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//	}
//
//}
