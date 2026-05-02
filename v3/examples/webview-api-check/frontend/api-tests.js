// Comprehensive Web API Tests
// Each test returns { supported: true|false|'partial', note?: string }

const API_TESTS = {
    "Storage APIs": {
        "localStorage": () => ({
            supported: typeof localStorage !== 'undefined',
            note: typeof localStorage !== 'undefined' ? `${localStorage.length} items` : undefined
        }),
        "sessionStorage": () => ({
            supported: typeof sessionStorage !== 'undefined'
        }),
        "IndexedDB": () => ({
            supported: typeof indexedDB !== 'undefined'
        }),
        "Cache API": () => ({
            supported: 'caches' in window
        }),
        "CookieStore API": () => ({
            supported: 'cookieStore' in window,
            note: !('cookieStore' in window) ? 'Chromium only' : undefined
        }),
        "Storage API": () => ({
            supported: navigator.storage !== undefined
        }),
        "Storage Access API": () => ({
            supported: 'hasStorageAccess' in document
        }),
        "File System Access": () => ({
            supported: 'showOpenFilePicker' in window,
            note: !('showOpenFilePicker' in window) ? 'Chromium only' : undefined
        }),
        "Origin Private File System": () => ({
            supported: navigator.storage && 'getDirectory' in navigator.storage
        })
    },

    "Network APIs": {
        "Fetch API": () => ({
            supported: typeof fetch !== 'undefined'
        }),
        "XMLHttpRequest": () => ({
            supported: typeof XMLHttpRequest !== 'undefined'
        }),
        "WebSocket": () => ({
            supported: typeof WebSocket !== 'undefined'
        }),
        "EventSource (SSE)": () => ({
            supported: typeof EventSource !== 'undefined'
        }),
        "Beacon API": () => ({
            supported: 'sendBeacon' in navigator
        }),
        "WebTransport": () => ({
            supported: typeof WebTransport !== 'undefined',
            note: typeof WebTransport === 'undefined' ? 'Experimental' : undefined
        }),
        "Background Fetch": () => ({
            supported: 'BackgroundFetchManager' in window,
            note: !('BackgroundFetchManager' in window) ? 'Chromium only' : undefined
        }),
        "Background Sync": () => ({
            supported: 'SyncManager' in window,
            note: !('SyncManager' in window) ? 'Chromium only' : undefined
        })
    },

    "Media APIs": {
        "Web Audio API": () => ({
            supported: typeof AudioContext !== 'undefined' || typeof webkitAudioContext !== 'undefined',
            note: typeof AudioContext === 'undefined' && typeof webkitAudioContext !== 'undefined' ? 'webkit prefix' : undefined
        }),
        "MediaDevices": () => ({
            supported: 'mediaDevices' in navigator
        }),
        "getUserMedia": () => ({
            supported: navigator.mediaDevices && 'getUserMedia' in navigator.mediaDevices
        }),
        "getDisplayMedia": () => ({
            supported: navigator.mediaDevices && 'getDisplayMedia' in navigator.mediaDevices
        }),
        "MediaRecorder": () => ({
            supported: typeof MediaRecorder !== 'undefined'
        }),
        "Media Session": () => ({
            supported: 'mediaSession' in navigator
        }),
        "Media Capabilities": () => ({
            supported: 'mediaCapabilities' in navigator
        }),
        "MediaSource Extensions": () => ({
            supported: typeof MediaSource !== 'undefined'
        }),
        "Picture-in-Picture": () => ({
            supported: 'pictureInPictureEnabled' in document
        }),
        "Audio Worklet": () => ({
            supported: typeof AudioWorkletNode !== 'undefined'
        }),
        "Web Speech (Recognition)": () => ({
            supported: 'SpeechRecognition' in window || 'webkitSpeechRecognition' in window
        }),
        "Web Speech (Synthesis)": () => ({
            supported: 'speechSynthesis' in window
        }),
        "Encrypted Media Extensions": () => ({
            supported: 'MediaKeys' in window
        })
    },

    "Graphics APIs": {
        "Canvas 2D": () => {
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            return { supported: ctx !== null };
        },
        "WebGL": () => {
            const canvas = document.createElement('canvas');
            const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
            return { supported: gl !== null };
        },
        "WebGL2": () => {
            const canvas = document.createElement('canvas');
            const gl = canvas.getContext('webgl2');
            return { supported: gl !== null };
        },
        "WebGPU": () => ({
            supported: 'gpu' in navigator,
            note: !('gpu' in navigator) ? 'Experimental' : undefined
        }),
        "OffscreenCanvas": () => ({
            supported: typeof OffscreenCanvas !== 'undefined'
        }),
        "ImageBitmap": () => ({
            supported: typeof createImageBitmap !== 'undefined'
        }),
        "CSS Painting API": () => ({
            supported: 'paintWorklet' in CSS,
            note: !('paintWorklet' in CSS) ? 'Houdini API' : undefined
        }),
        "Web Animations": () => ({
            supported: typeof Element.prototype.animate !== 'undefined'
        }),
        "View Transitions": () => ({
            supported: 'startViewTransition' in document,
            note: !('startViewTransition' in document) ? 'Chromium 111+' : undefined
        })
    },

    "Device APIs": {
        "Geolocation": () => ({
            supported: 'geolocation' in navigator
        }),
        "Device Orientation": () => ({
            supported: 'DeviceOrientationEvent' in window
        }),
        "Device Motion": () => ({
            supported: 'DeviceMotionEvent' in window
        }),
        "Accelerometer": () => ({
            supported: 'Accelerometer' in window,
            note: !('Accelerometer' in window) ? 'Requires secure context' : undefined
        }),
        "Gyroscope": () => ({
            supported: 'Gyroscope' in window,
            note: !('Gyroscope' in window) ? 'Requires secure context' : undefined
        }),
        "Magnetometer": () => ({
            supported: 'Magnetometer' in window,
            note: !('Magnetometer' in window) ? 'Chromium only' : undefined
        }),
        "Ambient Light Sensor": () => ({
            supported: 'AmbientLightSensor' in window,
            note: !('AmbientLightSensor' in window) ? 'Limited support' : undefined
        }),
        "Battery Status": () => ({
            supported: 'getBattery' in navigator,
            note: !('getBattery' in navigator) ? 'Chromium only' : undefined
        }),
        "Device Memory": () => ({
            supported: 'deviceMemory' in navigator,
            note: 'deviceMemory' in navigator ? `${navigator.deviceMemory} GB` : 'Chromium only'
        }),
        "Screen Orientation": () => ({
            supported: 'orientation' in screen
        }),
        "Screen Wake Lock": () => ({
            supported: 'wakeLock' in navigator
        }),
        "Vibration": () => ({
            supported: 'vibrate' in navigator
        }),
        "Web MIDI": () => ({
            supported: 'requestMIDIAccess' in navigator
        }),
        "Web Serial": () => ({
            supported: 'serial' in navigator,
            note: !('serial' in navigator) ? 'Chromium only' : undefined
        }),
        "WebHID": () => ({
            supported: 'hid' in navigator,
            note: !('hid' in navigator) ? 'Chromium only' : undefined
        }),
        "WebUSB": () => ({
            supported: 'usb' in navigator,
            note: !('usb' in navigator) ? 'Chromium only' : undefined
        }),
        "Web NFC": () => ({
            supported: 'NDEFReader' in window,
            note: !('NDEFReader' in window) ? 'Android Chrome only' : undefined
        }),
        "Web Bluetooth": () => ({
            supported: 'bluetooth' in navigator,
            note: !('bluetooth' in navigator) ? 'Limited support' : undefined
        }),
        "Gamepad API": () => ({
            supported: 'getGamepads' in navigator
        })
    },

    "Worker APIs": {
        "Web Workers": () => ({
            supported: typeof Worker !== 'undefined'
        }),
        "Shared Workers": () => ({
            supported: typeof SharedWorker !== 'undefined'
        }),
        "Service Worker": () => ({
            supported: 'serviceWorker' in navigator
        }),
        "Worklets": () => ({
            supported: typeof Worklet !== 'undefined'
        })
    },

    "Performance APIs": {
        "Performance API": () => ({
            supported: typeof performance !== 'undefined'
        }),
        "Performance Observer": () => ({
            supported: typeof PerformanceObserver !== 'undefined'
        }),
        "Navigation Timing": () => ({
            supported: typeof PerformanceNavigationTiming !== 'undefined'
        }),
        "Resource Timing": () => ({
            supported: typeof PerformanceResourceTiming !== 'undefined'
        }),
        "User Timing": () => ({
            supported: performance && 'mark' in performance && 'measure' in performance
        }),
        "Long Tasks API": () => ({
            supported: typeof PerformanceLongTaskTiming !== 'undefined'
        }),
        "Intersection Observer": () => ({
            supported: typeof IntersectionObserver !== 'undefined'
        }),
        "Resize Observer": () => ({
            supported: typeof ResizeObserver !== 'undefined'
        }),
        "Mutation Observer": () => ({
            supported: typeof MutationObserver !== 'undefined'
        }),
        "Reporting API": () => ({
            supported: typeof ReportingObserver !== 'undefined'
        }),
        "Compute Pressure": () => ({
            supported: 'PressureObserver' in window,
            note: !('PressureObserver' in window) ? 'Experimental' : undefined
        })
    },

    "Security APIs": {
        "Web Crypto": () => ({
            supported: typeof crypto !== 'undefined' && 'subtle' in crypto
        }),
        "Credentials API": () => ({
            supported: 'credentials' in navigator
        }),
        "Web Authentication": () => ({
            supported: typeof PublicKeyCredential !== 'undefined'
        }),
        "Permissions API": () => ({
            supported: 'permissions' in navigator
        }),
        "Trusted Types": () => ({
            supported: 'trustedTypes' in window
        }),
        "Content Security Policy": () => ({
            supported: typeof SecurityPolicyViolationEvent !== 'undefined'
        })
    },

    "UI & DOM APIs": {
        "Custom Elements": () => ({
            supported: 'customElements' in window
        }),
        "Shadow DOM": () => ({
            supported: 'attachShadow' in Element.prototype
        }),
        "HTML Templates": () => ({
            supported: 'content' in document.createElement('template')
        }),
        "Pointer Events": () => ({
            supported: 'PointerEvent' in window
        }),
        "Touch Events": () => ({
            supported: 'ontouchstart' in window || navigator.maxTouchPoints > 0
        }),
        "Pointer Lock": () => ({
            supported: 'requestPointerLock' in Element.prototype
        }),
        "Fullscreen API": () => ({
            supported: 'fullscreenEnabled' in document || 'webkitFullscreenEnabled' in document
        }),
        "Selection API": () => ({
            supported: typeof Selection !== 'undefined'
        }),
        "Clipboard API": () => ({
            supported: 'clipboard' in navigator
        }),
        "Clipboard (read)": async () => {
            if (!navigator.clipboard) return { supported: false };
            return { supported: 'read' in navigator.clipboard };
        },
        "Clipboard (write)": async () => {
            if (!navigator.clipboard) return { supported: false };
            return { supported: 'write' in navigator.clipboard };
        },
        "Drag and Drop": () => ({
            supported: 'draggable' in document.createElement('div')
        }),
        "EditContext": () => ({
            supported: 'EditContext' in window,
            note: !('EditContext' in window) ? 'Experimental' : undefined
        }),
        "Virtual Keyboard": () => ({
            supported: 'virtualKeyboard' in navigator,
            note: !('virtualKeyboard' in navigator) ? 'Chromium only' : undefined
        }),
        "Popover API": () => ({
            supported: 'popover' in HTMLElement.prototype
        }),
        "Dialog Element": () => ({
            supported: typeof HTMLDialogElement !== 'undefined'
        })
    },

    "Notifications & Messaging": {
        "Notifications API": () => ({
            supported: 'Notification' in window
        }),
        "Push API": () => ({
            supported: 'PushManager' in window
        }),
        "Channel Messaging": () => ({
            supported: typeof MessageChannel !== 'undefined'
        }),
        "Broadcast Channel": () => ({
            supported: typeof BroadcastChannel !== 'undefined'
        }),
        "postMessage": () => ({
            supported: 'postMessage' in window
        })
    },

    "Navigation & History": {
        "History API": () => ({
            supported: 'pushState' in history
        }),
        "Navigation API": () => ({
            supported: 'navigation' in window,
            note: !('navigation' in window) ? 'Chromium 102+' : undefined
        }),
        "URL API": () => ({
            supported: typeof URL !== 'undefined'
        }),
        "URLSearchParams": () => ({
            supported: typeof URLSearchParams !== 'undefined'
        }),
        "URLPattern": () => ({
            supported: typeof URLPattern !== 'undefined',
            note: typeof URLPattern === 'undefined' ? 'Limited support' : undefined
        })
    },

    "Sharing & Content": {
        "Share API": () => ({
            supported: 'share' in navigator
        }),
        "Web Share Target": () => ({
            supported: 'share' in navigator && 'canShare' in navigator
        }),
        "Badging API": () => ({
            supported: 'setAppBadge' in navigator,
            note: !('setAppBadge' in navigator) ? 'PWA context' : undefined
        }),
        "Content Index": () => ({
            supported: 'ContentIndex' in window,
            note: !('ContentIndex' in window) ? 'PWA context' : undefined
        }),
        "Contact Picker": () => ({
            supported: 'contacts' in navigator,
            note: !('contacts' in navigator) ? 'Android Chrome only' : undefined
        })
    },

    "Streams & Encoding": {
        "Streams API": () => ({
            supported: typeof ReadableStream !== 'undefined'
        }),
        "WritableStream": () => ({
            supported: typeof WritableStream !== 'undefined'
        }),
        "TransformStream": () => ({
            supported: typeof TransformStream !== 'undefined'
        }),
        "Compression Streams": () => ({
            supported: typeof CompressionStream !== 'undefined'
        }),
        "TextEncoder/Decoder": () => ({
            supported: typeof TextEncoder !== 'undefined' && typeof TextDecoder !== 'undefined'
        }),
        "Encoding API (streams)": () => ({
            supported: typeof TextEncoderStream !== 'undefined'
        }),
        "Blob": () => ({
            supported: typeof Blob !== 'undefined'
        }),
        "File API": () => ({
            supported: typeof File !== 'undefined' && typeof FileReader !== 'undefined'
        }),
        "FileReader": () => ({
            supported: typeof FileReader !== 'undefined'
        }),
        "ArrayBuffer": () => ({
            supported: typeof ArrayBuffer !== 'undefined'
        }),
        "DataView": () => ({
            supported: typeof DataView !== 'undefined'
        }),
        "Typed Arrays": () => ({
            supported: typeof Uint8Array !== 'undefined'
        })
    },

    "Payment APIs": {
        "Payment Request": () => ({
            supported: 'PaymentRequest' in window
        }),
        "Payment Handler": () => ({
            supported: 'PaymentManager' in window,
            note: !('PaymentManager' in window) ? 'Limited support' : undefined
        })
    },

    "Extended/Experimental": {
        "WebXR": () => ({
            supported: 'xr' in navigator,
            note: !('xr' in navigator) ? 'VR/AR devices' : undefined
        }),
        "Presentation API": () => ({
            supported: 'presentation' in navigator,
            note: !('presentation' in navigator) ? 'Cast-like APIs' : undefined
        }),
        "Remote Playback": () => ({
            supported: 'remote' in HTMLMediaElement.prototype
        }),
        "Window Management": () => ({
            supported: 'getScreenDetails' in window,
            note: !('getScreenDetails' in window) ? 'Multi-screen' : undefined
        }),
        "Document Picture-in-Picture": () => ({
            supported: 'documentPictureInPicture' in window,
            note: !('documentPictureInPicture' in window) ? 'Chromium only' : undefined
        }),
        "EyeDropper": () => ({
            supported: 'EyeDropper' in window,
            note: !('EyeDropper' in window) ? 'Chromium only' : undefined
        }),
        "File Handling": () => ({
            supported: 'launchQueue' in window,
            note: !('launchQueue' in window) ? 'PWA only' : undefined
        }),
        "Launch Handler": () => ({
            supported: 'LaunchParams' in window,
            note: !('LaunchParams' in window) ? 'PWA only' : undefined
        }),
        "Idle Detection": () => ({
            supported: 'IdleDetector' in window,
            note: !('IdleDetector' in window) ? 'Chromium only' : undefined
        }),
        "Keyboard Lock": () => ({
            supported: 'keyboard' in navigator && 'lock' in navigator.keyboard,
            note: !('keyboard' in navigator) ? 'Fullscreen only' : undefined
        }),
        "Local Font Access": () => ({
            supported: 'queryLocalFonts' in window,
            note: !('queryLocalFonts' in window) ? 'Chromium only' : undefined
        }),
        "Screen Capture": () => ({
            supported: navigator.mediaDevices && 'getDisplayMedia' in navigator.mediaDevices
        }),
        "Scheduler API": () => ({
            supported: 'scheduler' in window
        }),
        "Task Attribution": () => ({
            supported: typeof TaskAttributionTiming !== 'undefined'
        }),
        "Web Codecs (Video)": () => ({
            supported: typeof VideoEncoder !== 'undefined'
        }),
        "Web Codecs (Audio)": () => ({
            supported: typeof AudioEncoder !== 'undefined'
        }),
        "Web Locks": () => ({
            supported: 'locks' in navigator
        }),
        "Prioritized Task Scheduling": () => ({
            supported: 'scheduler' in window && 'postTask' in scheduler
        })
    },

    "CSS APIs": {
        "CSSOM": () => ({
            supported: typeof CSSStyleSheet !== 'undefined'
        }),
        "Constructable Stylesheets": () => ({
            supported: 'adoptedStyleSheets' in document
        }),
        "CSS Typed OM": () => ({
            supported: 'attributeStyleMap' in Element.prototype
        }),
        "CSS Properties & Values": () => ({
            supported: CSS && 'registerProperty' in CSS
        }),
        "CSS.supports": () => ({
            supported: CSS && 'supports' in CSS
        }),
        "CSS Font Loading": () => ({
            supported: 'fonts' in document
        }),
        "CSS Container Queries": () => ({
            supported: CSS && CSS.supports && CSS.supports('container-type', 'inline-size')
        }),
        "@layer support": () => ({
            supported: CSS && CSS.supports && CSS.supports('@layer test { }')
        }),
        "Subgrid": () => ({
            supported: CSS && CSS.supports && CSS.supports('grid-template-columns', 'subgrid')
        }),
        ":has() selector": () => ({
            supported: CSS && CSS.supports && CSS.supports('selector(:has(a))')
        }),
        "color-mix()": () => ({
            supported: CSS && CSS.supports && CSS.supports('color', 'color-mix(in srgb, red, blue)')
        }),
        "Scroll-driven Animations": () => ({
            supported: CSS && CSS.supports && CSS.supports('animation-timeline', 'scroll()'),
            note: !(CSS && CSS.supports && CSS.supports('animation-timeline', 'scroll()')) ? 'Chromium 115+' : undefined
        })
    },

    "JavaScript Features": {
        "ES Modules": () => ({
            supported: 'noModule' in document.createElement('script')
        }),
        "Import Maps": () => ({
            supported: HTMLScriptElement.supports && HTMLScriptElement.supports('importmap')
        }),
        "Dynamic Import": async () => {
            try {
                await import('data:text/javascript,export default 1');
                return { supported: true };
            } catch {
                return { supported: false };
            }
        },
        "Top-level Await": () => ({
            supported: true, // If we're running, it's supported in modules
            note: 'Module context'
        }),
        "WeakRef": () => ({
            supported: typeof WeakRef !== 'undefined'
        }),
        "FinalizationRegistry": () => ({
            supported: typeof FinalizationRegistry !== 'undefined'
        }),
        "BigInt": () => ({
            supported: typeof BigInt !== 'undefined'
        }),
        "globalThis": () => ({
            supported: typeof globalThis !== 'undefined'
        }),
        "Optional Chaining": () => {
            try {
                eval('const x = null?.foo');
                return { supported: true };
            } catch {
                return { supported: false };
            }
        },
        "Nullish Coalescing": () => {
            try {
                eval('const x = null ?? "default"');
                return { supported: true };
            } catch {
                return { supported: false };
            }
        },
        "Private Class Fields": () => {
            try {
                eval('class C { #x = 1 }');
                return { supported: true };
            } catch {
                return { supported: false };
            }
        },
        "Static Class Blocks": () => {
            try {
                eval('class C { static { } }');
                return { supported: true };
            } catch {
                return { supported: false };
            }
        },
        "Temporal (Stage 3)": () => ({
            supported: typeof Temporal !== 'undefined',
            note: typeof Temporal === 'undefined' ? 'Proposal' : undefined
        }),
        "Iterator Helpers": () => ({
            supported: typeof Iterator !== 'undefined' && 'from' in Iterator,
            note: !(typeof Iterator !== 'undefined') ? 'Proposal' : undefined
        }),
        "Array.at()": () => ({
            supported: 'at' in Array.prototype
        }),
        "Object.hasOwn()": () => ({
            supported: 'hasOwn' in Object
        }),
        "structuredClone": () => ({
            supported: typeof structuredClone !== 'undefined'
        }),
        "Atomics.waitAsync": () => ({
            supported: typeof Atomics !== 'undefined' && 'waitAsync' in Atomics
        }),
        "Array.fromAsync": () => ({
            supported: 'fromAsync' in Array,
            note: !('fromAsync' in Array) ? 'ES2024' : undefined
        }),
        "Promise.withResolvers": () => ({
            supported: 'withResolvers' in Promise,
            note: !('withResolvers' in Promise) ? 'ES2024' : undefined
        }),
        "RegExp v flag": () => {
            try {
                new RegExp('.', 'v');
                return { supported: true };
            } catch {
                return { supported: false, note: 'ES2024' };
            }
        }
    }
};
