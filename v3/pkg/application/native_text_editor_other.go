//go:build !darwin || ios || server

package application

func macTextEditorApplyText(*MacTextEditor, string) bool  { return false }
func macTextEditorReadText(*MacTextEditor) (string, bool) { return "", false }
func macTextEditorApplyEditable(*MacTextEditor, bool)     {}
func macTextEditorFocus(*MacTextEditor)                   {}
