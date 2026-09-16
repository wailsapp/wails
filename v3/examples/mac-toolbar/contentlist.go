package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// daymarkContentList is the middle column: a native content-list table of
// the notes, in Mail's rich-row style. Each row shows the note title, the
// first line of its body as the subtitle, the modification time as the
// trailing detail, and the priority as a badge. Selecting a row opens the
// note in the editor; a context menu offers Pin, Delete, and sort orders.
//
// Rows are keyed by the sidebar handle of the same note, so drag reorders,
// deletions, and new notes in the sidebar reconcile into the list without
// rebuilding it.
type daymarkContentList struct {
	split *daymarkSplit
	list  *application.MacContentList
	pane  *application.MacSplitPane

	entries []*daymarkListEntry
	byItem  map[*application.MacSidebarItem]*daymarkListEntry
	byRow   map[*application.MacContentListRow]*daymarkListEntry

	// sortMode is "manual" (sidebar order), "title", or "date".
	sortMode      string
	sortAscending bool
	syncing       bool
}

type daymarkListEntry struct {
	item      *application.MacSidebarItem
	row       *application.MacContentListRow
	modified  time.Time
	signature string
	title     string
}

func newDaymarkContentList(split *daymarkSplit) *daymarkContentList {
	result := &daymarkContentList{
		split:         split,
		list:          application.NewMacContentList(),
		byItem:        map[*application.MacSidebarItem]*daymarkListEntry{},
		byRow:         map[*application.MacContentListRow]*daymarkListEntry{},
		sortMode:      "manual",
		sortAscending: true,
	}
	result.list.
		SetStyle(application.MacContentListStyleInset).
		SetEmptyText("No notes match").
		OnSelectionChange(result.handleSelection).
		OnActivate(result.handleActivate).
		OnContextMenu(result.contextMenuForRow)

	// The empty area offers New Note and the sort orders.
	emptyArea := application.NewMenu()
	emptyArea.Add("New Note").OnClick(func(*application.Context) { split.NewNote() })
	emptyArea.AddSeparator()
	result.addSortItems(emptyArea)
	result.list.SetContextMenu(emptyArea)
	return result
}

func (c *daymarkContentList) NativeContentList() *application.MacContentList { return c.list }

// addSortItems appends the sort-order choices shared by both context menus.
func (c *daymarkContentList) addSortItems(menu *application.Menu) {
	sortMenu := menu.AddSubmenu("Sort By")
	sortMenu.Add("Manual Order").SetChecked(c.sortMode == "manual").OnClick(func(*application.Context) {
		c.setSort("manual", true)
	})
	sortMenu.Add("Title").SetChecked(c.sortMode == "title").OnClick(func(*application.Context) {
		c.setSort("title", true)
	})
	sortMenu.Add("Date Modified").SetChecked(c.sortMode == "date").OnClick(func(*application.Context) {
		c.setSort("date", false)
	})
}

func (c *daymarkContentList) setSort(mode string, ascending bool) {
	c.split.stateLock.Lock()
	c.sortMode = mode
	c.sortAscending = ascending
	c.split.stateLock.Unlock()
	c.split.syncContentList()
}

// contextMenuForRow builds the right-click menu for one row. The empty area
// returns nil so the list's fallback menu appears instead.
func (c *daymarkContentList) contextMenuForRow(_ *application.Context, row *application.MacContentListRow) *application.Menu {
	if row == nil {
		return nil
	}
	c.split.stateLock.Lock()
	entry := c.byRow[row]
	var item *application.MacSidebarItem
	var note daymarkNote
	found := false
	if entry != nil {
		item = entry.item
		if index := c.split.indexOfItemLocked(item); index >= 0 {
			note = c.split.notes[index]
			found = true
		}
	}
	c.split.stateLock.Unlock()
	if !found {
		return nil
	}
	menu := application.NewMenu()
	pinLabel := "Pin Note"
	if note.Pinned {
		pinLabel = "Unpin Note"
	}
	menu.Add(pinLabel).OnClick(func(*application.Context) { c.split.setPinned(item, !note.Pinned) })
	menu.Add("Delete Note").OnClick(func(*application.Context) { c.split.deleteNote(item) })
	menu.AddSeparator()
	c.addSortItems(menu)
	return menu
}

// handleSelection opens the selected note in the editor and mirrors the
// choice into the sidebar.
func (c *daymarkContentList) handleSelection(_ *application.Context, rows []*application.MacContentListRow) {
	if len(rows) == 0 {
		return
	}
	c.split.stateLock.Lock()
	if c.syncing {
		c.split.stateLock.Unlock()
		return
	}
	entry := c.byRow[rows[0]]
	index := -1
	if entry != nil {
		index = c.split.indexOfItemLocked(entry.item)
	}
	if index < 0 {
		c.split.stateLock.Unlock()
		return
	}
	c.split.active = index
	item := entry.item
	c.split.stateLock.Unlock()
	c.split.sidebar.SetSelectedItem(item)
	c.split.pushSidebarStateToInspector()
	c.split.emitActiveNote()
}

// handleActivate reveals the inspector for a double-clicked row.
func (c *daymarkContentList) handleActivate(ctx *application.Context, row *application.MacContentListRow) {
	c.handleSelection(ctx, []*application.MacContentListRow{row})
	c.split.SetInspectorCollapsed(false)
}

// syncContentList reconciles the native rows with the notes: it adds rows
// for new notes, updates changed ones in place, removes rows for deleted
// notes, applies the sort order and filters, and selects the active note.
func (s *daymarkSplit) syncContentList() {
	c := s.contentList
	if c == nil {
		return
	}
	s.stateLock.Lock()
	c.syncing = true
	notes := append([]daymarkNote(nil), s.notes...)
	items := append([]*application.MacSidebarItem(nil), s.items...)
	active, query, category := s.active, s.query, s.category
	sortMode, ascending := c.sortMode, c.sortAscending
	s.stateLock.Unlock()
	defer func() {
		s.stateLock.Lock()
		c.syncing = false
		s.stateLock.Unlock()
	}()

	now := time.Now()
	seen := make(map[*application.MacSidebarItem]bool, len(items))
	entries := make([]*daymarkListEntry, 0, len(items))
	for index, item := range items {
		note := notes[index]
		entry := c.byItem[item]
		if entry == nil {
			entry = &daymarkListEntry{
				item: item,
				row:  c.list.AddRow(note.Title),
				// Seed the demo with staggered times so Date Modified has an order.
				modified: now.Add(-time.Duration(index) * 26 * time.Hour),
			}
			c.byItem[item] = entry
			c.byRow[entry.row] = entry
		}
		signature := fmt.Sprintf("%s\x00%s\x00%s\x00%v\x00%d", note.Title, note.Subtitle, note.Body, note.Pinned, note.Priority)
		if entry.signature != "" && entry.signature != signature {
			entry.modified = now
		}
		entry.signature = signature
		entry.title = note.Title
		haystack := strings.ToLower(note.Title + " " + note.Subtitle + " " + note.Body)
		hidden := (query != "" && !strings.Contains(haystack, query)) ||
			(category != "" && !strings.HasPrefix(note.Category, category))
		entry.row.
			SetTitle(note.Title).
			SetSubtitle(firstBodyLine(note)).
			SetDetail(relativeTime(entry.modified, now)).
			SetSymbol(noteSymbol(note)).
			SetBadge(note.Priority).
			SetTooltip(note.Category).
			SetHidden(hidden)
		seen[item] = true
		entries = append(entries, entry)
	}
	for item, entry := range c.byItem {
		if seen[item] {
			continue
		}
		entry.row.Remove()
		delete(c.byItem, item)
		delete(c.byRow, entry.row)
	}
	c.entries = entries

	position := make(map[*application.MacContentListRow]int, len(entries))
	for index, entry := range entries {
		position[entry.row] = index
	}
	c.list.SortRows(func(a, b *application.MacContentListRow) bool {
		left, right := c.byRow[a], c.byRow[b]
		if left == nil || right == nil {
			return position[a] < position[b]
		}
		switch sortMode {
		case "title":
			l, r := strings.ToLower(left.title), strings.ToLower(right.title)
			if l != r {
				return (l < r) == ascending
			}
		case "date":
			if !left.modified.Equal(right.modified) {
				return left.modified.Before(right.modified) == ascending
			}
		}
		return position[a] < position[b]
	})

	if active >= 0 && active < len(entries) {
		c.list.SetSelectedRow(entries[active].row)
	}
}

// firstBodyLine returns the first non-empty body line, or the subtitle.
func firstBodyLine(note daymarkNote) string {
	for _, line := range strings.Split(note.Body, "\n") {
		line = strings.TrimSpace(strings.TrimLeft(line, "# "))
		if line != "" {
			return line
		}
	}
	return note.Subtitle
}

// relativeTime formats a modification time the way Mail's list does: a clock
// time today, "Yesterday", or a short date.
func relativeTime(when, now time.Time) string {
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	switch {
	case !when.Before(today):
		return when.Format("15:04")
	case !when.Before(today.Add(-24 * time.Hour)):
		return "Yesterday"
	case when.Year() == now.Year():
		return when.Format("2 Jan")
	}
	return when.Format("2 Jan 2006")
}
