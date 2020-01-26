package expose

import (
	"sgbd4/go/db"
)

func GetTables() *db.Tables {
	// fmt.Println(db.DB())

	// fmt.Println(db.DB().Tables())
	return db.DB().Tables()

}
