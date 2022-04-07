//go:build go1.18
// +build go1.18

package doctor

import (
	"fmt"
	"runtime/debug"
	"text/tabwriter"
)

func printBuildSettings(w *tabwriter.Writer) {
	if buildInfo, _ := debug.ReadBuildInfo(); buildInfo != nil {
		buildSettingToName := map[string]string{
			"vcs.revision": "Revision",
			"vcs.modified": "Modified",
		}
		for _, buildSetting := range buildInfo.Settings {
			name := buildSettingToName[buildSetting.Key]
			if name == "" {
				continue
			}

			fmt.Fprintf(w, "%s:\t%s\n", name, buildSetting.Value)
		}
	}
}
