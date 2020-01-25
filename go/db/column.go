package db

type Column struct {
	Name        string
	Constraints []Constrain
	Type        string
}

func (c *Column) LoadConstrains() error {
	return nil
}

func (c *Column) AddConstrain(name, constrainType, sql string) {
	c.Constraints = append(c.Constraints, Constrain{Name: name, Type: constrainType, SQL: sql})
}

func (c *Column) HasUnique() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsUnique()
	})
}

func (c *Column) HasNull() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsNull() && !c.IsPrimary()
	})
}

func (c *Column) HasPrimaryKey() bool {
	return c.iterateConstrains(func(c *Constrain) bool {
		return c.IsPrimary()
	})
}
func (c *Column) iterateConstrains(s func(c *Constrain) bool) bool {
	for i := range c.Constraints {
		if s(&c.Constraints[i]) {
			return true
		}
	}
	return false
}
