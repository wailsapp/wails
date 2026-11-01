import { useState, type ReactNode } from 'react';
import { checkWakeRequirements, complete, installDependency, saveWakeSetup } from '../api';
import type { DependencyStatus, WakeRequirements as Requirements, WakeSetupConfig, WakeSetupResult, WakeSetupState } from '../types';
import { targetKey, wakeSteps, wakeSigningPlatforms } from '../wake';
import SigningStep from './SigningStep';
import WakeSetupChoices from './WakeSetupChoices';
import WakeRequirements, { RequirementMessages } from './WakeRequirements';

interface Props {
  initial: WakeSetupState;
  sidebar: (progress: { labels: string[]; index: number }) => ReactNode;
}
const buttonClass = 'rounded-lg border border-gray-300 dark:border-gray-600 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed';
const nextClass = 'rounded-lg border border-red-500 px-5 py-2 text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-500/10 disabled:opacity-50 disabled:cursor-not-allowed';
const titles: Record<string, string> = { Project: 'About your project', Application: 'What are you building?', Targets: 'Where should it run?', Packages: 'Choose release outputs', Tools: 'Tools for your build', Signing: 'Configure mobile signing', Save: 'Save your build configuration' };

function Layout({ sidebar, children }: { sidebar: ReactNode; children: ReactNode }) {
  return <div className="fixed inset-0 overflow-auto bg-gray-50 dark:bg-[#0f0f0f]">
    <div className="fixed inset-0 bg-center bg-cover bg-no-repeat pointer-events-none opacity-50 dark:opacity-40" style={{ backgroundImage: "url('/digital_wales_master.webp')" }} />
    <div className="relative min-h-full min-w-full w-fit flex items-center justify-center p-4">
      <div className="glass-card rounded-2xl flex overflow-hidden relative z-10" style={{ aspectRatio: '3 / 2', width: 'clamp(54rem, min(100vw - 2rem, (100vh - 2rem) * 1.5), 75rem)' }}>
        {sidebar}<div className="flex-1 flex flex-col min-w-[36rem] bg-white/85 dark:bg-[#0a0e16]/85 backdrop-blur-[30px] backdrop-saturate-150 relative">{children}</div>
      </div>
    </div>
  </div>;
}

export default function WakeSetupFlow(props: Props) {
  if (props.initial.error) {
    return <Layout sidebar={props.sidebar({ labels: ['Project'], index: 0 })}>
      <main className="m-auto max-w-lg p-8 space-y-4">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Unable to load this project</h1>
        <p role="alert" className="text-red-600 dark:text-red-400">{props.initial.error}</p>
        <p className="text-sm text-gray-500 break-all">{props.initial.projectRoot}</p>
        <button className={buttonClass} onClick={() => window.location.reload()}>Retry</button>
      </main>
    </Layout>;
  }
  return <WakeSetupForm {...props} />;
}

function WakeSetupForm({ initial, sidebar }: Props) {
  const [config, setConfig] = useState(initial.config);
  const [step, setStep] = useState(initial.manifestPath ? 'Application' : 'Project');
  const [requirements, setRequirements] = useState<Requirements | null>(null);
  const [result, setResult] = useState<WakeSetupResult | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [installOutput, setInstallOutput] = useState('');
  const [finished, setFinished] = useState(false);
  const labels = wakeSteps(Boolean(initial.manifestPath), config.sign);
  const signingPlatforms = wakeSigningPlatforms(config.targets);
  const mobileSigning = config.targets.some((target) => target.os === 'ios' || target.os === 'android');
  const index = result ? labels.length : labels.indexOf(step);
  const update = (next: WakeSetupConfig) => { setConfig(next); setRequirements(null); setError(''); };
  const reportError = (cause: unknown) => setError(cause instanceof Error ? cause.message : 'The setup request failed. Please try again.');

  const check = async (draft = config) => {
    setBusy(true); setError('');
    try { setRequirements(await checkWakeRequirements(draft)); }
    catch (cause) { reportError(cause); }
    finally { setBusy(false); }
  };
  const next = async () => {
    const destination = labels[index + 1];
    if (!destination) return;
    setStep(destination);
    if (destination === 'Tools' || destination === 'Save') await check();
  };
  const install = async (dependency: DependencyStatus) => {
    if (!dependency.installCommand) return;
    setBusy(true); setError(''); setInstallOutput('');
    try {
      const installed = await installDependency(dependency.installCommand);
      setInstallOutput(installed.output || '');
      if (!installed.success) setError(installed.error || `Could not install ${dependency.name}.`);
      setRequirements(await checkWakeRequirements(config));
    } catch (cause) { reportError(cause); }
    finally { setBusy(false); }
  };
  const finishSigning = async (skip: boolean) => {
    const draft = skip ? { ...config, sign: false } : config;
    if (skip) update(draft);
    setStep('Save');
    await check(draft);
  };
  const save = async () => {
    setBusy(true); setError('');
    try {
      const checked = await checkWakeRequirements(config);
      setRequirements(checked);
      if (!checked.errors.length) setResult(await saveWakeSetup(config));
    } catch (cause) { reportError(cause); }
    finally { setBusy(false); }
  };
  const finish = async () => {
    setBusy(true); setError('');
    try { await complete(); setFinished(true); }
    catch (cause) { reportError(cause); }
    finally { setBusy(false); }
  };
  const metadataIncomplete = step === 'Project' && (!config.project?.info.productName.trim() || !config.project.info.productIdentifier.trim() || !config.project.info.version.trim());
  const cannotContinue = busy || metadataIncomplete || (step === 'Targets' && config.targets.length === 0);

  return <Layout sidebar={sidebar({ labels, index })}>
    {step === 'Signing' && !result && signingPlatforms.length > 0 ? <SigningStep platforms={signingPlatforms} wakeMode mobileSigning={mobileSigning} onNext={() => void finishSigning(false)} onSkip={() => void finishSigning(true)} onBack={() => setStep('Tools')} canGoBack /> : (
      <main className="flex-1 flex flex-col min-h-0">
        <header className="px-10 pt-8 pb-5 shrink-0">
          <p className="text-xs font-medium text-red-500 mb-2">WAKE SETUP</p>
          <h1 className="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">{result ? 'Build configuration saved' : titles[step]}</h1>
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400 break-all">{initial.projectRoot}</p>
        </header>
        <div className="flex-1 min-h-0 overflow-y-auto px-10 pb-4 space-y-5">
          {error && <p role="alert" className="rounded-lg bg-red-50 dark:bg-red-900/20 p-3 text-sm text-red-700 dark:text-red-300">{error}</p>}
          {result ? <>
            <p className="text-gray-700 dark:text-gray-200">Saved to <code className="break-all">{result.path}</code>.</p>
            <p className="text-sm text-gray-500">Run this command from the project directory:</p>
            <pre className="p-4 rounded-lg bg-gray-100 dark:bg-gray-950 text-sm text-gray-800 dark:text-gray-200 overflow-x-auto"><code>{result.command}</code></pre>
            {result.config.prebuilt && <p className="text-sm text-gray-600 dark:text-gray-300">The target binaries must already exist in the configured output directory. This command packages those files; future builds still compile normally.</p>}
            {requirements && <RequirementMessages requirements={requirements} />}
            {finished && <p role="status" className="text-green-700 dark:text-green-400">Setup is complete. You can close this window.</p>}
          </> : <>
            <WakeSetupChoices step={step} config={config} initial={initial} update={update} />
            {step === 'Signing' && <>
              <p className="text-sm text-gray-600 dark:text-gray-300">Mobile signing uses the existing ios.signing or android.signing fields in wails.yaml. This page does not configure mobile credentials.</p>
              <button className={buttonClass} onClick={() => void finishSigning(true)}>Set up signing later</button>
            </>}
            {step === 'Tools' && <>
              <div className="flex items-center justify-between gap-4"><p className="text-sm text-gray-500">Requirements are resolved for your selections on {initial.hostOS}/{initial.hostArch}.</p><button className={buttonClass} disabled={busy} onClick={() => void check()}>{busy ? 'Checking…' : 'Recheck'}</button></div>
              {requirements && <WakeRequirements requirements={requirements} busy={busy} onInstall={(dependency) => void install(dependency)} />}
              {installOutput && <pre className="text-xs p-3 bg-gray-100 dark:bg-gray-950 text-gray-600 dark:text-gray-300 whitespace-pre-wrap max-h-40 overflow-auto">{installOutput}</pre>}
            </>}
            {step === 'Save' && <>
              <p className="text-sm text-gray-600 dark:text-gray-300">Save selected targets and packages to <code>profiles.setup</code> in wails.yaml. The backend applies at project level; other profiles may inherit it. Existing hooks and signing configuration are kept.</p>
              <dl className="grid grid-cols-2 gap-3 text-sm text-gray-700 dark:text-gray-300"><dt>Application</dt><dd>{config.applicationType}</dd><dt>Backend</dt><dd>{config.backend === 'purego' ? 'Pure Go (experimental)' : 'Standard'}</dd><dt>Operation</dt><dd>{config.prebuilt ? 'Package existing binaries' : 'Build from source'}</dd><dt>Signing</dt><dd>{config.sign ? 'Enabled' : 'Disabled'}</dd></dl>
              <ul className="space-y-2 text-sm text-gray-700 dark:text-gray-300">{config.targets.map((target) => <li key={targetKey(target)}>{targetKey(target)} — {target.formats.join(', ') || 'default output'}</li>)}</ul>
              {requirements && <RequirementMessages requirements={requirements} />}
              {requirements?.command && <pre className="p-3 rounded-lg bg-gray-100 dark:bg-gray-950 text-sm text-gray-700 dark:text-gray-300 overflow-x-auto"><code>{requirements.command}</code></pre>}
              <button className={buttonClass} disabled={busy} onClick={() => void check()}>Recheck requirements</button>
            </>}
          </>}
        </div>
        <footer className="shrink-0 px-10 py-5 flex items-center justify-center gap-3">
          {result ? <button className={nextClass} disabled={busy || finished} onClick={() => void finish()}>{finished ? 'Finished' : busy ? 'Finishing…' : 'Finish'}</button> : <>
            {index > 0 && <button className={buttonClass} disabled={busy} onClick={() => { setError(''); setStep(labels[index - 1]); }}>Back</button>}
            <button className={nextClass} disabled={cannotContinue || (step === 'Save' && (!requirements || requirements.errors.length > 0))} onClick={() => void (step === 'Save' ? save() : next())}>{busy ? 'Checking…' : step === 'Save' ? 'Save wails.yaml' : 'Continue'}</button>
          </>}
        </footer>
      </main>
    )}
  </Layout>;
}
