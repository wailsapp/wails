// Package hostinfo collects the hardware fields used by the doctor reports.
package hostinfo

type Processor struct {
	Model string
	Cores int
}

type GraphicsCard struct {
	Name   string
	Vendor string
	Driver string
}
