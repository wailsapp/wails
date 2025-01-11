package flags

import (
	"errors"
	"slices"
	"strings"
	"unicode/utf8"
)

type GenerateBindingsOptions struct {
	BuildFlagsString  string `name:"f" description:"A list of additional space-separated Go build flags. Flags (or parts of them) can be wrapped in single or double quotes to include spaces"`
	OutputDirectory   string `name:"d" description:"The output directory" default:"frontend/bindings"`
	ModelsFilename    string `name:"models" description:"File name for exported JS/TS models (excluding the extension)" default:"models"`
	InternalFilename  string `name:"internal" description:"File name for unexported JS/TS models (excluding the extension)" default:"internal"`
	IndexFilename     string `name:"index" description:"File name for JS/TS package indexes (excluding the extension)" default:"index"`
	TS                bool   `name:"ts" description:"Generate Typescript bindings"`
	UseInterfaces     bool   `name:"i" description:"Generate Typescript interfaces instead of classes"`
	UseBundledRuntime bool   `name:"b" description:"Use the bundled runtime instead of importing the npm package"`
	UseNames          bool   `name:"names" description:"Use names instead of IDs for the binding calls"`
	NoIndex           bool   `name:"noindex" description:"Do not generate JS/TS index files"`
	DryRun            bool   `name:"dry" description:"Do not write output files"`
	Silent            bool   `name:"silent" description:"Silent mode"`
	Verbose           bool   `name:"v" description:"Enable debug output"`
	Clean             bool   `name:"clean" description:"Clean output directory before generation" default:"true"`
}

var ErrUnmatchedQuote = errors.New("build flags contain an unmatched quote")

func isWhitespace(r rune) bool {
	// We use Go's definition of whitespace instead of the Unicode ones
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func isNonWhitespace(r rune) bool {
	return !isWhitespace(r)
}

func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}

func isQuoteOrWhitespace(r rune) bool {
	return isQuote(r) || isWhitespace(r)
}

func (options *GenerateBindingsOptions) BuildFlags() (flags []string, err error) {
	str := options.BuildFlagsString

	// temporary buffer for flag assembly
	flag := make([]byte, 0, 32)

	for start := strings.IndexFunc(str, isNonWhitespace); start >= 0; start = strings.IndexFunc(str, isNonWhitespace) {
		// each iteration starts at the beginning of a flag
		// skip initial whitespace
		str = str[start:]

		// iterate over all quoted and unquoted parts of the flag and join them
		for {
			breakpoint := strings.IndexFunc(str, isQuoteOrWhitespace)
			if breakpoint < 0 {
				breakpoint = len(str)
			}

			// append everything up to the breakpoint
			flag = append(flag, str[:breakpoint]...)
			str = str[breakpoint:]

			quote, quoteSize := utf8.DecodeRuneInString(str)
			if !isQuote(quote) {
				// if the breakpoint is not a quote, we reached the end of the flag
				break
			}

			// otherwise, look for the closing quote
			str = str[quoteSize:]
			closingQuote := strings.IndexRune(str, quote)

			// closing quote not found, append everything to the last flag and raise an error
			if closingQuote < 0 {
				flag = append(flag, str...)
				str = ""
				err = ErrUnmatchedQuote
				break
			}

			// closing quote found, append quoted content to the flag and restart after the quote
			flag = append(flag, str[:closingQuote]...)
			str = str[closingQuote+quoteSize:]
		}

		// append a clone of the flag to the result, then reuse buffer space
		flags = append(flags, string(slices.Clone(flag)))
		flag = flag[:0]
	}

	return
}
