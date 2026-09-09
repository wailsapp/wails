import { useState, ReactNode } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { InitData } from '../types';
import { createProject } from '../api';
import wailsLogoWhite from '../assets/wails-logo-white-text.svg';
import wailsLogoBlack from '../assets/wails-logo-black-text.svg';

type Theme = 'light' | 'dark';
type Step = 'project' | 'framework' | 'language' | 'bindings' | 'details' | 'done';

const pageVariants = { initial: { opacity: 0 }, animate: { opacity: 1 }, exit: { opacity: 0 } };

const inputCls =
  'w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500';

const FRAMEWORKS = [
  { id: 'vanilla', name: 'Vanilla', description: 'Plain JavaScript / TypeScript', icon: 'javascript' },
  { id: 'react', name: 'React', description: 'React with Vite', icon: 'react' },
  { id: 'vue', name: 'Vue', description: 'Vue 3 with Vite', icon: 'vue' },
  { id: 'svelte', name: 'Svelte', description: 'Svelte with Vite', icon: 'svelte' },
];

function slugIdentifier(name: string): string {
  const slug = name.toLowerCase().replace(/[^a-z0-9]+/g, '');
  return `com.example.${slug || 'app'}`;
}

// One focused choice per page, presented like the setup wizard.
function Page({ title, subtitle, children, onBack, onNext, nextLabel = 'Continue', nextDisabled, busy }: {
  title: string; subtitle: string; children: ReactNode;
  onBack?: () => void; onNext: () => void; nextLabel?: string; nextDisabled?: boolean; busy?: boolean;
}) {
  return (
    <motion.main variants={pageVariants} initial="initial" animate="animate" exit="exit" transition={{ duration: 0.25 }} className="flex-1 flex flex-col">
      <header className="text-center mb-6 flex-shrink-0 px-10 pt-10">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-1.5 tracking-tight">{title}</h1>
        <p className="text-base text-gray-500 dark:text-gray-400">{subtitle}</p>
      </header>
      <div className="flex-1 overflow-y-auto scrollbar-thin min-h-0 px-10 flex flex-col justify-center">{children}</div>
      <div className="flex-shrink-0 pt-4 pb-6 flex items-center justify-center gap-3">
        {onBack && (
          <button onClick={onBack} disabled={busy}
            className="px-4 py-2 rounded-lg text-sm font-medium border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 disabled:opacity-50">Back</button>
        )}
        <button onClick={onNext} disabled={nextDisabled || busy}
          className="px-5 py-2 rounded-lg text-sm font-medium border border-red-500 text-red-600 dark:text-red-400 hover:bg-red-500/10 disabled:opacity-50 disabled:cursor-not-allowed">{nextLabel}</button>
      </div>
    </motion.main>
  );
}

export default function InitFlow({ data, theme, toggleTheme }: { data: InitData; theme: Theme; toggleTheme: () => void; }) {
  const has = (n: string) => data.templates.some((t) => t.name === n);
  const resolveTemplate = (fw: string, ts: boolean) => (!ts && has(`${fw}-js`) ? `${fw}-js` : fw);

  const baseDefault = (data.defaultTemplate || data.templateName || 'vanilla').replace(/-js$/, '');
  const defaultTS = !(data.defaultTemplate || data.templateName || '').endsWith('-js');
  const defaultFwName = FRAMEWORKS.find((f) => f.id === baseDefault)?.name ?? 'Vanilla';
  // Human-readable summary of the configured defaults, e.g. "React · TypeScript · Interfaces".
  const defaultSummary = [defaultFwName, defaultTS ? 'TypeScript' : 'JavaScript', ...(defaultTS ? [data.useInterfaces ? 'Interfaces' : 'Classes'] : [])].join(' · ');

  const [framework, setFramework] = useState(FRAMEWORKS.some((f) => f.id === baseDefault) ? baseDefault : 'vanilla');
  const [preferTS, setPreferTS] = useState(defaultTS);

  const [step, setStep] = useState<Step>('project');
  const [form, setForm] = useState<InitData>({ ...data });
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState('');
  const set = <K extends keyof InitData>(k: K, v: InitData[K]) => setForm((f) => ({ ...f, [k]: v }));

  // Linear flow; the Bindings page only applies to TypeScript.
  const flow: Step[] = ['project', 'framework', 'language', ...(preferTS ? (['bindings'] as Step[]) : []), 'details'];
  const goNext = () => { const i = flow.indexOf(step); if (i >= 0 && i < flow.length - 1) setStep(flow[i + 1]); };
  const goBack = () => { const i = flow.indexOf(step); if (i > 0) setStep(flow[i - 1]); };

  // Derive name-based config defaults when leaving the choices for Details.
  const enterDetails = () => {
    const name = form.projectName.trim();
    setForm((f) => ({
      ...f,
      templateName: resolveTemplate(framework, preferTS),
      productName: f.productName.trim() ? f.productName : name,
      productDescription: f.productDescription.trim() ? f.productDescription : `A ${name} application`,
      productIdentifier: f.productIdentifier.trim() ? f.productIdentifier : slugIdentifier(name),
    }));
    setStep('details');
  };
  const advance = () => { (flow[flow.indexOf(step) + 1] === 'details') ? enterDetails() : goNext(); };

  // Apply the configured global defaults and skip straight to Details.
  const useGlobalDefaults = () => {
    setFramework(baseDefault);
    setPreferTS(defaultTS);
    const name = form.projectName.trim();
    setForm((f) => ({
      ...f,
      useInterfaces: data.useInterfaces,
      templateName: resolveTemplate(baseDefault, defaultTS),
      productName: f.productName.trim() ? f.productName : name,
      productDescription: f.productDescription.trim() ? f.productDescription : `A ${name} application`,
      productIdentifier: f.productIdentifier.trim() ? f.productIdentifier : slugIdentifier(name),
    }));
    setStep('details');
  };

  const create = async () => {
    setError(''); setCreating(true);
    try {
      const result = await createProject({ ...form, templateName: resolveTemplate(framework, preferTS) });
      if (result.success) setStep('done');
      else { setError(result.error || 'Failed to create project'); }
    } catch { setError('Failed to create project'); }
    finally { setCreating(false); }
  };

  const stages: { key: Step; label: string }[] = [
    { key: 'project', label: 'Project' },
    { key: 'framework', label: 'Framework' },
    { key: 'language', label: 'Language' },
    ...(preferTS ? [{ key: 'bindings' as Step, label: 'Bindings' }] : []),
    { key: 'details', label: 'Details' },
    { key: 'done', label: 'Create' },
  ];
  const currentIndex = stages.findIndex((s) => s.key === step);

  return (
    <div className="fixed inset-0 overflow-auto bg-gray-50 dark:bg-[#0f0f0f] transition-colors">
      <div className="fixed inset-0 bg-center bg-cover bg-no-repeat pointer-events-none opacity-50 dark:opacity-40" style={{ backgroundImage: "url('/digital_wales_master.webp')" }} />
      <div className="relative min-h-full min-w-full w-fit flex items-center justify-center p-4">
        <div className="glass-card rounded-2xl flex overflow-hidden relative z-10" style={{ aspectRatio: '3 / 2', width: 'clamp(48rem, min(100vw - 2rem, (100vh - 2rem) * 1.5), 75rem)' }}>
          {/* Sidebar */}
          <aside className="w-48 flex-shrink-0 bg-gray-100/50 dark:bg-[#0a0e16]/40 backdrop-blur-[50px] backdrop-saturate-[1.8] border-r border-gray-200 dark:border-white/10 flex flex-col">
            <div className="p-6 flex justify-center">
              <img src={theme === 'dark' ? wailsLogoWhite : wailsLogoBlack} alt="Wails logo" className="h-24 object-contain" />
            </div>
            <nav className="flex-1 px-4 py-2">
              <ol className="space-y-1">
                {stages.map((s, i) => {
                  // On the done screen every step — including Create — is complete.
                  const isCurrent = s.key === step && step !== 'done';
                  const isDone = i < currentIndex || step === 'done';
                  return (
                    <li key={s.key}>
                      <div className={`flex items-center gap-3 px-3 py-2.5 rounded-lg ${isCurrent ? 'bg-white dark:bg-gray-800/80' : ''}`}>
                        <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium flex-shrink-0 ${isDone ? 'bg-green-500 text-white' : isCurrent ? 'bg-red-500 text-white' : 'bg-gray-300 dark:bg-gray-700 text-gray-700 dark:text-gray-400'}`}>{isDone ? '✓' : i + 1}</div>
                        <span className={`text-sm font-medium ${isCurrent ? 'text-gray-900 dark:text-white' : 'text-gray-600 dark:text-gray-300'}`}>{s.label}</span>
                      </div>
                    </li>
                  );
                })}
              </ol>
            </nav>
            <div className="p-4 flex justify-center">
              <button onClick={toggleTheme} className="p-1 hover:opacity-70 transition-opacity rounded" aria-label="Toggle theme">
                {theme === 'dark' ? (
                  <svg className="w-4 h-4 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.36 6.36l-.7-.7M6.34 6.34l-.7-.7m12.72 0l-.7.7M6.34 17.66l-.7.7M16 12a4 4 0 11-8 0 4 4 0 018 0z" /></svg>
                ) : (
                  <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" /></svg>
                )}
              </button>
            </div>
          </aside>

          {/* Content */}
          <div className="flex-1 flex flex-col min-w-[36rem] bg-white/85 dark:bg-[#0a0e16]/85 backdrop-blur-[30px] backdrop-saturate-150 relative">
            <AnimatePresence mode="wait">
              {step === 'project' && (
                <Page key="project" title="Create a Wails project" subtitle="What's your project called?" onNext={advance} nextDisabled={!form.projectName.trim()}>
                  <div className="max-w-md mx-auto w-full">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Project Name</label>
                    <input autoFocus type="text" value={form.projectName} onChange={(e) => set('projectName', e.target.value)} placeholder="my-app" className={inputCls} />
                    <p className="text-xs text-gray-400 dark:text-gray-500 mt-2 truncate text-center">
                      Will be created in <code className="font-mono">{form.baseDir}/{form.projectName.trim() || 'my-app'}</code>
                    </p>
                    <div className="mt-6 pt-5 border-t border-gray-200 dark:border-white/10 text-center">
                      <button type="button" onClick={useGlobalDefaults} disabled={!form.projectName.trim()}
                        className="px-4 py-2 rounded-lg text-sm font-medium border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed">
                        Use global defaults
                      </button>
                      <p className="text-xs text-gray-400 dark:text-gray-500 mt-1.5">{defaultSummary} — skips straight to details</p>
                    </div>
                  </div>
                </Page>
              )}

              {step === 'framework' && (
                <Page key="framework" title="Choose a framework" subtitle="The frontend stack for your app" onBack={goBack} onNext={advance}>
                  <div className="flex flex-col gap-3 max-w-md mx-auto w-full" role="radiogroup" aria-label="Framework">
                    {FRAMEWORKS.map((fw) => {
                      const selected = framework === fw.id;
                      const logo = fw.id === 'vanilla' ? (preferTS ? 'typescript' : 'javascript') : fw.icon;
                      return (
                        <button key={fw.id} type="button" role="radio" aria-checked={selected} onClick={() => setFramework(fw.id)}
                          className={`flex items-center gap-4 text-left rounded-xl p-4 transition-all border-2 ${selected ? 'border-red-500 bg-red-500/10 shadow-lg shadow-red-500/10' : 'border-gray-200 dark:border-white/10 bg-gray-100 dark:bg-white/5 hover:bg-gray-200 dark:hover:bg-white/10'}`}>
                          <img src={`/logos/${logo}.svg`} alt="" aria-hidden="true" className="w-10 h-10 flex-shrink-0" />
                          <span className="min-w-0">
                            <span className="block text-sm font-semibold text-gray-900 dark:text-white">{fw.name}</span>
                            <span className="block text-xs text-gray-500 dark:text-gray-400 truncate">{fw.description}</span>
                          </span>
                        </button>
                      );
                    })}
                  </div>
                </Page>
              )}

              {step === 'language' && (
                <Page key="language" title="Language Preference" subtitle="Choose your preferred language" onBack={goBack} onNext={advance}>
                  <div className="flex justify-center gap-4" role="radiogroup" aria-label="Language">
                    {([['javascript', 'JavaScript', 'Dynamic typing', false], ['typescript', 'TypeScript', 'Type safety', true]] as [string, string, string, boolean][]).map(([logo, name, sub, ts]) => {
                      const selected = preferTS === ts;
                      const ring = ts ? 'border-blue-400 bg-blue-400/10 shadow-blue-400/20' : 'border-yellow-400 bg-yellow-400/10 shadow-yellow-400/20';
                      return (
                        <button key={name} type="button" role="radio" aria-checked={selected} onClick={() => setPreferTS(ts)}
                          className={`w-40 h-48 rounded-xl p-5 flex flex-col items-center justify-center gap-3 transition-all border-2 ${selected ? `${ring} shadow-lg` : 'border-gray-200 dark:border-white/10 bg-gray-100 dark:bg-white/5 hover:bg-gray-200 dark:hover:bg-white/10'}`}>
                          <img src={`/logos/${logo}.svg`} alt="" aria-hidden="true" className="w-14 h-14" />
                          <span className="text-lg font-semibold text-gray-900 dark:text-white">{name}</span>
                          <span className="text-xs text-gray-500 dark:text-gray-400">{sub}</span>
                        </button>
                      );
                    })}
                  </div>
                </Page>
              )}

              {step === 'bindings' && (
                <Page key="bindings" title="TypeScript Bindings" subtitle="How Go structs are represented in TypeScript" onBack={goBack} onNext={advance}>
                  <div className="flex justify-center gap-4" role="radiogroup" aria-label="Binding style">
                    <button type="button" role="radio" aria-checked={form.useInterfaces} onClick={() => set('useInterfaces', true)}
                      className={`w-60 rounded-xl p-4 flex flex-col items-start gap-2 transition-all border-2 text-left ${form.useInterfaces ? 'border-blue-400 bg-blue-400/10 shadow-lg shadow-blue-400/20' : 'border-gray-200 dark:border-white/10 bg-gray-100 dark:bg-white/5 hover:bg-gray-200 dark:hover:bg-white/10'}`}>
                      <span className="text-base font-semibold text-gray-900 dark:text-white">Interfaces</span>
                      <pre className="text-[11px] leading-tight text-gray-700 dark:text-white/70 font-mono bg-gray-100 dark:bg-black/30 p-2 rounded-lg w-full overflow-x-auto" aria-hidden="true">{`interface Person {
  name: string;
  age: number;
}`}</pre>
                      <ul className="text-[11px] text-gray-500 dark:text-white/50 space-y-0.5 list-disc list-inside">
                        <li>Lightweight types</li>
                        <li>No runtime code</li>
                        <li>Simpler output</li>
                      </ul>
                    </button>
                    <button type="button" role="radio" aria-checked={!form.useInterfaces} onClick={() => set('useInterfaces', false)}
                      className={`w-60 rounded-xl p-4 flex flex-col items-start gap-2 transition-all border-2 text-left ${!form.useInterfaces ? 'border-purple-400 bg-purple-400/10 shadow-lg shadow-purple-400/20' : 'border-gray-200 dark:border-white/10 bg-gray-100 dark:bg-white/5 hover:bg-gray-200 dark:hover:bg-white/10'}`}>
                      <span className="text-base font-semibold text-gray-900 dark:text-white">Classes</span>
                      <pre className="text-[11px] leading-tight text-gray-700 dark:text-white/70 font-mono bg-gray-100 dark:bg-black/30 p-2 rounded-lg w-full overflow-x-auto" aria-hidden="true">{`class Person {
  name: string;
  static createFrom(src) {
    return new Person(src);
  }
}`}</pre>
                      <ul className="text-[11px] text-gray-500 dark:text-white/50 space-y-0.5 list-disc list-inside">
                        <li>Factory methods</li>
                        <li>Default initialization</li>
                        <li>More verbose</li>
                      </ul>
                    </button>
                  </div>
                </Page>
              )}

              {step === 'details' && (
                <Page key="details" title="Project details" subtitle="Written to build/config.yml" onBack={goBack} onNext={create} nextLabel={creating ? 'Creating…' : 'Create Project'} busy={creating} nextDisabled={!form.projectName.trim()}>
                  <div className="max-w-md mx-auto w-full space-y-3">
                    {([
                      ['productName', 'Product Name', 'My App'],
                      ['productCompany', 'Company', 'My Company'],
                      ['productIdentifier', 'Bundle Identifier', 'com.mycompany.myapp'],
                      ['productDescription', 'Description', 'A Wails application'],
                      ['productVersion', 'Version', '0.1.0'],
                      ['productCopyright', 'Copyright', '© 2026, My Company'],
                      ['productComments', 'Comments', 'Optional notes'],
                    ] as [keyof InitData, string, string][]).map(([k, label, ph]) => (
                      <div key={k}>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{label}</label>
                        <input type="text" value={String(form[k] ?? '')} onChange={(e) => set(k, e.target.value as never)} placeholder={ph} className={inputCls} />
                      </div>
                    ))}
                    <p className="text-xs text-gray-400 dark:text-gray-500">Advanced config (icons, file associations, custom URL schemes, iOS overrides) can be edited in <code className="font-mono">build/config.yml</code> after creation.</p>
                    {error && (
                      <div className="p-3 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800"><p className="text-sm text-red-600 dark:text-red-400">{error}</p></div>
                    )}
                  </div>
                </Page>
              )}

              {step === 'done' && (
                <motion.main key="done" variants={pageVariants} initial="initial" animate="animate" exit="exit" transition={{ duration: 0.25 }} className="flex-1 flex flex-col items-center justify-center px-10">
                  <div className="w-20 h-20 rounded-full bg-green-500/20 flex items-center justify-center mb-6">
                    <svg className="w-10 h-10 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
                  </div>
                  <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">Project created</h2>
                  <p className="text-gray-500 dark:text-gray-400 text-center max-w-sm">
                    <code className="font-mono">{form.projectName}</code> was scaffolded from the <code className="font-mono">{resolveTemplate(framework, preferTS)}</code> template at <code className="font-mono break-all">{form.baseDir}/{form.projectName}</code>. You can close this window and return to the terminal.
                  </p>
                </motion.main>
              )}
            </AnimatePresence>
          </div>
        </div>
      </div>
    </div>
  );
}
