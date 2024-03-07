package db

import (
	"encoding/json"
	"github.com/sagernet/sing-box/api/db/entity"
	"github.com/sagernet/sing-box/api/db/mysql_config"
	"github.com/sagernet/sing-box/option"
)

type ImplementationDb struct {
	Interface DbInterface
}

func (pr *ImplementationDb) GetVlessUsers() ([]option.VLESSUser, error) {
	var user_jsons []entity.DbUser
	query := mysql_config.GetTableVless().Find(&user_jsons)
	var result []option.VLESSUser
	for _, user_json := range user_jsons {
		temp := option.VLESSUser{}
		err := json.Unmarshal([]byte(user_json.UserJson), &temp)
		if err != nil {
			return nil, err
		} else {
			result = append(result, temp)
		}
	}
	return result, query.Error
}

func (pr *ImplementationDb) AddVlessUser(users []option.VLESSUser) error {
	var result []entity.DbUser
	for _, user := range users {
		temp, err := json.Marshal(user)
		if err != nil {
			return err
		} else {
			result = append(result, entity.DbUser{
				UserJson: string(temp),
			})
		}
	}
	query := mysql_config.GetTableVless().Create(&result)
	return query.Error
}
