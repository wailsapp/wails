package application

import "context"

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

	JSClient() []byte

	// Stop gracefully shuts down the transport.
	Stop() error
}
