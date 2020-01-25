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

func Run(sql string) string {

	r := new(run)
	json.Unmarshal([]byte(sql), &r)
	if db.DB().Conx() == nil {
		return (&response.Message{
			Type:    legend.TypeError,
			Message: translate.T("no_connection"),
		}).String()
	}

	var (
		result interface{}
		err    error
	)
	r.Run = strings.TrimSpace(strings.ToUpper(r.Run))
	if strings.HasPrefix(r.Run, legend.SELECT) {

		result, err = db.DB().Conx().Query(r.Run)

		if err == nil {
			return (&response.Message{
				Data: result,
			}).String()
		}

	} else {

		z, e := db.DB().Conx().Exec(r.Run)
		nr, _ := z.RowsAffected()

		if err == nil {
			return (&response.Message{
				Type:    legend.TypeSucces,
				Message: translate.T("rows_affected", nr),
				Data:    result,
			}).String()
		}

		result, err = z, e
	}

	return (&response.Message{
		Type:    legend.TypeError,
		Message: err.Error(),
	}).String()

}
