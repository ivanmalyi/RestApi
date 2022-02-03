package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestDb(t *testing.T, databaseUrl string)(*sql.DB, func(...string))  {
	t.Helper()

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		t.Fatal(err)
	}
	
	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err = db.Exec(fmt.Sprintf(`truncate %s cascade`, strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}
		_ = db.Close()
	}
}
