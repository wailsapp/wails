package dbus

//go:generate dbus-codegen-go -prefix org.kde -package notifier -output notifier/status_notifier_item.go StatusNotifierItem.xml
//go:generate dbus-codegen-go -prefix com.canonical -package menu -output menu/dbus_menu.go DbusMenu.xml
