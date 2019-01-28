package model

type ForeignKey struct {
	FieldName string `db:"field_name"`
	FkToTable string `db:"fk_to_table"`
	PkField string 	`db:"pk_field"`
}