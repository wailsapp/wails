import assert from "node:assert/strict";
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import { Application, ReflectionKind } from "typedoc";
import { load } from "./type-reference-docs.mjs";

test("helper documentation stays scoped to each conversion and application", async t => {
    const directory = await mkdtemp(path.join(os.tmpdir(), "wails-type-docs-"));
    t.after(() => rm(directory, { recursive: true, force: true }));
    const entry = path.join(directory, "index.ts");
    const config = path.join(directory, "tsconfig.json");
    await writeFile(config, JSON.stringify({
        compilerOptions: { target: "ESNext", types: [], strict: true },
        files: ["index.ts"],
    }));
    await writeFile(entry, `
interface Nested { value: string }
interface Hidden { nested: Nested; timestamp: Date }
class PrivateFields {
    private secret = 1;
    public visible = 2;
}
/** A public operation using private helper declarations. */
export declare function run(value: Record<string, Hidden>): PrivateFields;
`);
    const options = { entryPoints: [entry], tsconfig: config, logLevel: "Error" };
    const app = await Application.bootstrap(options);
    load(app);
    load(app);
    const project = await app.convert();
    assert.ok(project);
    const helpers = project.children.find(child => child.name === "<internal>");
    assert.ok(helpers);
    assert.deepEqual(helpers.children.map(child => child.name).sort(), [
        "Hidden", "Nested", "PrivateFields", "Record",
    ]);
    const fields = helpers.children.find(child => child.name === "PrivateFields");
    assert.equal(fields.kind, ReflectionKind.Class);
    assert.ok(fields.children.some(child => child.name === "visible"));
    assert.ok(!fields.children.some(child => child.name === "secret"));
    const signature = project.children.find(child => child.name === "run").signatures[0];
    assert.equal(signature.type.reflection, fields);
    const hidden = helpers.children.find(child => child.name === "Hidden");
    assert.equal(signature.parameters[0].type.typeArguments[1].reflection, hidden);

    const other = await Application.bootstrap(options);
    const plain = await other.convert();
    assert.ok(plain);
    assert.ok(!plain.children.some(child => child.name === "<internal>"));

    await writeFile(entry, `
interface Replacement { value: number }
export declare function next(): Replacement;
`);
    const repeated = await app.convert();
    assert.ok(repeated);
    const nextHelpers = repeated.children.find(child => child.name === "<internal>");
    assert.deepEqual(nextHelpers.children.map(child => child.name), ["Replacement"]);
});
