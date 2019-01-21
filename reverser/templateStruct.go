package reverser

func getTemplateStruct() string {
	return `
    type {{.Name}} {
		{{range.Fields}}
		{{.Name}}: {{.Type}}! {{.IsPrimary}} {{.IsNullable}} {{.MaxLength}}
	{{end}}
	}
`
}
