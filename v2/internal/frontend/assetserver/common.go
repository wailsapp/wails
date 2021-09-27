package assetserver

import (
	"bytes"
	"fmt"
	"strings"
)

func injectHTML(input string, html string) ([]byte, error) {
	splits := strings.Split(input, "</body>")
	if len(splits) != 2 {
		return nil, fmt.Errorf("unable to locate a </body> tag in your html")
	}

	var result bytes.Buffer
	result.WriteString(splits[0])
	result.WriteString(html)
	result.WriteString("</body>")
	result.WriteString(splits[1])
	return result.Bytes(), nil
}
