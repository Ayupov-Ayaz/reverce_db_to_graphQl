package main

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/daemon"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/errors"
	"github.com/viper"
	"os"
)


func getScanningParams() (tables []string, flags map[string]bool) {
	viper.SetConfigType("yaml")
	fileName := "tables.yaml"
	in, err := os.Open(fileName)

	if err != nil {
		errors.PrintError(fmt.Sprintf("Пожалуйста скопируйте файл %s под именем %s ",
			"tables.yaml.example", fileName), true)
	}

	if err := viper.ReadConfig(in); err != nil {
		errors.PrintFatalError(fmt.Sprintf("Возникла ошибка при чтении файла %s :\n %s",
			fileName, err.Error()), true)
	}

	tables = viper.GetStringSlice("tables")
	flagsSlice := viper.GetStringSlice("flags")
	flags = make(map[string]bool, 0)
	for _, flag := 	range flagsSlice {
		if flag == "d" {
			flags["d"] = true
		}
		if flag == "l" {
			flags["l"] = true
		}
		if flag == "*" {
			flags["*"] = true
		}
		if flag == "f" {
			flags["f"] = true
		}
	}
	return
}

func main() {
	tables, flags := getScanningParams()
	if len(tables) < 1 && !flags["*"]{
		errors.PrintError("Не заданы таблицы для сканирования!", true)

	}
	daemon.Run(tables, flags)
}