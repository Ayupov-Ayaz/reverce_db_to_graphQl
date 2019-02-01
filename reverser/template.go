package reverser

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"log"
	"os"
	"text/template"
)

/**
	Выгружает данные в шаблон и создает в папке results graphql структуры
 */
func (r *Reverser) sendToTemplate(tablesFromSearching map[string]*model.Table, flags map[string]bool) {
	// Если нужно отобразить таблицы без зависимостей
	var tables = make(map[string]*model.Table, 0)
	if flags["l"] && !flags["d"] {
		for _, tableName := range r.TablesForShow {
			if tableModel, exist := tablesFromSearching[tableName]; exist {
				tables[tableName] = tableModel
			}
		}
	} else {
		tables = tablesFromSearching
	}

	field := model.Field{}
	funcMap :=  template.FuncMap{
		"IsPrimaryKey"      		: field.IsPrimaryKey,
		"IsNullableField"			: field.IsNullableField,
		"GetDirectories"       		: field.GetDirectories,
		"GetGraphQlType"			: field.GetGraphQlType,
		"GetForeignType" 			: model.GetForeignType,
		"GetTableDirectivesByTable"	: model.GetTableDirectivesByTable,
	}

	t, err := template.New("tablesFromSearching").Funcs(funcMap).Parse(getTemplate())
	if err != nil {
		log.Printf("| SYS.ERROR | Не удалось создать шаблон \n")
		panic(err)
	}
	for _, table := range tables {
		fileName := "results/" + table.Name + ".graphql"
		fo, err := os.Create(fileName)
		if err != nil {
			log.Printf("| SYS.ERROR | Не удалось открыть файл %s, \n", fileName)
			panic(err)
		}

		if err := t.Execute(fo, table); err != nil {
			log.Printf("| SYS.ERROR | При попытке записать в файл %s структуру таблицы %s произошла ошибка:\n",
				fileName, table.Name)
			panic(err)
		}

		if err := fo.Close(); err != nil {
			log.Printf("| SYS.ERROR | Не удалось закрыть файл %s \n", fileName)
		}
	}

}

func getTemplate() string {
	return `
	type {{.Name}} {{GetTableDirectivesByTable .}} { {{range.Fields}}
		{{.Name}}: {{.IsPrimaryKey}}{{if .IsForeign}} {{GetForeignType .}}{{else}} {{.GetGraphQlType}}{{end}}{{.IsNullableField}} {{.GetDirectories}}{{end}}
	}
`
}
