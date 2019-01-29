package db

import (
	"fmt"
	"github.com/viper"
	"log"
	"os"
	"strconv"
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
		log.Printf("| ERROR | Пожалуйста скопируйте файл \"%s\" под именем \"%s\" ", "example.db.yaml", fileName)
		os.Exit(-1)
	}

	if err := viper.ReadConfig(in); err != nil {
		log.Printf("| SYS.ERROR | Возникла ошибка при чтении файла %s :\n %s", fileName, err.Error())
		os.Exit(-1)
	}

	dbParams := viper.GetStringMapString("db")

	if valid, errors := validateDbParams(dbParams); !valid {
		log.Printf("У Вас имеются ошибки в заполнении файла конфигурации: \"%s\"", fileName)
		for _, err := range errors {
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
	if len(s) == 0 || s == "******" {
		return true
	}
	return false
}