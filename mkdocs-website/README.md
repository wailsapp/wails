# v3 Docs

This is the documentation for Wails v3. It is currently a work in progress.

If you do not wish to build it locally, it is available online at
[https://wailsapp.github.io/wails/](https://wailsapp.github.io/wails/).

## Setup Steps

1. Install the wails3 CLI if you haven't already:

    ```shell
    git clone https://github.com/wailsapp/wails.git
    cd wails
    git checkout v3-alpha
    cd v3/cmd/wails3
    go install
    ```
2. Install [docker](https://www.docker.com)
3. Run the following command to build the docker container:

    ```shell
    wails3 task docs:setup
    ```
4. Serve the documentation locally:

    ```shell
    wails3 task docs:serve
    ```

5. Open your browser to [http://127.0.0.1:8000](http://127.0.0.1:8000)

6. For a complete build, run:

    ```shell
    wails3 task docs:build
    ```

## Contributing

If you would like to contribute to the documentation, please feel free to open a
PR!
