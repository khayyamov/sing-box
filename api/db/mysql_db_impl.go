package db

import (
	"github.com/sagernet/sing-box/api/db/entity"
	"github.com/sagernet/sing-box/api/db/mysql_config"
	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	shadowtls "github.com/sagernet/sing-shadowtls"
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/json"
)

type ImplementationDb struct {
	Interface DbInterface
}

func (pr *ImplementationDb) EditDbUser(users []entity.DbUser, protocolType string, delete bool) error {
	if protocolType == constant.TypeVLESS {
		if !delete {
			query := mysql_config.GetTableVless().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.VLESSUser{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableVless().Where("user_json LIKE ?", "%\""+model.UUID+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeVMess {
		if !delete {
			query := mysql_config.GetTableVmess().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.VMessUser{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableVmess().Where("user_json LIKE ?", "%\""+model.UUID+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeTrojan {
		if !delete {
			query := mysql_config.GetTableTrojan().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.TrojanUser{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableTrojan().Where("user_json LIKE ?", "%\""+model.Password+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeNaive {
		if !delete {
			query := mysql_config.GetTableNaive().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = auth.User{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableNaive().Where("user_json LIKE ?", "%\""+model.Password+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeHysteria {
		if !delete {
			query := mysql_config.GetTableHysteria().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.HysteriaUser{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableHysteria().Where("user_json LIKE ?", "%\""+model.AuthString+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeHysteria2 {
		if !delete {
			query := mysql_config.GetTableHysteria2().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.Hysteria2User{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableHysteria2().Where("user_json LIKE ?", "%\""+model.Password+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeShadowsocksMulti {
		if !delete {
			query := mysql_config.GetTableShadowsocksMulti().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.ShadowsocksUser{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableShadowsocksMulti().Where("user_json LIKE ?", "%\""+model.Password+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeShadowsocksRelay {
		if !delete {
			query := mysql_config.GetTableShadowsocksRelay().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.ShadowsocksDestination{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableShadowsocksRelay().Where("user_json LIKE ?", "%\""+model.Password+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeShadowTLS {
		if !delete {
			query := mysql_config.GetTableShadowtls().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.ShadowTLSUser{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableShadowtls().Where("user_json LIKE ?", "%\""+model.Password+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	} else if protocolType == constant.TypeTUIC {
		if !delete {
			query := mysql_config.GetTableTuic().Create(&users)
			return query.Error
		} else {
			for i := range users {
				var model = option.TUICUser{}
				err := json.Unmarshal([]byte(users[i].UserJson), &model)
				if err != nil {
					return err
				}
				query := mysql_config.GetTableTuic().Where("user_json LIKE ?", "%\""+model.Password+"\"%").Delete(&users)
				if query.Error != nil {
					return query.Error
				}
			}
			return nil
		}
	}
	return nil
}

func getUser[T any](protocolType string) ([]T, error) {
	var userJsons []entity.DbUser

	if protocolType == constant.TypeVLESS {
		query := mysql_config.GetTableVless().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeVMess {
		query := mysql_config.GetTableVmess().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeTrojan {
		query := mysql_config.GetTableTrojan().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeTUIC {
		query := mysql_config.GetTableTuic().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeHysteria2 {
		query := mysql_config.GetTableNaive().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeHysteria {
		query := mysql_config.GetTableHysteria().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeShadowsocksMulti {
		query := mysql_config.GetTableShadowsocksMulti().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeShadowsocksRelay {
		query := mysql_config.GetTableShadowsocksRelay().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeShadowTLS {
		query := mysql_config.GetTableShadowtls().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeTUIC {
		query := mysql_config.GetTableTuic().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	} else if protocolType == constant.TypeNaive {
		query := mysql_config.GetTableNaive().Find(&userJsons)
		result, _ := ConvertDbUserToProtocolModel[T](userJsons)
		return result, query.Error
	}
	query := mysql_config.GetTableVless().Find(&userJsons)
	result, _ := ConvertDbUserToProtocolModel[T](userJsons)
	return result, query.Error
}

func (pr *ImplementationDb) GetVlessUsers() ([]option.VLESSUser, error) {
	return getUser[option.VLESSUser](constant.TypeVLESS)
}

func (pr *ImplementationDb) GetVmessUsers() ([]option.VMessUser, error) {
	return getUser[option.VMessUser](constant.TypeVMess)
}

func (pr *ImplementationDb) GetTrojanUsers() ([]option.TrojanUser, error) {
	return getUser[option.TrojanUser](constant.TypeTrojan)
}

func (pr *ImplementationDb) GetHysteria2Users() ([]option.Hysteria2User, error) {
	return getUser[option.Hysteria2User](constant.TypeHysteria2)
}

func (pr *ImplementationDb) GetHysteriaUsers() ([]option.HysteriaUser, error) {
	return getUser[option.HysteriaUser](constant.TypeHysteria)
}

func (pr *ImplementationDb) GetNaiveUsers() ([]auth.User, error) {
	return getUser[auth.User](constant.TypeNaive)
}

func (pr *ImplementationDb) GetTuicUsers() ([]option.TUICUser, error) {
	return getUser[option.TUICUser](constant.TypeTUIC)
}

func (pr *ImplementationDb) GetShadowtlsUsers() ([]shadowtls.User, error) {
	return getUser[shadowtls.User](constant.TypeShadowTLS)
}

func (pr *ImplementationDb) GetShadowsocksMultiUsers() ([]option.ShadowsocksUser, error) {
	return getUser[option.ShadowsocksUser](constant.TypeShadowsocksMulti)
}

func (pr *ImplementationDb) GetShadowsocksRelayUsers() ([]option.ShadowsocksDestination, error) {
	return getUser[option.ShadowsocksDestination](constant.TypeShadowsocksRelay)
}