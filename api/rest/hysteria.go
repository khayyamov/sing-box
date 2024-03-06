package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"net/http"
)

func AddUserToHysteria(c *gin.Context) {
	var rq option.HysteriaUser
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUsers := []option.HysteriaUser{rq}

	userList := make([]int, 0, len(newUsers))
	userNameList := make([]string, 0, len(newUsers))
	userPasswordList := make([]string, 0, len(newUsers))
	for index, user := range newUsers {
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		var password string
		if user.AuthString != "" {
			password = user.AuthString
		} else {
			password = string(user.Auth)
		}
		userPasswordList = append(userPasswordList, password)
	}
	inbound.HysteriaPtr.Service.AddUser(userList, userPasswordList)
}
