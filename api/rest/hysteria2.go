package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sagernet/sing-box/inbound"
	"github.com/sagernet/sing-box/option"
	"net/http"
)

func AddUserToHysteria2(c *gin.Context) {
	var rq option.Hysteria2User
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUsers := []option.Hysteria2User{rq}

	userList := make([]int, 0, len(newUsers))
	userNameList := make([]string, 0, len(newUsers))
	userPasswordList := make([]string, 0, len(newUsers))
	for index, user := range newUsers {
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		userPasswordList = append(userPasswordList, user.Password)
	}
	inbound.Hysteria2Ptr.Service.AddUser(userList, userPasswordList)
}
