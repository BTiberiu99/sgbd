package db

import (
	"context"
	"fmt"
	"log"
	"sgbd4/go/legend"
	"sgbd4/go/translate"
	"sgbd4/go/utils"
	"strings"
	"sync"
)

type Column struct {
	Name        string
	Constraints []*Constraint
	Type        string
	Position    int
	sync        func(func())
	WithoutNULL bool
}

func (c *Column) existOrCreateSync() {
	if c.sync == nil {
		c.sync = utils.CreateSyncFunc()
	}
}

//LoadConstrains ... Load all informations about constrains for a column from database
func (c *Column) Load(table string) {

	group := &sync.WaitGroup{}

	group.Add(1)
	group.Add(1)
	group.Add(1)

	go c.loadConstrains(table, group)
	go c.loadCheckConstrains(table, group)
	go c.checkWithoutNull(table, group)

	group.Wait()

}

//LoadConstrains ... Load all informations about constrains for a column from database
func (c *Column) loadConstrains(table string, group *sync.WaitGroup) {

	query, _ := translate.QT(legend.QueryCONSTRAINTS, table, c.Name)

	c.loadQuery(query)

	group.Done()
}

func (c *Column) loadCheckConstrains(table string, group *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	query, _ := translate.QT(legend.QueryCHECKCONSTRAINTS, table, c.Name)

	c.loadQuery(query)

	group.Done()
}

func (c *Column) checkWithoutNull(table string, group *sync.WaitGroup) {
	var count int
	query, _ := translate.QT(legend.QueryCOUNTNOTNULL, table, c.Name)

	row := db.Conx().QueryRowContext(context.Background(), query)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	c.WithoutNULL = count == 0

	group.Done()
}

func (c *Column) loadQuery(query string) {

	rows, err := db.Conx().QueryContext(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	for rows.Next() {

		values := make([]string, len(cols))
		pointers := make([]interface{}, len(cols))

		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)

			// c.Constraints = append(c.Constraints, Constraint{})

		}

		c.AddConstrain(values...)

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

func (c *Column) AddNotNull(table string) error {

	query, _ := translate.QT(legend.QuerySETNOTNULL, table, c.Name)

	_, err := db.Conx().ExecContext(context.Background(), query)

	return err
}

func (c *Column) AddConstrain(items ...string) {
	c.existOrCreateSync()

	c.sync(func() {
		constr := &Constraint{}
		if len(items) > 0 {
			constr.Name = items[0]
		}
		if len(items) > 1 {
			constr.Type = strings.ToUpper(items[1])
		}

		if constr.IsForeignKey() {
			if len(items) > 2 {
				constr.ForeingTableName = items[2]
			}

			if len(items) > 3 {
				constr.ForeingColumnName = items[3]
			}

			if len(items) > 4 {
				constr.UpdateRule = items[4]
			}
			if len(items) > 5 {
				constr.DeleteRule = items[5]
			}
		}

		c.Constraints = append(c.Constraints, constr)
	})

}

func (c *Column) HasUnique() bool {
	return c.iterateConstrains(func(c *Constraint) bool {
		return c.IsUnique()
	})
}

func (c *Column) HasNotNull() bool {
	return c.iterateConstrains(func(c *Constraint) bool {
		return c.IsNotNull() && !c.IsPrimaryKey()
	})
}

func (c *Column) HasPrimaryKey() bool {
	return c.iterateConstrains(func(c *Constraint) bool {
		return c.IsPrimaryKey()
	})
}

func (c *Column) HasForeignKey() bool {
	return c.iterateConstrains(func(c *Constraint) bool {
		return c.IsForeignKey()
	})
}

func (c *Column) HasCheck() bool {
	return c.iterateConstrains(func(c *Constraint) bool {
		return c.IsCheck()
	})
}

func (c *Column) iterateConstrains(s func(c *Constraint) bool) bool {
	for i := range c.Constraints {
		if s(c.Constraints[i]) {
			return true
		}
	}
	return false
}
