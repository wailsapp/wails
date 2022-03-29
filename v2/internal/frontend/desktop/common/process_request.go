package common

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
)

type RequestRespone struct {
	Body       []byte
	MimeType   string
	StatusCode int
}

func (r RequestRespone) StatusText() string {
	return http.StatusText(r.StatusCode)
}

func (r RequestRespone) String() string {
	return fmt.Sprintf("Body: '%s', StatusCode: %d", string(r.Body), r.StatusCode)
}

func ProcessRequest(uri string, assets *assetserver.DesktopAssetServer, expectedScheme string, expectedHosts ...string) (RequestRespone, error) {
	// Translate URI to file
	file, err := translateUriToFile(uri, expectedScheme, expectedHosts...)
	if err != nil {
		if err == ErrUnexpectedHost {
			body := fmt.Sprintf("expected host one of \"%s\"", strings.Join(expectedHosts, ","))
			return textResponse(body, http.StatusInternalServerError), err
		}

		return RequestRespone{StatusCode: http.StatusInternalServerError}, err
	}

	content, mimeType, err := assets.Load(file)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if os.IsNotExist(err) {
			statusCode = http.StatusNotFound
		}
		return RequestRespone{StatusCode: statusCode}, err
	}

	return RequestRespone{Body: content, MimeType: mimeType, StatusCode: http.StatusOK}, nil
}

func textResponse(body string, statusCode int) RequestRespone {
	if body == "" {
		return RequestRespone{StatusCode: statusCode}
	}
	return RequestRespone{Body: []byte(body), MimeType: "text/plain;charset=UTF-8", StatusCode: statusCode}
}
