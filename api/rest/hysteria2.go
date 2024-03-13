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

func AddUserToHysteria2(c *gin.Context) {
	domesticLogicHysteria2(c, false)
}

func DeleteUserToHysteria2(c *gin.Context) {
	domesticLogicHysteria2(c, true)
}

func domesticLogicHysteria2(c *gin.Context, delete bool) {
	var rq option.Hysteria2User
	var rqArr []option.Hysteria2User

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditHysteria2Users(c, []option.Hysteria2User{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditHysteria2Users(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditHysteria2Users(c *gin.Context, newUsers []option.Hysteria2User, delete bool) {
	userList := make([]int, 0, len(newUsers))
	userNameList := make([]string, 0, len(newUsers))
	userPasswordList := make([]string, 0, len(newUsers))
	for index, user := range newUsers {
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		userPasswordList = append(userPasswordList, user.Password)
	}
	for i := range inbound.Hysteria2Ptr {
		if !delete {
			inbound.Hysteria2Ptr[i].Service.AddUser(userList, userPasswordList)
		} else {
			inbound.Hysteria2Ptr[i].Service.DeleteUser(userList, userPasswordList)
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeHysteria2, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
