package model

type Table struct {
	Name string
	Fields *[]Field
}

type Field struct {
	Name string `db:"name"`
	Type string `db:"type"`
	IsPrimary bool `db:"primary_key"`
	IsNullable bool `db:"is_nullable"`
	MaxLength *int `db:"max_length"`
}