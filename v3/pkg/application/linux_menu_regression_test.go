package application

import (
	"os"
	"strings"
	"testing"
)

func TestLinuxGTK4MenuUpdateRebuildsProcessedMenus(t *testing.T) {
	data, err := os.ReadFile("menu_linux.go")
	if err != nil {
		t.Skip("menu_linux.go not available")
	}

	body := functionBody(string(data), "func (m *linuxMenu) processMenu(menu *Menu)")

	if strings.Contains(body, "if impl.processed") && strings.Contains(body, "return") {
		t.Fatal("GTK4 processMenu must rebuild processed menus on Update, not return early")
	}

	if !strings.Contains(body, "menuClear(menu)") {
		t.Fatal("GTK4 processMenu must clear existing native items before rebuilding")
	}

	data, err = os.ReadFile("linux_cgo.go")
	if err != nil {
		t.Skip("linux_cgo.go not available")
	}
	if !strings.Contains(string(data), "func menuClear(menu *Menu)") {
		t.Fatal("GTK4 cgo path must provide menuClear for processMenu rebuilds")
	}
}

func functionBody(source, signature string) string {
	start := strings.Index(source, signature)
	if start == -1 {
		return ""
	}

	source = source[start:]
	depth := 0
	for i, r := range source {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return source[:i+1]
			}
		}
	}

	return source
}
