package operatingsystem

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func platformInfo() (*OS, error) {
	// Default value
	var result OS
	result.ID = "Unknown"
	result.Name = "Windows"
	result.Version = "Unknown"

	// Credit: https://stackoverflow.com/a/33288328
	// Ignore errors as it isn't a showstopper
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)

	defer key.Close()

	fmt.Printf("%+v\n", key)

	// Ignore errors as it isn't a showstopper
	productName, _, _ := key.GetStringValue("ProductName")
	fmt.Println(productName)

	return nil, nil
}
