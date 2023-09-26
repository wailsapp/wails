//go:build linux
// +build linux

package linux

import (
	"encoding/json"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v2/pkg/options"
	"os"
	"strings"
)

type dbusHandler func(string)

func (f dbusHandler) SendMessage(message string) *dbus.Error {
	f(message)
	return nil
}

func SetupSingleInstance(uniqueID string, activateAppOnSubsequentLaunch bool, callback func(data options.SecondInstanceData)) {
	id := "wails_app_" + strings.ReplaceAll(uniqueID, "-", "_")

	dbusName := "org." + id + ".SingleInstance"
	dbusPath := "/org/" + id + "/SingleInstance"
	var err error

	conn, err := dbus.ConnectSessionBus()
	println("dbus connected")
	if err != nil {

	}

	f := dbusHandler(func(message string) {
		println("got message", message)
		var secondInstanceData options.SecondInstanceData

		err := json.Unmarshal([]byte(message), &secondInstanceData)

		if err != nil {
			go callback(secondInstanceData)
		}
	})

	println("try to export callback")
	conn.Export(f, dbus.ObjectPath(dbusPath), dbusName)

	println("try to request name")
	reply, err := conn.RequestName(dbusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		println("error on requesting name")
		panic(err)
	}

	println("got reply", reply)

	// name already taken, so send args using dbus to existing instance and exit
	if reply == dbus.RequestNameReplyExists {
		println("name already taken, sending request to existing instance")
		data := options.SecondInstanceData{
			Args: os.Args[1:],
		}
		serialized, _ := json.Marshal(data)

		err = conn.Object(dbusName, dbus.ObjectPath(dbusPath)).Call(dbusName+".SendMessage", 0, string(serialized)).Store()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to call Foo function (is the server example running?):", err)
		}
		println("sent successfully", string(serialized))
		os.Exit(1)
	}
}
