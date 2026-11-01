import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import test from 'node:test';
import ts from 'typescript';

async function loadSource(name) {
  const source = await readFile(new URL(`../src/${name}.ts`, import.meta.url), 'utf8');
  const { outputText } = ts.transpileModule(source, {
    compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2020 },
  });
  return import(`data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`);
}
const { toggleWakeTarget, toggleWakeFormat, selectWakeBackend, wakeSteps, wakeSigningPlatforms } = await loadSource('wake');
const { getWizardBootstrap, checkWakeRequirements, saveWakeSetup } = await loadSource('api');
const base = () => ({ applicationType: 'gui', backend: 'cgo', experimentalPureGo: false, prebuilt: false, sign: false, targets: [{ os: 'linux', arch: 'amd64', formats: ['deb'] }] });

function mockAPI(t, handler) {
  const original = globalThis.fetch;
  t.after(() => { globalThis.fetch = original; });
  globalThis.fetch = async (url, options) => {
    const response = await handler(url, options);
    return new Response(response.text ?? JSON.stringify(response.body), { status: response.status || 200, headers: { 'Content-Type': response.text === undefined ? 'application/json' : 'text/plain' } });
  };
}

test('target and format changes preserve other architectures and never enable purego', () => {
  const original = base();
  const added = toggleWakeTarget(original, { os: 'linux', arch: 'arm64', formats: ['rpm'] });
  assert.deepEqual(original.targets, [{ os: 'linux', arch: 'amd64', formats: ['deb'] }]);
  assert.equal(added.experimentalPureGo, false);
  assert.deepEqual(added.targets[1].formats, []);
  const packaged = toggleWakeFormat(added, added.targets[1], 'rpm');
  assert.deepEqual(packaged.targets.map((item) => item.formats), [['deb'], ['rpm']]);
  const removed = toggleWakeTarget(packaged, packaged.targets[0]);
  assert.deepEqual(removed.targets, [{ os: 'linux', arch: 'arm64', formats: ['rpm'] }]);
});

test('only an explicit backend choice enables purego and switching back clears it', () => {
  const experimental = selectWakeBackend(base(), 'purego');
  assert.equal(experimental.experimentalPureGo, true);
  assert.equal(experimental.backend, 'purego');
  assert.equal(selectWakeBackend(experimental, 'cgo').experimentalPureGo, false);
});

test('existing metadata and unsigned flows skip only their own setup pages', () => {
  assert.deepEqual(wakeSteps(true, false), ['Application', 'Targets', 'Packages', 'Tools', 'Save']);
  assert.deepEqual(wakeSteps(false, true), ['Project', 'Application', 'Targets', 'Packages', 'Tools', 'Signing', 'Save']);
});

test('server state selects Wake without probing blanket legacy or Docker dependencies', async (t) => {
  const calls = [];
  mockAPI(t, (url) => {
    calls.push(url);
    if (url === '/api/init') return { body: null };
    if (url === '/api/state') return { body: { wakeEnabled: true } };
    if (url === '/api/wake/setup') return { body: { enabled: true, config: base() } };
    assert.fail(`Unexpected request ${url}`);
  });
  assert.equal((await getWizardBootstrap()).wake.enabled, true);
  assert.deepEqual(calls.sort(), ['/api/init', '/api/state', '/api/wake/setup']);
});

test('flag-off and init mode preserve legacy routing', async (t) => {
  let init = null, enabled = false;
  mockAPI(t, (url) => {
    if (url === '/api/init') return { body: init };
    if (url === '/api/state') return { body: { wakeEnabled: enabled } };
    assert.fail(`Legacy/init flow requested ${url}`);
  });
  assert.equal((await getWizardBootstrap()).wake, null);
  init = { mode: 'init' }; enabled = true;
  const loaded = await getWizardBootstrap();
  assert.equal(loaded.wake, null);
  assert.deepEqual(loaded.init, init);
});

test('failed mode detection cannot fall back to the legacy flow', async (t) => {
  mockAPI(t, (url) => url === '/api/state' ? { status: 500, body: { error: 'bad project' } } : { body: null });
  await assert.rejects(getWizardBootstrap(), /Could not load wizard state/);
});

test('requirements and atomic save send exact choices including one-time prebuilt and signing', async (t) => {
  const selected = { ...base(), applicationType: 'server', prebuilt: true, sign: true };
  const calls = [];
  mockAPI(t, (url, options) => {
    assert.equal(options.method, 'POST');
    assert.deepEqual(JSON.parse(options.body), selected);
    calls.push(url);
    return { body: url.endsWith('/requirements') ? { requirements: [], errors: [], warnings: [], command: 'package --skip-build' } : { path: '/project/wails.yaml', config: selected, command: 'package --skip-build' } };
  });
  assert.equal((await checkWakeRequirements(selected)).command, 'package --skip-build');
  assert.equal((await saveWakeSetup(selected)).config.sign, true);
  assert.deepEqual(calls, ['/api/wake/setup/requirements', '/api/wake/setup']);
});

test('rejected save surfaces the server error without claiming success', async (t) => {
  mockAPI(t, () => ({ status: 400, body: { error: 'Invalid target selection' } }));
  await assert.rejects(saveWakeSetup(base()), /Invalid target selection/);
});

test('signing tabs follow selected targets, retaining order and removing duplicate architectures', () => {
  assert.deepEqual(wakeSigningPlatforms([{ os: 'windows', arch: 'amd64' }, { os: 'windows', arch: 'arm64' }, { os: 'linux', arch: 'amd64' }]), ['windows', 'linux']);
});

test('plain-text backend validation errors retain the diagnostic', async (t) => {
  mockAPI(t, () => ({ status: 400, text: 'invalid project metadata\n' }));
  await assert.rejects(saveWakeSetup(base()), /invalid project metadata/);
});

test('selecting iOS sends the device destination displayed by the initial control', () => {
  const selected = toggleWakeTarget(base(), { os: 'ios', arch: 'arm64', formats: ['ipa'] });
  assert.equal(selected.targets[1].destination, 'device');
});
