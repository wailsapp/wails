package main

import (
	"context"
	"fmt"

	"github.com/leaanthony/u"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Navigation Gestures Example",
		Description: "Example showing allowsBackForwardNavigationGestures feature",
	})

	// Create a new window with navigation gestures enabled
	window := app.NewWebviewWindow()
	window.SetOptions(application.WebviewWindowOptions{
		Title:  "Navigation Gestures Demo",
		Width:  800,
		Height: 600,
		Mac: application.MacWindow{
			WebviewPreferences: application.MacWebviewPreferences{
				// Enable horizontal swipe gestures for back/forward navigation
				AllowsBackForwardNavigationGestures: u.True(),
			},
		},
		HTML: `
<!DOCTYPE html>
<html>
<head>
    <title>Navigation Gestures Demo</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, sans-serif; 
            padding: 40px; 
            line-height: 1.6;
        }
        .page { margin: 20px 0; }
        a { 
            display: inline-block; 
            margin: 10px 15px 10px 0; 
            padding: 8px 16px; 
            background: #007AFF; 
            color: white; 
            text-decoration: none; 
            border-radius: 6px; 
        }
        a:hover { background: #0051D0; }
        .instruction {
            background: #f0f0f0;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <h1>Navigation Gestures Demo</h1>
    
    <div class="instruction">
        <strong>Mac users:</strong> Try using two-finger horizontal swipe gestures to navigate back and forward between these pages!
    </div>

    <div class="page">
        <h2>Page 1</h2>
        <p>This is the first page. Click the links below to navigate to other pages, then try swiping left/right with two fingers to go back and forward.</p>
        <a href="#page2">Go to Page 2</a>
        <a href="#page3">Go to Page 3</a>
    </div>

    <div class="page" id="page2" style="display: none;">
        <h2>Page 2</h2>
        <p>This is the second page. You can swipe right to go back to Page 1, or click the links below.</p>
        <a href="#page1">Back to Page 1</a>
        <a href="#page3">Go to Page 3</a>
    </div>

    <div class="page" id="page3" style="display: none;">
        <h2>Page 3</h2>
        <p>This is the third page. Try swiping right to go back through your navigation history!</p>
        <a href="#page1">Back to Page 1</a>
        <a href="#page2">Go to Page 2</a>
    </div>

    <script>
        // Simple page navigation for demo
        function showPage(pageId) {
            document.querySelectorAll('.page').forEach(page => {
                page.style.display = 'none';
            });
            const targetPage = document.getElementById(pageId) || document.querySelector('.page');
            targetPage.style.display = 'block';
        }

        // Handle navigation
        window.addEventListener('hashchange', () => {
            const hash = window.location.hash.substring(1);
            showPage(hash);
        });

        // Handle initial page load
        if (window.location.hash) {
            showPage(window.location.hash.substring(1));
        }

        // Add some navigation history for gestures to work with
        document.querySelectorAll('a[href^="#"]').forEach(link => {
            link.addEventListener('click', (e) => {
                const href = link.getAttribute('href');
                history.pushState(null, null, href);
                showPage(href.substring(1));
                e.preventDefault();
            });
        });
    </script>
</body>
</html>
		`,
	})

	app.OnReady(func(ctx context.Context) {
		window.Show()
		fmt.Println("Navigation gestures demo is ready!")
		fmt.Println("On macOS, try two-finger horizontal swipe gestures to navigate back/forward")
	})

	err := app.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}