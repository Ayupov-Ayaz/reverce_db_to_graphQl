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
		"IsPrimaryKey"      : field.IsPrimaryKey,
		"IsNullableField"	: field.IsNullableField,
		"GetValidate"       : field.GetValidate,
		"GetGraphQlType"	: field.GetGraphQlType,
		"GetForeignType" 	: model.GetForeignType,
	}

	t, err := template.New("tables").Funcs(funcMap).Parse(getTemplateStruct())
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

func getTemplateStruct() string {
	return `
	type {{.Name}} { {{range.Fields}}
		{{.Name}}: {{if .IsForeign}}{{GetForeignType .}}{{else}}{{.GetGraphQlType}}{{end}}{{.IsNullableField}} {{.IsPrimaryKey}} {{.GetValidate}}{{end}}
	}
`
}
