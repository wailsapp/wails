//go:build !go1.18
// +build !go1.18

package doctor

import "text/tabwriter"

func printBuildSettings(_ *tabwriter.Writer) {
	return
}
