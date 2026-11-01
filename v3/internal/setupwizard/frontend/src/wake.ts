import type { InitData, WakeSetupConfig, WakeTarget } from './types';

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

export function initSteps(typescript: boolean, wakeEnabled: boolean): string[] {
  return ['project', 'framework', 'language', ...(typescript ? ['bindings'] : []), 'details', ...(wakeEnabled ? ['Application', 'Targets', 'Packages', 'Tools'] : [])];
}

export function initWakeConfig(form: InitData, config: WakeSetupConfig, typescript: boolean): WakeSetupConfig {
  return { ...config, prebuilt: false, project: {
    info: { projectName: form.projectName, productName: form.productName, productIdentifier: form.productIdentifier, companyName: form.productCompany, description: form.productDescription, version: form.productVersion, copyright: form.productCopyright, comments: form.productComments },
    bindings: { typescript, interfaces: typescript && form.useInterfaces },
  } };
}
