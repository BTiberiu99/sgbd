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

//AddPrimaryKey... set to a table without any primary key a new primary key that is autonumber
func AddPrimaryKey(table, primaryKeyName string) response.Message {

	if db.DB() == nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessageNoConnection),
		}
	}

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

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T(legend.MessagePrimaryKeySuccess, primaryKeyName),
		Data:    db.DB().ResetTables().Tables(),
	}
}

//FixPrimaryKey ... solves two cases of the normalization of the primary key,
//first is the one when an table has multiple primary kers the other one is when the primary key is not numeric
func FixPrimaryKey(tableName, primaryKeyName string) (res response.Message) {
	if db.DB() == nil {
		return response.Message{
			Type:    legend.TypeError,
			Message: translate.T(legend.MessageNoConnection),
		}
	}
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

	for _, key := range primaryKeys {
		fmt.Println(key)
		for _, constr := range key.Constraints {
			if constr.IsForeignKey() {
				query, err := translate.QT(legend.QueryREMOVECONSTRAINT, constr.ForeingTableName, constr.Name)

				if err != nil {
					panicOnError(queryNotDefined)
				}

				panicOnError(executeQuery(tx, query))
			}

		}
	}

	//Take first, because is same constraint
	primaryKey := primaryKeys[0]

	var constr *db.Constraint

	for _, con := range primaryKey.Constraints {
		if con.IsPrimaryKey() {
			constr = con
			break
		}
	}

	fmt.Println(*constr)

	panicOnError(removeAndAddPrimaryKey(tx, tableName, primaryKeyName, constr.Name))

	isMultiplePrimaryKeys := len(primaryKeys) > 1

	//DROP PRIMARY KEYS
	if isMultiplePrimaryKeys {

		//Remake constraints
		panicOnError(remakeConstraints(tx, &tables, tableName, isInPrimaryKeys))

	} else {

		//Remake foreing constraints
		panicOnError(remakeColumns(tx, &tables, primaryKey, primaryKeyName, tableName))

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
	fmt.Println(err, query)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return transactionNotPrepared
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		fmt.Println(err)
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

func remakeColumns(tx *sql.Tx, tables *db.Tables, primaryKey *db.Column, primaryKeyName, tableName string) error {
	//Load column with primary key to take type of new primary key

	for _, constr := range primaryKey.Constraints {

		if constr.IsForeignKey() {

			//ADD HELPER
			query, err := translate.QT(legend.QueryADDHELPER, constr.ForeingTableName)

			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)
			if err != nil {
				return err
			}

			//Create view with values temporary from the new primary key column
			const alias = "value"

			viewName := strings.Join([]string{constr.ForeingTableName, constr.ForeingColumnName}, "_")

			query, err = translate.QT(legend.QueryCREATEVIEW, viewName, tableName, primaryKeyName,
				primaryKey.Name, constr.ForeingTableName, constr.ForeingColumnName, alias)

			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)

			if err != nil {
				return err
			}

			column := tables.FindColumn(constr.ForeingTableName, constr.ForeingColumnName)

			fmt.Println(constr.ForeingTableName, constr.ForeingColumnName, column.Name)
			//Drop all constraints on column
			for _, constrCol := range column.Constraints {

				if constrCol.IsForeignKey() || constrCol.IsPrimaryKey() {
					query, err := translate.QT(legend.QueryREMOVECONSTRAINT, constrCol.ForeingTableName, constrCol.Name)
					if err != nil {
						return queryNotDefined
					}

					err = executeQuery(tx, query)
					if err != nil {
						return err
					}
				}

			}

			//Drop this column
			query, err = translate.QT(legend.QueryREMOVECOLUMN, constr.ForeingTableName, constr.ForeingColumnName)

			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)
			if err != nil {
				return err
			}

			//Create column back
			query, err = translate.QT(legend.QueryADDCOLUMN, constr.ForeingTableName, constr.ForeingColumnName, "INT4")

			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)

			if err != nil {
				return err
			}
			fmt.Println("Am ajuns aici 3")
			//Take values from view and add them to the new column
			query, err = translate.QT(legend.QueryREMAKECOLUMNS, constr.ForeingTableName, constr.ForeingColumnName, viewName, alias)

			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)
			if err != nil {
				return err
			}

			//ADD HELPER
			query, err = translate.QT(legend.QueryREMOVEHELPER, constr.ForeingTableName)

			if err != nil {
				return queryNotDefined
			}

			err = executeQuery(tx, query)
			if err != nil {
				return err
			}

			//Remake constraints on column

			//Add back constraints
			for _, constrCol := range column.Constraints {
				var (
					query string
					err   error
				)
				if !constr.IsForeignKey() {
					//Add normal constraints, NOT NULL ,CHECK etc
					query, err = translate.QT(legend.QueryADDCONSTRAINT, constr.ForeingTableName, constrCol.Name, constrCol.Type)
				} else {

					query, err = translate.QT(legend.QueryADDFOREIGNKEY, constrCol.Name, constrCol.ForeingTableName, column.Name, constr.ForeingTableName)
				}

				if err != nil {
					return queryNotDefined
				}

				err = executeQuery(tx, query)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//panics if error is not null
func panicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
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

func removeAndAddPrimaryKey(tx *sql.Tx, tableName, newPrimaryKeyName, oldConstraintPrimaryKeyName string) error {
	//Remove primary key constraints
	query, err := translate.QT(legend.QueryREMOVECONSTRAINT, tableName, oldConstraintPrimaryKeyName)

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
