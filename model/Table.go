package model


type Table struct {
	Name          string
	Fields        []*Field
	ForeignKeys   []*ForeignKey
	ForeignFields map[string]*Field
	Directives    []string
}

func GetTableDirectiveByTable(t Table) string {
	if len(t.Directives) == 0 {
		return ""
	}
	str := ""
	for _, directive := range t.Directives {
		str += " "+directive
	}
	return str
}

func GetTableDirectiveCollection() map[string]string {
	return map[string]string{
		"del": 		"@SoftDelete",
		"deleted": 	"@SoftDelete",
		"delete": 	"@SoftDelete",
	}
}