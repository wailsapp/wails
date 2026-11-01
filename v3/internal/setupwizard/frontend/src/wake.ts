import type { WakeSetupConfig, WakeTarget } from './types';

export function targetKey(target: Pick<WakeTarget, 'os' | 'arch'>): string {
  return `${target.os}/${target.arch}`;
}

export function toggleWakeTarget(config: WakeSetupConfig, target: WakeTarget): WakeSetupConfig {
  const key = targetKey(target);
  const selected = config.targets.some((item) => targetKey(item) === key);
  return {
    ...config,
    targets: selected
      ? config.targets.filter((item) => targetKey(item) !== key)
      : [...config.targets, { os: target.os, arch: target.arch, formats: [], ...(target.os === 'ios' ? { destination: 'device' as const } : {}) }],
  };
}

export function toggleWakeFormat(config: WakeSetupConfig, target: WakeTarget, format: string): WakeSetupConfig {
  return {
    ...config,
    targets: config.targets.map((item) => targetKey(item) === targetKey(target)
      ? { ...item, formats: item.formats.includes(format) ? item.formats.filter((value) => value !== format) : [...item.formats, format] }
      : item),
  };
}

export function selectWakeBackend(config: WakeSetupConfig, backend: WakeSetupConfig['backend']): WakeSetupConfig {
  return { ...config, backend, experimentalPureGo: backend === 'purego' };
}

export function wakeSteps(hasManifest: boolean, sign: boolean): string[] {
  return [...(hasManifest ? [] : ['Project']), 'Application', 'Targets', 'Packages', 'Tools', ...(sign ? ['Signing'] : []), 'Save'];
}

export function wakeSigningPlatforms(targets: WakeTarget[]): ('darwin' | 'windows' | 'linux')[] {
  const platforms = targets.map((target) => target.os).filter((os): os is 'darwin' | 'windows' | 'linux' => os === 'darwin' || os === 'windows' || os === 'linux');
  return [...new Set(platforms)];
}
