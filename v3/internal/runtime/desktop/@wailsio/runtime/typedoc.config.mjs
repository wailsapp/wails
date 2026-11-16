// MDN targets used by the runtime API. TypeDoc resolves these without a plugin.
const mappings = {};
const mdn = "https://developer.mozilla.org/docs/Web/";
for (const name of ["AbortSignal", "Blob", "CloseEvent", "Event", "EventTarget", "MessageEvent", "URL", "WebSocket"]) {
    mappings[name] = `${mdn}API/${name}`;
}
for (const name of ["Error", "Function", "Iterator", "Promise"]) {
    mappings[name] = `${mdn}JavaScript/Reference/Global_Objects/${name}`;
}
for (const [name, members] of Object.entries({
    AbortSignal: ["reason"],
    Error: ["cause", "message", "name", "stack", "stackTraceLimit", "captureStackTrace", "isError"],
    EventTarget: ["addEventListener", "dispatchEvent", "removeEventListener"],
    Promise: ["all", "allSettled", "any", "catch", "finally", "race", "reject", "resolve", "then", "try", "withResolvers"],
})) {
    for (const member of members) {
        mappings[`${name}.${member}`] = `${mappings[name]}/${member}`;
    }
}

mappings.setTimeout = `${mdn}API/Window/setTimeout`;

export default {
    externalSymbolLinkMappings: { global: mappings, typescript: mappings },
};
