# {{.Name}} Service

This service provides a simple URL shortener functionality within your Wails application.

## Installation

To use this service in your Wails v3 application, add it to the `Services` slice in your application options:

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "path/to/{{.Name}}"
)

app := application.New(application.Options{
    // ...
    Services: []application.Service{
        application.NewService({{.Name}}.New(), application.ServiceOptions{
            Route: "/s",
        }),
    },
    // ...
})
```

## Usage

Once the service is registered, you can use it to shorten URLs and redirect to original URLs:

```javascript
// Shorten a URL using the bound method
const shortURL = await wails.Services.Service.ShortenURL("https://wails.io");
console.log(shortURL); // Outputs: "/s/Ab3x5Y"

// Alternatively, you can use fetch to create a short URL
fetch('/s', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url: 'https://wails.io' })
})
.then(response => response.json())
.then(data => console.log(data.shortURL));

// To use a short URL, simply navigate to it or use it in an anchor tag
document.getElementById('shortLink').href = shortURL;
```

When a user visits the short URL (e.g., `/s/Ab3x5Y`), they will be redirected to the original URL.

## API Reference

- `ShortenURL(url: string): Promise<string>`
  Returns a shortened URL for the given original URL.

## HTTP Handling

This service implements the `ServeHTTP(w http.ResponseWriter, r *http.Request)` method to handle two types of HTTP requests:

1. POST requests to create short URLs:
    - Endpoint: `/s`
    - Body: JSON object with a `url` field
    - Response: JSON object with a `shortURL` field

2. GET requests to redirect to original URLs:
    - Endpoint: `/s/{shortCode}`
    - Action: Redirects to the original URL associated with the short code

## Considerations

- This is a simple in-memory URL shortener. In a production environment, you'd want to use a persistent storage solution.
- There's no mechanism to prevent duplicate short codes. In a real-world scenario, you'd want to ensure uniqueness.
- This service doesn't include features like custom short codes or expiration dates, which could be added for a more full-featured URL shortener.

## Support

If you encounter any issues or have questions about this service, please raise a ticket in our [GitHub repository](https://github.com/path/to/repo/issues).

Please note that this is a community-contributed service. While we appreciate the Wails team's efforts, direct your support requests to this service's maintainers rather than the core Wails team.

## License

[Include license information here]

## Contributing

We welcome contributions to improve this service! Please see our [Contributing Guidelines](CONTRIBUTING.md) for more information on how to get started.