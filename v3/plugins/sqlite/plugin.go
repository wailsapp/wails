package sqlite

import (
	"database/sql"
	_ "embed"
	"errors"
	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

//go:embed plugin.js
var js string

// ---------------- Plugin Setup ----------------
// This is the main plugin struct. It can be named anything you like.
// It must implement the application.Plugin interface.
// Both the Init() and Shutdown() methods are called synchronously when the app starts and stops.

type Config struct {
	// Add any configuration options here
	CanOpenFromJS  bool
	CanCloseFromJS bool
	DBFile         string
}

type Plugin struct {
	config          *Config
	app             *application.App
	conn            *sql.DB
	callableMethods []string
}

func NewPlugin(config *Config) *Plugin {
	return &Plugin{
		config: config,
	}
}

// Shutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (p *Plugin) Shutdown() {
	if p.conn != nil {
		p.conn.Close()
	}
}

// Name returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/sqlite"
}

// Init is called when the app is starting up. You can use this to
// initialise any resources you need. You can also access the application
// instance via the app property.
func (p *Plugin) Init(app *application.App) error {
	p.app = app
	p.callableMethods = []string{"Execute", "Select"}
	if p.config.CanOpenFromJS {
		p.callableMethods = append(p.callableMethods, "Open")
	}
	if p.config.CanCloseFromJS {
		p.callableMethods = append(p.callableMethods, "Close")
	}
	if p.config.DBFile != "" {
		_, err := p.Open(p.config.DBFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// CallableByJS returns a list of exported methods that can be called from the frontend
func (p *Plugin) CallableByJS() []string {
	return p.callableMethods
}

func (p *Plugin) InjectJS() string {
	return js
}

// ---------------- Plugin Methods ----------------
// Plugin methods are just normal Go methods. You can add as many as you like.
// The only requirement is that they are exported (start with a capital letter).
// You can also return any type that is JSON serializable.
// Any methods that you want to be callable from the frontend must be returned by the
// Exported() method above.
// See https://golang.org/pkg/encoding/json/#Marshal for more information.

func (p *Plugin) Open(dbPath string) (string, error) {
	var err error
	p.conn, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return "", err
	}
	return "Database connection opened", nil
}

func (p *Plugin) Execute(query string, args ...any) error {
	if p.conn == nil {
		return errors.New("no open database connection")
	}

	_, err := p.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (p *Plugin) Select(query string, args ...any) ([]map[string]any, error) {
	if p.conn == nil {
		return nil, errors.New("no open database connection")
	}

	rows, err := p.conn.Query(query, args...)
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

func (p *Plugin) Close() error {
	if p.conn == nil {
		return errors.New("no open database connection")
	}

	err := p.conn.Close()
	if err != nil {
		return err
	}
	p.conn = nil
	return nil
}
