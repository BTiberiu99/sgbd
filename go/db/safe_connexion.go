package db

//SafeConnection ... models a safe connection that can be sent to the frontend
type SafeConnection struct {
	Name  string
	Index string
}

func NewSafeConnectionFromConnection(conn *Connection) SafeConnection {
	return SafeConnection{
		Name:  conn.Database,
		Index: conn.SafeString(),
	}
}
