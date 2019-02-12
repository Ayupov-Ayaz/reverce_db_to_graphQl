package model

import "strings"

type Table struct {
	Name          string
	Fields        []*Field
	ForeignKeys   []*ForeignKey
	ForeignFields map[string]*Field
	InverseRelations map[string]string
}

func GetTableDirectivesByTable(t Table) string {
	directiveCollections := GetTableDirectiveCollection()
	var tableDirectives = make([]string, 0)
	for i := 0; i < len(t.Fields); {
		field := t.Fields[i]
		directive, exist := directiveCollections["fieldName"][field.Name]
		if exist {
			tableDirectives = append(tableDirectives, directive)
		}

		if directive == "@softDelete" {
			t.Fields = append(t.Fields[:i], t.Fields[i+1:]...)
		} else {
			i++
		}
	}
	if len(directiveCollections["table"]) > 0 {
		// Дополнить правилом
	}
	return strings.Join(tableDirectives, " ")
}

func GetTableDirectiveCollection() map[string]map[string]string {
	softDelete := "@softDelete"
	return map[string]map[string]string{
		"table" : {
			// регулярка или еще что нибудь
		},
		"fieldName" : {
			"del": 		softDelete,
			"deleted": 	softDelete,
			"delete": 	softDelete,
		},
	}
}

func IssetInverseRelations(t *Table) bool {
	return len(t.InverseRelations) > 0
}