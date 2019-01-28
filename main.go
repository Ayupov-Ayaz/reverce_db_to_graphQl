package main

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/daemon"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"log"
	"os"
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

func getTablesNameForScanning() *[]string {
	return  &[]string{
		// set database tables names
	}
}

func main() {
	cfg := processFlags()
	tables := getTablesNameForScanning()
	if len(*tables) < 1 {
		log.Println("| ERROR | Не заданы таблицы для сканирования!")
		os.Exit(-1)
	}
	daemon.Run(cfg, tables)
}