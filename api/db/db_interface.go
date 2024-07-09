package db

import (
	"github.com/sagernet/sing-box/api/db/entity"
	"github.com/sagernet/sing-box/option"
	shadowtls "github.com/sagernet/sing-shadowtls"
	"github.com/sagernet/sing/common/auth"
)

type DbInterface interface {
	GetVlessUsers() ([]option.VLESSUser, error)

	GetVmessUsers() ([]option.VMessUser, error)

	GetTrojanUsers() ([]option.TrojanUser, error)

	GetHysteria2Users() ([]option.Hysteria2User, error)

	GetHysteriaUsers() ([]option.HysteriaUser, error)

	GetNaiveUsers() ([]auth.User, error)

	GetTuicUsers() ([]option.TUICUser, error)

	GetShadowtlsUsers() ([]shadowtls.User, error)

	GetShadowsocksMultiUsers() ([]option.ShadowsocksUser, error)

	GetShadowsocksRelayUsers() ([]option.ShadowsocksDestination, error)

	EditDbUser(v []entity.DbUser, protocolType string, delete bool) error

	EditInRamUsers(v []entity.DbUser, deleteUser bool)

	IsUserExistInRamUsers(user entity.DbUser) bool

	AddUserInRamUsersIfNotExist(user entity.DbUser) bool
}
