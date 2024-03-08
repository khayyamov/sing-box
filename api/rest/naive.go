package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing/common/auth"
	"net/http"
)

func AddUserTonNaive(c *gin.Context) {
	var rq auth.User
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []auth.User{rq}
	for i := range inbound.NaivePtr {
		inbound.NaivePtr[i].Authenticator.AddUserToAuthenticator(newUsers)
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().AddUserToDb(users, C.TypeNaive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
