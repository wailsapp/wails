//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

type EventHandler func(arg *Event)

type EventManager struct {
	handler EventHandler
}

func (evm *EventManager) Fire(arg *Event) {
	if evm.handler != nil {
		evm.handler(arg)
	}
}

func (evm *EventManager) Bind(handler EventHandler) {
	evm.handler = handler
}
