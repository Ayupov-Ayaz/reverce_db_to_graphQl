package db

import (
	"log"
	"os"
	"regexp"
)

type Ip string

func (ip *Ip) validate() bool {
	reg, err := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if err != nil {
		log.Printf("| SYS.ERROR | При создании правила валидации для Ip произошла ошибка:\n %s", err.Error())
		os.Exit(-1)
	}

	if reg.Match([]byte(string(*ip))) {
		return true
	}
	return false
}