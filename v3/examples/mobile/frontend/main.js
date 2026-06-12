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

function applyPlatform() {
    $("platform").textContent = platform;
    document.querySelectorAll(".platform-only").forEach((el) => {
        el.style.display = el.dataset.platform === platform ? "" : "none";
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

// ---- System events (native -> JS) ---------------------------------------
// The framework emits "system:*" custom events from native OS signals on iOS
// and Android: battery/power, network, theme, lock and app lifecycle. We keep
// the latest value per category as a state card and a rolling event log.
const SYS_KINDS = ["battery", "network", "theme", "lock", "appstate", "memory"];
const SYS_ICON = { battery: "🔋", network: "📶", theme: "🎨", lock: "🔒", appstate: "📱", memory: "⚠️" };
const sysState = {};

function fmtSystem(kind, d) {
    if (d == null || typeof d !== "object") return String(d ?? "");
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
        case "theme": return d.dark ? "dark" : "light";
        case "lock": return d.locked ? "locked" : "unlocked";
        case "appstate": return d.state || "—";
        case "memory": return "low-memory warning";
        default: return JSON.stringify(d);
    }
}

const titleCase = (s) => s.charAt(0).toUpperCase() + s.slice(1);

function renderSysState() {
    const target = $("sysState");
    const keys = SYS_KINDS.filter((k) => k in sysState && k !== "memory");
    if (keys.length === 0) {
        target.innerHTML = '<p class="hint">Waiting for system events…</p>';
        return;
    }
    target.innerHTML = `<div class="metric-card">${keys.map((k) =>
        `<div class="metric-row"><span class="metric-key">${SYS_ICON[k]} ${titleCase(k)}</span>` +
        `<span class="metric-val">${esc(fmtSystem(k, sysState[k]))}</span></div>`).join("")}</div>`;
}

function logSystem(kind, d) {
    const log = $("sysLog");
    if (!log) return;
    const line = document.createElement("div");
    line.className = "event-line";
    line.innerHTML =
        `<span class="event-time">${esc(new Date().toLocaleTimeString())}</span>` +
        `<span class="event-name">${SYS_ICON[kind]} system:${kind}</span>` +
        `<span class="event-detail">${esc(fmtSystem(kind, d))}</span>`;
    log.prepend(line);
    while (log.childElementCount > 8) log.removeChild(log.lastChild);
}

SYS_KINDS.forEach((kind) => {
    Events.On("system:" + kind, (e) => {
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

applyPlatform();
