package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/pterm/pterm"
)

type LinuxCapabilities struct {
	GTK4Available        bool   `json:"gtk4_available"`
	GTK3Available        bool   `json:"gtk3_available"`
	WebKitGTK6Available  bool   `json:"webkitgtk_6_available"`
	WebKit2GTK4Available bool   `json:"webkit2gtk_4_1_available"`
	Recommended          string `json:"recommended"`
}

type Capabilities struct {
	Platform string             `json:"platform"`
	Arch     string             `json:"arch"`
	Linux    *LinuxCapabilities `json:"linux,omitempty"`
}

type ToolCapabilitiesOptions struct{}

func ToolCapabilities(_ *ToolCapabilitiesOptions) error {
	caps := Capabilities{
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	switch runtime.GOOS {
	case "linux":
		caps.Linux = detectLinuxCapabilities()
	}

	pterm.Println(capsToJSON(caps))
	return nil
}

func detectLinuxCapabilities() *LinuxCapabilities {
	caps := &LinuxCapabilities{}

	caps.GTK4Available = pkgConfigExists("gtk4")
	caps.WebKitGTK6Available = pkgConfigExists("webkitgtk-6.0")
	caps.GTK3Available = pkgConfigExists("gtk+-3.0")
	caps.WebKit2GTK4Available = pkgConfigExists("webkit2gtk-4.1")

	if caps.GTK4Available && caps.WebKitGTK6Available {
		caps.Recommended = "gtk4"
	} else if caps.GTK3Available && caps.WebKit2GTK4Available {
		caps.Recommended = "gtk3"
	} else {
		caps.Recommended = "none"
	}

	return caps
}

func pkgConfigExists(pkg string) bool {
	cmd := exec.Command("pkg-config", "--exists", pkg)
	return cmd.Run() == nil
}

func capsToJSON(caps Capabilities) string {
	result := fmt.Sprintf(`{"platform":"%s","arch":"%s"`, caps.Platform, caps.Arch)
	if caps.Linux != nil {
		l := caps.Linux
		result += fmt.Sprintf(`,"linux":{"gtk4_available":%t,"gtk3_available":%t,"webkitgtk_6_available":%t,"webkit2gtk_4_1_available":%t,"recommended":"%s"}`,
			l.GTK4Available, l.GTK3Available, l.WebKitGTK6Available, l.WebKit2GTK4Available, l.Recommended)
	}
	result += "}"
	return result
}
