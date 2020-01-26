package db

import (
	"context"
	"fmt"
	"log"
	"sgbd4/go/translate"
	"sgbd4/go/utils"
	"strings"
)

type Column struct {
	Name        string
	Constraints []Constrain
	Type        string
	Position    int
	sync        func(func())
}

func (c *Column) existOrCreateSync() {
	if c.sync == nil {
		c.sync = utils.CreateSyncFunc()
	}
}

//LoadConstrains ... Load all informations about constrains for a column from database
func (c *Column) LoadConstrains(table string) {
	query, _ := translate.QT("constrains", table, c.Name)
	c.loadQuery(query)
}

func (c *Column) LoadCheckConstrains(table string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	query, _ := translate.QT("check_constrains", table, c.Name)

	c.loadQuery(query)
}

func (c *Column) loadQuery(query string) {

	rows, err := db.Conx().QueryContext(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		var (
			name          string
			constrainType string
		)
		if err := rows.Scan(&name, &constrainType); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)

			// c.Constraints = append(c.Constraints, Constrain{})

		}

		c.AddConstrain(name, constrainType)

	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(err)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func (c *Column) AddConstrain(name, constrainType string) {
	c.existOrCreateSync()
	c.sync(func() {
		c.Constraints = append(c.Constraints, Constrain{Name: name, Type: strings.ToUpper(constrainType)})
	})

}

func (c *Column) HasUnique() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsUnique()
	})
}

func (c *Column) HasNotNull() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsNotNull() && !c.IsPrimaryKey()
	})
}

func (c *Column) HasPrimaryKey() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsPrimaryKey()
	})
}

func (c *Column) HasForeignKey() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsForeignKey()
	})
}

func (c *Column) HasCheck() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsCheck()
	})
}

func (c *Column) iterateConstrains(s func(c *Constrain) bool) bool {
	for i := range c.Constraints {
		if s(&c.Constraints[i]) {
			return true
		}
	}
	return false
}
