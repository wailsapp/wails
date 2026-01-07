import { useState, useEffect, useRef } from 'react';
import { motion } from 'framer-motion';
import type { SigningStatus } from '../types';
import { getSigningStatus } from '../api';

const pageVariants = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 }
};

type Platform = 'darwin' | 'windows' | 'linux';

const platformInfo: Record<Platform, { name: string; icon: string }> = {
  darwin: { name: 'macOS', icon: '🍎' },
  windows: { name: 'Windows', icon: '🪟' },
  linux: { name: 'Linux', icon: '🐧' }
};

interface Props {
  onNext: () => void;
  onSkip: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}

export default function SigningStep({ onNext, onSkip, onBack, canGoBack }: Props) {
  const [status, setStatus] = useState<SigningStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedPlatform, setSelectedPlatform] = useState<Platform>('darwin');
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
    loadStatus();
  }, []);

  const loadStatus = async () => {
    try {
      const s = await getSigningStatus();
      setStatus(s);
    } catch (e) {
      console.error('Failed to load signing status:', e);
    } finally {
      setLoading(false);
    }
  };

  const renderPlatformStatus = () => {
    if (!status) return null;
    
    if (selectedPlatform === 'darwin') {
      const darwin = status.darwin;
      return (
        <div className="space-y-4">
          <div className="flex items-center gap-3 p-4 rounded-lg bg-gray-100 dark:bg-gray-900/50">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center ${darwin.hasIdentity ? 'bg-green-500/20' : 'bg-gray-200 dark:bg-gray-800'}`}>
              {darwin.hasIdentity ? (
                <svg className="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <svg className="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              )}
            </div>
            <div className="flex-1">
              <div className="text-sm font-medium text-gray-900 dark:text-white">Code Signing Identity</div>
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {darwin.hasIdentity ? darwin.identity : 'Not configured'}
              </div>
            </div>
            {darwin.configSource && (
              <span className="text-xs px-2 py-1 rounded-full bg-gray-200 dark:bg-gray-800 text-gray-600 dark:text-gray-400">
                {darwin.configSource}
              </span>
            )}
          </div>

          <div className="flex items-center gap-3 p-4 rounded-lg bg-gray-100 dark:bg-gray-900/50">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center ${darwin.hasNotarization ? 'bg-green-500/20' : 'bg-gray-200 dark:bg-gray-800'}`}>
              {darwin.hasNotarization ? (
                <svg className="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <svg className="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              )}
            </div>
            <div className="flex-1">
              <div className="text-sm font-medium text-gray-900 dark:text-white">Notarization</div>
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {darwin.hasNotarization 
                  ? `Team ID: ${darwin.teamID || 'Configured'}` 
                  : 'Not configured'}
              </div>
            </div>
          </div>

          {darwin.identities && darwin.identities.length > 1 && (
            <div className="text-xs text-gray-500 dark:text-gray-400 p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
              <span className="font-medium">{darwin.identities.length} signing identities</span> found in keychain
            </div>
          )}
        </div>
      );
    }

    if (selectedPlatform === 'windows') {
      const windows = status.windows;
      return (
        <div className="space-y-4">
          <div className="flex items-center gap-3 p-4 rounded-lg bg-gray-100 dark:bg-gray-900/50">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center ${windows.hasCertificate ? 'bg-green-500/20' : 'bg-gray-200 dark:bg-gray-800'}`}>
              {windows.hasCertificate ? (
                <svg className="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <svg className="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              )}
            </div>
            <div className="flex-1">
              <div className="text-sm font-medium text-gray-900 dark:text-white">Code Signing Certificate</div>
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {windows.hasCertificate 
                  ? `Type: ${windows.certificateType}` 
                  : 'Not configured'}
              </div>
            </div>
            {windows.configSource && (
              <span className="text-xs px-2 py-1 rounded-full bg-gray-200 dark:bg-gray-800 text-gray-600 dark:text-gray-400">
                {windows.configSource}
              </span>
            )}
          </div>

          <div className="flex items-center gap-3 p-4 rounded-lg bg-gray-100 dark:bg-gray-900/50">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center ${windows.hasSignTool ? 'bg-green-500/20' : 'bg-gray-200 dark:bg-gray-800'}`}>
              {windows.hasSignTool ? (
                <svg className="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <svg className="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              )}
            </div>
            <div className="flex-1">
              <div className="text-sm font-medium text-gray-900 dark:text-white">SignTool</div>
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {windows.hasSignTool ? 'Available' : 'Not found (Windows SDK required)'}
              </div>
            </div>
          </div>

          {windows.timestampServer && (
            <div className="text-xs text-gray-500 dark:text-gray-400 p-3 rounded-lg bg-gray-50 dark:bg-gray-800/50">
              Timestamp server: <code className="font-mono">{windows.timestampServer}</code>
            </div>
          )}
        </div>
      );
    }

    if (selectedPlatform === 'linux') {
      const linux = status.linux;
      return (
        <div className="space-y-4">
          <div className="flex items-center gap-3 p-4 rounded-lg bg-gray-100 dark:bg-gray-900/50">
            <div className={`w-8 h-8 rounded-full flex items-center justify-center ${linux.hasGpgKey ? 'bg-green-500/20' : 'bg-gray-200 dark:bg-gray-800'}`}>
              {linux.hasGpgKey ? (
                <svg className="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <svg className="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              )}
            </div>
            <div className="flex-1">
              <div className="text-sm font-medium text-gray-900 dark:text-white">GPG Signing Key</div>
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {linux.hasGpgKey 
                  ? `Key ID: ${linux.gpgKeyID}` 
                  : 'Not configured'}
              </div>
            </div>
            {linux.configSource && (
              <span className="text-xs px-2 py-1 rounded-full bg-gray-200 dark:bg-gray-800 text-gray-600 dark:text-gray-400">
                {linux.configSource}
              </span>
            )}
          </div>
        </div>
      );
    }

    return null;
  };

  const getOverallStatus = () => {
    if (!status) return { configured: 0, total: 3 };
    let configured = 0;
    if (status.darwin.hasIdentity) configured++;
    if (status.windows.hasCertificate) configured++;
    if (status.linux.hasGpgKey) configured++;
    return { configured, total: 3 };
  };

  const overallStatus = getOverallStatus();

  return (
    <motion.main
      variants={pageVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className="flex-1 flex flex-col"
      aria-labelledby="signing-title"
    >
      <header className="text-center mb-6 flex-shrink-0 px-10 pt-10">
        <h1
          ref={headingRef}
          id="signing-title"
          className="text-2xl font-semibold text-gray-900 dark:text-white mb-1.5 tracking-tight focus:outline-none"
          tabIndex={-1}
        >
          Code Signing
        </h1>
        <p className="text-base text-gray-500 dark:text-gray-400">
          {overallStatus.configured > 0 
            ? `${overallStatus.configured} of ${overallStatus.total} platforms configured`
            : 'Sign your apps for distribution'}
        </p>
      </header>

      <div className="flex-1 overflow-y-auto scrollbar-thin min-h-0 px-10">
        {loading ? (
          <div className="flex items-center justify-center h-48">
            <motion.div
              className="w-8 h-8 border-2 border-gray-300 dark:border-gray-600 border-t-red-500 rounded-full"
              animate={{ rotate: 360 }}
              transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
            />
          </div>
        ) : (
          <div className="max-w-xl mx-auto">
            <div className="flex gap-2 mb-6" role="tablist">
              {(['darwin', 'windows', 'linux'] as Platform[]).map((platform) => {
                const info = platformInfo[platform];
                const isActive = selectedPlatform === platform;
                const hasConfig = status && (
                  (platform === 'darwin' && status.darwin.hasIdentity) ||
                  (platform === 'windows' && status.windows.hasCertificate) ||
                  (platform === 'linux' && status.linux.hasGpgKey)
                );
                
                return (
                  <button
                    key={platform}
                    role="tab"
                    aria-selected={isActive}
                    onClick={() => setSelectedPlatform(platform)}
                    className={`flex-1 flex items-center justify-center gap-2 px-4 py-3 rounded-lg text-sm font-medium transition-all ${
                      isActive
                        ? 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white'
                        : 'text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800/50'
                    }`}
                  >
                    <span className="text-lg">{info.icon}</span>
                    <span>{info.name}</span>
                    {hasConfig && (
                      <span className="w-2 h-2 rounded-full bg-green-500" />
                    )}
                  </button>
                );
              })}
            </div>

            <div role="tabpanel" aria-label={`${platformInfo[selectedPlatform].name} signing status`}>
              {renderPlatformStatus()}
            </div>

            <p className="text-xs text-gray-500 dark:text-gray-400 mt-6 text-center">
              Code signing ensures your app is trusted and hasn't been tampered with
            </p>
          </div>
        )}
      </div>

      <div className="flex-shrink-0 pt-4 pb-6 flex flex-col items-center gap-1.5">
        <div className="flex items-center gap-3">
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
            className="px-5 py-2 rounded-lg text-sm font-medium transition-colors border border-red-500 text-red-600 dark:text-red-400 hover:bg-red-500/10 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
          >
            Continue
          </button>
        </div>
        <button
          onClick={onSkip}
          className="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 rounded"
        >
          Set up later
        </button>
      </div>
    </motion.main>
  );
}
