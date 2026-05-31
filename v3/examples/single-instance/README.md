# Single Instance Example

This example demonstrates the single instance functionality in Wails v3. It shows how to:

1. Ensure only one instance of your application can run at a time
2. Notify the first instance when a second instance is launched
3. Pass data between instances
4. Handle command line arguments and working directory information from second instances

## Running the Example

1. Build and run the application:
   ```bash
   go build
   ./single-instance
   ```

2. Try launching a second instance of the application. You'll notice:
   - The second instance will exit immediately
   - The first instance will receive and display:
     - Command line arguments from the second instance
     - Working directory of the second instance
     - Additional data passed from the second instance

3. Check the application logs to see the information received from second instances.

## Features Demonstrated

- Setting up single instance lock with a unique identifier
- Handling second instance launches through callbacks
- Passing custom data between instances
- Displaying instance information in a web UI
- Cross-platform support (Windows, macOS, Linux)

## Code Overview

The example consists of:

- `main.go`: The main application code demonstrating single instance setup
- A simple web UI showing current instance information
- Callback handling for second instance launches

## Implementation Details

The application uses the Wails v3 single instance feature:

```go
app := application.New(&application.Options{
    SingleInstance: &application.SingleInstanceOptions{
        UniqueID: "com.wails.example.single-instance",
        OnSecondInstance: func(data application.SecondInstanceData) {
            // Handle second instance launch
        },
		AdditionalData: map[string]string{
		},
    },
})
```

The implementation uses platform-specific mechanisms:
- Windows: Named mutex and window messages
- Unix (Linux/macOS): File locking with flock and signals
