import { useState, useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { SigningStatus, SigningDefaults } from '../types';
import { getSigningStatus, getSigning, saveSigning, getState } from '../api';

const pageVariants = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 }
};

type Platform = 'darwin' | 'windows' | 'linux';
type HostOS = 'darwin' | 'windows' | 'linux';

interface Props {
  onNext: () => void;
  onSkip: () => void;
  onBack?: () => void;
  canGoBack?: boolean;
}

export default function SigningStep({ onNext, onSkip, onBack, canGoBack }: Props) {
  const [status, setStatus] = useState<SigningStatus | null>(null);
  const [config, setConfig] = useState<SigningDefaults | null>(null);
  const [loading, setLoading] = useState(true);
  const [hostOS, setHostOS] = useState<HostOS>('linux');
  const [selectedPlatform, setSelectedPlatform] = useState<Platform>('darwin');
  const [configuring, setConfiguring] = useState(false);
  const [saving, setSaving] = useState(false);
  const headingRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => {
    headingRef.current?.focus();
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [s, c, state] = await Promise.all([getSigningStatus(), getSigning(), getState()]);
      setStatus(s);
      setConfig(c || { darwin: {}, windows: {}, linux: {} });
      if (state.system?.os) {
        setHostOS(state.system.os as HostOS);
      }
    } catch (e) {
      console.error('Failed to load signing data:', e);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!config) return;
    setSaving(true);
    try {
      await saveSigning(config);
      await loadData();
      setConfiguring(false);
    } catch (e) {
      console.error('Failed to save signing config:', e);
    } finally {
      setSaving(false);
    }
  };

  const renderConfigForm = () => {
    if (!config) return null;

    if (selectedPlatform === 'darwin') {
      const isOnMac = hostOS === 'darwin';
      
      return (
        <div className="space-y-4">
          {!isOnMac && (
            <div className="p-3 rounded-lg bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 text-sm">
              <p className="text-amber-800 dark:text-amber-200 font-medium mb-1">Cross-platform signing</p>
              <p className="text-amber-700 dark:text-amber-300 text-xs">
                You can sign macOS apps from {hostOS === 'linux' ? 'Linux' : 'Windows'} using{' '}
                <a href="https://github.com/indygreg/apple-platform-rs/tree/main/apple-codesign" 
                   target="_blank" 
                   rel="noopener noreferrer"
                   className="underline hover:no-underline">rcodesign</a>.
                You'll need a .p12 certificate file exported from a Mac.
              </p>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Signing Identity
            </label>
            <input
              type="text"
              value={config.darwin?.identity || ''}
              onChange={(e) => setConfig({
                ...config,
                darwin: { ...config.darwin, identity: e.target.value }
              })}
              placeholder="Developer ID Application: Your Name (TEAMID)"
              className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
            />
            {isOnMac && (
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Find with: <code className="bg-gray-100 dark:bg-gray-800 px-1 rounded">security find-identity -v -p codesigning</code>
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Team ID
            </label>
            <input
              type="text"
              value={config.darwin?.teamID || ''}
              onChange={(e) => setConfig({
                ...config,
                darwin: { ...config.darwin, teamID: e.target.value }
              })}
              placeholder="ABCD1234EF"
              className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
            />
          </div>

          {!isOnMac && (
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                P12 Certificate Path
              </label>
              <input
                type="text"
                value={config.darwin?.p12Path || ''}
                onChange={(e) => setConfig({
                  ...config,
                  darwin: { ...config.darwin, p12Path: e.target.value }
                })}
                placeholder="/path/to/certificate.p12"
                className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
              />
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Export from Keychain Access on a Mac, or generate via Apple Developer Portal
              </p>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Notarization Profile {!isOnMac && '(Mac only)'}
            </label>
            <input
              type="text"
              value={config.darwin?.keychainProfile || ''}
              onChange={(e) => setConfig({
                ...config,
                darwin: { ...config.darwin, keychainProfile: e.target.value }
              })}
              placeholder="notarytool-profile"
              disabled={!isOnMac}
              className={`w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500 ${!isOnMac ? 'opacity-50 cursor-not-allowed' : ''}`}
            />
            {isOnMac && (
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Create with: <code className="bg-gray-100 dark:bg-gray-800 px-1 rounded">xcrun notarytool store-credentials</code>
              </p>
            )}
            {!isOnMac && (
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                For cross-platform notarization, use App Store Connect API keys instead
              </p>
            )}
          </div>

          {!isOnMac && (
            <>
              <div className="border-t border-gray-200 dark:border-gray-700 pt-4 mt-4">
                <p className="text-xs text-gray-500 dark:text-gray-400 mb-3 font-medium">
                  App Store Connect API (for notarization from {hostOS === 'linux' ? 'Linux' : 'Windows'})
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  API Key ID
                </label>
                <input
                  type="text"
                  value={config.darwin?.apiKeyID || ''}
                  onChange={(e) => setConfig({
                    ...config,
                    darwin: { ...config.darwin, apiKeyID: e.target.value }
                  })}
                  placeholder="ABC123DEF4"
                  className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Issuer ID
                </label>
                <input
                  type="text"
                  value={config.darwin?.apiIssuerID || ''}
                  onChange={(e) => setConfig({
                    ...config,
                    darwin: { ...config.darwin, apiIssuerID: e.target.value }
                  })}
                  placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
                  className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  API Key Path (.p8 file)
                </label>
                <input
                  type="text"
                  value={config.darwin?.apiKeyPath || ''}
                  onChange={(e) => setConfig({
                    ...config,
                    darwin: { ...config.darwin, apiKeyPath: e.target.value }
                  })}
                  placeholder="/path/to/AuthKey_ABC123DEF4.p8"
                  className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Create at{' '}
                  <a href="https://appstoreconnect.apple.com/access/api" 
                     target="_blank" 
                     rel="noopener noreferrer"
                     className="text-blue-500 hover:underline">App Store Connect → Users and Access → Keys</a>
                </p>
              </div>
            </>
          )}
        </div>
      );
    }

    if (selectedPlatform === 'windows') {
      return (
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Certificate Path (PFX/P12)
            </label>
            <input
              type="text"
              value={config.windows?.certificatePath || ''}
              onChange={(e) => setConfig({
                ...config,
                windows: { ...config.windows, certificatePath: e.target.value }
              })}
              placeholder="/path/to/certificate.pfx"
              className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
            />
          </div>

          <div className="text-center text-xs text-gray-500 dark:text-gray-400">— or —</div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Certificate Thumbprint (Windows Store)
            </label>
            <input
              type="text"
              value={config.windows?.thumbprint || ''}
              onChange={(e) => setConfig({
                ...config,
                windows: { ...config.windows, thumbprint: e.target.value }
              })}
              placeholder="ABC123DEF456..."
              className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Timestamp Server
            </label>
            <input
              type="text"
              value={config.windows?.timestampServer || 'http://timestamp.digicert.com'}
              onChange={(e) => setConfig({
                ...config,
                windows: { ...config.windows, timestampServer: e.target.value }
              })}
              className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
            />
          </div>
        </div>
      );
    }

    if (selectedPlatform === 'linux') {
      return (
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              GPG Key ID
            </label>
            <input
              type="text"
              value={config.linux?.gpgKeyID || ''}
              onChange={(e) => setConfig({
                ...config,
                linux: { ...config.linux, gpgKeyID: e.target.value }
              })}
              placeholder="ABCD1234EFGH5678"
              className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
            />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Find with: <code className="bg-gray-100 dark:bg-gray-800 px-1 rounded">gpg --list-secret-keys --keyid-format long</code>
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              GPG Key Path (optional)
            </label>
            <input
              type="text"
              value={config.linux?.gpgKeyPath || ''}
              onChange={(e) => setConfig({
                ...config,
                linux: { ...config.linux, gpgKeyPath: e.target.value }
              })}
              placeholder="~/.gnupg/private-key.asc"
              className="w-full px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
            />
          </div>
        </div>
      );
    }

    return null;
  };

  const renderPlatformStatus = () => {
    if (!status) return null;
    
    if (selectedPlatform === 'darwin') {
      const darwin = status.darwin;
      return (
        <div className="space-y-4">
          <StatusRow
            label="Code Signing Identity"
            configured={darwin.hasIdentity}
            value={darwin.hasIdentity ? (darwin.identity || 'Configured') : 'Not configured'}
            source={darwin.configSource}
          />
          <StatusRow
            label="Notarization"
            configured={darwin.hasNotarization}
            value={darwin.hasNotarization ? `Team ID: ${darwin.teamID || 'Configured'}` : 'Not configured'}
          />
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
          <StatusRow
            label="Code Signing Certificate"
            configured={windows.hasCertificate}
            value={windows.hasCertificate ? `Type: ${windows.certificateType}` : 'Not configured'}
            source={windows.configSource}
          />
          <StatusRow
            label="SignTool"
            configured={windows.hasSignTool}
            value={windows.hasSignTool ? 'Available' : 'Not found (Windows SDK required)'}
          />
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
          <StatusRow
            label="GPG Signing Key"
            configured={linux.hasGpgKey}
            value={linux.hasGpgKey ? `Key ID: ${linux.gpgKeyID}` : 'Not configured'}
            source={linux.configSource}
          />
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
              <PlatformTab
                platform="darwin"
                label="macOS"
                isActive={selectedPlatform === 'darwin'}
                hasConfig={status?.darwin.hasIdentity}
                onClick={() => { setSelectedPlatform('darwin'); setConfiguring(false); }}
              />
              <PlatformTab
                platform="windows"
                label="Windows"
                isActive={selectedPlatform === 'windows'}
                hasConfig={status?.windows.hasCertificate}
                onClick={() => { setSelectedPlatform('windows'); setConfiguring(false); }}
              />
              <PlatformTab
                platform="linux"
                label="Linux"
                isActive={selectedPlatform === 'linux'}
                hasConfig={status?.linux.hasGpgKey}
                onClick={() => { setSelectedPlatform('linux'); setConfiguring(false); }}
              />
            </div>

            <AnimatePresence mode="wait">
              {configuring ? (
                <motion.div
                  key="config"
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ duration: 0.2 }}
                >
                  {renderConfigForm()}
                  <div className="flex gap-3 mt-6">
                    <button
                      onClick={() => setConfiguring(false)}
                      className="flex-1 px-4 py-2 rounded-lg text-sm font-medium border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={handleSave}
                      disabled={saving}
                      className="flex-1 px-4 py-2 rounded-lg text-sm font-medium bg-red-500 text-white hover:bg-red-600 disabled:opacity-50"
                    >
                      {saving ? 'Saving...' : 'Save'}
                    </button>
                  </div>
                </motion.div>
              ) : (
                <motion.div
                  key="status"
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ duration: 0.2 }}
                  role="tabpanel"
                >
                  {renderPlatformStatus()}
                  <button
                    onClick={() => setConfiguring(true)}
                    className="w-full mt-4 px-4 py-2 rounded-lg text-sm font-medium border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
                  >
                    Configure {selectedPlatform === 'darwin' ? 'macOS' : selectedPlatform === 'windows' ? 'Windows' : 'Linux'} Signing
                  </button>
                </motion.div>
              )}
            </AnimatePresence>

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

function StatusRow({ label, configured, value, source }: { 
  label: string; 
  configured: boolean; 
  value: string; 
  source?: string;
}) {
  return (
    <div className="flex items-center gap-3 p-4 rounded-lg bg-gray-100 dark:bg-gray-900/50">
      <div className={`w-8 h-8 rounded-full flex items-center justify-center ${configured ? 'bg-green-500/20' : 'bg-gray-200 dark:bg-gray-800'}`}>
        {configured ? (
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
        <div className="text-sm font-medium text-gray-900 dark:text-white">{label}</div>
        <div className="text-xs text-gray-500 dark:text-gray-400">{value}</div>
      </div>
      {source && (
        <span className="text-xs px-2 py-1 rounded-full bg-gray-200 dark:bg-gray-800 text-gray-600 dark:text-gray-400">
          {source}
        </span>
      )}
    </div>
  );
}

function PlatformTab({ platform, label, isActive, hasConfig, onClick }: {
  platform: 'darwin' | 'windows' | 'linux';
  label: string;
  isActive: boolean;
  hasConfig?: boolean;
  onClick: () => void;
}) {
  const iconClass = `w-5 h-5 ${isActive ? 'text-gray-900 dark:text-white' : 'text-gray-400 dark:text-gray-500'}`;
  
  return (
    <button
      role="tab"
      aria-selected={isActive}
      onClick={onClick}
      className={`flex-1 flex items-center justify-center gap-2 px-4 py-3 rounded-lg text-sm font-medium transition-all ${
        isActive
          ? 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white'
          : 'text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800/50'
      }`}
    >
      {platform === 'darwin' && (
        <svg className={iconClass} viewBox="0 0 24 24" fill="currentColor">
          <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
        </svg>
      )}
      {platform === 'windows' && (
        <svg className={iconClass} viewBox="0 0 24 24" fill="currentColor">
          <path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>
        </svg>
      )}
      {platform === 'linux' && (
        <svg className={iconClass} viewBox="0 0 24 24" fill="currentColor">
          <path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a.424.424 0 00-.11.135c-.26.268-.45.6-.663.839-.199.199-.485.267-.797.4-.313.136-.658.269-.864.68-.09.189-.136.394-.132.602 0 .199.027.4.055.536.058.399.116.728.04.97-.249.68-.28 1.145-.106 1.484.174.334.535.47.94.601.81.2 1.91.135 2.774.6.926.466 1.866.67 2.616.47.526-.116.97-.464 1.208-.946.587-.003 1.23-.269 2.26-.334.699-.058 1.574.267 2.577.2.025.134.063.198.114.333l.003.003c.391.778 1.113 1.132 1.884 1.071.771-.06 1.592-.536 2.257-1.306.631-.765 1.683-1.084 2.378-1.503.348-.199.629-.469.649-.853.023-.4-.2-.811-.714-1.376v-.097l-.003-.003c-.17-.2-.25-.535-.338-.926-.085-.401-.182-.786-.492-1.046h-.003c-.059-.054-.123-.067-.188-.135a.357.357 0 00-.19-.064c.431-1.278.264-2.55-.173-3.694-.533-1.41-1.465-2.638-2.175-3.483-.796-1.005-1.576-1.957-1.56-3.368.026-2.152.236-6.133-3.544-6.139zm.529 3.405h.013c.213 0 .396.062.584.198.19.135.33.332.438.533.105.259.158.459.166.724 0-.02.006-.04.006-.06v.105a.086.086 0 01-.004-.021l-.004-.024a1.807 1.807 0 01-.15.706.953.953 0 01-.213.335.71.71 0 00-.088-.042c-.104-.045-.198-.064-.284-.133a1.312 1.312 0 00-.22-.066c.05-.06.146-.133.183-.198.053-.128.082-.264.088-.402v-.02a1.21 1.21 0 00-.061-.4c-.045-.134-.101-.2-.183-.333-.084-.066-.167-.132-.267-.132h-.016c-.093 0-.176.03-.262.132a.8.8 0 00-.205.334 1.18 1.18 0 00-.09.4v.019c.002.089.008.179.02.267-.193-.067-.438-.135-.607-.202a1.635 1.635 0 01-.018-.2v-.02a1.772 1.772 0 01.15-.768c.082-.22.232-.406.43-.533a.985.985 0 01.594-.2zm-2.962.059h.036c.142 0 .27.048.399.135.146.129.264.288.344.465.09.199.14.4.153.667v.004c.007.134.006.2-.002.266v.08c-.03.007-.056.018-.083.024-.152.055-.274.135-.393.2.012-.09.013-.18.003-.267v-.015c-.012-.133-.04-.2-.082-.333a.613.613 0 00-.166-.267.248.248 0 00-.183-.064h-.021c-.071.006-.13.04-.186.132a.552.552 0 00-.12.27.944.944 0 00-.023.33v.015c.012.135.037.2.08.267a.86.86 0 00.153.2c.071.085.178.135.305.178l.056.02a.398.398 0 00-.104.078c-.09.088-.198.2-.318.267-.145.085-.232.135-.39.135a1.04 1.04 0 01-.507-.151c-.106-.067-.199-.135-.285-.202l-.072-.053c-.239-.2-.439-.401-.618-.535a2.494 2.494 0 01-.393-.4c-.078-.1-.143-.199-.2-.298l-.06-.135-.048.066c-.078.133-.127.266-.127.465 0 .2.049.4.127.535.078.133.2.265.35.331.148.068.313.135.47.202.234.1.438.2.59.331.15.135.234.27.234.402 0 .135-.063.265-.198.332-.142.065-.32.102-.578.102-.232 0-.465-.037-.67-.1-.204-.068-.378-.17-.51-.301-.135-.135-.237-.301-.305-.5-.066-.199-.103-.432-.103-.699 0-.265.037-.5.106-.698.068-.2.166-.366.3-.5.135-.135.301-.234.5-.3.2-.067.432-.1.699-.1.266 0 .5.033.699.1.199.066.365.165.5.3.135.134.233.3.3.5.068.198.101.433.101.698 0 .267-.033.5-.1.7-.068.199-.166.365-.301.5-.135.134-.301.233-.5.3-.199.067-.433.1-.699.1-.267 0-.5-.033-.7-.1a1.379 1.379 0 01-.5-.3c-.134-.135-.233-.301-.3-.5-.066-.2-.1-.433-.1-.7 0-.266.034-.5.1-.698.067-.2.166-.366.3-.5.135-.135.301-.234.5-.3.2-.067.433-.1.7-.1z"/>
        </svg>
      )}
      <span>{label}</span>
      {hasConfig && (
        <span className="w-2 h-2 rounded-full bg-green-500" />
      )}
    </button>
  );
}
