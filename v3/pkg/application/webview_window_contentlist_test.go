package application

import (
	"strings"
	"testing"
)

// newRegisteredTestContentList wires a content list between a sidebar and
// the primary pane and registers the pane so pane-scoped native callbacks can
// be exercised without AppKit.
func newRegisteredTestContentList(t *testing.T) (*MacContentList, uint64) {
	t.Helper()
	list := NewMacContentList()
	split := NewMacSplitView()
	split.AddSidebar(NewMacSidebar())
	pane := split.AddContentList(list)
	split.AddPrimaryContent()
	internal := split.paneSnapshot()[1]
	registerMacSplitPane(internal)
	t.Cleanup(func() {
		unregisterMacSplitPane(internal.internalID)
		for _, row := range list.Rows() {
			unregisterMacContentListRow(row.internalID)
		}
	})
	return list, pane.internalID
}

func contentListTitles(rows []*MacContentListRow) string {
	titles := make([]string, 0, len(rows))
	for _, row := range rows {
		titles = append(titles, row.Title())
	}
	return strings.Join(titles, ",")
}

func TestMacContentListSnapshotRichRows(t *testing.T) {
	list := NewMacContentList().
		SetEmptyText("No notes").
		SetStyle(MacContentListStyleInset).
		SetAlternatingRowBackgrounds(true).
		SetRowHeight(52)
	row := list.AddRow("Saturday").
		SetSubtitle("A slow day").
		SetDetail("Yesterday").
		SetSymbol("doc.text").
		SetBadge(2).
		SetTooltip("Open").
		SetEnabled(false)
	list.AddRow("Hidden").SetHidden(true)
	list.SetSelectedRow(row)

	snapshot := list.snapshot()
	if len(snapshot.columns) != 0 || snapshot.headerVisible {
		t.Fatal("rich-row mode has no columns and hides the header by default")
	}
	if snapshot.emptyText != "No notes" || snapshot.style != MacContentListStyleInset ||
		!snapshot.alternatingRows || snapshot.rowHeight != 52 || snapshot.sortColumn != -1 {
		t.Fatalf("unexpected list snapshot: %#v", snapshot)
	}
	if len(snapshot.rows) != 2 {
		t.Fatalf("expected two rows, got %d", len(snapshot.rows))
	}
	first := snapshot.rows[0]
	if first.internalID != row.internalID || first.title != "Saturday" || first.subtitle != "A slow day" ||
		first.detail != "Yesterday" || first.symbol != "doc.text" || first.badge != 2 ||
		first.tooltip != "Open" || !first.disabled || first.hidden {
		t.Fatalf("unexpected row snapshot: %#v", first)
	}
	if !snapshot.rows[1].hidden {
		t.Fatal("hidden rows must reach the snapshot flagged hidden")
	}
	if len(snapshot.selectedIDs) != 1 || snapshot.selectedIDs[0] != row.internalID {
		t.Fatalf("selection should carry the selected row id, got %v", snapshot.selectedIDs)
	}
	if row.SetBadge(-4).Badge() != 0 {
		t.Fatal("negative badges clear the badge")
	}
}

func TestMacContentListSnapshotTableMode(t *testing.T) {
	list := NewMacContentList().SetColumns(
		MacContentListColumn{Title: "Name", Width: 200, MinWidth: 100, Sortable: true},
		MacContentListColumn{Title: "Size", Width: 80, Alignment: MacContentListAlignTrailing},
	)
	list.AddRow("report.pdf").SetCells("report.pdf", "12 KB")
	list.AddRow("untitled")
	list.SetSortable(true).SetHeaderVisible(false)

	snapshot := list.snapshot()
	if len(snapshot.columns) != 2 || snapshot.columns[0].Title != "Name" || !snapshot.columns[0].Sortable ||
		snapshot.columns[1].Alignment != MacContentListAlignTrailing {
		t.Fatalf("unexpected columns: %#v", snapshot.columns)
	}
	if snapshot.headerVisible {
		t.Fatal("SetHeaderVisible(false) must override the table-mode default")
	}
	if !snapshot.sortable {
		t.Fatal("sortable flag should reach the snapshot")
	}
	if got := snapshot.rows[0].cells; len(got) != 2 || got[1] != "12 KB" {
		t.Fatalf("unexpected cells: %v", got)
	}
	if len(snapshot.rows[1].cells) != 0 || snapshot.rows[1].title != "untitled" {
		t.Fatal("rows without cells keep their title for the first column")
	}
	if len(list.Columns()) != 2 {
		t.Fatal("Columns should return the configured columns")
	}
	list.SetColumns()
	if list.snapshot().headerVisible {
		t.Fatal("hiding the header explicitly persists after returning to rich-row mode")
	}
}

func TestMacContentListInsertRemoveOrdering(t *testing.T) {
	list := NewMacContentList()
	b := list.AddRow("b")
	d := list.AddRow("d")
	a := list.InsertRow(0, "a")
	c := list.InsertRow(2, "c")
	e := list.InsertRow(99, "e")
	if got := contentListTitles(list.Rows()); got != "a,b,c,d,e" {
		t.Fatalf("rows = %s", got)
	}
	list.SetAllowsMultipleSelection(true)
	list.SetSelectedRows([]*MacContentListRow{e, c, a})
	if got := contentListTitles(list.SelectedRows()); got != "a,c,e" {
		t.Fatalf("selection should follow display order, got %s", got)
	}

	c.Remove()
	if got := contentListTitles(list.Rows()); got != "a,b,d,e" {
		t.Fatalf("rows after remove = %s", got)
	}
	if got := contentListTitles(list.SelectedRows()); got != "a,e" {
		t.Fatalf("removed rows leave the selection, got %s", got)
	}
	if c.SetTitle("zombie").Title() != "c" {
		t.Fatal("removed rows are inert")
	}
	e.SetHidden(true)
	if got := contentListTitles(list.SelectedRows()); got != "a" {
		t.Fatalf("hidden rows leave the selection, got %s", got)
	}
	if list.SelectedRow() != a {
		t.Fatal("SelectedRow returns the first selected row")
	}

	list.RemoveAll()
	if len(list.Rows()) != 0 || len(list.SelectedRows()) != 0 {
		t.Fatal("RemoveAll clears rows and selection")
	}
	if b.SetTitle("x").Title() != "b" || d.SetSubtitle("x").Subtitle() != "" {
		t.Fatal("handles are inert after RemoveAll")
	}
	if list.InsertRow(0, "fresh") == nil {
		t.Fatal("the list stays usable after RemoveAll")
	}
}

func TestMacContentListSelectionCallbacksViaRegistry(t *testing.T) {
	list, paneID := newRegisteredTestContentList(t)
	first := list.AddRow("First")
	second := list.AddRow("Second")
	other := NewMacContentList().AddRow("Elsewhere")
	list.registerRows()
	registerMacContentListRow(other)
	t.Cleanup(func() { unregisterMacContentListRow(other.internalID) })

	var selection []*MacContentListRow
	calls := 0
	list.SetAllowsMultipleSelection(true).OnSelectionChange(func(_ *Context, rows []*MacContentListRow) {
		calls++
		selection = rows
	})
	activated := 0
	var activatedRow *MacContentListRow
	list.OnActivate(func(_ *Context, row *MacContentListRow) {
		activated++
		activatedRow = row
	})

	handleMacContentListSelectionChanged(macContentListSelectionEvent{
		paneID: paneID,
		rowIDs: []uint64{first.internalID, other.internalID, second.internalID, 424242},
	})
	if calls != 1 || contentListTitles(selection) != "First,Second" {
		t.Fatalf("calls=%d selection=%s", calls, contentListTitles(selection))
	}
	if contentListTitles(list.SelectedRows()) != "First,Second" {
		t.Fatal("the list should mirror the native selection")
	}

	handleMacContentListSelectionChanged(macContentListSelectionEvent{paneID: paneID})
	if calls != 2 || len(selection) != 0 || list.SelectedRow() != nil {
		t.Fatal("clearing the selection should notify with an empty slice")
	}

	handleMacContentListRowActivated(macContentListActivateEvent{paneID: paneID, rowID: second.internalID})
	handleMacContentListRowActivated(macContentListActivateEvent{paneID: paneID, rowID: other.internalID})
	handleMacContentListRowActivated(macContentListActivateEvent{paneID: paneID + 1000, rowID: first.internalID})
	if activated != 1 || activatedRow != second {
		t.Fatalf("activated=%d row=%v", activated, activatedRow)
	}

	second.Remove()
	handleMacContentListSelectionChanged(macContentListSelectionEvent{paneID: paneID, rowIDs: []uint64{second.internalID}})
	if calls != 3 || len(selection) != 0 {
		t.Fatal("removed rows must not resolve through the registry")
	}
}

func TestMacContentListSortCallbacks(t *testing.T) {
	list, paneID := newRegisteredTestContentList(t)
	list.SetColumns(
		MacContentListColumn{Title: "Name", Sortable: true},
		MacContentListColumn{Title: "Size", Sortable: true},
	).SetSortable(true)
	list.AddRow("b").SetCells("banana", "20")
	list.AddRow("a").SetCells("Apple", "5")
	list.AddRow("c").SetCells("cherry", "10")

	// Without OnSort the list sorts itself by the clicked column.
	handleMacContentListSortChanged(macContentListSortEvent{paneID: paneID, column: 0, ascending: true})
	if got := contentListTitles(list.Rows()); got != "a,b,c" {
		t.Fatalf("default sort by name = %s", got)
	}
	if column, ascending := list.SortOrder(); column != 0 || !ascending {
		t.Fatalf("sort order = %d/%v", column, ascending)
	}
	handleMacContentListSortChanged(macContentListSortEvent{paneID: paneID, column: 0, ascending: false})
	if got := contentListTitles(list.Rows()); got != "c,b,a" {
		t.Fatalf("descending sort by name = %s", got)
	}

	// With OnSort the callback owns the ordering; the indicator is recorded.
	var sorted []int
	list.OnSort(func(_ *Context, column int, ascending bool) {
		sorted = append(sorted, column)
		list.SortRows(func(a, b *MacContentListRow) bool {
			return a.Cells()[1] < b.Cells()[1]
		})
	})
	handleMacContentListSortChanged(macContentListSortEvent{paneID: paneID, column: 1, ascending: true})
	if len(sorted) != 1 || sorted[0] != 1 {
		t.Fatalf("OnSort calls = %v", sorted)
	}
	if got := contentListTitles(list.Rows()); got != "c,b,a" {
		t.Fatalf("SortRows ordering = %s", got)
	}
	if column, _ := list.SortOrder(); column != 1 {
		t.Fatal("the clicked column should be recorded before OnSort runs")
	}
	if snapshot := list.snapshot(); snapshot.sortColumn != 1 || !snapshot.sortAscending {
		t.Fatalf("snapshot sort state = %d/%v", snapshot.sortColumn, snapshot.sortAscending)
	}
	handleMacContentListSortChanged(macContentListSortEvent{paneID: paneID + 1000, column: 0, ascending: true})
	if len(sorted) != 1 {
		t.Fatal("unknown panes must not reach OnSort")
	}

	// Rich-row mode sorts by title, subtitle, or detail.
	rich := NewMacContentList()
	rich.AddRow("Zed").SetSubtitle("2")
	rich.AddRow("Amy").SetSubtitle("1")
	rich.AddRow("Kim").SetSubtitle("3")
	rich.SortBy(1, false)
	if got := contentListTitles(rich.Rows()); got != "Kim,Zed,Amy" {
		t.Fatalf("SortBy subtitle descending = %s", got)
	}
	rich.SortBy(-1, true)
	if column, _ := rich.SortOrder(); column != -1 || contentListTitles(rich.Rows()) != "Kim,Zed,Amy" {
		t.Fatal("a negative column clears the indicator without reordering")
	}
}

func TestMacContentListContextMenuResolution(t *testing.T) {
	list, paneID := newRegisteredTestContentList(t)
	first := list.AddRow("First")
	second := list.AddRow("Second")
	list.registerRows()

	rowMenu := NewMenu()
	dynamicMenu := NewMenu()
	fallback := NewMenu()
	first.SetContextMenu(rowMenu)
	list.SetContextMenu(fallback)
	var asked *MacContentListRow
	askedCalls := 0
	list.OnContextMenu(func(_ *Context, row *MacContentListRow) *Menu {
		askedCalls++
		asked = row
		if row == second {
			return dynamicMenu
		}
		return nil
	})

	if resolveMacContentListContextMenu(paneID, first.internalID) != rowMenu || askedCalls != 0 {
		t.Fatal("a row's own menu should win without consulting the callback")
	}
	if resolveMacContentListContextMenu(paneID, second.internalID) != dynamicMenu || asked != second {
		t.Fatal("the callback should receive the clicked row")
	}
	if resolveMacContentListContextMenu(paneID, 0) != fallback || asked != nil {
		t.Fatal("empty-area clicks should pass nil and fall back to the list menu")
	}
	if resolveMacContentListContextMenu(paneID+1000, 0) != nil {
		t.Fatal("unknown panes should show no menu")
	}
}

func TestMacContentListWrongModeAndDeadHandles(t *testing.T) {
	list := NewMacContentList()
	row := list.AddRow("Row")
	if row.SetCells("a", "b").Cells()[1] != "b" {
		t.Fatal("cells are stored even in rich-row mode so switching modes keeps them")
	}
	if list.SetRowHeight(-1).snapshot().rowHeight != 0 {
		t.Fatal("negative row heights are rejected")
	}
	if list.SetStyle(MacContentListStyle(99)).snapshot().style != MacContentListStyleAutomatic {
		t.Fatal("unknown styles are rejected")
	}
	if list.SortRows(nil) != list {
		t.Fatal("a nil comparator is a no-op")
	}
	if list.SetSelectedRow(NewMacContentList().AddRow("other")).SelectedRow() != nil {
		t.Fatal("rows from another list are ignored by SetSelectedRow")
	}
	list.SetSelectedRow(row)
	list.SetAllowsMultipleSelection(false)
	list.SetSelectedRows([]*MacContentListRow{row, list.AddRow("Second")})
	if len(list.SelectedRows()) != 1 {
		t.Fatal("single selection keeps only the first row")
	}

	var nilList *MacContentList
	if nilList.AddRow("x") != nil || nilList.SetColumns() != nil || nilList.Rows() != nil {
		t.Fatal("nil lists are safe")
	}
	var nilRow *MacContentListRow
	nilRow.Remove()
	if nilRow.SetTitle("x") != nil || nilRow.Title() != "" {
		t.Fatal("nil rows are safe")
	}

	split := NewMacSplitView()
	split.AddContentList(list)
	split.AddPrimaryContent()
	pane := split.paneSnapshot()[0]
	pane.markDead()
	if !list.dead || list.pane != nil {
		t.Fatal("markDead should detach the list")
	}
	if list.AddRow("late") != nil || row.SetTitle("late").Title() != "Row" {
		t.Fatal("dead lists reject mutation")
	}
	if list.SetEmptyText("x").snapshot().emptyText != "" {
		t.Fatal("dead lists ignore option changes")
	}
}

func TestMacSplitViewContentListOrderingValidation(t *testing.T) {
	valid := NewMacSplitView()
	valid.AddSidebar(NewMacSidebar())
	valid.AddContentList(NewMacContentList())
	valid.AddPrimaryContent()
	inspector, _, _, _ := newTestInspector()
	valid.AddInspector(inspector)
	if err := validateMacSplitView(valid); err != nil {
		t.Fatalf("sidebar, content list, primary, inspector should validate: %v", err)
	}
	if got := valid.paneSnapshot()[1]; got.role != macSplitPaneContentList || got.contentList == nil || got.primary {
		t.Fatal("the content list pane should carry the model without a WebView")
	}

	withoutSidebar := NewMacSplitView()
	withoutSidebar.AddContentList(NewMacContentList())
	withoutSidebar.AddPrimaryContent()
	if err := validateMacSplitView(withoutSidebar); err != nil {
		t.Fatalf("content list before primary without a sidebar should validate: %v", err)
	}

	tests := []struct {
		name string
		make func() *MacSplitView
		want string
	}{
		{"before the sidebar", func() *MacSplitView {
			s := NewMacSplitView()
			s.AddContentList(NewMacContentList())
			s.AddSidebar(NewMacSidebar())
			s.AddPrimaryContent()
			return s
		}, "after the sidebar"},
		{"after the primary", func() *MacSplitView {
			s := NewMacSplitView()
			s.AddSidebar(NewMacSidebar())
			s.AddPrimaryContent()
			s.AddContentList(NewMacContentList())
			return s
		}, "before the primary"},
		{"two content lists", func() *MacSplitView {
			s := NewMacSplitView()
			s.AddContentList(NewMacContentList())
			s.AddContentList(NewMacContentList())
			s.AddPrimaryContent()
			return s
		}, "at most one content list"},
		{"content list without a primary", func() *MacSplitView {
			s := NewMacSplitView()
			s.AddSidebar(NewMacSidebar())
			s.AddContentList(NewMacContentList())
			return s
		}, "exactly one primary"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateMacSplitView(test.make())
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want mention of %q", err, test.want)
			}
		})
	}

	list := NewMacContentList()
	split := NewMacSplitView()
	if split.AddContentList(list) == nil || split.AddContentList(list) != nil {
		t.Fatal("one content list cannot be attached twice")
	}
	if split.AddContentList(nil) != nil || len(split.paneSnapshot()) != 1 {
		t.Fatal("nil content lists add no pane")
	}
}
