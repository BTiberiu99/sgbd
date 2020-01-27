package db

import "strings"

import "sgbd4/go/legend"

type Constraint struct {
	Name string
	Type string
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
