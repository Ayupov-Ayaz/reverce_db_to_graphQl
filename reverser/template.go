package reverser

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/errors"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"os"
	"text/template"
)

/**
	Выгружает данные в шаблон и создает в папке results graphql структуры
 */
func (r *Reverser) sendToTemplate(tablesFromSearching map[string]*model.Table, flags map[string]bool) {
	// Если нужно отобразить таблицы без зависимостей
	var tables = make(map[string]*model.Table, 0)
	if flags["l"] && !flags["d"] && !flags["*"] {
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
		"IssetInverseRelations"		: model.IssetInverseRelations,
	}

	t, err := template.New("tablesFromSearching").Funcs(funcMap).Parse(getTemplate())
	if err != nil {
		errors.PrintFatalError(fmt.Sprintf("Не удалось создать шаблон \n %s", err.Error()), true)
	}
	var i int
	for _, table := range tables {
		i++
		fileName := "results/" + table.Name + ".graphql"
		fo, err := os.Create(fileName)
		if err != nil {
			errors.PrintFatalError(fmt.Sprintf("Не удалось открыть файл %s, \n", fileName), true)

		}

		if err := t.Execute(fo, table); err != nil {
			errors.PrintFatalError(fmt.Sprintf("При попытке записать в файл %s структуру таблицы %s произошла ошибка:\n",
				fileName, table.Name), true)
		}

		if err := fo.Close(); err != nil {
			errors.PrintError(fmt.Sprintf("Не удалось закрыть файл %s \n", fileName), false)
		}
	}
	var res string
	if i > 1 {
		res = "Создано %d типа(ов) \n"
	} else {
		res = "Создан %d тип \n"
	}
	fmt.Printf(res, i)

}

func getTemplate() string {
	return `
	type {{.Name}} {{GetTableDirectivesByTable .}} { {{range.Fields}}
		{{.Name}}: {{.IsPrimaryKey}}{{if .IsForeign}} {{GetForeignType .}}{{else}} {{.GetGraphQlType}}{{end}}{{.IsNullableField}} {{.GetDirectories}}{{end}}
		{{if IssetInverseRelations .}}
		#inverse relations{{range $key, $value := .InverseRelations}}
		{{$key}}: {{$value}}{{end}}{{end}}
	}
`
}
