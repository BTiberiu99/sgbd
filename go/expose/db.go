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

func CreateConnection(str string) string {
	conn := &db.Connection{}

	json.Unmarshal([]byte(str), &conn)

	err := conn.CheckConnection()

	if err != nil {
		return (&response.Message{
			Type:    legend.TypeError,
			Message: translate.T("fail_connection", err.Error()),
		}).String()
	}

	if _, exist := store.NewStore().Get(conn.SafeString()); exist {
		return (&response.Message{
			Type:    legend.TypeWarning,
			Message: translate.T("exist_connection", conn.Database),
		}).String()
	}

	store.NewStore().Add(*conn)

	db.UpdateConnection(conn)

	return (&response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T("succes_connection"),
		Data:    db.NewSafeConnectionFromConnection(conn),
	}).String()

}

func RemoveConnection(str string) string {

	safeConn := &db.SafeConnection{}

	json.Unmarshal([]byte(str), &safeConn)

	store.NewStore().Remove(utils.DecryptString(safeConn.Index))

	return (&response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T("succes_remove_connection", safeConn.Name),
		Data:    GetConnections(),
	}).String()
}

func GetConnections() string {

	connections := store.NewStore().Connections()

	safeConnections := make([]db.SafeConnection, len(connections))

	j := 0
	for _, con := range connections {
		safeConnections[j] = db.NewSafeConnectionFromConnection(&con)
		j++
	}
	data, err := json.Marshal(safeConnections)

	if err != nil {
		return ""
	}

	return string(data)
}

func SwitchConnection(str string) string {
	safeConn := &db.SafeConnection{}

	json.Unmarshal([]byte(str), &safeConn)

	conn, exist := store.NewStore().Get(utils.DecryptString(safeConn.Index))

	if !exist {
		return (&response.Message{
			Type:    legend.TypeError,
			Message: translate.T("not_exist_connection", safeConn.Name),
		}).String()
	}

	db.UpdateConnection(&conn)

	return (&response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T("succes_connection"),
		Data:    safeConn,
	}).String()
}
