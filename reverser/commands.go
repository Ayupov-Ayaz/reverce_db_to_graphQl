package reverser

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/model"
)

func getTableStruct(tableName string, db *db.DB)  (table *model.Table, err error) {

	fields := &[]model.Field{}

	t := &model.Table{
		Name: tableName,
		Fields: fields,
	}

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

	if  err := db.Select(fields, query, tableName); err != nil {
		return nil, err
	}

	return t, nil
}
