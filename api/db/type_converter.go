package db

import (
	"encoding/json"
	"github.com/sagernet/sing-box/api/db/entity"
)

func ConvertProtocolModelToDbUser[T any](arr []T) ([]entity.DbUser, error) {
	retArr := make([]entity.DbUser, 0, len(arr))
	for _, user := range arr {
		temp, err := json.Marshal(user)
		if err != nil {
			return nil, err
		} else {
			retArr = append(retArr, entity.DbUser{
				UserJson: string(temp),
			})
		}
	}
	return retArr, nil
}

func ConvertDbUserToProtocolModel[T any](arr []entity.DbUser) ([]T, error) {
	retArr := make([]T, 0, len(arr))
	for _, user := range arr {
		var userT T
		err := json.Unmarshal([]byte(user.UserJson), &userT)
		if err != nil {
			return nil, err
		}
		retArr = append(retArr, userT)
	}
	return retArr, nil
}