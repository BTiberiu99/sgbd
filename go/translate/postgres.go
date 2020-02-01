package translate

import (
	"fmt"
	"sgbd4/go/legend"
)

var (
	postgres = map[string]func(...string) string{

		legend.QuerySETNOTNULL: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului si al coloanei")
			}
			return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;", in[0], in[1])
		},
		legend.QueryTABLES: func(in ...string) string {

			return fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
		},

		legend.QueryCOLUMNS: func(in ...string) string {
			if len(in) < 1 {
				panic("Trebuie sa introduceti numele tabelului")
			}

			return fmt.Sprintf(`SELECT column_name,ordinal_position,udt_name
								FROM information_schema.columns 
								WHERE table_schema = 'public' AND table_name  = '%s';`, in[0])
		},

		legend.QueryCONSTRAINTS: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului si numele coloanei")
			}

			return fmt.Sprintf(`SELECT distinct
				ccu.constraint_name,
				tc.constraint_type,
				ccu.table_name AS foreign_table_name,
				ccu.column_name AS foreign_column_name,
				coalesce(rc.update_rule, '') AS update_rule,
				coalesce(rc.delete_rule, '')  AS delete_rule
				FROM
				information_schema.constraint_column_usage as ccu
				JOIN information_schema.key_column_usage as kcu
				ON  kcu.constraint_name = ccu.constraint_name
				AND  kcu.table_schema = ccu.table_schema
				JOIN information_schema.table_constraints as tc
				ON  tc.constraint_name = ccu.constraint_name
				AND tc.table_schema = ccu.table_schema
				LEFT JOIN information_schema.referential_constraints AS rc
				ON rc.constraint_name = ccu.constraint_name
				WHERE ccu.column_name='%s' AND  ccu.table_name='%s';`, in[1], in[0])
		},

		legend.QueryCOUNTNOTNULL: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului si numele coloanei")
			}
			return fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s IS NULL", in[0], in[1])
		},

		legend.QueryCHECKCONSTRAINTS: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului si numele coloanei")
			}

			return fmt.Sprintf(`SELECT constraint_name,check_clause
								FROM information_schema.check_constraints
								WHERE constraint_name  IN (SELECT constraint_name FROM information_schema.table_constraints WHERE table_name = '%s')
								AND check_clause  LIKE '%%%s%%';`, in[0], in[1])
		},

		legend.QueryADDPRIMARYKEY: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului si numele coloanei")
			}
			return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s SERIAL PRIMARY KEY;", in[0], in[1])
		},

		legend.QueryREMOVECOLUMN: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului si numele coloanei")
			}
			return fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s;", in[0], in[1])
		},
		legend.QueryADDCOLUMN: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului, numele coloanei si tipul coloanei")
			}
			return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", in[0], in[1], in[2])
		},
		legend.QueryREMOVECONSTRAINT: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului si numele constrangeri")
			}
			return fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;", in[0], in[1])
		},
		legend.QueryADDCONSTRAINT: func(in ...string) string {
			if len(in) < 2 {
				panic("Trebuie sa introduceti numele tabelului, numele constrangeri si definitia constrangeri")
			}
			return fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s;", in[0], in[1], in[2])
		},
		legend.QueryADDFOREIGNKEY: func(in ...string) string {
			if len(in) < 5 {
				panic("Trebuie sa introduceti numele tabelului,numele constrangeri,coloana care face referinta , tabelu si coloana catre care se face referinta")
			}
			simple := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)", in[0], in[1], in[2], in[3], in[4])

			if len(in) > 5 && in[5] != "" {
				simple += fmt.Sprintf(" ON UPDATE %s", in[5])
			}

			if len(in) > 6 && in[6] != "" {
				simple += fmt.Sprintf(" ON DELETE %s", in[6])
			}

			return simple
		},

		legend.QueryREMAKECOLUMNS: func(in ...string) string {
			if len(in) < 4 {
				panic("Trebuie sa introduceti numele tabelului,numele coloanei careia i se face update , numele coloanei din view  si numele view-ului din care se face update")
			}
			simple := fmt.Sprintf(`UPDATE %s t1 SET %s = t2.%s FROM %s t2 WHERE t1.row_number() = t2.row_number() `, in[0], in[1], in[3], in[2])

			return simple
		},
		legend.QueryCREATEVIEW: func(in ...string) string {
			if len(in) < 6 {
				panic("Trebuie sa introduceti numele view-ului,numele tabelului din care se iau datele, numele coloanei, numele coloanei din join, numele tabelului cu care se face join si numele coloanei care se inlocuieste, si alias pentru valoare")
			}

			return fmt.Sprintf(`CREATE TEMPORARY VIEW temp%s AS (SELECT t.%s AS %s FROM %s t JOIN %s t2 ON t2.%s = t.%s);`, in[0], in[2], in[6], in[1], in[4], in[5], in[3])

		},
	}
)
