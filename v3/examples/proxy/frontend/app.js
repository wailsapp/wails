// Function to wait for Wails runtime to be available
function waitForWailsRuntime(timeout = 5000) {
    return new Promise((resolve, reject) => {
        const startTime = Date.now();

        const checkRuntime = () => {
            if (typeof wails !== 'undefined' && wails.Call && wails.Call.ByName) {
                console.log('✅ Wails runtime detected!');
                resolve();
            } else if (Date.now() - startTime > timeout) {
                reject(new Error('Timeout waiting for Wails runtime'));
            } else {
                setTimeout(checkRuntime, 100);
            }
        };

        checkRuntime();
    });
}

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', async () => {
    console.log('DOM loaded, waiting for Wails runtime to be injected...');

    // Wait for Wails runtime to be injected by the WebView
    try {
        await waitForWailsRuntime(10000); // Wait up to 10 seconds
        console.log('Wails runtime is ready!');
    } catch (error) {
        console.error('Failed to load Wails runtime:', error);
    }
    // Display current origin
    document.getElementById('currentOrigin').textContent = window.location.origin;

    // Get UI elements
    const statusIndicator = document.getElementById('statusIndicator');
    const statusText = document.getElementById('statusText');
    const results = document.getElementById('results');
    const nameInput = document.getElementById('nameInput');
    const greetBtn = document.getElementById('greetBtn');
    const timeBtn = document.getElementById('timeBtn');
    const proxyBtn = document.getElementById('proxyBtn');

    // Check if Wails runtime is available
    async function checkWailsRuntime() {
        try {
            // Check if wails object exists
            if (typeof wails === 'undefined') {
                throw new Error('Wails runtime not found');
            }

            // Try to make a test call
            const response = await wails.Call.ByName("main.GreetService.GetTime");

            statusIndicator.className = 'status-indicator success';
            statusText.textContent = '✅ Wails runtime connected successfully!';

            console.log('Wails runtime test - Server time:', response);
            return true;
        } catch (error) {
            statusIndicator.className = 'status-indicator error';
            statusText.textContent = '❌ Wails runtime not available: ' + error.message;
            console.error('Wails runtime check failed:', error);

            // Disable buttons if runtime not available
            greetBtn.disabled = true;
            timeBtn.disabled = true;
            proxyBtn.disabled = true;

            results.innerHTML = `
                <div class="error">
                    <strong>Error:</strong> Unable to connect to Wails runtime.<br>
                    Make sure this page is loaded inside a Wails WebView window.
                </div>
            `;
            return false;
        }
    }

    // Format result for display
    function displayResult(title, data, isError = false) {
        const className = isError ? 'error' : 'success';
        const icon = isError ? '❌' : '✅';

        if (typeof data === 'object') {
            results.innerHTML = `
                <div class="${className}">
                    <strong>${icon} ${title}</strong>
                    <pre>${JSON.stringify(data, null, 2)}</pre>
                </div>
            `;
        } else {
            results.innerHTML = `
                <div class="${className}">
                    <strong>${icon} ${title}</strong>
                    <p>${data}</p>
                </div>
            `;
        }
    }

    // Display loading state
    function showLoading(message) {
        results.innerHTML = `<p class="loading">⏳ ${message}...</p>`;
    }

    // Greet button handler
    greetBtn.addEventListener('click', async () => {
        const name = nameInput.value.trim() || 'World';
        showLoading('Calling Greet method');

        try {
            const result = await wails.Call.ByName("main.GreetService.Greet", name);
            displayResult('Greet Response', result);
        } catch (error) {
            displayResult('Greet Error', error.message, true);
            console.error('Greet error:', error);
        }
    });

    // Get Time button handler
    timeBtn.addEventListener('click', async () => {
        showLoading('Getting server time');

        try {
            const result = await wails.Call.ByName("main.GreetService.GetTime");
            displayResult('Server Time', result);
        } catch (error) {
            displayResult('Time Error', error.message, true);
            console.error('GetTime error:', error);
        }
    });

    // Test Proxy button handler
    proxyBtn.addEventListener('click', async () => {
        showLoading('Testing proxy configuration');

        try {
            // The fact this page loaded proves the proxy works for HTML/JS
            // Now test if we can fetch additional resources through the proxy
            console.log('Testing proxy by fetching /health endpoint from external server...');

            let proxyWorking = false;
            let healthData = null;
            let proxyError = null;

            try {
                // This fetches from the external server through the proxy
                const response = await fetch('/health');
                if (response.ok) {
                    healthData = await response.json();
                    proxyWorking = true;
                } else {
                    proxyError = `HTTP ${response.status}: ${response.statusText}`;
                }
            } catch (fetchError) {
                proxyError = fetchError.message;
            }

            // Build comprehensive result
            const result = {
                proxy_status: {
                    html_js_loaded: "✅ Success (this page loaded from external server)",
                    runtime_injected: "✅ Success (Wails runtime is working)",
                    health_endpoint: proxyWorking ? "✅ Success" : "❌ Failed",
                    summary: proxyWorking
                        ? "✅ Proxy fully operational! All content served through local proxy."
                        : "⚠️  Proxy partially working. Main assets loaded but /health failed."
                },
                health_data: healthData,
                error: proxyError,
                location_info: {
                    origin: window.location.origin,
                    protocol: window.location.protocol,
                    host: window.location.host,
                    note: "Running from local asset server, proxying to external"
                }
            };

            displayResult('Proxy Test Results', result, !!proxyError);

            // Log detailed info
            console.log('Proxy Test Results:', result);
        } catch (error) {
            displayResult('Proxy Test Error', error.message, true);
            console.error('Proxy test error:', error);
        }
    });

    // Log detailed runtime information
    console.log('=== Wails Proxy Example ===');
    console.log('Page Origin:', window.location.origin);
    console.log('Protocol:', window.location.protocol);
    console.log('Host:', window.location.host);
    console.log('Pathname:', window.location.pathname);

    // Check for Wails global objects
    console.log('Wails Runtime Available:', typeof wails !== 'undefined');
    if (typeof wails !== 'undefined') {
        console.log('Wails Modules:', Object.keys(wails));
    }

    // Initialize - check runtime availability
    await checkWailsRuntime();

    // Test keyboard shortcuts if available
    if (typeof wails !== 'undefined' && wails.Window) {
        // Register Ctrl+R to reload (for development)
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.key === 'r') {
                e.preventDefault();
                if (wails.Window && wails.Window.Reload) {
                    wails.Window.Reload();
                } else {
                    window.location.reload();
                }
            }
        });
    }
});