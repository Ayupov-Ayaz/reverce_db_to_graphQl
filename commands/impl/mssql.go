package impl

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"log"
	"os"
)

type MssqlCommands struct {
	Params db.Params
}

func NewMssqlCommands() *MssqlCommands {
	params := db.Params{
		Version: "10.50.4000.0",
		Level: "SP2",
	}
	return &MssqlCommands{
		Params: params,
	}
}

func (mc *MssqlCommands) GetParams() *db.Params {
	return &mc.Params
}

func (mc *MssqlCommands) GetSupportedVersions() []string {
	return []string{
		"10.50.4000.1",
		//
	}
}
/**
	Получение структуры таблицы в базе данных
 */
func (mc *MssqlCommands) GetTableStruct(tableName string, db *db.DB) (table *model.Table) {

	var fields []*model.Field

	query := `
		declare @table_name varchar(255) = ?
		select distinct
	  		t_c.COLUMN_NAME                             					    as name
			, t_c.DATA_TYPE                                                     as type
			, MAX(isnull(t_c.CHARACTER_MAXIMUM_LENGTH, 0))						as max_length
			, MAX(case when lower(t_c.IS_NULLABLE) = 'yes' then 1 else 0 end)	as is_nullable
  			, MAX(isnull(s.is_unique, 0) )										as is_unique
  			, MAX(isnull(s.is_foreign, 0)) 										as is_foreign
  			, MAX(isnull(s.is_primary, 0)) 										as is_primary
	from INFORMATION_SCHEMA.COLUMNS												as t_c
	left join INFORMATION_SCHEMA.KEY_COLUMN_USAGE								as all_keys
	                 on all_keys.TABLE_NAME = t_c.TABLE_NAME
	                 and t_c.COLUMN_NAME = all_keys.COLUMN_NAME
	left join (
		select CONSTRAINT_NAME as keys
			, case when CONSTRAINT_TYPE = 'UNIQUE'      then 1 else 0 end 		as is_unique
			, case when CONSTRAINT_TYPE = 'FOREIGN KEY' then 1 else 0 end 		as is_foreign
			, case when CONSTRAINT_TYPE = 'PRIMARY KEY' then 1 else 0 end 		as is_primary
		from INFORMATION_SCHEMA.TABLE_CONSTRAINTS as  f_key 
		where TABLE_NAME = @table_name
	) as s ON s.keys = all_keys.CONSTRAINT_NAME
	where t_c.TABLE_NAME = @table_name 
	group by  t_c.COLUMN_NAME, t_c.DATA_TYPE`

	if  err := db.Select(&fields, query, tableName); err != nil {
		log.Printf("| DB.ERROR | Ошибка при получении структуры таблицы = %s :\n %s \n", tableName, err.Error())
		log.Printf("| NOTICE | Таблица %s пропущена", tableName)
		return nil
	}
	t := &model.Table{
		Name: tableName,
		Fields: fields,
	}
	return t
}

/**
	Получение внешних ключей таблицы
 */
func (mc *MssqlCommands) GetForeignKeys(tableName string, db *db.DB) []*model.ForeignKey {

	var foreignKeys []*model.ForeignKey

	query := `
	select
    	col.name			as field_name
    	,tab_prim.name		as fk_in_table
    	,col_prim.name		as refers_to_field
	
	from sys.tables													as tab
    	inner join sys.foreign_keys									as fk
    		on tab.object_id = fk.parent_object_id
    inner join sys.foreign_key_columns								as fkc
    		on fk.object_id = fkc.constraint_object_id
    	inner join sys.columns										as col
    		on fkc.parent_object_id = col.object_id
    		and fkc.parent_column_id = col.column_id
    	inner join sys.columns										as col_prim
    		on fkc.referenced_object_id = col_prim.object_id
    		and fkc.referenced_column_id = col_prim.column_id
    	inner join sys.tables										as tab_prim
    		on fk.referenced_object_id = tab_prim.object_id
	where tab.name= ? `

	if err := db.Select(&foreignKeys, query, tableName); err != nil {
		log.Printf("| DB.ERROR | Ошибка при получение внешних ключей таблицы %s: \n %s \n",
			tableName, err.Error())
		log.Printf("| NOTICE | Таблица %s пропущена", tableName)
	}
	return foreignKeys
}

/**
	Получает названия всех таблиц в базе данных
 */
func (mc *MssqlCommands) GetAllTableNames(db *db.DB) []string {
	query := `
		SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE'
		`
	tables, err :=  mc.getTables(query, db, nil)
	if err != nil {
		log.Printf("| DB.ERROR |При попытке получить названия всех таблиц в бд произошла ошибка:\n %s \n",
			err.Error())
		os.Exit(-1)
	}
	return tables
}

func (mc *MssqlCommands) GetTableByLike(tables []string, db *db.DB) []string {
	query := `
	select table_name from INFORMATION_SCHEMA.TABLES where TABLE_TYPE='BASE TABLE' and TABLE_NAME like ?
	`
	var allTables = make([]string, 0)
	for i := 0; i < len(tables); i++ {
		result, err := mc.getTables(query, db, tables[i])
		if err != nil {
			log.Printf("При попытке получить имена таблиц с помощью ключа \"-l\" произошла ошибка:\n%s",
				err.Error())
			os.Exit(-1)
		}
		allTables = append(allTables, result...)
	}
	return allTables
}

/**
	Получение таблиц
 */
func (mc *MssqlCommands) getTables(query string, db *db.DB, params interface{}) (tables []string, err error) {
	var res []struct {
		Name string `db:"table_name"`
	}

	if params == nil {
		err = db.Select(&res, query)
	} else {
		err = db.Select(&res, query, params)
	}

	if  err != nil {
		return nil, err
	}

	// TODO: pull workers
	for _, table := range res {
		tables = append(tables, table.Name)
	}
	return
}