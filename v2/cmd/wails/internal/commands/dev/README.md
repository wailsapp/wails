# Dev

The dev command allows you to develop your application through a standard browser. 

## Usage

`wails dev <flags>`

### Flags

| Flag           | Details      | Default |
| :------------- | :----------- | :------ |
| -compiler path/to/compiler  | Use a different go compiler, eg go1.15beta1 | go |
| -ldflags "custom ld flags" | Use given ldflags | | 
| -e list,of,extensions | File extensions to trigger rebuilds | go |
| -w | Show warnings | false |
| -v int | Verbosity level (0 - silent, 1 - default, 2 - verbose) | 1 |
| -loglevel  | Loglevel to pass to the application - Trace, Debug, Info, Warning, Error | Debug |

## How it works

The project is build