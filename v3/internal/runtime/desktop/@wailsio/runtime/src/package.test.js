import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const packageJson = JSON.parse(readFileSync("package.json", "utf8"));

describe("@wailsio/runtime package exports", () => {
    it("exposes cancellable without loading the root runtime", () => {
        expect(packageJson.exports["./cancellable"]).toEqual({
            types: "./types/cancellable.d.ts",
            default: "./dist/cancellable.js",
        });
    });
});
