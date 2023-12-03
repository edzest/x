package test

import "fmt"

// TestStore is an interface that provides reading and writing
// a Test in a datastore.
type TestStore interface {
	insert(t Test)
	get(id string) (Test, error)
}

// TempTestDb is an in-memory DB that satisfies TestStore interface.
type TempTestDB struct {
	records map[string]Test
}

func NewTempTestDB() *TempTestDB {
	return &TempTestDB{
		make(map[string]Test),
	}
}

var ErrorTestNotFound = fmt.Errorf("test not found")

func (db *TempTestDB) insert(t Test) {
	db.records[t.ID] = t
}

func (db *TempTestDB) get(id string) (Test, error) {
	if t, ok := db.records[id]; ok {
		return t, nil
	}
	return Test{}, ErrorTestNotFound
}
