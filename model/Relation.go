package model

type Relation struct {
	/* { map[TableName] -> map[  fk -> field_name, pk -> field_name ] } */
	LinkedTo map[string]map[string]string
	LinkToMe *map[string]map[string]string
}

func (r *Relation) AddLinkedTo(tableToLink, fkField, pkField string){
	newLinkedTo := make(map[string]string)
	newLinkedTo["fk"] = fkField
	newLinkedTo["pk"] = pkField
	if len(r.LinkedTo) == 0 {
		r.LinkedTo = make(map[string]map[string]string)
	}
	r.LinkedTo[tableToLink] = newLinkedTo
}
