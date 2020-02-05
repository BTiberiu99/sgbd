package expose

import (
	"encoding/json"
	"fmt"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/store"
	"sgbd4/go/translate"
	"sgbd4/go/utils"
)

var (
	runSafe = utils.CreateSyncFunc()
)

//CreateConnection ... creates new connection with a database
func CreateConnection(str string) response.Message {
	var msg response.Message

	runSafe(func() {
		defer func() {
			if r := recover(); r != nil {
				msg = response.Message{
					Type:    legend.TypeError,
					Message: fmt.Sprint(r),
				}
			}
		}()

		conn := &db.Connection{}

		json.Unmarshal([]byte(str), &conn)

		err := conn.CheckConnection()

		if err != nil {
			msg = response.Message{
				Type:    legend.TypeError,
				Message: translate.T(legend.MessageConnectionFail, err.Error()),
			}

			return
		}

		if _, exist := store.GetInstance().Get(conn.SafeString()); exist {
			msg = response.Message{
				Type:    legend.TypeWarning,
				Message: translate.T(legend.MessageConnectionExist, conn.Database),
			}
			return
		}

		store.GetInstance().Add(*conn)

		db.UpdateConnection(conn)

		msg = response.Message{
			Type:    legend.TypeSucces,
			Message: translate.T(legend.MessageConnectionSuccess),
			Data:    db.NewSafeConnectionFromConnection(conn),
		}
	})

	return msg

}

//RemoveConnection ... removes coonection from the store
func RemoveConnection(str string) response.Message {

	var msg response.Message

	runSafe(func() {
		defer func() {
			if r := recover(); r != nil {
				msg = response.Message{
					Type:    legend.TypeError,
					Message: fmt.Sprint(r),
				}
			}
		}()
		safeConn := &db.SafeConnection{}

		json.Unmarshal([]byte(str), &safeConn)
		_, is := store.GetInstance().Get(safeConn.Index)

		if is {
			store.GetInstance().Remove(safeConn.Index)
		}

		msg = response.Message{
			Type:    legend.TypeSucces,
			Message: translate.T(legend.MessageConnectionSuccessRemove, safeConn.Name),
			Data:    GetConnections().Data,
		}
	})

	return msg
}

//GetConnections ... get all the conections and the curent active connection
func GetConnections() response.Message {

	connections := store.GetInstance().Connections()

	safeConnections := make([]db.SafeConnection, len(connections))

	j := 0
	for _, con := range connections {
		safeConnections[j] = db.NewSafeConnectionFromConnection(&con)
		j++
	}

	return response.Message{
		Data: map[string]interface{}{
			"Index":       db.ActiveIndex,
			"Connections": safeConnections,
		},
	}
}

//SwitchConnection ... switches between the connections from the store
func SwitchConnection(str string) response.Message {

	var msg response.Message

	safeConn := &db.SafeConnection{}

	json.Unmarshal([]byte(str), &safeConn)

	runSafe(func() {

		defer func() {
			if r := recover(); r != nil {
				msg = response.Message{
					Type:    legend.TypeError,
					Message: fmt.Sprint(r),
				}
			}
		}()
		conn, exist := store.GetInstance().Get(safeConn.Index)

		if !exist {
			msg = response.Message{
				Type:    legend.TypeError,
				Message: translate.T(legend.MessageConnectionNotExist, safeConn.Name),
			}
			return
		}

		if db.DB().SafeString() != conn.SafeString() {

			err := db.UpdateConnection(&conn)

			store.GetInstance().Save()

			if err != nil {
				msg = response.Message{
					Type:    legend.TypeError,
					Message: translate.T(legend.MessageConnectionFail, err.Error()),
				}

				return
			}
		} else {
			db.DB().ResetTables()
		}

		msg = response.Message{
			Type:    legend.TypeSucces,
			Message: translate.T(legend.MessageConnectionSuccess),
			Data:    safeConn,
		}
	})

	return msg

}
