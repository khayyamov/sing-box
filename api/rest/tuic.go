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

func EditTuicUsers(c *gin.Context, newUsers []rq.GlobalModel, deletee bool) {
	utils.CurrentInboundName = "Tuic"
	for _, user := range newUsers {
		for i := range inbound.TUICPtr {
			convertedUser := option.TUICUser{
				Name:     user.Name,
				UUID:     user.UUID,
				Password: user.Password,
			}
			dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.TUICUser](convertedUser)
			var userList []string
			var userNameList []string
			var userUUIDList [][16]byte
			var userPasswordList []string
			if len(user.ReplacementField) > 0 {
				for _, model := range user.ReplacementField {
					if inbound.TUICPtr[i].Tag() == model.Tag {
						if len(model.Name) > 0 {
							convertedUser.Name = model.Name
						}
						if len(model.UUID) > 0 {
							convertedUser.UUID = model.UUID
						}
						if len(model.Password) > 0 {
							convertedUser.Password = model.Password
						}
						break
					}
				}
			}
			if !deletee {
				userUUID, err := uuid.FromString(convertedUser.UUID)
				if len(convertedUser.UUID) == 0 ||
					len(convertedUser.Password) == 0 ||
					len(convertedUser.Name) == 0 ||
					err != nil {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.TUICPtr[i].Tag() + "] User failed to add name or password invalid: " + dbUser.UserJson)
					continue
				}

				founded := false
				for _, inboundUsers := range inbound.TUICPtr[i].Users {
					if inboundUsers.UUID == convertedUser.UUID {
						founded = true
						break
					}
				}
				if !founded {
					dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.TUICUser](convertedUser)
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.TUICPtr[i].Tag() + "] User Added: " + dbUser.UserJson)
					userList = append(userList, convertedUser.UUID)
					userNameList = append(userNameList, convertedUser.Name)
					userUUIDList = append(userUUIDList, userUUID)
					userPasswordList = append(userPasswordList, convertedUser.Password)
					inbound.TUICPtr[i].Service.AddUser(userList, userUUIDList, userPasswordList)
					inbound.TUICPtr[i].Users[convertedUser.UUID] = convertedUser
				} else {
					utils.ApiLogInfo(utils.CurrentInboundName + "[" + inbound.TUICPtr[i].Tag() + "] User already exist: " + dbUser.UserJson)
				}
			} else {
				userUUID, err := uuid.FromString(convertedUser.UUID)
				if err != nil {
					utils.ApiLogError(utils.CurrentInboundName + "[" + inbound.TUICPtr[i].Tag() + "] User failed to deletee uuid invalid: " + dbUser.UserJson)
					continue
				}
				userList = append(userList, convertedUser.UUID)
				userNameList = append(userNameList, convertedUser.Name)
				userUUIDList = append(userUUIDList, userUUID)
				userPasswordList = append(userPasswordList, convertedUser.Password)
				inbound.TUICPtr[i].Service.DeleteUser(userList, userUUIDList)
				delete(inbound.TUICPtr[i].Users, convertedUser.UUID)
			}

			box.EditUserInV2rayApi(convertedUser.Name, deletee)
			db.GetDb().EditDbUser([]entity.DbUser{dbUser}, C.TypeTUIC, deletee)
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
