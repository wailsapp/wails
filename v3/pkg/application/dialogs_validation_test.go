package application

import (
	"testing"

	"github.com/matryer/is"
)

func countDefaultCancel(buttons []*Button) (defaults int, cancels int) {
	for _, b := range buttons {
		if b == nil {
			continue
		}
		if b.IsDefault {
			defaults++
		}
		if b.IsCancel {
			cancels++
		}
	}
	return defaults, cancels
}

func TestMessageDialogNormalizeButtonList_MultipleDefaultButtons_LastWins(t *testing.T) {
	i := is.New(t)

	d := newMessageDialog(QuestionDialogType).Buttons(
		Button{Label: "A", IsDefault: true},
		Button{Label: "B", IsDefault: true},
	)

	defaults, cancels := countDefaultCancel(d.ButtonList)
	i.Equal(defaults, 1)
	i.Equal(cancels, 0)
	i.Equal(d.ButtonList[1].IsDefault, true)
}

func TestMessageDialogNormalizeButtonList_MultipleCancelButtons_LastWins(t *testing.T) {
	i := is.New(t)

	d := newMessageDialog(QuestionDialogType).Buttons(
		Button{Label: "A", IsCancel: true},
		Button{Label: "B", IsCancel: true},
	)

	defaults, cancels := countDefaultCancel(d.ButtonList)
	i.Equal(defaults, 0)
	i.Equal(cancels, 1)
	i.Equal(d.ButtonList[1].IsCancel, true)
}

func TestMessageDialogAddDefaultButtonClearsPrevious(t *testing.T) {
	i := is.New(t)

	d := newMessageDialog(QuestionDialogType).
		AddDefaultButton("A").
		AddDefaultButton("B")

	defaults, cancels := countDefaultCancel(d.ButtonList)
	i.Equal(defaults, 1)
	i.Equal(cancels, 0)
}

func TestMessageDialogAddCancelButtonClearsPrevious(t *testing.T) {
	i := is.New(t)

	d := newMessageDialog(QuestionDialogType).
		AddCancelButton("A").
		AddCancelButton("B")

	defaults, cancels := countDefaultCancel(d.ButtonList)
	i.Equal(defaults, 0)
	i.Equal(cancels, 1)
}

func TestMessageDialogDefaultAndCancelOK(t *testing.T) {
	i := is.New(t)

	d := newMessageDialog(QuestionDialogType).Buttons(
		Button{Label: "Run"},
		Button{Label: "Cancel"},
	).Default("Run").Cancel("Cancel")

	defaults, cancels := countDefaultCancel(d.ButtonList)
	i.Equal(defaults, 1)
	i.Equal(cancels, 1)
}
