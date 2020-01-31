package expose

import (
	"context"
	"database/sql"
	"errors"
	"sgbd4/go/db"
	"sgbd4/go/legend"
	"sgbd4/go/response"
	"sgbd4/go/translate"
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

type tableForeingKeys struct {
	TableName   string
	ForeingKeys []*db.Constraint
}

func executeQuery(tx *sql.Tx, query string) error {
	stmt, err := tx.PrepareContext(context.Background(), query)
	defer stmt.Close()

	if err != nil {
		tx.Rollback()
		return errors.New("Nu s-a putut pregatit tranzactia")
	}

	_, err = stmt.Exec()

	if err != nil {
		tx.Rollback()
		return errors.New("Nu s-a putut executa tranzactia")
	}

	return nil
}

func FixPrimaryKey(tableName, primaryKeyName string) response.Message {
	tables := *db.DB().Tables()

	isInPrimaryKeys, primaryKeys, err := takePrimaryKeys(tables, tableName)

	errMessage := response.Message{
		Type:    legend.TypeError,
		Message: translate.T(legend.MessagePrimaryKeyFailFix),
	}

	if len(primaryKeys) < 1 {
		return errMessage
	}

	isMultiplePrimaryKeys := len(primaryKeys) > 1
	//Check foreign key column is in primary keys name

	tx, err := db.DB().Conx().BeginTx(context.Background(), nil)

	if err != nil {
		return errMessage
	}

	//Drop constraints foreignkeys
	err = tables.Iterate(func(table *db.Table, column *db.Column, constraint *db.Constraint) error {
		if !constraint.IsForeignKey() || constraint.ForeingTableName != tableName || !isInPrimaryKeys(constraint.ForeingColumnName) {
			return nil
		}

		query, err := translate.QT(legend.QueryREMOVECONSTRAINT, table.Name, constraint.Name)

		if err != nil {
			return errors.New("Query nu a putut fii definit")
		}

		return executeQuery(tx, query)

	})

	//Take first, because is same constraint
	primaryKey := primaryKeys[0]

	//Remove primary key constraints
	query, err := translate.QT(legend.QueryREMOVECONSTRAINT, tableName, primaryKey.Name)

	if err != nil {
		return errMessage
	}

	err = executeQuery(tx, query)

	if err != nil {

		return errMessage
	}

	//Add primary key constraint
	query, err = translate.QT(legend.QueryADDPRIMARYKEY, tableName, primaryKeyName)

	if err != nil {
		return errMessage
	}

	err = executeQuery(tx, query)

	if err != nil {
		return errMessage
	}

	//DROP PRIMARY KEYS
	if !isMultiplePrimaryKeys {

		//Remake connections
		err = tables.Iterate(func(table *db.Table, column *db.Column, constraint *db.Constraint) error {
			if !constraint.IsForeignKey() || constraint.ForeingTableName != tableName || !isInPrimaryKeys(constraint.ForeingColumnName) {
				return nil
			}

			query, err := translate.QT(legend.QueryREMOVECONSTRAINT, table.Name, constraint.Name)

			if err != nil {
				return errors.New("Query nu a putut fii definit")
			}

			return executeQuery(tx, query)

		})

		for _, primaryKey := range primaryKeys {
			query, err := translate.QT(legend.QueryREMOVECOLUMN, tableName, primaryKey.Name)

			if err != nil {
				return errMessage
			}

			err = executeQuery(tx, query)

			if err != nil {

				return errMessage
			}
		}
	} else {
		//Remake foreing constraints
		err = tables.Iterate(func(table *db.Table, column *db.Column, constraint *db.Constraint) error {
			if !constraint.IsForeignKey() || constraint.ForeingTableName != tableName || !isInPrimaryKeys(constraint.ForeingColumnName) {
				return nil
			}
			query, err := translate.QT(legend.QueryADDFOREIGNKEY, table.Name, constraint.Name,
				column.Name, constraint.ForeingTableName, constraint.ForeingColumnName,
				constraint.UpdateRule, constraint.DeleteRule)

			if err != nil {
				return errors.New("Query nu a putut fii definit")
			}

			return executeQuery(tx, query)

		})
	}

	if err != nil {
		return errMessage
	}

	err = tx.Commit()

	if err != nil {
		return errMessage
	}

	return response.Message{
		Type:    legend.TypeSucces,
		Message: translate.T(legend.MessagePrimaryKeySuccessFix),
		Data:    db.DB().ResetTables().Tables(),
	}
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
		return nil, nil, errors.New("Nu exista tabelu")
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
