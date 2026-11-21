const host = window.wailsOpenHarmony;
window.wails = { invoke: message => host.invoke(message) };
const runtime = await import('/wails/runtime.js');
const pending = new Map();
const session = crypto.randomUUID();
let nextID = 0;
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
        const id = `${session}:${++nextID}`;
        const payload = JSON.stringify({ object, method, windowName, args, clientId: runtime.clientId });
        return new Promise((resolve, reject) => {
            const timer = setTimeout(() => {
                pending.delete(id);
                reject(new Error('Go call timed out'));
            }, 5000);
            pending.set(id, { resolve, reject, timer });
            if (!host.request(id, payload)) {
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
await greet();
