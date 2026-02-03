package main

import (
	"log"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "File Input Test",
		Description: "Test for HTML file input (#4862)",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			Handler: http.HandlerFunc(handler),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "File Input Test",
		Width:  700,
		Height: 500,
		URL:    "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

const html = `<!DOCTYPE html>
<html>
<head>
    <title>File Input Test</title>
    <script src="/wails/runtime.js"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, sans-serif;
            padding: 20px;
            background: #f5f5f5;
        }
        h1 { color: #333; font-size: 18px; margin-bottom: 20px; }
        .grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
        }
        .card {
            background: white;
            padding: 15px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .card h2 { font-size: 14px; margin: 0 0 10px 0; }
        input[type="file"] { font-size: 12px; width: 100%; }
        button { font-size: 12px; padding: 6px 12px; cursor: pointer; }
        #result {
            margin-top: 15px;
            padding: 10px;
            background: #e8f5e9;
            border-radius: 4px;
            white-space: pre-wrap;
            font-size: 12px;
            max-height: 120px;
            overflow-y: auto;
        }
    </style>
</head>
<body>
    <h1>File Input Test (#4862)</h1>
    <div class="grid">
        <div class="card">
            <h2>1. Single File</h2>
            <input type="file" onchange="show(this)">
        </div>
        <div class="card">
            <h2>2. Multiple Files</h2>
            <input type="file" multiple onchange="show(this)">
        </div>
        <div class="card">
            <h2>3. Files or Directories</h2>
            <input type="file" webkitdirectory onchange="show(this)">
        </div>
        <div class="card">
            <h2>4. Directory Only (Wails API)</h2>
            <button onclick="openDir()">Choose Directory</button>
        </div>
        <div class="card">
            <h2>5. Filtered .txt (Wails API)</h2>
            <button onclick="openFiltered()">Choose .txt File</button>
        </div>
    </div>
    <div id="result">Click a file input or button to test...</div>
    <script>
        function show(input) {
            const r = document.getElementById('result');
            if (!input.files.length) {
                r.textContent = 'Cancelled';
                return;
            }
            let t = 'Selected ' + input.files.length + ' file(s):\n';
            for (const f of input.files) {
                t += '• ' + f.name + ' (' + f.size + ' bytes)\n';
            }
            r.textContent = t;
        }
        async function openDir() {
            const r = document.getElementById('result');
            try {
                const path = await wails.Dialogs.OpenFile({
                    Title: 'Select Directory',
                    CanChooseDirectories: true,
                    CanChooseFiles: false
                });
                r.textContent = path || 'Cancelled';
            } catch (e) {
                r.textContent = 'Error: ' + e;
            }
        }
        async function openFiltered() {
            const r = document.getElementById('result');
            try {
                const path = await wails.Dialogs.OpenFile({
                    Title: 'Select Text File',
                    Filters: [{ DisplayName: 'Text Files', Pattern: '*.txt' }]
                });
                r.textContent = path || 'Cancelled';
            } catch (e) {
                r.textContent = 'Error: ' + e;
            }
        }
    </script>
</body>
</html>`
