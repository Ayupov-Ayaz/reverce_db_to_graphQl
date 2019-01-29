package db

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type DB struct {
	*sqlx.DB
}

type Params struct {
	Version string `db:"product_version"`
	Level string `db:"product_level"`
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
		log.Printf("| ERROR | Вы указали неподдерживаемый sql driver: %s", cfg.Driver )
		os.Exit(-1)
	}

	db, err := sqlx.Connect(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg.Driver, " connected")
	return &DB{db}, nil
}

func (db *DB) GetParams() *Params{
	p := &Params{}
	query := `
		SELECT SERVERPROPERTY('ProductVersion') AS product_version
     	, SERVERPROPERTY('ProductLevel')   AS product_level;
	`
	if err := db.Get(p, query); err != nil {
		log.Printf("| ERROR | %s \n", err.Error())
	}
	return p
}

func (db *DB) CompareDbParams(dbCommands Paramser, flags map[string]bool) bool {
	realized := dbCommands.GetParams()
	server_params := db.GetParams()
	err_message := "| ERROR | Вы используете %s: %s, реализация построена для - \"%s\""
	ok := true
	// Проверка версий базы данных
	if realized.Version != server_params.Version {
		ok = false
		// проверка совместимости
		for _, version := range dbCommands.GetSupportedVersions() {
			if version == server_params.Version {
				ok = true
				break
			}
		}
		if flags["f"] {
			ok = true
		}else if !ok {
			log.Printf(err_message, "версию", server_params.Version, realized.Version)
		}
	}
	// проверка пакета обновлений
	if realized.Level != server_params.Level {
		log.Printf(err_message, " пакет обновлений ", server_params.Level, realized.Level)
		ok = false
	}
	if !ok && !flags["f"]{
		log.Printf("Если Вы всё равно желаете запусть программу, передайте ключ \"%s\"", "f")
	}
	return ok
}

func (db *DB)GetSupportedVersions() []string {
	return []string{
		// заглушка
	}
}