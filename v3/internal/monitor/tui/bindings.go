package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
	"github.com/atterpac/dado/core"

	monitor "github.com/wailsapp/wails/v3/internal/monitor"
)

// bindingsView is a static explorer of the bound methods the frontend can call:
// a master list on the left, signature + doc detail on the right.
type bindingsView struct {
	*components.ComponentBase
	m *Model

	split   *components.Split
	table   *components.Table
	detail  *components.Label
	body    *components.CodeView
	visible []monitor.BindingInfo
}

func newBindingsView(m *Model) *bindingsView {
	v := &bindingsView{m: m}

	v.table = components.NewTable()
	v.table.SetHeaders("METHOD", "SERVICE", "IN", "OUT")
	v.table.ConfigureEmpty("ƒ", "No bindings", "Snapshot not yet received")

	v.detail = components.NewLabel("").SetDynamicColors(true)
	v.body = components.NewCodeView()
	v.body.SetLanguage(components.LangGo).SetShowLineNumbers(false)

	detailPane := core.NewFlex().SetDirection(core.Column).
		AddItem(v.detail, 4, 0, false).
		AddItem(v.body, 0, 1, false)

	v.split = components.NewSplit()
	v.split.SetDirection(components.SplitHorizontal).SetRatio(0.5).
		SetLeft(v.table).SetRight(detailPane)

	v.table.SetSelectionChangedFunc(func(row, col int) { v.showDetail(row) })

	v.ComponentBase = components.NewComponentBase(v.split).
		SetName("Bindings").
		AddHint("j/k", "Move").
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

func (v *bindingsView) rebuild() {
	v.table.ClearRows()
	v.visible = nil
	if v.m.snapshot == nil {
		return
	}
	v.visible = v.m.snapshot.Bindings
	for _, b := range v.visible {
		v.table.AddRow(
			b.Name,
			b.Service,
			fmt.Sprintf("%d", len(b.Inputs)),
			fmt.Sprintf("%d", len(b.Outputs)),
		)
	}
}

func (v *bindingsView) showDetail(tableRow int) {
	idx := tableRow - 1
	if idx < 0 || idx >= len(v.visible) {
		v.detail.SetText("")
		v.body.SetCode("")
		return
	}
	b := v.visible[idx]
	v.detail.SetText(fmt.Sprintf(
		"[::b]%s[-:-:-]\n[gray]fqn[-] %s\n[gray]id[-] %d",
		escapeTags(b.Name), escapeTags(b.FQN), b.ID,
	))

	var sb strings.Builder
	sb.WriteString(signature(b))
	if b.Comment != "" {
		sb.WriteString("\n\n// ")
		sb.WriteString(strings.ReplaceAll(b.Comment, "\n", "\n// "))
	}
	v.body.SetCode(sb.String())
}

// signature renders a Go-like method signature from the binding params.
func signature(b monitor.BindingInfo) string {
	var in []string
	for _, p := range b.Inputs {
		if p.Name != "" {
			in = append(in, p.Name+" "+p.Type)
		} else {
			in = append(in, p.Type)
		}
	}
	var out []string
	for _, p := range b.Outputs {
		out = append(out, p.Type)
	}
	sig := fmt.Sprintf("func %s(%s)", b.Name, strings.Join(in, ", "))
	switch len(out) {
	case 0:
	case 1:
		sig += " " + out[0]
	default:
		sig += " (" + strings.Join(out, ", ") + ")"
	}
	return sig
}
