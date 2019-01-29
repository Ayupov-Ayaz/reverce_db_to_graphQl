package reverser

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/commands"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"log"
	"os"
)
/**
  Structure for reversing DataBase structure to GraphQl structure
 */
type Reverser struct {
	OutputFileName string
	Tables []string
}

func NewReverser(tables []string) *Reverser {
	return &Reverser{
		Tables: tables,
	}
}

/**
	Функция которая запускает процессы:
	1) сборку данных из бд по таблицам
	2) получение из бд внешних ключей таблиц
	3) создание карты отношений между таблицами
	4) выгрузка в шаблон результатов
 */
func (r *Reverser) Reverse(db *db.DB, com commands.DbCommander, flags map[string]bool) {
	var tableCollection = make(map[string]*model.Table, 0)
	var tableRelations = make(map[string]*model.Relation, 0)

	tableCollection, tableRelations = r.getTableData(tableCollection, tableRelations, com, db, flags)

	SpecialTypeDefinition(tableCollection, tableRelations)
	// отправляем в шаблон
	sendToTemplate(tableCollection)
}

func (r *Reverser) getTableData(tCollection map[string]*model.Table, tRelations map[string]*model.Relation,
	com commands.DbCommander, db *db.DB, flags map[string]bool) (
	tableCollection map[string]*model.Table, tableRelation map[string]*model.Relation){

	tableStructs := r.getTablestructs(db, com, flags)

	if len(r.Tables) != len(tableStructs) {
		deleteNotFoundTables(r.Tables, tableStructs)
	}

	if len(tableStructs) == 0 && len(tCollection) == 0 { // сработает при первом проходе
		log.Println("| ERROR | Все таблицы которые вы указали, не были найдены в бд. Проверьте названия таблиц.")
		os.Exit(-1)
	}

	// Создаем карту отношений таблиц
	tableRelation = DefiningTableRelations(tableStructs)
	// Делаем из среза карту
	tableCollection = makeTableCollection(tableStructs)
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
			for _, fk_field := range table.ForeignKeys {
				// ищем поле в котором есть этот внешний ключ
				for _, field := range table.Fields {
					if field.Name == fk_field.FieldName {

						// relation
						relation, ok := tableRelation[table.Name]
						if !ok{
							relation = &model.Relation{}
							tableRelation[table.Name] = relation
						}
						relation.AddLinkedTo(fk_field.FkToTable, fk_field.FieldName, fk_field.PkField)

						// inverseRelation
						inverseRelation, ok := tableRelation[fk_field.FkToTable]
						if !ok{
							inverseRelation = &model.Relation{}
							tableRelation[fk_field.FkToTable] = inverseRelation
						}
						inverseRelation.AddLinkToMe(table.Name, fk_field.FieldName, fk_field.PkField)
						
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
			// получаем все отношения для нашей таблицы
			relLinkedTo := relations[table.Name].LinkedTo
			// пробегаем по отношениям таблиц к которым ссылается наша таблица
			for toTable, keyMap := range relLinkedTo {
				// пробегаем по всем полям которые ссылаются на другие таблицы
				for _, field := range table.ForeignFields {
					if field.Name == keyMap["fk"] {
						/* 1. Один к одному:
							Если f_key по нашей табличке(table.Name) уникальное и у таблицы к которой у нас
							отношение(toTable) есть обратное отношение к нашей таблице, и ее поле f_key тоже уникальное

						  2. Один ко многим:
							Если f_key по нашей таблице(table.Name) уникальное, а таблица к которой мы
							ссылкаемся(toTable) не ссылается на нас
						*/

						// получаем таблицу к которой у нас отношение
						if linkTable, ok := tables[toTable]; ok {
							if len(linkTable.ForeignKeys) > 0 {
								if self, ok := relations[toTable]; ok {
									// ищем обратные отношения
									if inverseRel, foundInverseRelation := self.LinkedTo[tableName];
									foundInverseRelation {
										for _, inverseTableField := range tables[toTable].ForeignFields {
											if inverseRel["fk"] == inverseTableField.Name {
												if field.IsUnique && inverseTableField.IsUnique {
													field.FkType = toTable  // OneToOne
												} else if field.IsUnique && !inverseTableField.IsUnique {
													field.FkType = "[" + toTable + getNullSign(tables,
														tableName, field.Name) + "]" // OneToMany
												} else if !field.IsUnique && inverseTableField.IsUnique {
													field.FkType = toTable // ManyToOne
												} else {
													field.FkType = "[" + toTable + getNullSign(tables,
														tableName, field.Name) + "]"
												}
											} else {
												err := "| NOTICE | Для строки %s.%s не найдены внешние ключи." +
													"Постройте отношение в ручную для поля %s.%s\n"
												fmt.Printf(err , toTable, inverseTableField.Name, table.Name, field.Name)
											}
										}
									} else {
										// Если нету обратных отношений (например, главная таблица,
										// к которой все ссылаются)
										if field.IsUnique {
											field.FkType = toTable // OneToMany
										} else {
											/* 3. Многие ко многим:
											Если f_key по нашей таблице(table.Name) не уникальное и таблица к которой
											мы ссылаемся(toTable) сама не ссылается никуда
										*/
											field.FkType = "[" + toTable + getNullSign(tables,
												tableName, field.Name) + "]" // ManyToMany
										}
									}
								} else {
									log.Printf("| ERROR | По таблице %s, не найдены отношения \n", toTable)
								}
							} else {
								// Отношений у таблицы нету.
								if field.IsUnique {
									field.FkType = toTable
								} else {
									field.FkType = "[" + toTable +   getNullSign(tables, tableName, field.Name) + "]"
								}
							}
						} else {
							relationError := "| NOTICE | Не удалось определить отношение %s.%s => %s " +
								"(для поля %s прописан тип NO_TABLE_SPECIFIED)"
							tableNotFound := relationError + ", Таблица '%s' не была указана \n"
							relationError2 := relationError + ", Причина не ясна =( "
							if _, tableExist := tables[toTable]; !tableExist {
								log.Printf(tableNotFound, tableName, field.Name, toTable, field.Name, toTable)
							} else {
								log.Printf(relationError2, tableName, field.Name, toTable, field.Name)
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

// Функция которая проверяет значение isNullable по полю внешнего ключа таблицы
func getNullSign(tables map[string]*model.Table ,tableName, fkName string) string {
	if table, ok := tables[tableName]; ok {
		for _, field := range table.ForeignFields {
			if !field.IsNullable {
				return "!"
			} else {
				return ""
			}
		}
	}
	fmt.Printf("| NOTICE | Не задана таблица %s, не удалось проверить на \"nullable\" поле %s ", tableName, fkName)
	return "NOT_FOUND_TABLE"
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

func (r *Reverser) getTablestructs(db *db.DB, com commands.DbCommander, flags map[string]bool) []*model.Table {
	var tableStructs = make([]*model.Table, 0)

	// Получаем структуры таблиц
	for _, table := range r.Tables {
		tableStruct := com.GetTableStruct(table, db)
		if tableStruct != nil && len(tableStruct.Fields) > 0 {
			// достаем внешние ключи
			tableStruct.ForeignKeys = com.GetForeignKeys(table, db)
			tableStructs = append(tableStructs, tableStruct)
		}
	}
	return tableStructs
}