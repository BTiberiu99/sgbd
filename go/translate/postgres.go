package translate

import "fmt"

var (
	postgres = map[string]func(...string) string{
		"not_null": func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduci numele tabelului si al coloanei")
			}
			return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;", in[0], in[1])
		},
		"tables": func(in ...string) string {

			return fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
		},

		"columns": func(in ...string) string {
			if len(in) < 1 {
				panic("Trebuie sa introduci numele tabelului")
			}

			return fmt.Sprintf(`SELECT column_name,ordinal_position,udt_name
								FROM information_schema.columns 
								WHERE table_schema = 'public' AND table_name  = '%s';`, in[0])
		},

		"constrains": func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduci numele tabelului si numele coloanei")
			}

			return fmt.Sprintf(`SELECT constraint_name,constraint_type
							FROM information_schema.table_constraints
							WHERE constraint_name IN (SELECT constraint_name from information_schema.constraint_column_usage 
													WHERE table_schema = 'public' AND table_name = '%s' and column_name = '%s')
							AND table_schema = 'public'
							AND table_name   = '%s';`, in[0], in[1], in[0])
		},

		"check_constrains": func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduci numele tabelului si numele coloanei")
			}

			return fmt.Sprintf(`SELECT constraint_name,check_clause
								FROM information_schema.check_constraints
								WHERE constraint_name  IN (SELECT constraint_name FROM information_schema.table_constraints WHERE table_name = '%s')
								AND check_clause  LIKE '%%%s%%';`, in[0], in[1])
		},
	}
)
