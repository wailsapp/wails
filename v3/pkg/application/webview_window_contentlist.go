package application

import (
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

// MacContentListStyle selects the NSTableView presentation on macOS 11 and
// newer. Older releases ignore it.
type MacContentListStyle uint8

const (
	// MacContentListStyleAutomatic lets AppKit choose from the table's context.
	MacContentListStyleAutomatic MacContentListStyle = iota
	// MacContentListStyleInset uses rounded, inset selection like Mail's
	// message list.
	MacContentListStyleInset
	// MacContentListStyleSourceList uses the source-list appearance.
	MacContentListStyleSourceList
	// MacContentListStylePlain uses edge-to-edge rows with square selection.
	MacContentListStylePlain
	// MacContentListStyleFullWidth uses edge-to-edge selection with no inset.
	MacContentListStyleFullWidth
)

func validMacContentListStyle(style MacContentListStyle) bool {
	return style <= MacContentListStyleFullWidth
}

// MacContentListAlignment positions a column's text.
type MacContentListAlignment uint8

const (
	MacContentListAlignLeading MacContentListAlignment = iota
	MacContentListAlignCenter
	MacContentListAlignTrailing
)

// MacContentListColumn describes one NSTableColumn in table mode. Widths are
// points; zero leaves AppKit's defaults. Sortable columns show a sort
// indicator and report header clicks through MacContentList.OnSort.
type MacContentListColumn struct {
	Title     string
	Width     float64
	MinWidth  float64
	MaxWidth  float64
	Sortable  bool
	Alignment MacContentListAlignment
}

// MacContentList is a native NSTableView hosted by AppKit's content-list
// split-item role: the middle column of Finder, Mail, and document browsers.
//
// Without columns the list shows Mail-style rich rows: a title, a secondary
// subtitle line, an optional leading SF Symbol, a trailing detail such as a
// date, and a count badge. Calling SetColumns switches to a multi-column
// table whose rows are filled with SetCells. No WebView is created.
type MacContentList struct {
	lock sync.RWMutex

	columns []MacContentListColumn
	rows    []*MacContentListRow
	// selectedRows mirrors the native selection in row order.
	selectedRows            []*MacContentListRow
	allowsMultipleSelection bool
	sortable                bool
	// sortColumn is -1 while no column is sorted.
	sortColumn        int
	sortAscending     bool
	rowHeight         float64
	style             MacContentListStyle
	alternatingRows   bool
	emptyText         string
	headerVisible     bool
	headerVisibleSet  bool
	contextMenu       *Menu
	onContextMenu     func(*Context, *MacContentListRow) *Menu
	onSelectionChange func(*Context, []*MacContentListRow)
	onActivate        func(*Context, *MacContentListRow)
	onSort            func(*Context, int, bool)
	pane              *MacSplitPane
	dead              bool
}

// MacContentListRow is one selectable table row. Identifiers are generated
// internally; retain the returned handle to update, select, or remove it.
type MacContentListRow struct {
	lock sync.RWMutex

	internalID  uint64
	title       string
	subtitle    string
	detail      string
	symbol      string
	tooltip     string
	badge       int
	cells       []string
	disabled    bool
	hidden      bool
	contextMenu *Menu
	removed     bool
	list        *MacContentList
}

var macContentListRowID uint64

func nextMacContentListRowID() uint64 {
	return atomic.AddUint64(&macContentListRowID, 1)
}

// NewMacContentList creates an empty native content list in rich-row mode.
func NewMacContentList() *MacContentList {
	return &MacContentList{sortColumn: -1}
}

// SetColumns switches the list to table mode with the given columns. Rows
// then display their SetCells values, one per column. Passing no columns
// returns to rich-row mode.
func (l *MacContentList) SetColumns(columns ...MacContentListColumn) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.columns = append([]MacContentListColumn(nil), columns...)
	if l.sortColumn >= len(l.columns) && len(l.columns) > 0 {
		l.sortColumn = -1
	}
	l.lock.Unlock()
	l.reload()
	return l
}

// Columns returns the configured table columns; nil in rich-row mode.
func (l *MacContentList) Columns() []MacContentListColumn {
	if l == nil {
		return nil
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	return append([]MacContentListColumn(nil), l.columns...)
}

// AddRow appends a row. In table mode the title fills the first column until
// SetCells is called.
func (l *MacContentList) AddRow(title string) *MacContentListRow {
	if l == nil {
		return nil
	}
	l.lock.RLock()
	count := len(l.rows)
	l.lock.RUnlock()
	return l.InsertRow(count, title)
}

// InsertRow adds a row at index, clamped to the list's bounds.
func (l *MacContentList) InsertRow(index int, title string) *MacContentListRow {
	if l == nil {
		return nil
	}
	row := &MacContentListRow{
		internalID: nextMacContentListRowID(),
		title:      title,
		list:       l,
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return nil
	}
	if index < 0 {
		index = 0
	}
	if index > len(l.rows) {
		index = len(l.rows)
	}
	rows := make([]*MacContentListRow, 0, len(l.rows)+1)
	rows = append(rows, l.rows[:index]...)
	rows = append(rows, row)
	rows = append(rows, l.rows[index:]...)
	l.rows = rows
	l.lock.Unlock()
	l.reload()
	return row
}

// Rows returns every row in display order, hidden rows included.
func (l *MacContentList) Rows() []*MacContentListRow {
	if l == nil {
		return nil
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	return append([]*MacContentListRow(nil), l.rows...)
}

// RemoveAll removes every row and clears the selection. The handles are
// inert afterwards.
func (l *MacContentList) RemoveAll() *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	rows := l.rows
	l.rows = nil
	l.selectedRows = nil
	l.lock.Unlock()
	for _, row := range rows {
		row.retire()
	}
	l.reload()
	return l
}

// SetAllowsMultipleSelection lets the user select several rows with Command
// and Shift clicks.
func (l *MacContentList) SetAllowsMultipleSelection(allowed bool) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.allowsMultipleSelection = allowed
	if !allowed && len(l.selectedRows) > 1 {
		l.selectedRows = l.selectedRows[:1]
	}
	l.lock.Unlock()
	l.reload()
	return l
}

// SelectedRows returns the selected rows in display order.
func (l *MacContentList) SelectedRows() []*MacContentListRow {
	if l == nil {
		return nil
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	return append([]*MacContentListRow(nil), l.selectedRows...)
}

// SelectedRow returns the first selected row, or nil when nothing is
// selected.
func (l *MacContentList) SelectedRow() *MacContentListRow {
	if l == nil {
		return nil
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	if len(l.selectedRows) == 0 {
		return nil
	}
	return l.selectedRows[0]
}

// SetSelectedRow selects one row without invoking OnSelectionChange. Passing
// nil clears the selection.
func (l *MacContentList) SetSelectedRow(row *MacContentListRow) *MacContentList {
	if row == nil {
		return l.SetSelectedRows(nil)
	}
	return l.SetSelectedRows([]*MacContentListRow{row})
}

// SetSelectedRows replaces the selection without invoking OnSelectionChange.
// Rows from another list or removed rows are ignored; with single selection
// only the first row is kept.
func (l *MacContentList) SetSelectedRows(rows []*MacContentListRow) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	selected := make([]*MacContentListRow, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.list != l || row.isRemoved() {
			continue
		}
		selected = append(selected, row)
		if !l.allowsMultipleSelection {
			break
		}
	}
	l.selectedRows = orderMacContentListRows(l.rows, selected)
	ids := macContentListRowIDs(l.selectedRows)
	l.lock.Unlock()
	macContentListApplySelection(l, ids)
	return l
}

// OnSelectionChange sets the callback invoked whenever the user changes the
// selection, including clearing it. Passing nil clears the callback.
func (l *MacContentList) OnSelectionChange(callback func(*Context, []*MacContentListRow)) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	l.onSelectionChange = callback
	l.lock.Unlock()
	return l
}

// OnActivate sets the callback invoked when the user double-clicks a row or
// presses Return with a single row selected. Passing nil clears it.
func (l *MacContentList) OnActivate(callback func(*Context, *MacContentListRow)) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	l.onActivate = callback
	l.lock.Unlock()
	return l
}

// SetSortable enables header-click sorting for columns marked Sortable. A
// click reports through OnSort; without an OnSort callback the list sorts
// itself with SortBy.
func (l *MacContentList) SetSortable(sortable bool) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.sortable = sortable
	l.lock.Unlock()
	l.reload()
	return l
}

// OnSort sets the callback invoked when the user clicks a sortable column
// header. The list records the column and direction so the native indicator
// persists; reorder the rows with SortBy or SortRows. Passing nil restores
// the default SortBy behaviour.
func (l *MacContentList) OnSort(callback func(*Context, int, bool)) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	l.onSort = callback
	l.lock.Unlock()
	return l
}

// SortBy reorders the rows by one column's text, case-insensitively, and
// shows the native sort indicator on that column. In rich-row mode column 0
// is the title, 1 the subtitle, and 2 the detail text. A negative column
// clears the indicator without reordering.
func (l *MacContentList) SortBy(column int, ascending bool) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	if column < 0 {
		l.sortColumn = -1
		l.lock.Unlock()
		l.reload()
		return l
	}
	l.sortColumn = column
	l.sortAscending = ascending
	l.lock.Unlock()
	return l.SortRows(func(a, b *MacContentListRow) bool {
		left := strings.ToLower(a.sortText(column))
		right := strings.ToLower(b.sortText(column))
		if ascending {
			return left < right
		}
		return left > right
	})
}

// SortRows reorders the rows with a stable sort using less. It leaves the
// sort indicator unchanged; pair it with OnSort to sort by values that are
// not displayed, such as dates shown as relative text.
func (l *MacContentList) SortRows(less func(a, b *MacContentListRow) bool) *MacContentList {
	if l == nil || less == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	rows := append([]*MacContentListRow(nil), l.rows...)
	l.lock.Unlock()
	sort.SliceStable(rows, func(i, j int) bool { return less(rows[i], rows[j]) })
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.rows = rows
	l.selectedRows = orderMacContentListRows(rows, l.selectedRows)
	l.lock.Unlock()
	l.reload()
	return l
}

// SortOrder returns the sorted column and direction, or -1 when unsorted.
func (l *MacContentList) SortOrder() (column int, ascending bool) {
	if l == nil {
		return -1, false
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.sortColumn, l.sortAscending
}

// SetContextMenu sets the menu shown when the user right-clicks a row without
// its own menu or the empty area below the rows.
func (l *MacContentList) SetContextMenu(menu *Menu) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	l.contextMenu = menu
	l.lock.Unlock()
	return l
}

// OnContextMenu sets a callback that builds the context menu for a
// right-click. row is nil for the empty area. Returning nil falls back to the
// menu set with SetContextMenu; a row's own SetContextMenu takes precedence.
// The callback runs synchronously on the application thread while AppKit
// waits to show the menu, so keep it quick.
func (l *MacContentList) OnContextMenu(callback func(*Context, *MacContentListRow) *Menu) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	l.onContextMenu = callback
	l.lock.Unlock()
	return l
}

// SetRowHeight sets the row height in points. Zero restores the default:
// two-line rows in rich-row mode and AppKit's standard height in table mode.
func (l *MacContentList) SetRowHeight(points float64) *MacContentList {
	if l == nil {
		return l
	}
	if !isFiniteMacSplitNumber(points) || points < 0 {
		reportMacSplitError(l.ownerWindow(), "SetRowHeight: value %.1f must be a finite, non-negative number", points)
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.rowHeight = points
	l.lock.Unlock()
	l.reload()
	return l
}

// SetStyle selects the NSTableView style on macOS 11 and newer.
func (l *MacContentList) SetStyle(style MacContentListStyle) *MacContentList {
	if l == nil {
		return l
	}
	if !validMacContentListStyle(style) {
		reportMacSplitError(l.ownerWindow(), "SetStyle: unknown content list style %d", style)
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.style = style
	l.lock.Unlock()
	l.reload()
	return l
}

// SetAlternatingRowBackgrounds toggles AppKit's striped row backgrounds.
func (l *MacContentList) SetAlternatingRowBackgrounds(alternating bool) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.alternatingRows = alternating
	l.lock.Unlock()
	l.reload()
	return l
}

// SetEmptyText sets the placeholder shown centred in the list when it has no
// visible rows. An empty string shows nothing.
func (l *MacContentList) SetEmptyText(text string) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.emptyText = text
	l.lock.Unlock()
	l.reload()
	return l
}

// SetHeaderVisible shows or hides the column header. The header is shown by
// default in table mode and hidden in rich-row mode.
func (l *MacContentList) SetHeaderVisible(visible bool) *MacContentList {
	if l == nil {
		return l
	}
	l.lock.Lock()
	if l.dead {
		l.lock.Unlock()
		return l
	}
	l.headerVisible = visible
	l.headerVisibleSet = true
	l.lock.Unlock()
	l.reload()
	return l
}

func (l *MacContentList) ownerWindow() macSplitWindow {
	if l == nil {
		return nil
	}
	l.lock.RLock()
	pane := l.pane
	l.lock.RUnlock()
	if pane == nil || pane.split == nil {
		return nil
	}
	return pane.split.ownerWindow()
}

// Row setters. Text and badge changes update the native row in place;
// visibility changes reload the list.

// SetTitle updates the row's primary text. In table mode it also fills the
// first column when no cells are set.
func (r *MacContentListRow) SetTitle(title string) *MacContentListRow {
	return r.update(func() { r.title = title })
}

// SetSubtitle sets the secondary line shown beneath the title in rich-row
// mode.
func (r *MacContentListRow) SetSubtitle(subtitle string) *MacContentListRow {
	return r.update(func() { r.subtitle = subtitle })
}

// SetDetail sets the trailing secondary text shown beside the title in
// rich-row mode, for example a date.
func (r *MacContentListRow) SetDetail(detail string) *MacContentListRow {
	return r.update(func() { r.detail = detail })
}

// SetSymbol sets a leading SF Symbol on macOS 11 and newer in rich-row mode.
// Passing an empty string removes the image.
func (r *MacContentListRow) SetSymbol(symbol string) *MacContentListRow {
	return r.update(func() { r.symbol = symbol })
}

// SetBadge shows a trailing count in rich-row mode, like Mail's unread
// counts. Zero or a negative count removes the badge.
func (r *MacContentListRow) SetBadge(count int) *MacContentListRow {
	if count < 0 {
		count = 0
	}
	return r.update(func() { r.badge = count })
}

// SetCells sets the row's column values for table mode, in column order.
// Missing trailing values show as empty cells.
func (r *MacContentListRow) SetCells(values ...string) *MacContentListRow {
	copied := append([]string(nil), values...)
	return r.update(func() { r.cells = copied })
}

// SetEnabled controls whether the row can be selected.
func (r *MacContentListRow) SetEnabled(enabled bool) *MacContentListRow {
	return r.update(func() { r.disabled = !enabled })
}

// SetTooltip sets the native row tooltip.
func (r *MacContentListRow) SetTooltip(tooltip string) *MacContentListRow {
	return r.update(func() { r.tooltip = tooltip })
}

// SetHidden includes or removes the row from the native table. Hidden rows
// keep their handle and leave the selection.
func (r *MacContentListRow) SetHidden(hidden bool) *MacContentListRow {
	if r == nil || r.isDead() {
		return r
	}
	r.lock.Lock()
	changed := r.hidden != hidden
	r.hidden = hidden
	r.lock.Unlock()
	if !changed {
		return r
	}
	if hidden {
		r.list.pruneSelection()
	}
	r.list.reload()
	return r
}

// SetContextMenu sets the native menu shown when the row is right-clicked.
// Passing nil removes it; the list's OnContextMenu callback and fallback menu
// are consulted afterwards.
func (r *MacContentListRow) SetContextMenu(menu *Menu) *MacContentListRow {
	if r == nil || r.isDead() {
		return r
	}
	r.lock.Lock()
	r.contextMenu = menu
	r.lock.Unlock()
	return r
}

// Remove detaches the row from the list. The handle is inert afterwards.
func (r *MacContentListRow) Remove() {
	if r == nil || r.isDead() {
		return
	}
	list := r.list
	list.lock.Lock()
	for index, candidate := range list.rows {
		if candidate == r {
			list.rows = append(list.rows[:index:index], list.rows[index+1:]...)
			break
		}
	}
	list.lock.Unlock()
	r.retire()
	list.pruneSelection()
	list.reload()
}

// Title returns the row's primary text.
func (r *MacContentListRow) Title() string {
	if r == nil {
		return ""
	}
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.title
}

// Subtitle returns the row's secondary text.
func (r *MacContentListRow) Subtitle() string {
	if r == nil {
		return ""
	}
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.subtitle
}

// Detail returns the row's trailing detail text.
func (r *MacContentListRow) Detail() string {
	if r == nil {
		return ""
	}
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.detail
}

// Badge returns the row's badge count, or zero when none is shown.
func (r *MacContentListRow) Badge() int {
	if r == nil {
		return 0
	}
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.badge
}

// Cells returns the row's table-mode column values.
func (r *MacContentListRow) Cells() []string {
	if r == nil {
		return nil
	}
	r.lock.RLock()
	defer r.lock.RUnlock()
	return append([]string(nil), r.cells...)
}

// IsHidden reports whether the row is excluded from the native table.
func (r *MacContentListRow) IsHidden() bool {
	if r == nil {
		return false
	}
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.hidden
}

// update applies a field mutation and pushes the row to the native table.
func (r *MacContentListRow) update(mutate func()) *MacContentListRow {
	if r == nil || r.isDead() {
		return r
	}
	r.lock.Lock()
	mutate()
	r.lock.Unlock()
	macContentListApplyRow(r)
	return r
}

// sortText returns the text SortBy compares for column.
func (r *MacContentListRow) sortText(column int) string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if r.list != nil && r.list.hasColumns() {
		if column < len(r.cells) {
			return r.cells[column]
		}
		if column == 0 {
			return r.title
		}
		return ""
	}
	switch column {
	case 0:
		return r.title
	case 1:
		return r.subtitle
	case 2:
		return r.detail
	}
	return ""
}

func (l *MacContentList) hasColumns() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.columns) > 0
}

func (r *MacContentListRow) isRemoved() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.removed
}

func (r *MacContentListRow) isDead() bool {
	if r == nil || r.list == nil {
		return true
	}
	if r.isRemoved() {
		return true
	}
	r.list.lock.RLock()
	defer r.list.lock.RUnlock()
	return r.list.dead
}

// retire marks the row removed, drops its menu, and unregisters it.
func (r *MacContentListRow) retire() {
	if r == nil {
		return
	}
	r.lock.Lock()
	r.removed = true
	r.contextMenu = nil
	r.lock.Unlock()
	unregisterMacContentListRow(r.internalID)
}

// pruneSelection drops removed and hidden rows from the selection.
func (l *MacContentList) pruneSelection() {
	l.lock.Lock()
	defer l.lock.Unlock()
	kept := l.selectedRows[:0:0]
	for _, row := range l.selectedRows {
		row.lock.RLock()
		gone := row.removed || row.hidden
		row.lock.RUnlock()
		if !gone {
			kept = append(kept, row)
		}
	}
	l.selectedRows = kept
}

// orderMacContentListRows returns selected filtered to rows present in rows,
// in display order.
func orderMacContentListRows(rows, selected []*MacContentListRow) []*MacContentListRow {
	if len(selected) == 0 {
		return nil
	}
	wanted := make(map[*MacContentListRow]bool, len(selected))
	for _, row := range selected {
		wanted[row] = true
	}
	result := make([]*MacContentListRow, 0, len(selected))
	for _, row := range rows {
		if wanted[row] {
			result = append(result, row)
		}
	}
	return result
}

func macContentListRowIDs(rows []*MacContentListRow) []uint64 {
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.internalID)
	}
	return ids
}

type macContentListRowSnapshot struct {
	internalID uint64
	title      string
	subtitle   string
	detail     string
	symbol     string
	tooltip    string
	badge      int
	cells      []string
	disabled   bool
	hidden     bool
}

type macContentListSnapshot struct {
	columns                 []MacContentListColumn
	rows                    []macContentListRowSnapshot
	selectedIDs             []uint64
	allowsMultipleSelection bool
	sortable                bool
	sortColumn              int
	sortAscending           bool
	rowHeight               float64
	style                   MacContentListStyle
	alternatingRows         bool
	emptyText               string
	headerVisible           bool
}

func (l *MacContentList) snapshot() macContentListSnapshot {
	if l == nil {
		return macContentListSnapshot{sortColumn: -1}
	}
	l.lock.RLock()
	rows := append([]*MacContentListRow(nil), l.rows...)
	result := macContentListSnapshot{
		columns:                 append([]MacContentListColumn(nil), l.columns...),
		rows:                    make([]macContentListRowSnapshot, 0, len(rows)),
		selectedIDs:             macContentListRowIDs(l.selectedRows),
		allowsMultipleSelection: l.allowsMultipleSelection,
		sortable:                l.sortable,
		sortColumn:              l.sortColumn,
		sortAscending:           l.sortAscending,
		rowHeight:               l.rowHeight,
		style:                   l.style,
		alternatingRows:         l.alternatingRows,
		emptyText:               l.emptyText,
		headerVisible:           len(l.columns) > 0,
	}
	if l.headerVisibleSet {
		result.headerVisible = l.headerVisible
	}
	l.lock.RUnlock()
	for _, row := range rows {
		result.rows = append(result.rows, snapshotMacContentListRow(row))
	}
	return result
}

func snapshotMacContentListRow(row *MacContentListRow) macContentListRowSnapshot {
	row.lock.RLock()
	defer row.lock.RUnlock()
	return macContentListRowSnapshot{
		internalID: row.internalID,
		title:      row.title,
		subtitle:   row.subtitle,
		detail:     row.detail,
		symbol:     row.symbol,
		tooltip:    row.tooltip,
		badge:      row.badge,
		cells:      append([]string(nil), row.cells...),
		disabled:   row.disabled,
		hidden:     row.hidden,
	}
}

func (l *MacContentList) reload() {
	if l != nil {
		macContentListApplySnapshot(l)
	}
}

func (l *MacContentList) markDead() {
	if l == nil {
		return
	}
	for _, row := range l.Rows() {
		unregisterMacContentListRow(row.internalID)
		row.lock.Lock()
		row.contextMenu = nil
		row.lock.Unlock()
	}
	l.lock.Lock()
	l.dead = true
	l.pane = nil
	l.selectedRows = nil
	l.contextMenu = nil
	l.onContextMenu = nil
	l.onSelectionChange = nil
	l.onActivate = nil
	l.onSort = nil
	l.lock.Unlock()
}

func (l *MacContentList) registerRows() {
	for _, row := range l.Rows() {
		registerMacContentListRow(row)
	}
}

var macContentListRowRegistry = make(map[uint64]*MacContentListRow)
var macContentListRowRegistryLock sync.RWMutex

func registerMacContentListRow(row *MacContentListRow) {
	if row == nil {
		return
	}
	macContentListRowRegistryLock.Lock()
	macContentListRowRegistry[row.internalID] = row
	macContentListRowRegistryLock.Unlock()
}

func unregisterMacContentListRow(id uint64) {
	macContentListRowRegistryLock.Lock()
	delete(macContentListRowRegistry, id)
	macContentListRowRegistryLock.Unlock()
}

func macContentListRowByID(id uint64) *MacContentListRow {
	macContentListRowRegistryLock.RLock()
	defer macContentListRowRegistryLock.RUnlock()
	return macContentListRowRegistry[id]
}

// macContentListForPane resolves the content list hosted by a native split
// pane.
func macContentListForPane(paneID uint64) *MacContentList {
	pane := macSplitPaneByID(paneID)
	if pane == nil || pane.contentList == nil || pane.isDead() {
		return nil
	}
	return pane.contentList
}

// liveRow resolves a registered row that belongs to this list.
func (l *MacContentList) liveRow(id uint64) *MacContentListRow {
	row := macContentListRowByID(id)
	if row == nil || row.list != l || row.isDead() {
		return nil
	}
	return row
}

// macContentListSelectionEvent carries the native selection in row order.
type macContentListSelectionEvent struct {
	paneID uint64
	rowIDs []uint64
}

var macContentListSelectionEvents = make(chan macContentListSelectionEvent, 64)

// handleMacContentListSelectionChanged updates the list's selection and
// invokes OnSelectionChange with the whole set.
func handleMacContentListSelectionChanged(event macContentListSelectionEvent) {
	defer handlePanic()
	list := macContentListForPane(event.paneID)
	if list == nil {
		return
	}
	rows := make([]*MacContentListRow, 0, len(event.rowIDs))
	for _, id := range event.rowIDs {
		if row := list.liveRow(id); row != nil {
			rows = append(rows, row)
		}
	}
	list.lock.Lock()
	if list.dead {
		list.lock.Unlock()
		return
	}
	list.selectedRows = append([]*MacContentListRow(nil), rows...)
	callback := list.onSelectionChange
	list.lock.Unlock()
	if callback != nil {
		callback(newContext(), append([]*MacContentListRow(nil), rows...))
	}
}

type macContentListActivateEvent struct {
	paneID uint64
	rowID  uint64
}

var macContentListActivateEvents = make(chan macContentListActivateEvent, 64)

func handleMacContentListRowActivated(event macContentListActivateEvent) {
	defer handlePanic()
	list := macContentListForPane(event.paneID)
	if list == nil {
		return
	}
	row := list.liveRow(event.rowID)
	if row == nil {
		return
	}
	list.lock.RLock()
	callback := list.onActivate
	list.lock.RUnlock()
	if callback != nil {
		callback(newContext(), row)
	}
}

type macContentListSortEvent struct {
	paneID    uint64
	column    int
	ascending bool
}

var macContentListSortEvents = make(chan macContentListSortEvent, 64)

// handleMacContentListSortChanged records the header click and either
// notifies OnSort or applies the default SortBy ordering.
func handleMacContentListSortChanged(event macContentListSortEvent) {
	defer handlePanic()
	list := macContentListForPane(event.paneID)
	if list == nil || event.column < 0 {
		return
	}
	list.lock.Lock()
	if list.dead {
		list.lock.Unlock()
		return
	}
	list.sortColumn = event.column
	list.sortAscending = event.ascending
	callback := list.onSort
	list.lock.Unlock()
	if callback != nil {
		callback(newContext(), event.column, event.ascending)
		return
	}
	list.SortBy(event.column, event.ascending)
}

// resolveMacContentListContextMenu picks the menu for a right-click: the
// row's own menu, then the list's OnContextMenu callback, then the fallback
// menu. It runs synchronously on the application thread.
func resolveMacContentListContextMenu(paneID, rowID uint64) (menu *Menu) {
	defer handlePanic()
	list := macContentListForPane(paneID)
	if list == nil {
		return nil
	}
	var row *MacContentListRow
	if rowID != 0 {
		row = list.liveRow(rowID)
	}
	if row != nil {
		row.lock.RLock()
		menu = row.contextMenu
		row.lock.RUnlock()
		if menu != nil {
			return menu
		}
	}
	list.lock.RLock()
	dynamic := list.onContextMenu
	fallback := list.contextMenu
	list.lock.RUnlock()
	if dynamic != nil {
		if menu = dynamic(newContext(), row); menu != nil {
			return menu
		}
	}
	return fallback
}

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macContentListSelectionEvents
			go handleMacContentListSelectionChanged(event)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macContentListActivateEvents
			go handleMacContentListRowActivated(event)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macContentListSortEvents
			go handleMacContentListSortChanged(event)
		}
	})
}
