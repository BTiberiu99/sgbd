package db

import (
	"context"
	"log"
	"sgbd4/go/legend"
	"sgbd4/go/translate"
	"sync"
	"time"
)

//Table ... one ore many Table of a database
type Tables []*Table

//Load all informations about tables from database
func CreateTables() Tables {
	query, _ := translate.QT(legend.QueryTABLES)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.Conx().QueryContext(ctx, query)

	if err != nil {
		log.Fatal(err)
	}

	t := make([]*Table, 0)

	defer rows.Close()

	group := &sync.WaitGroup{}
	i := 0
	for rows.Next() {

		var name string
		if err := rows.Scan(&name); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)

		}

		t = append(t, &Table{
			Name:    name,
			Columns: []*Column{},
		})

		group.Add(1)

		go func(t *Table) {

			t.LoadTable()

			group.Done()

		}(t[i])

		i++

	}

	group.Wait()

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

	return Tables(t)
}

func (t Tables) Iterate(call func(*Table, *Column, *Constraint) error) error {

	var err error
	for _, table := range t {
		for _, column := range table.Columns {
			for _, constraint := range column.Constraints {
				err = call(table, column, constraint)

				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

func (t Tables) FindColumn(tableName, columnName string) *Column {
	for _, table := range t {
		for _, column := range table.Columns {
			if table.Name == tableName && column.Name == columnName {
				return column
			}
		}
	}

	return nil

}
