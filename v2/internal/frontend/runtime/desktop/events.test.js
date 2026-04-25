import { EventsOnMultiple, EventsNotify, eventListeners, EventsOn, EventsEmit, EventsOffAll, EventsOnce, EventsOff } from './events'
import { expect, describe, it, beforeAll, vi, afterEach, beforeEach } from 'vitest'
// Edit an assertion and save to see HMR in action

beforeAll(() => {
  window.WailsInvoke = vi.fn(() => {})
})

afterEach(() => {
  EventsOffAll();
  vi.resetAllMocks()
})

describe('EventsOnMultiple', () => {
  it('should stop after a specified number of times', () => {
    const cb = vi.fn()
    EventsOnMultiple('a', cb, 5)
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    expect(cb).toBeCalledTimes(5);
    expect(window.WailsInvoke).toBeCalledTimes(1);
    expect(window.WailsInvoke).toHaveBeenLastCalledWith('EXa');
  })

  it('should return a cancel fn', () => {
    const cb = vi.fn()
    const cancel = EventsOnMultiple('a', cb, 5)
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    cancel()
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    EventsNotify(JSON.stringify({name: 'a', data: {}}))
    expect(cb).toBeCalledTimes(2)
    expect(window.WailsInvoke).toBeCalledTimes(1);
    expect(window.WailsInvoke).toHaveBeenLastCalledWith('EXa');
  })
})

describe('EventsOn', () => {
  it('should create a listener with a count of -1', () => {
    EventsOn('a', () => {})
    expect(eventListeners['a'][0].maxCallbacks).toBe(-1)
  })

  it('should return a cancel fn', () => {
    const cancel = EventsOn('a', () => {})
    cancel();
    expect(window.WailsInvoke).toBeCalledTimes(1);
    expect(window.WailsInvoke).toHaveBeenLastCalledWith('EXa');
  })
})

describe('EventsOnce', () => {
  it('should create a listener with a count of 1', () => {
    EventsOnce('a', () => {})
    expect(eventListeners['a'][0].maxCallbacks).toBe(1)
  })

  it('should return a cancel fn', () => {
    const cancel = EventsOn('a', () => {})
    cancel();
    expect(window.WailsInvoke).toBeCalledTimes(1);
    expect(window.WailsInvoke).toHaveBeenLastCalledWith('EXa');
  })
})

describe('EventsNotify', () => {
  it('should inform a listener', () => {
    const cb = vi.fn()
    EventsOn('a', cb)
    EventsNotify(JSON.stringify({name: 'a', data: ["one", "two", "three"]}))
    expect(cb).toBeCalledTimes(1);
    expect(cb).toHaveBeenLastCalledWith("one", "two", "three");
    expect(window.WailsInvoke).toBeCalledTimes(0);
  })
})

describe('EventsEmit', () => {
  it('should emit an event', () => {
    EventsEmit('a', 'one', 'two', 'three')
    expect(window.WailsInvoke).toBeCalledTimes(1);
    const calledWith = window.WailsInvoke.calls[0][0];
    expect(calledWith.slice(0, 2)).toBe('EE')
    expect(JSON.parse(calledWith.slice(2))).toStrictEqual({data: ["one", "two", "three"], name: "a"})
  })
})

describe('EventsOff', () => {
  beforeEach(() => {
    EventsOn('a', () => {})
    EventsOn('a', () => {})
    EventsOn('a', () => {})
    EventsOn('b', () => {})
    EventsOn('c', () => {})
  })

  it('should cancel all event listeners for a single type', () => {
    EventsOff('a')
    expect(eventListeners['a']).toBeUndefined()
    expect(eventListeners['b']).not.toBeUndefined()
    expect(eventListeners['c']).not.toBeUndefined()
    expect(window.WailsInvoke).toBeCalledTimes(1);
    expect(window.WailsInvoke).toHaveBeenLastCalledWith('EXa');
  })

  it('should cancel all event listeners for multiple types', () => {
    EventsOff('a', 'b')
    expect(eventListeners['a']).toBeUndefined()
    expect(eventListeners['b']).toBeUndefined()
    expect(eventListeners['c']).not.toBeUndefined()
    expect(window.WailsInvoke).toBeCalledTimes(2);
    expect(window.WailsInvoke.calls).toStrictEqual([['EXa'], ['EXb']]);
  })
})

describe('notifyListeners - listenerOff during callback', () => {
  it('should keep listener removed when cancelled during notification (issue #4393)', () => {
    const cb1 = vi.fn()
    const cb2 = vi.fn()
    const cb3 = vi.fn()

    let cancelCb2 = null

    EventsOn('a', cb1)
    cancelCb2 = EventsOn('a', () => {
      cb2()
      cancelCb2()
    })
    EventsOn('a', cb3)

    EventsNotify(JSON.stringify({name: 'a', data: []}))

    expect(cb1).toBeCalledTimes(1)
    expect(cb2).toBeCalledTimes(1)
    expect(cb3).toBeCalledTimes(1)

    expect(eventListeners['a'].length).toBe(2)

    EventsNotify(JSON.stringify({name: 'a', data: []}))

    expect(cb1).toBeCalledTimes(2)
    expect(cb2).toBeCalledTimes(1)
    expect(cb3).toBeCalledTimes(2)
  })

  it('should handle all listeners being removed during notification', () => {
    const cb1 = vi.fn()
    const cb2 = vi.fn()

    let cancelCb1 = null
    let cancelCb2 = null

    cancelCb1 = EventsOn('a', () => {
      cb1()
      cancelCb1()
    })
    cancelCb2 = EventsOn('a', () => {
      cb2()
      cancelCb2()
    })

    EventsNotify(JSON.stringify({name: 'a', data: []}))

    expect(cb1).toBeCalledTimes(1)
    expect(cb2).toBeCalledTimes(1)

    expect(eventListeners['a']).toBeUndefined()

    expect(window.WailsInvoke).toBeCalledTimes(1);
    expect(window.WailsInvoke).toHaveBeenLastCalledWith('EXa');
  })
})

describe('EventsOffAll', () => {
  it('should cancel all event listeners', () => {
    EventsOn('a', () => {})
    EventsOn('a', () => {})
    EventsOn('a', () => {})
    EventsOn('b', () => {})
    EventsOn('c', () => {})
    EventsOffAll()
    expect(eventListeners).toStrictEqual({})
    expect(window.WailsInvoke).toBeCalledTimes(3);
    expect(window.WailsInvoke.calls).toStrictEqual([['EXa'], ['EXb'], ['EXc']]);
  })
})
