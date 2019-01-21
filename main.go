package main

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/daemon"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
)

func processFlags() *daemon.Config {
	var cfg = daemon.Config{}
	cfg.Db = &db.Config{
		Driver:    "mssql",
		User: 	  "********",
		Password: "********",
		Host:	  "********",
		Database: "********",
		Port:	   1433,
	}
	return &cfg
}

func main() {
	cfg := processFlags()
	daemon.Run(cfg)
}