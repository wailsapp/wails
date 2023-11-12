package assetserver

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang.org/x/net/html"
)

func BuildAssetServerConfig(appOptions *options.App) (assetserver.Options, error) {
	var options assetserver.Options
	if opt := appOptions.AssetServer; opt != nil {
		if appOptions.Assets != nil || appOptions.AssetsHandler != nil {
			panic("It's not possible to use the deprecated Assets and AssetsHandler options and the new AssetServer option at the same time. Please migrate all your Assets options to the AssetServer option.")
		}

		options = *opt
	} else {
		options = assetserver.Options{
			Assets:  appOptions.Assets,
			Handler: appOptions.AssetsHandler,
		}
	}

	return options, options.Validate()
}

const (
	HeaderHost          = "Host"
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"
	HeaderUserAgent     = "User-Agent"
	HeaderCacheControl  = "Cache-Control"
	HeaderUpgrade       = "Upgrade"

	WailsUserAgentValue = "wails.io"
)

func serveFile(rw http.ResponseWriter, filename string, blob []byte) error {
	header := rw.Header()
	header.Set(HeaderContentLength, strconv.Itoa(len(blob)))
	if mimeType := header.Get(HeaderContentType); mimeType == "" {
		mimeType = GetMimetype(filename, blob)
		header.Set(HeaderContentType, mimeType)
	}

	rw.WriteHeader(http.StatusOK)
	_, err := io.Copy(rw, bytes.NewReader(blob))
	return err
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

func isWebSocket(req *http.Request) bool {
	upgrade := req.Header.Get(HeaderUpgrade)
	return strings.EqualFold(upgrade, "websocket")
}
