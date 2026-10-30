import { readFileSync, statSync } from "node:fs";
import { dirname, join, relative, resolve } from "node:path";
import { pathToFileURL } from "node:url";
import ts from "typescript";

function buildImportGraph(directory) {
    const root = resolve(directory);
    if (!statSync(root).isDirectory()) {
        throw new Error(`Expected an output directory: ${root}`);
    }
    const files = new Set(ts.sys.readDirectory(root, [".js"], ["**/node_modules/**"]).map(file => resolve(file)));
    const graph = new Map();
    for (const file of files) {
        const imports = ts.preProcessFile(readFileSync(file, "utf8"), true, true).importedFiles;
        const dependencies = new Set();
        for (const { fileName } of imports) {
            if (!fileName.startsWith(".")) continue;
            const target = resolve(dirname(file), fileName);
            const dependency = [target, `${target}.js`, join(target, "index.js")].find(candidate => files.has(candidate));
            if (dependency) dependencies.add(dependency);
        }
        graph.set(file, dependencies);
    }
    return graph;
}

// Check the generated JavaScript files; TypeScript resolution is checked by tsc.
export function findImportCycle(directory) {
    const graph = buildImportGraph(directory);
    const visited = new Set();
    const stack = [];
    const active = new Map();
    function visit(file) {
        if (active.has(file)) return [...stack.slice(active.get(file)), file];
        if (visited.has(file)) return null;
        active.set(file, stack.length);
        stack.push(file);
        for (const dependency of graph.get(file)) {
            const cycle = visit(dependency);
            if (cycle) return cycle;
        }
        stack.pop();
        active.delete(file);
        visited.add(file);
        return null;
    }
    for (const file of graph.keys()) {
        const cycle = visit(file);
        if (cycle) return cycle;
    }
    return null;
}

if (process.argv[1] && import.meta.url === pathToFileURL(resolve(process.argv[1])).href) {
    const directory = process.argv[2] ?? "output";
    const cycle = findImportCycle(directory);
    if (cycle) {
        console.error(`Circular JavaScript imports: ${cycle.map(file => relative(directory, file)).join(" -> ")}`);
        process.exitCode = 1;
    } else {
        console.log("No circular JavaScript imports found.");
    }
}
