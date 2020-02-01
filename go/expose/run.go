package expose

import (
	"context"
	"encoding/json"
	"fmt"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/translate"
	"strconv"
	"strings"
)

type run struct {
	Run string `json:"run"`
}

type ReturnRows struct {
	Rows    interface{}
	Columns []string
}

func Run(sql string) response.Message {

	r := new(run)
	json.Unmarshal([]byte(sql), &r)
	if db.DB().Conx() == nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessageNoConnection),
		}
	}

	//Normalize query
	r.Run = strings.TrimSpace(strings.ToUpper(r.Run))

	//Check select
	if strings.HasPrefix(r.Run, legend.SELECT) {

		rows, err := db.DB().Conx().QueryContext(context.Background(), r.Run)

		if err != nil {
			return response.Message{
				Type:    legend.TypeError,
				Message: err.Error(),
			}
		}

		defer rows.Close()

		cols, _ := rows.Columns()
		data := make([][]interface{}, 0)
		t, _ := rows.ColumnTypes()
		for i := range t {
			fmt.Println(*t[i])
		}

		for rows.Next() {

			values := make([]interface{}, len(cols))
			pointers := make([]interface{}, len(cols))

			for i := range values {
				pointers[i] = &values[i]
			}

			if err := rows.Scan(pointers...); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				return response.Message{
					Type:    legend.TypeError,
					Message: err.Error(),
				}

				// c.Constraints = append(c.Constraints, Constraint{})

			}
			for key, val := range values {

				switch val.(type) {

				case []uint8:
					f, _ := strconv.ParseFloat(string(val.([]byte)), 64)
					values[key] = f
				}

			}

			data = append(data, values)

		}

		// If the database is being written to ensure to check for Close
		// errors that may be returned from the driver. The query may
		// encounter an auto-commit error and be forced to rollback changes.
		err = rows.Close()

		if err != nil {
			return response.Message{
				Type:    legend.TypeError,
				Message: err.Error(),
			}
		}

		err = rows.Err()
		if err != nil {
			return response.Message{
				Type:    legend.TypeError,
				Message: err.Error(),
			}
		}

		return response.Message{
			Data: ReturnRows{
				Rows:    data,
				Columns: cols,
			},
		}
		//Check others types
	} else {

		z, err := db.DB().Conx().Exec(r.Run)

		if err != nil {
			return response.Message{
				Type:    legend.TypeError,
				Message: err.Error(),
			}
		}
		nr, _ := z.RowsAffected()

		return response.Message{
			Type:    legend.TypeSucces,
			Message: translate.T(legend.MessageRowsAffected, nr),
			Data:    z,
		}

	}

}
