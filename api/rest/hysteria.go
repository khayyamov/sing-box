package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"io"
	"net/http"
)

func AddUserToHysteria(c *gin.Context) {
	domesticLogicHysteria(c, false)
}

func DeleteUserToHysteria(c *gin.Context) {
	domesticLogicHysteria(c, true)
}

func domesticLogicHysteria(c *gin.Context, delete bool) {
	var rq option.HysteriaUser
	var rqArr []option.HysteriaUser

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditHysteriaUsers(c, []option.HysteriaUser{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditHysteriaUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
func EditHysteriaUsers(c *gin.Context, newUsers []option.HysteriaUser, delete bool) {
	userList := make([]int, 0, len(newUsers))
	userNameList := make([]string, 0, len(newUsers))
	userPasswordList := make([]string, 0, len(newUsers))
	for index, user := range newUsers {
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		var password string
		if user.AuthString != "" {
			password = user.AuthString
		} else {
			password = string(user.Auth)
		}
		userPasswordList = append(userPasswordList, password)
	}
	for i := range inbound.HysteriaPtr {
		if !delete {
			inbound.HysteriaPtr[i].Service.AddUser(userList, userPasswordList)
		} else {
			inbound.HysteriaPtr[i].Service.DeleteUser(userList, userPasswordList)
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeHysteria, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
