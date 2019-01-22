package reverser

func getTemplateStruct() string {
	return `type {{.Name}} { {{range.Fields}}
				{{.Name}}: {{.GetGraphQlType}}{{.IsNullableField}} {{.IsPrimaryKey}} {{.GetValidate}} {{end}}
		   }
`
}
