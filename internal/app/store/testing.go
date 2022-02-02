package store

import (
	"fmt"
	"strings"
	"testing"
)

func TestStore(t *testing.T, databaseUrl string)(*Store, func(...string))  {
	t.Helper()

	config := NewConfig()
	config.DatabaseUrl = databaseUrl
	store := New(config)
	err := store.Open()
	if err != nil {
		t.Fatal(err)
	}
	
	return store, func(tables ...string) {
		if len(tables) > 0 {
			_, err = store.db.Exec(fmt.Sprintf(`truncate %s cascade`, strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}
		store.Close()
	}
}
