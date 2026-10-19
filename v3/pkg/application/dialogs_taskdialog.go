package application

import (
	"encoding/binary"
	"strings"
)

// Platform-independent helpers for the Windows message dialog implementation.
//
// MessageBox can only show fixed, localised button captions (Yes/No, ...), so
// question dialogs with custom buttons (#5794) are shown with TaskDialogIndirect.
// The logic that maps buttons to dialog IDs and back, and the byte layout of the
// (1-byte packed) TASKDIALOGCONFIG structure, lives here so it can be unit tested
// on any platform.

const (
	// taskDialogButtonIDBase is the first ID given to custom buttons. It must not
	// collide with the standard IDOK..IDCONTINUE (1-11) return values.
	taskDialogButtonIDBase int32 = 1000

	// idCancel is returned by TaskDialogIndirect when the dialog is dismissed
	// with Escape or the close button (TDF_ALLOW_DIALOG_CANCELLATION).
	idCancel int32 = 2

	// TASKDIALOGCONFIG flags
	tdfAllowDialogCancellation  uint32 = 0x0008
	tdfPositionRelativeToWindow uint32 = 0x1000
)

type taskDialogButton struct {
	ID    int32
	Label string
}

// taskDialogPlan describes the custom buttons of a task dialog and how to map
// the ID returned by TaskDialogIndirect back to a Button.
type taskDialogPlan struct {
	buttons   []taskDialogButton
	defaultID int32 // 0 when no default was requested
	allByID   map[int32]*Button
	cancel    *Button
}

func planTaskDialog(buttons []*Button) taskDialogPlan {
	p := taskDialogPlan{allByID: map[int32]*Button{}}
	for i, b := range buttons {
		if b == nil {
			continue
		}
		id := taskDialogButtonIDBase + int32(i)
		p.buttons = append(p.buttons, taskDialogButton{ID: id, Label: b.Label})
		p.allByID[id] = b
		if b.IsDefault && p.defaultID == 0 {
			p.defaultID = id
		}
		if b.IsCancel && p.cancel == nil {
			p.cancel = b
		}
	}
	return p
}

// buttonFor returns the button that was activated for the given result ID.
func (p taskDialogPlan) buttonFor(result int32) *Button {
	if b, ok := p.allByID[result]; ok {
		return b
	}
	if result == idCancel {
		return p.cancel
	}
	return nil
}

// messageBoxButtonFor maps a MessageBox result to a button of the dialog. It is
// used when TaskDialogIndirect is unavailable (no comctl32 v6 manifest). The
// button is found by caption first (case-insensitive, so "OK" matches "Ok"); as
// MessageBox cannot show custom captions, Yes/No results then fall back to the
// first/second non-cancel button. No and Cancel prefer the cancel button.
func messageBoxButtonFor(buttons []*Button, result int32) *Button {
	responses := []string{"", "Ok", "Cancel", "Abort", "Retry", "Ignore", "Yes", "No", "", "", "Try Again", "Continue"}
	if result > 0 && int(result) < len(responses) && responses[result] != "" {
		for _, b := range buttons {
			if b != nil && strings.EqualFold(strings.TrimSpace(b.Label), responses[result]) {
				return b
			}
		}
	}
	present, cancel := messageBoxFallbackButtons(buttons)
	switch result {
	case 6: // IDYES
		if len(present) > 0 {
			return present[0]
		}
	case 7: // IDNO
		if cancel != nil {
			return cancel
		}
		if len(present) > 1 {
			return present[1]
		}
	case 2: // IDCANCEL
		return cancel
	}
	return nil
}

// messageBoxFallbackButtons separates the cancel role from positional responses.
func messageBoxFallbackButtons(buttons []*Button) ([]*Button, *Button) {
	var present []*Button
	var cancel *Button
	for _, b := range buttons {
		if b == nil {
			continue
		}
		if b.IsCancel {
			if cancel == nil {
				cancel = b
			}
		} else {
			present = append(present, b)
		}
	}
	return present, cancel
}

// taskDialogConfig holds the values for a TASKDIALOGCONFIG. Pointer fields are
// passed as raw addresses; the caller keeps the pointed-to data alive.
type taskDialogConfig struct {
	hwndParent    uint64
	hInstance     uint64
	mainIcon      uint64
	callback      uint64
	flags         uint32
	title         uint64 // PCWSTR
	content       uint64 // PCWSTR
	buttonCount   uint32
	buttons       uint64 // *TASKDIALOG_BUTTON
	defaultButton int32
}

// encodeTaskDialogButtons lays out an array of TASKDIALOG_BUTTON {int; PCWSTR}.
// The structure is 1-byte packed in commctrl.h, so Go structs cannot be used.
func encodeTaskDialogButtons(ptrSize int, ids []int32, labels []uint64) []byte {
	size := 4 + ptrSize
	buf := make([]byte, size*len(ids))
	for i := range ids {
		off := i * size
		binary.LittleEndian.PutUint32(buf[off:], uint32(ids[i]))
		putPtr(buf[off+4:], ptrSize, labels[i])
	}
	return buf
}

// encodeTaskDialogConfig lays out a TASKDIALOGCONFIG, which commctrl.h declares
// inside pshpack1.h (1-byte packing).
func encodeTaskDialogConfig(ptrSize int, c taskDialogConfig) []byte {
	var buf []byte
	u32 := func(v uint32) { buf = binary.LittleEndian.AppendUint32(buf, v) }
	ptr := func(v uint64) {
		b := make([]byte, ptrSize)
		putPtr(b, ptrSize, v)
		buf = append(buf, b...)
	}
	u32(0) // cbSize, patched below
	ptr(c.hwndParent)
	ptr(c.hInstance)
	u32(c.flags)
	u32(0) // dwCommonButtons
	ptr(c.title)
	ptr(c.mainIcon) // pszMainIcon resource ID
	ptr(0)          // pszMainInstruction
	ptr(c.content)
	u32(c.buttonCount)
	ptr(c.buttons)
	u32(uint32(c.defaultButton))
	u32(0)          // cRadioButtons
	ptr(0)          // pRadioButtons
	u32(0)          // nDefaultRadioButton
	ptr(0)          // pszVerificationText
	ptr(0)          // pszExpandedInformation
	ptr(0)          // pszExpandedControlText
	ptr(0)          // pszCollapsedControlText
	ptr(0)          // hFooterIcon / pszFooterIcon
	ptr(0)          // pszFooter
	ptr(c.callback) // pfCallback
	ptr(0)          // lpCallbackData
	u32(0)          // cxWidth
	binary.LittleEndian.PutUint32(buf, uint32(len(buf)))
	return buf
}

func putPtr(b []byte, ptrSize int, v uint64) {
	if ptrSize == 8 {
		binary.LittleEndian.PutUint64(b, v)
	} else {
		binary.LittleEndian.PutUint32(b, uint32(v))
	}
}
