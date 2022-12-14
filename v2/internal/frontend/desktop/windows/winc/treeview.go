//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

// TreeItem represents an item in a TreeView widget.
type TreeItem interface {
	Text() string    // Text returns the text of the item.
	ImageIndex() int // ImageIndex is used only if SetImageList is called on the treeview
}

type treeViewItemInfo struct {
	handle       w32.HTREEITEM
	child2Handle map[TreeItem]w32.HTREEITEM
}

// StringTreeItem is helper for basic string lists.
type StringTreeItem struct {
	Data  string
	Image int
}

func (s StringTreeItem) Text() string    { return s.Data }
func (s StringTreeItem) ImageIndex() int { return s.Image }

type TreeView struct {
	ControlBase

	iml         *ImageList
	item2Info   map[TreeItem]*treeViewItemInfo
	handle2Item map[w32.HTREEITEM]TreeItem
	currItem    TreeItem

	onSelectedChange EventManager
	onExpand         EventManager
	onCollapse       EventManager
	onViewChange     EventManager
}

func NewTreeView(parent Controller) *TreeView {
	tv := new(TreeView)

	tv.InitControl("SysTreeView32", parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|
		w32.WS_BORDER|w32.TVS_HASBUTTONS|w32.TVS_LINESATROOT|w32.TVS_SHOWSELALWAYS|
		w32.TVS_TRACKSELECT /*|w32.WS_EX_CLIENTEDGE*/)

	tv.item2Info = make(map[TreeItem]*treeViewItemInfo)
	tv.handle2Item = make(map[w32.HTREEITEM]TreeItem)

	RegMsgHandler(tv)

	tv.SetFont(DefaultFont)
	tv.SetSize(200, 400)

	if err := tv.SetTheme("Explorer"); err != nil {
		// theme error is ignored
	}
	return tv
}

func (tv *TreeView) EnableDoubleBuffer(enable bool) {
	if enable {
		w32.SendMessage(tv.hwnd, w32.TVM_SETEXTENDEDSTYLE, 0, w32.TVS_EX_DOUBLEBUFFER)
	} else {
		w32.SendMessage(tv.hwnd, w32.TVM_SETEXTENDEDSTYLE, w32.TVS_EX_DOUBLEBUFFER, 0)
	}
}

// SelectedItem is current selected item after OnSelectedChange event.
func (tv *TreeView) SelectedItem() TreeItem {
	return tv.currItem
}

func (tv *TreeView) SetSelectedItem(item TreeItem) bool {
	var handle w32.HTREEITEM
	if item != nil {
		if info := tv.item2Info[item]; info == nil {
			return false // invalid item
		} else {
			handle = info.handle
		}
	}

	if w32.SendMessage(tv.hwnd, w32.TVM_SELECTITEM, w32.TVGN_CARET, uintptr(handle)) == 0 {
		return false // set selected failed
	}
	tv.currItem = item
	return true
}

func (tv *TreeView) ItemAt(x, y int) TreeItem {
	hti := w32.TVHITTESTINFO{Pt: w32.POINT{int32(x), int32(y)}}
	w32.SendMessage(tv.hwnd, w32.TVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))
	if item, ok := tv.handle2Item[hti.HItem]; ok {
		return item
	}
	return nil
}

func (tv *TreeView) Items() (list []TreeItem) {
	for item := range tv.item2Info {
		list = append(list, item)
	}
	return list
}

func (tv *TreeView) InsertItem(item, parent, insertAfter TreeItem) error {
	var tvins w32.TVINSERTSTRUCT
	tvi := &tvins.Item

	tvi.Mask = w32.TVIF_TEXT                                                     // w32.TVIF_CHILDREN | w32.TVIF_TEXT
	tvi.PszText = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item.Text()))) // w32.LPSTR_TEXTCALLBACK
	tvi.CChildren = 0                                                            // w32.I_CHILDRENCALLBACK

	if parent == nil {
		tvins.HParent = w32.TVI_ROOT
	} else {
		info := tv.item2Info[parent]
		if info == nil {
			return errors.New("winc: invalid parent")
		}
		tvins.HParent = info.handle
	}

	if insertAfter == nil {
		tvins.HInsertAfter = w32.TVI_LAST
	} else {
		info := tv.item2Info[insertAfter]
		if info == nil {
			return errors.New("winc: invalid prev item")
		}
		tvins.HInsertAfter = info.handle
	}

	tv.applyImage(tvi, item)

	hItem := w32.HTREEITEM(w32.SendMessage(tv.hwnd, w32.TVM_INSERTITEM, 0, uintptr(unsafe.Pointer(&tvins))))
	if hItem == 0 {
		return errors.New("winc: TVM_INSERTITEM failed")
	}
	tv.item2Info[item] = &treeViewItemInfo{hItem, make(map[TreeItem]w32.HTREEITEM)}
	tv.handle2Item[hItem] = item
	return nil
}

func (tv *TreeView) UpdateItem(item TreeItem) bool {
	it := tv.item2Info[item]
	if it == nil {
		return false
	}

	tvi := &w32.TVITEM{
		Mask:    w32.TVIF_TEXT,
		HItem:   it.handle,
		PszText: uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item.Text()))),
	}
	tv.applyImage(tvi, item)

	if w32.SendMessage(tv.hwnd, w32.TVM_SETITEM, 0, uintptr(unsafe.Pointer(tvi))) == 0 {
		return false
	}
	return true
}

func (tv *TreeView) DeleteItem(item TreeItem) bool {
	it := tv.item2Info[item]
	if it == nil {
		return false
	}

	if w32.SendMessage(tv.hwnd, w32.TVM_DELETEITEM, 0, uintptr(it.handle)) == 0 {
		return false
	}

	delete(tv.item2Info, item)
	delete(tv.handle2Item, it.handle)
	return true
}

func (tv *TreeView) DeleteAllItems() bool {
	if w32.SendMessage(tv.hwnd, w32.TVM_DELETEITEM, 0, 0) == 0 {
		return false
	}

	tv.item2Info = make(map[TreeItem]*treeViewItemInfo)
	tv.handle2Item = make(map[w32.HTREEITEM]TreeItem)
	return true
}

func (tv *TreeView) Expand(item TreeItem) bool {
	if w32.SendMessage(tv.hwnd, w32.TVM_EXPAND, w32.TVE_EXPAND, uintptr(tv.item2Info[item].handle)) == 0 {
		return false
	}
	return true
}

func (tv *TreeView) Collapse(item TreeItem) bool {
	if w32.SendMessage(tv.hwnd, w32.TVM_EXPAND, w32.TVE_COLLAPSE, uintptr(tv.item2Info[item].handle)) == 0 {
		return false
	}
	return true
}

func (tv *TreeView) EnsureVisible(item TreeItem) bool {
	if info := tv.item2Info[item]; info != nil {
		return w32.SendMessage(tv.hwnd, w32.TVM_ENSUREVISIBLE, 0, uintptr(info.handle)) != 0
	}
	return false
}

func (tv *TreeView) SetImageList(imageList *ImageList) {
	w32.SendMessage(tv.hwnd, w32.TVM_SETIMAGELIST, 0, uintptr(imageList.Handle()))
	tv.iml = imageList
}

func (tv *TreeView) applyImage(tc *w32.TVITEM, item TreeItem) {
	if tv.iml != nil {
		tc.Mask |= w32.TVIF_IMAGE | w32.TVIF_SELECTEDIMAGE
		tc.IImage = int32(item.ImageIndex())
		tc.ISelectedImage = int32(item.ImageIndex())
	}
}

func (tv *TreeView) OnSelectedChange() *EventManager {
	return &tv.onSelectedChange
}

func (tv *TreeView) OnExpand() *EventManager {
	return &tv.onExpand
}

func (tv *TreeView) OnCollapse() *EventManager {
	return &tv.onCollapse
}

func (tv *TreeView) OnViewChange() *EventManager {
	return &tv.onViewChange
}

// Message processer
func (tv *TreeView) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_NOTIFY:
		nm := (*w32.NMHDR)(unsafe.Pointer(lparam))

		switch nm.Code {
		case w32.TVN_ITEMEXPANDED:
			nmtv := (*w32.NMTREEVIEW)(unsafe.Pointer(lparam))

			switch nmtv.Action {
			case w32.TVE_COLLAPSE:
				tv.onCollapse.Fire(NewEvent(tv, nil))

			case w32.TVE_COLLAPSERESET:

			case w32.TVE_EXPAND:
				tv.onExpand.Fire(NewEvent(tv, nil))

			case w32.TVE_EXPANDPARTIAL:

			case w32.TVE_TOGGLE:
			}

		case w32.TVN_SELCHANGED:
			nmtv := (*w32.NMTREEVIEW)(unsafe.Pointer(lparam))
			tv.currItem = tv.handle2Item[nmtv.ItemNew.HItem]
			tv.onSelectedChange.Fire(NewEvent(tv, nil))

		case w32.TVN_GETDISPINFO:
			tv.onViewChange.Fire(NewEvent(tv, nil))
		}

	}
	return w32.DefWindowProc(tv.hwnd, msg, wparam, lparam)
}
