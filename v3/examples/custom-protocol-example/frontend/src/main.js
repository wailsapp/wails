import './style.css'
import { Events } from '@wailsio/runtime'

document.addEventListener('DOMContentLoaded', () => {
  const appDiv = document.getElementById('app');
  if (appDiv) {
    appDiv.innerHTML = `
      <div class="container">
        <h1>Custom Protocol / Deep Link Test</h1>
        <p>
            This page demonstrates handling custom URL schemes (deep links).
        </p>
        <p>
            <span class="label">Example Link:</span>
            Try opening this URL (e.g., by pasting it into your browser's address bar or using <code>open your-app-scheme://...</code> in terminal):
            <br>
            <a href="wailsexample://test/path?value=123&message=hello" id="example-url">wailsexample://test/path?value=123&message=hello</a>
        </p>

        <div class="url-display">
            <span class="label">Received URL:</span>
            <div id="received-url"><em>Waiting for application to be opened via a custom URL...</em></div>
        </div>
      </div>
    `;
  } else {
    console.error('Element with ID "app" not found.');
  }
});

// Listen for the event from Go
Events.On('frontend:ShowURL', (e) => {
    console.log('frontend:ShowURL event received, data:', e);
    displayUrl(e.data); 
}); 

// Make displayUrl available globally just in case, though direct call from event is better
window.displayUrl = function(url) {
    const urlElement = document.getElementById('received-url');
    if (urlElement) {
        urlElement.textContent = url || "No URL received or an error occurred.";
    } else {
        console.error("Element with ID 'received-url' not found in displayUrl.");
    }
}
