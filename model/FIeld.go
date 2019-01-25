package model

import (
	"log"
	"strconv"
)

type Field struct {
	Name string `db:"name"`
	Type string `db:"type"`
	IsPrimary bool `db:"is_primary"`
	IsForeign bool `db:"is_foreign"`
	IsNullable bool `db:"is_nullable"`
	IsUnique bool `db:"is_unique"`
	MaxLength int `db:"max_length"`
	FkType string
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
	case "char":	 return "String"
	case "bit":		 return "Boolean"
	case "bigint":	 return "Int"
	case "nchar":	 return "String"
	case "nvarchar": return "String"
	case "time":	 return "Time"
	case "real":	 return "Int"


	default:
		log.Printf("\n| ERROR | Не указан преобразователь типа для %s, тип найден у поля %s\n", f.Type, f.Name)
		return f.Type
	}
}

func GetForeignType(f *Field)  string {
	return "// TODO:Foreign "
}

func (f *Field) GetValidate() string {
	if f.MaxLength != 0 {
		return "@validate(max:" + strconv.Itoa(f.MaxLength)+")"
	}
	return ""
}