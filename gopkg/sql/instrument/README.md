# Instrument

Instrument is a function that return an Instrument object to observe/monitor your query/action into database.

It will also automaticall use QueryContext, SelectContext and all context type to execute operations to database.

Note: **All instrumen object is expected to used once. Do not use the instrument object twice, create a new instrumentation object from parent context instead.**

For example:
```go
q := sqldb.InstrumentQuery(ctx, db, "query_for_something")
rows, err := q.Query("SELECT something FROM table_something")
if err != nil {
    return err
}

// note that if you want to perform a new query, you need to create another instrumentation
q := ctx(ctx, db, "query_for_something_esle")
rows, err := q.Query*("SELECT something_else FROM table_something_else")
if err != nil {
    return err
}
```

 or in a more complex condition:
```go
span, ctx := opentracing.StartSpanFromContext(ctx, "function_something")
defer span.Finish()
timeout1, cancel := context.WithTimeout(ctx, time.Second*3)
q := sqldb.InstrumentQuery(timeout1, db, "query_for_something")
 rows, err := q.Query("SELECT something FROM table_something")
if err != nil {
    return err
}

// DO NOT REUSE timeout1
timeout2, cancel := context.WithTimeout(ctx, time.Second*3)
defer cancel()
// create a new instrumentation for each query
q := sqldb.InstrumentQuery(timeout2, db, "query_for_something_esle")
rows, err := q.Query*("SELECT something_else FROM table_something_else")
if err != nil {
    return err
}
```