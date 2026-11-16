import assert from "node:assert/strict";
import { mkdtempSync, mkdirSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { test } from "node:test";
import { findImportCycle } from "./check-cycles.mjs";

function fixture(t, files) {
    const root = mkdtempSync(join(tmpdir(), "wails import check "));
    t.after(() => rmSync(root, { recursive: true, force: true }));
    for (const [name, content] of Object.entries(files)) {
        const file = join(root, name);
        mkdirSync(dirname(file), { recursive: true });
        writeFileSync(file, content);
    }
    return root;
}

for (const [name, files, cyclic] of [
    ["static imports", { "a.js": 'import "./b.js"', "b.js": 'import "./a.js"' }, true],
    ["reexports", { "a.js": 'export * from "./b.js"', "b.js": 'export { x } from "./a.js"' }, true],
    ["dynamic imports", { "a.js": 'import("./b.js")', "b.js": 'import("./a.js")' }, true],
    ["commonjs", { "a.js": 'require("./b.js")', "b.js": 'require("./a.js")' }, true],
    ["self import", { "a.js": 'import "./a.js"' }, true],
    ["directory index", { "a.js": 'import "./folder"', "folder/index.js": 'import "../a"' }, true],
    ["diamond", { "a.js": 'import "./b.js"; import "./c.js"', "b.js": 'import "./d.js"', "c.js": 'import "./d.js"', "d.js": "export {}" }, false],
    ["source text", { "a.js": 'const text = \'import "./b.js"\'; // import "./b.js"', "b.js": 'import "./a.js"' }, false],
    ["external runtime", { "a.js": 'import "/wails/runtime.js"; import "external"; import "./missing.js"' }, false],
    ["ignore node modules", { "a.js": "export {}", "node_modules/module/a.js": 'import "./a.js"' }, false],
    ["JavaScript scope", { "a.ts": 'import "./b.js"', "b.ts": 'import "./a.js"' }, false],
    ["empty directory", {}, false],
]) {
    test(name, t => {
        const root = fixture(t, files);
        const cycle = findImportCycle(root);
        assert.equal(cycle !== null, cyclic);
        if (cycle) assert.equal(cycle[0], cycle.at(-1));
    });
}

test("CLI status and diagnostic identify cycles", t => {
    const root = fixture(t, { "a.js": 'import "./b.js"', "b.js": 'import "./a.js"' });
    const result = spawnSync(process.execPath, [fileURLToPath(new URL("./check-cycles.mjs", import.meta.url)), root], { encoding: "utf8" });
    assert.equal(result.status, 1);
    assert.match(result.stderr, /Circular JavaScript imports:/);
    assert.match(result.stderr, /a\.js/);
    assert.match(result.stderr, /b\.js/);
});

test("missing input fails", t => {
    const root = fixture(t, {});
    assert.throws(() => findImportCycle(join(root, "missing")), /ENOENT/);
});
