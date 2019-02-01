package commands

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/commands/impl"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/errors"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
)

type DbCommander interface {
	GetTableStruct(tableName string, db *db.DB) (table *model.Table)
	GetForeignKeys(tableName string, db *db.DB) []*model.ForeignKey
	TableSearcher
	db.Paramser
}

type TableSearcher interface {
	GetAllTableNames(db *db.DB) []string
	GetTableByLike(tables []string, db *db.DB) []string
}

func GetDbCommander(cfg *db.Config) DbCommander{
	switch cfg.Driver {
	case db.DriverMSSQL:
		return impl.NewMssqlCommands()

	default:
		errors.PrintError(fmt.Sprintf("Для %s нет реализаций", cfg.Driver), true)
		return nil
	}
}