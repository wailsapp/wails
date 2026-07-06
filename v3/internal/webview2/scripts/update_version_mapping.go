package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
	"updater/generator"
)

const URL = "https://raw.githubusercontent.com/MicrosoftDocs/edge-developer/master/microsoft-edge/webview2/release-notes/index.md"

//go:embed latest_version.txt
var latestVersionProcessed string

type Version struct {
	Number         string
	ReleaseNotes   string
	RuntimeVersion string
	Notes          []string
}

const debug = false

func getDoc() []byte {
	if debug {
		data, err := os.ReadFile("test.md")
		if err != nil {
			log.Fatal(err)
		}
		return data
	}
	// GET the URL
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body()
}

func extractVersion(in string) string {
	// Match the version numbers by format: 1.0.774.44, 0.9.515-prerelease,
	regex := regexp.MustCompile(`\d+\.\d+\.\d+(\.\d+|-prerelease)`)
	version := regex.Find([]byte(in))
	return string(version)
}

var latestVersion string

func main() {

	var forced bool
	if len(os.Args) > 1 {
		forced = os.Args[1] == "-forced"
	}

	var buf bytes.Buffer
	data := getDoc()
	buf.Write(data)
	var s scanner.Scanner
	s.Init(&buf)

	r := bufio.NewReader(&buf)

	var err error
	var line []byte
	var versions []*Version
	var currentVersion *Version
	var nomnom bool
	for err == nil {
		line, _, err = r.ReadLine()

		// Check if line starts with `[NuGet package for WebView2 `
		l := string(line)

		if strings.HasPrefix(l, `## `) {
			nomnom = false
			continue
		}

		if currentVersion != nil && nomnom {
			currentVersion.Notes = append(currentVersion.Notes, l)
		}

		if strings.HasPrefix(l, `[NuGet package for WebView2 `) {
			version := extractVersion(l)
			if version == "" {
				continue
			}
			if currentVersion != nil {
				versions = append(versions, currentVersion)
			} else {
				latestVersion = version
			}
			currentVersion = &Version{
				Number: version,
			}
			continue
		}
		if strings.HasSuffix(strings.TrimSpace(l), "or higher.") {
			if currentVersion != nil {
				currentVersion.RuntimeVersion = extractVersion(l)
				currentVersion.ReleaseNotes = `https://learn.microsoft.com/en-us/microsoft-edge/webview2/release-notes?tabs=win32cpp#` + strings.Replace(currentVersion.Number, ".", "", -1)
				nomnom = true
			}
			continue
		}

		if strings.HasPrefix(l, `Release Date:`) {
			if currentVersion != nil {
				currentVersion.Notes = append(currentVersion.Notes, strings.Trim(l, " "))
			}
		}
	}

	var buffer strings.Builder
	buffer.WriteString("//go:build windows\n\n")
	buffer.WriteString("package edge\n\n")
	buffer.WriteString("type Version struct {\n")
	buffer.WriteString("	SDKVersion         string\n")
	buffer.WriteString("	ReleaseNotes           string\n")
	buffer.WriteString("	RuntimeVersion string\n")
	buffer.WriteString("	Notes          string\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("var versionMapping = map[string]Version{\n")
	for _, version := range versions {
		buffer.WriteString(fmt.Sprintf("	\"%s\": {\n", version.Number))
		buffer.WriteString(fmt.Sprintf("		SDKVersion:     \"%s\",\n", version.Number))
		buffer.WriteString(fmt.Sprintf("		ReleaseNotes:   \"%s\",\n", version.ReleaseNotes))
		buffer.WriteString(fmt.Sprintf("		RuntimeVersion: \"%s\",\n", version.RuntimeVersion))
		buffer.WriteString("		Notes: ")
		buffer.WriteString(fmt.Sprintf("			`%s`,\n", strings.Replace(strings.Join(version.Notes, "\n"), "`", "'", -1)))
		buffer.WriteString("	},\n")
	}
	buffer.WriteString("}\n")

	// Write the buffer to ../pkg/edge/version_map.go
	err = os.WriteFile("../pkg/edge/version_map.go", []byte(buffer.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Save the latest release notes to a file
	if len(versions) > 0 {
		latestReleaseNotes := fmt.Sprintf("Version: %s\nRuntime Version: %s\nRelease Notes URL: %s\n\nNotes:\n%s",
			versions[0].Number,
			versions[0].RuntimeVersion,
			versions[0].ReleaseNotes,
			strings.Join(versions[0].Notes, "\n"))
		err = os.WriteFile("latest_release_notes.txt", []byte(latestReleaseNotes), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !forced {
		// Check if the latest version is different from the last time we ran this script
		latest, err := CompareBrowserVersions(latestVersion, latestVersionProcessed)
		if err != nil {
			log.Fatal(err)
		}
		if latest != 1 {
			println("No new version found")
			os.Exit(0)
		}
	}

	println("Processing version: ", latestVersion)
	// Download Webview2 IDL for this version
	idlData, err := DownloadIDL(latestVersion)
	if err != nil {
		log.Fatal(err)
	}

	files, err := generator.ParseIDL(idlData)
	if err != nil {
		log.Fatal(err)
	}

	// Write the files to the ../pkg/webview2 directory
	_ = os.Mkdir("../pkg/webview2", 0755)
	for _, file := range files {
		fileName := "../pkg/webview2/" + file.FileName
		println("Writing: ", fileName)
		err = os.WriteFile(fileName, file.Content.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Save the version to latest_version.txt
	err = os.WriteFile("latest_version.txt", []byte(latestVersion), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func DownloadIDL(version string) ([]byte, error) {

	// Look for the file locally: WebView2.version.idl
	data, err := os.ReadFile("WebView2." + version + ".idl")
	if err == nil {
		return data, nil
	}

	// URL for the nuget package: https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/<version>
	// Download the package to the current directory
	client := resty.New()
	println("Downloading: ", fmt.Sprintf("https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/%s", version))
	resp, err := client.R().
		EnableTrace().
		Get(fmt.Sprintf("https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/%s", version))
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(resp.Body())
	zr, err := zip.NewReader(reader, int64(reader.Len()))
	if err != nil {
		return nil, err
	}

	var idlData []byte
	for _, file := range zr.File {
		if file.Name == "WebView2.idl" {
			r, err := file.Open()
			if err != nil {
				return nil, err
			}
			idlData, err = io.ReadAll(r)
			if err != nil {
				return nil, err
			}
		}
	}
	// Write IDL to disk
	err = os.WriteFile("WebView2."+version+".idl", idlData, 0755)
	if err != nil {
		return nil, err
	}
	return idlData, nil
}

// CompareBrowserVersions will compare the 2 given versions and return:
//
//	-1 = v1 < v2
//	 0 = v1 == v2
//	 1 = v1 > v2
func CompareBrowserVersions(v1 string, v2 string) (int, error) {
	v, err := parseVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("v1 invalid: %w", err)
	}

	w, err := parseVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("v2 invalid: %w", err)
	}

	return v.compare(w), nil
}

type version struct {
	major int
	minor int
	patch int
	build int

	channel string
}

func (v version) String() string {
	vv := fmt.Sprintf("%d.%d.%d.%d", v.major, v.minor, v.patch, v.build)
	if v.channel != "" {
		vv += " " + v.channel
	}

	return vv
}

func (v version) compare(o version) int {
	if c := compareInt(v.major, o.major); c != 0 {
		return c
	}
	if c := compareInt(v.minor, o.minor); c != 0 {
		return c
	}
	if c := compareInt(v.patch, o.patch); c != 0 {
		return c
	}
	return compareInt(v.build, o.build)
}

func parseVersion(v string) (version, error) {
	var p version

	// Split away channel information...
	if i := strings.Index(v, " "); i > 0 {
		p.channel = v[i+1:]
		v = v[:i]
	}

	vv := strings.Split(v, ".")
	if len(vv) > 4 {
		return p, fmt.Errorf("too many version parts")
	}

	var err error
	vv, p.major, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad major version: %w", err)
	}

	vv, p.minor, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad minor version: %w", err)
	}

	vv, p.patch, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad patch version: %w", err)
	}

	_, p.build, err = parseInt(vv)
	if err != nil {
		return p, fmt.Errorf("bad build version: %w", err)
	}

	return p, nil
}

func parseInt(v []string) ([]string, int, error) {
	if len(v) == 0 {
		return nil, 0, nil
	}

	p, err := strconv.ParseInt(v[0], 10, 32)
	if err != nil {
		return nil, 0, err
	}
	return v[1:], int(p), nil
}

func compareInt(v1, v2 int) int {
	if v1 == v2 {
		return 0
	}
	if v1 < v2 {
		return -1
	} else {
		return +1
	}
}
