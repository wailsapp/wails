const host = window.wailsOpenHarmony;
window.wails = { invoke: message => host.invoke(message) };
const runtime = await import('/wails/runtime.js');
const pending = new Map();
const session = crypto.randomUUID();
let nextID = 0;
const params = new URLSearchParams(location.search);
const probe = params.get('probe');
const phase = params.get('phase');
let reuseID = false;
let probeTicks = 0;
let lastTick = -1;
let badTick = false;
window.wailsOpenHarmonyResolve = (id, payload) => {
    const reply = pending.get(id);
    if (!reply) return;
    pending.delete(id);
    clearTimeout(reply.timer);
    try {
        const envelope = JSON.parse(payload);
        if (!envelope.ok) throw new Error(envelope.error);
        reply.resolve(envelope.text ?? envelope.data);
    } catch (error) {
        reply.reject(error);
    }
};
runtime.setTransport({
    call(object, method, windowName, args) {
        const id = reuseID ? 'lifecycle-reused-id' : String(++nextID);
        const payload = JSON.stringify({ object, method, windowName, args, clientId: runtime.clientId });
        return new Promise((resolve, reject) => {
            const timer = setTimeout(() => {
                pending.delete(id);
                reject(new Error('Go call timed out'));
            }, reuseID ? 15000 : 5000);
            pending.set(id, { resolve, reject, timer });
            if (!host.request(session, id, payload)) {
                pending.delete(id);
                clearTimeout(timer);
                reject(new Error('Go is not ready or the request queue is full'));
            }
        });
    }
});

const result = document.querySelector('#result');
const event = document.querySelector('#event');
let greeted = false;
let ticked = false;
let reported = false;
function report() {
    const styled = getComputedStyle(document.body).color === 'rgb(20, 34, 52)';
    if (greeted && ticked && styled && !reported) {
        reported = true;
        host.report('M3 PASS: local HTML/JS/CSS, Greet and Go event');
    }
}
runtime.Events.On('hello:tick', ({ data }) => {
    const value = Array.isArray(data) ? data[0] : data;
    event.textContent = `Go event: ${value}`;
    if (Number.isInteger(value)) {
        badTick ||= value <= lastTick;
        lastTick = value;
        ++probeTicks;
    }
    ticked ||= Number.isInteger(value);
    report();
});
async function greet() {
    try {
        const message = await runtime.Call.ByName('main.GreetService.Greet', 'OpenHarmony');
        result.textContent = message;
        greeted = message === 'Hello OpenHarmony!';
        report();
    } catch (error) {
        result.textContent = String(error);
        host.report(`Greet failed: ${error}`);
    }
}
document.querySelector('#greet').addEventListener('click', greet);
Object.defineProperty(window, 'wailsOpenHarmonySession', { value: session });
if (!host.ready(session)) throw new Error('ArkWeb document session rejected');
await greet();

if (probe === 'api') {
    try {
        const zoomBefore = await runtime.Window.GetZoom();
        await runtime.Window.SetTitle('M3 NATIVE CHECK');
        await runtime.Window.SetBackgroundColour(20, 34, 52, 255);
        await runtime.Window.ZoomIn();
        const zoomAfter = await runtime.Window.GetZoom();
        await runtime.Window.ZoomReset();
        const appInfo = await runtime.Call.ByName('main.NativeProbe.AppInfo');
        const accent = await runtime.Call.ByName('main.NativeProbe.AccentColor');
        const secure = await runtime.Call.ByName('main.NativeProbe.SecureRoundTrip');
        if (!(zoomAfter > zoomBefore)) {
            throw new Error(`zoom ${zoomBefore} -> ${zoomAfter}`);
        }
        if (!accent || !appInfo || secure !== 'm3-native-value') {
            throw new Error(`accent=${accent}, appInfo=${appInfo}, secure=${secure}`);
        }
        host.report(`M3 API PASS zoom=${zoomBefore}->${zoomAfter} accent=${accent} secure=${secure}`);
    } catch (error) {
        host.report(`M3 API FAIL: ${error}`);
    }
}

if (probe === 'refresh' || probe === 'rebuild' || probe === 'stage') {
    host.report(`M4 ${probe} ${phase} session=${session} id=lifecycle-reused-id`);
    reuseID = true;
    const reply = runtime.Call.ByName('main.GreetService.LifecycleProbe', phase);
    reuseID = false;
    if (phase === 'old') {
        reply.catch(() => {});
        if (probe !== 'stage') {
            host.advanceLifecycleProbe();
        }
    } else {
        try {
            const value = await reply;
            if (value !== 'new' || !greeted || badTick || probeTicks < 2) {
                throw new Error(`reply=${value}, greeted=${greeted}, ticks=${probeTicks}, badTick=${badTick}`);
            }
            host.report(`M4 PASS ${probe} session=${session}`);
        } catch (error) {
            host.report(`M4 FAIL ${probe}: ${error}`);
        }
    }
}
