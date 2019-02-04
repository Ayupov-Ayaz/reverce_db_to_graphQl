package reverser

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/commands"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/errors"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"log"
)
/**
  Structure for reversing DataBase structure to GraphQl structure
 */
type Reverser struct {
	OutputFileName  string
	TablesForSearch []string
	TablesForShow []string // используется если передан флаг l и не передан флаг d
}

func NewReverser(tables []string) *Reverser {
	return &Reverser{
		TablesForSearch: tables,
	}
}

/**
	Функция которая запускает процессы:
	1) Получает структуры таблиц в виде map[tableName]*model.Table
	2) Получает отношения таблиц в виде map[tableName]*mode.Relation
	3) Проставляет спец типы полей опираясь на эти отношения
	4) выгрузка в шаблон результатов
 */
func (r *Reverser) Reverse(db *db.DB, com commands.DbCommander, flags map[string]bool) {
	var tableCollection = make(map[string]*model.Table, 0)
	var tableRelations = make(map[string]*model.Relation, 0)

	tableCollection, tableRelations = r.getTableData(tableCollection, tableRelations, com, db, flags)

	SpecialTypeDefinition(tableCollection, tableRelations)
	// отправляем в шаблон
	r.sendToTemplate(tableCollection, flags)
}

/**
	Получение структуры таблиц и их отношений рекурсивно получает все таблицы к которым есть внешние ключи
 */
func (r *Reverser) getTableData(tCollection map[string]*model.Table, tRelations map[string]*model.Relation,
	com commands.DbCommander, db *db.DB, flags map[string]bool) (
	tableCollection map[string]*model.Table, tableRelation map[string]*model.Relation){

	tableStructs := r.getTableStructs(db, com, flags)

	if len(r.TablesForSearch) != len(tableStructs) {
		deleteNotFoundTables(r.TablesForSearch, tableStructs)
	}

	if len(tableStructs) == 0 && len(tCollection) == 0 { // сработает при первом проходе
		errors.PrintError(fmt.Sprintln("Ни одна таблица,которую вы указали, не была найдена." +
			" Проверьте названия таблиц."), true)
	}

	// Создаем карту отношений таблиц
	tableRelation = DefiningTableRelations(tableStructs)
	// Делаем из среза карту
	tableCollection = makeTableCollection(tableStructs)

	if !flags["*"] {
		for {
			dependencies := r.searchDependenciesTable(tableCollection, tableRelation, db)
			if len(dependencies) == 0 {
				break
			}
			r.TablesForSearch = dependencies
			tCol, tRel := r.getTableData(tableCollection, tableRelation, com, db, flags)
			r.addedDependenciesToTableCollection(tableCollection, tCol)
			r.addedDependenciesToTableRelation(tableRelation, tRel)
		}
	}
	return
}
/**
	Просматривает все таблички и создает карту отношений таблиц друг к другу
 */
func DefiningTableRelations(tables []*model.Table) map[string]*model.Relation{

	tableRelation := make(map[string]*model.Relation)
	// пробегаем по всем табличкам
	for _, table := range tables {
		// если есть внешние ключи
		if len(table.ForeignKeys) > 0 {
			// пробегаем по полям у которых эти внешние ключи есть
			for _, fkField := range table.ForeignKeys {
				// ищем поле в котором есть этот внешний ключ
				for _, field := range table.Fields {
					if field.Name == fkField.FieldName {

						// relation
						relation, ok := tableRelation[table.Name]
						if !ok{
							relation = &model.Relation{}
							tableRelation[table.Name] = relation
						}
						relation.AddLinkedTo(fkField.FkInTable, fkField.FieldName, fkField.RefersToField)

						// inverseRelation
						inverseRelation, ok := tableRelation[fkField.FkInTable]
						if !ok{
							inverseRelation = &model.Relation{}
							tableRelation[fkField.FkInTable] = inverseRelation
						}
						inverseRelation.AddLinkToMe(table.Name, fkField.FieldName, fkField.RefersToField)
						
					}
				}
			}
		}
	}
	return tableRelation
}

/**
	Определение специальных типов у полей исходя из внешних ключей к таблицам которым они относятся
 */
func SpecialTypeDefinition(tables map[string]*model.Table, relations map[string]*model.Relation)  {
	// пробегаем по всем табличкам
	for tableName, table := range tables {
		// если есть внешние ключи
		if len(table.ForeignKeys) > 0 {
			// получаем все ссылки на другие таблицы по текущей таблице
			linkedTo := relations[table.Name].LinkedTo
			// пробегаем по отношениям таблиц к которым ссылается текущая таблица
			for toTable, relKey := range linkedTo {
				// пробегаем по всем полям нашей таблицы  которые яв-ся внешними ключами
				for _, field := range table.ForeignFields {
					// проверяем, есть ли какие нибудь отношения по внешнему ключу
					if _, exist :=  relKey.FieldsRk[field.Name];	exist {
						// получаем таблицу к которой у нас отношение
						if linkToTable, ok := tables[toTable]; ok {

								// ищем поле на которое ссылается наш внешний ключ
								for _, toField := range linkToTable.Fields {
									// находим поле на которое ссылается наш ForeignKey
									if relKey.FieldsRk[field.Name] == toField.Name {
										if field.IsUnique {
											field.FkType = toTable
										} else if !field.IsUnique && toField.IsUnique ||
											!field.IsUnique && toField.IsPrimary {
											field.FkType = "[" + toTable + "!]"
										}  else if !field.IsUnique && !toField.IsUnique {
											field.FkType = "[" + toTable + "]"
										} else {
											fmt.Printf("| NOTICE | Проверить отношение: %s.%s => %s.%s  %s\n",
												tableName, field.Name, toTable, toField.Name,
												"(для поля %s прописан тип NO_TABLE_SPECIFIED)")
											field.FkType = "NO_TABLE_SPECIFIED"
										}
									}
								}
						} else {
							relationError := " Не удалось определить отношение %s.%s => %s " +
								"(для поля %s прописан тип NO_TABLE_SPECIFIED)"
							tableNotFound := relationError + ", Таблица '%s' не была указана \n"
							relationError2 := relationError + ", Причина не ясна =( "
							if _, tableExist := tables[toTable]; !tableExist {
								errors.PrintNotice(fmt.Sprintf(tableNotFound, tableName, field.Name, toTable,
									field.Name, toTable))
							} else {
								errors.PrintNotice(fmt.Sprintf(relationError2, tableName, field.Name, toTable,
									field.Name))
							}
							field.FkType = "NO_TABLE_SPECIFIED"
						}
					}
				}
			}
		}
	}
}

// Преобразует slice таблиц в карту таблиц, где ключом является имя таблицы
func makeTableCollection(tablesSlice []*model.Table) map[string]*model.Table{
	tableCollection := make(map[string]*model.Table)
	for _, modelTable := range tablesSlice {
		tableCollection[modelTable.Name] = modelTable
		if len(modelTable.ForeignKeys) > 0 {
			foreignFields := make(map[string]*model.Field)
			for _, field := range modelTable.Fields {
				if field.IsForeign {
					foreignFields[field.Name] = field
				}
			}
			modelTable.ForeignFields = foreignFields
		}
	}
	return tableCollection
}

/**
	Если не найдена какая-нибудь таблица, ищем ее из общего списка, оповещаем пользователя, удаляем из списка,
	что бы не проводить дальнейшие манипуляции с ней
 */
func deleteNotFoundTables(searchingTables []string, tables []*model.Table) []string{

	for i := 0; i < len(searchingTables); {
		founded := false
		for j := 0; j < len(tables); j++ {
			if searchingTables[i] == tables[j].Name {
				i++
				founded = true
				break
			}
		}
		if !founded {
			log.Printf("| NOTICE | Таблица %s не была найдена в бд. Проверьте название!", searchingTables[i])
			searchingTables = append(searchingTables[:i], searchingTables[i+1:]...)
		}
	}
	return searchingTables
}

func (r *Reverser) getTableStructs(db *db.DB, com commands.DbCommander, flags map[string]bool) []*model.Table {
	var tableStructs = make([]*model.Table, 0)

	// Получаем структуры таблиц
	for _, table := range r.TablesForSearch {
		tableStruct := com.GetTableStruct(table, db)
		if tableStruct != nil && len(tableStruct.Fields) > 0 {
			// достаем внешние ключи
			tableStruct.ForeignKeys = com.GetForeignKeys(table, db)
			tableStructs = append(tableStructs, tableStruct)
		}
	}
	return tableStructs
}