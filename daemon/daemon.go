package daemon

import (
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/commands"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/db"
	"github.com/Ayupov-Ayaz/reverse_db_to_graphql/reverser"
	"log"
	"os"
)

func Run(tables []string, flags map[string]bool) {
	cfg := getConfigs()
	db, err := db.InitDB(cfg.Db)
	if err != nil {
		log.Printf("| SYS.ERROR | Ошибка при подключении к БД.(Проверьте данные подключения): \n %s",
			err.Error())
		os.Exit(-1)
	}
	// scanning
	rev := reverser.NewReverser(tables)
	com := commands.GetDbCommander(cfg.Db)

	if flags["*"] { // Если указан флаг "*" получаем все наши таблицы
		rev.Tables = com.GetAllTableNames(db)
	} else if flags["l"] { // Если указан флаг "l" получаем таблицы через оператор like
		rev.Tables = com.GetTableByLike(tables, db)
	}

	if db.CompareDbParams(com, flags) {
		rev.Reverse(db, com, flags)
	}
}