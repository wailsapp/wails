/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
// The following utilities have been factored out of ./events.ts
// for testing purposes.
export const eventListeners = new Map();
export class Listener {
    constructor(eventName, callback, maxCallbacks) {
        this.eventName = eventName;
        this.callback = callback;
        this.maxCallbacks = maxCallbacks || -1;
    }
    dispatch(data) {
        try {
            this.callback(data);
        }
        catch (err) {
            console.error(err);
        }
        if (this.maxCallbacks === -1)
            return false;
        this.maxCallbacks -= 1;
        return this.maxCallbacks === 0;
    }
}
export function listenerOff(listener) {
    let listeners = eventListeners.get(listener.eventName);
    if (!listeners) {
        return;
    }
    listeners = listeners.filter(l => l !== listener);
    if (listeners.length === 0) {
        eventListeners.delete(listener.eventName);
    }
    else {
        eventListeners.set(listener.eventName, listeners);
    }
}
