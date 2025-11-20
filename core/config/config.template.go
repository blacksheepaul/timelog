package config

type Config struct {
	Database struct {
		Host string `yaml:"host"`
	} `yaml:"database"`
	Server struct {
		Addr         string   `yaml:"addr" env-default:""`
		Port         int      `yaml:"port" env-required:"true"`
		AllowOrigins []string `yaml:"allow_origins"`
	} `yaml:"server"`
	Log struct {
		Level    string `yaml:"level" env-default:"info"`
		Path     string `yaml:"path" env-default:"app.log"`
		Rotation struct {
			MaxSize    int `yaml:"max_size"`
			MaxBackups int `yaml:"max_backups"`
			MaxAge     int `yaml:"max_age"`
		} `yaml:"rotation"`
		ORMLogLevel int `yaml:"orm_log_level"`
		MCP         struct {
			Enabled bool   `yaml:"enabled" env-default:"false"`
			Level   string `yaml:"level" env-default:"debug"`
			Path    string `yaml:"path" env-default:"logs/mcp.log"`
		} `yaml:"mcp"`
	} `yaml:"log"`
	Test struct {
		Flush bool `yaml:"flush" env-default:"false"`
	} `yaml:"test"`
}
