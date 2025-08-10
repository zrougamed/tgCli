package models

// Config represents the application configuration
type Config struct {
	TGCloud  TGCloudConfig            `mapstructure:"tgcloud"`
	Machines map[string]MachineConfig `mapstructure:"machines"`
	Default  string                   `mapstructure:"default"`
}

type TGCloudConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type MachineConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	GSPort   string `mapstructure:"gsPort"`
	RestPort string `mapstructure:"restPort"`
}

// GSQLCookie represents GSQL session cookies
type GSQLCookie struct {
	ClientCommit                   string `json:"clientCommit"`
	FromGsqlClient                 bool   `json:"fromGsqlClient"`
	FromGraphStudio                bool   `json:"fromGraphStudio"`
	GShellTest                     bool   `json:"gShellTest"`
	FromGsqlServer                 bool   `json:"fromGsqlServer"`
	ApplicationGatewayAffinity     string `json:"ApplicationGatewayAffinity,omitempty"`
	ApplicationGatewayAffinityCORS string `json:"ApplicationGatewayAffinityCORS,omitempty"`
}

// TGCloudResponse represents API responses from TigerGraph Cloud
type TGCloudResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
	Token   string      `json:"token,omitempty"`
}

// Machine represents a TigerGraph Cloud instance
type Machine struct {
	ID        string `json:"ID"`
	Name      string `json:"Name"`
	Tag       string `json:"Tag"`
	State     string `json:"State"`
	CreatedAt string `json:"CreatedAt"`
}
