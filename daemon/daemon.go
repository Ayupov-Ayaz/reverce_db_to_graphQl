package daemon

import (
	"fmt"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/commands"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/errors"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/reverser"
)

func Run(tables []string, flags map[string]bool) {
	cfg := getConfigs()
	dbCon, err := db.InitDB(cfg.Db)
	if err != nil {
		errors.PrintFatalError(fmt.Sprintf("Ошибка при подключении к БД.(Проверьте данные подключения): \n %s",
			err.Error()), true)
	}
	// scanning
	rev := reverser.NewReverser(tables)
	com := commands.GetDbCommander(cfg.Db)

	if flags["*"] { // Если указан флаг "*" получаем все наши таблицы
		rev.TablesForSearch = com.GetAllTableNames(dbCon)
	} else if flags["l"] { // Если указан флаг "l" получаем таблицы через оператор like
		rev.TablesForSearch = com.GetTableByLike(tables, dbCon)
		if !flags["d"] {
			var needTables = make([]string, 0)
			// нужно будет вывести только те таблицы которые были запрошены, по этому записываем их названия в срез
			for _, table := range rev.TablesForSearch {
				needTables = append(needTables, table)
			}
			rev.TablesForShow = needTables
		}
	}

	if dbCon.CompareDbParams(com, flags) {
		rev.Reverse(dbCon, com, flags)
	}
}