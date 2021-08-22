package assetserver

import (
	"bytes"
	"fmt"
	"strings"
)

func injectScript(input string, script string) ([]byte, error) {
	splits := strings.Split(input, "<head>")
	if len(splits) != 2 {
		return nil, fmt.Errorf("unable to locate a </head> tag in your html")
	}

	var result bytes.Buffer
	result.WriteString(splits[0])
	result.WriteString("<head>")
	result.WriteString(script)
	result.WriteString(splits[1])
	return result.Bytes(), nil
}
