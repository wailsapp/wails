# GTK3 vs GTK4 Benchmark

This benchmark suite compares the performance of Wails applications running on GTK3 vs GTK4.

## Building

Build both versions:

```bash
# Build GTK4 version (default)
go build -tags gtk4 -o benchmark-gtk4 .

# Build GTK3 version
go build -tags gtk3 -o benchmark-gtk3 .
```

## Running Benchmarks

Run each version to generate a report:

```bash
# Run GTK4 benchmark
./benchmark-gtk4

# Run GTK3 benchmark
./benchmark-gtk3
```

Each run will:
1. Display results in the console
2. Save a JSON report file (e.g., `benchmark-GTK4-WebKitGTK-6.0-20240115-143052.json`)

## Comparing Results

Use the comparison tool to analyze two reports:

```bash
go run compare.go benchmark-GTK3-*.json benchmark-GTK4-*.json
```

This will output a side-by-side comparison showing:
- Average times for each benchmark
- Percentage change between GTK3 and GTK4
- Summary of improvements and regressions

## Benchmarks Included

| Benchmark | Description |
|-----------|-------------|
| Screen Enumeration | Query all connected screens |
| Primary Screen Query | Get the primary display |
| Window Create/Destroy | Create and close windows |
| Window SetSize | Resize window operations |
| Window SetTitle | Update window title |
| Window Size Query | Get current window dimensions |
| Window Position Query | Get current window position |
| Window Center | Center window on screen |
| Window Show/Hide | Toggle window visibility |
| Menu Creation (Simple) | Create basic menus |
| Menu Creation (Complex) | Create nested menu structures |
| Menu with Accelerators | Menus with keyboard shortcuts |
| Event Emit | Dispatch custom events |
| Event Emit+Receive | Round-trip event handling |
| Dialog Setup (Info) | Create info dialog |
| Dialog Setup (Question) | Create question dialog |

## Expected Results

GTK4 improvements typically include:
- Better Wayland support
- Improved GPU rendering pipeline
- More efficient event dispatch
- Better fractional scaling support

Performance varies by operation - some may be faster in GTK4, others similar to GTK3.
