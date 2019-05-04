package yuwiki

import (
	"github.com/BurntSushi/toml"
	"log"
)

type AppConfig struct {
	Debug         bool   `toml:"debug"`
	LogFile       string `toml:"log_file"`
	SessionCookie string `toml:"session_cookie"`
	Secret        string `toml:"secret"`
	SessionAuth   string `toml:"session_auth"`
	Http          struct {
		Port            string `toml:"port"`
		Favicon         string `toml:"favicon"`
		StaticPath      string `toml:"static_path"`
		HtmlPathPattern string `toml:"html_path_pattern"`
	} `toml:"http"`
	DataSource struct {
		Dialect   string `toml:"dialect"`
		Url       string `toml:"url"`
		DdlUpdate bool   `toml:"ddl_update"`
	} `toml:"dataSource"`
	Cron struct {
		Backup string `toml:"backup"`
	} `toml:"cron"`
}

func InitConfig() *AppConfig {
	var config *AppConfig
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return config
}
