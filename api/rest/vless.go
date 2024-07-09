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

func AddUserToVless(c *gin.Context) {
	domesticLogicVless(c, false)
}

func DeleteUserToVless(c *gin.Context) {
	domesticLogicVless(c, true)
}

func domesticLogicVless(c *gin.Context, delete bool) {
	var rq option.VLESSUser
	var rqArr []option.VLESSUser

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditVlessUsers(c, []option.VLESSUser{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditVlessUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditVlessUsers(c *gin.Context, newUsers []option.VLESSUser, delete bool) {
	for _, user := range newUsers {
		box.EditUserInV2rayApi(user.Name, delete)
	}
	for i := range inbound.VLESSPtr {
		dbUser, _ := db.ConvertSingleProtocolModelToDbUser[option.VLESSUser](newUsers[i])
		if !delete {
			if db.GetDb().AddUserInRamUsersIfNotExist(dbUser) {
				inbound.VLESSPtr[i].Service.AddUser(
					common.MapIndexedString(newUsers, func(index any, it option.VLESSUser) string {
						return it.UUID
					}), common.Map(newUsers, func(it option.VLESSUser) string {
						return it.UUID
					}), common.Map(newUsers, func(it option.VLESSUser) string {
						return it.Flow
					}))
			}
		} else {
			if db.GetDb().IsUserExistInRamUsers(dbUser) {
				inbound.VLESSPtr[i].Service.DeleteUser(
					common.MapIndexedString(newUsers, func(index any, it option.VLESSUser) string {
						return it.UUID
					}), common.Map(newUsers, func(it option.VLESSUser) string {
						return it.UUID
					}), common.Map(newUsers, func(it option.VLESSUser) string {
						return it.Flow
					}))
			}
		}
	}
}

func AddVlessInbound(c *gin.Context) {
	var rq option.VLESSInboundOptions
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := box.AddInbound(box.Options{
		Options: option.Options{
			Inbounds: []option.Inbound{
				{
					Type:         C.TypeVLESS,
					VLESSOptions: rq,
				},
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
