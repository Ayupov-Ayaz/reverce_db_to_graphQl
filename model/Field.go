package model

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/errors"
	"strconv"
	"strings"
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
	if (!f.IsNullable || f.IsPrimary) && !f.IsForeign {
		// для внешних ключей вычисляется в функции reverser.DefiningTableRelations()
		return "!"
	}
	return ""
}

func (f *Field) GetGraphQlType() string {
	switch f.Type {
	case "varchar", "nvarchar", "nchar", "char", "text", "timestamp", "ntext", "image", "varbinary":  return "String"
	case "int", "bigint", "real", "numeric", "smallint" : 	 return "Int"
	case "datetime", "datetime2", "smalldatetime": return "Datetime"
	case "tinyint", "bit":  return "Boolean"
	case "date": 	 return "Date"
	case "float", "money", "decimal": 	 return "Float"
	case "time":	 return "Time"

	default:
		errors.PrintError(fmt.Sprintf("| ERROR | Не указан преобразователь типа для %s, тип найден у поля %s\n",
			f.Type, f.Name), false)
		return f.Type
	}
}

func GetForeignType(f *Field)  string {
	return f.FkType
}

func (f *Field) GetDirectories() string {
	directories := make([]string, 0)
	if f.MaxLength != 0 && f.MaxLength != -1 {
		directories = append(directories, "@validate(max:" + strconv.Itoa(f.MaxLength)+")")
	}
	return strings.Join(directories, " ")
}