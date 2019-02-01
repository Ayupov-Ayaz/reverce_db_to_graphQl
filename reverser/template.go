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
func sendToTemplate(tables map[string]*model.Table) {

	field := model.Field{}
	funcMap :=  template.FuncMap{
		"IsPrimaryKey"      		: field.IsPrimaryKey,
		"IsNullableField"			: field.IsNullableField,
		"GetDirectories"       		: field.GetDirectories,
		"GetGraphQlType"			: field.GetGraphQlType,
		"GetForeignType" 			: model.GetForeignType,
		"GetTableDirectivesByTable"	: model.GetTableDirectivesByTable,
	}

	t, err := template.New("tables").Funcs(funcMap).Parse(getTemplate())
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
