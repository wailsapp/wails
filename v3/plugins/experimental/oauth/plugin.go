package oauth

import (
	"context"
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
	Success   = "wails:oauth:success"
	Error     = "wails:oauth:error"
	LoggedOut = "wails:oauth:loggedout"
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

func (p *Plugin) Shutdown() error {
	if p.server != nil {
		return p.server.Close()
	}
	return nil
}

func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/oauth"
}

func (p *Plugin) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	store := sessions.NewCookieStore([]byte(p.config.SessionSecret))
	store.MaxAge(p.config.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

	gothic.Store = store
	goth.UseProviders(p.config.Providers...)

	return nil
}

func (p *Plugin) CallableByJS() []string {
	return []string{
		"Github",
		"LogoutGithub",
		"Amazon",
		"LogoutAmazon",
		"Apple",
		"LogoutApple",
		"Auth0",
		"LogoutAuth0",
		"AzureAD",
		"LogoutAzureAD",
		"BattleNet",
		"LogoutBattleNet",
		"Bitbucket",
		"LogoutBitbucket",
		"Box",
		"LogoutBox",
		"Dailymotion",
		"LogoutDailymotion",
		"Deezer",
		"LogoutDeezer",
		"DigitalOcean",
		"LogoutDigitalOcean",
		"Discord",
		"LogoutDiscord",
		"Dropbox",
		"LogoutDropbox",
		"EveOnline",
		"LogoutEveOnline",
		"Facebook",
		"LogoutFacebook",
		"Fitbit",
		"LogoutFitbit",
		"Gitea",
		"LogoutGitea",
		"Gitlab",
		"LogoutGitlab",
		"Google",
		"LogoutGoogle",
		"GooglePlus",
		"LogoutGooglePlus",
		"Heroku",
		"LogoutHeroku",
		"Intercom",
		"LogoutIntercom",
		"Instagram",
		"LogoutInstagram",
		"Kakao",
		"LogoutKakao",
		"LastFM",
		"LogoutLastFM",
		"LinkedIn",
		"LogoutLinkedIn",
		"Line",
		"LogoutLine",
		"Mastodon",
		"LogoutMastodon",
		"Meetup",
		"LogoutMeetup",
		"MicrosoftOnline",
		"LogoutMicrosoftOnline",
		"Naver",
		"LogoutNaver",
		"NextCloud",
		"LogoutNextCloud",
		"Okta",
		"LogoutOkta",
		"Onedrive",
		"LogoutOnedrive",
		"OpenIDConnect",
		"LogoutOpenIDConnect",
		"Patreon",
		"LogoutPatreon",
		"PayPal",
		"LogoutPayPal",
		"SalesForce",
		"LogoutSalesForce",
		"SeaTalk",
		"LogoutSeaTalk",
		"Shopify",
		"LogoutShopify",
		"Slack",
		"LogoutSlack",
		"SoundCloud",
		"LogoutSoundCloud",
		"Spotify",
		"LogoutSpotify",
		"Steam",
		"LogoutSteam",
		"Strava",
		"LogoutStrava",
		"Stripe",
		"LogoutStripe",
		"TikTok",
		"LogoutTikTok",
		"Twitter",
		"LogoutTwitter",
		"TwitterV2",
		"LogoutTwitterV2",
		"Typetalk",
		"LogoutTypetalk",
		"Twitch",
		"LogoutTwitch",
		"Uber",
		"LogoutUber",
		"VK",
		"LogoutVK",
		"WeCom",
		"LogoutWeCom",
		"Wepay",
		"LogoutWepay",
		"Xero",
		"LogoutXero",
		"Yahoo",
		"LogoutYahoo",
		"Yammer",
		"LogoutYammer",
		"Yandex",
		"LogoutYandex",
		"Zoom",
		"LogoutZoom",
	}
}

func (p *Plugin) InjectJS() string {
	return ""
}

func (p *Plugin) start(provider string) error {
	if p.server != nil {
		return fmt.Errorf("server already processing request. Please wait for the current login to complete")
	}

	router := pat.New()
	router.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			application.Get().EmitEvent(Error, err.Error())
		} else {
			application.Get().EmitEvent(Success, user)
		}

		_ = p.server.Close()
		p.server = nil
	})

	router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: router,
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

	application.Get().OnEvent(Success, func(event *application.CustomEvent) {
		window.Close()
	})
	application.Get().OnEvent(Error, func(event *application.CustomEvent) {
		window.Close()
	})

	return nil
}

func (p *Plugin) logout(provider string) error {
	if p.server != nil {
		return fmt.Errorf("server already processing request. Please wait for the current operation to complete")
	}

	router := pat.New()
	router.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		err := gothic.Logout(res, req)
		if err != nil {
			application.Get().EmitEvent(Error, err.Error())
		} else {
			application.Get().EmitEvent(LoggedOut)
		}
		_ = p.server.Close()
		p.server = nil
	})

	p.server = &http.Server{
		Addr:    p.config.Address,
		Handler: router,
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
	p.config.WindowConfig.URL = "http://" + p.config.Address + "/logout/" + provider
	window := application.Get().NewWebviewWindowWithOptions(*p.config.WindowConfig)
	window.Show()

	application.Get().OnEvent(LoggedOut, func(event *application.CustomEvent) {
		window.Close()
	})
	application.Get().OnEvent(Error, func(event *application.CustomEvent) {
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

func (p *Plugin) LogoutAmazon() error {
	return p.logout("amazon")
}

func (p *Plugin) LogoutApple() error {
	return p.logout("apple")
}

func (p *Plugin) LogoutAuth0() error {
	return p.logout("auth0")
}

func (p *Plugin) LogoutAzureAD() error {
	return p.logout("azuread")
}

func (p *Plugin) LogoutBattleNet() error {
	return p.logout("battlenet")
}

func (p *Plugin) LogoutBitbucket() error {
	return p.logout("bitbucket")
}

func (p *Plugin) LogoutBox() error {
	return p.logout("box")
}

func (p *Plugin) LogoutDailymotion() error {
	return p.logout("dailymotion")
}

func (p *Plugin) LogoutDeezer() error {
	return p.logout("deezer")
}

func (p *Plugin) LogoutDigitalOcean() error {
	return p.logout("digitalocean")
}

func (p *Plugin) LogoutDiscord() error {
	return p.logout("discord")
}

func (p *Plugin) LogoutDropbox() error {
	return p.logout("dropbox")
}

func (p *Plugin) LogoutEveOnline() error {
	return p.logout("eveonline")
}

func (p *Plugin) LogoutFacebook() error {
	return p.logout("facebook")
}

func (p *Plugin) LogoutFitbit() error {
	return p.logout("fitbit")
}

func (p *Plugin) LogoutGitea() error {
	return p.logout("gitea")
}

func (p *Plugin) LogoutGitlab() error {
	return p.logout("gitlab")
}

func (p *Plugin) LogoutGithub() error {
	return p.logout("github")
}

func (p *Plugin) LogoutGoogle() error {
	return p.logout("google")
}

func (p *Plugin) LogoutGooglePlus() error {
	return p.logout("gplus")
}

func (p *Plugin) LogoutHeroku() error {
	return p.logout("heroku")
}

func (p *Plugin) LogoutIntercom() error {
	return p.logout("intercom")
}

func (p *Plugin) LogoutInstagram() error {
	return p.logout("instagram")
}

func (p *Plugin) LogoutKakao() error {
	return p.logout("kakao")
}

func (p *Plugin) LogoutLastFM() error {
	return p.logout("lastfm")
}

func (p *Plugin) LogoutLinkedIn() error {
	return p.logout("linkedin")
}

func (p *Plugin) LogoutLine() error {
	return p.logout("line")
}

func (p *Plugin) LogoutMastodon() error {
	return p.logout("mastodon")
}

func (p *Plugin) LogoutMeetup() error {
	return p.logout("meetup")
}

func (p *Plugin) LogoutMicrosoftOnline() error {
	return p.logout("microsoftonline")
}

func (p *Plugin) LogoutNaver() error {
	return p.logout("naver")
}

func (p *Plugin) LogoutNextCloud() error {
	return p.logout("nextcloud")
}

func (p *Plugin) LogoutOkta() error {
	return p.logout("okta")
}

func (p *Plugin) LogoutOnedrive() error {
	return p.logout("onedrive")
}

func (p *Plugin) LogoutOpenIDConnect() error {
	return p.logout("openid-connect")
}

func (p *Plugin) LogoutPatreon() error {
	return p.logout("patreon")
}

func (p *Plugin) LogoutPayPal() error {
	return p.logout("paypal")
}

func (p *Plugin) LogoutSalesForce() error {
	return p.logout("salesforce")
}

func (p *Plugin) LogoutSeaTalk() error {
	return p.logout("seatalk")
}

func (p *Plugin) LogoutShopify() error {
	return p.logout("shopify")
}

func (p *Plugin) LogoutSlack() error {
	return p.logout("slack")
}

func (p *Plugin) LogoutSoundCloud() error {
	return p.logout("soundcloud")
}

func (p *Plugin) LogoutSpotify() error {
	return p.logout("spotify")
}

func (p *Plugin) LogoutSteam() error {
	return p.logout("steam")
}

func (p *Plugin) LogoutStrava() error {
	return p.logout("strava")
}

func (p *Plugin) LogoutStripe() error {
	return p.logout("stripe")
}

func (p *Plugin) LogoutTikTok() error {
	return p.logout("tiktok")
}

func (p *Plugin) LogoutTwitter() error {
	return p.logout("twitter")
}

func (p *Plugin) LogoutTwitterV2() error {
	return p.logout("twitterv2")
}

func (p *Plugin) LogoutTypetalk() error {
	return p.logout("typetalk")
}

func (p *Plugin) LogoutTwitch() error {
	return p.logout("twitch")
}

func (p *Plugin) LogoutUber() error {
	return p.logout("uber")
}

func (p *Plugin) LogoutVK() error {
	return p.logout("vk")
}

func (p *Plugin) LogoutWeCom() error {
	return p.logout("wecom")
}

func (p *Plugin) LogoutWepay() error {
	return p.logout("wepay")
}

func (p *Plugin) LogoutXero() error {
	return p.logout("xero")
}

func (p *Plugin) LogoutYahoo() error {
	return p.logout("yahoo")
}

func (p *Plugin) LogoutYammer() error {
	return p.logout("yammer")
}

func (p *Plugin) LogoutYandex() error {
	return p.logout("yandex")
}

func (p *Plugin) LogoutZoom() error {
	return p.logout("zoom")
}
