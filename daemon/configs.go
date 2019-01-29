package daemon

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
)
type Config struct {
	Db *db.Config
}

func getConfigs() *Config {
	var cfg = Config{}
	cfg.Db = db.GetDbConfig()
	return &cfg
}
