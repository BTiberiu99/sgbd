package expose

import (
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/translate"
)

//GetTables ... takes all tables from the database
func GetTables() response.Message {

	if db.DB() == nil {
		return response.Message{
			Type:    legend.TypeWarning,
			Message: translate.T(legend.MessageConnectionNotExist),
		}
	}

	return response.Message{
		Data: db.DB().Tables(),
	}

}

//ResetTables ... resets the cache of the tables
func ResetTables() response.Message {

	if db.DB() == nil {
		return response.Message{
			Type:    legend.TypeWarning,
			Message: translate.T(legend.MessageConnectionNotExist),
		}
	}

	db.DB().ResetTables()
	return response.Message{}
}
