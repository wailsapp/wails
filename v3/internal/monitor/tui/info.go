package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
)

// infoView is a simple dashboard of app metadata and registered event listeners.
type infoView struct {
	*components.ComponentBase
	m     *Model
	label *components.Label
}

func newInfoView(m *Model) *infoView {
	v := &infoView{m: m}
	v.label = components.NewLabel("").SetDynamicColors(true).SetScrollable(true)

	v.ComponentBase = components.NewComponentBase(
		components.NewPanel().SetTitle("App Info").SetContent(v.label).SetFocused(true)).
		SetName("Info").
		AddHint("r", "Refresh").
		AddHint("q", "Back").
		SetOnStart(v.rebuild)
	v.ComponentBase.SetInputHandler(func(ev *tcell.EventKey) bool {
		if ev.Key() == tcell.KeyRune && ev.Rune() == 'r' {
			v.m.requestSnapshot()
			v.m.toasts.Info("Refreshing…")
			return true
		}
		return false
	})
	return v
}

func (v *infoView) rebuild() {
	if v.m.snapshot == nil {
		v.label.SetText("[gray]Waiting for snapshot…[-]")
		return
	}
	a := v.m.snapshot.App

	var b strings.Builder
	b.WriteString("[::b]Application[-:-:-]\n")
	row(&b, "Name", a.Name)
	row(&b, "PID", fmt.Sprintf("%d", a.PID))
	row(&b, "Platform", a.Platform)
	row(&b, "Debug mode", fmt.Sprintf("%v", a.DebugMode))
	if a.Transport != "" {
		row(&b, "Transport", a.Transport)
	}
	if a.DevServer != "" {
		row(&b, "Dev server", a.DevServer)
	}
	row(&b, "Windows", fmt.Sprintf("%d", len(v.m.snapshot.Windows)))
	row(&b, "Bindings", fmt.Sprintf("%d", len(v.m.snapshot.Bindings)))

	b.WriteString("\n[::b]Backend Event Listeners[-:-:-]\n")
	if len(v.m.snapshot.Listeners) == 0 {
		b.WriteString("[gray]  none registered[-]\n")
	} else {
		for _, l := range v.m.snapshot.Listeners {
			fmt.Fprintf(&b, "  %-40s [gray]×%d[-]\n", escapeTags(l.Name), l.Count)
		}
	}

	v.label.SetText(strings.TrimRight(b.String(), "\n"))
}

func row(b *strings.Builder, k, val string) {
	fmt.Fprintf(b, "  [gray]%-12s[-] %s\n", k, escapeTags(val))
}
