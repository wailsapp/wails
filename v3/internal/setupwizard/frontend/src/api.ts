import type { WizardState, DependencyStatus, DockerStatus, UserConfig, GlobalDefaults, SigningDefaults, SigningStatus } from './types';

const API_BASE = '/api';

export async function getState(): Promise<WizardState> {
  const response = await fetch(`${API_BASE}/state`);
  if (!response.ok) throw new Error(`Could not load wizard state (HTTP ${response.status})`);
  return response.json();
}

export async function checkDependencies(): Promise<DependencyStatus[]> {
  const response = await fetch(`${API_BASE}/dependencies/check`);
  return response.json();
}

export async function checkMobileDependencies(ios: boolean, android: boolean): Promise<DependencyStatus[]> {
  const response = await fetch(`${API_BASE}/dependencies/mobile?ios=${ios}&android=${android}`);
  return response.json();
}

export async function getDockerStatus(): Promise<DockerStatus> {
  const response = await fetch(`${API_BASE}/docker/status`);
  return response.json();
}

export function subscribeDockerStatus(onUpdate: (status: DockerStatus) => void): () => void {
  let eventSource: EventSource | null = null;
  let closed = false;

  const connect = () => {
    if (closed) return;
    
    eventSource = new EventSource(`${API_BASE}/docker/status/stream`);
    
    eventSource.onmessage = (event) => {
      try {
        const status = JSON.parse(event.data) as DockerStatus;
        onUpdate(status);
      } catch (e) {
        console.error('Failed to parse docker status:', e);
      }
    };

    eventSource.onerror = () => {
      eventSource?.close();
      if (!closed) {
        setTimeout(connect, 1000);
      }
    };
  };

  connect();

  return () => {
    closed = true;
    eventSource?.close();
  };
}

export async function buildDockerImage(): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/docker/build`, { method: 'POST' });
  return response.json();
}

export interface DockerStartBackgroundResponse {
  started: boolean;
  reason?: string;
  status: DockerStatus;
}

export async function startDockerBuildBackground(): Promise<DockerStartBackgroundResponse> {
  const response = await fetch(`${API_BASE}/docker/start-background`, { method: 'POST' });
  return response.json();
}

export async function detectConfig(): Promise<Partial<UserConfig>> {
  const response = await fetch(`${API_BASE}/config/detect`);
  return response.json();
}

export async function saveConfig(config: UserConfig): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/config/save`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config),
  });
  return response.json();
}

export async function complete(): Promise<{ status: string; duration: string; dockerBuilding: boolean }> {
  const response = await fetch(`${API_BASE}/complete`, { method: 'POST' });
  if (!response.ok) throw new Error(`Could not finish setup (HTTP ${response.status})`);
  return response.json();
}

export interface CloseResponse {
  status: string;
  dockerBuilding: boolean;
  message?: string;
}

export async function close(): Promise<CloseResponse> {
  const response = await fetch(`${API_BASE}/close`);
  return response.json();
}

export interface InstallResult {
  success: boolean;
  output: string;
  error?: string;
}

export async function installDependency(command: string): Promise<InstallResult> {
  const response = await fetch(`${API_BASE}/dependencies/install`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ command }),
  });
  return response.json();
}

export async function getDefaults(): Promise<GlobalDefaults> {
  const response = await fetch(`${API_BASE}/defaults`);
  return response.json();
}

export async function saveDefaults(defaults: GlobalDefaults): Promise<{ status: string; path: string }> {
  const response = await fetch(`${API_BASE}/defaults`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(defaults),
  });
  return response.json();
}

export async function getSigningStatus(): Promise<SigningStatus> {
  const response = await fetch(`${API_BASE}/signing/status`);
  return response.json();
}

export async function getSigning(): Promise<SigningDefaults | null> {
  const response = await fetch(`${API_BASE}/signing`);
  return response.json();
}

export async function saveSigning(signing: SigningDefaults): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/signing`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(signing),
  });
  return response.json();
}

export async function validateNotarizationProfile(profile: string): Promise<{ valid: boolean; error?: string }> {
  const response = await fetch(`${API_BASE}/signing/notarize/validate?profile=${encodeURIComponent(profile)}`);
  return response.json();
}

// createNotarizationProfile opens a Terminal window that runs `xcrun notarytool
// store-credentials`. The app-specific password is typed into that window, so it
// never reaches the browser; poll getNotarizationStatus until the window is done.
export async function createNotarizationProfile(data: {
  profileName: string;
  appleID: string;
  teamID: string;
}): Promise<{ success: boolean; pending?: boolean; command?: string; error?: string }> {
  const response = await fetch(`${API_BASE}/signing/notarize/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  return response.json();
}

export async function getNotarizationStatus(): Promise<{
  state: 'idle' | 'running' | 'succeeded' | 'failed';
  command?: string;
  error?: string;
}> {
  const response = await fetch(`${API_BASE}/signing/notarize/status`);
  return response.json();
}

export async function cancelNotarization(): Promise<{ success: boolean }> {
  const response = await fetch(`${API_BASE}/signing/notarize/cancel`, { method: 'POST' });
  return response.json();
}

export async function getInit(): Promise<import('./types').InitData | null> {
  const response = await fetch(`${API_BASE}/init`);
  if (!response.ok) throw new Error(`Could not detect wizard mode (HTTP ${response.status})`);
  return response.json();
}

export async function createProject(data: import('./types').InitData): Promise<{ success: boolean; error?: string }> {
  const response = await fetch(`${API_BASE}/init/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  const body = await response.text();
  let result;
  try { result = JSON.parse(body); }
  catch { throw new Error(body.trim() || `Project creation failed (HTTP ${response.status})`); }
  if (!response.ok) throw new Error(result.error || `Project creation failed (HTTP ${response.status})`);
  return result;
}

export async function createGPGKey(data: {
  name: string;
  email: string;
  passphrase: string;
}): Promise<{ success: boolean; keyID?: string; keyPath?: string; error?: string }> {
  const response = await fetch(`${API_BASE}/signing/gpg/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  return response.json();
}

export async function exportGPGKey(data: {
  keyID?: string;
  passphrase?: string;
}): Promise<{ success: boolean; keyID?: string; keyPath?: string; needsPassphrase?: boolean; error?: string }> {
  const response = await fetch(`${API_BASE}/signing/gpg/export`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  return response.json();
}

export async function createWindowsCert(data: {
  commonName: string;
  password: string;
}): Promise<{ success: boolean; path?: string; selfSign?: boolean; error?: string }> {
  const response = await fetch(`${API_BASE}/signing/windows/create-cert`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  return response.json();
}

export async function reportBug(currentStep: string): Promise<{ status: string; body?: string; url?: string }> {
  const response = await fetch(`${API_BASE}/report-bug?step=${encodeURIComponent(currentStep)}`);
  return response.json();
}

async function wakeRequest<T>(path: string, config?: import('./types').WakeSetupConfig): Promise<T> {
  const response = await fetch(`${API_BASE}/wake/setup${path}`, config === undefined ? undefined : {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config),
  });
  const body = await response.text();
  let result;
  try { result = JSON.parse(body); }
  catch {
    throw new Error(!response.ok && body.trim() ? body.trim() : `Unexpected Wake setup response (HTTP ${response.status})`);
  }
  if (!response.ok) throw new Error(typeof result?.error === 'string' ? result.error : typeof result === 'string' ? result : `Wake setup request failed (HTTP ${response.status})`);
  return result;
}

export function checkWakeRequirements(config: import('./types').WakeSetupConfig): Promise<import('./types').WakeRequirements> {
  return wakeRequest('/requirements', config);
}

export async function getWizardBootstrap() {
  const [init, state] = await Promise.all([getInit(), getState()]);
  if (init?.wakeEnabled && !init.wakeSetup?.enabled) throw new Error('Wake project options are unavailable for this init session.');
  return { init, state };
}
