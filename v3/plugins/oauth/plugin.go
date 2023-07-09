package oauth

import (
	"fmt"
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

	// WindowConfig is the configuration for the window that will be opened
	// to perform the OAuth login.
	WindowConfig *application.WebviewWindowOptions
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
	if result.config.WindowConfig == nil {
		result.config.WindowConfig = &application.WebviewWindowOptions{
			Title:  "OAuth Login",
			Width:  600,
			Height: 850,
			Hidden: true,
		}
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
		"Amazon",
		"Apple",
		"Auth0",
		"AzureAD",
		"BattleNet",
		"Bitbucket",
		"Box",
		"Dailymotion",
		"Deezer",
		"DigitalOcean",
		"Discord",
		"Dropbox",
		"EveOnline",
		"Facebook",
		"Fitbit",
		"Gitea",
		"Gitlab",
		"Github",
		"Google",
		"GooglePlus",
		"Heroku",
		"Intercom",
		"Instagram",
		"Kakao",
		"LastFM",
		"LinkedIn",
		"Line",
		"Mastodon",
		"Meetup",
		"MicrosoftOnline",
		"Naver",
		"NextCloud",
		"Okta",
		"Onedrive",
		"OpenIDConnect",
		"Patreon",
		"PayPal",
		"SalesForce",
		"SeaTalk",
		"Shopify",
		"Slack",
		"SoundCloud",
		"Spotify",
		"Steam",
		"Strava",
		"Stripe",
		"TikTok",
		"Twitter",
		"TwitterV2",
		"Typetalk",
		"Twitch",
		"Uber",
		"VK",
		"WeCom",
		"Wepay",
		"Xero",
		"Yahoo",
		"Yammer",
		"Yandex",
		"Zoom",
	}
}

func (p *Plugin) InjectJS() string {
	return ""
}

func (p *Plugin) start(provider string) error {
	if p.server != nil {
		return fmt.Errorf("server already processing request. Please wait for the current login to complete")
	}
	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: p.router,
	}

	go p.server.ListenAndServe()
	// Keep trying to connect until we succeed
	var keepTrying = true
	var connected = false

	go func() {
		time.Sleep(3 * time.Second)
		keepTrying = false
	}()

	for keepTrying {
		_, err := http.Get("http://" + p.config.Address)
		if err == nil {
			connected = true
			break
		}
	}

	if !connected {
		return fmt.Errorf("server failed to start")
	}

	// create a window
	p.config.WindowConfig.URL = "http://" + p.config.Address + "/auth/" + provider
	window := application.Get().NewWebviewWindowWithOptions(*p.config.WindowConfig)
	window.Show()

	application.Get().Events.On(Success, func(event *application.WailsEvent) {
		window.Close()
	})
	application.Get().Events.On(Error, func(event *application.WailsEvent) {
		window.Close()
	})

	return nil
}

// ---------------- Plugin Methods ----------------

func (p *Plugin) Amazon() error {
	return p.start("amazon")
}

func (p *Plugin) Apple() error {
	return p.start("apple")
}

func (p *Plugin) Auth0() error {
	return p.start("auth0")
}

func (p *Plugin) AzureAD() error {
	return p.start("azuread")
}

func (p *Plugin) BattleNet() error {
	return p.start("battlenet")
}

func (p *Plugin) Bitbucket() error {
	return p.start("bitbucket")
}

func (p *Plugin) Box() error {
	return p.start("box")
}

func (p *Plugin) Dailymotion() error {
	return p.start("dailymotion")
}

func (p *Plugin) Deezer() error {
	return p.start("deezer")
}

func (p *Plugin) DigitalOcean() error {
	return p.start("digitalocean")
}

func (p *Plugin) Discord() error {
	return p.start("discord")
}

func (p *Plugin) Dropbox() error {
	return p.start("dropbox")
}

func (p *Plugin) EveOnline() error {
	return p.start("eveonline")
}

func (p *Plugin) Facebook() error {
	return p.start("facebook")
}

func (p *Plugin) Fitbit() error {
	return p.start("fitbit")
}

func (p *Plugin) Gitea() error {
	return p.start("gitea")
}

func (p *Plugin) Gitlab() error {
	return p.start("gitlab")
}

func (p *Plugin) Github() error {
	return p.start("github")
}

func (p *Plugin) Google() error {
	return p.start("google")
}

func (p *Plugin) GooglePlus() error {
	return p.start("gplus")
}

func (p *Plugin) Heroku() error {
	return p.start("heroku")
}

func (p *Plugin) Intercom() error {
	return p.start("intercom")
}

func (p *Plugin) Instagram() error {
	return p.start("instagram")
}

func (p *Plugin) Kakao() error {
	return p.start("kakao")
}

func (p *Plugin) LastFM() error {
	return p.start("lastfm")
}

func (p *Plugin) LinkedIn() error {
	return p.start("linkedin")
}

func (p *Plugin) Line() error {
	return p.start("line")
}

func (p *Plugin) Mastodon() error {
	return p.start("mastodon")
}

func (p *Plugin) Meetup() error {
	return p.start("meetup")
}

func (p *Plugin) MicrosoftOnline() error {
	return p.start("microsoftonline")
}

func (p *Plugin) Naver() error {
	return p.start("naver")
}

func (p *Plugin) NextCloud() error {
	return p.start("nextcloud")
}

func (p *Plugin) Okta() error {
	return p.start("okta")
}

func (p *Plugin) Onedrive() error {
	return p.start("onedrive")
}

func (p *Plugin) OpenIDConnect() error {
	return p.start("openid-connect")
}

func (p *Plugin) Patreon() error {
	return p.start("patreon")
}

func (p *Plugin) PayPal() error {
	return p.start("paypal")
}

func (p *Plugin) SalesForce() error {
	return p.start("salesforce")
}

func (p *Plugin) SeaTalk() error {
	return p.start("seatalk")
}

func (p *Plugin) Shopify() error {
	return p.start("shopify")
}

func (p *Plugin) Slack() error {
	return p.start("slack")
}

func (p *Plugin) SoundCloud() error {
	return p.start("soundcloud")
}

func (p *Plugin) Spotify() error {
	return p.start("spotify")
}

func (p *Plugin) Steam() error {
	return p.start("steam")
}

func (p *Plugin) Strava() error {
	return p.start("strava")
}

func (p *Plugin) Stripe() error {
	return p.start("stripe")
}

func (p *Plugin) TikTok() error {
	return p.start("tiktok")
}

func (p *Plugin) Twitter() error {
	return p.start("twitter")
}

func (p *Plugin) TwitterV2() error {
	return p.start("twitterv2")
}

func (p *Plugin) Typetalk() error {
	return p.start("typetalk")
}

func (p *Plugin) Twitch() error {
	return p.start("twitch")
}

func (p *Plugin) Uber() error {
	return p.start("uber")
}

func (p *Plugin) VK() error {
	return p.start("vk")
}

func (p *Plugin) WeCom() error {
	return p.start("wecom")
}

func (p *Plugin) Wepay() error {
	return p.start("wepay")
}

func (p *Plugin) Xero() error {
	return p.start("xero")
}

func (p *Plugin) Yahoo() error {
	return p.start("yahoo")
}

func (p *Plugin) Yammer() error {
	return p.start("yammer")
}

func (p *Plugin) Yandex() error {
	return p.start("yandex")
}

func (p *Plugin) Zoom() error {
	return p.start("zoom")
}
