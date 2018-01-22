# SQLImporter CLI

SQLImporter command line tools used to import schema from *.sql files in a directory with a special rules in *.sql file. For reference: [sqlimporter library](https://github.com/lab46/example/blob/master/pkg/testutil/sqlimporter/README.md)

This tools might useful if your environment don't have any postgresql/mysql command installed but have your database running in a container.

## How to use it 

1. Install `go`.
2. `go install github.com/lab46/example/tools/sqlimporter/...`
3. Now you can type `sqlimporter` to get some help.
```shell
sqlimporter command line tools

Usage:
  sqlimporter [command]

Available Commands:
  help        Help about any command
  import      import postgresql/mysql schema from directory
  test        test command for sqlimporter

Flags:
  -h, --help      help for sqlimporter
  -v, --verbose   sqlimporter verbose output

Use "sqlimporter [command] --help" for more information about a command.
```

To import a schema into a database:

`sqlimporter [driver] [dbname] [dsn] [directory]`

or for example:

`sqlimporter import book postgres "postgres://exampleapp:exampleapp@localhost:5432?sslmode=disable" "files/dbschema/book/"`



