package db

import (
	"sgbd4/go/legend"
	"strings"
)

//Constraint ... models the constraints of a column
type Constraint struct {
	Name              string
	Type              string
	ForeingTableName  string `json:"foreing_table_name,omitempty"`
	ForeingColumnName string `json:"foreing_column_name,omitempty"`
	UpdateRule        string `json:"update_rule,omitempty"`
	DeleteRule        string `json:"delete_rule,omitempty"`
}

func (c *Constraint) IsPrimaryKey() bool {
	return strings.Contains(c.Type, legend.PRIMARYKEY)
}

func (c *Constraint) IsForeignKey() bool {
	return strings.Contains(c.Type, legend.FOREIGNKEY)
}

func (c *Constraint) IsNotNull() bool {
	return strings.Contains(c.Type, legend.NOTNULL)
}

func (c *Constraint) IsCheck() bool {
	return strings.Contains(c.Type, legend.CHECK)
}

func (c *Constraint) IsUnique() bool {
	return strings.Contains(c.Type, legend.UNIQUE)
}
