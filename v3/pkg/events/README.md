# Events

This package is used to generate the event management code and to allow quick addition of events.

## Usage

1. Add events to `events.txt`
2. Run `task generate:events`

## Notes

For events that you want to handle manually, add a `!` to the end of the event name and 
add custom code into the appropriate platform. See [this PR](https://github.com/wailsapp/wails/pull/2991)
for an example of how to do this.