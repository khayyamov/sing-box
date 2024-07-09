package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"io"
	"net/http"
)

func AddUserToVmess(c *gin.Context) {
	domesticLogicVmess(c, false)
}

func DeleteUserToVmess(c *gin.Context) {
	domesticLogicVmess(c, true)
}

func domesticLogicVmess(c *gin.Context, delete bool) {
	var rq option.VMessUser
	var rqArr []option.VMessUser

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditVmessUsers(c, []option.VMessUser{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditVmessUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditVmessUsers(c *gin.Context, newUsers []option.VMessUser, delete bool) {
	for _, user := range newUsers {
		box.EditUserInV2rayApi(user.Name, delete)
	}
	for i := range inbound.VMessPtr {
		if !delete {
			err := inbound.VMessPtr[i].Service.AddUser(
				common.MapIndexedString(newUsers, func(index any, it option.VMessUser) string {
					return it.UUID
				}), common.Map(newUsers, func(it option.VMessUser) string {
					return it.UUID
				}), common.Map(newUsers, func(it option.VMessUser) int {
					return it.AlterId
				}))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			err := inbound.VMessPtr[i].Service.DeleteUser(
				common.MapIndexedString(newUsers, func(index any, it option.VMessUser) string {
					return it.UUID
				}), common.Map(newUsers, func(it option.VMessUser) string {
					return it.UUID
				}), common.Map(newUsers, func(it option.VMessUser) int {
					return it.AlterId
				}))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeVMess, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
