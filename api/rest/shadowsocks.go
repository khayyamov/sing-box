package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/inbound"
	shadowtls "github.com/sagernet/sing-shadowtls"
	"net/http"
)

func AddUserToShadowsocks(c *gin.Context) {
	var rq shadowtls.User
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUsers := []shadowtls.User{rq}
	inbound.ShadowTlsPtr.Service.AddUser(newUsers)
}
