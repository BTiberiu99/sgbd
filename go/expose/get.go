package expose

import (
	"sgbd4/go/db"
)

func GetTables() *db.Tables {
	// fmt.Println(db.DB())

	return db.DB().Tables()

}
