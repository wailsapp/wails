import ts from 'typescript';
import { readdirSync } from 'node:fs';
import { join } from 'node:path';

// Each output configuration is a separate application. Loading them into one
// program merges incompatible global event declarations (classes vs interfaces).
const configFile = ts.readConfigFile('tsconfig.json', ts.sys.readFile);
const host = { getCanonicalFileName: p => p, getCurrentDirectory: ts.sys.getCurrentDirectory, getNewLine: () => '\n' };
if (configFile.error) {
    console.error(ts.formatDiagnosticsWithColorAndContext([configFile.error], host));
    process.exit(1);
}
let failed = false;
for (const lang of readdirSync('output')) {
    for (const interfaces of readdirSync(join('output', lang))) {
        for (const names of readdirSync(join('output', lang, interfaces))) {
            const root = join('output', lang, interfaces, names).replaceAll('\\', '/');
            const config = ts.parseJsonConfigFileContent({
                ...configFile.config,
                include: [`${root}/**/*.js`, `${root}/**/*.ts`],
            }, ts.sys, process.cwd());
            const program = ts.createProgram({ rootNames: config.fileNames, options: config.options, projectReferences: config.projectReferences });
            const diagnostics = [...config.errors, ...ts.getPreEmitDiagnostics(program)];
            if (diagnostics.length) {
                console.error(ts.formatDiagnosticsWithColorAndContext(diagnostics, host));
                failed = true;
            }
        }
    }
}
process.exitCode = failed ? 1 : 0;
