package reportcopy

// Write copies text using the system clipboard utility.
func Write(text string) error { return writeCommand(text, "pbcopy") }
