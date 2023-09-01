# v3 Docs

This is the documentation for Wails v3. It is currently a work in progress.

## Recommended Setup Steps

Install the wails3 CLI if you haven't already:

```shell
go install github.com/wailsapp/wails/v2/cmd/wails3@latest
```

The documentation uses mkdocs, so you will need to install [Python](https://www.python.org/). Once installed, 
you can setup the documentation by running the following command:

```bash
wails3 task docs:setup
```

This will install the required dependencies for you. 

If you have installed the wails3 CLI, you can run the following command to build the documentation and serve it locally:

```bash
wails3 task docs:serve
```

### Manual Setup

To install manually, you will need to do the following:

  - Install [Python](https://www.python.org/)
  - Run `pip install -r requirements.txt` to install the required dependencies
  - Run `mkdocs serve` to serve the documentation locally
  - Run `mkdocs build` to build the documentation


## Contributing

If you would like to contribute to the documentation, please feel free to open a PR!

