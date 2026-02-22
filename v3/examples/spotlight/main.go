package main

import (
	"log"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// This example demonstrates how to create a Spotlight-like launcher window
// that appears on all macOS Spaces and can overlay fullscreen applications.
//
// Key features:
// - Window appears on all Spaces (virtual desktops)
// - Can overlay fullscreen applications
// - Floating window level keeps it above other windows
// - Accessory activation policy hides from Dock
// - Frameless design with translucent backdrop

func main() {
	app := application.New(application.Options{
		Name:        "Spotlight Example",
		Description: "A Spotlight-like launcher demonstrating CollectionBehavior",
		Mac: application.MacOptions{
			// Accessory apps don't appear in the Dock
			ActivationPolicy: application.ActivationPolicyAccessory,
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(spotlightHTML))
			}),
		},
	})

	// Create a Spotlight-like window
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Spotlight",
		Width:     680,
		Height:    80,
		Frameless: true,
		// Center the window
		InitialPosition: application.WindowCentered,
		// Prevent resizing
		DisableResize: true,
		Mac: application.MacWindow{
			// Combine multiple behaviors using bitwise OR:
			// - CanJoinAllSpaces: window appears on ALL Spaces (virtual desktops)
			// - FullScreenAuxiliary: window can overlay fullscreen applications
			CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
				application.MacWindowCollectionBehaviorFullScreenAuxiliary,
			// Float above other windows
			WindowLevel: application.MacWindowLevelFloating,
			// Translucent vibrancy effect
			Backdrop: application.MacBackdropTranslucent,
			// Hidden title bar for clean look
			TitleBar: application.MacTitleBar{
				AppearsTransparent: true,
				Hide:               true,
			},
		},
		URL: "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

const spotlightHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Spotlight</title>
    <script type="module" src="/wails/runtime.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", sans-serif;
            background: transparent;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
            padding: 16px;
            user-select: none;
            -webkit-user-select: none;
        }
        .search-container {
            width: 100%;
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 12px 16px;
            background: rgba(255, 255, 255, 0.1);
            border-radius: 12px;
        }
        .search-icon {
            width: 24px;
            height: 24px;
            opacity: 0.6;
        }
        .search-input {
            flex: 1;
            background: transparent;
            border: none;
            outline: none;
            font-size: 24px;
            color: white;
            font-weight: 300;
        }
        .search-input::placeholder {
            color: rgba(255, 255, 255, 0.5);
        }
        @media (prefers-color-scheme: light) {
            .search-container {
                background: rgba(0, 0, 0, 0.05);
            }
            .search-input {
                color: #333;
            }
            .search-input::placeholder {
                color: rgba(0, 0, 0, 0.4);
            }
        }
    </style>
</head>
<body>
    <div class="search-container">
        <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <path d="M21 21l-4.35-4.35"/>
        </svg>
        <input type="text" class="search-input" placeholder="Spotlight Search" autofocus>
    </div>
</body>
</html>`
