package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/api/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/inbound"
	shadowtls "github.com/sagernet/sing-shadowtls"
	"net/http"
)

func AddUserToShadowtls(c *gin.Context) {
	var rq shadowtls.User
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []shadowtls.User{rq}
	for i := range inbound.ShadowTlsPtr {
		inbound.ShadowTlsPtr[i].Service.AddUser(newUsers)
	}
	users, err := db.ConvertProtocolModelToDbUser(newUsers)
	err = db.GetDb().AddUserToDb(users, C.TypeShadowTLS)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
