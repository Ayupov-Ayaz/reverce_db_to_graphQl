package daemon

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/commands"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/reverser"
	"log"
	"os"
)

type Config struct {
	Db *db.Config
}

func Run(cfg *Config, tables *[]string) {
	db, err := db.InitDB(cfg.Db)
	if err != nil {
		log.Printf("| SYS.ERROR | Ошибка при подключении к БД.(Проверьте данные подключения): \n %s",
			err.Error())
		os.Exit(-1)
	}
	// scanning
	rev := reverser.NewReverser( tables)
	com := commands.GetDbCommander(cfg.Db)
	rev.Reverse(db, com)
}