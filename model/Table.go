package model


type Table struct {
	Name string
	Fields []*Field
	ForeignKeys []*ForeignKey
}

