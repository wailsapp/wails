package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "WebviewPanel Demo",
		Description: "Demonstrates embedding multiple webview panels within a single window",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create the main window with a simple header
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "WebviewPanel Demo - Multi-Panel Layout",
		Width:            1200,
		Height:           700,
		BackgroundType:   application.BackgroundTypeSolid,
		BackgroundColour: application.NewRGB(45, 45, 45),
		HTML: `<!DOCTYPE html>
<html>
<head>
	<title>WebviewPanel Demo</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
			background: #2d2d2d;
			color: #fff;
		}
		.header {
			background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
			padding: 12px 20px;
			display: flex;
			align-items: center;
			gap: 15px;
			height: 50px;
			border-bottom: 1px solid #0f3460;
		}
		.header h1 {
			font-size: 15px;
			font-weight: 500;
			color: #e94560;
		}
		.header .subtitle {
			font-size: 12px;
			color: #888;
		}
		.info {
			margin-left: auto;
			font-size: 11px;
			color: #666;
			padding: 4px 10px;
			background: rgba(0,0,0,0.3);
			border-radius: 4px;
		}
	</style>
</head>
<body>
	<div class="header">
		<h1>🖥️ WebviewPanel Demo</h1>
		<span class="subtitle">Multiple independent webviews in one window</span>
		<span class="info">Panels render below this header area</span>
	</div>
</body>
</html>`,
	})

	// =====================================================================
	// Example 1: Using explicit coordinates (traditional approach)
	// =====================================================================
	
	// Create a sidebar panel on the left with explicit positioning
	sidebarPanel := window.NewPanel(application.WebviewPanelOptions{
		Name:   "sidebar",
		X:      0,
		Y:      50,    // Start below the 50px header
		Width:  220,
		Height: 650,
		HTML: `<!DOCTYPE html>
<html>
<head>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
			background: linear-gradient(180deg, #1e1e2e 0%, #181825 100%);
			color: #cdd6f4;
			padding: 15px;
			height: 100vh;
		}
		h3 {
			font-size: 10px;
			color: #6c7086;
			text-transform: uppercase;
			letter-spacing: 1px;
			margin-bottom: 12px;
			padding-bottom: 8px;
			border-bottom: 1px solid #313244;
		}
		ul { list-style: none; }
		li {
			padding: 10px 12px;
			margin: 4px 0;
			border-radius: 8px;
			cursor: pointer;
			font-size: 13px;
			transition: all 0.2s;
			display: flex;
			align-items: center;
			gap: 10px;
		}
		li:hover { background: #313244; }
		li.active {
			background: linear-gradient(135deg, #89b4fa 0%, #74c7ec 100%);
			color: #1e1e2e;
			font-weight: 500;
		}
		.section { margin-bottom: 25px; }
		.icon { font-size: 16px; }
	</style>
</head>
<body>
	<div class="section">
		<h3>Navigation</h3>
		<ul>
			<li class="active"><span class="icon">🏠</span> Dashboard</li>
			<li><span class="icon">📊</span> Analytics</li>
			<li><span class="icon">📁</span> Projects</li>
			<li><span class="icon">👥</span> Team</li>
			<li><span class="icon">⚙️</span> Settings</li>
		</ul>
	</div>
	<div class="section">
		<h3>Favorites</h3>
		<ul>
			<li><span class="icon">⭐</span> Starred</li>
			<li><span class="icon">🕐</span> Recent</li>
			<li><span class="icon">📌</span> Pinned</li>
		</ul>
	</div>
</body>
</html>`,
		BackgroundColour: application.NewRGB(30, 30, 46),
		Visible:          boolPtr(true),
		ZIndex:           1,
	})

	// =====================================================================
	// Example 2: Content panel showing an external website
	// This demonstrates loading external URLs in an embedded webview
	// =====================================================================
	
	contentPanel := window.NewPanel(application.WebviewPanelOptions{
		Name:             "content",
		X:                220, // Right of sidebar
		Y:                50,  // Below header
		Width:            980,
		Height:           650,
		URL:              "https://wails.io", // External website
		DevToolsEnabled:  boolPtr(true),
		Visible:          boolPtr(true),
		BackgroundColour: application.NewRGB(255, 255, 255),
		ZIndex:           1,
	})

	// Log panel creation
	log.Printf("✅ Created sidebar panel: %s (ID: %d)", sidebarPanel.Name(), sidebarPanel.ID())
	log.Printf("✅ Created content panel: %s (ID: %d)", contentPanel.Name(), contentPanel.ID())
	
	// =====================================================================
	// Alternative: Using layout helper methods (commented examples)
	// =====================================================================
	//
	// // Create a panel and dock it to the left
	// sidebar := window.NewPanel(opts).DockLeft(200)
	//
	// // Create a panel and fill space beside another
	// content := window.NewPanel(opts).FillBeside(sidebar, "right")
	//
	// // Create a panel that fills the entire window
	// fullscreen := window.NewPanel(opts).FillWindow()
	//
	// // Dock panels to different edges
	// toolbar := window.NewPanel(opts).DockTop(50)
	// statusBar := window.NewPanel(opts).DockBottom(30)
	//

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
