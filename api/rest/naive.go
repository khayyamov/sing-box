package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing/common/auth"
	"io"
	"net/http"
)

func AddUserTonNaive(c *gin.Context) {
	domesticLogicNaive(c, false)
}
func DeleteUserTonNaive(c *gin.Context) {
	domesticLogicNaive(c, true)
}

func domesticLogicNaive(c *gin.Context, delete bool) {
	var rq auth.User
	var rqArr []auth.User

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditNaiveUsers(c, []auth.User{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditNaiveUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditNaiveUsers(c *gin.Context, newUsers []auth.User, delete bool) {
	for i := range inbound.NaivePtr {
		if !delete {
			inbound.NaivePtr[i].Authenticator.AddUserToAuthenticator(newUsers)
		} else {
			inbound.NaivePtr[i].Authenticator.DeleteUserToAuthenticator(newUsers)
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeNaive, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
