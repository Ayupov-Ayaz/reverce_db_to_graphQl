package model

type Relation struct {
	/* { map[TableName] -> map[  fk -> field_name, pk -> field_name ] } */
	LinkedTo map[string]*RelationKey
	LinkToMe map[string]*RelationKey
}

type RelationKey struct {
	FieldsRk map[string]string // ForeignKey -> PrimaryKey
}

func (r *Relation) AddLinkedTo(tableToLink, fkField, pkField string) {
	if len(r.LinkedTo) == 0 {
		r.LinkedTo = make(map[string]*RelationKey)
	}
	newLinkedTo := r.newRelationKey(fkField, pkField)

	if _, exist := r.LinkedTo[tableToLink]; !exist || r.LinkedTo[tableToLink] == nil {
		r.LinkedTo[tableToLink] = newLinkedTo
	} else  { // если уже существует такая таблица, то проверяем внешние ключи
		r.LinkedTo[tableToLink].MergeRelationKey(newLinkedTo)
	}
}

// Добавление новой ссылки на связь таблицы
func (r *Relation) AddLinkToMe(tableToLink, fkField, pkField string) {
	newLinkToMe := r.newRelationKey(fkField, pkField)
	if len(r.LinkToMe) == 0 {
		r.LinkToMe =  make(map[string]*RelationKey)
	}

	if _,exist := r.LinkedTo[tableToLink]; !exist || r.LinkToMe[tableToLink] == nil {
		r.LinkToMe[tableToLink] = newLinkToMe
	} else {
		r.LinkToMe[tableToLink].MergeRelationKey(newLinkToMe)
	}
}

func (r *Relation) newRelationKey(fkField, pkField string) *RelationKey {
	return &RelationKey{
		map[string]string{
			fkField :  pkField,
		},
	}
}

func (rk *RelationKey) MergeRelationKey(depRelKey *RelationKey ) {
	for fkField, pkField := range depRelKey.FieldsRk {
		if _, exist := rk.FieldsRk[fkField]; !exist { //fkField != rk.FieldsRk[fkField]
			rk.FieldsRk[fkField] = pkField
		}
	}
}