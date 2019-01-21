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
			 select
			 ISC.COLUMN_NAME                               as name
			 ,ISC.DATA_TYPE   							   as type
			 ,isnull(ISC.CHARACTER_MAXIMUM_LENGTH, 0)      as max_length
			 ,case when lower(ISC.IS_NULLABLE) = 'yes' 
					then  1 else 0 end 	       			   as is_nullable
			 ,case when 
				 ISC.COLUMN_NAME = ISKCU.COLUMN_NAME 
				 then 1 else 0 end 					       as primary_key

			 from INFORMATION_SCHEMA.COLUMNS ISC
			 left join INFORMATION_SCHEMA.KEY_COLUMN_USAGE ISKCU on
			 ISKCU.TABLE_NAME =  ISC.TABLE_NAME
			 	   and  ISKCU.TABLE_CATALOG = ISC.TABLE_CATALOG
			 	   and  ISKCU.TABLE_SCHEMA = ISC.TABLE_SCHEMA
			 
			 where ISC.TABLE_NAME = 'tc_trip' `

	if  err := db.Select(fields, query); err != nil {
		return nil, err
	}

	return t, nil
}
