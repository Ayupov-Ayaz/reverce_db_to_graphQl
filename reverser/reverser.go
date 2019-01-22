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
		tableStructs = append(tableStructs, *tableStruct)
	}
	sendToTemplate(&tableStructs)
	return nil
}

func sendToTemplate(tables *[]model.Table) {
	t, err := template.New("tables").Funcs(model.GetFuncMap()).Parse(getTemplateStruct())
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