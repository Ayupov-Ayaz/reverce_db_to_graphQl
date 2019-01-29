package main

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/daemon"
	"log"
	"os"
)


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