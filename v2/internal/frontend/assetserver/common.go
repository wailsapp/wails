package assetserver

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

type optionType string

const (
	noAutoInject        optionType = "noautoinject"
	noAutoInjectRuntime optionType = "noautoinjectruntime"
	noAutoInjectIPC     optionType = "noautoinjectipc"
)

type Options struct {
	disableRuntimeInjection bool
	disableIPCInjection     bool
}

func newOptions(optionString string) *Options {
	var result = &Options{}
	optionString = strings.ToLower(optionString)
	options := strings.Split(optionString, ",")
	for _, option := range options {
		switch optionType(strings.TrimSpace(option)) {
		case noAutoInject:
			result.disableRuntimeInjection = true
			result.disableIPCInjection = true
		case noAutoInjectIPC:
			result.disableIPCInjection = true
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

func createScriptNode(scriptName string) *html.Node {
	return &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{
			{
				Key: "src",
				Val: scriptName,
			},
		},
	}
}

func createDivNode(id string) *html.Node {
	return &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{
				Namespace: "",
				Key:       "id",
				Val:       id,
			},
		},
	}
}

func insertScriptInHead(htmlNode *html.Node, scriptName string) error {
	headNode := findFirstTag(htmlNode, "head")
	if headNode == nil {
		return errors.New("cannot find head in HTML")
	}
	scriptNode := createScriptNode(scriptName)
	if headNode.FirstChild != nil {
		headNode.InsertBefore(scriptNode, headNode.FirstChild)
	} else {
		headNode.AppendChild(scriptNode)
	}
	return nil
}

func appendSpinnerToBody(htmlNode *html.Node) error {
	bodyNode := findFirstTag(htmlNode, "body")
	if bodyNode == nil {
		return errors.New("cannot find body in HTML")
	}
	scriptNode := createDivNode("wails-spinner")
	bodyNode.AppendChild(scriptNode)
	return nil
}

func getHTMLNode(htmldata []byte) (*html.Node, error) {
	return html.Parse(bytes.NewReader(htmldata))
}

func findFirstTag(htmlnode *html.Node, tagName string) *html.Node {
	var extractor func(*html.Node) *html.Node
	var result *html.Node
	extractor = func(node *html.Node) *html.Node {
		if node.Type == html.ElementNode && node.Data == tagName {
			return node
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			result := extractor(child)
			if result != nil {
				return result
			}
		}
		return nil
	}
	result = extractor(htmlnode)
	return result
}
