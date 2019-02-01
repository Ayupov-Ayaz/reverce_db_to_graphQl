package errors

import (
	"log"
	"os"
)

type dbError struct {
	s string
}

type fatalError struct {
	s string
}

type notice struct {
	s string
}

type error struct {
	s string
}

func (e *notice) Error() string {
	return "| NOTICE | " + e.s
}

func (e *fatalError) Error() string {
	return "| FATAL_ERROR | " + e.s
}

func (e *dbError) Error() string {
	return "| DB_ERROR | " + e.s
}

func (e *error) Error() string {
	return "| ERROR | " + e.s
}

func NewError(message string) *error {
	return &error{ s: message }
}

func NewFatalError(message string) *fatalError {
	return &fatalError{ s: message }
}

func NewNotice(message string) *notice {
	return &notice{ s: message }
}

func NewDbError(message string) *dbError {
	return &dbError{ s: message }
}

func PrintNotice(message string) {
	log.Println(NewNotice(message).Error())
}

func PrintError(message string, die bool) {
	log.Println(NewError(message).Error())
	if die {
		os.Exit(-1)
	}
}

func PrintFatalError(message string, die bool) {
	log.Println(NewFatalError(message).Error())
	if die {
		os.Exit(-1)
	}
}

func PrintDbError(message string, die bool) {
	log.Println(NewDbError(message).Error())
	if die {
		os.Exit(-1)
	}
}