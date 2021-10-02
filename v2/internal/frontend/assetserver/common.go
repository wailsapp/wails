package assetserver

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

type optionType string

const (
	noAutoInject         optionType = "noautoinject"
	noAutoInjectRuntime  optionType = "noautoinjectruntime"
	noAutoInjectBindings optionType = "noautoinjectbindings"
)

type Options struct {
	disableRuntimeInjection  bool
	disableBindingsInjection bool
}

func newOptions(optionString string) *Options {
	var result = &Options{}
	optionString = strings.ToLower(optionString)
	options := strings.Split(optionString, ",")
	for _, option := range options {
		switch optionType(strings.TrimSpace(option)) {
		case noAutoInject:
			result.disableRuntimeInjection = true
			result.disableBindingsInjection = true
		case noAutoInjectBindings:
			result.disableBindingsInjection = true
		case noAutoInjectRuntime:
			result.disableRuntimeInjection = true
		}
	}
	return result
}

func injectHTML(input string, html string) ([]byte, error) {
	splits := strings.Split(input, "</head>")
	if len(splits) != 2 {
		return nil, fmt.Errorf("unable to locate a </head> tag in your html")
	}

	var result bytes.Buffer
	result.WriteString(splits[0])
	result.WriteString(html)
	result.WriteString("</head>")
	result.WriteString(splits[1])
	return result.Bytes(), nil
}

func extractOptions(htmldata []byte) (*Options, error) {
	doc, err := html.Parse(bytes.NewReader(htmldata))
	if err != nil {
		return nil, err
	}
	var extractor func(*html.Node) *Options
	extractor = func(node *html.Node) *Options {
		if node.Type == html.ElementNode && node.Data == "meta" {
			isWailsOptionsTag := false
			wailsOptions := ""
			for _, attr := range node.Attr {
				if isWailsOptionsTag && attr.Key == "content" {
					wailsOptions = attr.Val
				}
				if attr.Val == "wails-options" {
					isWailsOptionsTag = true
				}
			}
			return newOptions(wailsOptions)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			result := extractor(child)
			if result != nil {
				return result
			}
		}
		return nil
	}
	result := extractor(doc)
	if result == nil {
		result = &Options{}
	}
	return result, nil
}
