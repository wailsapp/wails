# Your First Application

Creating your first application with Wails v3 Alpha is an exciting journey into the world of modern desktop app development. This guide will walk you through the process of creating a basic application, showcasing the power and simplicity of Wails.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go (version 1.21 or later)
- Node.js (LTS version)
- Wails v3 Alpha (see the [installation guide](installation.md) for instructions)

## Step 1: Creating a New Project

Open your terminal and run the following command to create a new Wails project:

`wails3 init -n myfirstapp`

This command creates a new directory called `myfirstapp` with all the necessary files.

## Step 2: Exploring the Project Structure

Navigate to the `myfirstapp` directory. You'll find several files and folders:

- `build`: Contains files used by the build process.
- `frontend`: Contains your web frontend code.
- `go.mod` & `go.sum`: Go module files.
- `main.go`: The entry point for your Wails application.
- `Taskfile.yml`: Defines all the tasks used by the build system. Learn more at the [Task](https://taskfile.dev/) website.

Take a moment to explore these files and familiarize yourself with the structure.

!!! note
    Although Wails v3 uses [Task](https://taskfile.dev/) as its default build system, there is nothing stopping you from using `make` or any other alternative build system.  

## Step 3: Building Your Application

To build your application, execute:

`wails3 build`

This command compiles a debug version of your application and saves it in a new `bin` directory. 
You can run this like you would any normal application:


=== "Mac"

    `./bin/myfirstapp`

=== "Windows"

    `bin\myfirstapp.exe`

=== "Linux"

    `./bin/myfirstapp`

You'll see a simple UI, the starting point for your application. As it is the debug version, you'll also see logs in the console window. This is useful for debugging purposes.

## Step 4: Dev Mode

We can also run the application in development mode. This mode allows you to make changes to your frontend code and see the changes reflected in the running application without having to rebuild the entire application.

1. Open a new terminal window.
2. Run `wails3 dev`.
3. Open `frontend/main.js`.
4. Change the line that has `<h1>Hello Wails!</h1>` to `<h1>Hello World!</h1>`.
5. Save the file.

The application will update automatically, and you'll see the changes reflected in the running application. 

## Step 5: Building the Application Again

Once you're happy with your changes, build your application again:

`wails3 build`

You'll notice that the build time was faster this time. That's because the new build system only builds the parts of your application that have changed.

You should see a new executable in the `build` directory.

## Conclusion

Congratulations! You've just created and built your first Wails application. This is just the beginning of what you can achieve with Wails v3 Alpha. Explore the documentation, experiment with different features, and start building amazing applications!
