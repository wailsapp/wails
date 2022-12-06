//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

// ListItem represents an item in a ListView widget.
type ListItem interface {
	Text() []string  // Text returns the text of the multi-column item.
	ImageIndex() int // ImageIndex is used only if SetImageList is called on the listview
}

// ListItemChecker is used for checkbox support in ListView.
type ListItemChecker interface {
	Checked() bool
	SetChecked(checked bool)
}

// ListItemSetter is used in OnEndLabelEdit event.
type ListItemSetter interface {
	SetText(s string) // set first item in the array via LabelEdit event
}

// StringListItem is helper for basic string lists.
type StringListItem struct {
	ID    int
	Data  string
	Check bool
}

func (s StringListItem) Text() []string          { return []string{s.Data} }
func (s StringListItem) Checked() bool           { return s.Check }
func (s StringListItem) SetChecked(checked bool) { s.Check = checked }
func (s StringListItem) ImageIndex() int         { return 0 }

type ListView struct {
	ControlBase

	iml       *ImageList
	lastIndex int
	cols      int // count of columns

	item2Handle map[ListItem]uintptr
	handle2Item map[uintptr]ListItem

	onEndLabelEdit EventManager
	onDoubleClick  EventManager
	onClick        EventManager
	onKeyDown      EventManager
	onItemChanging EventManager
	onItemChanged  EventManager
	onCheckChanged EventManager
	onViewChange   EventManager
	onEndScroll    EventManager
}

func NewListView(parent Controller) *ListView {
	lv := new(ListView)

	lv.InitControl("SysListView32", parent /*w32.WS_EX_CLIENTEDGE*/, 0,
		w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.LVS_REPORT|w32.LVS_EDITLABELS|w32.LVS_SHOWSELALWAYS)

	lv.item2Handle = make(map[ListItem]uintptr)
	lv.handle2Item = make(map[uintptr]ListItem)

	RegMsgHandler(lv)

	lv.SetFont(DefaultFont)
	lv.SetSize(200, 400)

	if err := lv.SetTheme("Explorer"); err != nil {
		// theme error is ignored
	}
	return lv
}

// FIXME: Changes the state of an item in a list-view control. Refer LVM_SETITEMSTATE message.
func (lv *ListView) setItemState(i int, state, mask uint) {
	var item w32.LVITEM
	item.State, item.StateMask = uint32(state), uint32(mask)
	w32.SendMessage(lv.hwnd, w32.LVM_SETITEMSTATE, uintptr(i), uintptr(unsafe.Pointer(&item)))
}

func (lv *ListView) EnableSingleSelect(enable bool) {
	SetStyle(lv.hwnd, enable, w32.LVS_SINGLESEL)
}

func (lv *ListView) EnableSortHeader(enable bool) {
	SetStyle(lv.hwnd, enable, w32.LVS_NOSORTHEADER)
}

func (lv *ListView) EnableSortAscending(enable bool) {
	SetStyle(lv.hwnd, enable, w32.LVS_SORTASCENDING)
}

func (lv *ListView) EnableEditLabels(enable bool) {
	SetStyle(lv.hwnd, enable, w32.LVS_EDITLABELS)
}

func (lv *ListView) EnableFullRowSelect(enable bool) {
	if enable {
		w32.SendMessage(lv.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, 0, w32.LVS_EX_FULLROWSELECT)
	} else {
		w32.SendMessage(lv.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, w32.LVS_EX_FULLROWSELECT, 0)
	}
}

func (lv *ListView) EnableDoubleBuffer(enable bool) {
	if enable {
		w32.SendMessage(lv.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, 0, w32.LVS_EX_DOUBLEBUFFER)
	} else {
		w32.SendMessage(lv.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, w32.LVS_EX_DOUBLEBUFFER, 0)
	}
}

func (lv *ListView) EnableHotTrack(enable bool) {
	if enable {
		w32.SendMessage(lv.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, 0, w32.LVS_EX_TRACKSELECT)
	} else {
		w32.SendMessage(lv.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, w32.LVS_EX_TRACKSELECT, 0)
	}
}

func (lv *ListView) SetItemCount(count int) bool {
	return w32.SendMessage(lv.hwnd, w32.LVM_SETITEMCOUNT, uintptr(count), 0) != 0
}

func (lv *ListView) ItemCount() int {
	return int(w32.SendMessage(lv.hwnd, w32.LVM_GETITEMCOUNT, 0, 0))
}

func (lv *ListView) ItemAt(x, y int) ListItem {
	hti := w32.LVHITTESTINFO{Pt: w32.POINT{int32(x), int32(y)}}
	w32.SendMessage(lv.hwnd, w32.LVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))
	return lv.findItemByIndex(int(hti.IItem))
}

func (lv *ListView) Items() (list []ListItem) {
	for item := range lv.item2Handle {
		list = append(list, item)
	}
	return list
}

func (lv *ListView) AddColumn(caption string, width int) {
	var lc w32.LVCOLUMN
	lc.Mask = w32.LVCF_TEXT
	if width != 0 {
		lc.Mask = lc.Mask | w32.LVCF_WIDTH
		lc.Cx = int32(width)
	}
	lc.PszText = syscall.StringToUTF16Ptr(caption)
	lv.insertLvColumn(&lc, lv.cols)
	lv.cols++
}

// StretchLastColumn makes the last column take up all remaining horizontal
// space of the *ListView.
// The effect of this is not persistent.
func (lv *ListView) StretchLastColumn() error {
	if lv.cols == 0 {
		return nil
	}
	if w32.SendMessage(lv.hwnd, w32.LVM_SETCOLUMNWIDTH, uintptr(lv.cols-1), w32.LVSCW_AUTOSIZE_USEHEADER) == 0 {
		//panic("LVM_SETCOLUMNWIDTH failed")
	}
	return nil
}

// CheckBoxes returns if the *TableView has check boxes.
func (lv *ListView) CheckBoxes() bool {
	return w32.SendMessage(lv.hwnd, w32.LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)&w32.LVS_EX_CHECKBOXES > 0
}

// SetCheckBoxes sets if the *TableView has check boxes.
func (lv *ListView) SetCheckBoxes(value bool) {
	exStyle := w32.SendMessage(lv.hwnd, w32.LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	oldStyle := exStyle
	if value {
		exStyle |= w32.LVS_EX_CHECKBOXES
	} else {
		exStyle &^= w32.LVS_EX_CHECKBOXES
	}
	if exStyle != oldStyle {
		w32.SendMessage(lv.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, 0, exStyle)
	}

	mask := w32.SendMessage(lv.hwnd, w32.LVM_GETCALLBACKMASK, 0, 0)
	if value {
		mask |= w32.LVIS_STATEIMAGEMASK
	} else {
		mask &^= w32.LVIS_STATEIMAGEMASK
	}

	if w32.SendMessage(lv.hwnd, w32.LVM_SETCALLBACKMASK, mask, 0) == w32.FALSE {
		panic("SendMessage(LVM_SETCALLBACKMASK)")
	}
}

func (lv *ListView) applyImage(lc *w32.LVITEM, imIndex int) {
	if lv.iml != nil {
		lc.Mask |= w32.LVIF_IMAGE
		lc.IImage = int32(imIndex)
	}
}

func (lv *ListView) AddItem(item ListItem) {
	lv.InsertItem(item, lv.ItemCount())
}

func (lv *ListView) InsertItem(item ListItem, index int) {
	text := item.Text()
	li := &w32.LVITEM{
		Mask:    w32.LVIF_TEXT | w32.LVIF_PARAM,
		PszText: syscall.StringToUTF16Ptr(text[0]),
		IItem:   int32(index),
	}

	lv.lastIndex++
	ix := new(int)
	*ix = lv.lastIndex
	li.LParam = uintptr(*ix)
	lv.handle2Item[li.LParam] = item
	lv.item2Handle[item] = li.LParam

	lv.applyImage(li, item.ImageIndex())
	lv.insertLvItem(li)

	for i := 1; i < len(text); i++ {
		li.Mask = w32.LVIF_TEXT
		li.PszText = syscall.StringToUTF16Ptr(text[i])
		li.ISubItem = int32(i)
		lv.setLvItem(li)
	}
}

func (lv *ListView) UpdateItem(item ListItem) bool {
	lparam, ok := lv.item2Handle[item]
	if !ok {
		return false
	}

	index := lv.findIndexByItem(item)
	if index == -1 {
		return false
	}

	text := item.Text()
	li := &w32.LVITEM{
		Mask:    w32.LVIF_TEXT | w32.LVIF_PARAM,
		PszText: syscall.StringToUTF16Ptr(text[0]),
		LParam:  lparam,
		IItem:   int32(index),
	}

	lv.applyImage(li, item.ImageIndex())
	lv.setLvItem(li)

	for i := 1; i < len(text); i++ {
		li.Mask = w32.LVIF_TEXT
		li.PszText = syscall.StringToUTF16Ptr(text[i])
		li.ISubItem = int32(i)
		lv.setLvItem(li)
	}
	return true
}

func (lv *ListView) insertLvColumn(lvColumn *w32.LVCOLUMN, iCol int) {
	w32.SendMessage(lv.hwnd, w32.LVM_INSERTCOLUMN, uintptr(iCol), uintptr(unsafe.Pointer(lvColumn)))
}

func (lv *ListView) insertLvItem(lvItem *w32.LVITEM) {
	w32.SendMessage(lv.hwnd, w32.LVM_INSERTITEM, 0, uintptr(unsafe.Pointer(lvItem)))
}

func (lv *ListView) setLvItem(lvItem *w32.LVITEM) {
	w32.SendMessage(lv.hwnd, w32.LVM_SETITEM, 0, uintptr(unsafe.Pointer(lvItem)))
}

func (lv *ListView) DeleteAllItems() bool {
	if w32.SendMessage(lv.hwnd, w32.LVM_DELETEALLITEMS, 0, 0) == w32.TRUE {
		lv.item2Handle = make(map[ListItem]uintptr)
		lv.handle2Item = make(map[uintptr]ListItem)
		return true
	}
	return false
}

func (lv *ListView) DeleteItem(item ListItem) error {
	index := lv.findIndexByItem(item)
	if index == -1 {
		return errors.New("item not found")
	}

	if w32.SendMessage(lv.hwnd, w32.LVM_DELETEITEM, uintptr(index), 0) == 0 {
		return errors.New("SendMessage(TVM_DELETEITEM) failed")
	}

	h := lv.item2Handle[item]
	delete(lv.item2Handle, item)
	delete(lv.handle2Item, h)
	return nil
}

func (lv *ListView) findIndexByItem(item ListItem) int {
	lparam, ok := lv.item2Handle[item]
	if !ok {
		return -1
	}

	it := &w32.LVFINDINFO{
		Flags:  w32.LVFI_PARAM,
		LParam: lparam,
	}
	var i int = -1
	return int(w32.SendMessage(lv.hwnd, w32.LVM_FINDITEM, uintptr(i), uintptr(unsafe.Pointer(it))))
}

func (lv *ListView) findItemByIndex(i int) ListItem {
	it := &w32.LVITEM{
		Mask:  w32.LVIF_PARAM,
		IItem: int32(i),
	}

	if w32.SendMessage(lv.hwnd, w32.LVM_GETITEM, 0, uintptr(unsafe.Pointer(it))) == w32.TRUE {
		if item, ok := lv.handle2Item[it.LParam]; ok {
			return item
		}
	}
	return nil
}

func (lv *ListView) EnsureVisible(item ListItem) bool {
	if i := lv.findIndexByItem(item); i != -1 {
		return w32.SendMessage(lv.hwnd, w32.LVM_ENSUREVISIBLE, uintptr(i), 1) == 0
	}
	return false
}

func (lv *ListView) SelectedItem() ListItem {
	if items := lv.SelectedItems(); len(items) > 0 {
		return items[0]
	}
	return nil
}

func (lv *ListView) SetSelectedItem(item ListItem) bool {
	if i := lv.findIndexByItem(item); i > -1 {
		lv.SetSelectedIndex(i)
		return true
	}
	return false
}

// mask is used to set the LVITEM.Mask for ListView.GetItem which indicates which attributes you'd like to receive
// of LVITEM.
func (lv *ListView) SelectedItems() []ListItem {
	var items []ListItem

	var i int = -1
	for {
		if i = int(w32.SendMessage(lv.hwnd, w32.LVM_GETNEXTITEM, uintptr(i), uintptr(w32.LVNI_SELECTED))); i == -1 {
			break
		}

		if item := lv.findItemByIndex(i); item != nil {
			items = append(items, item)
		}
	}
	return items
}

func (lv *ListView) SelectedCount() uint {
	return uint(w32.SendMessage(lv.hwnd, w32.LVM_GETSELECTEDCOUNT, 0, 0))
}

// GetSelectedIndex first selected item index. Returns -1 if no item is selected.
func (lv *ListView) SelectedIndex() int {
	var i int = -1
	return int(w32.SendMessage(lv.hwnd, w32.LVM_GETNEXTITEM, uintptr(i), uintptr(w32.LVNI_SELECTED)))
}

// Set i to -1 to select all items.
func (lv *ListView) SetSelectedIndex(i int) {
	lv.setItemState(i, w32.LVIS_SELECTED, w32.LVIS_SELECTED)
}

func (lv *ListView) SetImageList(imageList *ImageList) {
	w32.SendMessage(lv.hwnd, w32.LVM_SETIMAGELIST, w32.LVSIL_SMALL, uintptr(imageList.Handle()))
	lv.iml = imageList
}

// Event publishers
func (lv *ListView) OnEndLabelEdit() *EventManager {
	return &lv.onEndLabelEdit
}

func (lv *ListView) OnDoubleClick() *EventManager {
	return &lv.onDoubleClick
}

func (lv *ListView) OnClick() *EventManager {
	return &lv.onClick
}

func (lv *ListView) OnKeyDown() *EventManager {
	return &lv.onKeyDown
}

func (lv *ListView) OnItemChanging() *EventManager {
	return &lv.onItemChanging
}

func (lv *ListView) OnItemChanged() *EventManager {
	return &lv.onItemChanged
}

func (lv *ListView) OnCheckChanged() *EventManager {
	return &lv.onCheckChanged
}

func (lv *ListView) OnViewChange() *EventManager {
	return &lv.onViewChange
}

func (lv *ListView) OnEndScroll() *EventManager {
	return &lv.onEndScroll
}

// Message processer
func (lv *ListView) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	/*case w32.WM_ERASEBKGND:
	lv.StretchLastColumn()
	println("case w32.WM_ERASEBKGND")
	return 1*/

	case w32.WM_NOTIFY:
		nm := (*w32.NMHDR)(unsafe.Pointer(lparam))
		code := int32(nm.Code)

		switch code {
		case w32.LVN_BEGINLABELEDITW:
			// println("Begin label edit")
		case w32.LVN_ENDLABELEDITW:
			nmdi := (*w32.NMLVDISPINFO)(unsafe.Pointer(lparam))
			if nmdi.Item.PszText != nil {
				fmt.Println(nmdi.Item.PszText, nmdi.Item)
				if item, ok := lv.handle2Item[nmdi.Item.LParam]; ok {
					lv.onEndLabelEdit.Fire(NewEvent(lv,
						&LabelEditEventData{Item: item,
							Text: w32.UTF16PtrToString(nmdi.Item.PszText)}))
				}
				return w32.TRUE
			}
		case w32.NM_DBLCLK:
			lv.onDoubleClick.Fire(NewEvent(lv, nil))

		case w32.NM_CLICK:
			ac := (*w32.NMITEMACTIVATE)(unsafe.Pointer(lparam))
			var hti w32.LVHITTESTINFO
			hti.Pt = w32.POINT{ac.PtAction.X, ac.PtAction.Y}
			w32.SendMessage(lv.hwnd, w32.LVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))

			if hti.Flags == w32.LVHT_ONITEMSTATEICON {
				if item := lv.findItemByIndex(int(hti.IItem)); item != nil {
					if item, ok := item.(ListItemChecker); ok {
						checked := !item.Checked()
						item.SetChecked(checked)
						lv.onCheckChanged.Fire(NewEvent(lv, item))

						if w32.SendMessage(lv.hwnd, w32.LVM_UPDATE, uintptr(hti.IItem), 0) == w32.FALSE {
							panic("SendMessage(LVM_UPDATE)")
						}
					}
				}
			}

			hti.Pt = w32.POINT{ac.PtAction.X, ac.PtAction.Y}
			w32.SendMessage(lv.hwnd, w32.LVM_SUBITEMHITTEST, 0, uintptr(unsafe.Pointer(&hti)))
			lv.onClick.Fire(NewEvent(lv, hti.ISubItem))

		case w32.LVN_KEYDOWN:
			nmkey := (*w32.NMLVKEYDOWN)(unsafe.Pointer(lparam))
			if nmkey.WVKey == w32.VK_SPACE && lv.CheckBoxes() {
				if item := lv.SelectedItem(); item != nil {
					if item, ok := item.(ListItemChecker); ok {
						checked := !item.Checked()
						item.SetChecked(checked)
						lv.onCheckChanged.Fire(NewEvent(lv, item))
					}

					index := lv.findIndexByItem(item)
					if w32.SendMessage(lv.hwnd, w32.LVM_UPDATE, uintptr(index), 0) == w32.FALSE {
						panic("SendMessage(LVM_UPDATE)")
					}
				}
			}
			lv.onKeyDown.Fire(NewEvent(lv, nmkey.WVKey))
			key := nmkey.WVKey
			w32.SendMessage(lv.Parent().Handle(), w32.WM_KEYDOWN, uintptr(key), 0)

		case w32.LVN_ITEMCHANGING:
			// This event also fires when listview has changed via code.
			nmlv := (*w32.NMLISTVIEW)(unsafe.Pointer(lparam))
			item := lv.findItemByIndex(int(nmlv.IItem))
			lv.onItemChanging.Fire(NewEvent(lv, item))

		case w32.LVN_ITEMCHANGED:
			// This event also fires when listview has changed via code.
			nmlv := (*w32.NMLISTVIEW)(unsafe.Pointer(lparam))
			item := lv.findItemByIndex(int(nmlv.IItem))
			lv.onItemChanged.Fire(NewEvent(lv, item))

		case w32.LVN_GETDISPINFO:
			nmdi := (*w32.NMLVDISPINFO)(unsafe.Pointer(lparam))
			if nmdi.Item.StateMask&w32.LVIS_STATEIMAGEMASK > 0 {
				if item, ok := lv.handle2Item[nmdi.Item.LParam]; ok {
					if item, ok := item.(ListItemChecker); ok {

						checked := item.Checked()
						if checked {
							nmdi.Item.State = 0x2000
						} else {
							nmdi.Item.State = 0x1000
						}
					}
				}
			}

			lv.onViewChange.Fire(NewEvent(lv, nil))

		case w32.LVN_ENDSCROLL:
			lv.onEndScroll.Fire(NewEvent(lv, nil))
		}
	}
	return w32.DefWindowProc(lv.hwnd, msg, wparam, lparam)
}
