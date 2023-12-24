package test

import (
	"errors"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"golang.org/x/net/context"
)

var client *db.Client

// initializes firesbase app and returns a db client
func initializeFirebaseDb() (*db.Client, error) {

	if client != nil {
		log.Println("reusing existing firebase db client")
		return client, nil
	}

	// initialize firebase app
	// config and options are read automatically if
	// GOOGLE_APPLICATION_CREDENTIALS environment variable is set.
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}

	databaseURL := os.Getenv("FIREBASE_DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("FIREBASE_DATABASE_URL not set")
	}
	client, err = app.DatabaseWithURL(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("error initializing database client: %v\n", err)
		return nil, err
	}

	return client, nil
}

type FbTestDB struct {
	client *db.Client
}

func NewFbTestDB(client *db.Client) *FbTestDB {
	return &FbTestDB{
		client: client,
	}
}

func (db *FbTestDB) insert(t Test) {
	ref := db.client.NewRef("tests/" + t.ID)
	if err := ref.Set(context.TODO(), t); err != nil {
		log.Fatal("error setting value:", err)
	}
	log.Println("inserted test:", t.ID)
}

func (db *FbTestDB) get(id string) (Test, error) {
	ref := db.client.NewRef("tests/" + id)
	var t Test
	if err := ref.Get(context.Background(), &t); err != nil {
		log.Fatal("error getting value:", err)
		return Test{}, err
	}
	log.Println("got test:", t.ID)
	return t, nil
}

func (db *FbTestDB) list() ([]Test, error) {
	ref := db.client.NewRef("tests/")
	var res map[string]Test
	if err := ref.Get(context.Background(), &res); err != nil {
		log.Println("error getting tests:", err)
		return []Test{}, err
	}
	tests := make([]Test, 0)
	for _, v := range res {
		tests = append(tests, v)
	}
	return tests, nil
}
