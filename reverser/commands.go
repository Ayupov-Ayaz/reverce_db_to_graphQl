package reverser

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
	"log"
)

/**
	Получение структуры таблицы в базе данных
 */
func getTableStruct(tableName string, db *db.DB)  (table *model.Table, err error) {

	fields := []*model.Field{}

	query := `
		select distinct
     			   t_c.COLUMN_NAME                             					      as name
     			  , t_c.DATA_TYPE                                                     as type
     			  , isnull(t_c.CHARACTER_MAXIMUM_LENGTH, 0)					          as max_length
			  
     			  , case when lower(t_c.IS_NULLABLE) = 'yes' then  1 else 0 end       as is_nullable
			  
     			  , case when p_keys.CONSTRAINT_TYPE is not null then 1 else 0 end    as is_primary
			  
     			  , case when f_keys.CONSTRAINT_TYPE is not null then 1 else 0 end    as is_foreign

        from INFORMATION_SCHEMA.COLUMNS t_c
	   
        left join INFORMATION_SCHEMA.KEY_COLUMN_USAGE all_keys on all_keys.TABLE_NAME =  t_c.TABLE_NAME
                                                              and t_c.COLUMN_NAME = all_keys.COLUMN_NAME

        left join INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS p_keys on p_keys.TABLE_NAME = all_keys.TABLE_NAME
                                                                and p_keys.CONSTRAINT_NAME = all_keys.CONSTRAINT_NAME
                                                                and p_keys.CONSTRAINT_TYPE = 'PRIMARY KEY'

        left join INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS f_keys on f_keys.TABLE_NAME = all_keys.TABLE_NAME
                                                                and f_keys.CONSTRAINT_NAME = all_keys.CONSTRAINT_NAME
                                                                and f_keys.CONSTRAINT_TYPE = 'FOREIGN KEY'

	    where t_c.TABLE_NAME = ? `

	if  err := db.Select(&fields, query, tableName); err != nil {
		return nil, err
	}

	t := &model.Table{
		Name: tableName,
		Fields: fields,
	}

	return t, nil
}

/**
	Получение внешних ключей таблицы
 */
func GetForeignKeys(tableName string, db *db.DB) []*model.ForeignKey {

	foreignKeyses := []*model.ForeignKey{}

	query := `
	select
       col.name 			as field_name
       ,tab_prim.name 		as fk_to_table
       ,col_prim.name 		as pk_field
	
	from sys.tables as tab
       inner join sys.foreign_keys as fk
                  on tab.object_id = fk.parent_object_id
       inner join sys.foreign_key_columns as fkc
                  on fk.object_id = fkc.constraint_object_id
       inner join sys.columns as col
                  on fkc.parent_object_id = col.object_id
                  and fkc.parent_column_id = col.column_id
       inner join sys.columns as col_prim
                  on fkc.referenced_object_id = col_prim.object_id
                  and fkc.referenced_column_id = col_prim.column_id
       inner join sys.tables as tab_prim
                  on fk.referenced_object_id = tab_prim.object_id
	where tab.name= ?
	`
	if err := db.Select(&foreignKeyses, query, tableName); err != nil {
		log.Printf("| ERROR | GetForeignKeys( %s ) : \n %s \n", tableName, err.Error())
		panic(err)
	}

	return foreignKeyses
}
