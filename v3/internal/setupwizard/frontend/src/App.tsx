import { useState, useEffect, createContext, useContext, ReactNode } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { DependencyStatus, SystemInfo, DockerStatus, GlobalDefaults } from './types';
import { checkDependencies, getState, getDockerStatus, buildDockerImage, close, getDefaults, saveDefaults } from './api';
import wailsLogoWhite from './assets/wails-logo-white-text.svg';
import wailsLogoBlack from './assets/wails-logo-black-text.svg';

// OOBE Steps - branching state machine
type OOBEStep =
  | 'splash'
  | 'checking'
  | 'deps-ready'
  | 'deps-missing'
  | 'cross-platform'
  | 'docker-setup'
  | 'author'
  | 'project-defaults'
  | 'complete';

type Theme = 'light' | 'dark';

// Theme context
const ThemeContext = createContext<{ theme: Theme; toggleTheme: () => void }>({
  theme: 'dark',
  toggleTheme: () => {}
});

const useTheme = () => useContext(ThemeContext);

// Bottom controls - theme toggle + sponsor (inside content area, bottom left)
function BottomControls() {
  const { theme, toggleTheme } = useTheme();

  const handleSponsorClick = () => {
    window.open('https://github.com/sponsors/leaanthony', '_blank', 'noopener,noreferrer');
  };

  return (
    <div className="flex items-center gap-3">
      <button
        onClick={handleSponsorClick}
        className="p-1 hover:opacity-70 transition-opacity"
        title="Sponsor Wails"
      >
        <svg className="w-4 h-4 text-red-500" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
        </svg>
      </button>
      <button
        onClick={toggleTheme}
        className="p-1 hover:opacity-70 transition-opacity"
        title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
      >
        {theme === 'dark' ? (
          <svg className="w-4 h-4 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        ) : (
          <svg className="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        )}
      </button>
    </div>
  );
}

// Page animation variants - fade only
const pageVariants = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 }
};

// Logo footer - shown on all pages
function LogoFooter() {
  const { theme } = useTheme();

  return (
    <img
      src={theme === 'dark' ? wailsLogoWhite : wailsLogoBlack}
      alt="Wails"
      width={72}
      className="object-contain opacity-50"
    />
  );
}

// Page template - header + subheader, content, optional buttons
function PageTemplate({
  title,
  subtitle,
  children,
  primaryAction,
  primaryLabel,
  secondaryAction,
  secondaryLabel,
  primaryDisabled = false
}: {
  title: string;
  subtitle: string;
  children?: ReactNode;
  primaryAction?: () => void;
  primaryLabel?: string;
  secondaryAction?: () => void;
  secondaryLabel?: string;
  primaryDisabled?: boolean;
}) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col"
    >
      {/* Header - centered */}
      <div className="text-center mb-6 flex-shrink-0">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">{title}</h1>
        <p className="text-gray-500 dark:text-gray-400">{subtitle}</p>
      </div>

      {/* Scrollable content area */}
      <div className="flex-1 overflow-y-auto scrollbar-thin min-h-0">
        {children}
      </div>

      {/* Actions - only shown if provided */}
      {(primaryAction || secondaryAction) && (
        <div className="flex-shrink-0 pt-4 flex flex-col items-center gap-1.5">
          {primaryAction && primaryLabel && (
            <button
              onClick={primaryAction}
              disabled={primaryDisabled}
              className={`px-5 py-2 rounded-lg text-sm font-medium transition-colors ${
                primaryDisabled
                  ? 'bg-gray-200 dark:bg-gray-700 text-gray-400 cursor-not-allowed'
                  : 'bg-red-600 text-white hover:bg-red-500'
              }`}
            >
              {primaryLabel}
            </button>
          )}
          {secondaryAction && secondaryLabel && (
            <button
              onClick={secondaryAction}
              className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
            >
              {secondaryLabel}
            </button>
          )}
        </div>
      )}
    </motion.div>
  );
}

// Splash Page - simple welcome with Let's Start
function SplashPage({ onNext }: { onNext: () => void }) {
  const { theme } = useTheme();

  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col items-center justify-center"
    >
      {/* Logo with glow effect */}
      <motion.div
        className="text-center mb-10"
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
      >
        <div className="flex justify-center">
          <img
            src={theme === 'dark' ? wailsLogoWhite : wailsLogoBlack}
            alt="Wails"
            width={280}
            className="object-contain"
            style={{ filter: 'drop-shadow(0 0 60px rgba(239, 68, 68, 0.4))' }}
          />
        </div>
      </motion.div>

      {/* Welcome text */}
      <motion.div
        className="text-center px-8 max-w-lg"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4 tracking-tight">
          Welcome to Wails
        </h1>
        <p className="text-base text-gray-600 dark:text-gray-300 leading-relaxed mb-8">
          Build beautiful cross-platform apps using Go and web technologies
        </p>
      </motion.div>

      {/* Let's Start button */}
      <motion.button
        onClick={onNext}
        className="px-6 py-2.5 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-500 transition-colors"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.4 }}
      >
        Let's Start
      </motion.button>
    </motion.div>
  );
}

// Checking Screen - brief loading while checking dependencies
function CheckingPage() {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col items-center justify-center"
    >
      <motion.div
        className="w-12 h-12 border-3 border-gray-300 dark:border-gray-600 border-t-red-500 rounded-full mb-6"
        animate={{ rotate: 360 }}
        transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
      />
      <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
        Checking your system...
      </h2>
      <p className="text-gray-500 dark:text-gray-400">
        This will only take a moment
      </p>
    </motion.div>
  );
}

// Deps Ready Page - simple checkmark, deps are good
function DepsReadyPage({ onNext }: { onNext: () => void }) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col items-center justify-center"
    >
      {/* Animated checkmark */}
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

      <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">
        You're all set!
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-sm">
        Your system has everything needed to build Wails apps
      </p>

      <button
        onClick={onNext}
        className="px-5 py-2 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-500 transition-colors"
      >
        Continue
      </button>
    </motion.div>
  );
}

// Deps Missing Page - show what's missing with install command
function DepsMissingPage({
  dependencies,
  onRetry,
  onContinue
}: {
  dependencies: DependencyStatus[];
  onRetry: () => void;
  onContinue: () => void;
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
  onYes,
  onSkip
}: {
  onYes: () => void;
  onSkip: () => void;
}) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col items-center justify-center"
    >
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
        <button
          onClick={onYes}
          className="px-5 py-2 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-500 transition-colors"
        >
          Yes, set this up
        </button>
        <button
          onClick={onSkip}
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
        >
          Not right now
        </button>
      </div>
    </motion.div>
  );
}

// Docker Setup Page - handles install/not running/building states
function DockerSetupPage({
  dockerStatus,
  buildingImage,
  onBuildImage,
  onCheckAgain,
  onContinueBackground,
  onSkip
}: {
  dockerStatus: DockerStatus | null;
  buildingImage: boolean;
  onBuildImage: () => void;
  onCheckAgain: () => void;
  onContinueBackground: () => void;
  onSkip: () => void;
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
        className="h-full flex flex-col items-center justify-center"
      >
        <div className="w-16 h-16 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6">
          <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
            <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27Z"/>
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

        <p className="text-xs text-gray-400 dark:text-gray-500 mb-6 text-center">
          After installing, come back and we'll continue setting up
        </p>

        <div className="flex flex-col items-center gap-1.5">
          <button
            onClick={onCheckAgain}
            className="px-5 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 text-sm font-medium hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
          >
            Check Again
          </button>
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
        className="h-full flex flex-col items-center justify-center"
      >
        <div className="w-16 h-16 rounded-2xl bg-gray-200 dark:bg-gray-800 flex items-center justify-center mb-6 opacity-50">
          <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
            <path fill="#6b7280" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27Z"/>
          </svg>
        </div>

        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Start Docker
        </h2>
        <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-sm">
          Please start Docker Desktop to continue
        </p>

        <div className="flex flex-col items-center gap-1.5">
          <button
            onClick={onCheckAgain}
            className="px-5 py-2 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-500 transition-colors"
          >
            Check Again
          </button>
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
    return (
      <motion.div
        variants={pageVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        transition={{ duration: 0.3 }}
        className="h-full flex flex-col items-center justify-center"
      >
        <div className="w-16 h-16 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6">
          <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
            <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27Z"/>
          </svg>
        </div>

        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Building cross-compiler image
        </h2>

        <div className="w-64 mb-4">
          <div className="flex items-center justify-between text-sm text-gray-500 mb-1">
            <span>Progress</span>
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
          className="px-5 py-2 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-500 transition-colors"
        >
          Continue in background
        </button>
      </motion.div>
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
        className="h-full flex flex-col items-center justify-center"
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

        <button
          onClick={onContinueBackground}
          className="px-5 py-2 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-500 transition-colors"
        >
          Continue
        </button>
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
      className="h-full flex flex-col items-center justify-center"
    >
      <div className="w-16 h-16 rounded-2xl bg-blue-500/20 flex items-center justify-center mb-6">
        <svg className="w-10 h-10" viewBox="0 0 756.26 596.9">
          <path fill="#1d63ed" d="M743.96,245.25c-18.54-12.48-67.26-17.81-102.68-8.27-1.91-35.28-20.1-65.01-53.38-90.95l-12.32-8.27-8.21,12.4c-16.14,24.5-22.94,57.14-20.53,86.81,1.9,18.28,8.26,38.83,20.53,53.74-46.1,26.74-88.59,20.67-276.77,20.67H.06c-.85,42.49,5.98,124.23,57.96,190.77,5.74,7.35,12.04,14.46,18.87,21.31,42.26,42.32,106.11,73.35,201.59,73.44,145.66.13,270.46-78.6,346.37-268.97,24.98.41,90.92,4.48,123.19-57.88.79-1.05,8.21-16.54,8.21-16.54l-12.3-8.27Z"/>
        </svg>
      </div>

      <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
        Docker is ready!
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-2 text-center max-w-sm">
        Build the cross-compilation image to enable building for all platforms
      </p>
      <p className="text-xs text-gray-400 dark:text-gray-500 mb-8 text-center">
        This will download ~800MB and may take several minutes
      </p>

      <div className="flex flex-col items-center gap-2">
        <button
          onClick={onBuildImage}
          className="px-5 py-2 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-500 transition-colors"
        >
          Build Image
        </button>
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

// Author Page
function AuthorPage({
  defaults,
  onDefaultsChange,
  onNext
}: {
  defaults: GlobalDefaults;
  onDefaultsChange: (defaults: GlobalDefaults) => void;
  onNext: () => void;
}) {
  return (
    <PageTemplate
      title="Tell us about yourself"
      subtitle="This information will be used in your apps"
      primaryAction={onNext}
      primaryLabel="Continue"
    >
      <div className="space-y-3 max-w-md mx-auto">
        <div>
          <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">Your Name</label>
          <input
            type="text"
            value={defaults.author.name}
            onChange={(e) => onDefaultsChange({
              ...defaults,
              author: { ...defaults.author, name: e.target.value }
            })}
            placeholder="Jane Developer"
            className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500/20"
          />
        </div>
        <div>
          <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">Company <span className="text-gray-400">(optional)</span></label>
          <input
            type="text"
            value={defaults.author.company}
            onChange={(e) => onDefaultsChange({
              ...defaults,
              author: { ...defaults.author, company: e.target.value }
            })}
            placeholder="Acme Corp"
            className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500/20"
          />
        </div>
      </div>
    </PageTemplate>
  );
}

// Project Defaults Page
function ProjectDefaultsPage({
  defaults,
  onDefaultsChange,
  onNext,
  saving
}: {
  defaults: GlobalDefaults;
  onDefaultsChange: (defaults: GlobalDefaults) => void;
  onNext: () => void;
  saving: boolean;
}) {
  return (
    <PageTemplate
      title="Project defaults"
      subtitle="These will be used when creating new apps"
      primaryAction={onNext}
      primaryLabel={saving ? "Saving..." : "Continue"}
      primaryDisabled={saving}
    >
      <div className="space-y-3 max-w-md mx-auto">
        <div>
          <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">Bundle ID Prefix</label>
          <input
            type="text"
            value={defaults.project.productIdentifierPrefix}
            onChange={(e) => onDefaultsChange({
              ...defaults,
              project: { ...defaults.project, productIdentifierPrefix: e.target.value }
            })}
            placeholder="com.acme"
            className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500/20 font-mono"
          />
          <p className="text-xs text-gray-400 mt-0.5">Example: com.acme.myapp</p>
        </div>

        <div>
          <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">Default Version</label>
          <input
            type="text"
            value={defaults.project.defaultVersion}
            onChange={(e) => onDefaultsChange({
              ...defaults,
              project: { ...defaults.project, defaultVersion: e.target.value }
            })}
            placeholder="0.1.0"
            className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500/20 font-mono"
          />
        </div>

        <div>
          <label className="block text-xs text-gray-600 dark:text-gray-400 mb-1">Preferred Template</label>
          <select
            value={defaults.project.defaultTemplate}
            onChange={(e) => onDefaultsChange({
              ...defaults,
              project: { ...defaults.project, defaultTemplate: e.target.value }
            })}
            className="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-gray-200 focus:border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500/20"
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
    </PageTemplate>
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
      <p className="text-gray-600 dark:text-gray-400 mb-1 text-sm">{label}</p>
      <div className="flex items-center gap-2">
        <code className="flex-1 text-green-600 dark:text-green-400 font-mono text-xs bg-gray-100 dark:bg-gray-900 px-3 py-2 rounded-lg">
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
      transition={{ duration: 0.3 }}
      className="h-full flex flex-col items-center justify-center"
    >
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

      <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">
        You're ready to build!
      </h2>
      <p className="text-gray-500 dark:text-gray-400 mb-8 text-center max-w-sm">
        Your development environment is all set up
      </p>

      <div className="bg-gray-100 dark:bg-gray-900/50 rounded-xl p-5 mb-8 w-full max-w-sm">
        <div className="space-y-4">
          <CopyableCommand command="wails3 init -n myapp" label="Create your first app:" />
          <CopyableCommand command="cd myapp && wails3 dev" label="Start developing:" />
          <CopyableCommand command="wails3 build" label="Build for production:" />
        </div>
      </div>

      <button
        onClick={onClose}
        className="px-5 py-2 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-500 transition-colors"
      >
        Start Building
      </button>

      <a
        href="https://wails.io/docs"
        target="_blank"
        rel="noopener noreferrer"
        className="mt-3 text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
      >
        Read the documentation
      </a>
    </motion.div>
  );
}

// Persistent Docker status indicator
function DockerStatusIndicator({
  status,
  visible
}: {
  status: DockerStatus | null;
  visible: boolean;
}) {
  if (!visible || !status) return null;
  if (!status.installed || !status.running) return null;
  if (status.imageBuilt && status.pullStatus !== 'pulling') return null;

  const isPulling = status.pullStatus === 'pulling';
  const progress = status.pullProgress || 0;

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: 10 }}
      className="bg-white/95 dark:bg-gray-900/95 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg px-3 py-2 backdrop-blur-sm min-w-[200px]"
    >
        <div className="flex items-center gap-3">
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
                  <span className="truncate">Building image...</span>
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
    </motion.div>
  );
}

// Main App
export default function App() {
  const [step, setStep] = useState<OOBEStep>('splash');
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
      defaultVersion: '0.1.0'
    }
  });
  const [savingDefaults, setSavingDefaults] = useState(false);
  const [backgroundDockerStarted, setBackgroundDockerStarted] = useState(false);
  const [theme, setTheme] = useState<Theme>(() => {
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

  useEffect(() => {
    init();
  }, []);

  const init = async () => {
    const state = await getState();
    setSystem(state.system);
  };

  // Handle splash -> checking -> deps result
  const handleSplashNext = async () => {
    setStep('checking');

    // Check dependencies
    const deps = await checkDependencies();
    setDependencies(deps);

    // Determine next step based on deps
    const missingRequired = deps.filter(d => d.required && !d.installed);
    if (missingRequired.length === 0) {
      setStep('deps-ready');
    } else {
      setStep('deps-missing');
    }
  };

  const handleDepsReadyNext = async () => {
    // Check if docker is installed but image is not built
    const dockerDep = dependencies.find(d => d.name === 'docker');
    if (dockerDep?.installed && dockerDep.imageBuilt === false) {
      // Docker installed but no image - ask about cross-compilation
      setStep('cross-platform');
    } else {
      // Either no docker or image already built - go to author
      const loadedDefaults = await getDefaults();
      setDefaults(loadedDefaults);
      setStep('author');
    }
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
    // Check if docker is installed but image is not built
    const dockerDep = dependencies.find(d => d.name === 'docker');
    if (dockerDep?.installed && dockerDep.imageBuilt === false) {
      // Docker installed but no image - ask about cross-compilation
      setStep('cross-platform');
    } else {
      // Either no docker or image already built - go to author
      const loadedDefaults = await getDefaults();
      setDefaults(loadedDefaults);
      setStep('author');
    }
  };

  const handleCrossPlatformYes = async () => {
    // Check Docker status
    const docker = await getDockerStatus();
    setDockerStatus(docker);
    setStep('docker-setup');
  };

  const handleCrossPlatformSkip = async () => {
    // Load defaults and go to author
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setStep('author');
  };

  const handleDockerCheckAgain = async () => {
    const docker = await getDockerStatus();
    setDockerStatus(docker);
  };

  const handleDockerBuildImage = async () => {
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

  const handleDockerContinueBackground = async () => {
    // If build is in progress, let it continue in background
    if (buildingImage || (dockerStatus && dockerStatus.pullStatus === 'pulling')) {
      setBackgroundDockerStarted(true);
    }
    // Load defaults and go to author
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setStep('author');
  };

  const handleDockerSkip = async () => {
    const loadedDefaults = await getDefaults();
    setDefaults(loadedDefaults);
    setStep('author');
  };

  const handleAuthorNext = () => {
    setStep('project-defaults');
  };

  const handleProjectDefaultsNext = async () => {
    setSavingDefaults(true);
    await saveDefaults(defaults);
    setSavingDefaults(false);
    setStep('complete');
  };

  const handleClose = async () => {
    await close();
    window.close();
  };

  // Poll Docker status in background
  useEffect(() => {
    if (backgroundDockerStarted && (buildingImage || (dockerStatus && dockerStatus.pullStatus === 'pulling'))) {
      const poll = async () => {
        const status = await getDockerStatus();
        setDockerStatus(status);
        if (status.pullStatus === 'pulling') {
          setTimeout(poll, 2000);
        } else {
          setBuildingImage(false);
        }
      };
      const timer = setTimeout(poll, 2000);
      return () => clearTimeout(timer);
    }
  }, [backgroundDockerStarted, buildingImage, dockerStatus?.pullStatus]);

  const showDockerIndicator = backgroundDockerStarted && step !== 'docker-setup' && step !== 'splash' && step !== 'checking';

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

        {/* Main content card */}
        <div className="w-full max-w-2xl bg-white dark:bg-gray-900/80 border border-gray-200 dark:border-gray-800 rounded-xl shadow-2xl h-[85vh] flex flex-col overflow-hidden relative z-10">
          <div className="flex-1 flex flex-col p-6 min-h-0">
            <AnimatePresence mode="wait">
              {step === 'splash' && (
                <SplashPage key="splash" onNext={handleSplashNext} />
              )}
              {step === 'checking' && (
                <CheckingPage key="checking" />
              )}
              {step === 'deps-ready' && (
                <DepsReadyPage key="deps-ready" onNext={handleDepsReadyNext} />
              )}
              {step === 'deps-missing' && (
                <DepsMissingPage
                  key="deps-missing"
                  dependencies={dependencies}
                  onRetry={handleDepsMissingRetry}
                  onContinue={handleDepsMissingContinue}
                />
              )}
              {step === 'cross-platform' && (
                <CrossPlatformPage
                  key="cross-platform"
                  onYes={handleCrossPlatformYes}
                  onSkip={handleCrossPlatformSkip}
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
                />
              )}
              {step === 'author' && (
                <AuthorPage
                  key="author"
                  defaults={defaults}
                  onDefaultsChange={setDefaults}
                  onNext={handleAuthorNext}
                />
              )}
              {step === 'project-defaults' && (
                <ProjectDefaultsPage
                  key="project-defaults"
                  defaults={defaults}
                  onDefaultsChange={setDefaults}
                  onNext={handleProjectDefaultsNext}
                  saving={savingDefaults}
                />
              )}
              {step === 'complete' && (
                <CompletePage key="complete" onClose={handleClose} />
              )}
            </AnimatePresence>
          </div>

          {/* Bottom controls - inside content card */}
          <div className="flex-shrink-0 px-6 pb-6 pt-2 flex items-end justify-between">
            <LogoFooter />
            {/* Center: Docker status indicator */}
            <div className="flex-1 flex justify-center">
              <AnimatePresence>
                {showDockerIndicator && (
                  <DockerStatusIndicator
                    status={dockerStatus}
                    visible={showDockerIndicator}
                  />
                )}
              </AnimatePresence>
            </div>
            <BottomControls />
          </div>
        </div>
      </div>
    </ThemeContext.Provider>
  );
}
