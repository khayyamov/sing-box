package db

import "github.com/sagernet/sing-box/option"

type DbInterface interface {
	GetVlessUsers() ([]option.VLESSUser, error)
	AddVlessUser(users []option.VLESSUser) error
}
