package colour

import (
	"fmt"
	"strings"

	"github.com/wzshiming/ctc"
)

var ColourEnabled = true

func Col(col ctc.Color, text string) string {
	if !ColourEnabled {
		return text
	}
	return fmt.Sprintf("%s%s%s", col, text, ctc.Reset)
}

func Yellow(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightYellow, text)
}

func Red(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightRed, text)
}

func Blue(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightBlue, text)
}

func Green(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightGreen, text)
}

func Cyan(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightCyan, text)
}

func Magenta(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightMagenta, text)
}

func White(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightWhite, text)
}

func Black(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBrightBlack, text)
}

func DarkYellow(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundYellow, text)
}

func DarkRed(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundRed, text)
}

func DarkBlue(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBlue, text)
}

func DarkGreen(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundGreen, text)
}

func DarkCyan(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundCyan, text)
}

func DarkMagenta(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundMagenta, text)
}

func DarkWhite(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundWhite, text)
}

func DarkBlack(text string) string {
	if !ColourEnabled {
		return text
	}
	return Col(ctc.ForegroundBlack, text)
}

var rainbowCols = []func(string) string{Red, Yellow, Green, Cyan, Blue, Magenta}

func Rainbow(text string) string {
	if !ColourEnabled {
		return text
	}
	var builder strings.Builder

	for i := 0; i < len(text); i++ {
		fn := rainbowCols[i%len(rainbowCols)]
		builder.WriteString(fn(text[i : i+1]))
	}

	return builder.String()
}
