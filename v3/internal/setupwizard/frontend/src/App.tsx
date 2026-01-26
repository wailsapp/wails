import { useState, useEffect, createContext, useContext, ReactNode, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { DependencyStatus, SystemInfo, DockerStatus, GlobalDefaults } from './types';
import { checkDependencies, getState, getDockerStatus, buildDockerImage, getDefaults, saveDefaults, subscribeDockerStatus } from './api';
import wailsLogoWhite from './assets/wails-logo-white-text.svg';
import wailsLogoBlack from './assets/wails-logo-black-text.svg';
import SigningStep from './components/SigningStep';

type OOBEStep =
  | 'splash'
  | 'checking'
  | 'deps-ready'
  | 'deps-missing'
  | 'cross-platform'
  | 'docker-options'
  | 'sdk-license'
  | 'docker-setup'
  | 'projects'
  | 'signing'
  | 'language-select'
  | 'binding-style'
  | 'template-select'
  | 'complete';

// Framework template info
type FrameworkTemplate = {
  id: string;
  name: string;
  description: string;
  color: string;
  icon: string;
};

// Wizard stages for sidebar progress display
type WizardStage = 'welcome' | 'dependencies' | 'platform' | 'identity' | 'templates' | 'complete';

function getWizardStage(step: OOBEStep): WizardStage {
  switch (step) {
    case 'splash':
    case 'checking':
      return 'welcome';
    case 'deps-ready':
    case 'deps-missing':
      return 'dependencies';
    case 'cross-platform':
    case 'docker-options':
    case 'sdk-license':
    case 'docker-setup':
      return 'platform';
    case 'projects':
    case 'signing':
      return 'identity';
    case 'language-select':
    case 'binding-style':
    case 'template-select':
      return 'templates';
    case 'complete':
      return 'complete';
    default:
      return 'welcome';
  }
}

// Get stage index for progress tracking (1-6)
function getStageIndex(stage: WizardStage): number {
  const stages: WizardStage[] = ['welcome', 'dependencies', 'platform', 'identity', 'templates', 'complete'];
  return stages.indexOf(stage) + 1;
}

type Theme = 'light' | 'dark';

// Theme context
const ThemeContext = createContext<{ theme: Theme; toggleTheme: () => void }>({
  theme: 'dark',
  toggleTheme: () => {}
});

const useTheme = () => useContext(ThemeContext);

// Page animation variants - fade only
const pageVariants = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 }
};


// Sidebar with progress steps (1:4 ratio - sidebar is 20% of total width)
function Sidebar({ currentStep, dockerStatus, buildingDocker }: {
  currentStep: OOBEStep;
  dockerStatus: DockerStatus | null;
  buildingDocker: boolean;
}) {
  const { theme, toggleTheme } = useTheme();
  const currentStage = getWizardStage(currentStep);
  const currentIndex = getStageIndex(currentStage);
  const [showBugOverlay, setShowBugOverlay] = useState(false);
  const [bugReportUrl, setBugReportUrl] = useState('');

  const stages = [
    { key: 'welcome' as const, label: 'Welcome' },
    { key: 'dependencies' as const, label: 'Dependencies' },
    { key: 'platform' as const, label: 'Platform' },
    { key: 'identity' as const, label: 'Projects' },
    { key: 'templates' as const, label: 'Templates' },
    { key: 'complete' as const, label: 'Complete' },
  ];

  const handleSponsorClick = () => {
    window.open('https://github.com/sponsors/leaanthony', '_blank', 'noopener,noreferrer');
  };

  const handleReportBug = async () => {
    try {
      const { reportBug } = await import('./api');
      const result = await reportBug(currentStep);
      if (result.body && result.url) {
        await navigator.clipboard.writeText(result.body);
        setBugReportUrl(result.url);
        setShowBugOverlay(true);
      }
    } catch (e) {
      console.error('Failed to report bug:', e);
    }
  };

  const handleOpenGitHub = () => {
    window.open(bugReportUrl, '_blank', 'noopener,noreferrer');
    setShowBugOverlay(false);
  };

  const isDockerBuilding = buildingDocker;

  return (
    <aside
      className="w-48 flex-shrink-0 bg-gray-100/80 dark:bg-transparent dark:glass-sidebar border-r border-gray-200 dark:border-transparent flex flex-col"
      aria-label="Setup progress"
    >
      <div className="p-6 flex justify-center">
        <img
          src={theme === 'dark' ? wailsLogoWhite : wailsLogoBlack}
          alt="Wails logo"
          className="h-24 object-contain"
        />
      </div>

      <nav className="flex-1 px-4 py-2" aria-label="Setup steps">
        <ol className="space-y-1">
          {stages.map((stage, index) => {
            const stageIndex = index + 1;
            const isCurrent = stage.key === currentStage;
            const isCompleted = stageIndex < currentIndex;
            const stepStatus = isCompleted ? 'completed' : isCurrent ? 'current' : 'upcoming';
            const showDockerSubtask = stage.key === 'platform' && isDockerBuilding;

            return (
              <li
                key={stage.key}
                aria-current={isCurrent ? 'step' : undefined}
              >
                <div
                  className={`flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors ${
                    isCurrent ? 'bg-white dark:bg-gray-800/80' : ''
                  }`}
                  aria-label={`Step ${stageIndex}: ${stage.label}, ${stepStatus}`}
                >
                  <div
                    className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium flex-shrink-0 ${
                      isCompleted
                        ? 'bg-green-500 text-white'
                        : isCurrent
                          ? 'bg-red-500 text-white'
                          : 'bg-gray-300 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
                    }`}
                    aria-hidden="true"
                  >
                    {isCompleted ? (
                      <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                    ) : (
                      stageIndex
                    )}
                  </div>

                  <span
                    className={`text-sm font-medium ${
                      isCurrent
                        ? 'text-gray-900 dark:text-white'
                        : isCompleted
                          ? 'text-green-700 dark:text-gray-400'
                          : 'text-gray-400 dark:text-gray-600'
                    }`}
                  >
                    {stage.label}
                  </span>
                </div>

                {showDockerSubtask && (
                  <div className="ml-9 mt-1 mb-3">
                    <div className="flex items-center justify-between text-[10px] text-gray-500 dark:text-gray-400 mb-1">
                      <span>{dockerStatus?.pullMessage || 'Connecting'}</span>
                      <span>{dockerStatus?.pullProgress || 0}%</span>
                    </div>
                    <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-1">
                      <div
                        className="bg-blue-500 h-1 rounded-full transition-all duration-300"
                        style={{ width: `${dockerStatus?.pullProgress || 0}%` }}
                      />
                    </div>
                  </div>
                )}
              </li>
            );
          })}
        </ol>
      </nav>



      <div className="p-4 flex justify-center gap-3">
        <button
          onClick={handleReportBug}
          className="p-1 hover:opacity-70 transition-opacity focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 rounded"
          aria-label="Report a bug"
          title="Report a bug"
        >
          <svg className="w-4 h-4 text-gray-500 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </button>
        <button
          onClick={handleSponsorClick}
          className="p-1 hover:opacity-70 transition-opacity focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 rounded"
          aria-label="Sponsor Wails on GitHub"
        >
          <svg className="w-4 h-4 text-red-500" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </button>
        <button
          onClick={toggleTheme}
          className="p-1 hover:opacity-70 transition-opacity focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 rounded"
          aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
        >
          {theme === 'dark' ? (
            <svg className="w-4 h-4 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          ) : (
            <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          )}
        </button>
      </div>

      {/* Bug Report Overlay */}
      {showBugOverlay && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-sm mx-4 shadow-xl">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-10 h-10 rounded-full bg-green-100 dark:bg-green-900 flex items-center justify-center">
                <svg className="w-5 h-5 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Template Copied!</h3>
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
              The issue template has been copied to your clipboard. Click below to open the GitHub issue and paste it into a new comment.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setShowBugOverlay(false)}
                className="flex-1 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleOpenGitHub}
                className="flex-1 px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
              >
                Open GitHub
              </button>
            </div>
          </div>
        </div>
      )}
    </aside>
  );
}

function PageTemplate({
  title,
  subtitle,
  children,
  primaryAction,
  primaryLabel,
  secondaryAction,
  secondaryLabel,
  primaryDisabled = false,
  onBack,
  canGoBack = false
}: {
  title: string;
  subtitle: string;
  children?: ReactNode;
  primaryAction?: () => void;
  primaryLabel?: string;
  secondaryAction?: () => void;
  secondaryLabel?: string;
  primaryDisabled?: boolean;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, [title]);

  const actionsElement = (primaryAction || secondaryAction) ? (
    <div className="flex-shrink-0 pt-4 pb-6 flex flex-col items-center gap-1.5" role="group" aria-label="Page actions">
      <div className="flex items-center gap-3">
        {canGoBack && onBack && (
          <button
            onClick={onBack}
            className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
          >
            Back
          </button>
        )}
        {primaryAction && primaryLabel && (
          <button
            onClick={primaryAction}
            disabled={primaryDisabled}
            className={`px-5 py-2 rounded-lg text-sm font-medium transition-colors border focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 ${
              primaryDisabled
                ? 'border-gray-300 dark:border-gray-700 text-gray-400 cursor-not-allowed'
                : 'border-red-500 text-red-600 dark:text-red-400 hover:bg-red-500/10'
            }`}
            aria-disabled={primaryDisabled}
          >
            {primaryLabel}
          </button>
        )}
      </div>
      {secondaryAction && secondaryLabel && (
        <button
          onClick={secondaryAction}
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 rounded"
        >
          {secondaryLabel}
        </button>
      )}
    </div>
  ) : null;

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col"
      aria-labelledby="page-title"
    >
      <header className="text-center mb-6 flex-shrink-0 px-10 pt-10">
        <h1
          ref={headingRef}
          id="page-title"
          className="text-2xl font-semibold text-gray-900 dark:text-white mb-1.5 tracking-tight focus:outline-none"
          tabIndex={-1}
        >
          {title}
        </h1>
        <p className="text-base text-gray-500 dark:text-gray-400">{subtitle}</p>
      </header>

      <div className="flex-1 overflow-y-auto scrollbar-thin min-h-0 px-10">
        {children}
      </div>

      {actionsElement}
    </motion.main>
  );
}

// Splash Page - simple welcome with Let's Start
function SplashPage({ onNext }: { onNext: () => void }) {
  const { theme } = useTheme();
  const headingRef = useRef<HTMLHeadingElement>(null);
  const startButtonRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, []);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && document.activeElement === startButtonRef.current) {
      onNext();
    }
  };

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
      aria-labelledby="splash-title"
      onKeyDown={handleKeyDown}
    >
      <motion.div
        className="text-center mb-10"
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
      >
        <div className="flex justify-center">
          <img
            src={theme === 'dark' ? wailsLogoWhite : wailsLogoBlack}
            alt=""
            width={280}
            className="object-contain"
            style={{ filter: 'drop-shadow(0 0 60px rgba(239, 68, 68, 0.4))' }}
            aria-hidden="true"
          />
        </div>
      </motion.div>

      <motion.div
        className="text-center px-8 max-w-lg"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <div className="flex items-center justify-center gap-2 mb-3">
          <h1
            ref={headingRef}
            id="splash-title"
            className="text-2xl font-semibold text-gray-900 dark:text-white tracking-tight focus:outline-none"
            tabIndex={-1}
          >
            Welcome to Wails
          </h1>
          <span
            className="px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide rounded-full bg-amber-500/20 text-amber-600 dark:text-amber-400 border border-amber-500/30"
            role="status"
            aria-label="This setup wizard is experimental"
          >
            Experimental
          </span>
        </div>
        <p className="text-base text-gray-600 dark:text-gray-300 leading-relaxed mb-8">
          Build beautiful cross-platform apps using Go and web technologies
        </p>
      </motion.div>

      <motion.button
        ref={startButtonRef}
        onClick={onNext}
        className="px-6 py-2.5 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.4 }}
      >
        Let's Start
      </motion.button>
    </motion.main>
  );
}

function CheckingPage() {
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, []);

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-start pt-[15%]"
      aria-labelledby="checking-title"
      aria-busy="true"
    >
      <motion.div
        className="w-12 h-12 border-3 border-gray-300 dark:border-gray-600 border-t-red-500 rounded-full mb-6"
        animate={{ rotate: 360 }}
        transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
        role="status"
        aria-label="Loading"
      />
      <h2
        ref={headingRef}
        id="checking-title"
        className="text-xl font-semibold text-gray-900 dark:text-white mb-2 focus:outline-none"
        tabIndex={-1}
      >
        Checking your system...
      </h2>
      <p className="text-gray-500 dark:text-gray-400" aria-live="polite">
        This will only take a moment
      </p>
    </motion.main>
  );
}

function DepsReadyPage({ onNext, onBack, canGoBack }: { onNext: () => void; onBack?: () => void; canGoBack?: boolean }) {
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, []);

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
      aria-labelledby="deps-ready-title"
    >
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', stiffness: 200, damping: 15 }}
        className="w-20 h-20 rounded-full bg-green-500/20 flex items-center justify-center mb-6"
        aria-hidden="true"
      >
        <svg className="w-10 h-10 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
        </svg>
      </motion.div>

      <h2
        ref={headingRef}
        id="deps-ready-title"
        className="text-2xl font-semibold text-gray-900 dark:text-white mb-2 focus:outline-none"
        tabIndex={-1}
      >
        All dependencies installed
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-sm">
        Your system has everything needed to build Wails apps
      </p>

      <div className="flex items-center gap-3" role="group" aria-label="Navigation">
        {canGoBack && onBack && (
          <button
            onClick={onBack}
            className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
          >
            Back
          </button>
        )}
        <button
          onClick={onNext}
          className="px-5 py-2 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
        >
          Continue
        </button>
      </div>
    </motion.main>
  );
}

// Deps Missing Page - show what's missing with install command
function DepsMissingPage({
  dependencies,
  onRetry,
  onContinue,
  onBack,
  canGoBack
}: {
  dependencies: DependencyStatus[];
  onRetry: () => void;
  onContinue: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  const [copied, setCopied] = useState(false);
  const missingDeps = dependencies.filter(d => !d.installed && d.required);

  // Build combined install command
  const combinedInstallCommand = (() => {
    const systemCommands = missingDeps
      .filter(d => d.installCommand?.startsWith('sudo '))
      .map(d => d.installCommand!);

    if (systemCommands.length === 0) return null;

    const pacmanPkgs: string[] = [];
    const aptPkgs: string[] = [];
    const dnfPkgs: string[] = [];

    for (const cmd of systemCommands) {
      if (cmd.includes('pacman -S')) {
        const match = cmd.match(/pacman -S\s+(.+)/);
        if (match) pacmanPkgs.push(...match[1].split(/\s+/));
      } else if (cmd.includes('apt install')) {
        const match = cmd.match(/apt install\s+(.+)/);
        if (match) aptPkgs.push(...match[1].split(/\s+/));
      } else if (cmd.includes('dnf install')) {
        const match = cmd.match(/dnf install\s+(.+)/);
        if (match) dnfPkgs.push(...match[1].split(/\s+/));
      }
    }

    if (pacmanPkgs.length > 0) return `sudo pacman -S ${pacmanPkgs.join(' ')}`;
    if (aptPkgs.length > 0) return `sudo apt install ${aptPkgs.join(' ')}`;
    if (dnfPkgs.length > 0) return `sudo dnf install ${dnfPkgs.join(' ')}`;
    return null;
  })();

  const copyCommand = () => {
    if (combinedInstallCommand) {
      navigator.clipboard.writeText(combinedInstallCommand);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <PageTemplate
      title="Almost there!"
      subtitle="A few things need to be installed first"
      primaryAction={onRetry}
      primaryLabel="Check Again"
      secondaryAction={onContinue}
      secondaryLabel="Continue anyway"
      onBack={onBack}
      canGoBack={canGoBack}
    >
      {/* Missing dependencies list */}
      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-4 mb-4">
        {missingDeps.map(dep => (
          <div key={dep.name} className="flex items-start gap-3 py-2 border-b border-gray-200/50 dark:border-gray-800/50 last:border-0">
            <div className="w-5 h-5 rounded-full bg-red-500/20 flex items-center justify-center flex-shrink-0 mt-0.5">
              <svg className="w-3 h-3 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <div>
              <div className="text-sm font-medium text-gray-900 dark:text-white">{dep.name}</div>
              {dep.message && (
                <p className="text-xs text-gray-500 mt-0.5">{dep.message}</p>
              )}
              {dep.helpUrl && (
                <a
                  href={dep.helpUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 text-xs text-blue-500 dark:text-blue-400 hover:text-blue-600 dark:hover:text-blue-300 mt-1"
                >
                  Install instructions
                  <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* Combined install command */}
      {combinedInstallCommand && (
        <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-4">
          <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">Run this command to install everything:</p>
          <div className="flex items-center gap-2">
            <code className="flex-1 text-xs bg-gray-200 dark:bg-gray-900 text-gray-700 dark:text-gray-300 px-3 py-2 rounded font-mono overflow-x-auto">
              {combinedInstallCommand}
            </code>
            <button
              onClick={copyCommand}
              className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 transition-colors p-2"
              title="Copy command"
            >
              {copied ? (
                <svg className="w-5 h-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              )}
            </button>
          </div>
        </div>
      )}
    </PageTemplate>
  );
}

// Cross-Platform Question Page
function CrossPlatformPage({
  dockerDep,
  onYes,
  onSkip,
  onBack,
  canGoBack
}: {
  dockerDep: DependencyStatus | undefined;
  onYes: () => void;
  onSkip: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  // If Docker is installed and image is already built, show ready state
  const isReady = dockerDep?.installed && dockerDep?.imageBuilt === true;

  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
    >
      {isReady ? (
        <>
          {/* Green checkmark for ready state */}
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: 'spring', stiffness: 200, damping: 15 }}
            className="w-20 h-20 rounded-full bg-green-500/20 flex items-center justify-center mb-6"
          >
            <svg className="w-10 h-10 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
            </svg>
          </motion.div>

          {/* Platform icons - smaller for ready state */}
          <div className="flex items-center gap-4 mb-4">
            {/* Windows */}
            <svg className="w-8 h-8 text-gray-600 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
              <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
            </svg>
            {/* Apple */}
            <svg className="w-8 h-8 text-gray-600 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
            </svg>
            {/* Linux */}
            <svg className="w-8 h-8 text-gray-600 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a.424.424 0 00-.11.135c-.26.268-.45.6-.663.839-.199.199-.485.267-.797.4-.313.136-.658.269-.864.68-.09.189-.136.394-.132.602 0 .199.027.4.055.536.058.399.116.728.04.97-.249.68-.28 1.145-.106 1.484.174.334.535.47.94.601.81.2 1.91.135 2.774.6.926.466 1.866.67 2.616.47.526-.116.97-.464 1.208-.946.587-.003 1.23-.269 2.26-.334.699-.058 1.574.267 2.577.2.025.134.063.198.114.333l.003.003c.391.778 1.113 1.132 1.884 1.071.771-.06 1.592-.536 2.257-1.306.631-.765 1.683-1.084 2.378-1.503.348-.199.629-.469.649-.853.023-.4-.2-.811-.714-1.376v-.097l-.003-.003c-.17-.2-.25-.535-.338-.926-.085-.401-.182-.786-.492-1.046h-.003c-.059-.054-.123-.067-.188-.135a.357.357 0 00-.19-.064c.431-1.278.264-2.55-.173-3.694-.533-1.41-1.465-2.638-2.175-3.483-.796-1.005-1.576-1.957-1.56-3.368.026-2.152.236-6.133-3.544-6.139zm.529 3.405h.013c.213 0 .396.062.584.198.19.135.33.332.438.533.105.259.158.459.166.724 0-.02.006-.04.006-.06v.105a.086.086 0 01-.004-.021l-.004-.024a1.807 1.807 0 01-.15.706.953.953 0 01-.213.335.71.71 0 00-.088-.042c-.104-.045-.198-.064-.284-.133a1.312 1.312 0 00-.22-.066c.05-.06.146-.133.183-.198.053-.128.082-.264.088-.402v-.02a1.21 1.21 0 00-.061-.4c-.045-.134-.101-.2-.183-.333-.084-.066-.167-.132-.267-.132h-.016c-.093 0-.176.03-.262.132a.8.8 0 00-.205.334 1.18 1.18 0 00-.09.4v.019c.002.089.008.179.02.267-.193-.067-.438-.135-.607-.202a1.635 1.635 0 01-.018-.2v-.02a1.772 1.772 0 01.15-.768c.082-.22.232-.406.43-.533a.985.985 0 01.594-.2zm-2.962.059h.036c.142 0 .27.048.399.135.146.129.264.288.344.465.09.199.14.4.153.667v.004c.007.134.006.2-.002.266v.08c-.03.007-.056.018-.083.024-.152.055-.274.135-.393.2.012-.09.013-.18.003-.267v-.015c-.012-.133-.04-.2-.082-.333a.613.613 0 00-.166-.267.248.248 0 00-.183-.064h-.021c-.071.006-.13.04-.186.132a.552.552 0 00-.12.27.944.944 0 00-.023.33v.015c.012.135.037.2.08.334.046.134.098.2.166.268.01.009.02.018.034.024-.07.057-.117.07-.176.136a.304.304 0 01-.131.068 2.62 2.62 0 01-.275-.402 1.772 1.772 0 01-.155-.667 1.759 1.759 0 01.08-.668 1.43 1.43 0 01.283-.535c.128-.133.26-.2.418-.2zm1.37 1.706c.332 0 .733.065 1.216.399.293.2.523.269 1.052.468h.003c.255.136.405.266.478.399v-.131a.571.571 0 01.016.47c-.123.31-.516.643-1.063.842v.002c-.268.135-.501.333-.775.465-.276.135-.588.292-1.012.267a1.139 1.139 0 01-.448-.067 3.566 3.566 0 01-.322-.198c-.195-.135-.363-.332-.612-.465v-.005h-.005c-.4-.246-.616-.512-.686-.71-.07-.268-.005-.47.193-.6.224-.135.38-.271.483-.336.104-.074.143-.102.176-.131h.002v-.003c.169-.202.436-.47.839-.601.139-.036.294-.065.466-.065zm2.8 2.142c.358 1.417 1.196 3.475 1.735 4.473.286.534.855 1.659 1.102 3.024.156-.005.33.018.513.064.646-1.671-.546-3.467-1.089-3.966-.22-.2-.232-.335-.123-.335.59.534 1.365 1.572 1.646 2.757.13.535.16 1.104.021 1.67.067.028.135.06.205.067 1.032.534 1.413.938 1.23 1.537v-.002c-.06-.135-.12-.2-.09-.267.046-.134.078-.333-.201-.465-.57-.267-.96-.4-1.18-.535a.98.98 0 01-.36-.4c-.298.533-.648.868-.94 1.002-.04-.2-.021-.4.09-.6a.71.71 0 01.381-.267c.376-.202.559-.47.646-.869.067-.399.024-.733-.135-1.135-.15-.4-.396-.665-.794-.933a2.01 2.01 0 00-.92-.267c-.435-.064-.747.048-.988.135-.075.022-.155.04-.239.054a2.56 2.56 0 01.106-.858c.09-.335.2-.6.323-.868a.262.262 0 01-.09-.134c-.067-.267-.2-.2-.33-.002a1.763 1.763 0 00-.172.535 2.114 2.114 0 00-.038.467c-.065.065-.132.135-.198.199-.257.193-.52.398-.737.601a2.71 2.71 0 01-.18-.202c-.27-.332-.393-.667-.354-1.067a.89.89 0 01.11-.334c.031-.053.067-.067.1-.135a.065.065 0 01.016-.023.09.09 0 01.015-.023v-.003a5.59 5.59 0 01.166-.267c.126-.2.27-.4.461-.602.14-.134.274-.267.41-.4.069-.066.14-.135.21-.2.07-.066.136-.135.203-.2.069-.134.202-.2.37-.266a.33.33 0 00.14-.067c-.12-.067-.137-.2-.061-.336.134-.332.453-.668.785-.933.332-.265.66-.4.875-.4.232.003.325.068.227.403z"/>
            </svg>
          </div>

          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2 text-center">
            Cross-platform builds ready!
          </h2>
          <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-md">
            You can build for Windows, macOS, and Linux from this machine
          </p>

          <div className="flex items-center gap-3">
            {canGoBack && onBack && (
              <button
                onClick={onBack}
                className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
              >
                Back
              </button>
            )}
            <button
              onClick={onSkip}
              className="px-5 py-2 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors"
            >
              Continue
            </button>
          </div>
        </>
      ) : (
        <>
          {/* Platform icons */}
          <div className="flex items-center gap-6 mb-8">
            {/* Windows */}
            <svg className="w-12 h-12 text-gray-600 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
              <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
            </svg>
            {/* Apple */}
            <svg className="w-12 h-12 text-gray-600 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
            </svg>
            {/* Linux */}
            <svg className="w-12 h-12 text-gray-600 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a.424.424 0 00-.11.135c-.26.268-.45.6-.663.839-.199.199-.485.267-.797.4-.313.136-.658.269-.864.68-.09.189-.136.394-.132.602 0 .199.027.4.055.536.058.399.116.728.04.97-.249.68-.28 1.145-.106 1.484.174.334.535.47.94.601.81.2 1.91.135 2.774.6.926.466 1.866.67 2.616.47.526-.116.97-.464 1.208-.946.587-.003 1.23-.269 2.26-.334.699-.058 1.574.267 2.577.2.025.134.063.198.114.333l.003.003c.391.778 1.113 1.132 1.884 1.071.771-.06 1.592-.536 2.257-1.306.631-.765 1.683-1.084 2.378-1.503.348-.199.629-.469.649-.853.023-.4-.2-.811-.714-1.376v-.097l-.003-.003c-.17-.2-.25-.535-.338-.926-.085-.401-.182-.786-.492-1.046h-.003c-.059-.054-.123-.067-.188-.135a.357.357 0 00-.19-.064c.431-1.278.264-2.55-.173-3.694-.533-1.41-1.465-2.638-2.175-3.483-.796-1.005-1.576-1.957-1.56-3.368.026-2.152.236-6.133-3.544-6.139zm.529 3.405h.013c.213 0 .396.062.584.198.19.135.33.332.438.533.105.259.158.459.166.724 0-.02.006-.04.006-.06v.105a.086.086 0 01-.004-.021l-.004-.024a1.807 1.807 0 01-.15.706.953.953 0 01-.213.335.71.71 0 00-.088-.042c-.104-.045-.198-.064-.284-.133a1.312 1.312 0 00-.22-.066c.05-.06.146-.133.183-.198.053-.128.082-.264.088-.402v-.02a1.21 1.21 0 00-.061-.4c-.045-.134-.101-.2-.183-.333-.084-.066-.167-.132-.267-.132h-.016c-.093 0-.176.03-.262.132a.8.8 0 00-.205.334 1.18 1.18 0 00-.09.4v.019c.002.089.008.179.02.267-.193-.067-.438-.135-.607-.202a1.635 1.635 0 01-.018-.2v-.02a1.772 1.772 0 01.15-.768c.082-.22.232-.406.43-.533a.985.985 0 01.594-.2zm-2.962.059h.036c.142 0 .27.048.399.135.146.129.264.288.344.465.09.199.14.4.153.667v.004c.007.134.006.2-.002.266v.08c-.03.007-.056.018-.083.024-.152.055-.274.135-.393.2.012-.09.013-.18.003-.267v-.015c-.012-.133-.04-.2-.082-.333a.613.613 0 00-.166-.267.248.248 0 00-.183-.064h-.021c-.071.006-.13.04-.186.132a.552.552 0 00-.12.27.944.944 0 00-.023.33v.015c.012.135.037.2.08.334.046.134.098.2.166.268.01.009.02.018.034.024-.07.057-.117.07-.176.136a.304.304 0 01-.131.068 2.62 2.62 0 01-.275-.402 1.772 1.772 0 01-.155-.667 1.759 1.759 0 01.08-.668 1.43 1.43 0 01.283-.535c.128-.133.26-.2.418-.2zm1.37 1.706c.332 0 .733.065 1.216.399.293.2.523.269 1.052.468h.003c.255.136.405.266.478.399v-.131a.571.571 0 01.016.47c-.123.31-.516.643-1.063.842v.002c-.268.135-.501.333-.775.465-.276.135-.588.292-1.012.267a1.139 1.139 0 01-.448-.067 3.566 3.566 0 01-.322-.198c-.195-.135-.363-.332-.612-.465v-.005h-.005c-.4-.246-.616-.512-.686-.71-.07-.268-.005-.47.193-.6.224-.135.38-.271.483-.336.104-.074.143-.102.176-.131h.002v-.003c.169-.202.436-.47.839-.601.139-.036.294-.065.466-.065zm2.8 2.142c.358 1.417 1.196 3.475 1.735 4.473.286.534.855 1.659 1.102 3.024.156-.005.33.018.513.064.646-1.671-.546-3.467-1.089-3.966-.22-.2-.232-.335-.123-.335.59.534 1.365 1.572 1.646 2.757.13.535.16 1.104.021 1.67.067.028.135.06.205.067 1.032.534 1.413.938 1.23 1.537v-.002c-.06-.135-.12-.2-.09-.267.046-.134.078-.333-.201-.465-.57-.267-.96-.4-1.18-.535a.98.98 0 01-.36-.4c-.298.533-.648.868-.94 1.002-.04-.2-.021-.4.09-.6a.71.71 0 01.381-.267c.376-.202.559-.47.646-.869.067-.399.024-.733-.135-1.135-.15-.4-.396-.665-.794-.933a2.01 2.01 0 00-.92-.267c-.435-.064-.747.048-.988.135-.075.022-.155.04-.239.054a2.56 2.56 0 01.106-.858c.09-.335.2-.6.323-.868a.262.262 0 01-.09-.134c-.067-.267-.2-.2-.33-.002a1.763 1.763 0 00-.172.535 2.114 2.114 0 00-.038.467c-.065.065-.132.135-.198.199-.257.193-.52.398-.737.601a2.71 2.71 0 01-.18-.202c-.27-.332-.393-.667-.354-1.067a.89.89 0 01.11-.334c.031-.053.067-.067.1-.135a.065.065 0 01.016-.023.09.09 0 01.015-.023v-.003a5.59 5.59 0 01.166-.267c.126-.2.27-.4.461-.602.14-.134.274-.267.41-.4.069-.066.14-.135.21-.2.07-.066.136-.135.203-.2.069-.134.202-.2.37-.266a.33.33 0 00.14-.067c-.12-.067-.137-.2-.061-.336.134-.332.453-.668.785-.933.332-.265.66-.4.875-.4.232.003.325.068.227.403z"/>
            </svg>
          </div>

          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2 text-center">
            Build for multiple platforms?
          </h2>
          <p className="text-gray-500 dark:text-gray-400 mb-2 text-center max-w-md">
            Wails can compile your app for Windows, macOS, and Linux from a single machine
          </p>
          <p className="text-xs text-gray-400 dark:text-gray-500 mb-8 text-center">
            Requires Docker for cross-compilation
          </p>

          <div className="flex flex-col items-center gap-2">
            <div className="flex items-center gap-3">
              {canGoBack && onBack && (
                <button
                  onClick={onBack}
                  className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
                >
                  Back
                </button>
              )}
              <button
                onClick={onYes}
                className="px-5 py-2 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors"
              >
                Yes, set this up
              </button>
            </div>
            <button
              onClick={onSkip}
              className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
            >
              Not right now
            </button>
          </div>
        </>
      )}
    </motion.div>
  );
}

// Docker Options Page - choose between official download or build your own
function DockerOptionsPage({
  onDownloadOfficial,
  onSkip,
  onBack,
  canGoBack
}: {
  onDownloadOfficial: () => void;
  onSkip: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
    >
      <div className="w-16 h-16 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6">
        <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
          <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27ZM189.67,206.39h-81.7v81.7h81.7v-81.7ZM295.22,206.39h-81.7v81.7h81.7v-81.7ZM400.77,206.39h-81.7v81.7h81.7v-81.7ZM506.32,206.39h-81.7v81.7h81.7v-81.7ZM84.12,206.39H2.42v81.7h81.7v-81.7ZM189.67,103.2h-81.7v81.7h81.7v-81.7ZM295.22,103.2h-81.7v81.7h81.7v-81.7ZM400.77,103.2h-81.7v81.7h81.7v-81.7ZM400.77,0h-81.7v81.7h81.7V0Z"/>
        </svg>
      </div>

      <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
        Set up cross-compiler
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-sm">
        Choose how to get the Docker cross-compilation image
      </p>

      <div className="flex flex-col gap-3 mb-6 w-full max-w-xs">
        <button
          onClick={onDownloadOfficial}
          className="w-full px-5 py-3 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors flex items-center justify-center gap-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          Download official image
        </button>

        <a
          href="https://wails.io/docs/guides/build/cross-platform#build-your-own-image"
          target="_blank"
          rel="noopener noreferrer"
          className="w-full px-5 py-3 rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 text-sm font-medium hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors flex items-center justify-center gap-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
          </svg>
          Build your own image
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
          </svg>
        </a>
      </div>

      <div className="flex flex-col items-center gap-2">
        {canGoBack && onBack && (
          <button
            onClick={onBack}
            className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
          >
            Back
          </button>
        )}
        <button
          onClick={onSkip}
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
        >
          Skip for now
        </button>
      </div>
    </motion.div>
  );
}

function DockerBuildError({
  onBuildImage,
  onSkip
}: {
  onBuildImage: () => void;
  onSkip: () => void;
}) {
  const [showLogs, setShowLogs] = useState(false);
  const [logs, setLogs] = useState<string | null>(null);

  const handleViewLogs = async () => {
    if (!logs) {
      const response = await fetch('/api/docker/logs');
      const text = await response.text();
      setLogs(text);
    }
    setShowLogs(true);
  };

  if (showLogs) {
    return (
      <motion.div
        variants={pageVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ duration: 0.3 }}
        className="flex-1 flex flex-col items-center justify-center max-w-4xl mx-auto w-full"
      >
        <div className="w-full flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Build Logs
          </h2>
          <button
            onClick={() => setShowLogs(false)}
            className="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
          >
            Back
          </button>
        </div>
        <pre className="w-full h-96 overflow-auto bg-gray-900 text-gray-100 p-4 rounded-lg text-xs font-mono whitespace-pre-wrap">
          {logs || 'No logs available'}
        </pre>
      </motion.div>
    );
  }

  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
    >
      <div className="w-16 h-16 rounded-2xl bg-amber-500/20 flex items-center justify-center mb-6">
        <svg className="w-8 h-8 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>

      <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
        Build failed
      </h2>

      <p className="text-sm text-gray-500 dark:text-gray-400 mb-2 text-center max-w-sm">
        Check your internet connection and try again, or download the SDK manually.
      </p>

      <button
        onClick={handleViewLogs}
        className="text-sm text-blue-500 hover:text-blue-600 mb-6"
      >
        View logs
      </button>

      <div className="flex flex-col gap-3 items-center">
        <div className="flex gap-3">
          <button
            onClick={onBuildImage}
            className="px-5 py-2.5 rounded-lg bg-blue-500 text-white text-sm font-medium hover:bg-blue-600 transition-colors"
          >
            Try again
          </button>
          <a
            href="https://wails.io/docs/guides/build/cross-platform#build-your-own-image"
            target="_blank"
            rel="noopener noreferrer"
            className="px-5 py-2.5 rounded-lg border border-blue-500 text-blue-600 dark:text-blue-400 text-sm font-medium hover:bg-blue-500/10 transition-colors"
          >
            Build your own
          </a>
        </div>
        <button
          onClick={onSkip}
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
        >
          Skip for now
        </button>
      </div>
    </motion.div>
  );
}

function SDKLicensePage({
  onAgree,
  onDecline,
  onBack,
  canGoBack
}: {
  onAgree: () => void;
  onDecline: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  const [agreed, setAgreed] = useState(false);

  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
    >
      <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-1">
        Apple SDK License Agreement
      </h2>

      <p className="text-sm text-gray-500 dark:text-gray-400 mb-4 text-center max-w-md">
        Cross-platform builds for macOS require the Apple SDK. Please review and accept the license terms.
      </p>

      <div className="w-full max-w-2xl h-72 mb-4 rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700 bg-white">
        <iframe
          src="/assets/apple-sdk-license.pdf#view=FitH&navpanes=0&toolbar=0"
          className="w-full h-full"
          title="Apple Xcode and SDK License Agreement"
        />
      </div>

      <label className="flex items-center gap-2 mb-5 cursor-pointer">
        <input
          type="checkbox"
          checked={agreed}
          onChange={(e) => setAgreed(e.target.checked)}
          className="w-4 h-4 rounded border-gray-300 text-blue-500 focus:ring-blue-500"
        />
        <span className="text-sm text-gray-600 dark:text-gray-300">
          I agree to Apple's Xcode and SDK License Agreement
        </span>
      </label>

      <div className="flex flex-col items-center gap-2">
        <div className="flex gap-3">
          {canGoBack && onBack && (
            <button
              onClick={onBack}
              className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              Back
            </button>
          )}
          <button
            onClick={onAgree}
            disabled={!agreed}
            className={`px-5 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              agreed
                ? 'bg-red-500 text-white hover:bg-red-600'
                : 'bg-gray-200 text-gray-400 cursor-not-allowed'
            }`}
          >
            Continue
          </button>
        </div>
        <button
          onClick={onDecline}
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
        >
          Skip for now
        </button>
      </div>
    </motion.div>
  );
}

function DockerSetupPage({
  dockerStatus,
  buildingImage,
  onBuildImage,
  onCheckAgain,
  onContinueBackground,
  onSkip,
  onBack,
  canGoBack
}: {
  dockerStatus: DockerStatus | null;
  buildingImage: boolean;
  onBuildImage: () => void;
  onCheckAgain: () => void;
  onContinueBackground: () => void;
  onSkip: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  // Docker not installed
  if (!dockerStatus || !dockerStatus.installed) {
    return (
      <motion.div
        variants={pageVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ duration: 0.3 }}
        className="flex-1 flex flex-col items-center justify-center"
      >
        <div className="w-16 h-16 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6">
          <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
            <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27ZM189.67,206.39h-81.7v81.7h81.7v-81.7ZM295.22,206.39h-81.7v81.7h81.7v-81.7ZM400.77,206.39h-81.7v81.7h81.7v-81.7ZM506.32,206.39h-81.7v81.7h81.7v-81.7ZM84.12,206.39H2.42v81.7h81.7v-81.7ZM189.67,103.2h-81.7v81.7h81.7v-81.7ZM295.22,103.2h-81.7v81.7h81.7v-81.7ZM400.77,103.2h-81.7v81.7h81.7v-81.7ZM400.77,0h-81.7v81.7h81.7V0Z"/>
          </svg>
        </div>

        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Install Docker
        </h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 text-center max-w-sm">
          Cross-platform builds require Docker Desktop
        </p>

        <a
          href="https://docs.docker.com/get-docker/"
          target="_blank"
          rel="noopener noreferrer"
          className="px-5 py-2 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-500 transition-colors inline-flex items-center gap-2 mb-4"
        >
          Download Docker Desktop
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
          </svg>
        </a>

        <p className="text-xs text-gray-400 dark:text-gray-500 mb-6 text-center max-w-xs">
          After installing, come back and we'll continue setting up.
          Some platforms may require a reboot.
        </p>

        <div className="flex flex-col items-center gap-1.5">
          <div className="flex items-center gap-3">
            {canGoBack && onBack && (
              <button
                onClick={onBack}
                className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
              >
                Back
              </button>
            )}
            <button
              onClick={onCheckAgain}
              className="px-5 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 text-sm font-medium hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
            >
              Check Again
            </button>
          </div>
          <button
            onClick={onSkip}
            className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
          >
            Skip for now
          </button>
        </div>
      </motion.div>
    );
  }

  // Docker not running
  if (!dockerStatus.running) {
    return (
      <motion.div
        variants={pageVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ duration: 0.3 }}
        className="flex-1 flex flex-col items-center justify-center"
      >
        <div className="w-16 h-16 rounded-2xl bg-gray-200 dark:bg-gray-800 flex items-center justify-center mb-6 opacity-50">
          <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
            <path fill="#6b7280" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27ZM189.67,206.39h-81.7v81.7h81.7v-81.7ZM295.22,206.39h-81.7v81.7h81.7v-81.7ZM400.77,206.39h-81.7v81.7h81.7v-81.7ZM506.32,206.39h-81.7v81.7h81.7v-81.7ZM84.12,206.39H2.42v81.7h81.7v-81.7ZM189.67,103.2h-81.7v81.7h81.7v-81.7ZM295.22,103.2h-81.7v81.7h81.7v-81.7ZM400.77,103.2h-81.7v81.7h81.7v-81.7ZM400.77,0h-81.7v81.7h81.7V0Z"/>
          </svg>
        </div>

        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Start Docker
        </h2>
        <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-sm">
          Please start Docker Desktop to continue
        </p>

        <div className="flex flex-col items-center gap-1.5">
          <div className="flex items-center gap-3">
            {canGoBack && onBack && (
              <button
                onClick={onBack}
                className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
              >
                Back
              </button>
            )}
            <button
              onClick={onCheckAgain}
              className="px-5 py-2 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors"
            >
              Check Again
            </button>
          </div>
          <button
            onClick={onSkip}
            className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
          >
            Skip for now
          </button>
        </div>
      </motion.div>
    );
  }

  // Building image
  if (buildingImage || dockerStatus.pullStatus === 'pulling') {
    const progress = dockerStatus.pullProgress || 0;
    const stage = dockerStatus.pullMessage || 'Connecting';
    return (
      <motion.div
        variants={pageVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ duration: 0.3 }}
        className="flex-1 flex flex-col items-center justify-center"
      >
        <div className="w-16 h-16 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6">
          <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
            <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27ZM189.67,206.39h-81.7v81.7h81.7v-81.7ZM295.22,206.39h-81.7v81.7h81.7v-81.7ZM400.77,206.39h-81.7v81.7h81.7v-81.7ZM506.32,206.39h-81.7v81.7h81.7v-81.7ZM84.12,206.39H2.42v81.7h81.7v-81.7ZM189.67,103.2h-81.7v81.7h81.7v-81.7ZM295.22,103.2h-81.7v81.7h81.7v-81.7ZM400.77,103.2h-81.7v81.7h81.7v-81.7ZM400.77,0h-81.7v81.7h81.7V0Z"/>
          </svg>
        </div>

        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Downloading cross-compiler image
        </h2>

        <div className="w-64 mb-4">
          <div className="flex items-center justify-between text-sm text-gray-500 mb-1">
            <span>{stage}</span>
            <span>{progress}%</span>
          </div>
          <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <motion.div
              className="h-full bg-blue-500"
              animate={{ width: `${progress}%` }}
            />
          </div>
        </div>

        <p className="text-xs text-gray-400 dark:text-gray-500 mb-8 text-center">
          This may take several minutes
        </p>

        <button
          onClick={onContinueBackground}
          className="px-5 py-2 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors"
        >
          Continue in background
        </button>
      </motion.div>
    );
  }

  if (dockerStatus.pullStatus === 'error') {
    return (
      <DockerBuildError
        onBuildImage={onBuildImage}
        onSkip={onSkip}
      />
    );
  }

  // Image already built
  if (dockerStatus.imageBuilt) {
    return (
      <motion.div
        variants={pageVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ duration: 0.3 }}
        className="flex-1 flex flex-col items-center justify-center"
      >
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: 'spring', stiffness: 200, damping: 15 }}
          className="w-16 h-16 rounded-2xl bg-green-500/20 flex items-center justify-center mb-6"
        >
          <svg className="w-8 h-8 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
          </svg>
        </motion.div>

        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Cross-platform builds ready!
        </h2>
        <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-sm">
          You can now build for Windows, macOS, and Linux
        </p>

        <div className="flex items-center gap-3">
          {canGoBack && onBack && (
            <button
              onClick={onBack}
              className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              Back
            </button>
          )}
          <button
            onClick={onContinueBackground}
            className="px-5 py-2 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors"
          >
            Continue
          </button>
        </div>
      </motion.div>
    );
  }

  // Docker ready, image not built yet - prompt to build
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
    >
      <div className="w-16 h-16 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6">
        <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
          <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27ZM189.67,206.39h-81.7v81.7h81.7v-81.7ZM295.22,206.39h-81.7v81.7h81.7v-81.7ZM400.77,206.39h-81.7v81.7h81.7v-81.7ZM506.32,206.39h-81.7v81.7h81.7v-81.7ZM84.12,206.39H2.42v81.7h81.7v-81.7ZM189.67,103.2h-81.7v81.7h81.7v-81.7ZM295.22,103.2h-81.7v81.7h81.7v-81.7ZM400.77,103.2h-81.7v81.7h81.7v-81.7ZM400.77,0h-81.7v81.7h81.7V0Z"/>
        </svg>
      </div>

      <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
        Docker is ready!
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-2 text-center max-w-sm">
        Download the cross-compilation image to enable building for all platforms
      </p>
      <p className="text-xs text-gray-400 dark:text-gray-500 mb-8 text-center">
        This will download ~800MB and may take several minutes
      </p>

      <div className="flex flex-col items-center gap-2">
        <div className="flex items-center gap-3">
          {canGoBack && onBack && (
            <button
              onClick={onBack}
              className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              Back
            </button>
          )}
          <button
            onClick={onBuildImage}
            className="px-5 py-2 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-500 transition-colors"
          >
            Download Image
          </button>
        </div>
        <button
          onClick={onSkip}
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
        >
          Skip, I'll do it later
        </button>
      </div>
    </motion.div>
  );
}

// Available framework templates
const FRAMEWORKS: FrameworkTemplate[] = [
  { id: 'vanilla', name: 'Vanilla', description: 'Plain JavaScript/TypeScript', color: '#f7df1e', icon: 'javascript' },
  { id: 'react', name: 'React', description: 'React with Vite', color: '#61dafb', icon: 'react' },
  { id: 'vue', name: 'Vue', description: 'Vue 3 with Vite', color: '#42b883', icon: 'vue' },
  { id: 'svelte', name: 'Svelte', description: 'Svelte with Vite', color: '#ff3e00', icon: 'svelte' },
  { id: 'preact', name: 'Preact', description: 'Lightweight React alternative', color: '#673ab8', icon: 'preact' },
  { id: 'lit', name: 'Lit', description: 'Web Components with Lit', color: '#324fff', icon: 'lit' },
  { id: 'solid', name: 'Solid', description: 'Solid.js with Vite', color: '#2c4f7c', icon: 'solid' },
  { id: 'qwik', name: 'Qwik', description: 'Qwik with Vite', color: '#18b6f6', icon: 'qwik' },
];

function LanguageSelectPage({
  preferTypeScript,
  onSelect,
  onNext,
  onBack,
  canGoBack
}: {
  preferTypeScript: boolean;
  onSelect: (useTypeScript: boolean) => void;
  onNext: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, []);

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center"
      aria-labelledby="language-title"
    >
      <h2
        ref={headingRef}
        id="language-title"
        className="text-2xl font-semibold text-gray-900 dark:text-white mb-2 text-center focus:outline-none"
        tabIndex={-1}
      >
        Language Preference
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-md">
        Choose your preferred language for new projects
      </p>

      <div className="flex gap-4 mb-8" role="radiogroup" aria-label="Programming language">
        <button
          onClick={() => onSelect(false)}
          role="radio"
          aria-checked={!preferTypeScript}
          className={`w-40 h-48 rounded-xl p-5 flex flex-col items-center justify-center gap-3 transition-all border-2 focus:outline-none focus:ring-2 focus:ring-yellow-400 focus:ring-offset-2 dark:focus:ring-offset-gray-900 ${
            !preferTypeScript
              ? 'border-yellow-400 bg-yellow-400/10 shadow-lg shadow-yellow-400/20'
              : 'border-gray-200 dark:border-white/10 bg-gray-100 dark:bg-white/5 hover:bg-gray-200 dark:hover:bg-white/10'
          }`}
        >
          <div className="w-16 h-16 flex items-center justify-center" aria-hidden="true">
            <img src="/logos/javascript.svg" alt="" className="w-14 h-14" />
          </div>
          <span className="text-lg font-semibold text-gray-900 dark:text-white">JavaScript</span>
          <span className="text-xs text-gray-500 dark:text-white/50">Dynamic typing</span>
        </button>

        <button
          onClick={() => onSelect(true)}
          role="radio"
          aria-checked={preferTypeScript}
          className={`w-40 h-48 rounded-xl p-5 flex flex-col items-center justify-center gap-3 transition-all border-2 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 dark:focus:ring-offset-gray-900 ${
            preferTypeScript
              ? 'border-blue-400 bg-blue-400/10 shadow-lg shadow-blue-400/20'
              : 'border-gray-200 dark:border-white/10 bg-gray-100 dark:bg-white/5 hover:bg-gray-200 dark:hover:bg-white/10'
          }`}
        >
          <div className="w-16 h-16 flex items-center justify-center" aria-hidden="true">
            <img src="/logos/typescript.svg" alt="" className="w-14 h-14" />
          </div>
          <span className="text-lg font-semibold text-gray-900 dark:text-white">TypeScript</span>
          <span className="text-xs text-gray-500 dark:text-white/50">Type safety</span>
        </button>
      </div>

      <div className="flex items-center gap-3" role="group" aria-label="Navigation">
        {canGoBack && onBack && (
          <button
            onClick={onBack}
            className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
          >
            Back
          </button>
        )}
        <button
          onClick={onNext}
          className="px-6 py-2.5 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
        >
          Continue
        </button>
      </div>
    </motion.main>
  );
}

function BindingStylePage({
  useInterfaces,
  onSelect,
  onNext,
  onBack,
  canGoBack
}: {
  useInterfaces: boolean;
  onSelect: (useInterfaces: boolean) => void;
  onNext: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, []);

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col items-center justify-center px-4 overflow-hidden"
      aria-labelledby="binding-title"
    >
      <h2
        ref={headingRef}
        id="binding-title"
        className="text-2xl font-semibold text-gray-900 dark:text-white mb-2 text-center focus:outline-none"
        tabIndex={-1}
      >
        TypeScript Binding Style
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-6 text-center max-w-lg">
        Choose how Go structs are represented in TypeScript
      </p>

      <div className="flex gap-4 mb-8 max-w-full overflow-x-auto p-1" role="radiogroup" aria-label="Binding style">
        <button
          onClick={() => onSelect(true)}
          role="radio"
          aria-checked={useInterfaces}
          className={`w-56 shrink-0 rounded-xl p-4 flex flex-col items-start gap-2 transition-all border-2 text-left focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 dark:focus:ring-offset-gray-900 ${
            useInterfaces
              ? 'border-blue-400 bg-blue-400/10 shadow-lg shadow-blue-400/20'
              : 'border-white/10 bg-white/5 hover:bg-white/10'
          }`}
        >
          <span className="text-base font-semibold text-gray-900 dark:text-white">Interfaces</span>
          <pre className="text-[10px] leading-tight text-gray-700 dark:text-white/70 font-mono bg-gray-100 dark:bg-black/30 p-2 rounded-lg w-full overflow-x-auto" aria-hidden="true">
{`interface Person {
  name: string;
  age: number;
}`}
          </pre>
          <ul className="text-[10px] text-gray-500 dark:text-white/50 space-y-0.5" aria-label="Features">
            <li>Lightweight types</li>
            <li>No runtime code</li>
            <li>Simpler output</li>
          </ul>
        </button>

        <button
          onClick={() => onSelect(false)}
          role="radio"
          aria-checked={!useInterfaces}
          className={`w-56 shrink-0 rounded-xl p-4 flex flex-col items-start gap-2 transition-all border-2 text-left focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 dark:focus:ring-offset-gray-900 ${
            !useInterfaces
              ? 'border-purple-400 bg-purple-400/10 shadow-lg shadow-purple-400/20'
              : 'border-white/10 bg-white/5 hover:bg-white/10'
          }`}
        >
          <span className="text-base font-semibold text-gray-900 dark:text-white">Classes</span>
          <pre className="text-[10px] leading-tight text-gray-700 dark:text-white/70 font-mono bg-gray-100 dark:bg-black/30 p-2 rounded-lg w-full overflow-x-auto" aria-hidden="true">
{`class Person {
  name: string;
  age: number;
  constructor(src) {
    Object.assign(this, src);
  }
  static createFrom(src) {
    return new Person(src);
  }
}`}
          </pre>
          <ul className="text-[10px] text-gray-500 dark:text-white/50 space-y-0.5" aria-label="Features">
            <li>Factory methods</li>
            <li>Default initialization</li>
            <li>More verbose</li>
          </ul>
        </button>
      </div>

      <div className="flex items-center gap-3" role="group" aria-label="Navigation">
        {canGoBack && onBack && (
          <button
            onClick={onBack}
            className="px-4 py-2 rounded-lg text-sm font-medium transition-colors border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
          >
            Back
          </button>
        )}
        <button
          onClick={onNext}
          className="px-6 py-2.5 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
        >
          Continue
        </button>
      </div>
    </motion.main>
  );
}

function TemplateSelectPage({
  selectedFramework,
  preferTypeScript,
  onSelect,
  onNext,
  onSkip,
  onBack,
  canGoBack
}: {
  selectedFramework: string;
  preferTypeScript: boolean;
  onSelect: (frameworkId: string) => void;
  onNext: () => void;
  onSkip: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  return (
    <PageTemplate
      title="Default Template"
      subtitle="Choose a framework for new projects"
      primaryAction={onNext}
      primaryLabel="Continue"
      secondaryAction={onSkip}
      secondaryLabel="Skip"
      onBack={onBack}
      canGoBack={canGoBack}
    >
      <div
        className="grid grid-cols-4 gap-3 max-w-2xl mx-auto p-1"
        role="radiogroup"
        aria-label="Framework templates"
      >
        {FRAMEWORKS.map((framework) => (
          <button
            key={framework.id}
            onClick={() => onSelect(framework.id)}
            role="radio"
            aria-checked={selectedFramework === framework.id}
            className={`aspect-square rounded-xl p-4 flex flex-col items-center justify-center gap-2 transition-all border-2 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 ${
              selectedFramework === framework.id
                ? 'border-red-500 bg-red-500/10 shadow-lg shadow-red-500/10'
                : 'border-gray-200 dark:border-white/10 bg-gray-100 dark:bg-white/5 hover:bg-gray-200 dark:hover:bg-white/10'
            }`}
          >
            <img
              src={`/logos/${framework.id === 'vanilla' ? (preferTypeScript ? 'typescript' : 'javascript') : framework.icon}.svg`}
              alt=""
              aria-hidden="true"
              className="w-12 h-12"
            />
            <span className="text-sm font-medium text-gray-900 dark:text-white">{framework.name}</span>
          </button>
        ))}
      </div>
    </PageTemplate>
  );
}

function ProjectsPage({
  defaults,
  onDefaultsChange,
  onNext,
  onSkip,
  saving,
  onBack,
  canGoBack
}: {
  defaults: GlobalDefaults;
  onDefaultsChange: (defaults: GlobalDefaults) => void;
  onNext: () => void;
  onSkip: () => void;
  saving: boolean;
  onBack?: () => void;
  canGoBack?: boolean;
}) {
  const [editingField, setEditingField] = useState<'name' | 'company' | 'bundleId' | null>(null);
  const [tempValue, setTempValue] = useState('');

  const handleRowClick = (field: 'name' | 'company' | 'bundleId') => {
    if (field === 'name') {
      setTempValue(defaults.author.name);
    } else if (field === 'company') {
      setTempValue(defaults.author.company);
    } else if (field === 'bundleId') {
      setTempValue(defaults.project.productIdentifierPrefix);
    }
    setEditingField(field);
  };

  const handleSaveField = () => {
    if (editingField === 'name') {
      onDefaultsChange({
        ...defaults,
        author: { ...defaults.author, name: tempValue }
      });
    } else if (editingField === 'company') {
      onDefaultsChange({
        ...defaults,
        author: { ...defaults.author, company: tempValue }
      });
    } else if (editingField === 'bundleId') {
      onDefaultsChange({
        ...defaults,
        project: { ...defaults.project, productIdentifierPrefix: tempValue }
      });
    }
    setEditingField(null);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSaveField();
    } else if (e.key === 'Escape') {
      setEditingField(null);
    }
  };

  return (
    <PageTemplate
      title="Project Defaults"
      subtitle="Set defaults for new Wails projects"
      primaryAction={onNext}
      primaryLabel={saving ? "Saving..." : "Continue"}
      primaryDisabled={saving}
      secondaryAction={onSkip}
      secondaryLabel="Skip"
      onBack={onBack}
      canGoBack={canGoBack}
    >
      <div className="max-w-xl mx-auto">
        <div className="settings-group" role="group" aria-label="Project default settings">
          {editingField === 'name' ? (
            <div className="settings-row">
              <label htmlFor="author-input" className="sr-only">Author name</label>
              <div className="flex-1">
                <input
                  id="author-input"
                  type="text"
                  value={tempValue}
                  onChange={(e) => setTempValue(e.target.value)}
                  onKeyDown={handleKeyDown}
                  onBlur={handleSaveField}
                  autoFocus
                  placeholder="Your Name"
                  aria-label="Author name"
                  className="w-full bg-transparent border-none text-sm text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 rounded px-1"
                />
              </div>
            </div>
          ) : (
            <button
              className="settings-row w-full text-left focus:outline-none focus:ring-2 focus:ring-inset focus:ring-red-500"
              onClick={() => handleRowClick('name')}
              aria-label={`Author: ${defaults.author.name || 'Not set'}. Click to edit.`}
            >
              <span className="text-sm font-medium text-gray-800 dark:text-white/90">Author</span>
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-white/65">
                <span>{defaults.author.name || 'Not set'}</span>
                <span className="text-gray-400 dark:text-white/40 text-xs" aria-hidden="true">&#9656;</span>
              </div>
            </button>
          )}

          {editingField === 'company' ? (
            <div className="settings-row">
              <label htmlFor="company-input" className="sr-only">Company name</label>
              <div className="flex-1">
                <input
                  id="company-input"
                  type="text"
                  value={tempValue}
                  onChange={(e) => setTempValue(e.target.value)}
                  onKeyDown={handleKeyDown}
                  onBlur={handleSaveField}
                  autoFocus
                  placeholder="Acme Corp"
                  aria-label="Company name"
                  className="w-full bg-transparent border-none text-sm text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 rounded px-1"
                />
              </div>
            </div>
          ) : (
            <button
              className="settings-row w-full text-left focus:outline-none focus:ring-2 focus:ring-inset focus:ring-red-500"
              onClick={() => handleRowClick('company')}
              aria-label={`Company: ${defaults.author.company || 'Not set'}. Click to edit.`}
            >
              <span className="text-sm font-medium text-gray-800 dark:text-white/90">Company</span>
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-white/65">
                <span>{defaults.author.company || 'Not set'}</span>
                <span className="text-gray-400 dark:text-white/40 text-xs" aria-hidden="true">&#9656;</span>
              </div>
            </button>
          )}

          {editingField === 'bundleId' ? (
            <div className="settings-row">
              <label htmlFor="bundle-input" className="sr-only">Bundle identifier</label>
              <div className="flex-1">
                <input
                  id="bundle-input"
                  type="text"
                  value={tempValue}
                  onChange={(e) => setTempValue(e.target.value)}
                  onKeyDown={handleKeyDown}
                  onBlur={handleSaveField}
                  autoFocus
                  placeholder="com.example"
                  aria-label="Bundle identifier"
                  className="w-full bg-transparent border-none text-sm text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 rounded px-1 font-mono"
                />
              </div>
            </div>
          ) : (
            <button
              className="settings-row w-full text-left focus:outline-none focus:ring-2 focus:ring-inset focus:ring-red-500"
              onClick={() => handleRowClick('bundleId')}
              aria-label={`Bundle identifier: ${defaults.project.productIdentifierPrefix || 'com.example'}. Click to edit.`}
            >
              <span className="text-sm font-medium text-gray-800 dark:text-white/90">Bundle identifier</span>
              <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-white/65">
                <span className="font-mono">{defaults.project.productIdentifierPrefix || 'com.example'}</span>
                <span className="text-gray-400 dark:text-white/40 text-xs" aria-hidden="true">&#9656;</span>
              </div>
            </button>
          )}
        </div>
        <p className="text-xs text-gray-500 dark:text-white/40 mt-3 text-center" id="settings-description">
          These defaults are used when creating new projects
        </p>
      </div>
    </PageTemplate>
  );
}

function CompletePage() {
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
  }, []);

  const handleStartBuilding = () => {
    window.open('https://v3alpha.wails.io/quick-start/first-app/', '_blank', 'noopener,noreferrer');
  };

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      aria-labelledby="complete-title"
      className="flex-1 flex flex-col items-center justify-center px-8"
    >
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', stiffness: 200, damping: 15 }}
        className="w-16 h-16 rounded-full bg-green-500/20 flex items-center justify-center mb-4"
        aria-hidden="true"
      >
        <svg className="w-8 h-8 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
        </svg>
      </motion.div>
      <h2
        ref={headingRef}
        id="complete-title"
        className="text-xl font-semibold text-gray-900 dark:text-white mb-6 focus:outline-none"
        tabIndex={-1}
      >
        You're ready to build!
      </h2>

      <button
        onClick={handleStartBuilding}
        className="px-5 py-2 rounded-lg border border-red-500 text-red-600 dark:text-red-400 text-sm font-medium hover:bg-red-500/10 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
      >
        Start Building
      </button>
    </motion.main>
  );
}

// Main App
export default function App() {
  const [step, setStep] = useState<OOBEStep>('splash');
  const [stepHistory, setStepHistory] = useState<OOBEStep[]>([]);
  const [dependencies, setDependencies] = useState<DependencyStatus[]>([]);
  const [_system, setSystem] = useState<SystemInfo | null>(null);
  const [dockerStatus, setDockerStatus] = useState<DockerStatus | null>(null);
  const [buildingImage, setBuildingImage] = useState(false);
  const [defaults, setDefaults] = useState<GlobalDefaults>({
    author: { name: '', company: '' },
    project: {
      productIdentifierPrefix: 'com.example',
      defaultTemplate: 'vanilla',
      copyrightTemplate: '(c) {year}, {company}',
      descriptionTemplate: 'A {name} application',
      defaultVersion: '0.1.0',
      useInterfaces: true
    }
  });
  const [savingDefaults, setSavingDefaults] = useState(false);
  const [backgroundDockerStarted, setBackgroundDockerStarted] = useState(false);
  const [preferTypeScript, setPreferTypeScript] = useState(true);
  const [selectedFramework, setSelectedFramework] = useState('vanilla');
  const [useInterfaces, setUseInterfaces] = useState(true);
  const [showDockerToast, setShowDockerToast] = useState(false);
  const [prevDockerPullStatus, setPrevDockerPullStatus] = useState<string | null>(null);
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('wails-setup-theme');
      if (saved === 'light' || saved === 'dark') return saved;
      if (window.matchMedia('(prefers-color-scheme: light)').matches) return 'light';
    }
    return 'dark';
  });

  const navigateTo = (newStep: OOBEStep) => {
    setStepHistory(prev => [...prev, step]);
    setStep(newStep);
  };

  const goBack = () => {
    if (stepHistory.length === 0) return;
    const newHistory = [...stepHistory];
    let previousStep = newHistory.pop()!;
    
    // Skip transient steps like 'checking'
    while (previousStep === 'checking' && newHistory.length > 0) {
      previousStep = newHistory.pop()!;
    }
    
    setStepHistory(newHistory);
    setStep(previousStep);
  };

  const canGoBack = stepHistory.length > 0 && step !== 'splash' && step !== 'checking';

  const toggleTheme = () => {
    setTheme(prev => {
      const next = prev === 'dark' ? 'light' : 'dark';
      localStorage.setItem('wails-setup-theme', next);
      return next;
    });
  };

  // Apply theme class to document
  useEffect(() => {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme]);

  useEffect(() => {
    init();
  }, []);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
        return;
      }

      if (e.key === 'Escape' && canGoBack) {
        e.preventDefault();
        goBack();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [canGoBack, step]);

  const init = async () => {
    const state = await getState();
    setSystem(state.system);
  };

  const handleSplashNext = async () => {
    navigateTo('checking');

    const deps = await checkDependencies();
    setDependencies(deps);

    const missingRequired = deps.filter(d => d.required && !d.installed);
    if (missingRequired.length === 0) {
      setStep('deps-ready');
    } else {
      setStep('deps-missing');
    }
  };

  const handleDepsReadyNext = async () => {
    navigateTo('cross-platform');
  };

  const handleDepsMissingRetry = async () => {
    setStep('checking');
    const deps = await checkDependencies();
    setDependencies(deps);

    const missingRequired = deps.filter(d => d.required && !d.installed);
    if (missingRequired.length === 0) {
      setStep('deps-ready');
    } else {
      setStep('deps-missing');
    }
  };

  const handleDepsMissingContinue = async () => {
    navigateTo('cross-platform');
  };

  const handleCrossPlatformYes = async () => {
    navigateTo('docker-options');
  };

  const handleDockerOptionsDownload = async () => {
    navigateTo('sdk-license');
  };

  const handleDockerOptionsSkip = async () => {
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setUseInterfaces(loadedDefaults.project?.useInterfaces ?? true);
    navigateTo('projects');
  };

  const handleSDKLicenseAgree = async () => {
    const docker = await getDockerStatus();
    setDockerStatus(docker);
    navigateTo('docker-setup');
  };

  const handleSDKLicenseDecline = async () => {
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setUseInterfaces(loadedDefaults.project?.useInterfaces ?? true);
    navigateTo('projects');
  };

  const handleCrossPlatformSkip = async () => {
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setUseInterfaces(loadedDefaults.project?.useInterfaces ?? true);
    navigateTo('projects');
  };

  const handleDockerCheckAgain = async () => {
    const docker = await getDockerStatus();
    setDockerStatus(docker);
  };

  const handleDockerBuildImage = async () => {
    setBuildingImage(true);
    await buildDockerImage();

    const unsubscribe = subscribeDockerStatus((status) => {
      setDockerStatus(status);
      if (status.pullStatus !== 'pulling') {
        setBuildingImage(false);
        unsubscribe();
      }
    });
  };

  const handleDockerContinueBackground = async () => {
    if (buildingImage || (dockerStatus && dockerStatus.pullStatus === 'pulling')) {
      setBackgroundDockerStarted(true);
    }
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setUseInterfaces(loadedDefaults.project?.useInterfaces ?? true);
    navigateTo('projects');
  };

  const handleDockerSkip = async () => {
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setUseInterfaces(loadedDefaults.project?.useInterfaces ?? true);
    navigateTo('projects');
  };

  const handleProjectsNext = () => {
    navigateTo('signing');
  };

  const handleProjectsSkip = () => {
    navigateTo('signing');
  };

  const handleSigningNext = () => {
    navigateTo('language-select');
  };

  const handleSigningSkip = () => {
    navigateTo('language-select');
  };

  const handleLanguageSelectNext = () => {
    if (preferTypeScript) {
      navigateTo('binding-style');
    } else {
      navigateTo('template-select');
    }
  };

  const handleBindingStyleNext = () => {
    navigateTo('template-select');
  };

  const handleTemplateSelectNext = async () => {
    const templateName = preferTypeScript && selectedFramework !== 'vanilla'
      ? `${selectedFramework}-ts`
      : preferTypeScript && selectedFramework === 'vanilla'
        ? 'vanilla-ts'
        : selectedFramework;

    const updatedDefaults = {
      ...defaults,
      project: {
        ...defaults.project,
        defaultTemplate: templateName,
        useInterfaces: preferTypeScript ? useInterfaces : true,
      }
    };

    setSavingDefaults(true);
    await saveDefaults(updatedDefaults);
    setSavingDefaults(false);
    navigateTo('complete');
  };

  const handleTemplateSelectSkip = async () => {
    const updatedDefaults = {
      ...defaults,
      project: {
        ...defaults.project,
        useInterfaces: preferTypeScript ? useInterfaces : true,
      }
    };
    setSavingDefaults(true);
    await saveDefaults(updatedDefaults);
    setSavingDefaults(false);
    navigateTo('complete');
  };

  // Stream Docker status in background via SSE
  useEffect(() => {
    if (backgroundDockerStarted && (buildingImage || (dockerStatus && dockerStatus.pullStatus === 'pulling'))) {
      const unsubscribe = subscribeDockerStatus((status) => {
        setDockerStatus(status);
        if (status.pullStatus !== 'pulling') {
          setBuildingImage(false);
        }
      });
      return unsubscribe;
    }
  }, [backgroundDockerStarted, buildingImage, dockerStatus?.pullStatus]);

  useEffect(() => {
    if (prevDockerPullStatus === 'pulling' && dockerStatus?.pullStatus === 'complete' && step !== 'docker-setup') {
      setShowDockerToast(true);
      const timer = setTimeout(() => setShowDockerToast(false), 3000);
      return () => clearTimeout(timer);
    }
    setPrevDockerPullStatus(dockerStatus?.pullStatus || null);
  }, [dockerStatus?.pullStatus, step]);

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      <div className="min-h-screen bg-gray-50 dark:bg-[#0f0f0f] flex items-center justify-center p-4 transition-colors relative overflow-hidden">
        {/* Scrolling background - shown on all pages */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="scrolling-bg w-full h-[200%] opacity-[0.08] dark:opacity-[0.06]">
            <img src="/showcase/montage.png" alt="" className="w-full h-1/2 object-cover object-center" />
            <img src="/showcase/montage.png" alt="" className="w-full h-1/2 object-cover object-center" />
          </div>
        </div>

        <div className="w-[75vw] max-w-[1200px] h-[75vh] max-h-[800px] glass-card rounded-2xl flex overflow-hidden relative z-10">
          {/* Sidebar */}
          <Sidebar currentStep={step} dockerStatus={dockerStatus} buildingDocker={backgroundDockerStarted && (buildingImage || dockerStatus?.pullStatus === 'pulling')} />

          {/* Content area - distinct from sidebar in dark mode */}
          <div className="flex-1 flex flex-col min-w-0 bg-white/50 dark:bg-white/[0.03] relative">
            <AnimatePresence>
              {showDockerToast && (
                <motion.div
                  initial={{ opacity: 0, y: -20, scale: 0.95 }}
                  animate={{ opacity: 1, y: 0, scale: 1 }}
                  exit={{ opacity: 0, y: -10, scale: 0.95 }}
                  transition={{ duration: 0.15, ease: "easeOut" }}
                  className="absolute top-4 right-4 z-50 flex items-center gap-2 px-3 py-2 bg-green-500 text-white rounded-lg shadow-lg"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                  <span className="text-sm font-medium">Docker image ready</span>
                </motion.div>
              )}
            </AnimatePresence>
            <div className="flex-1 flex flex-col min-h-0">
              <AnimatePresence mode="wait">
                {step === 'splash' && (
                  <SplashPage key="splash" onNext={handleSplashNext} />
                )}
                {step === 'checking' && (
                  <CheckingPage key="checking" />
                )}
                {step === 'deps-ready' && (
                  <DepsReadyPage key="deps-ready" onNext={handleDepsReadyNext} onBack={goBack} canGoBack={canGoBack} />
                )}
                {step === 'deps-missing' && (
                  <DepsMissingPage
                    key="deps-missing"
                    dependencies={dependencies}
                    onRetry={handleDepsMissingRetry}
                    onContinue={handleDepsMissingContinue}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'cross-platform' && (
                  <CrossPlatformPage
                    key="cross-platform"
                    dockerDep={dependencies.find(d => d.name === 'docker')}
                    onYes={handleCrossPlatformYes}
                    onSkip={handleCrossPlatformSkip}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'docker-options' && (
                  <DockerOptionsPage
                    key="docker-options"
                    onDownloadOfficial={handleDockerOptionsDownload}
                    onSkip={handleDockerOptionsSkip}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'sdk-license' && (
                  <SDKLicensePage
                    key="sdk-license"
                    onAgree={handleSDKLicenseAgree}
                    onDecline={handleSDKLicenseDecline}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'docker-setup' && (
                  <DockerSetupPage
                    key="docker-setup"
                    dockerStatus={dockerStatus}
                    buildingImage={buildingImage}
                    onBuildImage={handleDockerBuildImage}
                    onCheckAgain={handleDockerCheckAgain}
                    onContinueBackground={handleDockerContinueBackground}
                    onSkip={handleDockerSkip}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'projects' && (
                  <ProjectsPage
                    key="projects"
                    defaults={defaults}
                    onDefaultsChange={setDefaults}
                    onNext={handleProjectsNext}
                    onSkip={handleProjectsSkip}
                    saving={savingDefaults}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'signing' && (
                  <SigningStep
                    key="signing"
                    onNext={handleSigningNext}
                    onSkip={handleSigningSkip}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'language-select' && (
                  <LanguageSelectPage
                    key="language-select"
                    preferTypeScript={preferTypeScript}
                    onSelect={setPreferTypeScript}
                    onNext={handleLanguageSelectNext}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'binding-style' && (
                  <BindingStylePage
                    key="binding-style"
                    useInterfaces={useInterfaces}
                    onSelect={setUseInterfaces}
                    onNext={handleBindingStyleNext}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'template-select' && (
                  <TemplateSelectPage
                    key="template-select"
                    selectedFramework={selectedFramework}
                    preferTypeScript={preferTypeScript}
                    onSelect={setSelectedFramework}
                    onNext={handleTemplateSelectNext}
                    onSkip={handleTemplateSelectSkip}
                    onBack={goBack}
                    canGoBack={canGoBack}
                  />
                )}
                {step === 'complete' && (
                  <CompletePage key="complete" />
                )}
              </AnimatePresence>
            </div>
          </div>
        </div>
      </div>
    </ThemeContext.Provider>
  );
}
