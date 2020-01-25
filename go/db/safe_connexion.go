package db

import "sgbd4/go/utils"

type SafeConnection struct {
	Name  string
	Index string
}

func NewSafeConnectionFromConnection(conn *Connection) SafeConnection {
	return SafeConnection{
		Name:  conn.Database,
		Index: utils.EncryptString(conn.SafeString()),
	}
}
