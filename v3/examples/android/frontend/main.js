import {GreetService} from "./bindings/changeme";
import * as Runtime from "@wailsio/runtime";

const resultElement = document.getElementById('result');
const timeElement = document.getElementById('time');
const deviceInfoElement = document.getElementById('deviceInfo');
const Events = Runtime.Events;
const IOS = Runtime.IOS; // May be undefined in published package; we guard usages below.

window.doGreet = () => {
    let name = document.getElementById('name').value;
    if (!name) {
        name = 'anonymous';
    }
    GreetService.Greet(name).then((result) => {
        resultElement.innerText = result;
    }).catch((err) => {
        console.log(err);
    });
}

window.doHaptic = (style) => {
    if (!IOS || !IOS.Haptics?.Impact) {
        console.warn('IOS runtime not available in @wailsio/runtime. Skipping haptic call.');
        return;
    }
    IOS.Haptics.Impact(style).catch((err) => {
        console.error('Haptics error:', err);
    });
}

window.getDeviceInfo = async () => {
    if (!IOS || !IOS.Device?.Info) {
        deviceInfoElement.innerText = 'iOS runtime not available; cannot fetch device info.';
        return;
    }
    try {
        const info = await IOS.Device.Info();
        deviceInfoElement.innerText = JSON.stringify(info, null, 2);
    } catch (e) {
        deviceInfoElement.innerText = `Error: ${e?.message || e}`;
    }
}

// Generic caller for IOS.<Group>.<Method>(args)
window.iosJsSet = async (methodPath, args) => {
    if (!IOS) {
        console.warn('IOS runtime not available in @wailsio/runtime.');
        return;
    }
    try {
        const [group, method] = methodPath.split('.');
        const target = IOS?.[group];
        const fn = target?.[method];
        if (typeof fn !== 'function') {
            console.warn('IOS method not found:', methodPath);
            return;
        }
        await fn(args);
    } catch (e) {
        console.error('iosJsSet error for', methodPath, e);
    }
}

// Emit events for Go handlers
window.emitGo = (eventName, data) => {
    try {
        Events.Emit(eventName, data);
    } catch (e) {
        console.error('emitGo error:', e);
    }
}

// Toggle helpers for UI switches
window.setGoToggle = (eventName, enabled) => {
  emitGo(eventName, { enabled: !!enabled });
}

window.setJsToggle = (methodPath, enabled) => {
  iosJsSet(methodPath, { enabled: !!enabled });
}

Events.On('time', (payload) => {
    // payload may be a plain value or an object with a `data` field depending on emitter/runtime
    const value = (payload && typeof payload === 'object' && 'data' in payload) ? payload.data : payload;
    console.log('[frontend] time event:', payload, '->', value);
    timeElement.innerText = value;
});

// Simple pane switcher responding to native UITabBar
function showPaneByIndex(index) {
  const panes = [
    document.getElementById('screen-bindings'),
    document.getElementById('screen-go'),
    document.getElementById('screen-js'),
  ];
  panes.forEach((el, i) => {
    if (!el) return;
    if (i === index) el.classList.add('active');
    else el.classList.remove('active');
  });
}

// Listen for native tab selection events posted by the iOS layer
window.addEventListener('nativeTabSelected', (e) => {
  const idx = (e && e.detail && typeof e.detail.index === 'number') ? e.detail.index : 0;
  showPaneByIndex(idx);
});

// Ensure default pane is visible on load (index 0)
window.addEventListener('DOMContentLoaded', () => {
  showPaneByIndex(0);
});
