package daemon

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"log"
)

type Config struct {
	Db *db.Config
}

func Run(cfg *Config) error {
	db, err := db.InitDB(cfg.Db)
	if err != nil {
		log.Println("ERROR| Initialization db: \n %s", err.Error())
	}
	fmt.Println(db)
	return nil
}