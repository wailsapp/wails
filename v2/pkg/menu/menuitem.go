package menu

type MenuItem struct {
	Id          string `json:"Id,omitempty"`
	Label       string
	Role        Role   `json:"Role,omitempty"`
	Accelerator string `json:"Accelerator,omitempty"`
	Type        Type
	Disabled    bool
	Hidden      bool
	Checked     bool
	SubMenu     []*MenuItem `json:"SubMenu,omitempty"`
}

func Text(label string, id string) *MenuItem {
	return &MenuItem{
		Id:    id,
		Label: label,
		Type:  TextType,
	}
}

// Separator provides a menu separator
func Separator() *MenuItem {
	return &MenuItem{
		Type: SeparatorType,
	}
}
