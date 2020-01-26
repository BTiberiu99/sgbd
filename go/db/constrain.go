package db

import "strings"

import "sgbd4/go/legend"

type Constrain struct {
	Name string
	Type string
}

func (c *Constrain) IsPrimaryKey() bool {
	return strings.Contains(c.Type, legend.PRIMARYKEY)
}

func (c *Constrain) IsForeignKey() bool {
	return strings.Contains(c.Type, legend.FOREIGNKEY)
}

func (c *Constrain) IsNotNull() bool {
	return strings.Contains(c.Type, legend.NOTNULL)
}

func (c *Constrain) IsCheck() bool {
	return strings.Contains(c.Type, legend.CHECK)
}

func (c *Constrain) IsUnique() bool {
	return strings.Contains(c.Type, legend.UNIQUE)
}
