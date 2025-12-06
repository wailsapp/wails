import type { WizardState, DependencyStatus, DockerStatus, UserConfig, WailsConfig } from './types';

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

export async function buildDockerImage(): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/docker/build`, { method: 'POST' });
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

export async function close(): Promise<void> {
  await fetch(`${API_BASE}/close`);
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
