package impl

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"log"
)

type MssqlCommands struct {}

func NewMssqlCommands() *MssqlCommands {
	return &MssqlCommands{}
}
/**
	Получение структуры таблицы в базе данных
 */
func (mc *MssqlCommands) GetTableStruct(tableName string, db *db.DB) (table *model.Table) {

	fields := []*model.Field{}

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

	foreignKeyses := []*model.ForeignKey{}

	query := `
	select
    	col.name			as field_name
    	,tab_prim.name		as fk_to_table
    	,col_prim.name		as pk_field
	
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

	if err := db.Select(&foreignKeyses, query, tableName); err != nil {
		log.Printf("| DB.ERROR | Ошибка при получение внешних ключей таблицы %s: \n %s \n",
			tableName, err.Error())
		log.Printf("| NOTICE | Таблица %s пропущена", tableName)
	}
	return foreignKeyses
}
