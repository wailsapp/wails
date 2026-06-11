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
$("btnScreens").addEventListener("click", async () => {
    try {
        show("screensOut", await Screens.GetAll());
    } catch (e) {
        show("screensOut", "Error: " + (e?.message ?? e));
    }
});

// ---- Device info --------------------------------------------------------
$("btnDevice").addEventListener("click", async () => {
    try {
        if (platform === "ios") show("deviceOut", await IOS.Device.Info());
        else if (platform === "android") show("deviceOut", await Android.Device.Info());
        else show("deviceOut", "Device info is only available on iOS and Android.");
    } catch (e) {
        show("deviceOut", "Error: " + (e?.message ?? e));
    }
});

// ---- Native: iOS --------------------------------------------------------
document.querySelectorAll("[data-haptic]").forEach((btn) => {
    btn.addEventListener("click", () => IOS.Haptics.Impact(btn.dataset.haptic).catch(() => {}));
});
// Scroll is disabled by default. Enabling it turns the native WKWebView
// scrollView on AND appends a tall filler section so the page actually has
// somewhere to scroll/bounce; disabling removes the filler and the native
// scroll again, so there is nothing to scroll.
function setScrollEnabled(enabled) {
    Events.Emit("ios:setScrollEnabled", { enabled });
    let filler = $("scrollFiller");
    if (enabled && !filler) {
        filler = document.createElement("section");
        filler.id = "scrollFiller";
        filler.className = "scroll-filler";
        filler.innerHTML =
            "<h2>Scroll test</h2>" +
            '<p class="hint">Scrolling is on, so this tall section was added to ' +
            "give the page something to scroll and bounce against. Turn scrolling " +
            "off to remove it.</p>" +
            Array.from({ length: 12 }, (_, i) =>
                `<div class="scroll-block">Scroll block ${i + 1} / 12</div>`).join("");
        document.querySelector("main").appendChild(filler);
    } else if (!enabled && filler) {
        filler.remove();
    }
}

const iosScroll = $("iosScroll");
if (iosScroll) {
    // Apply the initial state (scroll off by default) on the native side.
    setScrollEnabled(iosScroll.checked);
    iosScroll.addEventListener("change", (e) => setScrollEnabled(e.target.checked));
}
$("iosBounce")?.addEventListener("change", (e) =>
    Events.Emit("ios:setBounceEnabled", { enabled: e.target.checked }));

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
