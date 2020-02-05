package db

import (
	"context"
	"log"
	"sgbd4/go/legend"
	"sgbd4/go/translate"
	"sort"
	"strings"
	"sync"
	"time"
)

//Table ... models the table of a database
type Table struct {
	Name    string
	Columns []*Column
}

//Load ... Load all informations about an table from database
func (t *Table) LoadTable() {
	query, _ := translate.QT(legend.QueryCOLUMNS, t.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := db.Conx().QueryContext(ctx, query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	group := &sync.WaitGroup{}
	i := 0
	for rows.Next() {

		var (
			name     string
			position int
			sqlType  string
		)
		if err := rows.Scan(&name, &position, &sqlType); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)

		}

		t.AddColumn(name, sqlType, position)

		group.Add(1)

		go func(c *Column) {
			c.Load(t.Name)
			group.Done()

		}(t.Columns[i])

		i++
	}

	group.Wait()

	//Sort columns
	sort.Slice(t.Columns, func(i, j int) bool {
		return t.Columns[i].Position < t.Columns[j].Position
	})

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

func (t *Table) AddColumn(name, sqlType string, position int) {

	c := &Column{
		Name:        name,
		Position:    position,
		Type:        strings.ToUpper(sqlType),
		Constraints: []*Constraint{},
	}

	t.Columns = append(t.Columns, c)

}

func (t *Table) ConstrainNotNull(column string) string {
	for i := range t.Columns {
		if column == t.Columns[i].Name {
			// str := CNotNull(t.Name, t.Columns[i].Name)
			break
		}
	}

	return ""
}
