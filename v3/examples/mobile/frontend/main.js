import * as Runtime from "@wailsio/runtime";
import { SystemService } from "./bindings/mobile-example";

const { Events, Dialogs, Clipboard, Screens, IOS, Android } = Runtime;

const $ = (id) => document.getElementById(id);
const show = (id, value) =>
    ($(id).textContent = typeof value === "string" ? value : JSON.stringify(value, null, 2));

// ---- Platform detection -------------------------------------------------
// The Android JavascriptInterface bridge exposes window.wails; iOS exposes a
// WKScriptMessageHandler. A cheap synchronous check keeps the UI responsive
// on first paint and drives which platform-specific controls are shown.
function detectPlatform() {
    if (typeof window.wails?.platform === "function") {
        try { return window.wails.platform(); } catch (_) { /* ignore */ }
    }
    if (window.webkit?.messageHandlers?.external) return "ios";
    return "desktop";
}

const platform = detectPlatform();

const isMobile = platform === "ios" || platform === "android";

function applyPlatform() {
    $("platform").textContent = platform;
    document.querySelectorAll(".platform-only").forEach((el) => {
        el.style.display = el.dataset.platform === platform ? "" : "none";
    });
    // .mobile-only blocks are shown on both iOS and Android (the Go side routes
    // each common:* event to the matching platform bridge), hidden on desktop.
    document.querySelectorAll(".mobile-only").forEach((el) => {
        el.style.display = isMobile ? "" : "none";
    });
    const hints = {
        ios: "Running on iOS — haptics and WKWebView toggles are available.",
        android: "Running on Android — vibration and toasts are available.",
        desktop: "Running on desktop.",
    };
    $("nativeHint").textContent = hints[platform] || "";
}

// ---- Tabs ---------------------------------------------------------------
document.querySelectorAll(".tab").forEach((tab) => {
    tab.addEventListener("click", () => {
        document.querySelectorAll(".tab").forEach((t) => t.classList.remove("active"));
        document.querySelectorAll(".pane").forEach((p) => p.classList.remove("active"));
        tab.classList.add("active");
        $("pane-" + tab.dataset.tab).classList.add("active");
    });
});

// ---- Bindings (JS -> Go) ------------------------------------------------
$("btnGreet").addEventListener("click", async () => {
    try {
        show("bindingsOut", await SystemService.Greet($("greetName").value));
    } catch (e) {
        show("bindingsOut", "Error: " + (e?.message ?? e));
    }
});

$("btnAdd").addEventListener("click", async () => {
    const a = parseInt($("addA").value || "0", 10);
    const b = parseInt($("addB").value || "0", 10);
    show("bindingsOut", `${a} + ${b} = ${await SystemService.Add(a, b)}`);
});

$("btnDivide").addEventListener("click", async () => {
    const a = parseFloat($("divA").value || "0");
    const b = parseFloat($("divB").value || "0");
    try {
        show("bindingsOut", `${a} ÷ ${b} = ${await SystemService.Divide(a, b)}`);
    } catch (e) {
        // Demonstrates a Go error surfacing in JS
        show("bindingsOut", "Go returned an error: " + (e?.message ?? e));
    }
});

// ---- Events (Go <-> JS) -------------------------------------------------
// A single Emit("name", value) delivers the value on event.data; multi-arg
// emits deliver an array. Handle both.
const eventValue = (e) => (Array.isArray(e?.data) ? e.data[0] : (e?.data ?? e));

Events.On("time", (e) => {
    $("clock").textContent = eventValue(e);
});

Events.On("pong", (e) => {
    show("eventsOut", "pong @ " + eventValue(e));
});

$("btnPing").addEventListener("click", () => {
    Events.Emit("ping", { from: platform });
    $("eventsOut").textContent = "ping sent, waiting for pong…";
});

// ---- System events ------------------------------------------------------
// Native OS signals arrive in Go as common: application events; the app
// forwards them here as "sys:*" custom events. We keep the latest value per
// category as a state card plus a rolling event log.
const SYS = {
    battery: { icon: "🔋", label: "Battery" },
    network: { icon: "📶", label: "Network" },
    theme: { icon: "🎨", label: "Theme" },
    lock: { icon: "🔒", label: "Lock" },
    memory: { icon: "⚠️", label: "Memory" },
};
const sysState = {};

function fmtSystem(kind, d) {
    d = d || {};
    switch (kind) {
        case "battery": {
            const pct = typeof d.level === "number" && d.level >= 0
                ? Math.round(d.level * 100) + "%" : "—";
            const tags = [d.state, d.lowPowerMode ? "low-power" : null].filter(Boolean).join(", ");
            return tags ? `${pct} · ${tags}` : pct;
        }
        case "network": {
            if (!d.connected) return "offline";
            const tags = [d.type,
                d.metered ? "metered" : null,
                typeof d.signal === "number" ? `${d.signal} dBm` : null,
                d.expensive ? "expensive" : null,
                d.constrained ? "low-data" : null].filter(Boolean).join(", ");
            return tags ? `online · ${tags}` : "online";
        }
        case "theme": return d.isDarkMode ? "dark" : "light";
        case "lock": return d.locked ? "locked" : "unlocked";
        case "memory": return "low-memory warning";
        default: return JSON.stringify(d);
    }
}

function renderSysState() {
    const target = $("sysState");
    const keys = Object.keys(SYS).filter((k) => k in sysState && k !== "memory");
    if (keys.length === 0) {
        target.innerHTML = '<p class="hint">Waiting for system events…</p>';
        return;
    }
    target.innerHTML = `<div class="metric-card">${keys.map((k) =>
        `<div class="metric-row"><span class="metric-key">${SYS[k].icon} ${SYS[k].label}</span>` +
        `<span class="metric-val">${esc(fmtSystem(k, sysState[k]))}</span></div>`).join("")}</div>`;
}

function logSystem(kind, d) {
    const log = $("sysLog");
    if (!log) return;
    const line = document.createElement("div");
    line.className = "event-line";
    line.innerHTML =
        `<span class="event-time">${esc(new Date().toLocaleTimeString())}</span>` +
        `<span class="event-name">${SYS[kind].icon} sys:${kind}</span>` +
        `<span class="event-detail">${esc(fmtSystem(kind, d))}</span>`;
    log.prepend(line);
    while (log.childElementCount > 8) log.removeChild(log.lastChild);
}

Object.keys(SYS).forEach((kind) => {
    Events.On("sys:" + kind, (e) => {
        const d = eventValue(e);
        if (kind !== "memory") sysState[kind] = d;
        renderSysState();
        logSystem(kind, d);
    });
});

// ---- Dialogs ------------------------------------------------------------
async function runDialog(kind, fn) {
    try {
        const result = await fn();
        show("dialogsOut", `${kind} → ${result === undefined ? "dismissed" : result}`);
    } catch (e) {
        show("dialogsOut", `${kind} error: ${e?.message ?? e}`);
    }
}

$("btnInfo").addEventListener("click", () =>
    runDialog("Info", () => Dialogs.Info({ Title: "Hello", Message: "This is an info dialog." })));
$("btnWarning").addEventListener("click", () =>
    runDialog("Warning", () => Dialogs.Warning({ Title: "Careful", Message: "This is a warning." })));
$("btnError").addEventListener("click", () =>
    runDialog("Error", () => Dialogs.Error({ Title: "Oops", Message: "Something went wrong." })));
$("btnQuestion").addEventListener("click", () =>
    runDialog("Question", () => Dialogs.Question({
        Title: "Confirm",
        Message: "Do you like Wails?",
        Buttons: [{ Label: "Yes", IsDefault: true }, { Label: "No", IsCancel: true }],
    })));

// ---- Clipboard ----------------------------------------------------------
$("btnClipSet").addEventListener("click", async () => {
    await Clipboard.SetText($("clipText").value);
    show("clipOut", "Copied: " + $("clipText").value);
});
$("btnClipGet").addEventListener("click", async () => {
    const text = await Clipboard.Text();
    show("clipOut", "Clipboard: " + (text || "(empty)"));
});

// ---- Screens ------------------------------------------------------------
const esc = (s) => String(s).replace(/[&<>"]/g, (c) =>
    ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]));

// Render screens as labelled key/value cards rather than raw JSON, so the
// metrics read cleanly and fit without needing a scrollable section.
function renderScreens(screens) {
    const out = $("screensOut");
    if (!Array.isArray(screens) || screens.length === 0) {
        out.innerHTML = '<p class="hint">No screens reported.</p>';
        return;
    }
    const dim = (r) => (r ? `${r.Width} × ${r.Height}` : "—");
    out.innerHTML = screens.map((s, i) => {
        const rows = [
            ["Name", s.Name || s.ID || `Screen ${i + 1}`],
            ["Scale", `${s.ScaleFactor}×`],
            ["Size", `${dim(s.Size)} pt`],
            ["Physical", `${dim(s.PhysicalBounds)} px`],
            ["Work area", `${dim(s.WorkArea)} pt`],
            ["Primary", s.IsPrimary ? "Yes" : "No"],
        ];
        const body = rows.map(([k, v]) =>
            `<div class="metric-row"><span class="metric-key">${esc(k)}</span>` +
            `<span class="metric-val">${esc(v)}</span></div>`).join("");
        return `<div class="metric-card">${body}</div>`;
    }).join("");
}

$("btnScreens").addEventListener("click", async () => {
    try {
        renderScreens(await Screens.GetAll());
    } catch (e) {
        $("screensOut").innerHTML = `<p class="hint">Error: ${esc(e?.message ?? e)}</p>`;
    }
});

// ---- Device info --------------------------------------------------------
// Render an object as labelled rows (camelCase keys → Title Case, booleans →
// Yes/No), matching the screen metrics card instead of dumping JSON.
function renderKeyVals(target, obj) {
    if (!obj || typeof obj !== "object") {
        target.innerHTML = '<p class="hint">No data.</p>';
        return;
    }
    const pretty = (k) => k.replace(/([a-z0-9])([A-Z])/g, "$1 $2")
        .replace(/^./, (c) => c.toUpperCase());
    const fmt = (v) =>
        typeof v === "boolean" ? (v ? "Yes" : "No")
            : v !== null && typeof v === "object" ? JSON.stringify(v)
                : String(v);
    const rows = Object.entries(obj).map(([k, v]) =>
        `<div class="metric-row"><span class="metric-key">${esc(pretty(k))}</span>` +
        `<span class="metric-val">${esc(fmt(v))}</span></div>`).join("");
    target.innerHTML = `<div class="metric-card">${rows}</div>`;
}

$("btnDevice").addEventListener("click", async () => {
    try {
        if (platform === "ios") renderKeyVals($("deviceOut"), await IOS.Device.Info());
        else if (platform === "android") renderKeyVals($("deviceOut"), await Android.Device.Info());
        else $("deviceOut").innerHTML = '<p class="hint">Device info is only available on iOS and Android.</p>';
    } catch (e) {
        $("deviceOut").innerHTML = `<p class="hint">Error: ${esc(e?.message ?? e)}</p>`;
    }
});

// ---- Native: iOS --------------------------------------------------------
document.querySelectorAll("[data-haptic]").forEach((btn) => {
    btn.addEventListener("click", () => IOS.Haptics.Impact(btn.dataset.haptic).catch(() => {}));
});
// Scroll is enabled by default so content taller than the viewport (e.g. the
// System tab's screen/device cards) is always reachable. Toggling scroll sets the
// native WKWebView scrollView's scrollEnabled and appends/removes a tall filler
// section inside the Native pane (next to the toggle) so there is an obvious area
// to scroll and bounce against — without cluttering the other tabs.
function setScrollEnabled(enabled) {
    Events.Emit("ios:setScrollEnabled", { enabled });
    let filler = $("scrollFiller");
    if (enabled && !filler) {
        filler = document.createElement("section");
        filler.id = "scrollFiller";
        filler.className = "scroll-filler";
        filler.innerHTML =
            "<h2>Scroll test</h2>" +
            '<p class="hint">Scroll is on, so this tall section gives the page ' +
            "something to scroll and bounce against. Turn Scroll off to remove it.</p>" +
            Array.from({ length: 12 }, (_, i) =>
                `<div class="scroll-block">Scroll block ${i + 1} / 12</div>`).join("");
        $("pane-native").appendChild(filler);
    } else if (!enabled && filler) {
        filler.remove();
    }
}

const iosScroll = $("iosScroll");
const iosBounce = $("iosBounce");

function setBounceEnabled(enabled) {
    Events.Emit("ios:setBounceEnabled", { enabled });
}

// Bounce acts on the native WKWebView scrollView, which only tracks drags while
// scroll is enabled — so a disabled scrollView can never bounce. Reflect that by
// dimming and disabling the bounce toggle until scroll is turned on.
function syncBounceAvailability() {
    if (!iosBounce) return;
    const scrollOn = iosScroll ? iosScroll.checked : false;
    iosBounce.disabled = !scrollOn;
    iosBounce.closest(".switch")?.classList.toggle("disabled", !scrollOn);
}

if (iosScroll) {
    // Apply the initial state (scroll off by default) on the native side.
    setScrollEnabled(iosScroll.checked);
    iosScroll.addEventListener("change", (e) => {
        setScrollEnabled(e.target.checked);
        syncBounceAvailability();
    });
}
if (iosBounce) {
    // Push the initial bounce state too (the change handler alone never fires for
    // the default), then keep it in sync with the scroll toggle.
    setBounceEnabled(iosBounce.checked);
    iosBounce.addEventListener("change", (e) => setBounceEnabled(e.target.checked));
}
syncBounceAvailability();

// ---- Native: Android ----------------------------------------------------
$("btnVibrate")?.addEventListener("click", async () => {
    await Android.Haptics.Vibrate(200);
    show("nativeOut", "Vibrated 200ms");
});
$("btnToast")?.addEventListener("click", async () => {
    await Android.Toast.Show("Hello from Wails 👋");
    show("nativeOut", "Toast shown");
});

// ---- Mobile features (iOS + Android) ------------------------------------
// Each control emits a "common:*" event; the Go side routes it to the matching
// platform bridge. Asynchronous results arrive back as "common:*" events.
const logMobile = (msg) => console.log("[mobile]", msg);

$("btnShare")?.addEventListener("click", () => {
    Events.Emit("common:share", {
        text: "Check out Wails — the Go + Web framework for native apps.",
        url: "https://wails.io",
    });
    show("outShare", "Opened share sheet");
});

$("btnOpenUrl")?.addEventListener("click", () => {
    const url = $("openUrl").value.trim();
    if (!url) return;
    Events.Emit("common:openURL", { url });
    show("outShare", "Opening " + url);
});

$("mfKeepAwake")?.addEventListener("change", (e) => {
    Events.Emit("common:keepAwake", { enabled: e.target.checked });
    show("outScreen", "Keep awake: " + (e.target.checked ? "on" : "off"));
});

const mfTorch = $("mfTorch");
mfTorch?.addEventListener("change", (e) => {
    Events.Emit("common:torch", { enabled: e.target.checked });
});
Events.On("common:torch", (e) => {
    const d = eventValue(e) || {};
    if (d.available === false) {
        if (mfTorch) mfTorch.checked = false;
        show("outScreen", "Torch not available on this device");
    } else {
        show("outScreen", "Torch: " + (d.on ? "on" : "off"));
    }
});

// Brightness (0-100 in the slider, 0-1 over the wire)
const mfBrightness = $("mfBrightness");
mfBrightness?.addEventListener("change", (e) => {
    Events.Emit("common:setBrightness", { value: e.target.value / 100 });
    show("outScreen", "Brightness set to " + e.target.value + "%");
});
$("btnGetBrightness")?.addEventListener("click", () => Events.Emit("common:getBrightness", {}));
Events.On("common:brightness", (e) => {
    const v = (eventValue(e) || {}).value;
    if (typeof v === "number" && v >= 0) {
        if (mfBrightness) mfBrightness.value = Math.round(v * 100);
        show("outScreen", "Brightness is " + Math.round(v * 100) + "%");
    }
});

// Safe-area insets + app info → metrics card
$("btnSafeArea")?.addEventListener("click", () => Events.Emit("common:getSafeArea", {}));
Events.On("common:safeArea", (e) => {
    renderKeyVals($("mobileMetrics"), eventValue(e));
});
$("btnAppInfo")?.addEventListener("click", () => Events.Emit("common:getAppInfo", {}));
Events.On("common:appInfo", (e) => {
    renderKeyVals($("mobileMetrics"), eventValue(e));
});

// Orientation
document.querySelectorAll("[data-orient]").forEach((btn) => {
    btn.addEventListener("click", () => {
        Events.Emit("common:setOrientation", { mode: btn.dataset.orient });
        setTimeout(() => Events.Emit("common:getOrientation", {}), 400);
    });
});
Events.On("common:orientation", (e) => {
    const o = (eventValue(e) || {}).orientation;
    if (o) $("mfOrientation").textContent = o;
});

// Status bar
document.querySelectorAll("[data-statusbar]").forEach((btn) => {
    btn.addEventListener("click", () => {
        const v = btn.dataset.statusbar;
        if (v === "hidden") Events.Emit("common:setStatusBar", { hidden: true });
        else if (v === "shown") Events.Emit("common:setStatusBar", { hidden: false });
        else Events.Emit("common:setStatusBar", { style: v });
        show("outStatusbar", "Status bar: " + v);
    });
});

// Biometric authentication
$("btnBiometric")?.addEventListener("click", () => {
    Events.Emit("common:authenticate", { reason: "Unlock the kitchen sink" });
    show("outSecurity", "Requesting authentication…");
});
Events.On("common:biometric", (e) => {
    const d = eventValue(e) || {};
    show("outSecurity", d.ok ? "✓ Authenticated" : "✗ Authentication failed: " + (d.error || "unknown"));
});

// Local notification
$("btnNotify")?.addEventListener("click", () => {
    Events.Emit("common:notify", {
        title: "Wails Kitchen Sink",
        body: "This is a local notification 👋",
        delay: 2,
    });
    show("outSecurity", "Scheduling notification…");
});
Events.On("common:notification", (e) => {
    const d = eventValue(e) || {};
    if (d.presented) { show("outSecurity", "📬 Notification shown (foreground)"); return; }
    if (d.tapped) { show("outSecurity", "👆 Notification tapped"); return; }
    show("outSecurity", d.ok ? "Notification posted" + (d.scheduled ? " (in " + d.scheduled + "s)" : "")
                   : "Notification failed: " + (d.error || "denied"));
});

// Secure storage
$("btnSecSet")?.addEventListener("click", () => {
    Events.Emit("common:secureSet", { key: $("secKey").value, value: $("secVal").value });
    show("outSecure", "Saved '" + $("secKey").value + "' securely");
});
$("btnSecGet")?.addEventListener("click", () =>
    Events.Emit("common:secureGet", { key: $("secKey").value }));
$("btnSecDel")?.addEventListener("click", () => {
    Events.Emit("common:secureDelete", { key: $("secKey").value });
    show("outSecure", "Deleted '" + $("secKey").value + "'");
});
Events.On("common:secureValue", (e) => {
    const d = eventValue(e) || {};
    show("outSecure", "Loaded '" + d.key + "' = " + (d.value ? "\"" + d.value + "\"" : "(empty)"));
});

// ---- Hardware: sensors & device capabilities ----------------------------
const logHardware = (msg) => console.log("[hardware]", msg);

// Haptics
document.querySelectorAll("[data-haptic2]").forEach((btn) => {
    btn.addEventListener("click", () => {
        Events.Emit("common:haptic", { type: btn.dataset.haptic2 });
        show("outHaptics", "Haptic: " + btn.dataset.haptic2);
    });
});

// Location (one-shot)
$("btnLocation")?.addEventListener("click", () => {
    Events.Emit("common:getLocation", {});
    show("outLocation", "Requesting location…");
});
Events.On("common:location", (e) => {
    const d = eventValue(e) || {};
    if (d.error) {
        show("outLocation", "Location error: " + d.error);
    } else {
        show("outLocation", `Location: ${d.lat?.toFixed(5)}, ${d.lng?.toFixed(5)} (±${Math.round(d.accuracy)}m)`);
    }
});

// Accelerometer stream
$("mfMotion")?.addEventListener("change", (e) => {
    Events.Emit("common:watchMotion", { enabled: e.target.checked });
    show("outMotion", "Accelerometer: " + (e.target.checked ? "on" : "off"));
});
Events.On("common:motion", (e) => {
    const d = eventValue(e) || {};
    if (d.available === false) {
        show("outMotion", "Accelerometer not available");
        if ($("mfMotion")) $("mfMotion").checked = false;
        return;
    }
    show("outMotion", `Motion  x:${d.x?.toFixed(2)}  y:${d.y?.toFixed(2)}  z:${d.z?.toFixed(2)}`);
});

// Proximity
$("mfProximity")?.addEventListener("change", (e) => {
    Events.Emit("common:watchProximity", { enabled: e.target.checked });
    show("outMotion", "Proximity: " + (e.target.checked ? "watching" : "off"));
});
Events.On("common:proximity", (e) => {
    const d = eventValue(e) || {};
    if (d.available === false) {
        show("outMotion", "Proximity sensor not available");
        if ($("mfProximity")) $("mfProximity").checked = false;
        return;
    }
    show("outMotion", "Proximity: " + (d.near ? "near" : "far"));
});

// Text-to-speech
$("btnSpeak")?.addEventListener("click", () => {
    Events.Emit("common:speak", { text: $("speakText").value });
    show("outSpeech", "Speaking…");
});
$("btnStopSpeak")?.addEventListener("click", () => {
    Events.Emit("common:stopSpeak", {});
    show("outSpeech", "Speech stopped");
});

// Device state queries → metrics card
const bytesToGB = (b) => (b / 1e9).toFixed(2) + " GB";
$("btnStorage")?.addEventListener("click", () => Events.Emit("common:getStorage", {}));
Events.On("common:storage", (e) => {
    const d = eventValue(e) || {};
    renderKeyVals($("hardwareMetrics"), {
        free: bytesToGB(d.free || 0),
        total: bytesToGB(d.total || 0),
        used: bytesToGB((d.total || 0) - (d.free || 0)),
    });
});
$("btnPower")?.addEventListener("click", () => Events.Emit("common:getPower", {}));
Events.On("common:power", (e) => {
    const d = eventValue(e) || {};
    renderKeyVals($("hardwareMetrics"), {
        battery: typeof d.level === "number" && d.level >= 0 ? Math.round(d.level * 100) + "%" : "unknown",
        charging: !!d.charging,
        lowPowerMode: !!d.lowPower,
    });
});
$("btnNetwork")?.addEventListener("click", () => Events.Emit("common:getNetwork", {}));
Events.On("common:network", (e) => {
    const d = eventValue(e) || {};
    renderKeyVals($("hardwareMetrics"), { connected: !!d.connected, type: d.type || "none" });
});

// Keyboard insets
$("mfKeyboard")?.addEventListener("change", (e) => {
    Events.Emit("common:watchKeyboard", { enabled: e.target.checked });
    show("outKeyboard", "Keyboard watch: " + (e.target.checked ? "on" : "off"));
});
Events.On("common:keyboard", (e) => {
    const d = eventValue(e) || {};
    show("outKeyboard", `Keyboard ${d.visible ? "shown" : "hidden"} (height ${d.height || 0}px)`);
});

// Screen-capture protection / detection
$("mfScreenProtect")?.addEventListener("change", (e) => {
    Events.Emit("common:setScreenProtect", { enabled: e.target.checked });
    show("outKeyboard", "Screen protection: " + (e.target.checked ? "on" : "off"));
});
Events.On("common:screenCapture", (e) => {
    const d = eventValue(e) || {};
    if (d.screenshot) show("outKeyboard", "⚠ Screenshot detected");
    else if (d.recording !== undefined) show("outKeyboard", "Screen recording: " + (d.recording ? "active" : "inactive"));
    else show("outKeyboard", "Screen capture " + (d.protected ? "blocked (FLAG_SECURE)" : "allowed"));
});

// ---- Camera & background ------------------------------------------------
const logCamera = (msg) => show("outCamera", msg);

// Photo/Video pill tabs drive a single capture button and a switched preview
// area: photos show in an <img>, recorded videos play back in a <video>.
const captureImg = $("captureImg");
const captureVideo = $("captureVideo");
const capturePlaceholder = $("capturePlaceholder");
let captureModeVal = "photo";

function setCaptureMode(mode) {
    captureModeVal = mode;
    document.querySelectorAll("#captureMode .pill-tab").forEach((t) =>
        t.classList.toggle("active", t.dataset.mode === mode));
    if ($("btnCapture")) $("btnCapture").textContent = mode === "video" ? "Record video" : "Take photo";
    const hasPhoto = !!(captureImg && captureImg.getAttribute("src"));
    const hasVideo = !!(captureVideo && captureVideo.getAttribute("src"));
    if (mode === "video") {
        if (captureImg) captureImg.hidden = true;
        if (captureVideo) captureVideo.hidden = !hasVideo;
    } else {
        if (captureVideo) { try { captureVideo.pause(); } catch (_) {} captureVideo.hidden = true; }
        if (captureImg) captureImg.hidden = !hasPhoto;
    }
    if (capturePlaceholder) {
        const has = mode === "video" ? hasVideo : hasPhoto;
        capturePlaceholder.hidden = has;
        capturePlaceholder.textContent = mode === "video"
            ? "No video yet — tap Record video." : "No photo yet — tap Take photo.";
    }
}
document.querySelectorAll("#captureMode .pill-tab").forEach((t) =>
    t.addEventListener("click", () => setCaptureMode(t.dataset.mode)));

$("btnCapture")?.addEventListener("click", () => {
    if (captureModeVal === "video") {
        Events.Emit("common:captureVideo", {});
        logCamera("Opening camera (video)…");
    } else {
        Events.Emit("common:capturePhoto", {});
        logCamera("Opening camera…");
    }
});

Events.On("common:capture", (e) => {
    const d = eventValue(e) || {};
    if (d.error) { logCamera("Capture error: " + d.error); return; }
    if (d.cancelled) { logCamera("Capture cancelled"); return; }
    const kb = d.size ? Math.round(d.size / 1024) : 0;
    if (d.type === "photo") {
        if (d.thumb && captureImg) captureImg.src = d.thumb;
        setCaptureMode("photo"); // reveal the photo preview
        logCamera(`Photo captured (${kb} KB)\n${d.path || ""}`);
    } else {
        // Prefer the streamed URL (served from the cache with Range support, any
        // size); fall back to an inline data URL if that's all we got.
        const vsrc = d.streamUrl || d.dataUrl;
        if (vsrc && captureVideo) captureVideo.src = vsrc;
        setCaptureMode("video"); // reveal the <video> for playback
        logCamera(vsrc
            ? `Video captured (${kb} KB) — tap ▶ to play\n${d.path || ""}`
            : `Video captured (${kb} KB)\n${d.path || ""}`);
    }
});

// Foreground service (Android)
$("mfForegroundService")?.addEventListener("change", (e) => {
    if (e.target.checked) {
        Events.Emit("common:startForegroundService", {
            title: "Wails Kitchen Sink",
            text: "Background work running",
        });
    } else {
        Events.Emit("common:stopForegroundService", {});
    }
});
Events.On("android:foregroundService", (e) => {
    const d = eventValue(e) || {};
    if (d.error) show("outBg", "Foreground service error: " + d.error);
    else show("outBg", "Foreground service: " + (d.running ? "running" : "stopped"));
});

// Background task window (iOS)
$("btnBgTask")?.addEventListener("click", () => {
    Events.Emit("ios:beginBackgroundTask", { seconds: 20 });
    show("outBg", "Requested a background-task window…");
});
Events.On("ios:backgroundTask", (e) => {
    const d = eventValue(e) || {};
    show("outBg", "Background task: " + (d.message || JSON.stringify(d)));
});

// Ask for the current orientation once the page is up.
if (isMobile) setTimeout(() => Events.Emit("common:getOrientation", {}), 600);

applyPlatform();
