//go:build android

package application

import "unsafe"

// Android doesn't have traditional menu items like desktop platforms
// These are placeholder implementations

func (m *MenuItem) handleStyleChange() {}

func (m *MenuItem) handleLabelChange() {}

func (m *MenuItem) handleCheckedChange() {}

func (m *MenuItem) handleEnabledChange() {}

func (m *MenuItem) handleTooltipChange() {}

func (m *MenuItem) handleSubmenuChange() {}

func (m *MenuItem) nativeMenuItem() unsafe.Pointer {
	return nil
}
