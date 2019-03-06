package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/browser"

	"github.com/wailsapp/wails/cmd"
)

func init() {

	commandDescription := `Generates an issue in Github using the given title, description and system report.`

	initCommand := app.Command("issue", "Generates an issue in Github.").
		LongDescription(commandDescription)

	initCommand.Action(func() error {

		logger.PrintSmallBanner("Generate Issue")
		fmt.Println()
		message := `Thanks for taking the time to submit an issue!

To help you in this process, we will ask for some information, add Go/Wails details automatically, then prepare the issue for your editing and submission.
`

		logger.Yellow(message)

		title := cmd.Prompt("Issue Title")
		description := cmd.Prompt("Issue Description")

		var str strings.Builder

		gomodule, exists := os.LookupEnv("GO111MODULE")
		if !exists {
			gomodule = "(Not Set)"
		}

		str.WriteString("\n| Name   | Value |\n| ----- | ----- |\n")
		str.WriteString(fmt.Sprintf("| Wails Version | %s |\n", cmd.Version))
		str.WriteString(fmt.Sprintf("| Go Version    | %s |\n", runtime.Version()))
		str.WriteString(fmt.Sprintf("| Platform      | %s |\n", runtime.GOOS))
		str.WriteString(fmt.Sprintf("| Arch          | %s |\n", runtime.GOARCH))
		str.WriteString(fmt.Sprintf("| GO111MODULE   | %s |\n", gomodule))

		fmt.Println()
		fmt.Println("Processing template and preparing for upload.")

		// Grab issue template
		resp, err := http.Get("https://raw.githubusercontent.com/wailsapp/wails/master/.github/ISSUE_TEMPLATE/bug_report.md")
		if err != nil {
			logger.Red("Unable to read in issue template. Are you online?")
			os.Exit(1)
		}
		defer resp.Body.Close()
		template, _ := ioutil.ReadAll(resp.Body)
		body := string(template)
		body = "**Description**\n" + (strings.Split(body, "**Description**")[1])
		fullURL := "https://github.com/wailsapp/wails/issues/new?"
		body = strings.Replace(body, "A clear and concise description of what the bug is.", description, -1)
		body = strings.Replace(body, "Please paste the output of `wails report` here.", str.String(), -1)
		params := "title=" + title + "&body=" + body

		fmt.Println("Opening browser to file issue.")
		browser.OpenURL(fullURL + url.PathEscape(params))
		return nil
	})
}
