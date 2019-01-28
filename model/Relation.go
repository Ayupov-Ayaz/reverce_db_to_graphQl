package model

type Relation struct {
	/* { map[TableName] -> map[  fk -> field_name, pk -> field_name ] } */
	LinkedTo map[string]map[string]string
	LinkToMe map[string]map[string]string
}

func (r *Relation) AddLinkedTo(tableToLink, fkField, pkField string) {
	newLinkedTo := r.makeLinkMap(fkField, pkField)
	if len(r.LinkedTo) == 0 {
		r.LinkedTo = make(map[string]map[string]string)
	}
	r.LinkedTo[tableToLink] = newLinkedTo
}

func (r *Relation) AddLinkToMe(tableToLink, fkField, pkField string) {
	newLinkToMe := r.makeLinkMap(fkField, pkField)
	if len(r.LinkToMe) == 0 {
		r.LinkToMe = make(map[string]map[string]string)
	}
	r.LinkToMe[tableToLink] = newLinkToMe
}

func (r *Relation) makeLinkMap(fkField, pkField string) map[string]string {
	newLinkMap := make(map[string]string)
	newLinkMap["fk"] = fkField
	newLinkMap["pk"] = pkField
	return newLinkMap
}