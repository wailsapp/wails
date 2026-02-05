# Self-Update Example

This example demonstrates how to use the selfupdate service to implement automatic application updates in a Wails v3 application.

## Features

- Check for updates from GitHub Releases
- Support for pre-release versions
- Download progress tracking
- One-click update installation
- Dark mode support (follows system preference)
- Responsive UI with CSS variables for theming

## Configuration

Edit `main.go` to configure the update service for your repository:

```go
updateService := selfupdate.New(&selfupdate.Config{
    CurrentVersion: Version,
    Provider:       "github",
    GitHub: &selfupdate.GitHubConfig{
        Owner: "your-org",     // Your GitHub username or organization
        Repo:  "your-repo",    // Your repository name
        Token: os.Getenv("GITHUB_TOKEN"), // Optional: for private repos
    },
    AutoCheck: true, // Check for updates on startup
})
```

## Building with Version

Set the version at build time using ldflags:

```bash
go build -ldflags "-X main.Version=1.0.0" -o myapp
```

## Asset Naming

The selfupdate service expects release assets to follow a naming convention:

```
{repo}_{goos}_{goarch}{ext}
```

For example:
- `myapp_darwin_amd64.tar.gz`
- `myapp_darwin_arm64.tar.gz`
- `myapp_linux_amd64.tar.gz`
- `myapp_windows_amd64.zip`

You can customize this with the `AssetPattern` configuration option.

## Events

The service emits the following events that you can listen to in your frontend:

- `selfupdate:available` - Emitted when an update is available (with AutoCheck)
- `selfupdate:progress` - Emitted during download with progress information

## API

The service exposes these methods to the frontend:

- `Check()` - Check for updates
- `CheckWithPrerelease()` - Check for updates including pre-releases
- `Download()` - Download the available update
- `Install()` - Install the downloaded update
- `Restart()` - Restart the application
- `GetCurrentVersion()` - Get the current version
- `CanUpdate()` - Check if the app has permission to update itself

## Security

For production use, it's recommended to enable signature verification:

1. Generate a key pair:
   ```go
   pub, priv, _ := selfupdate.GenerateKeyPair()
   ```

2. Sign your release artifacts in CI:
   ```go
   signature, _ := selfupdate.SignData(binaryData, privateKey)
   ```

3. Configure the public key:
   ```go
   selfupdate.New(&selfupdate.Config{
       PublicKey: "your-base64-encoded-public-key",
       // ...
   })
   ```
