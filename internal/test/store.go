package test

import "fmt"

// TestStore is an interface that provides reading from and writing to
// a datastore.
type TestStore interface {
	insert(t Test)
	get(id string) (Test, error)
}

var db map[string]Test

// Initialises a global map of tests
func init() {
	db = make(map[string]Test)
}

// TempTestDb is an in-memory DB that satisfies TestStore interface.
type TempTestDB struct {
	records map[string]Test
}

func NewTempTestDB() *TempTestDB {
	return &TempTestDB{
		records: db,
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
