package reverser

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"html/template"
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

func NewReverser(tables *[]string) *Reverser {
	return &Reverser{
		Tables: *tables,
	}
}

func (r *Reverser) Reverse(db *db.DB) error {
	var tableStructs = make([]*model.Table, 0)

	// Получаем структуры таблиц
	for _, table := range r.Tables {
		tableStruct, err := getTableStruct(table, db)
		if err != nil {
			panic(err)
		}

		// достаем внешние ключи
		tableStruct.ForeignKeys = GetForeignKeys(table, db)
		tableStructs = append(tableStructs, tableStruct)
	}

	// Создаем карту отношений таблиц
	relations := DefiningTableRelations(tableStructs)
	fmt.Println(relations)

	tableCollection := make(map[string]*model.Table)
	for _, modelTable := range tableStructs{
		tableCollection[modelTable.Name] = modelTable
	}

	SpecialTypeDefinition(tableCollection, relations)
	// отправляем в шаблон
	sendToTemplate(tableStructs)
	return nil
}

func sendToTemplate(tables []*model.Table) {

	field := model.Field{}
	funcMap :=  template.FuncMap{
		"IsPrimaryKey"      : field.IsPrimaryKey,
		"IsNullableField"	: field.IsNullableField,
		"GetValidate"       : field.GetValidate,
		"GetGraphQlType"	: field.GetGraphQlType,
		"GetForeignType" 	: model.GetForeignType,
	}

	t, err := template.New("tables").Funcs(funcMap).Parse(getTemplateStruct())
	if err != nil { panic(err) }
	for _,table := range tables {
		fileName := "results/" + table.Name + ".graphql"
		fo, err := os.Create(fileName)
		if err != nil {log.Printf("| ERROR | Не удалось открыть файл %s, \n", fileName)}

		defer func() {
			if err := fo.Close(); err != nil {log.Printf("| ERROR | Не удалось закрыть файл %s \n", fileName)}
		}()

		if err := t.Execute(fo, table); err != nil {
			panic(err)
		}
	}

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

func SpecialTypeDefinition(tables map[string]*model.Table, relations map[string]*model.Relation)  {
	// пробегаем по всем табличкам
	for tableName, table := range tables {
		// если есть внешние ключи
		if len(table.ForeignKeys) > 0 {
			relLinkedTo := relations[table.Name].LinkedTo
			// пробегаем по всем таблицам к которым ссылается наша таблица
			for toTable, keyMap := range relLinkedTo {
				// определяем отношения
				for _, field := range table.Fields {
					if field.Name == keyMap["fk"] {
						/* 1. Один к одному:
							Если f_key по нашей табличке(table.Name) уникальное и у таблицы к которой у нас отношение(toTable)
							есть отношение  к нашей таблице и ее поле f_key тоже уникальное
						 */
						if field.IsUnique {
							// получаем таблицу к которой у нас отношение
							if linkTable, ok := tables[toTable]; ok {
								if len(linkTable.ForeignKeys) > 0 {
									if self, ok := relations[toTable]; ok {
										if inversRel, foundInverseRelation := self.LinkedTo[tableName];
										foundInverseRelation {
											for _, inverseTableField := range tables[toTable].Fields {
												if inversRel["fk"] == inverseTableField.Name {
													if inverseTableField.IsUnique {
														// 1 к 1
														field.FkType = toTable
													}
												}
											}
										}
									}
								}
							}
						}
						/** 2. Один ко многим:
							 Если f_key по нашей таблице(table.Name) уникальное, а таблица к которой мы ссылкаемся(toTable)
							не ссылается на нас
						 */

						/* 3. Многие ко многим:
							Если f_key по нашей таблице(table.Name) не уникальное и таблица к которой мы ссылаемся(toTable)
							сама не ссылается никуда, когда на нее ссылаются  2 или более таблиц
						*/


					}
				}
			}
		}
	}
}