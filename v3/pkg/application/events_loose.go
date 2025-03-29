//go:build production || !strictevents

package application

func eventRegistered(name string) {}

func warnAboutUnregisteredEvent(name string) {}
