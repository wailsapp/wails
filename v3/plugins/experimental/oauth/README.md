# oauth Plugin

This plugin provides the ability to initiate an OAuth authentication flow with a wide range of OAuth providers:

  - Amazon
  - Apple
  - Auth0
  - AzureAD
  - BattleNet
  - Bitbucket
  - Box
  - Dailymotion
  - Deezer
  - DigitalOcean
  - Discord
  - Dropbox
  - EveOnline
  - Facebook
  - Fitbit
  - Gitea
  - Gitlab
  - Github
  - Google
  - GooglePlus
  - Heroku
  - Intercom
  - Instagram
  - Kakao
  - LastFM
  - LinkedIn
  - Line
  - Mastodon
  - Meetup
  - MicrosoftOnline
  - Naver
  - NextCloud
  - Okta
  - Onedrive
  - OpenIDConnect
  - Patreon
  - PayPal
  - SalesForce
  - SeaTalk
  - Shopify
  - Slack
  - SoundCloud
  - Spotify
  - Steam
  - Strava
  - Stripe
  - TikTok
  - Twitter
  - TwitterV2
  - Typetalk
  - Twitch
  - Uber
  - VK
  - WeCom
  - Wepay
  - Xero
  - Yahoo
  - Yammer
  - Yandex
  - Zoom

## Installation

Add the plugin to the `Plugins` option in the Applications options. This example we are using the github provider:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/plugins/browser"
)

func main() {
    oAuthPlugin := oauth.NewPlugin(oauth.Config{
        Providers: []goth.Provider{
            github.New(
                os.Getenv("clientkey"),
                os.Getenv("secret"),
                "http://localhost:9876/auth/github/callback",
                "email",
                "profile"),
        },
    })

    app := application.New(application.Options{
    // ...
    Plugins: map[string]application.Plugin{
        "oauth": oAuthPlugin,
    },
})
```

### Configuration

The plugin takes a `Config` struct as a parameter. This struct has the following fields:

```go
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
```

If you don't specify a `WindowConfig`, the plugin will use the default window configuration:

```go
&application.WebviewWindowOptions{
    Title:  "OAuth Login",
    Width:  600,
    Height: 850,
    Hidden: true,
}
```

## Usage

### Go

You can start the flow by calling one of the provider methods:

```go
err := oAuthPlugin.Github()
```

In this example, we send an event from the frontend to start the process, so we listen for the event in the backend:

```go
app.Events.On("github-login", func(e *application.WailsEvent) {
    err := oAuthPlugin.Github()
    if err != nil {
        // process error
    }
})
```

### JavaScript

You can start the flow by calling one of the provider methods:

```javascript
await wails.Plugin("oauth","Github")
```

### Handling Success & Failure

When the OAuth flow completes, the plugin will send one of 2 events:

    - `wails:oauth:success` - The OAuth flow completed successfully. The event data will contain the user information.
    - `wails:oauth:error` - The OAuth flow failed. The event data will contain the error message.

In Javascript, we can listen for these events like so:

```javascript
window.wails.Events.On("wails:oauth:success", (event) => {
    document.getElementById("main").style.display = "none";
    document.getElementById("name").innerText = event.data.Name;
    document.getElementById("logo").src = event.data.AvatarURL;
    document.body.style.backgroundColor = "#000";
    document.body.style.color = "#FFF";
});
```

If you want to handle them in Go, you can do so like this:

```go
app.Events.On("wails:oauth:success", func(e *application.WailsEvent) {
    // Do something with the user data
})
```

Both these events are constants in the plugin:

```go
const (
    Success = "wails:oauth:success"
    Error   = "wails:oauth:error"
)
```

There is a working example of GitHub auth in the `v3/examples/oauth` directory.

## Logging Out

To log out, you can call the relevant `Logout` method for the provider:

```go
    err := oAuthPlugin.GithubLogout()
```

On success, the plugin will send a `wails:oauth:loggedout` event. On failure, it will send a `wails:oauth:error` event.

## Support

If you find a bug in this plugin, please raise a ticket on the Wails [Issue Tracker](https://github.com/wailsapp/wails/issues). 
