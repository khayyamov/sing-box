package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"io"
	"net/http"
)

func AddUserToTrojan(c *gin.Context) {
	domesticLoginTrojan(c, false)
}

func DeleteUserToTrojan(c *gin.Context) {
	domesticLoginTrojan(c, true)
}

func domesticLoginTrojan(c *gin.Context, delete bool) {
	var rq option.TrojanUser
	var rqArr []option.TrojanUser

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditTrojanUsers(c, []option.TrojanUser{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditTrojanUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditTrojanUsers(c *gin.Context, newUsers []option.TrojanUser, delete bool) {
	for _, user := range newUsers {
		if !delete {
			AddUserToV2rayApi(user.Name)
		}
	}
	for i := range inbound.TrojanPtr {
		if !delete {
			err := inbound.TrojanPtr[i].Service.AddUser(
				common.MapIndexed(newUsers, func(index int, it option.TrojanUser) int {
					return index
				}), common.Map(newUsers, func(it option.TrojanUser) string {
					return it.Password
				}))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			err := inbound.TrojanPtr[i].Service.DeleteUser(
				common.MapIndexed(newUsers, func(index int, it option.TrojanUser) int {
					return index
				}), common.Map(newUsers, func(it option.TrojanUser) string {
					return it.Password
				}))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeTrojan, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
