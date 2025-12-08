import { useState, useEffect, createContext, useContext } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { DependencyStatus, SystemInfo, DockerStatus, GlobalDefaults } from './types';
import { checkDependencies, getState, getDockerStatus, buildDockerImage, close, getDefaults, saveDefaults, startDockerBuildBackground } from './api';
import WailsLogo from './components/WailsLogo';

type Step = 'splash' | 'dependencies' | 'defaults' | 'docker' | 'complete';
type Theme = 'light' | 'dark';

// Theme context
const ThemeContext = createContext<{ theme: Theme; toggleTheme: () => void }>({
  theme: 'dark',
  toggleTheme: () => {}
});

const useTheme = () => useContext(ThemeContext);

// Theme toggle button component
function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      onClick={toggleTheme}
      className="fixed top-4 left-4 z-50 p-2 rounded-lg bg-gray-200 dark:bg-gray-800 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors"
      title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {theme === 'dark' ? (
        // Sun icon for dark mode (click to switch to light)
        <svg className="w-5 h-5 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
      ) : (
        // Moon icon for light mode (click to switch to dark)
        <svg className="w-5 h-5 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
        </svg>
      )}
    </button>
  );
}

// Splash/Landing Page with scrolling background
function SplashPage({ onNext }: { onNext: () => void }) {
  const { theme, toggleTheme } = useTheme();

  return (
    <div className="h-full flex flex-col">
      {/* Main content area */}
      <div className="flex-1 min-h-0">
        <div className="h-full flex flex-col items-center justify-center">
          {/* Logo with glow effect */}
          <motion.div
            className="text-center mb-10"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.6, ease: "easeOut" }}
          >
            <div className="flex justify-center">
              <img
                src={theme === 'dark' ? '/assets/wails-logo-white-text-B284k7fX.svg' : '/assets/wails-logo-black-text-Cx-vsZ4W.svg'}
                alt="Wails"
                width={280}
                className="object-contain"
                style={{ filter: 'drop-shadow(0 0 60px rgba(239, 68, 68, 0.4))' }}
              />
            </div>
          </motion.div>

          {/* Apple-style welcome text */}
          <motion.div
            className="text-center px-8 max-w-lg"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4 tracking-tight">
              Welcome to Wails
            </h1>
            <p className="text-base text-gray-600 dark:text-gray-300 leading-relaxed mb-6">
              Let's get your development environment ready. We'll guide you through each step, making sure everything is set up perfectly.
            </p>
            <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed">
              This takes just a few minutes. You can skip any step and return later.
            </p>
          </motion.div>
        </div>
      </div>

      {/* Footer - matches TemplateFooter dimensions */}
      <div className="flex-shrink-0">
        <div className="w-full flex justify-between items-center pt-4 mt-4 border-t border-gray-200 dark:border-gray-800">
          {/* Left side: Theme toggle and Sponsor */}
          <div className="flex items-center gap-2">
            <button
              onClick={toggleTheme}
              className="p-1.5 rounded-md bg-gray-200/80 dark:bg-gray-800/80 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors"
              title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
            >
              {theme === 'dark' ? (
                <svg className="w-3.5 h-3.5 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              ) : (
                <svg className="w-3.5 h-3.5 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                </svg>
              )}
            </button>
            <a
              href="https://github.com/sponsors/leaanthony"
              target="_blank"
              rel="noopener noreferrer"
              className="p-1.5 rounded-md bg-gray-200/80 dark:bg-gray-800/80 hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors group"
              title="Sponsor Wails"
            >
              <svg className="w-3.5 h-3.5 text-red-500 group-hover:text-red-600 dark:group-hover:text-red-400" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
              </svg>
            </a>
          </div>

          {/* Get Started button - matches template button dimensions */}
          <button
            onClick={onNext}
            className="px-3 py-1.5 text-xs rounded-md bg-red-600 text-white hover:bg-red-500 transition-colors"
          >
            Get Started
          </button>
        </div>
      </div>
    </div>
  );
}

// Classic wizard page slide animation
const pageVariants = {
  initial: { opacity: 0, x: 50 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: -50 }
};

// Wizard step indicator
function StepIndicator({ steps, currentStep }: { steps: { id: Step; label: string }[]; currentStep: Step }) {
  const currentIndex = steps.findIndex(s => s.id === currentStep);

  return (
    <div className="flex items-center justify-center gap-1 text-[11px] text-gray-500 dark:text-gray-400">
      {steps.map((step, i) => (
        <div key={step.id} className="flex items-center">
          <span className={i <= currentIndex ? 'text-gray-900 dark:text-white font-medium' : 'text-gray-400 dark:text-gray-500'}>
            {step.label}
          </span>
          {i < steps.length - 1 && (
            <span className="mx-1.5 text-gray-400 dark:text-gray-600">&rsaquo;</span>
          )}
        </div>
      ))}
    </div>
  );
}

// Template footer with theme toggle + sponsor on left, navigation on right (matches saved design)
function TemplateFooter({
  onBack,
  onNext,
  nextLabel = 'Next',
  backLabel = '← Back',
  showBack = true,
  nextDisabled = false,
  showRetry = false,
  onRetry
}: {
  onBack?: () => void;
  onNext: () => void;
  nextLabel?: string;
  backLabel?: string;
  showBack?: boolean;
  nextDisabled?: boolean;
  showRetry?: boolean;
  onRetry?: () => void;
}) {
  const { theme, toggleTheme } = useTheme();

  return (
    <div className="w-full flex justify-between items-center pt-4 mt-4 border-t border-gray-200 dark:border-gray-800">
      {/* Left side: Theme toggle and Sponsor */}
      <div className="flex items-center gap-2">
        <button
          onClick={toggleTheme}
          className="p-1.5 rounded-md bg-gray-200/80 dark:bg-gray-800/80 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors"
          title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
        >
          {theme === 'dark' ? (
            <svg className="w-3.5 h-3.5 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          ) : (
            <svg className="w-3.5 h-3.5 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
          )}
        </button>
        <a
          href="https://github.com/sponsors/leaanthony"
          target="_blank"
          rel="noopener noreferrer"
          className="p-1.5 rounded-md bg-gray-200/80 dark:bg-gray-800/80 hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors group"
          title="Sponsor Wails"
        >
          <svg className="w-3.5 h-3.5 text-red-500 group-hover:text-red-600 dark:group-hover:text-red-400" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </a>
      </div>

      {/* Right side: Navigation buttons */}
      <div className="flex gap-2">
        {showBack && onBack && (
          <button
            onClick={onBack}
            className="px-3 py-1.5 text-xs rounded-md bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
          >
            {backLabel}
          </button>
        )}
        {showRetry && onRetry && (
          <button
            onClick={onRetry}
            className="px-3 py-1.5 text-xs rounded-md bg-blue-600 text-white hover:bg-blue-500 transition-colors flex items-center gap-1.5"
          >
            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Retry
          </button>
        )}
        <button
          onClick={onNext}
          disabled={nextDisabled}
          className={`px-3 py-1.5 text-xs rounded-md transition-colors ${
            nextDisabled
              ? 'bg-gray-200 dark:bg-gray-700 text-gray-400 cursor-not-allowed'
              : 'bg-red-600 text-white hover:bg-red-500'
          }`}
        >
          {nextLabel}
        </button>
      </div>
    </div>
  );
}

// Legacy wizard footer (kept for backwards compatibility)
function WizardFooter({
  onBack,
  onNext,
  onCancel: _onCancel,
  nextLabel = 'Next',
  backLabel = 'Back',
  showBack = true,
  nextDisabled = false,
  showRetry = false,
  onRetry
}: {
  onBack?: () => void;
  onNext: () => void;
  onCancel?: () => void;
  nextLabel?: string;
  backLabel?: string;
  showBack?: boolean;
  nextDisabled?: boolean;
  showRetry?: boolean;
  onRetry?: () => void;
}) {
  return (
    <TemplateFooter
      onBack={onBack}
      onNext={onNext}
      nextLabel={nextLabel}
      backLabel={showBack ? '← Back' : backLabel}
      showBack={showBack}
      nextDisabled={nextDisabled}
      showRetry={showRetry}
      onRetry={onRetry}
    />
  );
}

// Dependency row component
function DependencyRow({
  dep
}: {
  dep: DependencyStatus;
}) {
  return (
    <div className="flex items-start gap-2 py-1.5 border-b border-gray-200/50 dark:border-gray-800/50 last:border-0">
      {/* Status icon */}
      <div className="mt-0.5">
        {dep.installed ? (
          <div className="w-4 h-4 rounded-full bg-green-500/20 flex items-center justify-center">
            <svg className="w-2.5 h-2.5 text-green-500 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
            </svg>
          </div>
        ) : dep.required ? (
          <div className="w-4 h-4 rounded-full bg-red-500/20 flex items-center justify-center">
            <svg className="w-2.5 h-2.5 text-red-500 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
        ) : (
          <div className="w-4 h-4 rounded-full bg-gray-400/20 dark:bg-gray-600/20 flex items-center justify-center">
            <div className="w-1.5 h-1.5 rounded-full bg-gray-400 dark:bg-gray-500" />
          </div>
        )}
      </div>

      {/* Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className={`text-sm ${dep.installed ? 'text-gray-900 dark:text-white' : dep.required ? 'text-red-600 dark:text-red-300' : 'text-gray-500 dark:text-gray-400'}`}>
            {dep.name}
          </span>
          {!dep.required && (
            <span className="text-[10px] text-gray-500">(optional)</span>
          )}
          <span className="flex-1" />
          {dep.version && (
            <span className="text-[10px] text-gray-500 font-mono">{dep.version}</span>
          )}
        </div>
        {dep.message && (
          <p className="text-[11px] text-gray-500 mt-0.5">{dep.message}</p>
        )}

        {/* Help URL link for non-system installs */}
        {!dep.installed && dep.helpUrl && (
          <div className="mt-1">
            <a
              href={dep.helpUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-xs text-blue-500 dark:text-blue-400 hover:text-blue-600 dark:hover:text-blue-300"
            >
              Install from {new URL(dep.helpUrl).hostname}
              <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
          </div>
        )}
      </div>
    </div>
  );
}

// Dependencies Page
function DependenciesPage({
  dependencies,
  onNext,
  onBack,
  onRetry,
  checking
}: {
  dependencies: DependencyStatus[];
  onNext: () => void;
  onBack: () => void;
  onCancel: () => void;
  onRetry: () => void;
  checking: boolean;
}) {
  const { theme } = useTheme();
  const [copied, setCopied] = useState(false);
  const missingRequired = dependencies.filter(d => d.required && !d.installed);
  const allRequiredInstalled = dependencies.length > 0 && missingRequired.length === 0;
  const missingDeps = dependencies.filter(d => !d.installed);

  // Build combined install command from all missing deps that have system commands (starting with sudo)
  const combinedInstallCommand = (() => {
    const systemCommands = missingDeps
      .filter(d => d.installCommand?.startsWith('sudo '))
      .map(d => d.installCommand!);

    if (systemCommands.length === 0) return null;

    // Extract package names from "sudo pacman -S pkg" style commands
    // Group by package manager
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

    if (pacmanPkgs.length > 0) {
      return `sudo pacman -S ${pacmanPkgs.join(' ')}`;
    } else if (aptPkgs.length > 0) {
      return `sudo apt install ${aptPkgs.join(' ')}`;
    } else if (dnfPkgs.length > 0) {
      return `sudo dnf install ${dnfPkgs.join(' ')}`;
    }

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
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.2 }}
      className="h-full flex flex-col"
    >
      {/* Header: Logo left, title right */}
      <div className="flex items-center gap-6 mb-4 flex-shrink-0">
        <div className="flex-shrink-0">
          <img
            src={theme === 'dark' ? '/assets/wails-logo-white-text-B284k7fX.svg' : '/assets/wails-logo-black-text-Cx-vsZ4W.svg'}
            alt="Wails"
            width={80}
            className="object-contain"
            style={{ filter: 'drop-shadow(0 0 60px rgba(239, 68, 68, 0.4))' }}
          />
        </div>
        <div className="flex-1">
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">System Dependencies</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
            The following dependencies are needed to build Wails applications.
          </p>
        </div>
      </div>

      {/* Scrollable content area */}
      <div className="flex-1 overflow-y-auto scrollbar-thin min-h-0 px-4">
        {/* Status Summary - show above deps when all good OR show checking spinner */}
        {checking ? (
          <div className="rounded-lg p-3 bg-gray-100 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 mb-4">
            <div className="flex items-center gap-3 text-gray-600 dark:text-gray-400 text-sm">
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                className="w-4 h-4 border-2 border-gray-400 dark:border-gray-600 border-t-red-500 rounded-full"
              />
              Checking dependencies...
            </div>
          </div>
        ) : allRequiredInstalled && (
          <div className="rounded-lg p-3 bg-green-500/10 border border-green-500/20 mb-4">
            <div className="flex items-center gap-2 text-green-600 dark:text-green-400 text-sm">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              All required dependencies are installed.
            </div>
          </div>
        )}

        {/* All Dependencies */}
        <div className="mb-4">
          <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg px-4">
            {dependencies.length > 0 ? (
              dependencies.map(dep => (
                <DependencyRow
                  key={dep.name}
                  dep={dep}
                />
              ))
            ) : !checking && (
              <div className="py-4 text-center text-sm text-gray-500">
                No dependencies to check.
              </div>
            )}
          </div>
        </div>

        {/* Combined Install Command */}
        {combinedInstallCommand && (
          <div className="mb-4 p-3 bg-gray-100 dark:bg-gray-900/50 rounded-lg">
            <div className="text-xs text-gray-600 dark:text-gray-300 mb-2">Install all missing dependencies:</div>
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
                  <svg className="w-5 h-5 text-green-500 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
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
      </div>

      {/* Footer - grounded to bottom */}
      <div className="flex-shrink-0">
        <TemplateFooter
          onBack={onBack}
          onNext={onNext}
          nextLabel="Next"
          nextDisabled={checking}
          showRetry={!checking && !allRequiredInstalled && dependencies.length > 0}
          onRetry={onRetry}
        />
      </div>
    </motion.div>
  );
}

// Docker Page
function DockerPage({
  dockerStatus,
  buildingImage,
  onBuildImage,
  onNext,
  onBack,
  onCancel
}: {
  dockerStatus: DockerStatus | null;
  buildingImage: boolean;
  onBuildImage: () => void;
  onNext: () => void;
  onBack: () => void;
  onCancel: () => void;
}) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.2 }}
    >
      <div className="mb-6">
        <h2 className="text-xl font-bold mb-1 text-gray-900 dark:text-white">Cross-Platform Builds</h2>
        <p className="text-sm text-gray-600 dark:text-gray-300">
          Docker enables building for macOS, Windows, and Linux from any platform.
        </p>
      </div>

      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-4 mb-6">
        <div className="flex items-start gap-4">
          <div className="w-12 h-12 rounded-xl bg-blue-500/20 flex items-center justify-center flex-shrink-0">
            <svg className="w-7 h-7" viewBox="0 0 756.26 596.9">
              <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27ZM189.67,206.39h-81.7v81.7h81.7v-81.7ZM295.22,206.39h-81.7v81.7h81.7v-81.7ZM400.77,206.39h-81.7v81.7h81.7v-81.7ZM506.32,206.39h-81.7v81.7h81.7v-81.7ZM84.12,206.39H2.42v81.7h81.7v-81.7ZM189.67,103.2h-81.7v81.7h81.7v-81.7ZM295.22,103.2h-81.7v81.7h81.7v-81.7ZM400.77,103.2h-81.7v81.7h81.7v-81.7ZM400.77,0h-81.7v81.7h81.7V0Z"/>
            </svg>
          </div>

          <div className="flex-1">
            <h3 className="font-medium text-gray-900 dark:text-white mb-1">Docker Status</h3>

            {!dockerStatus ? (
              <div className="text-sm text-gray-600 dark:text-gray-300">Checking Docker...</div>
            ) : !dockerStatus.installed ? (
              <div>
                <div className="flex items-center gap-2 text-yellow-600 dark:text-yellow-400 text-sm mb-2">
                  <span className="w-2 h-2 rounded-full bg-yellow-500" />
                  Not installed
                </div>
                <p className="text-xs text-gray-500 mb-2">
                  Docker is optional but required for cross-platform builds.
                </p>
                <a
                  href="https://docs.docker.com/get-docker/"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 text-xs text-blue-500 dark:text-blue-400 hover:text-blue-600 dark:hover:text-blue-300"
                >
                  Install Docker Desktop
                  <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              </div>
            ) : !dockerStatus.running ? (
              <div>
                <div className="flex items-center gap-2 text-yellow-600 dark:text-yellow-400 text-sm mb-2">
                  <span className="w-2 h-2 rounded-full bg-yellow-500" />
                  Installed but not running
                </div>
                <p className="text-xs text-gray-500">
                  Start Docker Desktop to enable cross-platform builds.
                </p>
              </div>
            ) : dockerStatus.imageBuilt ? (
              <div>
                <div className="flex items-center gap-2 text-green-600 dark:text-green-400 text-sm mb-2">
                  <span className="w-2 h-2 rounded-full bg-green-500" />
                  Ready for cross-platform builds
                </div>
                <p className="text-xs text-gray-500">
                  Docker {dockerStatus.version} • wails-cross image installed
                </p>
              </div>
            ) : buildingImage ? (
              <div>
                <div className="flex items-center gap-2 text-blue-500 dark:text-blue-400 text-sm mb-2">
                  <motion.span
                    className="w-3 h-3 border-2 border-blue-500 dark:border-blue-400 border-t-transparent rounded-full"
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                  />
                  Building wails-cross image... {dockerStatus.pullProgress}%
                </div>
                <div className="h-1.5 bg-gray-300 dark:bg-gray-700 rounded-full overflow-hidden">
                  <motion.div
                    className="h-full bg-blue-500"
                    animate={{ width: `${dockerStatus.pullProgress}%` }}
                  />
                </div>
              </div>
            ) : (
              <div>
                <div className="flex items-center gap-2 text-gray-600 dark:text-gray-400 text-sm mb-2">
                  <span className="w-2 h-2 rounded-full bg-gray-400 dark:bg-gray-500" />
                  Cross-compilation image not installed
                </div>
                <p className="text-xs text-gray-500 mb-3">
                  Docker {dockerStatus.version} is running. Build the wails-cross image to enable cross-platform builds.
                </p>
                <button
                  onClick={onBuildImage}
                  className="text-sm px-4 py-2 rounded-lg bg-blue-500/20 text-blue-600 dark:text-blue-400 hover:bg-blue-500/30 transition-colors border border-blue-500/30"
                >
                  Build Cross-Compilation Image
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-4">
        <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">What you can build:</h3>
        <div className="grid grid-cols-3 gap-4 text-center text-sm">
          <div className="py-2">
            <div className="flex justify-center mb-2">
              {/* Apple logo */}
              <svg className="w-8 h-8 text-gray-700 dark:text-gray-300" viewBox="0 0 24 24" fill="currentColor">
                <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
              </svg>
            </div>
            <div className="text-gray-700 dark:text-gray-300">macOS</div>
            <div className="text-xs text-gray-500">.app / .dmg</div>
          </div>
          <div className="py-2">
            <div className="flex justify-center mb-2">
              {/* Windows logo */}
              <svg className="w-8 h-8 text-gray-700 dark:text-gray-300" viewBox="0 0 24 24" fill="currentColor">
                <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
              </svg>
            </div>
            <div className="text-gray-700 dark:text-gray-300">Windows</div>
            <div className="text-xs text-gray-500">.exe / .msi</div>
          </div>
          <div className="py-2">
            <div className="flex justify-center mb-2">
              {/* Tux - Linux penguin */}
              <svg className="w-8 h-8" viewBox="0 0 1024 1024" fill="currentColor">
                <path className="text-gray-700 dark:text-gray-300" fillRule="evenodd" clipRule="evenodd" d="M186.828,734.721c8.135,22.783-2.97,48.36-25.182,55.53c-12.773,4.121-27.021,5.532-40.519,5.145c-24.764-0.714-32.668,8.165-24.564,31.376c2.795,8.01,6.687,15.644,10.269,23.363c7.095,15.287,7.571,30.475-0.168,45.697c-2.572,5.057-5.055,10.168-7.402,15.337c-9.756,21.488-5.894,30.47,17.115,36.3c18.451,4.676,37.425,7.289,55.885,11.932c40.455,10.175,80.749,21,121.079,31.676c20.128,5.325,40.175,9.878,61.075,3.774c27.01-7.889,41.849-27.507,36.217-54.78c-4.359-21.112-10.586-43.132-21.634-61.314c-26.929-44.322-56.976-86.766-86.174-129.69c-5.666-8.329-12.819-15.753-19.905-22.987c-23.511-24.004-32.83-26.298-64.022-16.059c-7.589-15.327-5.198-31.395-2.56-47.076c1.384-8.231,4.291-16.796,8.718-23.821c18.812-29.824,29.767-62.909,41.471-95.738c13.545-37.999,30.87-73.47,57.108-105.131c21.607-26.074,38.626-55.982,57.303-84.44c6.678-10.173,6.803-21.535,6.23-33.787c-2.976-63.622-6.561-127.301-6.497-190.957c0.081-78.542,65.777-139.631,156.443-127.536c99.935,13.331,159.606,87.543,156.629,188.746c-2.679,91.191,27.38,170.682,89.727,239.686c62.132,68.767,91.194,153.119,96.435,245.38c0.649,11.46-1.686,23.648-5.362,34.583c-2.265,6.744-9.651,11.792-14.808,17.536c-6.984,7.781-14.497,15.142-20.959,23.328c-12.077,15.294-25.419,28.277-45.424,32.573c-30.163,6.475-50.177-2.901-63.81-30.468c-1.797-3.636-3.358-7.432-5.555-10.812c-5.027-7.741-10.067-18.974-20.434-15.568c-6.727,2.206-14.165,11.872-15.412,19.197c-2.738,16.079-5.699,33.882-1.532,49.047c11.975,43.604,9.224,86.688,3.062,130.371c-3.513,24.898-0.414,49.037,23.13,63.504c24.495,15.044,48.407,7.348,70.818-6.976c3.742-2.394,7.25-5.249,10.536-8.252c30.201-27.583,65.316-46.088,104.185-58.488c14.915-4.759,29.613-11.405,42.97-19.554c19.548-11.932,18.82-25.867-0.854-38.036c-7.187-4.445-14.944-8.5-22.984-10.933c-23.398-7.067-34.812-23.963-39.767-46.375c-3.627-16.398-4.646-32.782,4.812-51.731c1.689,10.577,2.771,17.974,4.062,25.334c5.242,29.945,20.805,52.067,48.321,66.04c8.869,4.5,17.161,10.973,24.191,18.055c10.372,10.447,10.407,22.541,0.899,33.911c-4.886,5.837-10.683,11.312-17.052,15.427c-11.894,7.685-23.962,15.532-36.92,21.056c-45.461,19.375-84.188,48.354-120.741,80.964c-19.707,17.582-44.202,15.855-68.188,13.395c-21.502-2.203-38.363-12.167-48.841-31.787c-6.008-11.251-15.755-18.053-28.35-18.262c-42.991-0.722-85.995-0.785-128.993-0.914c-8.92-0.026-17.842,0.962-26.769,1.1c-25.052,0.391-47.926,7.437-68.499,21.808c-5.987,4.186-12.068,8.24-17.954,12.562c-19.389,14.233-40.63,17.873-63.421,10.497c-25.827-8.353-51.076-18.795-77.286-25.591c-38.792-10.057-78.257-17.493-117.348-26.427c-43.557-9.959-51.638-24.855-33.733-65.298c8.605-19.435,8.812-38.251,3.55-58.078c-2.593-9.773-5.126-19.704-6.164-29.72c-1.788-17.258,4.194-24.958,21.341-27.812c12.367-2.059,25.069-2.132,37.423-4.255C165.996,776.175,182.158,759.821,186.828,734.721z M698.246,454.672c9.032,15.582,18.872,30.76,26.936,46.829c20.251,40.355,34.457,82.42,30.25,128.537c-0.871,9.573-2.975,19.332-6.354,28.313c-5.088,13.528-18.494,19.761-33.921,17.5c-13.708-2.007-15.566-12.743-16.583-23.462c-1.035-10.887-1.435-21.864-1.522-32.809c-0.314-39.017-7.915-76.689-22.456-112.7c-5.214-12.915-14.199-24.3-21.373-36.438c-2.792-4.72-6.521-9.291-7.806-14.435c-8.82-35.31-21.052-68.866-43.649-98.164c-11.154-14.454-14.638-31.432-9.843-49.572c1.656-6.269,3.405-12.527,4.695-18.875c3.127-15.406-1.444-22.62-15.969-28.01c-15.509-5.752-30.424-13.273-46.179-18.138c-12.963-4.001-15.764-12.624-15.217-23.948c0.31-6.432,0.895-13.054,2.767-19.159c3.27-10.672,9.56-18.74,21.976-19.737c12.983-1.044,22.973,4.218,28.695,16.137c5.661,11.8,6.941,23.856,1.772,36.459c-4.638,11.314-0.159,17.13,11.52,13.901c4.966-1.373,11.677-7.397,12.217-11.947c2.661-22.318,1.795-44.577-9.871-64.926c-11.181-19.503-31.449-27.798-52.973-21.69c-26.941,7.646-39.878,28.604-37.216,60.306c0.553,6.585,1.117,13.171,1.539,18.14c-15.463-1.116-29.71-2.144-44.146-3.184c-0.73-8.563-0.741-16.346-2.199-23.846c-1.843-9.481-3.939-19.118-7.605-27.993c-4.694-11.357-12.704-20.153-26.378-20.08c-13.304,0.074-20.082,9.253-25.192,19.894c-11.385,23.712-9.122,47.304,1.739,70.415c1.69,3.598,6.099,8.623,8.82,8.369c3.715-0.347,7.016-5.125,11.028-8.443c-17.322-9.889-25.172-30.912-16.872-46.754c3.016-5.758,10.86-10.391,17.474-12.498c8.076-2.575,15.881,2.05,18.515,10.112c3.214,9.837,4.66,20.323,6.051,30.641c0.337,2.494-1.911,6.161-4.06,8.031c-12.73,11.068-25.827,21.713-38.686,32.635c-2.754,2.339-5.533,4.917-7.455,7.921c-5.453,8.523-6.483,16.016,3.903,22.612c6.351,4.035,11.703,10.012,16.616,15.86c7.582,9.018,17.047,14.244,28.521,13.972c46.214-1.09,91.113-6.879,128.25-38.61c1.953-1.668,7.641-1.83,9.262-0.271c1.896,1.823,2.584,6.983,1.334,9.451c-1.418,2.797-5.315,4.806-8.555,6.139c-22.846,9.401-45.863,18.383-68.699,27.808c-22.67,9.355-45.875,13.199-70.216,8.43c-2.864-0.562-5.932-0.076-10.576-0.076c10.396,14.605,21.893,24.62,38.819,23.571c12.759-0.79,26.125-2.244,37.846-6.879c17.618-6.967,33.947-17.144,51.008-25.588c5.737-2.837,11.903-5.131,18.133-6.474c2.185-0.474,5.975,2.106,7.427,4.334c0.804,1.237-1.1,5.309-2.865,6.903c-2.953,2.667-6.796,4.339-10.227,6.488c-21.264,13.325-42.521,26.658-63.771,40.002c-8.235,5.17-16.098,11.071-24.745,15.408c-16.571,8.316-28.156,6.68-40.559-7.016c-10.026-11.072-18.225-23.792-27.376-35.669c-2.98-3.87-6.41-7.393-9.635-11.074c-1.543,26.454-14.954,46.662-26.272,67.665c-12.261,22.755-21.042,45.964-8.633,69.951c-4.075,4.752-7.722,8.13-10.332,12.18c-29.353,45.525-52.72,93.14-52.266,149.186c0.109,13.75-0.516,27.55-1.751,41.24c-0.342,3.793-3.706,9.89-6.374,10.287c-3.868,0.573-10.627-1.946-12.202-5.111c-6.939-13.938-14.946-28.106-17.81-43.101c-3.031-15.865-0.681-32.759-0.681-50.958c-2.558,5.441-5.907,9.771-6.539,14.466c-1.612,11.975-3.841,24.322-2.489,36.14c2.343,20.486,5.578,41.892,21.418,56.922c21.76,20.642,44.75,40.021,67.689,59.375c20.161,17.01,41.426,32.724,61.388,49.954c22.306,19.257,15.029,51.589-13.006,60.711c-2.144,0.697-4.25,1.513-8.117,2.9c20.918,28.527,40.528,56.508,38.477,93.371c23.886-27.406,2.287-47.712-10.241-69.677c6.972-6.97,12.504-8.75,21.861-1.923c10.471,7.639,23.112,15.599,35.46,16.822c62.957,6.229,123.157,2.18,163.56-57.379c2.57-3.788,8.177-5.519,12.37-8.205c1.981,4.603,5.929,9.354,5.596,13.78c-1.266,16.837-3.306,33.673-6.265,50.292c-1.978,11.097-6.572,21.71-8.924,32.766c-1.849,8.696,1.109,15.219,12.607,15.204c1.387-6.761,2.603-13.474,4.154-20.108c10.602-45.342,16.959-90.622,6.691-137.28c-3.4-15.454-2.151-32.381-0.526-48.377c2.256-22.174,12.785-32.192,33.649-37.142c2.765-0.654,6.489-3.506,7.108-6.002c4.621-18.597,18.218-26.026,35.236-28.913c19.98-3.386,39.191-0.066,59.491,10.485c-2.108-3.7-2.525-5.424-3.612-6.181c-8.573-5.968-17.275-11.753-25.307-17.164C776.523,585.58,758.423,514.082,698.246,454.672z M427.12,221.259c1.83-0.584,3.657-1.169,5.486-1.755c-2.37-7.733-4.515-15.555-7.387-23.097c-0.375-0.983-4.506-0.533-6.002-0.668C422.211,205.409,424.666,213.334,427.12,221.259z M565.116,212.853c5.3-12.117-1.433-21.592-14.086-20.792C555.663,198.899,560.315,205.768,565.116,212.853z"/>
              </svg>
            </div>
            <div className="text-gray-700 dark:text-gray-300">Linux</div>
            <div className="text-xs text-gray-500">.deb / .rpm / PKGBUILD</div>
          </div>
        </div>
      </div>

      <WizardFooter
        onBack={onBack}
        onNext={onNext}
        onCancel={onCancel}
        nextLabel="Next"
      />
    </motion.div>
  );
}

// Defaults Page - Configure global defaults for new projects
function DefaultsPage({
  defaults,
  onDefaultsChange,
  onNext,
  onBack,
  onCancel,
  saving
}: {
  defaults: GlobalDefaults;
  onDefaultsChange: (defaults: GlobalDefaults) => void;
  onNext: () => void;
  onBack: () => void;
  onCancel: () => void;
  saving: boolean;
}) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.2 }}
    >
      <div className="mb-3">
        <h2 className="text-lg font-bold mb-0.5 text-gray-900 dark:text-white">Project Defaults</h2>
        <p className="text-xs text-gray-600 dark:text-gray-300">
          Configure defaults for new Wails projects.
        </p>
      </div>

      {/* Author Information */}
      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-3 mb-3">
        <h3 className="text-[11px] font-medium text-gray-500 dark:text-gray-400 mb-2">Author Information</h3>
        <div className="grid grid-cols-2 gap-2">
          <div>
            <label className="block text-[10px] text-gray-500 mb-0.5">Your Name</label>
            <input
              type="text"
              value={defaults.author.name}
              onChange={(e) => onDefaultsChange({
                ...defaults,
                author: { ...defaults.author, name: e.target.value }
              })}
              placeholder="John Doe"
              className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none"
            />
          </div>
          <div>
            <label className="block text-[10px] text-gray-500 mb-0.5">Company</label>
            <input
              type="text"
              value={defaults.author.company}
              onChange={(e) => onDefaultsChange({
                ...defaults,
                author: { ...defaults.author, company: e.target.value }
              })}
              placeholder="My Company"
              className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none"
            />
          </div>
        </div>
      </div>

      {/* Project Defaults */}
      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-3 mb-3">
        <h3 className="text-[11px] font-medium text-gray-500 dark:text-gray-400 mb-2">Project Settings</h3>
        <div className="space-y-2">
          <div className="grid grid-cols-2 gap-2">
            <div>
              <label className="block text-[10px] text-gray-500 mb-0.5">Bundle ID Prefix</label>
              <input
                type="text"
                value={defaults.project.productIdentifierPrefix}
                onChange={(e) => onDefaultsChange({
                  ...defaults,
                  project: { ...defaults.project, productIdentifierPrefix: e.target.value }
                })}
                placeholder="com.mycompany"
                className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none font-mono"
              />
            </div>
            <div>
              <label className="block text-[10px] text-gray-500 mb-0.5">Default Version</label>
              <input
                type="text"
                value={defaults.project.defaultVersion}
                onChange={(e) => onDefaultsChange({
                  ...defaults,
                  project: { ...defaults.project, defaultVersion: e.target.value }
                })}
                placeholder="0.1.0"
                className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none font-mono"
              />
            </div>
          </div>
          <div>
            <label className="block text-[10px] text-gray-500 mb-0.5">Default Template</label>
            <select
              value={defaults.project.defaultTemplate}
              onChange={(e) => onDefaultsChange({
                ...defaults,
                project: { ...defaults.project, defaultTemplate: e.target.value }
              })}
              className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 focus:border-red-500 focus:outline-none"
            >
              <option value="vanilla">Vanilla (JavaScript)</option>
              <option value="vanilla-ts">Vanilla (TypeScript)</option>
              <option value="react">React</option>
              <option value="react-ts">React (TypeScript)</option>
              <option value="react-swc">React + SWC</option>
              <option value="react-swc-ts">React + SWC (TypeScript)</option>
              <option value="preact">Preact</option>
              <option value="preact-ts">Preact (TypeScript)</option>
              <option value="svelte">Svelte</option>
              <option value="svelte-ts">Svelte (TypeScript)</option>
              <option value="solid">Solid</option>
              <option value="solid-ts">Solid (TypeScript)</option>
              <option value="lit">Lit</option>
              <option value="lit-ts">Lit (TypeScript)</option>
              <option value="vue">Vue</option>
              <option value="vue-ts">Vue (TypeScript)</option>
            </select>
          </div>
        </div>
      </div>

      {/* macOS Signing */}
      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-3 mb-3">
        <div className="flex items-center gap-2 mb-1">
          <svg className="w-4 h-4 text-gray-500 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
            <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
          </svg>
          <h3 className="text-[11px] font-medium text-gray-500 dark:text-gray-400">macOS Code Signing</h3>
          <span className="text-[9px] text-gray-400 dark:text-gray-500">(optional)</span>
        </div>
        <p className="text-[9px] text-gray-400 dark:text-gray-500 mb-2 ml-6">These are public identifiers. App-specific passwords are stored securely in your Keychain.</p>
        <div className="space-y-2">
          <div>
            <label className="block text-[10px] text-gray-500 mb-0.5">Developer ID</label>
            <input
              type="text"
              value={defaults.signing?.macOS?.developerID || ''}
              onChange={(e) => onDefaultsChange({
                ...defaults,
                signing: {
                  ...defaults.signing,
                  macOS: { ...defaults.signing?.macOS, developerID: e.target.value, appleID: defaults.signing?.macOS?.appleID || '', teamID: defaults.signing?.macOS?.teamID || '' },
                  windows: defaults.signing?.windows || { certificatePath: '', timestampServer: '' }
                }
              })}
              placeholder="Developer ID Application: John Doe (TEAMID)"
              className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none font-mono"
            />
          </div>
          <div className="grid grid-cols-2 gap-2">
            <div>
              <label className="block text-[10px] text-gray-500 mb-0.5">Apple ID</label>
              <input
                type="email"
                value={defaults.signing?.macOS?.appleID || ''}
                onChange={(e) => onDefaultsChange({
                  ...defaults,
                  signing: {
                    ...defaults.signing,
                    macOS: { ...defaults.signing?.macOS, appleID: e.target.value, developerID: defaults.signing?.macOS?.developerID || '', teamID: defaults.signing?.macOS?.teamID || '' },
                    windows: defaults.signing?.windows || { certificatePath: '', timestampServer: '' }
                  }
                })}
                placeholder="you@example.com"
                className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none"
              />
            </div>
            <div>
              <label className="block text-[10px] text-gray-500 mb-0.5">Team ID</label>
              <input
                type="text"
                value={defaults.signing?.macOS?.teamID || ''}
                onChange={(e) => onDefaultsChange({
                  ...defaults,
                  signing: {
                    ...defaults.signing,
                    macOS: { ...defaults.signing?.macOS, teamID: e.target.value, developerID: defaults.signing?.macOS?.developerID || '', appleID: defaults.signing?.macOS?.appleID || '' },
                    windows: defaults.signing?.windows || { certificatePath: '', timestampServer: '' }
                  }
                })}
                placeholder="ABCD1234EF"
                className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none font-mono"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Windows Signing */}
      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-3 mb-3">
        <div className="flex items-center gap-2 mb-2">
          <svg className="w-4 h-4 text-gray-500 dark:text-gray-400" viewBox="0 0 24 24" fill="currentColor">
            <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
          </svg>
          <h3 className="text-[11px] font-medium text-gray-500 dark:text-gray-400">Windows Code Signing</h3>
          <span className="text-[9px] text-gray-400 dark:text-gray-500">(optional)</span>
        </div>
        <div className="space-y-2">
          <div>
            <label className="block text-[10px] text-gray-500 mb-0.5">Certificate Path (.pfx)</label>
            <input
              type="text"
              value={defaults.signing?.windows?.certificatePath || ''}
              onChange={(e) => onDefaultsChange({
                ...defaults,
                signing: {
                  ...defaults.signing,
                  macOS: defaults.signing?.macOS || { developerID: '', appleID: '', teamID: '' },
                  windows: { ...defaults.signing?.windows, certificatePath: e.target.value, timestampServer: defaults.signing?.windows?.timestampServer || '' }
                }
              })}
              placeholder="/path/to/certificate.pfx"
              className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none font-mono"
            />
          </div>
          <div>
            <label className="block text-[10px] text-gray-500 mb-0.5">Timestamp Server</label>
            <input
              type="text"
              value={defaults.signing?.windows?.timestampServer || ''}
              onChange={(e) => onDefaultsChange({
                ...defaults,
                signing: {
                  ...defaults.signing,
                  macOS: defaults.signing?.macOS || { developerID: '', appleID: '', teamID: '' },
                  windows: { ...defaults.signing?.windows, timestampServer: e.target.value, certificatePath: defaults.signing?.windows?.certificatePath || '' }
                }
              })}
              placeholder="http://timestamp.digicert.com"
              className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-2 py-1 text-xs text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none font-mono"
            />
          </div>
        </div>
      </div>

      {/* Info about where this is stored */}
      <div className="text-[10px] text-gray-500 dark:text-gray-600 mb-3">
        <span className="text-gray-400 dark:text-gray-500">Stored in:</span> ~/.config/wails/defaults.yaml
      </div>

      <WizardFooter
        onBack={onBack}
        onNext={onNext}
        onCancel={onCancel}
        nextLabel={saving ? "Saving..." : "Finish"}
        nextDisabled={saving}
      />
    </motion.div>
  );
}

// Persistent Docker status indicator - shown across all pages when Docker build is in progress
function DockerStatusIndicator({
  status,
  visible
}: {
  status: DockerStatus | null;
  visible: boolean;
}) {
  if (!visible || !status) return null;

  // Don't show if Docker is not installed/running or if image is already built
  if (!status.installed || !status.running) return null;
  if (status.imageBuilt && status.pullStatus !== 'pulling') return null;

  const isPulling = status.pullStatus === 'pulling';
  const progress = status.pullProgress || 0;

  return (
    <motion.div
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -10 }}
      className="fixed top-4 right-4 z-50"
    >
      <div className="bg-white/95 dark:bg-gray-900/95 border border-gray-200 dark:border-gray-700 rounded-lg shadow-xl px-4 py-3 backdrop-blur-sm min-w-[240px]">
        <div className="flex items-center gap-3">
          {/* Docker icon */}
          <div className="w-8 h-8 rounded-lg bg-blue-500/20 flex items-center justify-center flex-shrink-0">
            <svg className="w-5 h-5" viewBox="0 0 756.26 596.9">
              <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27Z"/>
            </svg>
          </div>

          <div className="flex-1 min-w-0">
            {isPulling ? (
              <>
                <div className="flex items-center gap-2 text-blue-600 dark:text-blue-400 text-sm mb-1">
                  <motion.span
                    className="w-3 h-3 border-2 border-blue-600 dark:border-blue-400 border-t-transparent rounded-full"
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                  />
                  <span className="truncate">Downloading cross-compile image...</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="flex-1 h-1.5 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                    <motion.div
                      className="h-full bg-blue-500"
                      animate={{ width: `${progress}%` }}
                    />
                  </div>
                  <span className="text-xs text-gray-500 tabular-nums">{progress}%</span>
                </div>
              </>
            ) : status.imageBuilt ? (
              <div className="flex items-center gap-2 text-green-600 dark:text-green-400 text-sm">
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span>Docker image ready</span>
              </div>
            ) : (
              <div className="text-sm text-gray-600 dark:text-gray-400">
                Preparing Docker build...
              </div>
            )}
          </div>
        </div>
      </div>
    </motion.div>
  );
}

// Copyable command component
function CopyableCommand({ command, label }: { command: string; label: string }) {
  const [copied, setCopied] = useState(false);

  const copyCommand = () => {
    navigator.clipboard.writeText(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div>
      <p className="text-gray-600 dark:text-gray-400 mb-1">{label}</p>
      <div className="flex items-center gap-2">
        <code className="flex-1 text-green-600 dark:text-green-400 font-mono text-xs bg-gray-100 dark:bg-gray-900 px-2 py-1 rounded">
          {command}
        </code>
        <button
          onClick={copyCommand}
          className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 transition-colors p-1"
          title="Copy command"
        >
          {copied ? (
            <svg className="w-4 h-4 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          ) : (
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          )}
        </button>
      </div>
    </div>
  );
}

// Complete Page
function CompletePage({ onClose }: { onClose: () => void }) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.2 }}
      className="text-center py-8"
    >
      <motion.div
        initial={{ scale: 0 }}
        animate={{ scale: 1 }}
        transition={{ type: 'spring', stiffness: 200, damping: 15 }}
        className="w-16 h-16 rounded-full bg-green-500/20 flex items-center justify-center mx-auto mb-6"
      >
        <svg className="w-8 h-8 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
        </svg>
      </motion.div>

      <h2 className="text-2xl font-bold mb-2 text-gray-900 dark:text-white">Setup Complete</h2>
      <p className="text-gray-600 dark:text-gray-300 mb-8">
        Your development environment is ready to use.
      </p>

      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-lg p-4 text-left mb-6 max-w-sm mx-auto">
        <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-3">Next Steps</h3>
        <div className="space-y-3 text-sm">
          <CopyableCommand command="wails3 init -n myapp" label="Create a new project:" />
          <CopyableCommand command="wails3 dev" label="Start development server:" />
          <CopyableCommand command="wails3 build" label="Build for production:" />
        </div>
      </div>

      <button
        onClick={onClose}
        className="px-6 py-2.5 rounded-lg bg-red-600 text-white font-medium hover:bg-red-500 transition-colors"
      >
        Close
      </button>
    </motion.div>
  );
}

// Main App
export default function App() {
  const [step, setStep] = useState<Step>('splash');
  const [dependencies, setDependencies] = useState<DependencyStatus[]>([]);
  const [_system, setSystem] = useState<SystemInfo | null>(null);
  const [dockerStatus, setDockerStatus] = useState<DockerStatus | null>(null);
  const [buildingImage, setBuildingImage] = useState(false);
  const [checkingDeps, setCheckingDeps] = useState(false);
  const [defaults, setDefaults] = useState<GlobalDefaults>({
    author: { name: '', company: '' },
    project: {
      productIdentifierPrefix: 'com.example',
      defaultTemplate: 'vanilla',
      copyrightTemplate: '© {year}, {company}',
      descriptionTemplate: 'A {name} application',
      defaultVersion: '0.1.0'
    }
  });
  const [savingDefaults, setSavingDefaults] = useState(false);
  const [backgroundDockerStarted, setBackgroundDockerStarted] = useState(false);
  const [theme, setTheme] = useState<Theme>(() => {
    // Default to dark, but check for saved preference or system preference
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('wails-setup-theme');
      if (saved === 'light' || saved === 'dark') return saved;
      if (window.matchMedia('(prefers-color-scheme: light)').matches) return 'light';
    }
    return 'dark';
  });

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

  const steps: { id: Step; label: string }[] = [
    { id: 'dependencies', label: 'Dependencies' },
    { id: 'docker', label: 'Docker' },
    { id: 'defaults', label: 'Defaults' },
    { id: 'complete', label: 'Complete' },
  ];

  useEffect(() => {
    init();
  }, []);

  const init = async () => {
    const state = await getState();
    setSystem(state.system);
  };

  // Trigger dependency check when entering dependencies page
  useEffect(() => {
    if (step === 'dependencies' && dependencies.length === 0 && !checkingDeps) {
      const check = async () => {
        setCheckingDeps(true);
        const deps = await checkDependencies();
        setDependencies(deps);
        setCheckingDeps(false);
      };
      check();
    }
  }, [step]);

  const handleNext = async () => {
    if (step === 'splash') {
      // Just transition to dependencies - checking happens there
      setStep('dependencies');
    } else if (step === 'dependencies') {
      // Check docker status and start background build if available
      const dockerDep = dependencies.find(d => d.name === 'docker');
      if (dockerDep?.installed) {
        const docker = await getDockerStatus();
        setDockerStatus(docker);
        // Start background Docker build (so it downloads while user configures defaults)
        startBackgroundDockerBuild(dependencies);
      }
      setStep('docker');
    } else if (step === 'docker') {
      // Load existing defaults when entering defaults page
      const loadedDefaults = await getDefaults();
      setDefaults(loadedDefaults);
      setStep('defaults');
    } else if (step === 'defaults') {
      // Save defaults before proceeding
      setSavingDefaults(true);
      await saveDefaults(defaults);
      setSavingDefaults(false);
      setStep('complete');
    }
  };

  const handleRetryDeps = async () => {
    setCheckingDeps(true);
    const deps = await checkDependencies();
    setDependencies(deps);
    setCheckingDeps(false);
  };

  const handleBack = () => {
    if (step === 'dependencies') setStep('splash');
    else if (step === 'docker') setStep('dependencies');
    else if (step === 'defaults') setStep('docker');
  };

  const handleBuildImage = async () => {
    setBuildingImage(true);
    await buildDockerImage();

    const poll = async () => {
      const status = await getDockerStatus();
      setDockerStatus(status);
      if (status.pullStatus === 'pulling') {
        setTimeout(poll, 1000);
      } else {
        setBuildingImage(false);
      }
    };
    poll();
  };

  // Start background Docker build after dependencies check
  const startBackgroundDockerBuild = async (deps: DependencyStatus[]) => {
    const dockerDep = deps.find(d => d.name === 'docker');
    if (!dockerDep?.installed || backgroundDockerStarted) return;

    setBackgroundDockerStarted(true);

    // Try to start background build
    const result = await startDockerBuildBackground();
    setDockerStatus(result.status);

    // If build started, poll for status
    if (result.started && result.status.pullStatus === 'pulling') {
      setBuildingImage(true);
      const poll = async () => {
        const status = await getDockerStatus();
        setDockerStatus(status);
        if (status.pullStatus === 'pulling') {
          setTimeout(poll, 1000);
        } else {
          setBuildingImage(false);
        }
      };
      setTimeout(poll, 1000);
    }
  };

  const handleClose = async () => {
    await close();
    window.close();
  };

  const handleCancel = handleClose;

  // Show Docker indicator on defaults page when Docker build is in progress (Docker now downloads while user configures)
  const showDockerIndicator = backgroundDockerStarted && step === 'defaults';

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      <div className="min-h-screen bg-gray-50 dark:bg-[#0f0f0f] flex items-center justify-center p-4 transition-colors relative overflow-hidden">
        {/* Scrolling background - only visible on splash */}
        {step === 'splash' && (
          <div className="absolute inset-0 overflow-hidden pointer-events-none">
            <div className="scrolling-bg w-full h-[200%] opacity-[0.08] dark:opacity-[0.06]">
              <img src="/showcase/montage.png" alt="" className="w-full h-1/2 object-cover object-center" />
              <img src="/showcase/montage.png" alt="" className="w-full h-1/2 object-cover object-center" />
            </div>
          </div>
        )}

        {/* Theme toggle - only show on pages that don't have their own footer */}
        {step !== 'splash' && step !== 'dependencies' && <ThemeToggle />}

        {/* Persistent Docker status indicator */}
        <AnimatePresence>
          {showDockerIndicator && (
            <DockerStatusIndicator
              status={dockerStatus}
              visible={showDockerIndicator}
            />
          )}
        </AnimatePresence>

        {/* Main content card */}
        <div className="w-full max-w-2xl bg-white dark:bg-gray-900/80 border border-gray-200 dark:border-gray-800 rounded-xl shadow-2xl h-[85vh] flex flex-col overflow-hidden relative z-10">
          <div className="flex-1 flex flex-col p-4 min-h-0">
            {/* Splash page - full height, no header */}
            {step === 'splash' && (
              <SplashPage onNext={handleNext} />
            )}

            {/* Dependencies page - uses PageTemplate layout (logo left, title right, footer at bottom) */}
            {step === 'dependencies' && (
              <DependenciesPage
                dependencies={dependencies}
                onNext={handleNext}
                onBack={handleBack}
                onCancel={handleCancel}
                onRetry={handleRetryDeps}
                checking={checkingDeps}
              />
            )}

            {/* Other pages with centered header */}
            {step !== 'splash' && step !== 'dependencies' && (
              <>
                {/* Header with logo and step indicator */}
                <div className="flex flex-col items-center mb-4 flex-shrink-0">
                  <WailsLogo size={120} theme={theme} />
                  <div className="mt-3">
                    <StepIndicator steps={steps} currentStep={step} />
                  </div>
                </div>

                {/* Page content */}
                <div className="flex-1 min-h-0 overflow-y-auto">
                  <AnimatePresence mode="wait">
                    {step === 'defaults' && (
                      <DefaultsPage
                        key="defaults"
                        defaults={defaults}
                        onDefaultsChange={setDefaults}
                        onNext={handleNext}
                        onBack={handleBack}
                        onCancel={handleCancel}
                        saving={savingDefaults}
                      />
                    )}
                    {step === 'docker' && (
                      <DockerPage
                        key="docker"
                        dockerStatus={dockerStatus}
                        buildingImage={buildingImage}
                        onBuildImage={handleBuildImage}
                        onNext={handleNext}
                        onBack={handleBack}
                        onCancel={handleCancel}
                      />
                    )}
                    {step === 'complete' && (
                      <CompletePage key="complete" onClose={handleClose} />
                    )}
                  </AnimatePresence>
                </div>
              </>
            )}
          </div>
        </div>
      </div>
    </ThemeContext.Provider>
  );
}
