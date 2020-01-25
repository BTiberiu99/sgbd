package db

var (
	db = &Connection{}
)

func UpdateConnection(conn *Connection) {
	db = conn
}

func DB() *Connection {
	return db
}
