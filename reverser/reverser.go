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

	// смотрим на влешние ключи таблиц
	SendForeignTypes(tableStructs)

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
		fileName := "results/" + table.Name + ".graphQl"
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
	Для полей с внешними ключами проставляем специальный тип
 */
func SendForeignTypes(tables []*model.Table) {

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

						// добавляем для поля дополнительный тип который указывает на другую таблицу
						field.FkType = &fk_field.FkToTable
					}
				}
			}
		}
	}

}

