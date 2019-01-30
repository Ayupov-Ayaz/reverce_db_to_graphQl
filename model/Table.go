package model


type Table struct {
	Name          string
	Fields        []*Field
	ForeignKeys   []*ForeignKey
	ForeignFields map[string]*Field
	Directives    []string
}


func GetTableDirectiveCollection() map[string]string {
	return map[string]string{
		"del": 		"@SoftDelete",
		"deleted": 	"@SoftDelete",
		"delete": 	"@SoftDelete",
	}
}