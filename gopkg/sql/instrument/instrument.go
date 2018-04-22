package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/lab46/monorepo/gopkg/tracing"
)

func finish(span opentracing.Span, t time.Time) {
	span.Finish()
}

// QueryInstrument instrument query by using sqlx.DB
type QueryInstrument struct {
	db        *sqlx.DB
	ctx       context.Context
	queryname string
}

// InstrumentQuery return QueryInstrument
func InstrumentQuery(ctx context.Context, db *sqlx.DB, queryname string) *QueryInstrument {
	q := QueryInstrument{
		db:        db,
		ctx:       ctx,
		queryname: queryname,
	}
	return &q
}

// Begin queryinstrument
func (qi QueryInstrument) Begin(ctx context.Context) (*sql.Tx, error) {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	tx, err := qi.db.BeginTx(ctx, nil)
	return tx, err
}

// Beginx queryinstrument
func (qi QueryInstrument) Beginx(ctx context.Context) (*sqlx.Tx, error) {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	tx, err := qi.db.BeginTxx(ctx, nil)
	return tx, err
}

// Query queryinstrument
func (qi *QueryInstrument) Query(query string, args ...interface{}) (*sql.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	rows, err := qi.db.QueryContext(ctx, query, args...)
	return rows, err
}

// Queryx queryinstrument
func (qi *QueryInstrument) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	rows, err := qi.db.QueryxContext(ctx, query, args...)
	return rows, err
}

// QueryRow queryinstrument
func (qi *QueryInstrument) QueryRow(query string, args ...interface{}) *sql.Row {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	row := qi.db.QueryRowContext(ctx, query, args...)
	return row
}

// QueryRowx queryinstrument
func (qi *QueryInstrument) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	row := qi.db.QueryRowxContext(ctx, query, args...)
	return row
}

// Get queryinstrument
func (qi *QueryInstrument) Get(dest interface{}, query string, args ...interface{}) error {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	err := qi.db.GetContext(ctx, dest, query, args...)
	return err
}

// Select queryinstruyment
func (qi *QueryInstrument) Select(dest interface{}, query string, args ...interface{}) error {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	err := qi.db.SelectContext(ctx, dest, query, args...)
	return err
}

// Exec queryinstrument
func (qi *QueryInstrument) Exec(query string, args ...interface{}) (sql.Result, error) {
	span, ctx := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	result, err := qi.db.ExecContext(ctx, query, args...)
	return result, err
}

// Rebind queryinstrument
func (qi *QueryInstrument) Rebind(query string) string {
	span, _ := tracing.StartSpanFromContext(qi.ctx, qi.queryname)
	defer finish(span, time.Now())
	return qi.db.Rebind(query)
}

// StatementInstrument instrument statement operations
type StatementInstrument struct {
	stmt          *sql.Stmt
	ctx           context.Context
	statementname string
}

// InstrumentStatement return StatementInstrument
func InstrumentStatement(ctx context.Context, stmt *sql.Stmt, statementname string) *StatementInstrument {
	st := StatementInstrument{
		stmt:          stmt,
		ctx:           ctx,
		statementname: statementname,
	}
	return &st
}

// Query of statement
func (sti *StatementInstrument) Query(args ...interface{}) (*sql.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(sti.ctx, sti.statementname)
	defer finish(span, time.Now())
	rows, err := sti.stmt.QueryContext(ctx, args...)
	return rows, err
}

// QueryRow queryinstrument
func (sti *StatementInstrument) QueryRow(args ...interface{}) *sql.Row {
	span, ctx := tracing.StartSpanFromContext(sti.ctx, sti.statementname)
	defer finish(span, time.Now())
	row := sti.stmt.QueryRowContext(ctx, args...)
	return row
}

// QueryRowx queryinstrument
func (sti *StatementInstrument) QueryRowx(args ...interface{}) *sql.Row {
	span, ctx := tracing.StartSpanFromContext(sti.ctx, sti.statementname)
	defer finish(span, time.Now())
	row := sti.stmt.QueryRowContext(ctx, args...)
	return row
}

// Exec queryinstrument
func (sti *StatementInstrument) Exec(args ...interface{}) (sql.Result, error) {
	span, ctx := tracing.StartSpanFromContext(sti.ctx, sti.statementname)
	defer finish(span, time.Now())
	result, err := sti.stmt.ExecContext(ctx, args...)
	return result, err
}

// StatementxInstrument struct
type StatementxInstrument struct {
	stmtx         *sqlx.Stmt
	ctx           context.Context
	statementname string
}

// InstrumentStatementx return StatementxInstrument
func InstrumentStatementx(ctx context.Context, stmtx *sqlx.Stmt, statementname string) *StatementxInstrument {
	stix := StatementxInstrument{
		stmtx:         stmtx,
		ctx:           ctx,
		statementname: statementname,
	}
	return &stix
}

// Query function
func (stix *StatementxInstrument) Query(args ...interface{}) (*sql.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(stix.ctx, stix.statementname)
	defer finish(span, time.Now())
	rows, err := stix.stmtx.QueryContext(ctx, args...)
	return rows, err
}

// Queryx function
func (stix *StatementxInstrument) Queryx(args ...interface{}) (*sqlx.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(stix.ctx, stix.statementname)
	defer finish(span, time.Now())
	rows, err := stix.stmtx.QueryxContext(ctx, args...)
	return rows, err
}

// Get function
func (stix *StatementxInstrument) Get(dest interface{}, args ...interface{}) error {
	span, ctx := tracing.StartSpanFromContext(stix.ctx, stix.statementname)
	defer finish(span, time.Now())
	err := stix.stmtx.GetContext(ctx, dest, args...)
	return err
}

// Select function
func (stix *StatementxInstrument) Select(dest interface{}, args ...interface{}) error {
	span, ctx := tracing.StartSpanFromContext(stix.ctx, stix.statementname)
	defer finish(span, time.Now())
	err := stix.stmtx.SelectContext(ctx, dest, args...)
	return err
}

// Exec function
func (stix *StatementxInstrument) Exec(dest interface{}, args ...interface{}) (sql.Result, error) {
	span, ctx := tracing.StartSpanFromContext(stix.ctx, stix.statementname)
	defer finish(span, time.Now())
	result, err := stix.stmtx.ExecContext(ctx, args...)
	return result, err
}

// TxInstrument struct
type TxInstrument struct {
	tx        *sql.Tx
	ctx       context.Context
	queryname string
}

// InstrumentTx return TxInstrument
func InstrumentTx(ctx context.Context, tx *sql.Tx, queryname string) *TxInstrument {
	txi := TxInstrument{
		tx:        tx,
		ctx:       ctx,
		queryname: queryname,
	}
	return &txi
}

// Query queryinstrument
func (tx *TxInstrument) Query(query string, args ...interface{}) (*sql.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(tx.ctx, tx.queryname)
	defer finish(span, time.Now())
	rows, err := tx.tx.QueryContext(ctx, query, args...)
	return rows, err
}

// QueryRow queryinstrument
func (tx *TxInstrument) QueryRow(query string, args ...interface{}) *sql.Row {
	span, ctx := tracing.StartSpanFromContext(tx.ctx, tx.queryname)
	defer finish(span, time.Now())
	row := tx.tx.QueryRowContext(ctx, query, args...)
	return row
}

// Exec queryinstrument
func (tx *TxInstrument) Exec(query string, args ...interface{}) (sql.Result, error) {
	span, ctx := tracing.StartSpanFromContext(tx.ctx, tx.queryname)
	defer finish(span, time.Now())
	result, err := tx.tx.ExecContext(ctx, query, args...)
	return result, err
}

// Commit queryinstrument
func (tx *TxInstrument) Commit() error {
	span, _ := tracing.StartSpanFromContext(tx.ctx, tx.queryname)
	defer finish(span, time.Now())
	err := tx.tx.Commit()
	return err
}

// TxxInstrument struct
type TxxInstrument struct {
	tx        *sqlx.Tx
	ctx       context.Context
	queryname string
}

// InstrumentTxx return TxInstrument
func InstrumentTxx(ctx context.Context, tx *sqlx.Tx, queryname string) *TxxInstrument {
	txix := TxxInstrument{
		tx:        tx,
		ctx:       ctx,
		queryname: queryname,
	}
	return &txix
}

// Query queryinstrument
func (txx *TxxInstrument) Query(query string, args ...interface{}) (*sql.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	rows, err := txx.tx.QueryContext(ctx, query, args...)
	return rows, err
}

// Queryx queryinstrument
func (txx *TxxInstrument) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	span, ctx := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	rows, err := txx.tx.QueryxContext(ctx, query, args...)
	return rows, err
}

// QueryRow queryinstrument
func (txx *TxxInstrument) QueryRow(query string, args ...interface{}) *sql.Row {
	span, ctx := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	row := txx.tx.QueryRowContext(ctx, query, args...)
	return row
}

// QueryRowx queryinstrument
func (txx *TxxInstrument) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	span, ctx := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	row := txx.tx.QueryRowxContext(ctx, query, args...)
	return row
}

// Get queryinstrument
func (txx *TxxInstrument) Get(dest interface{}, query string, args ...interface{}) error {
	span, ctx := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	err := txx.tx.GetContext(ctx, dest, query, args...)
	return err
}

// Select queryinstruyment
func (txx *TxxInstrument) Select(dest interface{}, query string, args ...interface{}) error {
	span, ctx := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	err := txx.tx.SelectContext(ctx, dest, query, args...)
	return err
}

// Exec queryinstrument
func (txx *TxxInstrument) Exec(query string, args ...interface{}) (sql.Result, error) {
	span, ctx := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	result, err := txx.tx.ExecContext(ctx, query, args...)
	return result, err
}

// Commit queryinstrument
func (txx *TxxInstrument) Commit() error {
	span, _ := tracing.StartSpanFromContext(txx.ctx, txx.queryname)
	defer finish(span, time.Now())
	err := txx.tx.Commit()
	return err
}
