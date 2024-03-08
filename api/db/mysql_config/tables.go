package mysql_config

import "gorm.io/gorm"

func GetTableVless() *gorm.DB            { return MySqlInstance().Table("vless") }
func GetTableVmess() *gorm.DB            { return MySqlInstance().Table("vmess") }
func GetTableTuic() *gorm.DB             { return MySqlInstance().Table("tuic") }
func GetTableTrojan() *gorm.DB           { return MySqlInstance().Table("trojan") }
func GetTableShadowtls() *gorm.DB        { return MySqlInstance().Table("shadowtls") }
func GetTableShadowsocksMulti() *gorm.DB { return MySqlInstance().Table("shadowsocks_multi") }
func GetTableShadowsocksRelay() *gorm.DB { return MySqlInstance().Table("shadowsocks_relay") }
func GetTableNaive() *gorm.DB            { return MySqlInstance().Table("naive") }
func GetTableHysteria2() *gorm.DB        { return MySqlInstance().Table("hysteria2") }
func GetTableHysteria() *gorm.DB         { return MySqlInstance().Table("hysteria") }
