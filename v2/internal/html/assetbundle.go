package html

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/assetdb"
	"golang.org/x/net/html"
)

// AssetBundle is a collection of Assets
type AssetBundle struct {
	assets        []*Asset
	basedirectory string
}

// NewAssetBundle creates a new AssetBundle struct containing
// the given html and all the assets referenced by it
func NewAssetBundle(pathToHTML string) (*AssetBundle, error) {

	// Create result
	result := &AssetBundle{
		basedirectory: filepath.Dir(pathToHTML),
	}

	err := result.loadAssets(pathToHTML)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// loadAssets processes the given html file and loads in
// all referenced assets
func (a *AssetBundle) loadAssets(pathToHTML string) error {

	// Save HTML
	htmlAsset := &Asset{
		Type: AssetTypes.HTML,
		Path: filepath.Base(pathToHTML),
	}
	err := htmlAsset.Load(a.basedirectory)
	if err != nil {
		return err
	}
	a.assets = append(a.assets, htmlAsset)

	return a.processHTML(htmlAsset.AsString())
}

// Credit to: https://drstearns.github.io/tutorials/tokenizing/
func (a *AssetBundle) processHTML(htmldata string) error {

	// Tokenize the html
	buf := bytes.NewBufferString(htmldata)
	tokenizer := html.NewTokenizer(buf)

	for {
		//get the next token type
		tokenType := tokenizer.Next()

		//if it's an error token, we either reached
		//the end of the file, or the HTML was malformed
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			//otherwise, there was an error tokenizing,
			//which likely means the HTML was malformed.
			//since this is a simple command-line utility,
			//we can just use log.Fatalf() to report the error
			//and exit the process with a non-zero status code
			return tokenizer.Err()
		}

		//process the token according to the token type...
		if tokenType == html.StartTagToken {
			//get the token
			token := tokenizer.Token()

			//if the name of the element is "title"
			if "link" == token.Data {
				//the next token should be the page title
				tokenType = tokenizer.Next()
				//just make sure it's actually a text token
				asset := &Asset{}
				for _, attr := range token.Attr {
					// Favicon
					if attr.Key == "rel" && attr.Val == "icon" {
						asset.Type = AssetTypes.FAVICON
					}
					if attr.Key == "href" {
						asset.Path = attr.Val
					}
					// stylesheet
					if attr.Key == "rel" && attr.Val == "stylesheet" {
						asset.Type = AssetTypes.CSS
					}
				}
				err := asset.Load(a.basedirectory)
				if err != nil {
					return err
				}
				a.assets = append(a.assets, asset)
			}
			if "script" == token.Data {

				tokenType = tokenizer.Next()
				//just make sure it's actually a text token
				asset := &Asset{Type: AssetTypes.JS}
				for _, attr := range token.Attr {
					if attr.Key == "src" {
						asset.Path = attr.Val
						break
					}
				}
				err := asset.Load(a.basedirectory)
				if err != nil {
					return err
				}
				a.assets = append(a.assets, asset)
			}
		}
	}

	return nil
}

// WriteToCFile dumps all the assets to C files in the given directory
func (a *AssetBundle) WriteToCFile(targetDir string) (string, error) {

	// Write out the assets.c file
	var cdata strings.Builder

	// Write header
	header := `// assets.c
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL.
// This file was auto-generated. DO NOT MODIFY.

`
	cdata.WriteString(header)

	// Loop over the Assets
	var err error
	assetVariables := slicer.String()
	var variableName string
	for index, asset := range a.assets {
		variableName = fmt.Sprintf("%s%d", asset.Type, index)
		assetCdata := fmt.Sprintf("const unsigned char %s[]={ %s0x00 };\n", variableName, asset.AsCHexData())
		cdata.WriteString(assetCdata)
		assetVariables.Add(variableName)
	}

	if assetVariables.Length() > 0 {
		cdata.WriteString(fmt.Sprintf("\nconst unsigned char *assets[] = { %s, 0x00 };", assetVariables.Join(", ")))
	} else {
		cdata.WriteString("\nconst unsigned char *assets[] = { 0x00 };")
	}

	// Save file
	assetsFile := filepath.Join(targetDir, "assets.c")
	err = ioutil.WriteFile(assetsFile, []byte(cdata.String()), 0600)
	if err != nil {
		return "", err
	}
	return assetsFile, nil
}

// ConvertToAssetDB returns an assetdb.AssetDB initialized with
// the items in the AssetBundle
func (a *AssetBundle) ConvertToAssetDB() (*assetdb.AssetDB, error) {
	assetdb := assetdb.NewAssetDB()

	// Loop over the Assets
	for _, asset := range a.assets {
		assetdb.AddAsset(asset.Path, []byte(asset.Data))
	}

	return assetdb, nil
}

func (a *AssetBundle) Dump() {
	println("Assets:")
	for _, asset := range a.assets {
		asset.Dump()
	}
}
