//go:build !wails_native

package application

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

type frontendState struct {
	assets *assetserver.AssetServer
}

func configureFrontend(app *App, options Options) {
	messageProcessor := NewMessageProcessor(app.Logger)
	app.messageProcessor = messageProcessor
	app.bindings = NewBindings(options.MarshalError, options.BindAliases)

	if options.NativeOnly {
		return
	}

	transport := options.Transport
	if transport == nil {
		transport = NewHTTPTransport(HTTPTransportWithLogger(app.Logger))
	}
	if err := transport.Start(app.ctx, messageProcessor); err != nil {
		app.fatal("failed to start custom transport: %w", err)
	}
	app.OnShutdown(func() {
		if err := transport.Stop(); err != nil {
			app.error("failed to stop custom transport: %w", err)
		}
	})

	if eventTransport, ok := transport.(WailsEventListener); ok {
		app.wailsEventListeners = append(app.wailsEventListeners, eventTransport)
	} else {
		app.wailsEventListeners = append(app.wailsEventListeners, &EventIPCTransport{app: app})
	}

	middlewares := []assetserver.Middleware{
		func(next http.Handler) http.Handler {
			if middleware := options.Assets.Middleware; middleware != nil {
				return middleware(next)
			}
			return next
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				path := request.URL.Path
				if strings.HasPrefix(path, eventPayloadPath) {
					app.serveEventPayload(response, request)
					return
				}
				if strings.HasPrefix(path, streamPath) {
					app.serveStream(response, request)
					return
				}
				switch path {
				case "/wails/runtime.js":
					if err := assetserver.ServeFile(response, path, runtimeJSWithPrelude()); err != nil {
						app.fatal("unable to serve runtime.js: %w", err)
					}
				case "/wails/transport.js":
					if err := assetserver.ServeFile(response, path, transport.JSClient()); err != nil {
						app.fatal("unable to serve transport.js: %w", err)
					}
				case "/wails/custom.js":
					http.NotFound(response, request)
				default:
					next.ServeHTTP(response, request)
				}
			})
		},
	}
	if handler, ok := transport.(TransportHTTPHandler); ok {
		middlewares = append(middlewares, handler.Handler())
	}

	assetOptions := &assetserver.Options{
		Handler:    options.Assets.Handler,
		Middleware: assetserver.ChainMiddleware(middlewares...),
		Logger:     app.Logger,
	}
	if options.Assets.DisableLogging {
		assetOptions.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	server, err := assetserver.NewAssetServer(assetOptions)
	if err != nil {
		app.fatal("application initialisation failed: %w", err)
	}
	app.assets = server
	app.assets.LogDetails()

	if assetTransport, ok := transport.(AssetServerTransport); ok {
		if err := assetTransport.ServeAssets(server); err != nil {
			app.fatal("failed to configure transport for serving assets: %w", err)
		}
		app.debug("Transport configured to serve assets")
	}
}

func startFrontendEventLoops(app *App) {
	if app.options.NativeOnly {
		return
	}
	go func() {
		for event := range windowEvents {
			go app.handleWindowEvent(event)
		}
	}()
	go func() {
		for request := range webviewRequests {
			go app.handleWebViewRequest(request)
		}
	}()
	go func() {
		for event := range windowMessageBuffer {
			go app.handleWindowMessage(event)
		}
	}()
	go func() {
		for event := range windowKeyEvents {
			go app.handleWindowKeyEvent(event)
		}
	}()
	go func() {
		for message := range windowDragAndDropBuffer {
			go app.handleDragAndDropMessage(message)
		}
	}()
	go func() {
		for event := range splitPaneLoadedEvents {
			handleMacSplitPaneLoaded(event)
		}
	}()
}

func bindFrontendService(app *App, service Service) error {
	return app.bindings.Add(service)
}

func attachServiceRoute(app *App, service Service) error {
	if service.options.Route == "" {
		return nil
	}
	handler, ok := service.Instance().(http.Handler)
	if !ok {
		handler = http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
			http.Error(response, fmt.Sprintf("Service '%s' does not handle HTTP requests", getServiceName(service)), http.StatusServiceUnavailable)
		})
	}
	app.assets.AttachServiceHandler(service.options.Route, handler)
	return nil
}
