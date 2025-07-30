# HTTP CORS Workaround Example

This example demonstrates how to use the Wails v3 HTTP runtime API to bypass CORS issues when making network requests from the frontend.

## The Problem

When using Wails with the `wails://wails` protocol on Linux/macOS, CORS restrictions prevent direct HTTP requests from the frontend to external APIs. This is because:

1. Wails uses a custom protocol (`wails://wails`) which is not recognized by external servers
2. The browser sends an empty Origin header for custom protocols
3. External servers reject requests without a valid Origin header

## The Solution

The Wails HTTP runtime API allows you to make HTTP requests through the Go backend, completely bypassing CORS restrictions.

## Usage

Instead of using `fetch` or `axios` directly:

```javascript
// This will fail with CORS error
const response = await fetch('https://api.example.com/data');
```

Use the Wails HTTP API:

```javascript
// This works - request is made from Go backend
const response = await wails.HTTP.Get('https://api.example.com/data');
```

## API Methods

- `wails.HTTP.Get(url, options)` - GET request
- `wails.HTTP.Post(url, body, options)` - POST request
- `wails.HTTP.Put(url, body, options)` - PUT request
- `wails.HTTP.Delete(url, options)` - DELETE request
- `wails.HTTP.Patch(url, body, options)` - PATCH request
- `wails.HTTP.Head(url, options)` - HEAD request
- `wails.HTTP.Fetch(options)` - Generic request with full options

## Example with Headers and Timeout

```javascript
const response = await wails.HTTP.Post('https://api.example.com/users', {
    name: 'John Doe',
    email: 'john@example.com'
}, {
    headers: {
        'Authorization': 'Bearer token123',
        'X-Custom-Header': 'value'
    },
    timeout: 10 // 10 seconds timeout
});

if (response.error) {
    console.error('Request failed:', response.error);
} else {
    console.log('Status:', response.statusCode);
    console.log('Data:', JSON.parse(response.body));
}
```

## Running the Example

1. Navigate to this directory
2. Run `wails3 dev`
3. Click the buttons to test various HTTP methods
4. Check the console for responses

## Notes

- All requests are made from the Go backend, so there are no CORS restrictions
- The response includes `statusCode`, `headers`, `body`, and `error` (if any)
- Request bodies are automatically JSON stringified for objects
- Default timeout is 30 seconds, but can be customized