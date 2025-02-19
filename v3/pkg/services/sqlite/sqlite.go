//wails:include stmt.js
package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

type Config struct {
	// DBSource is the database URI to use.
	// The string ":memory:" can be used to create an in-memory database.
	// The sqlite driver can be configured through query parameters.
	// For more details see https://pkg.go.dev/modernc.org/sqlite#Driver.Open
	DBSource string
}

//wails:inject export {
//wails:inject     ExecContext as Execute,
//wails:inject     QueryContext as Query
//wails:inject };
//wails:inject
//wails:inject import { Stmt } from "./stmt.js";
//wails:inject
//wails:inject **:/**
//wails:inject **: * Prepare creates a prepared statement for later queries or executions.
//wails:inject **: * Multiple queries or executions may be run concurrently from the returned statement.
//wails:inject **: *
//wails:inject **: * The caller must call the statement's Close method when it is no longer needed.
//wails:inject **: * Statements are closed automatically
//wails:inject **: * when the connection they are associated with is closed.
//wails:inject **: *
//wails:inject **: * Prepare supports early cancellation.
//wails:inject j*: *
//wails:inject j*: * @param {string} query
//wails:inject j*: * @returns {Promise<Stmt | null> & { cancel(): void }}
//wails:inject **: */
//wails:inject j*:export function Prepare(query) {
//wails:inject t*:export function Prepare(query: string): Promise<Stmt | null> & { cancel(): void } {
//wails:inject **:    const promise = PrepareContext(query);
//wails:inject j*:    const wrapper = /** @type {any} */(promise.then(function (id) {
//wails:inject t*:    const wrapper: any = (promise.then(function (id) {
//wails:inject **:        return id == null ? null : new Stmt(
//wails:inject **:            ClosePrepared.bind(null, id),
//wails:inject **:            ExecPrepared.bind(null, id),
//wails:inject **:            QueryPrepared.bind(null, id));
//wails:inject **:    }));
//wails:inject **:    wrapper.cancel = promise.cancel;
//wails:inject **:    return wrapper;
//wails:inject **:}
type Service struct {
	lock   sync.RWMutex
	config *Config
	conn   *sql.DB
	stmts  map[uint64]struct{}
}

// New initialises a sqlite service with the default configuration.
func New() *Service {
	return NewWithConfig(nil)
}

// NewWithConfig initialises a sqlite service with a custom configuration.
// If config is nil, it falls back to the default configuration, i.e. an in-memory database.
//
// The database connection is not opened right away.
// A call to [Service.Open] must succeed before using all other methods.
// If the service is registered with the application,
// [Service.Open] will be called automatically at startup.
func NewWithConfig(config *Config) *Service {
	result := &Service{}
	result.Configure(config)
	return result
}

// ServiceName returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (s *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/plugins/sqlite"
}

// ServiceStartup opens the database connection.
// It returns a non-nil error in case of failures.
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return s.Open()
}

// ServiceShutdown closes the database connection.
// It returns a non-nil error in case of failures.
func (s *Service) ServiceShutdown() error {
	return s.Close()
}

// Configure changes the database service configuration.
// The connection state at call time is preserved.
// Consumers will need to call [Service.Open] manually after Configure
// in order to reconnect with the new configuration.
//
// See [NewWithConfig] for details on configuration.
//
//wails:ignore
func (s *Service) Configure(config *Config) {
	if config == nil {
		config = &Config{DBSource: ":memory:"}
	} else {
		// Clone to prevent changes from the outside.
		clone := new(Config)
		*clone = *config
		config = clone
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.config = config
}

// Open validates the current configuration,
// closes the current connection if one is present,
// then opens and validates a new connection.
//
// Even when a non-nil error is returned,
// the database service is left in a consistent state,
// ready for a new call to Open.
func (s *Service) Open() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.config.DBSource == "" {
		return errors.New(`no database source specified; please set DBSource in the config to a filename or specify ":memory:" to use an in-memory database`)
	}

	if err := s.closeImpl(); err != nil {
		return err
	}

	conn, err := sql.Open("sqlite", s.config.DBSource)
	if err != nil {
		return errors.Wrap(err, "error opening database connection")
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return errors.Wrap(err, "error opening database connection")
	}

	s.conn = conn
	s.stmts = make(map[uint64]struct{})

	return nil
}

// Close closes the current database connection if one is open, otherwise has no effect.
// Additionally, Close closes all open prepared statements associated to the connection.
//
// Even when a non-nil error is returned,
// the database service is left in a consistent state,
// ready for a call to [Service.Open].
func (s *Service) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.closeImpl()
}

// closeImpl performs the close operation without acquiring the lock first.
// It is the caller's responsibility
// to ensure the lock is held exclusively (in write mode)
// for the entire duration of the call.
func (s *Service) closeImpl() error {
	if s.conn == nil {
		return nil
	}

	for id := range s.stmts {
		if stmt, ok := stmts.Load(id); ok {
			// WARN: do not delegate to [Stmt.Close], it would cause a deadlock.
			// Ignore errors, closing the connection should free up all resources.
			_ = stmt.(*Stmt).sqlStmt.Close()
		}
	}

	err := s.conn.Close()

	// Clear the connection even in case of errors:
	// if [sql.DB.Close] returns an error,
	// the connection becomes unusable.
	s.conn = nil
	s.stmts = nil

	return err
}

// Execute executes a query without returning any rows.
//
//wails:ignore
func (s *Service) Execute(query string, args ...any) error {
	return s.ExecContext(context.Background(), query, args...)
}

// ExecContext executes a query without returning any rows.
// It supports early cancellation.
//
//wails:internal
func (s *Service) ExecContext(ctx context.Context, query string, args ...any) error {
	s.lock.RLock()
	conn := s.conn
	s.lock.RUnlock()

	if conn == nil {
		return errors.New("no open database connection")
	}

	_, err := conn.ExecContext(ctx, query, args...)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// Query executes a query and returns a slice of key-value records,
// one per row, with column names as keys.
//
//wails:ignore
func (s *Service) Query(query string, args ...any) (Rows, error) {
	return s.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query and returns a slice of key-value records,
// one per row, with column names as keys.
// It supports early cancellation, returning the slice of results fetched so far.
//
//wails:internal
func (s *Service) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	s.lock.RLock()
	conn := s.conn
	s.lock.RUnlock()

	if conn == nil {
		return nil, errors.New("no open database connection")
	}

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return Rows{}, nil
		} else {
			return nil, err
		}
	}

	return parseRows(ctx, rows)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
//
// The caller should call the statement's Close method when it is no longer needed.
// Statements are closed automatically
// when the connection they are associated with is closed.
//
//wails:ignore
func (s *Service) Prepare(query string) (*Stmt, error) {
	return s.PrepareContext(context.Background(), query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
//
// The caller must call the statement's Close method when it is no longer needed.
// Statements are closed automatically
// when the connection they are associated with is closed.
//
// PrepareContext supports early cancellation.
//
//wails:internal
func (s *Service) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	s.lock.RLock()
	conn := s.conn
	s.lock.RUnlock()

	if conn == nil {
		return nil, errors.New("no open database connection")
	}

	id := nextId.Load()
	for id != 0 && !nextId.CompareAndSwap(id, id+1) {
	}
	if id == 0 {
		return nil, errors.New("prepared statement ids exhausted")
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.stmts[id] = struct{}{}
	}()

	wrapper := &Stmt{
		sqlStmt: stmt,
		db:      s,
		id:      id,
	}
	stmts.Store(id, wrapper)

	return wrapper, nil
}

// ClosePrepared closes a prepared statement
// obtained with [Service.Prepare] or [Service.PrepareContext].
// ClosePrepared is idempotent:
// it has no effect on prepared statements that are already closed.
//
//wails:internal
func (s *Service) ClosePrepared(stmt *Stmt) error {
	return stmt.Close()
}

// ExecPrepared executes a prepared statement
// obtained with [Service.Prepare] or [Service.PrepareContext]
// without returning any rows.
// It supports early cancellation.
//
//wails:internal
func (s *Service) ExecPrepared(ctx context.Context, stmt *Stmt, args ...any) error {
	if stmt == nil {
		return errors.New("no prepared statement provided")
	} else if stmt.sqlStmt == nil {
		return errors.New("prepared statement is not valid")
	}

	_, err := stmt.ExecContext(ctx, args...)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// QueryPrepared executes a prepared statement
// obtained with [Service.Prepare] or [Service.PrepareContext]
// and returns a slice of key-value records, one per row, with column names as keys.
// It supports early cancellation, returning the slice of results fetched so far.
//
//wails:internal
func (s *Service) QueryPrepared(ctx context.Context, stmt *Stmt, args ...any) (Rows, error) {
	if stmt == nil {
		return nil, errors.New("no prepared statement provided")
	} else if stmt.sqlStmt == nil {
		return nil, errors.New("prepared statement is not valid")
	}

	rows, err := stmt.sqlStmt.QueryContext(ctx, args...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return Rows{}, nil
		} else {
			return nil, err
		}
	}

	return parseRows(ctx, rows)
}

type (
	// Row holds a single row in the result of a query.
	// It is a key-value map where keys are column names.
	Row = map[string]any

	// Rows holds the result of a query
	// as an array of key-value maps where keys are column names.
	Rows = []Row
)

func parseRows(ctx context.Context, rows *sql.Rows) (Rows, error) {
	defer rows.Close()

	columns, _ := rows.Columns()
	values := make([]any, len(columns))
	pointers := make([]any, len(columns))
	results := []map[string]any{}

	for rows.Next() {
		select {
		default:
		case <-ctx.Done():
			return results, nil
		}

		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(columns))
		for i, column := range columns {
			row[column] = values[i]
		}

		results = append(results, row)
	}

	return results, nil
}

var (
	// stmts holds all currently active prepared statements,
	// for all [Service] instances.
	stmts sync.Map

	// nextId holds the next available prepared statement id.
	// We use a counter to make sure IDs are never reused.
	nextId atomic.Uint64
)

func init() {
	nextId.Store(1)
}

type (
	sqlStmt = *sql.Stmt

	// Stmt wraps a prepared sql statement pointer.
	// It provides the same methods as the [sql.Stmt] type.
	//
	//wails:internal
	Stmt struct {
		sqlStmt
		db *Service
		id uint64
	}
)

// Close closes the statement.
// It has no effect when the statement is already closed.
func (s *Stmt) Close() error {
	if s == nil || s.sqlStmt == nil {
		return nil
	}

	err := s.sqlStmt.Close()
	stmts.Delete(s.id)

	func() {
		s.db.lock.Lock()
		defer s.db.lock.Unlock()

		delete(s.db.stmts, s.id)
	}()

	return errors.Wrap(err, "error closing prepared statement")
}

func (s *Stmt) MarshalText() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(16)

	if _, err := fmt.Fprintf(&buf, "%016x", s.id); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Stmt) UnmarshalText(text []byte) error {
	if n, err := fmt.Fscanf(bytes.NewReader(text), "%x", &s.id); n < 1 || err != nil {
		return errors.New("invalid prepared statement id")
	}

	if stmt, ok := stmts.Load(s.id); ok {
		*s = *(stmt.(*Stmt))
	} else {
		s.sqlStmt = nil
		s.db = nil
	}

	return nil
}
