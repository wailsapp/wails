//go:build linux

package application

import (
	"errors"
	"os"
	"strings"
	"sync"
	"syscall"

	"github.com/godbus/dbus/v5"
)

type dbusHandler func(string)

var setup sync.Once

func (f dbusHandler) SendMessage(message string) *dbus.Error {
	f(message)
	return nil
}

type linuxLock struct {
	file     *os.File
	uniqueID string
	dbusPath string
	dbusName string
	manager  *singleInstanceManager
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &linuxLock{
		manager: manager,
	}, nil
}

func (l *linuxLock) acquire(uniqueID string) error {
	if uniqueID == "" {
		return errors.New("UniqueID is required for single instance lock")
	}

	id := "wails_app_" + strings.ReplaceAll(strings.ReplaceAll(uniqueID, "-", "_"), ".", "_")

	l.dbusName = "org." + id + ".SingleInstance"
	l.dbusPath = "/org/" + id + "/SingleInstance"

	conn, err := dbus.ConnectSessionBus()
	// if we will reach any error during establishing connection or sending message we will just continue.
	// It should not be the case that such thing will happen actually, but just in case.
	if err != nil {
		return err
	}

	setup.Do(func() {
		f := dbusHandler(func(message string) {
			secondInstanceBuffer <- message
		})

		err = conn.Export(f, dbus.ObjectPath(l.dbusPath), l.dbusName)
	})
	if err != nil {
		return err
	}

	reply, err := conn.RequestName(l.dbusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	// if name already taken, try to send args to existing instance, if no success just launch new instance
	if reply == dbus.RequestNameReplyExists {
		return alreadyRunningError
	}
	return nil
}

func (l *linuxLock) release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.file.Name())
		l.file = nil
	}
}

func (l *linuxLock) notify(data string) error {
	conn, err := dbus.ConnectSessionBus()
	// if we will reach any error during establishing connection or sending message we will just continue.
	// It should not be the case that such thing will happen actually, but just in case.
	if err != nil {
		return err
	}

	err = conn.Object(l.dbusName, dbus.ObjectPath(l.dbusPath)).Call(l.dbusName+".SendMessage", 0, data).Store()
	if err != nil {
		return err
	}
	os.Exit(l.manager.options.ExitCode)
	return nil
}
