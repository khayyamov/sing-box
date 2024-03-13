package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"io"
	"net/http"
)

func AddUserToTuic(c *gin.Context) {
	domesticLogicTuic(c, true)
}

func DeleteUserToTuic(c *gin.Context) {
	domesticLogicTuic(c, true)
}

func domesticLogicTuic(c *gin.Context, delete bool) {
	var rq option.TUICUser
	var rqArr []option.TUICUser

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditTuicUsers(c, []option.TUICUser{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditTuicUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditTuicUsers(c *gin.Context, newUsers []option.TUICUser, delete bool) {
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
		if !delete {
			inbound.TUICPtr[i].Service.AddUser(userList, userUUIDList, userPasswordList)
		} else {
			inbound.TUICPtr[i].Service.DeleteUser(userList, userUUIDList)
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeTUIC, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
