import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import test from 'node:test';
import ts from 'typescript';

const source = await readFile(new URL('../src/signing.ts', import.meta.url), 'utf8');
const { outputText } = ts.transpileModule(source, {
  compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2020 },
});
const { selectSigningIdentity } = await import(`data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`);

test('switching certificates updates the derived team and preserves other settings', () => {
  const next = selectSigningIdentity({
    identity: 'Apple Development: Example (TEAMAAAAAA)',
    teamID: 'TEAMAAAAAA',
    keychainProfile: 'profile',
  }, 'Developer ID Application: Example (TEAMBBBBBB)');
  assert.equal(next.teamID, 'TEAMBBBBBB');
  assert.equal(next.identity, 'Developer ID Application: Example (TEAMBBBBBB)');
  assert.equal(next.keychainProfile, 'profile');
});

test('a SHA-1 identity retains an explicitly supplied team', () => {
  const next = selectSigningIdentity({ identity: 'A'.repeat(40), teamID: 'MANUALTEAM' }, 'B'.repeat(40));
  assert.equal(next.teamID, 'MANUALTEAM');
});

test('an opaque identity does not inherit a team derived from another certificate', () => {
  const next = selectSigningIdentity({ identity: 'Developer ID: Example (TEAMAAAAAA)', teamID: 'TEAMAAAAAA' }, 'A'.repeat(40));
  assert.equal(next.teamID, '');
});
