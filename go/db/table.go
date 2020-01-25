package db


type Table struct {
	Name    string
	Columns []Column
}


func (t *Table) Load() {

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