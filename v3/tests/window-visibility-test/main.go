package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// WindowTestService provides methods for testing window visibility scenarios
type WindowTestService struct {
	app *application.App
}

// NewWindowTestService creates a new window test service
func NewWindowTestService() *WindowTestService {
	return &WindowTestService{}
}

// SetApp sets the application reference (internal method, not exposed to frontend)
func (w *WindowTestService) setApp(app *application.App) {
	w.app = app
}

// CreateNormalWindow creates a standard window - should show immediately
func (w *WindowTestService) CreateNormalWindow() string {
	log.Println("Creating normal window...")

	w.app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Normal Window - Should Show Immediately",
		Width:  600,
		Height: 400,
		X:      100,
		Y:      100,
		HTML:   "<html><head><title>Normal Window</title><style>body{font-family:Arial,sans-serif;padding:20px;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:white;}</style></head><body><h1>✅ Normal Window</h1><p>This window should have appeared immediately after clicking the button.</p><p>Timestamp: " + time.Now().Format("15:04:05") + "</p></body></html>",
	})

	return "Normal window created"
}

// CreateDelayedContentWindow creates a window with delayed content to test navigation timing
func (w *WindowTestService) CreateDelayedContentWindow() string {
	log.Println("Creating delayed content window...")

	// Use HTML that will take time to load (simulates heavy Vue app)
	delayedHTML := `
	<html>
	<head>
		<title>Delayed Content</title>
		<style>
			body { font-family: Arial, sans-serif; padding: 20px; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; }
			.spinner { border: 4px solid #f3f3f3; border-top: 4px solid #3498db; border-radius: 50%; width: 40px; height: 40px; animation: spin 2s linear infinite; margin: 20px auto; }
			@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
		</style>
	</head>
	<body>
		<h1>⏳ Delayed Content Window</h1>
		<p>This window tests navigation completion timing.</p>
		<div class="spinner"></div>
		<p>Loading... (simulates heavy content)</p>
		<script>
			// Simulate slow loading content
			setTimeout(function() {
				document.querySelector('.spinner').style.display = 'none';
				document.body.innerHTML += '<h2>✅ Content Loaded!</h2><p>Navigation completed at: ' + new Date().toLocaleTimeString() + '</p>';
			}, 3000);
		</script>
		<p>Window container should be visible immediately, even during load.</p>
	</body>
	</html>`

	w.app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Delayed Content Window - Test Navigation Timing",
		Width:  600,
		Height: 400,
		X:      150,
		Y:      150,
		HTML:   delayedHTML,
	})

	return "Delayed content window created"
}

// CreateHiddenThenShowWindow creates a hidden window then shows it after delay
func (w *WindowTestService) CreateHiddenThenShowWindow() string {
	log.Println("Creating hidden then show window...")

	window := w.app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Hidden Then Show Window - Test Show() Robustness",
		Width:  600,
		Height: 400,
		X:      200,
		Y:      200,
		HTML:   "<html><head><title>Hidden Then Show</title><style>body{font-family:Arial,sans-serif;padding:20px;background:linear-gradient(135deg,#a8edea 0%,#fed6e3 100%);color:#333;}</style></head><body><h1>🔄 Hidden Then Show Window</h1><p>This window was created hidden and then shown after 2 seconds.</p><p>Should test the robustness of the show() method.</p><p>Created at: " + time.Now().Format("15:04:05") + "</p></body></html>",
		Hidden: true, // Start hidden
	})

	// Show after 2 seconds to test delayed showing
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("Showing previously hidden window...")
		window.Show()
	}()

	return "Hidden window created, will show in 2 seconds"
}

// CreateMultipleWindows creates multiple windows simultaneously to test performance
func (w *WindowTestService) CreateMultipleWindows() string {
	log.Println("Creating multiple windows...")

	for i := 0; i < 3; i++ {
		bgColors := []string{"#ff9a9e,#fecfef", "#a18cd1,#fbc2eb", "#fad0c4,#ffd1ff"}
		content := fmt.Sprintf(`
		<html>
		<head>
			<title>Batch Window %d</title>
			<style>
				body { font-family: Arial, sans-serif; padding: 20px; background: linear-gradient(135deg, %s); color: #333; text-align: center; }
			</style>
		</head>
		<body>
			<h1>🔢 Batch Window %d</h1>
			<p>Part of multiple windows stress test</p>
			<p>All windows should appear quickly and simultaneously</p>
			<p>Created at: %s</p>
		</body>
		</html>`, i+1, bgColors[i], i+1, time.Now().Format("15:04:05"))

		w.app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
			Title:  fmt.Sprintf("Batch Window %d - Stress Test", i+1),
			Width:  400,
			Height: 300,
			X:      250 + (i * 50),
			Y:      250 + (i * 50),
			HTML:   content,
		})
	}

	return "Created 3 windows simultaneously"
}

// CreateEfficiencyModeTestWindow creates a window designed to trigger efficiency mode issues
func (w *WindowTestService) CreateEfficiencyModeTestWindow() string {
	log.Println("Creating efficiency mode test window...")

	// Create content that might trigger efficiency mode or WebView2 delays
	heavyHTML := `
	<html>
	<head>
		<title>Efficiency Mode Test</title>
		<style>
			body { font-family: Arial, sans-serif; padding: 20px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; }
			.test-section { margin: 20px 0; padding: 15px; background: rgba(255,255,255,0.1); border-radius: 8px; }
			.heavy-content { height: 200px; overflow-y: scroll; background: rgba(0,0,0,0.2); padding: 10px; border-radius: 4px; }
		</style>
	</head>
	<body>
		<h1>⚡ Efficiency Mode Test Window</h1>
		<p>This window tests the fix for Windows 10 Pro efficiency mode issue #2861</p>
		
		<div class="test-section">
			<h3>Window Container Status</h3>
			<p id="container-status">✅ Window container is visible (this text proves it)</p>
		</div>
		
		<div class="test-section">
			<h3>WebView2 Status</h3>
			<p id="webview-status">⏳ WebView2 navigation in progress...</p>
			<p id="navigation-time"></p>
		</div>
		
		<div class="test-section">
			<h3>Heavy Content (simulates Vue.js app)</h3>
			<div class="heavy-content" id="heavy-content">
				Loading heavy content...
			</div>
		</div>
		
		<script>
			// Track navigation completion
			var startTime = performance.now();
			document.getElementById('navigation-time').textContent = 'Navigation started at: ' + new Date().toLocaleTimeString();
			
			// Simulate heavy JavaScript processing
			function heavyProcessing() {
				var content = '';
				for (var i = 0; i < 1000; i++) {
					content += 'Line ' + i + ': Simulated heavy content processing...<br>';
				}
				document.getElementById('heavy-content').innerHTML = content;
			}
			
			// Simulate WebView2 navigation completion with delay
			setTimeout(function() {
				document.getElementById('webview-status').innerHTML = '✅ WebView2 navigation completed successfully';
				var endTime = performance.now();
				document.getElementById('navigation-time').innerHTML += '<br>Navigation completed at: ' + new Date().toLocaleTimeString() + 
					'<br>Total time: ' + Math.round(endTime - startTime) + 'ms';
				heavyProcessing();
			}, 2000);
			
			console.log('Window created at:', new Date().toLocaleTimeString());
		</script>
	</body>
	</html>`

	w.app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Efficiency Mode Test - Issue #2861 Reproduction",
		Width:  700,
		Height: 500,
		X:      300,
		Y:      300,
		HTML:   heavyHTML,
	})

	return "Efficiency mode test window created"
}

// GetWindowCount returns the current number of windows
func (w *WindowTestService) GetWindowCount() int {
	// This would need to be implemented based on the app's window tracking
	// For now, return a placeholder
	return 1 // Main window
}

//go:embed assets/*
var assets embed.FS

func main() {
	// Create the service
	service := NewWindowTestService()

	// Create application with menu
	app := application.New(application.Options{
		Name:        "Window Visibility Test",
		Description: "Test application for window visibility robustness (Issue #2861)",
		Services: []application.Service{
			application.NewService(service),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	// Set app reference in service
	service.setApp(app)

	// Create application menu
	menu := app.NewMenu()

	// File menu
	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("New Normal Window").OnClick(func(ctx *application.Context) {
		service.CreateNormalWindow()
	})
	fileMenu.Add("New Delayed Content Window").OnClick(func(ctx *application.Context) {
		service.CreateDelayedContentWindow()
	})
	fileMenu.AddSeparator()
	fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	// Test menu
	testMenu := menu.AddSubmenu("Tests")
	testMenu.Add("Hidden Then Show Window").OnClick(func(ctx *application.Context) {
		service.CreateHiddenThenShowWindow()
	})
	testMenu.Add("Multiple Windows Stress Test").OnClick(func(ctx *application.Context) {
		service.CreateMultipleWindows()
	})
	testMenu.Add("Efficiency Mode Test").OnClick(func(ctx *application.Context) {
		service.CreateEfficiencyModeTestWindow()
	})

	// Help menu
	helpMenu := menu.AddSubmenu("Help")
	helpMenu.Add("About").OnClick(func(ctx *application.Context) {
		app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
			Title:  "About Window Visibility Test",
			Width:  500,
			Height: 400,
			X:      400,
			Y:      300,
			HTML:   "<html><head><title>About</title><style>body{font-family:Arial,sans-serif;padding:20px;background:#f0f0f0;color:#333;text-align:center;}</style></head><body><h1>Window Visibility Test</h1><p>This application tests the fixes for Wails v3 issue #2861</p><p><strong>Windows 10 Pro Efficiency Mode Fix</strong></p><p>Tests window container vs WebView content visibility</p><hr><p><em>Created for testing robust window visibility patterns</em></p></body></html>",
		})
	})

	app.SetMenu(menu)

	// Create main window
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Window Visibility Test - Issue #2861",
		Width:  800,
		Height: 600,
		X:      50,
		Y:      50,
		URL:    "/index.html",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
