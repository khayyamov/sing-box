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

func AddUserToShadowsocksRelay(c *gin.Context) {
	domesticLogicShadowsocksRelay(c, false)
}

func DeleteUserToShadowsocksRelay(c *gin.Context) {
	domesticLogicShadowsocksRelay(c, true)
}

func domesticLogicShadowsocksRelay(c *gin.Context, delete bool) {
	var rq option.ShadowsocksDestination
	var rqArr []option.ShadowsocksDestination

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditShadowsocksRelayUsers(c, []option.ShadowsocksDestination{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditShadowsocksRelayUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditShadowsocksRelayUsers(c *gin.Context, newUsers []option.ShadowsocksDestination, delete bool) {
	for _, user := range newUsers {
		box.EditUserInV2rayApi(user.Name, delete)
	}
	for i := range inbound.ShadowsocksRelayPtr {
		if !delete {
			inbound.ShadowsocksRelayPtr[i].Destinations = append(
				inbound.ShadowsocksRelayPtr[i].Destinations, newUsers...)
			err := inbound.ShadowsocksRelayPtr[i].Service.AddUsersWithPasswords(
				common.MapIndexedString(newUsers, func(index any, user option.ShadowsocksDestination) string {
					return user.Name
				}), common.Map(newUsers, func(user option.ShadowsocksDestination) string {
					return user.Password
				}), common.Map(newUsers, option.ShadowsocksDestination.Build))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		} else {
			for j := range newUsers {
				for k := range inbound.ShadowsocksRelayPtr[i].Destinations {
					if newUsers[j].Password == inbound.ShadowsocksRelayPtr[i].Destinations[k].Password {
						inbound.ShadowsocksRelayPtr[i].Destinations = append(
							inbound.ShadowsocksRelayPtr[i].Destinations[:k],
							inbound.ShadowsocksRelayPtr[i].Destinations[k+1:]...)
						break
					}
				}
			}
			err := inbound.ShadowsocksRelayPtr[i].Service.DeleteUsersWithPasswords(
				common.MapIndexedString(newUsers, func(index any, user option.ShadowsocksDestination) string {
					return user.Name
				}), common.Map(newUsers, func(user option.ShadowsocksDestination) string {
					return user.Password
				}))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeShadowsocksRelay, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
