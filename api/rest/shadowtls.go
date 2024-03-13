package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	shadowtls "github.com/sagernet/sing-shadowtls"
	"io"
	"net/http"
)

func AddUserToShadowtls(c *gin.Context) {
	domesticLogicShadowtls(c, false)
}

func DeleteUserToShadowtls(c *gin.Context) {
	domesticLogicShadowtls(c, true)
}

func domesticLogicShadowtls(c *gin.Context, delete bool) {
	var rq shadowtls.User
	var rqArr []shadowtls.User

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditShadowtlsUsers(c, []shadowtls.User{rq}, delete)
		return
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditShadowtlsUsers(c, rqArr, delete)
		return
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditShadowtlsUsers(c *gin.Context, newUsers []shadowtls.User, delete bool) {
	for i := range inbound.ShadowTlsPtr {
		if !delete {
			inbound.ShadowTlsPtr[i].Service.AddUser(newUsers)
		} else {
			err := inbound.ShadowTlsPtr[i].Service.DeleteUser(newUsers)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().EditDbUser(users, C.TypeShadowTLS, delete)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
