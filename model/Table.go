package model

import (
	"html/template"
	"log"
	"strconv"
)

type Table struct {
	Name string
	Fields *[]Field
}

type Field struct {
	Name string `db:"name"`
	Type string `db:"type"`
	IsPrimary bool `db:"primary_key"`
	IsNullable bool `db:"is_nullable"`
	MaxLength int `db:"max_length"`
}

func (f *Field) IsPrimaryKey() string {
	if f.IsPrimary {
		return "@primary"
	}
	return ""
}

func (f *Field) IsNullableField() string {
	if f.IsNullable {
		return ""
	}
	return "!"
}

func (f *Field) GetGraphQlType() string {
	switch f.Type {
	case "varchar":  return "String"
	case "datetime": return "Datetime"
	case "date": 	 return "Date"
	case "int": 	 return "Int"
	case "float": 	 return "Float"
	case "tinyint":  return "Boolean"
	case "numeric":  return "Int"
	default:
		log.Printf("\n| ERROR | Не указан преобразователь типа для %s \n", f.Type)
		return f.Type
	}
}

func (f *Field) GetValidate() string {
	if f.MaxLength != 0 {
		return "max:" + strconv.Itoa(f.MaxLength)
	}
	return ""
}

func GetFuncMap() template.FuncMap {
	field := Field{}
	return template.FuncMap{
		"IsPrimaryKey"      : field.IsPrimaryKey,
		"IsNullableField"	: field.IsNullableField,
		"GetValidate"       : field.GetValidate,
		"GetGraphQlType"	: field.GetGraphQlType,
	}
}