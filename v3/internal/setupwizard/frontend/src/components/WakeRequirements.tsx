import type { DependencyStatus, WakeRequirements as Requirements } from '../types';

export function RequirementMessages({ requirements }: { requirements: Requirements }) {
  return <div className="space-y-2" aria-live="polite">
    {requirements.errors.map((message) => <p key={message} role="alert" className="text-sm text-red-700 dark:text-red-300">{message}</p>)}
    {requirements.warnings.map((message) => <p key={message} className="text-sm text-amber-700 dark:text-amber-300">{message}</p>)}
  </div>;
}

export default function WakeRequirements({ requirements, busy, onInstall }: {
  requirements: Requirements; busy: boolean; onInstall?: (dependency: DependencyStatus) => void;
}) {
  return <>
    <RequirementMessages requirements={requirements} />
    <div className="space-y-3">{requirements.requirements.map((dependency) => <div key={dependency.name} className="rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div className="flex items-center justify-between gap-4">
        <div><p className="font-medium text-gray-800 dark:text-gray-200">{dependency.name}</p><p className={`text-xs mt-1 ${dependency.status === 'installed' ? 'text-green-700 dark:text-green-400' : 'text-amber-700 dark:text-amber-300'}`}>{dependency.status === 'installed' ? dependency.version || 'Installed' : dependency.status === 'needs_config' ? 'Needs configuration' : dependency.status === 'needs_update' ? 'Update needed' : 'Not ready'}</p></div>
        {onInstall && dependency.status !== 'installed' && dependency.installCommand && <button className="rounded-lg border border-gray-300 dark:border-gray-600 px-3 py-2 text-sm text-gray-700 dark:text-gray-300 disabled:opacity-50" disabled={busy} onClick={() => onInstall(dependency)}>{dependency.status === 'needs_update' ? 'Update' : 'Install'}</button>}
      </div>
      {dependency.message && <p className="mt-2 text-sm text-gray-500">{dependency.message}</p>}
      {dependency.configCommand && <pre className="mt-2 p-2 rounded bg-gray-100 dark:bg-gray-950 text-xs text-gray-700 dark:text-gray-300 overflow-x-auto"><code>{dependency.configCommand}</code></pre>}
      {dependency.helpUrl && <a className="inline-block mt-2 text-sm underline text-red-600 dark:text-red-400" href={dependency.helpUrl} target="_blank" rel="noopener noreferrer">{dependency.helpLabel || 'Setup instructions'}</a>}
    </div>)}</div>
  </>;
}
