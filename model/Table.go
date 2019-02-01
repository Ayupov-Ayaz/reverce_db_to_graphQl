package model

import "strings"

type Table struct {
	Name          string
	Fields        []*Field
	ForeignKeys   []*ForeignKey
	ForeignFields map[string]*Field
}

func GetTableDirectivesByTable(t Table) string {
	directiveCollections := GetTableDirectiveCollection()
	var tableDirectives = make([]string, 0)
	for _, field := range t.Fields {
		if directive, exist := directiveCollections["fieldName"][field.Name]; exist {
			tableDirectives = append(tableDirectives, directive)
		}
	}
	if len(directiveCollections["table"]) > 0 {
		// Дополнить правилом
	}
	return strings.Join(tableDirectives, " ")
}

func GetTableDirectiveCollection() map[string]map[string]string {
	return map[string]map[string]string{
		"table" : {
			// регулярка или еще что нибудь
		},
		"fieldName" : {
			"del": 		"@SoftDelete",
			"deleted": 	"@SoftDelete",
			"delete": 	"@SoftDelete",
		},
	}
}