import type { ReactNode } from 'react';
import type { WakeSetupConfig, WakeSetupState } from '../types';
import { selectWakeBackend, targetKey, toggleWakeFormat, toggleWakeTarget } from '../wake';

const inputClass = 'w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-sm text-gray-900 dark:text-white focus:outline-hidden focus:ring-2 focus:ring-red-500';

function Choice({ checked, onChange, label, children, name, disabled }: { checked: boolean; onChange: () => void; label: string; children?: ReactNode; name?: string; disabled?: boolean }) {
  return <label className={`flex items-start gap-3 rounded-lg border p-3 ${disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'} ${checked ? 'border-red-500 bg-red-50 dark:bg-red-950/20' : 'border-gray-200 dark:border-gray-700'}`}>
    <input type={name ? 'radio' : 'checkbox'} name={name} checked={checked} onChange={onChange} disabled={disabled} className="mt-1 accent-red-500" />
    <span className="text-sm font-medium text-gray-800 dark:text-gray-200">{label}{children && <span className="block mt-1 text-xs font-normal text-gray-500 dark:text-gray-400">{children}</span>}</span>
  </label>;
}

export default function WakeSetupChoices({ step, config, initial, update }: {
  step: string; config: WakeSetupConfig; initial: WakeSetupState; update: (config: WakeSetupConfig) => void;
}) {
  if (step === 'Application') {
    return <>
      <fieldset className="space-y-2"><legend className="text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Application type</legend>
        {([['gui', 'Desktop application', 'A window and web frontend.'], ['cli', 'Command-line application', 'Not available with the selected frontend template.'], ['server', 'Server', 'Not available with the selected frontend template.']] as const).map(([kind, label, detail]) => <Choice key={kind} name="kind" disabled={kind !== 'gui'} checked={config.applicationType === kind} label={label} onChange={() => update({ ...config, applicationType: kind })}>{detail}</Choice>)}
      </fieldset>
      <fieldset className="space-y-2"><legend className="text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Build backend</legend>
        <p className="text-xs text-gray-500 dark:text-gray-400">Choose the build backend for your new project.</p>
        <Choice name="backend" checked={config.backend === 'cgo'} label="Standard backend" onChange={() => update(selectWakeBackend(config, 'cgo'))}>{config.applicationType === 'gui' ? 'CGO where the platform requires it; the tool check resolves the compiler.' : 'Go builds without a GUI compiler requirement.'}</Choice>
        <Choice name="backend" checked={config.backend === 'purego'} label="Pure Go (experimental)" onChange={() => update(selectWakeBackend(config, 'purego'))}>Explicitly enable the experimental Pure Go backend. Linux GUI applications still need GTK and WebKit at runtime.</Choice>
      </fieldset>
    </>;
  }
  if (step === 'Targets') {
    return <>
      <p className="text-sm text-gray-500">Select the platforms and architectures to include in the setup profile.</p>
      <div className="grid grid-cols-2 gap-3">{initial.catalog.map((target) => <Choice key={targetKey(target)} checked={config.targets.some((selected) => targetKey(selected) === targetKey(target))} label={targetKey(target)} onChange={() => update(toggleWakeTarget(config, target))}>{target.experimental ? 'Experimental target' : undefined}</Choice>)}</div>
    </>;
  }
  if (step === 'Packages') {
    return <>
      {config.targets.map((target) => <fieldset key={targetKey(target)} className="rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <legend className="px-1 font-medium text-gray-800 dark:text-gray-200">{targetKey(target)}</legend>
        <div className="flex flex-wrap gap-x-5 gap-y-3">{initial.catalog.find((entry) => targetKey(entry) === targetKey(target))?.formats.map((format) => <label key={format} className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300"><input type="checkbox" className="accent-red-500" checked={target.formats.includes(format)} onChange={() => update(toggleWakeFormat(config, target, format))} />{format}</label>)}</div>
        {target.formats.length === 0 && <p className="mt-3 text-xs text-gray-500">Use the platform's default output without an additional package.</p>}
        {target.os === 'ios' && <label className="block mt-3 text-sm text-gray-700 dark:text-gray-300">Destination<select className={`${inputClass} mt-1`} value={target.destination || 'device'} onChange={(event) => update({ ...config, targets: config.targets.map((item) => targetKey(item) === targetKey(target) ? { ...item, destination: event.target.value as 'device' | 'simulator' } : item) })}><option value="device">Device</option><option value="simulator">Simulator</option></select></label>}
      </fieldset>)}
      <Choice checked={config.sign} onChange={() => update({ ...config, sign: !config.sign })} label="Sign release outputs">Use your configured signing credentials. Manage global credentials with wails3 setup; project-specific and mobile signing can be configured in wails.yaml.</Choice>
    </>;
  }
  return null;
}
