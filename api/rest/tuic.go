package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"net/http"
)

func AddUserToTuic(c *gin.Context) {
	var rq option.TUICUser
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUsers := []option.TUICUser{rq}

	var userList []int
	var userNameList []string
	var userUUIDList [][16]byte
	var userPasswordList []string
	for index, user := range newUsers {
		if user.UUID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing uuid for user " + user.UUID})
		}
		userUUID, err := uuid.FromString(user.UUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid for user " + user.UUID})
		}
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		userUUIDList = append(userUUIDList, userUUID)
		userPasswordList = append(userPasswordList, user.Password)
	}
	for i := range inbound.TUICPtr {
		inbound.TUICPtr[i].Service.AddUser(userList, userUUIDList, userPasswordList)
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().AddUserToDb(users, C.TypeTUIC)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
