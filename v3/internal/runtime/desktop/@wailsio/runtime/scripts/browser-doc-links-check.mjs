import assert from "node:assert/strict";
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import { Application } from "typedoc";
import config from "../typedoc.config.mjs";

test("runtime browser and JavaScript types resolve through native TypeDoc mappings", async t => {
    const directory = await mkdtemp(path.join(os.tmpdir(), "wails-doc-links-"));
    t.after(() => rm(directory, { recursive: true, force: true }));
    const entry = path.join(directory, "index.ts");
    const tsconfig = path.join(directory, "tsconfig.json");
    await writeFile(tsconfig, JSON.stringify({
        compilerOptions: { target: "ESNext", types: [], strict: true },
        files: ["index.ts"],
    }));
    await writeFile(entry, `
/** See {@link AbortSignal#reason}, {@link setTimeout}, and {@link !Promise}. */
export declare function wait(signal: AbortSignal): Promise<Blob>;
`);
    const app = await Application.bootstrap({
        ...config, entryPoints: [entry], tsconfig, logLevel: "Error",
    });
    const project = await app.convert();
    assert.ok(project);
    const wait = project.children.find(child => child.name === "wait");
    const signature = wait.signatures[0];
    assert.equal(signature.parameters[0].type.externalUrl,
        "https://developer.mozilla.org/docs/Web/API/AbortSignal");
    assert.equal(signature.type.externalUrl,
        "https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Promise");
    assert.equal(signature.type.typeArguments[0].externalUrl,
        "https://developer.mozilla.org/docs/Web/API/Blob");
    const links = signature.comment.summary.filter(part => part.kind === "inline-tag");
    assert.deepEqual(links.map(part => part.target), [
        "https://developer.mozilla.org/docs/Web/API/AbortSignal/reason",
        "https://developer.mozilla.org/docs/Web/API/Window/setTimeout",
        "https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/Promise",
    ]);
    assert.equal(links[0].text, "AbortSignal#reason");
});
