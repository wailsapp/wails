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

The project is built using a special mode that starts a webserver and starts listening to port 34115. When the frontend project is run independently, so long as the JS is wrapped with the runtime method `ready`, then the frontend will connect to the backend code via websockets. The interface should be present in your browser, and you should be able to interact with the backend as you would in a desktop app.  