package expose

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/translate"
	"strings"
)

var (
	queryNotDefined        = errors.New("Interogarea nu a putut fii definita")
	transactionNotPrepared = errors.New("Nu s-a putut pregatit tranzactia")
	transactionNotExecuted = errors.New("Nu s-a putut executa tranzactia")
)

func AddPrimaryKey(table, primaryKeyName string) response.Message {
	query, err := translate.QT(legend.QueryADDPRIMARYKEY, table, primaryKeyName)

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessagePrimaryKeyFail, primaryKeyName),
		}
	}

	_, err = db.DB().Conx().ExecContext(context.Background(), query)

	if err != nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessagePrimaryKeyFail, primaryKeyName),
		}
	}

	//Reload table
	data := db.Table{
		Name: table,
	}

	data.LoadTable()

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T(legend.MessagePrimaryKeySuccess, primaryKeyName),
		Data:    data,
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func FixPrimaryKey(tableName, primaryKeyName string) (res response.Message) {

	defer func() {
		if r := recover(); r != nil {
			res = response.Message{
				Type:    legend.TypeError,
				Message: translate.T(legend.MessagePrimaryKeyFailFix, fmt.Sprint(r)),
			}
		}
	}()

	tables := *db.DB().Tables()

	isInPrimaryKeys, primaryKeys, err := takePrimaryKeys(tables, tableName)

	if len(primaryKeys) < 1 {
		panic("Nu exista nicio cheie primara")
	}

	tx, err := db.DB().Conx().BeginTx(context.Background(), nil)

	panicOnError(err)

	//Drop constraints foreignkeys

	panicOnError(tables.Iterate(func(table *db.Table, column *db.Column, constraint *db.Constraint) error {
		if !constraint.IsForeignKey() || constraint.ForeingTableName != tableName || !isInPrimaryKeys(constraint.ForeingColumnName) {
			return nil
		}

		query, err := translate.QT(legend.QueryREMOVECONSTRAINT, table.Name, constraint.Name)

		if err != nil {
			return queryNotDefined
		}

		return executeQuery(tx, query)

	}))

	//Take first, because is same constraint
	primaryKey := primaryKeys[0]

	panicOnError(removeAndAddPrimaryKey(tx, tableName, primaryKeyName, primaryKey.Name))

	isMultiplePrimaryKeys := len(primaryKeys) > 1

	//DROP PRIMARY KEYS
	if isMultiplePrimaryKeys {

		//Remake constraints
		panicOnError(remakeConstraints(tx, &tables, tableName, isInPrimaryKeys))

	} else {

		//Remake foreing constraints
		panicOnError(remakeColumns(tx, &tables, primaryKeys, primaryKeyName, tableName, isInPrimaryKeys))
	}

	err = tx.Commit()

	panicOnError(err)

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T(legend.MessagePrimaryKeySuccessFix),
		Data:    db.DB().ResetTables().Tables(),
	}
}

func executeQuery(tx *sql.Tx, query string) error {

	stmt, err := tx.PrepareContext(context.Background(), query)

	defer stmt.Close()

	if err != nil {
		tx.Rollback()
		return transactionNotPrepared
	}

	_, err = stmt.Exec()

	if err != nil {
		tx.Rollback()
		return transactionNotExecuted
	}

	return nil
}

func takePrimaryKeys(tables db.Tables, tableName string) (func(string) bool, []*db.Column, error) {
	var table *db.Table
	//Take table
	for i := range tables {
		if tables[i].Name == tableName {
			table = tables[i]
		}
	}

	if table == nil {
		panic("Nu exista tabelu")
	}

	//Take primaryKey/Keys
	primaryKeys := make([]*db.Column, 0)

	for i := range table.Columns {
		if table.Columns[i].HasPrimaryKey() {
			primaryKeys = append(primaryKeys, table.Columns[i])
		}
	}

	return func(name string) bool {
		for _, column := range primaryKeys {
			if name == column.Name {
				return true
			}
		}

		return false
	}, primaryKeys, nil

}

func remakeColumns(tx *sql.Tx, tables *db.Tables, primaryKeys []*db.Column, primaryKeyName, tableName string, isInPrimaryKeys func(string) bool) error {
	return tables.Iterate(func(table *db.Table, column *db.Column, constraint *db.Constraint) error {
		if !constraint.IsForeignKey() || constraint.ForeingTableName != tableName || !isInPrimaryKeys(constraint.ForeingColumnName) {
			return nil
		}

		//Remove any constraint for this column
		for _, constr := range column.Constraints {
			query, err := translate.QT(legend.QueryREMOVECONSTRAINT, table.Name, constr.Name)
			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)
			if err != nil {
				return err
			}
		}

		//Create view with values temporary from the new primary key column
		const alias = "value"
		query, err := translate.QT(legend.QueryCREATEVIEW, strings.Join([]string{column.Name, table.Name, constraint.Name}, "_"), tableName, primaryKeyName,
			primaryKeys[0].Name, column.Name, table.Name, alias)

		if err != nil {
			return queryNotDefined
		}

		err = executeQuery(tx, query)

		if err != nil {
			return err
		}

		//Drop this column

		query, err = translate.QT(legend.QueryREMOVECOLUMN, table.Name, column.Name)

		if err != nil {
			return queryNotDefined
		}

		err = executeQuery(tx, query)
		if err != nil {
			return err
		}

		//Load column with primary key to take type of new primary key

		aux := db.Column{
			Name: primaryKeyName,
		}

		aux.Load(tableName)

		//Create column back
		query, err = translate.QT(legend.QueryADDCOLUMN, table.Name, column.Name, aux.Type)

		if err != nil {
			return queryNotDefined
		}

		err = executeQuery(tx, query)

		if err != nil {
			return err
		}

		//Take values from view and add them to the new column
		query, err = translate.QT(legend.QueryREMAKECOLUMNS, table.Name, column.Name, alias, strings.Join([]string{column.Name, table.Name, constraint.Name}, "_"))

		if err != nil {
			return queryNotDefined
		}

		err = executeQuery(tx, query)
		if err != nil {
			return err
		}

		//Add back constraints
		for _, constr := range column.Constraints {
			var (
				query string
				err   error
			)
			if !constr.IsForeignKey() {
				//Add normal constraints, NOT NULL ,CHECK etc
				query, err = translate.QT(legend.QueryADDCONSTRAINT, table.Name, constr.Name)
			} else {

				name := primaryKeyName
				//Check to noy bew foreign key on another column
				if !isInPrimaryKeys(constr.ForeingColumnName) {
					name = constr.ForeingColumnName
				}
				query, err = translate.QT(legend.QueryADDFOREIGNKEY, table.Name, constr.Name, column.Name, constr.ForeingTableName, name)
			}

			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)

			if err != nil {
				return err
			}
		}

		return err

	})
}

func remakeConstraints(tx *sql.Tx, tables *db.Tables, tableName string, isInPrimaryKeys func(string) bool) error {
	return tables.Iterate(func(table *db.Table, column *db.Column, constraint *db.Constraint) error {
		if !constraint.IsForeignKey() || constraint.ForeingTableName != tableName || !isInPrimaryKeys(constraint.ForeingColumnName) {
			return nil
		}

		query, err := translate.QT(legend.QueryADDFOREIGNKEY, table.Name, constraint.Name, column.Name, constraint.ForeingTableName, constraint.ForeingColumnName)

		if err != nil {
			return queryNotDefined
		}

		return executeQuery(tx, query)

	})
}

func removeAndAddPrimaryKey(tx *sql.Tx, tableName, newPrimaryKeyName, oldPrimaryKeyName string) error {
	//Remove primary key constraints
	query, err := translate.QT(legend.QueryREMOVECONSTRAINT, tableName, oldPrimaryKeyName)

	if err != nil {
		return err
	}

	err = executeQuery(tx, query)

	if err != nil {
		return err
	}

	//Add primary key constraint
	query, err = translate.QT(legend.QueryADDPRIMARYKEY, tableName, newPrimaryKeyName)

	if err != nil {
		return err
	}

	err = executeQuery(tx, query)

	if err != nil {
		return err
	}

	return nil
}
