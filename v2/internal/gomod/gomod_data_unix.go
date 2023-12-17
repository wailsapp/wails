//go:build darwin || linux

package gomod

const basic string = `module changeme

go 1.17

require github.com/wailsapp/wails/v2 v2.0.0-beta.7

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => /home/lea/wails/v2
`

const basicUpdated string = `module changeme

go 1.17

require github.com/wailsapp/wails/v2 v2.0.0-beta.20

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => /home/lea/wails/v2
`

const multilineRequire = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => /home/lea/wails/v2
`

const multilineReplace = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)

replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => /home/lea/wails/v2
`

const multilineReplaceNoVersion = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)

replace github.com/wailsapp/wails/v2 => /home/lea/wails/v2
`

const multilineReplaceNoVersionBlock = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)

replace (
	github.com/wailsapp/wails/v2 => /home/lea/wails/v2
)
`

const multilineReplaceBlock = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7
)

replace (
	github.com/wailsapp/wails/v2 v2.0.0-beta.7 => /home/lea/wails/v2
)
`

const multilineRequireUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

//replace github.com/wailsapp/wails/v2 v2.0.0-beta.7 => /home/lea/wails/v2
`

const multilineReplaceUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

replace github.com/wailsapp/wails/v2 v2.0.0-beta.20 => /home/lea/wails/v2
`

const multilineReplaceNoVersionUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

replace github.com/wailsapp/wails/v2 => /home/lea/wails/v2
`

const multilineReplaceNoVersionBlockUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

replace (
	github.com/wailsapp/wails/v2 => /home/lea/wails/v2
)
`

const multilineReplaceBlockUpdated = `module changeme

go 1.17

require (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20
)

replace (
	github.com/wailsapp/wails/v2 v2.0.0-beta.20 => /home/lea/wails/v2
)
`
