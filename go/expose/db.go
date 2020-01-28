package expose

import (
	"encoding/json"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/store"
	"sgbd4/go/translate"
	"sgbd4/go/utils"
)

func CreateConnection(str string) response.Message {
	conn := &db.Connection{}

	json.Unmarshal([]byte(str), &conn)

	err := conn.CheckConnection()

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T("fail_connection", err.Error()),
		}
	}

	if _, exist := store.NewStore().Get(conn.SafeString()); exist {
		return response.Message{
			Type:    legend.TypeWarning,
			Message: translate.T("exist_connection", conn.Database),
		}
	}

	store.NewStore().Add(*conn)

	db.UpdateConnection(conn)

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T("succes_connection"),
		Data:    db.NewSafeConnectionFromConnection(conn),
	}

}

func RemoveConnection(str string) response.Message {

	safeConn := &db.SafeConnection{}

	json.Unmarshal([]byte(str), &safeConn)

	store.NewStore().Remove(utils.DecryptString(safeConn.Index))

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T("succes_remove_connection", safeConn.Name),
		Data:    GetConnections(),
	}
}

func GetConnections() []db.SafeConnection {

	connections := store.NewStore().Connections()

	safeConnections := make([]db.SafeConnection, len(connections))

	j := 0
	for _, con := range connections {
		safeConnections[j] = db.NewSafeConnectionFromConnection(&con)
		j++
	}

	return safeConnections
}

func SwitchConnection(str string) response.Message {
	safeConn := &db.SafeConnection{}

	json.Unmarshal([]byte(str), &safeConn)

	conn, exist := store.NewStore().Get(safeConn.Index)

	if !exist {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T("not_exist_connection", safeConn.Name),
		}
	}

	err := db.UpdateConnection(&conn)

	store.NewStore().Save()

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T("fail_connection", err.Error()),
		}
	}

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T("succes_connection"),
		Data:    safeConn,
	}
}

func AddNotNull(table, column string) response.Message {
	col := &db.Column{}

	err := json.Unmarshal([]byte(column), &col)

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T("fail_add_not_null", col.Name),
		}
	}

	err = col.AddNotNull(table)

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T("fail_add_not_null", col.Name),
		}
	}

	col.Constraints = []db.Constraint{}
	col.Load(table)

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T("succes_add_not_null", column),
		Data:    col,
	}
}
