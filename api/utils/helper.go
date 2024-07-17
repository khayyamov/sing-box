package utils

import (
	"encoding/json"
	"errors"
	"github.com/sagernet/sing-box/api/constant"
	log2 "github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/auth"
)

func UUIDFromDBUserJson(user string) (string, error) {

	//vless, vmess, tuic
	uuidUser := option.VMessUser{}
	err := json.Unmarshal([]byte(user), &uuidUser)
	if err == nil {
		if len(uuidUser.UUID) > 0 {
			return uuidUser.UUID, nil
		}
	}

	//naive
	usernameUser := auth.User{}
	err = json.Unmarshal([]byte(user), &usernameUser)
	if err == nil {
		if len(usernameUser.Username) > 0 {
			return usernameUser.Username, nil
		}
	}

	//hysteria1,2 and name base protocol
	nameUser := option.HysteriaUser{}
	err = json.Unmarshal([]byte(user), &nameUser)
	if err == nil {
		if len(nameUser.Name) > 0 {
			return nameUser.Name, nil
		}
	}
	return "", errors.New("")
}

func ApiLogInfo(log any) {
	if constant.ApiLog {
		log2.Info(log)
	}
}

func ApiLogError(log any) {
	if constant.ApiLog {
		log2.Error(log)
	}
}
