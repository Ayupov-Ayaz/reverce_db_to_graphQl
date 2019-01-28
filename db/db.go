package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type DB struct {
	*sqlx.DB
}


func InitDB(cfg *Config) (*DB, error){
	var dsn string

	switch cfg.Driver {
	case DriverMSSQL:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
	case DriverMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
	default:
		log.Printf("| ERROR | Вы указали не поддерживаемый sql driver: %s", cfg.Driver )
		os.Exit(-1)
	}

	db, err := sqlx.Connect(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg.Driver, " connected")
	return &DB{db}, nil
}