package main

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/daemon"
	"github.com/viper"
	"log"
	"os"
)


func getScanningParams() (tables []string, flags map[string]bool) {
	viper.SetConfigType("yaml")
	fileName := "tables.yaml"
	in, err := os.Open(fileName)

	if err != nil {
		log.Printf("| ERROR | Пожалуйста скопируйте файл %s под именем %s ", "tables.yaml.example", fileName)
		os.Exit(-1)
	}

	if err := viper.ReadConfig(in); err != nil {
		log.Printf("| SYS.ERROR | Возникла ошибка при чтении файла %s :\n %s", fileName, err.Error())
		os.Exit(-1)
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
		log.Println("| ERROR | Не заданы таблицы для сканирования!")
		os.Exit(-1)
	}
	daemon.Run(tables, flags)
}