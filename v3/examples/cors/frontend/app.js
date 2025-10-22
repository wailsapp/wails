// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', async () => {
    // Display current origin
    document.getElementById('currentOrigin').textContent = window.location.origin;

    // Get UI elements
    const statusIndicator = document.getElementById('statusIndicator');
    const statusText = document.getElementById('statusText');
    const results = document.getElementById('results');
    const nameInput = document.getElementById('nameInput');
    const greetBtn = document.getElementById('greetBtn');
    const timeBtn = document.getElementById('timeBtn');
    const corsBtn = document.getElementById('corsBtn');

    // Check if Wails runtime is available
    async function checkWailsRuntime() {
        try {
            // Check if wails object exists
            if (typeof wails === 'undefined') {
                throw new Error('Wails runtime not found');
            }

            // Try to make a test call
            const response = await wails.Call.ByName("main.GreetService.TestCORS");

            statusIndicator.className = 'status-indicator success';
            statusText.textContent = '✅ Wails runtime connected successfully!';

            console.log('Wails runtime test response:', response);
            return true;
        } catch (error) {
            statusIndicator.className = 'status-indicator error';
            statusText.textContent = '❌ Wails runtime not available: ' + error.message;
            console.error('Wails runtime check failed:', error);

            // Disable buttons if runtime not available
            greetBtn.disabled = true;
            timeBtn.disabled = true;
            corsBtn.disabled = true;

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

    // Test CORS button handler
    corsBtn.addEventListener('click', async () => {
        showLoading('Testing CORS configuration');

        try {
            const result = await wails.Call.ByName("main.GreetService.TestCORS");
            displayResult('CORS Test Response', result);

            // Also test runtime info
            console.log('CORS Headers Check:', {
                origin: window.location.origin,
                protocol: window.location.protocol,
                host: window.location.host,
                pathname: window.location.pathname
            });
        } catch (error) {
            displayResult('CORS Test Error', error.message, true);
            console.error('CORS test error:', error);
        }
    });

    // Log detailed runtime information
    console.log('=== Wails CORS Example ===');
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