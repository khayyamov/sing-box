package db

import (
	"github.com/sagernet/sing-box/api/db/entity"
	"github.com/sagernet/sing-box/api/db/mysql_config"
	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	shadowtls "github.com/sagernet/sing-shadowtls"
	"github.com/sagernet/sing/common/auth"
)

type ImplementationDb struct {
	Interface DbInterface
}

func (pr *ImplementationDb) AddUserToDb(users []entity.DbUser, protocolType string) error {
	if protocolType == constant.TypeVLESS {
		query := mysql_config.GetTableVless().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeVMess {
		query := mysql_config.GetTableVmess().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeTrojan {
		query := mysql_config.GetTableVless().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeVMess {
		query := mysql_config.GetTableTrojan().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeNaive {
		query := mysql_config.GetTableNaive().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeHysteria {
		query := mysql_config.GetTableHysteria().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeShadowsocksMulti {
		query := mysql_config.GetTableShadowsocksMulti().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeShadowsocksRelay {
		query := mysql_config.GetTableShadowsocksRelay().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeShadowTLS {
		query := mysql_config.GetTableShadowtls().Create(&users)
		return query.Error
	} else if protocolType == constant.TypeTUIC {
		query := mysql_config.GetTableTuic().Create(&users)
		return query.Error
	}
	return nil
}

func getUser[T any]() ([]T, error) {
	var userJsons []entity.DbUser
	query := mysql_config.GetTableVless().Find(&userJsons)
	result, _ := ConvertDbUserToProtocolModel[T](userJsons)
	return result, query.Error
}

func (pr *ImplementationDb) GetVlessUsers() ([]option.VLESSUser, error) {
	return getUser[option.VLESSUser]()
}

func (pr *ImplementationDb) GetVmessUsers() ([]option.VMessUser, error) {
	return getUser[option.VMessUser]()
}

func (pr *ImplementationDb) GetTrojanUsers() ([]option.TrojanUser, error) {
	return getUser[option.TrojanUser]()
}

func (pr *ImplementationDb) GetHysteria2Users() ([]option.Hysteria2User, error) {
	return getUser[option.Hysteria2User]()
}

func (pr *ImplementationDb) GetHysteriaUsers() ([]option.HysteriaUser, error) {
	return getUser[option.HysteriaUser]()
}

func (pr *ImplementationDb) GetNaiveUsers() ([]auth.User, error) {
	return getUser[auth.User]()
}

func (pr *ImplementationDb) GetTuicUsers() ([]option.TUICUser, error) {
	return getUser[option.TUICUser]()
}

func (pr *ImplementationDb) GetShadowtlsUsers() ([]shadowtls.User, error) {
	return getUser[shadowtls.User]()
}

func (pr *ImplementationDb) GetShadowsocksMultiUsers() ([]option.ShadowsocksUser, error) {
	return getUser[option.ShadowsocksUser]()
}

func (pr *ImplementationDb) GetShadowsocksRelayUsers() ([]option.ShadowsocksDestination, error) {
	return getUser[option.ShadowsocksDestination]()
}
