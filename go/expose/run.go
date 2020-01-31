package expose

import (
	"encoding/json"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/translate"
	"strings"
)

type run struct {
	Run string `json:"run"`
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

	var (
		result interface{}
		err    error
	)

	//Normalize query
	r.Run = strings.TrimSpace(strings.ToUpper(r.Run))

	//Check select
	if strings.HasPrefix(r.Run, legend.SELECT) {

		result, err = db.DB().Conx().Query(r.Run)

		if err == nil {
			return response.Message{
				Data: result,
			}
		}

		//Check others types
	} else {

		z, e := db.DB().Conx().Exec(r.Run)
		nr, _ := z.RowsAffected()

		if e == nil {
			return response.Message{
				Type:    legend.TypeSucces,
				Message: translate.T(legend.MessageRowsAffected, nr),
				Data:    result,
			}
		}

		result, err = z, e
	}

	return response.Message{
		Type:    legend.TypeError,
		Message: err.Error(),
	}

}
