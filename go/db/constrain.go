package db

import "strings"

type Constrain struct {
	Name string
	Type string
	SQL  string
}

func (c *Constrain) IsUnique() bool {
	return strings.Contains(strings.ToUpper(c.SQL), "UNIQUE")
}

func (c *Constrain) IsNull() bool {
	return !strings.Contains(strings.ToUpper(c.SQL), "NOT NULL")
}

func (c *Constrain) IsPrimary() bool {
	return strings.Contains(strings.ToUpper(c.SQL), "PRIMARY KEY")
}

func (c *Constrain) IsForeignKey() bool {
	return strings.Contains(strings.ToUpper(c.SQL), "FOREIGN KEY")
}
