package menu

type CallbackData struct {
	MenuItem *MenuItem
	// ContextData string
}

type Callback func(*CallbackData)
