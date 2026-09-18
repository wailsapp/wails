package main

import (
	"embed"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

// message is one row of the mailbox shown by the WebView. Everything that
// filters, sorts, or counts these rows is a native AppKit control living in
// the window's titlebar.
type message struct {
	Sender  string `json:"sender"`
	Subject string `json:"subject"`
	Folder  string `json:"folder"`
	Starred bool   `json:"starred"`
	Unread  bool   `json:"unread"`
	Minutes int    `json:"minutes"`
}

var folders = []string{"Inbox", "Starred", "Archive"}

// mailbox holds the demo state and the native accessory handles that reflect
// it. Native callbacks arrive on their own goroutines, so state is guarded.
type mailbox struct {
	app    *application.App
	window *application.WebviewWindow

	lock     sync.Mutex
	messages []message
	folder   int
	query    string
	sortBy   string

	// Titlebar accessories: the folder switcher sits next to the window
	// buttons, the search tools at the trailing edge, and the status strip
	// spans the full width beneath the titlebar.
	folderStrip *application.MacAccessory
	searchStrip *application.MacAccessory
	statusStrip *application.MacAccessory
	status      *application.MacAccessoryControl
	markRead    *application.MacAccessoryControl
	search      *application.MacAccessoryControl
}

func main() {
	app := application.New(application.Options{
		Name:        "Postbox",
		Description: "Native titlebar accessories built from AppKit controls",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	box := &mailbox{app: app, sortBy: "newest", messages: []message{
		{Sender: "Mara Quill", Subject: "Saturday plans, slowly", Folder: "Inbox", Unread: true, Minutes: 4},
		{Sender: "Field Notes", Subject: "Issue 42: paying attention", Folder: "Inbox", Starred: true, Minutes: 31},
		{Sender: "Tomás Reyes", Subject: "Peaches by scent, not colour", Folder: "Inbox", Unread: true, Minutes: 58},
		{Sender: "Bakery on Elm", Subject: "Your loaf is ready", Folder: "Inbox", Minutes: 130},
		{Sender: "Library", Subject: "A smaller promise (due Friday)", Folder: "Inbox", Starred: true, Unread: true, Minutes: 400},
		{Sender: "Mara Quill", Subject: "Re: the long way home", Folder: "Archive", Minutes: 1440},
		{Sender: "City Water", Subject: "Neighbourhood watering schedule", Folder: "Archive", Starred: true, Minutes: 2880},
	}}

	box.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Postbox",
		Width:  760,
		Height: 480,
		URL:    "/",
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropNormal,
			TitleBar: application.MacTitleBar{
				FullSizeContent: false,
			},
		},
	})

	box.buildAccessories()
	box.installMenu()

	app.Event.On("mail:ready", func(*application.CustomEvent) { box.publish() })
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// buildAccessories creates the three native strips and attaches them before
// the window exists; Wails queues them and installs them with the window.
func (m *mailbox) buildAccessories() {
	// Leading: a segmented folder switcher with SF Symbols.
	m.folderStrip = application.NewMacAccessory(application.MacAccessoryLayoutLeading)
	m.folderStrip.AddSegmented(folders, 0).
		SetSegmentSymbols("tray", "star", "archivebox").
		SetTooltip("Choose a folder").
		OnSelectionChange(func(_ *application.Context, index int, _ string) {
			m.lock.Lock()
			m.folder = index
			m.lock.Unlock()
			m.publish()
		})

	// Trailing: an incremental search field, a sort menu, and a compose button.
	m.searchStrip = application.NewMacAccessory(application.MacAccessoryLayoutTrailing)
	m.search = m.searchStrip.AddSearch("Search mail").
		SetIncremental(true).
		SetWidth(200).
		OnSearch(func(_ *application.Context, query string) {
			m.lock.Lock()
			m.query = strings.ToLower(strings.TrimSpace(query))
			m.lock.Unlock()
			m.publish()
		})
	sortMenu := application.NewMenu()
	for _, option := range []struct{ label, key string }{
		{"Newest First", "newest"}, {"Oldest First", "oldest"}, {"By Sender", "sender"},
	} {
		key := option.key
		sortMenu.AddRadio(option.label, key == "newest").OnClick(func(*application.Context) {
			m.lock.Lock()
			m.sortBy = key
			m.lock.Unlock()
			m.publish()
		})
	}
	m.searchStrip.AddMenuButton("", sortMenu).
		SetSymbol("arrow.up.arrow.down").
		SetTooltip("Sort messages")
	m.searchStrip.AddSymbolButton("square.and.pencil").
		SetTooltip("Compose a message").
		OnClick(func(*application.Context) { m.compose() })

	// Bottom: a full-width status strip with a secondary label, a flexible
	// space, and actions pinned to the trailing edge.
	m.statusStrip = application.NewMacAccessory(application.MacAccessoryLayoutBottom).
		SetHeight(30).
		SetPreferredScrollEdgeEffectStyle(application.MacScrollEdgeEffectStyleSoft)
	m.status = m.statusStrip.AddLabel("").SetSymbol("envelope")
	m.statusStrip.AddFlexibleSpace()
	m.markRead = m.statusStrip.AddButton("Mark All Read").
		SetTooltip("Mark every visible message as read").
		OnClick(func(*application.Context) { m.markAllRead() })
	m.statusStrip.AddButton("Clear Search").
		SetTooltip("Reset the search field").
		OnClick(func(*application.Context) {
			m.search.SetText("")
			m.lock.Lock()
			m.query = ""
			m.lock.Unlock()
			m.publish()
		})

	for _, accessory := range []*application.MacAccessory{m.folderStrip, m.searchStrip, m.statusStrip} {
		if err := m.window.AddTitlebarAccessory(accessory); err != nil {
			log.Fatal(err)
		}
	}
}

// installMenu adds a View menu that exercises the accessory lifecycle: the
// status strip can be hidden in place, and the search tools can be detached
// and reattached at runtime.
func (m *mailbox) installMenu() {
	menu := m.app.Menu.New()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	view := menu.AddSubmenu("View")
	view.Add("Toggle Status Bar").SetAccelerator("CmdOrCtrl+/").OnClick(func(*application.Context) {
		m.statusStrip.SetHidden(!m.statusStrip.IsHidden())
	})
	view.Add("Toggle Search Tools").SetAccelerator("CmdOrCtrl+Shift+F").OnClick(func(*application.Context) {
		if m.searchStrip.IsAttached() {
			m.searchStrip.Remove()
			return
		}
		if err := m.window.AddTitlebarAccessory(m.searchStrip); err != nil {
			m.window.Error("reattach search tools: %s", err)
		}
	})
	menu.AddRole(application.WindowMenu)
	m.app.Menu.Set(menu)
}

func (m *mailbox) compose() {
	m.lock.Lock()
	m.messages = append([]message{{
		Sender: "You", Subject: "Untitled draft", Folder: "Inbox", Unread: true, Minutes: 0,
	}}, m.messages...)
	m.folder = 0
	m.lock.Unlock()
	m.folderStrip.Controls()[0].SetSelectedSegment(0)
	m.publish()
}

func (m *mailbox) markAllRead() {
	m.lock.Lock()
	visible := m.visibleLocked()
	for _, index := range visible {
		m.messages[index].Unread = false
	}
	m.lock.Unlock()
	m.publish()
}

// visibleLocked returns the indexes of the messages matching the folder and
// query, in display order. The caller holds the lock.
func (m *mailbox) visibleLocked() []int {
	var visible []int
	for index, item := range m.messages {
		switch m.folder {
		case 0:
			if item.Folder != "Inbox" {
				continue
			}
		case 1:
			if !item.Starred {
				continue
			}
		case 2:
			if item.Folder != "Archive" {
				continue
			}
		}
		if m.query != "" && !strings.Contains(strings.ToLower(item.Sender+" "+item.Subject), m.query) {
			continue
		}
		visible = append(visible, index)
	}
	sort.SliceStable(visible, func(i, j int) bool {
		a, b := m.messages[visible[i]], m.messages[visible[j]]
		switch m.sortBy {
		case "oldest":
			return a.Minutes > b.Minutes
		case "sender":
			return a.Sender < b.Sender
		default:
			return a.Minutes < b.Minutes
		}
	})
	return visible
}

// publish pushes the visible messages to the WebView and refreshes the
// native status label and button state.
func (m *mailbox) publish() {
	m.lock.Lock()
	visible := m.visibleLocked()
	rows := make([]message, 0, len(visible))
	unread := 0
	for _, index := range visible {
		rows = append(rows, m.messages[index])
		if m.messages[index].Unread {
			unread++
		}
	}
	folder := folders[m.folder]
	m.lock.Unlock()

	payload, err := json.Marshal(map[string]any{"folder": folder, "messages": rows})
	if err == nil {
		m.app.Event.Emit("mail:state", string(payload))
	}
	status := folderSummary(folder, len(rows), unread)
	m.status.SetText(status)
	m.markRead.SetEnabled(unread > 0)
}

func folderSummary(folder string, shown, unread int) string {
	noun := "messages"
	if shown == 1 {
		noun = "message"
	}
	summary := folder + ": " + strconv.Itoa(shown) + " " + noun
	if unread > 0 {
		summary += ", " + strconv.Itoa(unread) + " unread"
	}
	return summary
}
