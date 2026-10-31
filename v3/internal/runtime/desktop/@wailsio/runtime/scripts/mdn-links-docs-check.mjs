import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import { Application } from "typedoc";
import config from "../typedoc.config.mjs";

function isMDNLink(value) {
    if (typeof value !== "string") return false;
    try {
        const url = new URL(value);
        return url.protocol === "https:" && url.hostname === "developer.mozilla.org" &&
            url.port === "" && url.username === "" && url.password === "";
    } catch {
        return false;
    }
}

function externalLinks(value, links, unlinked) {
    if (!value || typeof value !== "object") return;
    const add = (name, url) => {
        if (links.has(name)) assert.equal(links.get(name), url, `inconsistent link for ${name}`);
        links.set(name, url);
    };
    if (isMDNLink(value.externalUrl)) {
        add(`type:${value.name}`, value.externalUrl);
    }
    if (isMDNLink(value.target)) {
        // The built-in resolver keeps source # notation for instance members.
        add(`inline:${value.text.replaceAll("#", ".")}`, value.target);
    }
    if (value.type === "reference" && value.package === "typescript" && !value.externalUrl &&
        (typeof value.target === "object" || value.target === -1)) {
        unlinked.add(value.name);
    }
    for (const child of Object.values(value)) externalLinks(child, links, unlinked);
}

test("MDN links require the exact HTTPS origin without credentials", () => {
    const cases = [
        ["https://developer.mozilla.org/en-US/docs/Web/API/Window", true],
        ["https://DEVELOPER.MOZILLA.ORG:443/en-US/docs/Web/API/Window?q=window#methods", true],
        ["https://developer.mozilla.org.evil.example/docs", false],
        ["https://evil-developer.mozilla.org/docs", false],
        ["https://evil.example/developer.mozilla.org", false],
        ["https://evil.example/?host=developer.mozilla.org", false],
        ["https://evil.example/#developer.mozilla.org", false],
        ["https://developer.mozilla.org@evil.example/docs", false],
        ["https://user:secret@developer.mozilla.org/docs", false],
        ["https://developer.mozilla.org:8443/docs", false],
        ["http://developer.mozilla.org/docs", false],
        ["ftp://developer.mozilla.org/docs", false],
        ["javascript:developer.mozilla.org", false],
        ["//developer.mozilla.org/docs", false],
        ["not a URL: developer.mozilla.org", false],
        [null, false],
        [{ hostname: "developer.mozilla.org" }, false],
    ];
    for (const [url, accepted] of cases) {
        const links = new Map();
        externalLinks({ name: "Window", externalUrl: url, text: "Window#close", target: url },
            links, new Set());
        assert.deepEqual(Object.fromEntries(links), accepted
            ? { "type:Window": url, "inline:Window.close": url }
            : {}, `unexpected links for ${JSON.stringify(url)}`);
    }
});

test("runtime documentation retains MDN targets and exposes newly unlinked types", async () => {
    const app = await Application.bootstrapWithPlugins({
        ...config,
        entryPoints: ["./src/index.ts"],
        plugin: ["./scripts/type-reference-docs.mjs"],
        gitRevision: "master",
        logLevel: "Error",
    });
    const project = await app.convert();
    assert.ok(project);
    const model = app.serializer.projectToObject(project, process.cwd().replaceAll("\\", "/"));
    const links = new Map();
    const unlinked = new Set();
    externalLinks(model, links, unlinked);
    const expected = JSON.parse(await readFile(new URL("./testdata/mdn-links.json", import.meta.url), "utf8"));
    assert.deepEqual(Object.fromEntries([...links].sort()), expected.links);
    assert.deepEqual([...unlinked].sort(), expected.unlinkedTypes,
        "update MDN mappings when new public types need links");
});
