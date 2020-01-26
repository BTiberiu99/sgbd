package db

import "sgbd4/go/utils"

import "fmt"

var (
	db            *Connection
	ActiveIndex   = ""
	runThreadSafe = utils.CreateSyncFunc()
)

func UpdateConnection(conn *Connection) error {

	runThreadSafe(func() {
		db = conn
		ActiveIndex = db.SafeString()
	})
	fmt.Println(db, ActiveIndex)
	if conn.con != nil {
		return nil
	}

	err := db.createConnection()

	return err

}

func DB() *Connection {
	return db
}
