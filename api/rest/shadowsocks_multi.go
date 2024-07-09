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

func AddUserToShadowsocksMulti(c *gin.Context) {
	domesticLogicShadowsocksMulti(c, false)
}

func DeleteUserToShadowsocksMulti(c *gin.Context) {
	domesticLogicShadowsocksMulti(c, true)
}

func domesticLogicShadowsocksMulti(c *gin.Context, delete bool) {
	var rq option.ShadowsocksUser
	var rqArr []option.ShadowsocksUser

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditShadowsocksMultiUsers(c, []option.ShadowsocksUser{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditShadowsocksMultiUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditShadowsocksMultiUsers(c *gin.Context, newUsers []option.ShadowsocksUser, delete bool) {
	for _, user := range newUsers {
		box.EditUserInV2rayApi(user.Name, delete)
	}
	for i := range inbound.ShadowsocksMultiPtr {
		if !delete {
			inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users, newUsers...)
			err := inbound.ShadowsocksMultiPtr[i].Service.AddUsersWithPasswords(
				common.MapIndexedString(newUsers, func(index any, user option.ShadowsocksUser) string {
					return user.Name
				}), common.Map(newUsers, func(user option.ShadowsocksUser) string {
					return user.Password
				}))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			for j := range newUsers {
				for k := range inbound.ShadowsocksMultiPtr[i].Users {
					if newUsers[j].Password == inbound.ShadowsocksMultiPtr[i].Users[k].Password {
						inbound.ShadowsocksMultiPtr[i].Users = append(inbound.ShadowsocksMultiPtr[i].Users[:k],
							inbound.ShadowsocksMultiPtr[i].Users[k+1:]...)
						break
					}
				}
			}
			err := inbound.ShadowsocksMultiPtr[i].Service.DeleteUsersWithPasswords(
				common.MapIndexedString(newUsers, func(index any, user option.ShadowsocksUser) string {
					return user.Name
				}), common.Map(newUsers, func(user option.ShadowsocksUser) string {
					return user.Password
				}))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeShadowsocksMulti, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
