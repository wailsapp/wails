# Wails v3 CORS Example

This example demonstrates how to use Wails v3 with external URLs and proper CORS (Cross-Origin Resource Sharing) configuration. It shows how a Wails application can load its frontend from an external HTTPS server while maintaining secure communication with the Wails backend.

## 🎯 What This Example Demonstrates

- Loading a Wails WebView from an external URL (`https://app-local.wails-awesome.io:3000`)
- Configuring CORS to allow cross-origin communication between the external frontend and Wails backend
- Using Go to generate SSL certificates for cross-platform compatibility
- Secure HTTPS setup with self-signed certificates
- Making runtime calls from an external origin to the Wails backend

## 📋 Prerequisites

- Go 1.21 or higher
- Wails v3 CLI (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- Task (taskfile) - [Install Task](https://taskfile.dev/installation/)
- Administrator/root access (for trusting certificates and modifying hosts file)

## 🚀 Quick Start

### 1. Clone and Navigate to Example

```bash
cd v3/examples/cors
```

### 2. Run Automated Setup

```bash
# Run the complete setup
task setup
```

This will:
- Generate SSL certificates
- Check your hosts file
- Display setup instructions

### 3. Trust the CA Certificate

Based on your operating system:

**Windows:**
```bash
task trust-cert-windows
```

**macOS:**
```bash
task trust-cert-macos
```

**Linux:**
```bash
task trust-cert-linux
```

### 4. Add Hosts Entry

Add the following entry to your hosts file:

```bash
# Automated (requires admin/sudo)
task add-hosts-entry

# Or manually add to hosts file:
127.0.0.1    app-local.wails-awesome.io
```

Hosts file locations:
- Windows: `C:\Windows\System32\drivers\etc\hosts`
- macOS/Linux: `/etc/hosts`

### 5. Run the Example

You'll need two terminal windows:

**Terminal 1 - Start the external HTTPS server:**
```bash
task run-server
```

**Terminal 2 - Start the Wails application:**
```bash
task run-app
```

## 📁 Project Structure

```
cors/
├── main.go                 # Wails application with CORS configuration
├── external_server.go      # External HTTPS server
├── generate_certs.go       # Cross-platform certificate generator
├── Taskfile.yml           # Task automation
├── frontend/              # Frontend files served by external server
│   ├── index.html        # Main HTML file
│   └── app.js           # JavaScript application
├── certs/                # Generated certificates (git-ignored)
│   ├── ca.crt           # CA certificate (trust this)
│   ├── ca.key           # CA private key
│   ├── server.crt       # Server certificate
│   └── server.key       # Server private key
├── go.mod               # Go module file
└── README.md            # This file
```

## 🔧 How It Works

### 1. Certificate Generation

The `generate_certs.go` file creates:
- A self-signed Certificate Authority (CA)
- A server certificate signed by the CA for `app-local.wails-awesome.io`

This is done in pure Go for cross-platform compatibility.

### 2. External HTTPS Server

The `external_server.go` file:
- Serves the frontend files over HTTPS
- Uses the generated certificates
- Listens on port 3000
- Logs all incoming requests for debugging

### 3. Wails Application

The `main.go` file:
- Configures CORS to allow the external origin
- Creates a WebView window pointing to the external URL
- Provides backend services accessible via the Wails runtime

### 4. CORS Configuration

```go
CORS: application.CORSConfig{
    Enabled: true,
    AllowedOrigins: []string{
        "https://app-local.wails-awesome.io:3000",
        "https://localhost:3000",
    },
    AllowedMethods: []string{"GET", "POST", "OPTIONS"},
    AllowedHeaders: []string{
        "Content-Type",
        "X-Wails-Window-ID",
        "X-Wails-Window-Name",
        "X-Wails-Client-ID",
    },
    MaxAge: 5 * time.Minute,
}
```

## 🎮 Available Tasks

Run `task --list` to see all available tasks:

- `task setup` - Complete setup process
- `task generate-certs` - Generate SSL certificates
- `task clean-certs` - Remove generated certificates
- `task regenerate-certs` - Clean and regenerate certificates
- `task trust-cert-[os]` - Trust the CA certificate (OS-specific)
- `task check-hosts` - Check if hosts file has required entry
- `task add-hosts-entry` - Add entry to hosts file
- `task run-server` - Run the external HTTPS server
- `task run-app` - Run the Wails application
- `task build` - Build the Wails application
- `task test-cors` - Test CORS with curl

## 🧪 Testing the Example

Once running, you should see:

1. **External Server Console:** Logs showing incoming requests
2. **Wails Application:** A window loading from `https://app-local.wails-awesome.io:3000`
3. **Frontend Interface:** Three test buttons:
   - **Greet:** Calls the backend Greet method with a name
   - **Get Time:** Retrieves the current server time
   - **Test CORS:** Tests the CORS configuration

The browser developer console will show:
- Current page origin
- CORS headers in network requests
- Successful/failed runtime calls

## 🔒 Security Considerations

### Development vs Production

This example uses self-signed certificates for development. In production:

1. Use proper SSL certificates from a trusted CA
2. Configure CORS with specific allowed origins (no wildcards)
3. Use HTTPS for all communication
4. Validate and sanitize all inputs

### Example Production Configuration

```go
CORS: application.CORSConfig{
    Enabled: true,
    AllowedOrigins: []string{
        "https://app.mycompany.com",
        "https://cdn.mycompany.com",
    },
    AllowedMethods: []string{"GET", "POST"},
    AllowedHeaders: []string{
        "Content-Type",
        "X-Wails-Window-ID",
        "X-Wails-Window-Name",
        "Authorization",
    },
    MaxAge: 24 * time.Hour,
}
```

## ❓ Troubleshooting

### Certificate Errors

If you see certificate warnings:
1. Make sure you've trusted the CA certificate: `task trust-cert-[os]`
2. Restart your browser after trusting the certificate
3. Check that the certificates were generated: `ls certs/`

### "Cannot reach this site"

1. Ensure the hosts entry is added: `task check-hosts`
2. Verify the external server is running: `task run-server`
3. Check firewall settings for port 3000

### CORS Errors

1. Check the browser console for specific CORS error messages
2. Verify the origin in the CORS configuration matches exactly
3. Ensure the external server is using HTTPS (not HTTP)
4. Check that the Wails runtime is properly loaded

### "Wails runtime not found"

1. Make sure you're running the page inside the Wails WebView
2. Check that the Wails application is running: `task run-app`
3. Verify the URL in the WebView matches the external server URL

## 📚 Learn More

- [Wails v3 Documentation](https://v3alpha.wails.io)
- [CORS Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [Go TLS Documentation](https://pkg.go.dev/crypto/tls)

## 📝 License

This example is part of the Wails project and is licensed under the MIT License.