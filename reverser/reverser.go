package reverser

import (
	"bytes"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"io/ioutil"
	"os"
	"text/template"
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

	// TODO: передавать в шаблон
	sendToTemplate(&tableStructs)


	return nil
}

func sendToTemplate(tables *[]model.Table) {
	// Открытваем(Создаем) файл для записи с именем таблицы

	//fo, err := os.Create("result/1.txt")
	//if err != nil { panic(err) }
	//
	//defer func() {
	//	if err := fo.Close(); err != nil { panic(err) }
	//}()

	buff := bytes.NewBufferString("")

	template, err := template.New("Table").Parse("/templates/table.txt")
	if err != nil {
		panic(err)
	}
	if err := template.Execute(buff, tables); err != nil { panic(err) }
	if err := ioutil.WriteFile("result/1.txt", buff.Bytes(), os.ModePerm); err != nil { panic(err) }
	//err = tmpl.ExecuteTemplate(write, table)
}