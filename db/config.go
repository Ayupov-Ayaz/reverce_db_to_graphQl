package db

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	DriverMSSQL = "mssql"
	DriverMySQL  = "mysql"
)

type Config struct {
	Driver string
	User string
	Password string
	Host Ip
	Database string
	Port int
}

func GetDbConfig() *Config {
	viper.SetConfigType("yaml")
	fileName := ".db.yaml"
	in, err := os.Open(fileName)
	if err != nil {
		errors.PrintError(fmt.Sprintf("Пожалуйста скопируйте файл \"%s\" под именем \"%s\" ",
			"example.db.yaml", fileName), true)
	}

	if err := viper.ReadConfig(in); err != nil {
		errors.PrintFatalError(fmt.Sprintf("| SYS.ERROR | Возникла ошибка при чтении файла %s :\n %s",
			fileName, err.Error()), true)
	}

	dbParams := viper.GetStringMapString("db")

	if valid, errs := validateDbParams(dbParams); !valid {
		errors.PrintError(fmt.Sprintf("У Вас имеются ошибки в заполнении файла конфигурации: \"%s\"",
			fileName), false)
		for _, err := range errs {
			log.Println(err)
		}
		os.Exit(-1)
	}
	port, _ := strconv.Atoi(dbParams["port"]) // на ошибку проверяется в валидаторе

	cfg := &Config{
		Driver:   dbParams["driver"],
		User: 	  dbParams["user"],
		Password: dbParams["password"],
		Host:	  Ip(dbParams["host"]),
		Database: dbParams["database"],
		Port:	  port,
	}

	return cfg
}

func validateDbParams(dbParams map[string]string) (valid bool, errs []string) {
	errs = make([]string, 0)
	valid = true
	errorsMessages := make(map[string]string)
	errorsMessages["empty"] = "Значение \"%s\" не указано!"
	errorsMessages["not_valid"] = "Значение \"%s\" не валидно!"

	if isEmptyDbParam(dbParams["database"]) {
		errs = append(errs, fmt.Sprintf(errorsMessages["empty"], "database"))
		valid = false
	}

	if isEmptyDbParam(dbParams["driver"]) {
		errs = append(errs, fmt.Sprintf(errorsMessages["empty"], "driver"))
		valid = false
	}

	if isEmptyDbParam(dbParams["user"]) {
		errs = append(errs, fmt.Sprintf(errorsMessages["empty"], "user"))
		valid = false
	}

	if isEmptyDbParam(dbParams["password"]) {
		errs = append(errs, fmt.Sprintf(errorsMessages["empty"], "password"))
		valid = false
	}

	ip := Ip(dbParams["host"])
	if isEmptyDbParam(dbParams["host"]) {
		errs = append(errs, fmt.Sprintf(errorsMessages["empty"], "host"))
		valid = false
	} else if !ip.validate() {
		errs = append(errs, fmt.Sprintf(errorsMessages["not_valid"], "host"))
		valid = false
	}

	if isEmptyDbParam(dbParams["port"]) {
		errs = append(errs, fmt.Sprintf(errorsMessages["empty"], "port"))
		valid = false
	} else if len(dbParams["port"]) < 4 {
		errs = append(errs, fmt.Sprintf(errorsMessages["not_valid"], "port"))
		valid = false
	} else {
		if _, err := strconv.Atoi(dbParams["port"]); err != nil {
			errs = append(errs, fmt.Sprintf(errorsMessages["not_valid"], "port"))
			valid = false
		}
	}
	return
}

func isEmptyDbParam(s string) bool {
	if len(strings.TrimSpace(s)) == 0 {
		return true
	}
	return false
}