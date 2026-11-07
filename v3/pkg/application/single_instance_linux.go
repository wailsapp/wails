//go:build linux && !android && !server

package application

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/godbus/dbus/v5"
)

var setup sync.Once

// dbusHandler receives the second instance's call. The session bus is per user, but any process of that user can
// call the exported method, so a caller that is not this executable is dropped: its arguments would otherwise reach
// OnSecondInstanceLaunch as if a real second launch had sent them.
type dbusHandler struct {
	conn    *dbus.Conn
	deliver func(string)
}

func (h dbusHandler) SendMessage(sender dbus.Sender, message string) *dbus.Error {
	var pid uint32
	call := h.conn.BusObject().Call("org.freedesktop.DBus.GetConnectionUnixProcessID", 0, string(sender))
	if err := call.Store(&pid); err != nil {
		return dbus.MakeFailedError(errors.New("cannot identify the caller"))
	}
	if same, known := sameExecutable("/proc", int(pid)); known && !same {
		return dbus.MakeFailedError(errors.New("caller is not this application"))
	}
	h.deliver(message)
	return nil
}

// sameExecutable compares the executable of process pid with this process's. known is false when pid's executable
// cannot be read (a sandbox's bus proxy reports a pid from another namespace), which leaves the caller unverified
// rather than refused.
func sameExecutable(procRoot string, pid int) (same, known bool) {
	theirs, err := os.Stat(filepath.Join(procRoot, strconv.Itoa(pid), "exe"))
	if err != nil {
		return false, false
	}
	ours, err := os.Stat(filepath.Join(procRoot, "self", "exe"))
	if err != nil {
		return false, false
	}
	return os.SameFile(ours, theirs), true
}

type linuxLock struct {
	file          *os.File
	uniqueID      string
	dbusPath      string
	dbusName      string
	dbusInterface string
	manager       *singleInstanceManager
}

// singleInstanceNames derives the three D-Bus identifiers the lock needs from
// UniqueID. They cannot all be the same string, because D-Bus spells them
// differently: bus names may contain hyphens, interface names may not, and
// object paths separate elements with "/" and allow neither hyphens nor dots.
//
// The bus name keeps UniqueID verbatim. That matters for sandboxed builds: a
// flatpak may only own names prefixed with its own app id, so an app that
// follows the documented convention for UniqueID ("unique per application, e.g.
// com.myapp.myapplication") claims a name it is allowed to own and needs no
// extra portal permission.
func singleInstanceNames(uniqueID string) (busName, interfaceName, objectPath string, err error) {
	for _, element := range strings.Split(uniqueID, ".") {
		if element == "" {
			return "", "", "", fmt.Errorf("UniqueID %q has an empty element; it must be a dot-separated name such as com.myapp.myapplication", uniqueID)
		}
		if element[0] >= '0' && element[0] <= '9' {
			return "", "", "", fmt.Errorf("UniqueID %q has an element starting with a digit (%q), which D-Bus does not allow", uniqueID, element)
		}
		for _, r := range element {
			if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
				continue
			}
			return "", "", "", fmt.Errorf("UniqueID %q contains %q, which D-Bus does not allow in a name", uniqueID, r)
		}
	}

	// Hyphens are legal in a bus name but not in an interface name or an object
	// path, so those two are built from a hyphen-free form.
	unhyphenated := strings.ReplaceAll(uniqueID, "-", "_")

	busName = uniqueID + ".SingleInstance"
	interfaceName = unhyphenated + ".SingleInstance"
	objectPath = "/" + strings.ReplaceAll(unhyphenated, ".", "/") + "/SingleInstance"
	return busName, interfaceName, objectPath, nil
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

	busName, interfaceName, objectPath, err := singleInstanceNames(uniqueID)
	if err != nil {
		return err
	}

	l.uniqueID = uniqueID
	l.dbusName = busName
	l.dbusInterface = interfaceName
	l.dbusPath = objectPath

	conn, err := dbus.ConnectSessionBus()
	// if we will reach any error during establishing connection or sending message we will just continue.
	// It should not be the case that such thing will happen actually, but just in case.
	if err != nil {
		return err
	}

	setup.Do(func() {
		f := dbusHandler{conn: conn, deliver: func(message string) { secondInstanceBuffer <- message }}

		err = conn.Export(f, dbus.ObjectPath(l.dbusPath), l.dbusInterface)
	})
	if err != nil {
		return err
	}

	reply, err := conn.RequestName(l.dbusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		// A sandbox that refuses the name fails here. Say so, rather than
		// letting the caller read it as a second instance and exit silently.
		return fmt.Errorf("could not claim the single instance name %q: %w", l.dbusName, err)
	}

	switch reply {
	case dbus.RequestNameReplyPrimaryOwner, dbus.RequestNameReplyAlreadyOwner:
		return nil
	case dbus.RequestNameReplyExists:
		// Someone else holds the name, so this is a second instance. The caller
		// hands off to the first one and exits.
		return alreadyRunningError
	default:
		return fmt.Errorf("unexpected reply %d when claiming the single instance name %q", reply, l.dbusName)
	}
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

	err = conn.Object(l.dbusName, dbus.ObjectPath(l.dbusPath)).Call(l.dbusInterface+".SendMessage", 0, data).Store()
	if err != nil {
		return err
	}
	os.Exit(l.manager.options.ExitCode)
	return nil
}
