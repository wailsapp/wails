package oauth

import (
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/wailsapp/wails/v3/pkg/application"
	"net/http"
	"time"
)

// ---------------- Plugin Setup ----------------
// This is the main plugin struct. It can be named anything you like.
// It must implement the application.Plugin interface.
// Both the Init() and Shutdown() methods are called synchronously when the app starts and stops.

const (
	Success = "wails:oauth:success"
	Error   = "wails:oauth:error"
)

type Plugin struct {
	config Config
	server *http.Server
	router *pat.Router
}

type Config struct {

	// Address to bind the temporary webserver to
	// Defaults to localhost:9876
	Address string

	// SessionSecret is the secret used to encrypt the session store.
	SessionSecret string

	// MaxAge is the maximum age of the session in seconds.
	MaxAge int

	// Providers is a list of goth providers to use.
	Providers []goth.Provider
}

func NewPlugin(config Config) *Plugin {
	result := &Plugin{
		config: config,
	}
	if result.config.MaxAge == 0 {
		result.config.MaxAge = 86400 * 30 // 30 days
	}
	if result.config.Address == "" {
		result.config.Address = "localhost:9876"
	}
	return result
}

func (p *Plugin) Shutdown() {
	if p.server != nil {
		p.server.Close()
	}
}

func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/oauth"
}

func (p *Plugin) Init(_ *application.App) error {

	store := sessions.NewCookieStore([]byte(p.config.SessionSecret))
	store.MaxAge(p.config.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

	gothic.Store = store
	goth.UseProviders(p.config.Providers...)

	p.router = pat.New()
	p.router.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		println("Callback...")
		event := &application.WailsEvent{
			Name:   Success,
			Sender: "",
		}
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			event.Data = err.Error()
			event.Name = Error
		} else {
			event.Data = user
		}
		application.Get().Events.Emit(event)
		err = p.server.Close()
		if err != nil {
			return
		}
		p.server = nil
	})

	p.router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		println("Authenticating...")
		gothic.BeginAuthHandler(res, req)
	})

	return nil
}

func (p *Plugin) CallableByJS() []string {
	return []string{
		"Start",
	}
}

func (p *Plugin) InjectJS() string {
	return ""
}

// ---------------- Plugin Methods ----------------

func (p *Plugin) Start() {
	if p.server != nil {
		println("Already listening")
		return
	}
	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: p.router,
	}
	println("Starting server")

	go p.server.ListenAndServe()
	time.Sleep(1 * time.Second)
}
