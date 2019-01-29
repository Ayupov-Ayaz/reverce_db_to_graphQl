package reverser

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
)

/**
	Возвращает список с названиями зависимых таблиц
 */
func (r *Reverser) searchDependenciesTable(tables map[string]*model.Table,
									    relations map[string]*model.Relation, db *db.DB) []string {
	dependencies := make(map[string]bool, 0)
	for _, rel := range relations {
		if len(rel.LinkedTo) > 0 {
			for toTable := range rel.LinkedTo {
				if _, ok := tables[toTable]; !ok {
					dependencies[toTable] = true
				}
			}
		}
	}
	var results = make([]string, len(dependencies))
	for tName := range dependencies {
		results = append(results, tName)
	}
	return results
}
