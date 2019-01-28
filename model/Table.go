package model


type Table struct {
	Name string
	Fields []*Field
	ForeignKeys []*ForeignKey
	ForeignFields map[string]*Field
}

