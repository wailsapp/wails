//go:build linux
// +build linux

package linux

import (
	"encoding/json"
	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v2/pkg/options"
	"os"
	"strings"
)

func SetupSingleInstance(uniqueID string, activateAppOnSubsequentLaunch bool, callback func(data options.SecondInstanceData)) {
	id := "wails_app_" + strings.ReplaceAll(uniqueID, "-", "_")

	dbusName := "org." + id + ".SingleInstance"
	dbusPath := "/org/" + id + "/SingleInstance"
	var err error

	conn, err := dbus.ConnectSessionBus()
	if err != nil {

	}

	f := func(message string) {
		var secondInstanceData options.SecondInstanceData

		err := json.Unmarshal([]byte(message), &secondInstanceData)

		if err != nil {
			go callback(secondInstanceData)
		}
	}

	conn.Export(f, dbus.ObjectPath(dbusPath), dbusName)

	reply, err := conn.RequestName(dbusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}

	// name already taken, so send args using dbus to existing instance and exit
	if reply != dbus.RequestNameReplyPrimaryOwner {
		data := options.SecondInstanceData{
			Args: os.Args[1:],
		}
		serialized, _ := json.Marshal(data)
		conn.Object(dbusName, dbus.ObjectPath(dbusPath)).Call(dbusName, 0, serialized)
		os.Exit(1)
	}
}
