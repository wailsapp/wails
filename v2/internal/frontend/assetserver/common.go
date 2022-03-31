package assetserver

import (
	"bytes"
	"errors"

	"golang.org/x/net/html"
)

const (
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"
	HeaderUserAgent     = "User-Agent"
	HeaderCacheControl  = "Cache-Control"

	WailsUserAgentValue = "wails.io"
)

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
