package reverser

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"html/template"
	"log"
	"os"
)
/**
  Structure for reversing DataBase structure to GraphQl structure
 */
type Reverser struct {
	OutputFileName string
	Tables []string
}

func NewReverser(tables *[]string) *Reverser {
	return &Reverser{
		Tables: *tables,
	}
}

func (r *Reverser) Reverse(db *db.DB) error {
	var tableStructs = make([]model.Table, 0)

	// Получаем структуры таблиц
	for _, table := range r.Tables {
		tableStruct, err := getTableStruct(table, db)
		if err != nil {
			panic(err)
		}
		// достаем внешние ключи
		tableStruct.ForeignKeys = GetForeignKeys(table, db)
		tableStructs = append(tableStructs, *tableStruct)
	}
	// отправляем в шаблон
	sendToTemplate(&tableStructs)
	return nil
}

func sendToTemplate(tables *[]model.Table) {


	field := model.Field{}

	funcMap :=  template.FuncMap{
		"IsPrimaryKey"      : field.IsPrimaryKey,
		"IsNullableField"	: field.IsNullableField,
		"GetValidate"       : field.GetValidate,
		"GetGraphQlType"	: field.GetGraphQlType,
		"GetForeignType" 	: model.GetForeignType,
	}

	t, err := template.New("tables").Funcs(funcMap).Parse(getTemplateStruct())
	if err != nil { panic(err) }
	for _,table := range *tables {
		fileName := "results/" + table.Name + ".graphQl"
		fo, err := os.Create(fileName)
		if err != nil {log.Printf("| ERROR | Не удалось открыть файл %s, \n", fileName)}

		defer func() {
			if err := fo.Close(); err != nil {log.Printf("| ERROR | Не удалось закрыть файл %s \n", fileName)}
		}()

		if err := t.Execute(fo, table); err != nil {
			panic(err)
		}
	}

}