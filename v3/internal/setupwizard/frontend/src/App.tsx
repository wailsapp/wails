import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { DependencyStatus, SystemInfo, DockerStatus } from './types';
import { checkDependencies, getState, getDockerStatus, buildDockerImage, close, installDependency } from './api';
import WailsLogo from './components/WailsLogo';

type Step = 'welcome' | 'dependencies' | 'docker' | 'complete';

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
    <div className="flex items-center justify-center gap-1 text-xs text-gray-500 mb-6">
      {steps.map((step, i) => (
        <div key={step.id} className="flex items-center">
          <span className={i <= currentIndex ? 'text-white' : 'text-gray-600'}>
            {step.label}
          </span>
          {i < steps.length - 1 && (
            <span className="mx-2 text-gray-700">&rsaquo;</span>
          )}
        </div>
      ))}
    </div>
  );
}

// Wizard footer with navigation buttons
function WizardFooter({
  onBack,
  onNext,
  onCancel,
  nextLabel = 'Next',
  backLabel = 'Back',
  showBack = true,
  nextDisabled = false
}: {
  onBack?: () => void;
  onNext: () => void;
  onCancel?: () => void;
  nextLabel?: string;
  backLabel?: string;
  showBack?: boolean;
  nextDisabled?: boolean;
}) {
  return (
    <div className="flex justify-between items-center pt-6 mt-6 border-t border-gray-800">
      <div>
        {onCancel && (
          <button
            onClick={onCancel}
            className="px-4 py-2 text-gray-500 hover:text-gray-300 transition-colors"
          >
            Cancel
          </button>
        )}
      </div>
      <div className="flex gap-3">
        {showBack && onBack && (
          <button
            onClick={onBack}
            className="px-5 py-2 rounded-lg bg-gray-800 text-gray-300 hover:bg-gray-700 transition-colors"
          >
            {backLabel}
          </button>
        )}
        <button
          onClick={onNext}
          disabled={nextDisabled}
          className={`px-5 py-2 rounded-lg font-medium transition-colors ${
            nextDisabled
              ? 'bg-gray-700 text-gray-500 cursor-not-allowed'
              : 'bg-red-600 text-white hover:bg-red-500'
          }`}
        >
          {nextLabel}
        </button>
      </div>
    </div>
  );
}

// Welcome Page
function WelcomePage({ system, onNext, onCancel }: { system: SystemInfo | null; onNext: () => void; onCancel: () => void }) {
  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.2 }}
    >
      <div className="text-center mb-8">
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ delay: 0.1 }}
          className="inline-block"
        >
          <WailsLogo size={80} />
        </motion.div>
        <h1 className="text-2xl font-bold mt-4 mb-2">Welcome to Wails Setup</h1>
        <p className="text-gray-400">
          This wizard will help you set up your development environment.
        </p>
      </div>

      {system && (
        <div className="bg-gray-900/50 rounded-lg p-4 mb-6">
          <h3 className="text-sm font-medium text-gray-400 mb-3">System Information</h3>
          <div className="grid grid-cols-2 gap-y-2 text-sm">
            <span className="text-gray-500">Operating System</span>
            <span className="text-gray-200">{system.osName || system.os} ({system.arch})</span>
            <span className="text-gray-500">Wails Version</span>
            <span className="text-gray-200">v{system.wailsVersion}</span>
            <span className="text-gray-500">Go Version</span>
            <span className="text-gray-200">{system.goVersion}</span>
          </div>
        </div>
      )}

      <div className="bg-gray-900/50 rounded-lg p-4">
        <h3 className="text-sm font-medium text-gray-400 mb-2">Setup will check:</h3>
        <ul className="text-sm text-gray-300 space-y-1">
          <li className="flex items-center gap-2">
            <span className="text-gray-600">•</span>
            Required build dependencies (GTK, WebKit, GCC)
          </li>
          <li className="flex items-center gap-2">
            <span className="text-gray-600">•</span>
            Optional tools (npm, Docker)
          </li>
          <li className="flex items-center gap-2">
            <span className="text-gray-600">•</span>
            Cross-compilation capabilities
          </li>
        </ul>
      </div>

      <WizardFooter
        onNext={onNext}
        onCancel={onCancel}
        nextLabel="Check Dependencies"
        showBack={false}
      />
    </motion.div>
  );
}

// Dependency row component
function DependencyRow({
  dep,
  onInstall,
  installing
}: {
  dep: DependencyStatus;
  onInstall?: (cmd: string) => void;
  installing: boolean;
}) {
  const [copied, setCopied] = useState(false);

  const copyCommand = () => {
    if (dep.installCommand) {
      navigator.clipboard.writeText(dep.installCommand);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const isSystemCommand = dep.installCommand?.startsWith('sudo ');

  return (
    <div className="flex items-start gap-3 py-3 border-b border-gray-800/50 last:border-0">
      {/* Status icon */}
      <div className="mt-0.5">
        {dep.installed ? (
          <div className="w-5 h-5 rounded-full bg-green-500/20 flex items-center justify-center">
            <svg className="w-3 h-3 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
            </svg>
          </div>
        ) : dep.required ? (
          <div className="w-5 h-5 rounded-full bg-red-500/20 flex items-center justify-center">
            <svg className="w-3 h-3 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
        ) : (
          <div className="w-5 h-5 rounded-full bg-gray-600/20 flex items-center justify-center">
            <div className="w-2 h-2 rounded-full bg-gray-500" />
          </div>
        )}
      </div>

      {/* Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className={`font-medium ${dep.installed ? 'text-white' : dep.required ? 'text-red-300' : 'text-gray-400'}`}>
            {dep.name}
          </span>
          {dep.version && (
            <span className="text-xs text-gray-500 font-mono">{dep.version}</span>
          )}
          {dep.required && !dep.installed && (
            <span className="text-xs text-red-400 bg-red-500/10 px-1.5 py-0.5 rounded">Required</span>
          )}
        </div>
        {dep.message && (
          <p className="text-xs text-gray-500 mt-0.5">{dep.message}</p>
        )}

        {/* Install command */}
        {!dep.installed && dep.installCommand && (
          <div className="mt-2">
            {dep.helpUrl ? (
              <a
                href={dep.helpUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1 text-xs text-blue-400 hover:text-blue-300"
              >
                {dep.installCommand}
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
              </a>
            ) : (
              <div className="flex items-center gap-2">
                <code className="text-xs bg-gray-900 text-gray-300 px-2 py-1 rounded font-mono">
                  {dep.installCommand}
                </code>
                <button
                  onClick={copyCommand}
                  className="text-xs text-gray-500 hover:text-gray-300 transition-colors"
                  title="Copy command"
                >
                  {copied ? (
                    <svg className="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  ) : (
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                  )}
                </button>
                {isSystemCommand && onInstall && (
                  <button
                    onClick={() => onInstall(dep.installCommand!)}
                    disabled={installing}
                    className="text-xs px-2 py-1 rounded bg-blue-500/20 text-blue-400 hover:bg-blue-500/30 transition-colors disabled:opacity-50"
                  >
                    {installing ? 'Installing...' : 'Install'}
                  </button>
                )}
              </div>
            )}
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
  onCancel,
  onRefresh
}: {
  dependencies: DependencyStatus[];
  onNext: () => void;
  onBack: () => void;
  onCancel: () => void;
  onRefresh: () => void;
}) {
  const [installing, setInstalling] = useState(false);

  const required = dependencies.filter(d => d.required);
  const optional = dependencies.filter(d => !d.required);
  const missingRequired = required.filter(d => !d.installed);
  const allRequiredInstalled = missingRequired.length === 0;

  const handleInstall = async (command: string) => {
    setInstalling(true);
    try {
      const result = await installDependency(command);
      if (result.success) {
        // Refresh dependencies after install
        onRefresh();
      } else {
        alert(`Installation failed: ${result.error || result.output}`);
      }
    } catch {
      alert('Failed to run install command');
    }
    setInstalling(false);
  };

  return (
    <motion.div
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.2 }}
    >
      <div className="mb-6">
        <h2 className="text-xl font-bold mb-1">System Dependencies</h2>
        <p className="text-sm text-gray-400">
          The following dependencies are needed to build Wails applications.
        </p>
      </div>

      {/* Required Dependencies */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-sm font-medium text-gray-400">Required</h3>
          {!allRequiredInstalled && (
            <span className="text-xs text-red-400">
              {missingRequired.length} missing
            </span>
          )}
        </div>
        <div className="bg-gray-900/50 rounded-lg px-4">
          {required.map(dep => (
            <DependencyRow
              key={dep.name}
              dep={dep}
              onInstall={handleInstall}
              installing={installing}
            />
          ))}
        </div>
      </div>

      {/* Optional Dependencies */}
      {optional.length > 0 && (
        <div className="mb-6">
          <h3 className="text-sm font-medium text-gray-400 mb-2">Optional</h3>
          <div className="bg-gray-900/50 rounded-lg px-4">
            {optional.map(dep => (
              <DependencyRow
                key={dep.name}
                dep={dep}
                onInstall={handleInstall}
                installing={installing}
              />
            ))}
          </div>
        </div>
      )}

      {/* Status Summary */}
      <div className={`rounded-lg p-3 ${allRequiredInstalled ? 'bg-green-500/10 border border-green-500/20' : 'bg-yellow-500/10 border border-yellow-500/20'}`}>
        {allRequiredInstalled ? (
          <div className="flex items-center gap-2 text-green-400 text-sm">
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            All required dependencies are installed. You can proceed.
          </div>
        ) : (
          <div className="flex items-center gap-2 text-yellow-400 text-sm">
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            Install missing dependencies before continuing, or proceed anyway.
          </div>
        )}
      </div>

      <WizardFooter
        onBack={onBack}
        onNext={onNext}
        onCancel={onCancel}
        nextLabel={allRequiredInstalled ? 'Next' : 'Continue Anyway'}
      />
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
        <h2 className="text-xl font-bold mb-1">Cross-Platform Builds</h2>
        <p className="text-sm text-gray-400">
          Docker enables building for macOS, Windows, and Linux from any platform.
        </p>
      </div>

      <div className="bg-gray-900/50 rounded-lg p-4 mb-6">
        <div className="flex items-start gap-4">
          <div className="w-12 h-12 rounded-xl bg-blue-500/20 flex items-center justify-center flex-shrink-0">
            <svg className="w-6 h-6 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M13.983 11.078h2.119a.186.186 0 00.186-.185V9.006a.186.186 0 00-.186-.186h-2.119a.186.186 0 00-.185.186v1.887c0 .102.083.185.185.185zm-2.954-5.43h2.118a.186.186 0 00.186-.186V3.574a.186.186 0 00-.186-.185h-2.118a.186.186 0 00-.185.185v1.888c0 .102.082.185.185.186zm0 2.716h2.118a.187.187 0 00.186-.186V6.29a.186.186 0 00-.186-.185h-2.118a.186.186 0 00-.185.185v1.888c0 .102.082.185.185.186zm-2.93 0h2.12a.186.186 0 00.184-.186V6.29a.186.186 0 00-.185-.185H8.1a.186.186 0 00-.185.185v1.888c0 .102.083.185.185.186zm-2.964 0h2.119a.186.186 0 00.185-.186V6.29a.186.186 0 00-.185-.185H5.136a.186.186 0 00-.186.185v1.888c0 .102.084.185.186.186zm5.893 2.715h2.118a.186.186 0 00.186-.185V9.006a.186.186 0 00-.186-.186h-2.118a.186.186 0 00-.185.186v1.887c0 .102.082.185.185.185zm-2.93 0h2.12a.185.185 0 00.184-.185V9.006a.185.185 0 00-.184-.186h-2.12a.185.185 0 00-.184.186v1.887c0 .102.083.185.185.185zm-2.964 0h2.119a.185.185 0 00.185-.185V9.006a.185.185 0 00-.185-.186h-2.12a.185.185 0 00-.184.186v1.887c0 .102.083.185.185.185zm-2.92 0h2.12a.185.185 0 00.184-.185V9.006a.185.185 0 00-.184-.186h-2.12a.186.186 0 00-.185.186v1.887c0 .102.084.185.185.185z"/>
            </svg>
          </div>

          <div className="flex-1">
            <h3 className="font-medium text-white mb-1">Docker Status</h3>

            {!dockerStatus ? (
              <div className="text-sm text-gray-400">Checking Docker...</div>
            ) : !dockerStatus.installed ? (
              <div>
                <div className="flex items-center gap-2 text-yellow-400 text-sm mb-2">
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
                  className="inline-flex items-center gap-1 text-xs text-blue-400 hover:text-blue-300"
                >
                  Install Docker Desktop
                  <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              </div>
            ) : !dockerStatus.running ? (
              <div>
                <div className="flex items-center gap-2 text-yellow-400 text-sm mb-2">
                  <span className="w-2 h-2 rounded-full bg-yellow-500" />
                  Installed but not running
                </div>
                <p className="text-xs text-gray-500">
                  Start Docker Desktop to enable cross-platform builds.
                </p>
              </div>
            ) : dockerStatus.imageBuilt ? (
              <div>
                <div className="flex items-center gap-2 text-green-400 text-sm mb-2">
                  <span className="w-2 h-2 rounded-full bg-green-500" />
                  Ready for cross-platform builds
                </div>
                <p className="text-xs text-gray-500">
                  Docker {dockerStatus.version} • wails-cross image installed
                </p>
              </div>
            ) : buildingImage ? (
              <div>
                <div className="flex items-center gap-2 text-blue-400 text-sm mb-2">
                  <motion.span
                    className="w-3 h-3 border-2 border-blue-400 border-t-transparent rounded-full"
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                  />
                  Building wails-cross image... {dockerStatus.pullProgress}%
                </div>
                <div className="h-1.5 bg-gray-700 rounded-full overflow-hidden">
                  <motion.div
                    className="h-full bg-blue-500"
                    animate={{ width: `${dockerStatus.pullProgress}%` }}
                  />
                </div>
              </div>
            ) : (
              <div>
                <div className="flex items-center gap-2 text-gray-400 text-sm mb-2">
                  <span className="w-2 h-2 rounded-full bg-gray-500" />
                  Cross-compilation image not installed
                </div>
                <p className="text-xs text-gray-500 mb-3">
                  Docker {dockerStatus.version} is running. Build the wails-cross image to enable cross-platform builds.
                </p>
                <button
                  onClick={onBuildImage}
                  className="text-sm px-4 py-2 rounded-lg bg-blue-500/20 text-blue-400 hover:bg-blue-500/30 transition-colors border border-blue-500/30"
                >
                  Build Cross-Compilation Image
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="bg-gray-900/50 rounded-lg p-4">
        <h3 className="text-sm font-medium text-gray-400 mb-2">What you can build:</h3>
        <div className="grid grid-cols-3 gap-4 text-center text-sm">
          <div className="py-2">
            <div className="text-lg mb-1">🍎</div>
            <div className="text-gray-300">macOS</div>
            <div className="text-xs text-gray-500">.app / .dmg</div>
          </div>
          <div className="py-2">
            <div className="text-lg mb-1">🪟</div>
            <div className="text-gray-300">Windows</div>
            <div className="text-xs text-gray-500">.exe / .msi</div>
          </div>
          <div className="py-2">
            <div className="text-lg mb-1">🐧</div>
            <div className="text-gray-300">Linux</div>
            <div className="text-xs text-gray-500">.deb / .rpm / AppImage</div>
          </div>
        </div>
      </div>

      <WizardFooter
        onBack={onBack}
        onNext={onNext}
        onCancel={onCancel}
        nextLabel="Finish"
      />
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
      <p className="text-gray-500 mb-1">{label}</p>
      <div className="flex items-center gap-2">
        <code className="flex-1 text-green-400 font-mono text-xs bg-gray-900 px-2 py-1 rounded">
          {command}
        </code>
        <button
          onClick={copyCommand}
          className="text-gray-500 hover:text-gray-300 transition-colors p-1"
          title="Copy command"
        >
          {copied ? (
            <svg className="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
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
        <svg className="w-8 h-8 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
        </svg>
      </motion.div>

      <h2 className="text-2xl font-bold mb-2">Setup Complete</h2>
      <p className="text-gray-400 mb-8">
        Your development environment is ready to use.
      </p>

      <div className="bg-gray-900/50 rounded-lg p-4 text-left mb-6 max-w-sm mx-auto">
        <h3 className="text-sm font-medium text-gray-400 mb-3">Next Steps</h3>
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
  const [step, setStep] = useState<Step>('welcome');
  const [dependencies, setDependencies] = useState<DependencyStatus[]>([]);
  const [system, setSystem] = useState<SystemInfo | null>(null);
  const [dockerStatus, setDockerStatus] = useState<DockerStatus | null>(null);
  const [buildingImage, setBuildingImage] = useState(false);

  const steps: { id: Step; label: string }[] = [
    { id: 'welcome', label: 'Welcome' },
    { id: 'dependencies', label: 'Dependencies' },
    { id: 'docker', label: 'Docker' },
    { id: 'complete', label: 'Complete' },
  ];

  useEffect(() => {
    init();
  }, []);

  const init = async () => {
    const state = await getState();
    setSystem(state.system);
  };

  const refreshDependencies = async () => {
    const deps = await checkDependencies();
    setDependencies(deps);
  };

  const handleNext = async () => {
    if (step === 'welcome') {
      const deps = await checkDependencies();
      setDependencies(deps);
      setStep('dependencies');
    } else if (step === 'dependencies') {
      const dockerDep = dependencies.find(d => d.name === 'docker');
      if (dockerDep?.installed) {
        const docker = await getDockerStatus();
        setDockerStatus(docker);
      }
      setStep('docker');
    } else if (step === 'docker') {
      setStep('complete');
    }
  };

  const handleBack = () => {
    if (step === 'dependencies') setStep('welcome');
    else if (step === 'docker') setStep('dependencies');
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

  const handleClose = async () => {
    await close();
    window.close();
  };

  const handleCancel = handleClose;

  return (
    <div className="min-h-screen bg-[#0f0f0f] flex items-center justify-center p-4">
      <div className="w-full max-w-lg">
        {/* Wizard container */}
        <div className="bg-gray-900/80 border border-gray-800 rounded-xl p-6 shadow-2xl">
          <StepIndicator steps={steps} currentStep={step} />

          <AnimatePresence mode="wait">
            {step === 'welcome' && (
              <WelcomePage
                key="welcome"
                system={system}
                onNext={handleNext}
                onCancel={handleCancel}
              />
            )}
            {step === 'dependencies' && (
              <DependenciesPage
                key="dependencies"
                dependencies={dependencies}
                onNext={handleNext}
                onBack={handleBack}
                onCancel={handleCancel}
                onRefresh={refreshDependencies}
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

        {/* Footer */}
        <div className="text-center mt-4 text-xs text-gray-600">
          Wails • Build cross-platform apps with Go
        </div>
      </div>
    </div>
  );
}
