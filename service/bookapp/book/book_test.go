package book

import (
	"testing"

	"github.com/lab46/example/gopkg/sqldb"
	"github.com/lab46/example/gopkg/testutil/sqlimporter"
	"github.com/lab46/example/gopkg/testutil/testenv"
)

var (
	testDriver = "postgres"
	schemaPath = "../../files/dbschema/book"
	dataPath   = "testdata"
)

func TestAddBook(t *testing.T) {
	t.Parallel()
	db, drop, err := sqlimporter.CreateRandomDB(testDriver, testenv.EnvConfig.PostgresDSN)
	if err != nil {
		t.Error(err)
	}
	defer drop()
	// import schema
	err = sqlimporter.ImportSchemaFromFiles(db, schemaPath)
	if err != nil {
		t.Error(err)
	}

	Init(db, sqldb.NewLoadBalancer(db))
	err = AddBook(Book{
		Title:  "test1",
		Author: "author1",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestListOfBooks(t *testing.T) {
	t.Parallel()
	db, drop, err := sqlimporter.CreateRandomDB(testDriver, testenv.EnvConfig.PostgresDSN)
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
	Init(db, sqldb.NewLoadBalancer(db))
	books, err := ListOfBooks()
	if err != nil {
		t.Error(err)
	}
	if len(books) != 2 {
		t.Errorf("Got %d books, but expecting 2 from test file", len(books))
	}
}

func TestGetBookyID(t *testing.T) {
	t.Parallel()
	db, drop, err := sqlimporter.CreateRandomDB(testDriver, testenv.EnvConfig.PostgresDSN)
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
	Init(db, sqldb.NewLoadBalancer(db))
	b, err := GetBookByID(1)
	if err != nil {
		t.Error(err)
	}
	if b.Title == "" {
		t.Error("book title cannot be empty")
	}
}
