import {GreetService} from "./bindings/changeme";
import * as Runtime from "@wailsio/runtime";

const resultElement = document.getElementById('result');
const timeElement = document.getElementById('time');
const deviceInfoElement = document.getElementById('deviceInfo');
const Events = Runtime.Events;
const WML = Runtime.WML;

const Android = Runtime.Android; // May be undefined in published package; we guard usages below.

// Enable WML triggers (e.g. data-wml-openurl) so logo clicks call OpenURL.
WML?.Enable?.();

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
    if (!Android || !Android.Haptics?.Vibrate) {
        console.warn('Android runtime not available in @wailsio/runtime. Skipping haptic call.');
        return;
    }
    const duration = style === 'light' ? 50 : style === 'heavy' ? 200 : 100;
    Android.Haptics.Vibrate(duration).catch((err) => {
        console.error('Haptics error:', err);
    });
}

window.getDeviceInfo = async () => {
    try {
        if (!Android || !Android.Device?.Info) {
            deviceInfoElement.innerText = 'Android runtime not available; cannot fetch device info.';
            return;
        }
        const info = await Android.Device.Info();
        deviceInfoElement.innerText = JSON.stringify(info, null, 2);
    } catch (e) {
        deviceInfoElement.innerText = `Error: ${e?.message || e}`;
    }
}
// Generic caller for Android.<Group>.<Method>(args)
window.androidJsSet = async (methodPath, args) => {
    if (!Android) {
        console.warn('Android runtime not available in @wailsio/runtime.');
        return;
    }
    try {
        const [group, method] = methodPath.split('.');
        const target = Android?.[group];
        const fn = target?.[method];
        if (typeof fn !== 'function') {
            console.warn('Android method not found:', methodPath);
            return;
        }
        let payload = args;
        if (args && typeof args === 'object') {
            if ('enabled' in args) {
                payload = !!args.enabled;
            } else if ('ua' in args) {
                payload = args.ua;
            }
        }
        await fn(payload);
    } catch (e) {
        console.error('androidJsSet error for', methodPath, e);
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
    androidJsSet(methodPath, { enabled: !!enabled });
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

// Listen for native tab selection events posted by the native layer
window.addEventListener('nativeTabSelected', (e) => {
  const idx = (e && e.detail && typeof e.detail.index === 'number') ? e.detail.index : 0;
  showPaneByIndex(idx);
});

// Ensure default pane is visible on load (index 0)
window.addEventListener('DOMContentLoaded', () => {
  showPaneByIndex(0);
});
