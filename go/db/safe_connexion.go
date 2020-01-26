package db

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
