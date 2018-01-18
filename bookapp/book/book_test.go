package book

import (
	"testing"

	"github.com/lab46/example/pkg/rdbms"

	"github.com/lab46/example/pkg/testutil/sqlimporter"
)

var (
	testDSN    = "postgres://exampleapp:exampleapp@localhost:5432?sslmode=disable"
	testDriver = "postgres"
	schemaPath = "../../files/dbschema/book"
	dataPath   = "testdata"
)

func TestAddBook(t *testing.T) {
	db, err, drop := sqlimporter.CreateDB(testDriver, testDSN)
	if err != nil {
		t.Error(err)
	}
	defer drop()
	// import schema
	err = sqlimporter.ImportSchemaFromFiles(db, schemaPath)
	if err != nil {
		t.Error(err)
	}

	bs := Init(db, rdbms.NewLoadBalancer(db))
	err = bs.AddBook(Book{
		Title:  "test1",
		Author: "author1",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestListOfBooks(t *testing.T) {
	db, err, drop := sqlimporter.CreateDB(testDriver, testDSN)
	if err != nil {
		t.Error(err)
	}
	defer drop()
	// import schema
	err = sqlimporter.ImportSchemaFromFiles(db, schemaPath)
	if err != nil {
		t.Error(err)
	}
	// import data
	err = sqlimporter.ImportSchemaFromFiles(db, dataPath)
	if err != nil {
		t.Error(err)
	}
	bs := Init(db, rdbms.NewLoadBalancer(db))
	books, err := bs.ListOfBooks()
	if err != nil {
		t.Error(err)
	}
	if len(books) != 2 {
		t.Errorf("Got %d books, but expecting 2 from test file", len(books))
	}
}
