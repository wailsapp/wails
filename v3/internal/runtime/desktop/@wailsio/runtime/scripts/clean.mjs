import { rmSync } from "node:fs";

for (const output of ["dist", "docs", "types", "tsconfig.tsbuildinfo"]) {
    rmSync(new URL(`../${output}`, import.meta.url), {
        recursive: true,
        force: true,
        maxRetries: 3,
    });
}
