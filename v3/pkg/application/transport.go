package application

import (
	"context"
	"net/http"
)

// Transport defines the interface for custom IPC transport implementations.
// Developers can provide their own transport (e.g., WebSocket, custom protocol)
// while retaining all Wails generated bindings and event communication.
//
// The transport is responsible for:
//   - Receiving runtime call requests from the frontend
//   - Processing them through Wails' MessageProcessor
//   - Sending responses back to the frontend
//
// Example use case: Implementing WebSocket-based transport instead of HTTP fetch.
type Transport interface {
	// Start initializes and starts the transport layer.
	// The provided handler should be called to process Wails runtime requests.
	// The context is the application context and will be cancelled on shutdown.
	Start(ctx context.Context, messageProcessor *MessageProcessor) error

	// Stop gracefully shuts down the transport.
	Stop() error
}

// AssetServerTransport is an optional interface that transports can implement
// to serve assets over HTTP, enabling browser-based deployments.
//
// When a transport implements this interface, Wails will call ServeAssets()
// after Start() to provide the asset server handler. The transport should
// integrate this handler into its HTTP server to serve HTML, CSS, JS, and
// other static assets alongside the IPC transport.
//
// This is useful for:
//   - Running Wails apps in a browser instead of a webview
//   - Exposing the app over a network
//   - Custom server configurations with both assets and IPC
type AssetServerTransport interface {
	Transport

	// ServeAssets configures the transport to serve assets.
	// The assetHandler is Wails' internal asset server that handles:
	//   - All static assets (HTML, CSS, JS, images, etc.)
	//   - /wails/runtime.js - The Wails runtime library
	//
	// The transport should integrate this handler into its HTTP server.
	// Typically this means mounting it at "/" and ensuring the IPC endpoint
	// (e.g., /wails/ws for WebSocket) is handled separately.
	//
	// This method is called after Start() completes successfully.
	ServeAssets(assetHandler http.Handler) error
}

// TransportHTTPHandler is an optional interface that transports can implement
// to provide HTTP middleware for the Wails asset server in webview scenarios.
//
// When a transport implements this interface, Wails will use Handler() in
// asset server middlewares that may provide handling for request done from webview to wails:// URLs.
//
// This is used by the default HTTP transport to handle IPC endpoints.
type TransportHTTPHandler interface {
	Handler() func(next http.Handler) http.Handler
}
