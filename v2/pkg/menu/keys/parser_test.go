package keys

import (
	"testing"

	"github.com/matryer/is"
)

func TestParse(t *testing.T) {

	i := is.New(t)

	type args struct {
		Input    string
		Expected *Accelerator
	}

	gooddata := []args{
		{"CmdOrCtrl+A", CmdOrCtrl("A")},
		{"SHIFT+.", Shift(".")},
		{"CTRL+plus", Control("+")},
		{"CTRL+SHIFT+escApe", Combo("escape", ControlKey, ShiftKey)},
		{";", Key(";")},
		{"OptionOrAlt+Page Down", OptionOrAlt("Page Down")},
	}
	for _, tt := range gooddata {
		result, err := Parse(tt.Input)
		i.NoErr(err)
		i.Equal(result, tt.Expected)
	}
	baddata := []string{"CmdOrCrl+A", "SHIT+.", "CTL+plus", "CTRL+SHIF+esApe", "escap", "Sper+Tab", "OptionOrAlt"}
	for _, d := range baddata {
		result, err := Parse(d)
		i.True(err != nil)
		i.Equal(result, nil)
	}
}
