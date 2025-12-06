import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { DependencyStatus, SystemInfo } from './types';
import { checkDependencies, getState } from './api';
import WailsLogo from './components/WailsLogo';

type CheckState = 'idle' | 'checking' | 'complete';

interface DependencyWithState extends DependencyStatus {
  checkState: CheckState;
}

export default function App() {
  const [dependencies, setDependencies] = useState<DependencyWithState[]>([]);
  const [system, setSystem] = useState<SystemInfo | null>(null);
  const [started, setStarted] = useState(false);
  const [allComplete, setAllComplete] = useState(false);

  // Load system info on mount
  useEffect(() => {
    getState().then((s) => {
      setSystem(s.system);
    });
  }, []);

  // Start checking after logo animation
  useEffect(() => {
    const timer = setTimeout(() => {
      setStarted(true);
      startDependencyCheck();
    }, 1500);
    return () => clearTimeout(timer);
  }, []);

  const startDependencyCheck = async () => {
    // First get the list of dependencies
    const deps = await checkDependencies();

    // Initialize all as 'idle'
    const depsWithState: DependencyWithState[] = deps.map(d => ({
      ...d,
      checkState: 'idle'
    }));
    setDependencies(depsWithState);

    // Animate each dependency check one by one with staggered timing
    for (let i = 0; i < depsWithState.length; i++) {
      await new Promise(resolve => setTimeout(resolve, 300));

      setDependencies(prev => prev.map((d, idx) =>
        idx === i ? { ...d, checkState: 'checking' } : d
      ));

      // Simulate check time (the actual check is already done)
      await new Promise(resolve => setTimeout(resolve, 400 + Math.random() * 300));

      setDependencies(prev => prev.map((d, idx) =>
        idx === i ? { ...d, checkState: 'complete' } : d
      ));
    }

    // All complete
    await new Promise(resolve => setTimeout(resolve, 500));
    setAllComplete(true);
  };

  const allInstalled = dependencies.length > 0 &&
    dependencies.every(d => d.installed || !d.required);

  const getStatusIcon = (dep: DependencyWithState) => {
    if (dep.checkState === 'idle') {
      return (
        <div className="w-5 h-5 rounded-full border-2 border-gray-600" />
      );
    }

    if (dep.checkState === 'checking') {
      return (
        <motion.div
          className="w-5 h-5 rounded-full border-2 border-red-500 border-t-transparent"
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
        />
      );
    }

    // Complete
    if (dep.installed) {
      return (
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          className="w-5 h-5 rounded-full bg-green-500 flex items-center justify-center"
        >
          <motion.svg
            className="w-3 h-3 text-white"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            initial={{ pathLength: 0 }}
            animate={{ pathLength: 1 }}
            transition={{ duration: 0.3 }}
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
          </motion.svg>
        </motion.div>
      );
    } else {
      return (
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          className="w-5 h-5 rounded-full bg-red-500 flex items-center justify-center"
        >
          <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </motion.div>
      );
    }
  };

  return (
    <div className="min-h-screen bg-[#0f0f0f] grid-bg relative overflow-hidden">
      {/* Radial glow background */}
      <div className="absolute inset-0 radial-glow pointer-events-none" />

      {/* Main content */}
      <div className="relative z-10 min-h-screen flex flex-col items-center justify-center px-8">
        {/* Logo */}
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.6, ease: 'easeOut' }}
        >
          <WailsLogo size={160} />
        </motion.div>

        {/* Title */}
        <motion.div
          className="mt-6 text-center"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3, duration: 0.5 }}
        >
          <h1 className="text-3xl font-bold text-white">
            Wails Setup
          </h1>
          <p className="text-gray-400 mt-2">
            {system?.os && `${system.os}/${system.arch}`}
            {system?.wailsVersion && ` • v${system.wailsVersion}`}
          </p>
        </motion.div>

        {/* Dependencies section */}
        <motion.div
          className="mt-10 w-full max-w-md"
          initial={{ opacity: 0 }}
          animate={{ opacity: started ? 1 : 0 }}
          transition={{ duration: 0.4 }}
        >
          {/* Header */}
          <motion.div
            className="flex items-center gap-2 mb-4"
            initial={{ opacity: 0, x: -10 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.2 }}
          >
            <span className="text-sm font-medium text-gray-400 uppercase tracking-wider">
              Checking Dependencies
            </span>
            {!allComplete && dependencies.length > 0 && (
              <motion.span
                className="text-xs text-gray-500"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
              >
                ({dependencies.filter(d => d.checkState === 'complete').length}/{dependencies.length})
              </motion.span>
            )}
          </motion.div>

          {/* Dependency list */}
          <div className="space-y-2">
            <AnimatePresence mode="popLayout">
              {dependencies.map((dep, index) => (
                <motion.div
                  key={dep.name}
                  className="flex items-center gap-4 p-4 rounded-xl bg-gray-800/50 border border-gray-700/50"
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{
                    delay: index * 0.1,
                    duration: 0.3,
                    ease: 'easeOut'
                  }}
                  layout
                >
                  {/* Status icon */}
                  <div className="flex-shrink-0">
                    {getStatusIcon(dep)}
                  </div>

                  {/* Name and version */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="text-white font-medium">{dep.name}</span>
                      {dep.checkState === 'complete' && dep.version && (
                        <motion.span
                          className="text-xs text-gray-500 font-mono"
                          initial={{ opacity: 0 }}
                          animate={{ opacity: 1 }}
                        >
                          v{dep.version}
                        </motion.span>
                      )}
                    </div>

                    {/* Message for missing deps */}
                    {dep.checkState === 'complete' && !dep.installed && dep.message && (
                      <motion.p
                        className="text-xs text-red-400 mt-1 truncate"
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                      >
                        {dep.message}
                      </motion.p>
                    )}
                  </div>

                  {/* Required badge */}
                  {dep.required && dep.checkState === 'complete' && !dep.installed && (
                    <motion.span
                      className="text-xs px-2 py-1 rounded bg-red-500/20 text-red-400"
                      initial={{ opacity: 0, scale: 0.8 }}
                      animate={{ opacity: 1, scale: 1 }}
                    >
                      Required
                    </motion.span>
                  )}
                </motion.div>
              ))}
            </AnimatePresence>
          </div>

          {/* Summary */}
          <AnimatePresence>
            {allComplete && (
              <motion.div
                className="mt-6"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
              >
                {allInstalled ? (
                  <div className="text-center p-4 rounded-xl bg-green-500/10 border border-green-500/30">
                    <motion.div
                      initial={{ scale: 0 }}
                      animate={{ scale: 1 }}
                      transition={{ type: 'spring', stiffness: 200, damping: 15 }}
                    >
                      <div className="text-3xl mb-2">✓</div>
                    </motion.div>
                    <p className="text-green-400 font-medium">
                      All dependencies installed!
                    </p>
                    <p className="text-gray-400 text-sm mt-1">
                      You're ready to build Wails applications.
                    </p>
                  </div>
                ) : (
                  <div className="text-center p-4 rounded-xl bg-yellow-500/10 border border-yellow-500/30">
                    <div className="text-3xl mb-2">⚠</div>
                    <p className="text-yellow-400 font-medium">
                      Some dependencies are missing
                    </p>
                    <p className="text-gray-400 text-sm mt-1">
                      Install the missing dependencies to continue.
                    </p>
                  </div>
                )}

                {/* Close button */}
                <motion.button
                  className="mt-6 w-full btn-primary"
                  onClick={() => window.close()}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: 0.4 }}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  Close
                </motion.button>
              </motion.div>
            )}
          </AnimatePresence>
        </motion.div>

        {/* Loading dots while fetching deps */}
        {started && dependencies.length === 0 && (
          <motion.div
            className="mt-10 flex space-x-2"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
          >
            {[0, 1, 2].map((i) => (
              <motion.div
                key={i}
                className="w-2 h-2 rounded-full bg-red-500/50"
                animate={{
                  scale: [1, 1.3, 1],
                  opacity: [0.5, 1, 0.5]
                }}
                transition={{
                  duration: 1,
                  repeat: Infinity,
                  delay: i * 0.15
                }}
              />
            ))}
          </motion.div>
        )}
      </div>
    </div>
  );
}
