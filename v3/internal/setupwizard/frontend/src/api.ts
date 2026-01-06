import type { WizardState, DependencyStatus, DockerStatus, UserConfig, WailsConfig, GlobalDefaults } from './types';

const API_BASE = '/api';

export async function getState(): Promise<WizardState> {
  const response = await fetch(`${API_BASE}/state`);
  return response.json();
}

export async function checkDependencies(): Promise<DependencyStatus[]> {
  const response = await fetch(`${API_BASE}/dependencies/check`);
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

export async function complete(): Promise<{ status: string; duration: string }> {
  const response = await fetch(`${API_BASE}/complete`);
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

export async function getWailsConfig(): Promise<WailsConfig | null> {
  const response = await fetch(`${API_BASE}/wails-config`);
  return response.json();
}

export async function saveWailsConfig(config: WailsConfig): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/wails-config`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config),
  });
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
