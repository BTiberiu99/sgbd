package expose

import (
	"encoding/json"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/translate"
)

//AddNotNull... Set NOT NULL constraint to a column of a table
func AddNotNull(table, column string) response.Message {

	if db.DB() == nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessageNoConnection),
		}
	}

	col := &db.Column{}

	err := json.Unmarshal([]byte(column), &col)

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessageFailAddNotNULL, col.Name),
		}
	}

	err = col.AddNotNull(table)

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessageFailAddNotNULL, col.Name),
		}
	}

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T(legend.MessageSuccessAddNotNULL, col.Name),
	}
}
