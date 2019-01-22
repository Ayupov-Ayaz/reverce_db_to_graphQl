package model

import (
	"fmt"
	"log"
	"strconv"
)

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
	default:
		log.Printf("\n| ERROR | Не указан преобразователь типа для %s, тип найден у поля %s\n", f.Type, f.Name)
		return f.Type
	}
}

func GetForeignType(f *Field)  string {
	fmt.Println(f)
	return ""
}

func (f *Field) GetValidate() string {
	if f.MaxLength != 0 {
		return "@validate(max:" + strconv.Itoa(f.MaxLength)+")"
	}
	return ""
}