package model

type ForeignKey struct {
	FieldName     string `db:"field_name"`
	FkInTable     string `db:"fk_in_table"`
	RefersToField string `db:"refers_to_field"`
}