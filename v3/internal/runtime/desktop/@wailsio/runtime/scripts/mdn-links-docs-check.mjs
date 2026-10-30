import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import { Application } from "typedoc";
import config from "../typedoc.config.mjs";

function externalLinks(value, links, unlinked) {
    if (!value || typeof value !== "object") return;
    const add = (name, url) => {
        if (links.has(name)) assert.equal(links.get(name), url, `inconsistent link for ${name}`);
        links.set(name, url);
    };
    if (value.externalUrl?.includes("developer.mozilla.org")) {
        add(`type:${value.name}`, value.externalUrl);
    }
    if (typeof value.target === "string" && value.target.includes("developer.mozilla.org")) {
        // The built-in resolver keeps source # notation for instance members.
        add(`inline:${value.text.replaceAll("#", ".")}`, value.target);
    }
    if (value.type === "reference" && value.package === "typescript" && !value.externalUrl &&
        (typeof value.target === "object" || value.target === -1)) {
        unlinked.add(value.name);
    }
    for (const child of Object.values(value)) externalLinks(child, links, unlinked);
}

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
