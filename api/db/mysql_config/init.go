package mysql_config

import (
	"fmt"
	"github.com/sagernet/sing-box/api/constant"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var dbLock = &sync.Mutex{}
var DB *gorm.DB

func MySqlInstance() *gorm.DB {
	var err error
	if DB == nil {

		dbLock.Lock()
		defer dbLock.Unlock()

		if DB == nil {

			DB, err = gorm.Open(mysql.Open(constant.DbUsername + ":" + constant.DbPassword + "@" + constant.DbConnection + "(" + constant.DbHost + ":" + constant.DbPort + ")/" + constant.DbName + "?charset=" + constant.DbCharset + "&parseTime=True"))
			if err != nil {
				panic(fmt.Errorf("connect db fail: %w", err))
			}

		}

	}
	return DB
}
