package commands

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/commands/impl"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"log"
	"os"
)

type DbCommander interface {
	GetTableStruct(tableName string, db *db.DB) (table *model.Table)
	GetForeignKeys(tableName string, db *db.DB) []*model.ForeignKey
	db.Paramser
}

func GetDbCommander(cfg *db.Config) DbCommander{
	switch cfg.Driver {
	case db.DriverMSSQL:
		return impl.NewMssqlCommands()

	default:
		log.Printf("| ERROR | Для %s нет реализаций", cfg.Driver)
		os.Exit(-1)
		return nil
	}
}