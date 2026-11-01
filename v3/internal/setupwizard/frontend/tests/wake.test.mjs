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
const { toggleWakeTarget, toggleWakeFormat, selectWakeBackend, initSteps, initWakeConfig } = await loadSource('wake');
const { getWizardBootstrap, checkWakeRequirements, createProject } = await loadSource('api');
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

test('init creates only after Wake choices; flag-off flow ends at details', () => {
  assert.deepEqual(initSteps(false, true), ['project', 'framework', 'language', 'details', 'Application', 'Targets', 'Packages', 'Tools']);
  assert.deepEqual(initSteps(true, false), ['project', 'framework', 'language', 'bindings', 'details']);
});

test('setup remains global even with Wake enabled and never requests project routes', async (t) => {
  const calls = [];
  mockAPI(t, (url) => {
    calls.push(url);
    if (url === '/api/init') return { body: null };
    if (url === '/api/state') return { body: { wakeEnabled: true } };
    assert.fail(`Global setup requested ${url}`);
  });
  assert.equal((await getWizardBootstrap()).init, null);
  assert.deepEqual(calls.sort(), ['/api/init', '/api/state']);
});

test('init owns Wake choices and requires server-provided options', async (t) => {
  let init = { mode: 'init', wakeEnabled: true, wakeSetup: { enabled: true, config: base() } };
  mockAPI(t, (url) => {
    if (url === '/api/init') return { body: init };
    if (url === '/api/state') return { body: { wakeEnabled: true } };
    assert.fail(`Init requested ${url}`);
  });
  assert.deepEqual((await getWizardBootstrap()).init, init);
  init = { mode: 'init', wakeEnabled: true };
  await assert.rejects(getWizardBootstrap(), /options are unavailable/);
});

test('failed mode detection cannot fall back to global setup', async (t) => {
  mockAPI(t, (url) => url === '/api/state' ? { status: 500, body: { error: 'bad request' } } : { body: null });
  await assert.rejects(getWizardBootstrap(), /Could not load wizard state/);
});

test('preview sends current metadata and only final creation submits choices', async (t) => {
  const form = { mode: 'init', projectName: 'demo', productName: 'Demo', productIdentifier: 'com.example.demo', productCompany: 'Example', productDescription: 'Demo project', productVersion: '1.0.0', productCopyright: '', productComments: '', useInterfaces: true };
  const selected = initWakeConfig(form, { ...base(), sign: true, prebuilt: true }, false);
  assert.equal(selected.prebuilt, false);
  assert.deepEqual(selected.project.bindings, { typescript: false, interfaces: false });
  assert.equal(selected.project.info.productName, 'Demo');
  const calls = [];
  mockAPI(t, (url, options) => {
    calls.push(url);
    assert.equal(options.method, 'POST');
    const sent = JSON.parse(options.body);
    if (url.endsWith('/requirements')) { assert.deepEqual(sent, selected); return { body: { requirements: [], errors: [], warnings: [], command: 'wails3 build setup' } }; }
    assert.equal(url, '/api/init/create');
    assert.deepEqual(sent.wake, selected);
    return { body: { success: true } };
  });
  await checkWakeRequirements(selected);
  assert.deepEqual(calls, ['/api/wake/setup/requirements']);
  assert.equal((await createProject({ ...form, wake: selected })).success, true);
  assert.deepEqual(calls, ['/api/wake/setup/requirements', '/api/init/create']);
});

test('rejected creation preserves backend diagnostics', async (t) => {
  mockAPI(t, () => ({ status: 400, text: 'Invalid target selection\n' }));
  await assert.rejects(createProject({ wake: base() }), /Invalid target selection/);
});

test('plain-text preview errors retain the diagnostic', async (t) => {
  mockAPI(t, () => ({ status: 400, text: 'invalid project metadata\n' }));
  await assert.rejects(checkWakeRequirements(base()), /invalid project metadata/);
});

test('new iOS selection includes the displayed device destination', () => {
  const selected = toggleWakeTarget(base(), { os: 'ios', arch: 'arm64', formats: ['ipa'] });
  assert.equal(selected.targets[1].destination, 'device');
});
