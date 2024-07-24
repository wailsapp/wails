package runtime

import "github.com/wailsapp/wails/v2/internal/system"

func SystemReport() (*system.Report, error) {
	return system.GenerateReport()
}
