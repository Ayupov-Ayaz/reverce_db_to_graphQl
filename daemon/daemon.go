package daemon

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/reverser"
	"log"
)

type Config struct {
	Db *db.Config
}

func Run(cfg *Config, tables *[]string) error {
	db, err := db.InitDB(cfg.Db)
	if err != nil {
		log.Println("ERROR| Initialization db: \n %s", err.Error())
	}
	// scanning
	rev := reverser.NewReverser( tables)
	rev.Reverse(db)

	return nil
}