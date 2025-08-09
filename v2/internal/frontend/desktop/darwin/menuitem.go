//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "Application.h"
#import "WailsContext.h"

#include <stdlib.h>
*/
import "C"

import (
	"log"
	"math"
	"sync"
)

var (
	menuItemToID      = make(map[*MenuItem]uint)
	idToMenuItem      = make(map[uint]*MenuItem)
	menuItemLock      sync.Mutex
	menuItemIDCounter uint = 0
)

func createMenuItemID(item *MenuItem) uint {
	menuItemLock.Lock()
	defer menuItemLock.Unlock()
	counter := 0
	for {
		menuItemIDCounter++
		value := idToMenuItem[menuItemIDCounter]
		if value == nil {
			break
		}
		counter++
		if counter == math.MaxInt {
			log.Fatal("insane amounts of menuitems detected! Aborting before the collapse of the world!")
		}
	}
	idToMenuItem[menuItemIDCounter] = item
	menuItemToID[item] = menuItemIDCounter
	return menuItemIDCounter
}

func getMenuItemForID(id uint) *MenuItem {
	menuItemLock.Lock()
	defer menuItemLock.Unlock()
	return idToMenuItem[id]
}
