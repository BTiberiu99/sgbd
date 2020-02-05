package expose

import (
	"fmt"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/translate"
)

//GetTables ... takes all tables from the database
func GetTables() response.Message {
	msg := response.Message{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				msg = response.Message{
					Type:    legend.TypeError,
					Message: fmt.Sprint(r),
				}
			}
		}()

		if db.DB() == nil {
			msg = response.Message{
				Type:    legend.TypeWarning,
				Message: translate.T(legend.MessageConnectionNotExist),
			}
			return
		}

		msg = response.Message{
			Data: db.DB().Tables(),
		}

	}()

	// fmt.Println("DONE")

	return msg
}

//ResetTables ... resets the cache of the tables
func ResetTables() response.Message {
	msg := response.Message{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				msg = response.Message{
					Type:    legend.TypeError,
					Message: fmt.Sprint(r),
				}
			}
		}()

		if db.DB() == nil {
			msg = response.Message{
				Type:    legend.TypeWarning,
				Message: translate.T(legend.MessageConnectionNotExist),
			}
			return
		}

		db.DB().ResetTables()

	}()

	return msg
}
