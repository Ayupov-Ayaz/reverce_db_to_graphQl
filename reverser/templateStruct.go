package reverser

func getTemplateStruct() string {
	return `type {{.Name}} { {{range.Fields}}
				{{.Name}}: {{if .IsForeign}}{{GetForeignType .}}{{else}}{{.GetGraphQlType}}{{end}}{{.IsNullableField}} {{.IsPrimaryKey}} {{.GetValidate}}{{end}}
		   }
`
}
