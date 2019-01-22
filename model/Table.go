package model


type Table struct {
	Name string
	Fields *[]Field
	ForeignKeys *[]ForeignKey
}

type Field struct {
	Name string `db:"name"`
	Type string `db:"type"`
	IsPrimary bool `db:"is_primary"`
	IsForeign bool `db:"is_foreign"`
	IsNullable bool `db:"is_nullable"`
	MaxLength int `db:"max_length"`
}
