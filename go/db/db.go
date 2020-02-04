package db

import (
	"sgbd4/go/utils"
)

var (
	db            *Connection
	ActiveIndex   = ""
	runThreadSafe = utils.CreateSyncFunc()
)

var (
	reuse = map[string]*Connection{}
)

func UpdateConnection(conn *Connection) error {

	runThreadSafe(func() {
		if _, ok := reuse[conn.SafeString()]; !ok {
			reuse[conn.SafeString()] = conn
		} else {
			reuse[conn.SafeString()].ResetTables()
		}

		db = reuse[conn.SafeString()]
		ActiveIndex = db.SafeString()
	})
	// fmt.Println(db, ActiveIndex)
	if conn.con != nil {
		return nil
	}

	err := db.createConnection()

	return err

}

func DB() *Connection {
	//Block when changes are made

	return db
}
