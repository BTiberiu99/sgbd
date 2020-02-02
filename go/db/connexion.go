package db

import (
	"database/sql"
	"fmt"
	"sgbd4/go/utils"

	_ "github.com/lib/pq"
)

//Connection ... models the connection with a database
type Connection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	con      *sql.DB
	database *Tables
}

//string... returns the sql connection string
func (c *Connection) string() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

//SafeString ... returns an unique string of the connexion
func (c *Connection) SafeString() string {

	return utils.Sha512EmptyHash(fmt.Sprintf("%s_%d_%s_%s",
		c.Host, c.Port, c.User, c.Database))
}

//Conx ... returns the underline sql db connection on Connexion struct
func (c *Connection) Conx() *sql.DB {

	if c.con == nil {

		err := c.createConnection()

		if err != nil {
			panic(err)
		}
	}

	return c.con
}

//CheckConnection ... check if connection is still alive
func (c *Connection) CheckConnection() error {

	c.createConnection()

	err := c.con.Ping()

	if err != nil {
		return err
	}

	return err
}

func (c *Connection) createConnection() error {

	con, err := sql.Open("postgres", c.string())
	con.SetMaxOpenConns(2)
	con.SetMaxIdleConns(1)

	if err != nil {
		return err
	}

	c.con = con

	return nil
}

//Tables ... Take all tables from database or cache
func (c *Connection) Tables() *Tables {

	if c.database != nil {
		return c.database
	}

	database := CreateTables()

	c.database = &database

	return c.database
}

//ResetTables ... Reset cached tables
func (c *Connection) ResetTables() *Connection {
	c.database = nil
	return c
}
