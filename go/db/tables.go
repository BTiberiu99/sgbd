package db

import "sgbd4/go/translate"

import "fmt"

type Tables struct {
	Tables []Table
}

func (t *Tables) LoadTables() error {
	qry, _ := translate.QT("tables")

	tables, err := db.Conx().Query(qry)

	if err != nil {
		return err
	}

	table := map[string]interface{}{}

	for tables.Next() {
		tables.Scan(table)
		// fmt.Println(table)
	}

	fmt.Println(table)

	return nil
}
