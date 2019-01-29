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
	var results = make([]string, 0)
	for tName := range dependencies {
		results = append(results, tName)
	}
	return results
}

/**
	Добавление зависимой карты структур в главную карту структур
 */
func (r *Reverser) addedDependenciesToTableCollection(mainCollection map[string]*model.Table,
													    dependencies map[string]*model.Table) {
	for tableName := range dependencies {
		if _, exists := mainCollection[tableName]; !exists {
			mainCollection[tableName] = dependencies[tableName]
		}
	}
}

/**
	Добавление зависимой карты связей в главную карту связей
 */
func (r *Reverser) addedDependenciesToTableRelation(mainRelations map[string]*model.Relation,
	dependencies map[string]*model.Relation) {
	for tableName := range dependencies {
		if _, exists := mainRelations[tableName]; !exists {
			mainRelations[tableName] = dependencies[tableName]
		}
	}
}