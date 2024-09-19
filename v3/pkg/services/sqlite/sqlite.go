package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

// ---------------- Service Setup ----------------
// This is the main Service struct. It can be named anything you like.

type Config struct {
	DBFile string
}

type Service struct {
	config *Config
	conn   *sql.DB
}

func New(config *Config) *Service {
	return &Service{
		config: config,
	}
}

// OnShutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (s *Service) OnShutdown() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// Name returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (s *Service) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/sqlite"
}

// OnStartup is called when the app is starting up. You can use this to
// initialise any resources you need.
func (s *Service) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	if s.config.DBFile == "" {
		return errors.New(`no database file specified. Please set DBFile in the config to either a filename or use ":memory:" to use an in-memory database`)
	}
	db, err := s.Open(s.config.DBFile)
	if err != nil {
		return err
	}
	_ = db

	return nil
}

func (s *Service) Open(dbPath string) (string, error) {
	var err error
	s.conn, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return "", err
	}
	return "Database connection opened", nil
}

func (s *Service) Execute(query string, args ...any) error {
	if s.conn == nil {
		return errors.New("no open database connection")
	}

	_, err := s.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Select(query string, args ...any) ([]map[string]any, error) {
	if s.conn == nil {
		return nil, errors.New("no open database connection")
	}

	rows, err := s.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))

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

func (s *Service) Close() (string, error) {
	if s.conn == nil {
		return "", errors.New("no open database connection")
	}

	err := s.conn.Close()
	if err != nil {
		return "", err
	}
	s.conn = nil
	return "Database connection closed", nil
}
