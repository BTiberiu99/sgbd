package db

import (
	"database/sql"
	"fmt"
	"sgbd4/go/utils"

	_ "github.com/lib/pq"
)

type Connection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	con      *sql.DB
	database *Tables
}

func (c *Connection) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

func (c *Connection) SafeString() string {

	return utils.Sha512EmptyHash(fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Database))
}

func (c *Connection) Conx() *sql.DB {

	if c.con == nil {

		err := c.createConnection()

		if err != nil {
			panic(err)
		}
	}

	return c.con
}

func (c *Connection) CheckConnection() error {

	c.createConnection()

	err := c.con.Ping()

	if err != nil {
		return err
	}

	return err
}

func (c *Connection) createConnection() error {

	con, err := sql.Open("postgres", c.String())

	if err != nil {
		return err
	}

	c.con = con

	return nil
}

func (c *Connection) Tables() *Tables {

	if c.database == nil {
		c.database = new(Tables)
	}

	database := CreateTables()

	c.database = &database

	return c.database
}
