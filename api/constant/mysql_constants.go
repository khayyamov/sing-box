package constant

const DbConnection = "tcp"
const DbCharset = "utf8"

var (
	DbName     = "users_db"
	DbHost     string // init from argument mysql_host
	DbPort     string // init from argument mysql_port
	DbUsername string // init from argument mysql_username
	DbPassword string // init from argument mysql_password
	ApiHost    string // init from argument api_host
	ApiPort    string // init from argument api_port
)
