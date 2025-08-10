package constants

var (
	TGCLOUD_BASE_URL = "https://tgcloud.io/api"
	TIGERTOOL_URL    = "https://tigertool.tigergraph.com"
)

const (
	VERSION_CLI      = "0.1.1"
	GSQL_PATH        = "/gsqlserver/gsql/"
	GSQL_SEPARATOR   = "__GSQL__"
	GSQL_COOKIES     = "__GSQL__COOKIES__"
	COMMAND_ENDPOINT = "command"
	FILE_ENDPOINT    = "file"
	LOGIN_ENDPOINT   = "login"
)

var (
	HomeDir          string
	ConfigDir        string
	ConfigFile       string
	CredsFile        string
	Debug            bool
	AvailableVersion string
)
