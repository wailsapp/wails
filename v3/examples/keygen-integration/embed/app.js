// Event listener setup
document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    getVersion();
});

// Global variables
let currentUpdateInfo = null;
let downloadInProgress = false;

// Setup event listeners for Keygen events
function setupEventListeners() {
    // Listen for license status events
    window.wails.Event.On('keygen:license-status', (event) => {
        logEvent('License Status Changed', event);
        updateLicenseStatus(event.data);
    });

    // Listen for update available events
    window.wails.Event.On('keygen:update-available', (event) => {
        logEvent('Update Available', event);
        showUpdateAvailable(event.data);
    });

    // Listen for download progress events
    window.wails.Event.On('keygen:download-progress', (event) => {
        updateDownloadProgress(event.data);
    });

    // Listen for update installed events
    window.wails.Event.On('keygen:update-installed', (event) => {
        logEvent('Update Installed', event);
        alert('Update installed successfully! Please restart the application.');
    });
}

// Get current version
async function getVersion() {
    try {
        const version = await window.App.GetCurrentVersion();
        document.getElementById('version').textContent = version;
    } catch (error) {
        logError('Failed to get version', error);
    }
}

// License Management Functions
async function validateLicense() {
    const licenseKey = document.getElementById('licenseKey').value.trim();
    if (!licenseKey) {
        alert('Please enter a license key');
        return;
    }

    try {
        // Set the license key
        await window.App.SetLicenseKey(licenseKey);
        
        // Validate the license
        const result = await window.App.ValidateLicense();
        
        if (result.valid) {
            document.getElementById('licenseStatus').classList.remove('hidden');
            updateLicenseStatus({
                valid: true,
                status: 'active',
                key: licenseKey,
                message: result.message
            });
            
            // Enable machine activation if required
            if (result.requiresActivation) {
                document.getElementById('activateBtn').disabled = false;
            }
        } else {
            alert(`License validation failed: ${result.message}`);
        }
    } catch (error) {
        logError('License validation error', error);
        alert(`Error: ${error}`);
    }
}

async function getLicenseInfo() {
    try {
        const info = await window.App.GetLicenseInfo();
        document.getElementById('licenseDetails').classList.remove('hidden');
        document.getElementById('licenseDetailsContent').textContent = JSON.stringify(info, null, 2);
        logEvent('License Info Retrieved', info);
    } catch (error) {
        logError('Failed to get license info', error);
        alert(`Error: ${error}`);
    }
}

async function saveOfflineLicense() {
    try {
        await window.App.SaveOfflineLicense();
        alert('License saved for offline use');
        logEvent('Offline License Saved', {});
    } catch (error) {
        logError('Failed to save offline license', error);
        alert(`Error: ${error}`);
    }
}

async function clearLicenseCache() {
    if (confirm('Are you sure you want to clear the license cache?')) {
        try {
            await window.App.ClearLicenseCache();
            alert('License cache cleared');
            logEvent('License Cache Cleared', {});
            
            // Reset UI
            document.getElementById('licenseStatus').classList.add('hidden');
            document.getElementById('licenseDetails').classList.add('hidden');
            document.getElementById('licenseKey').value = '';
        } catch (error) {
            logError('Failed to clear license cache', error);
            alert(`Error: ${error}`);
        }
    }
}

// Machine Activation Functions
async function activateMachine() {
    try {
        const result = await window.App.ActivateMachine();
        if (result.success) {
            document.getElementById('activateBtn').disabled = true;
            document.getElementById('deactivateBtn').disabled = false;
            document.getElementById('machineInfo').classList.remove('hidden');
            
            document.getElementById('machineIdValue').textContent = result.machine?.id || 'N/A';
            document.getElementById('fingerprintValue').textContent = result.fingerprint;
            
            alert('Machine activated successfully!');
            logEvent('Machine Activated', result);
        } else {
            alert(`Machine activation failed: ${result.message}`);
        }
    } catch (error) {
        logError('Machine activation error', error);
        alert(`Error: ${error}`);
    }
}

async function deactivateMachine() {
    if (confirm('Are you sure you want to deactivate this machine?')) {
        try {
            await window.App.DeactivateMachine();
            
            document.getElementById('activateBtn').disabled = false;
            document.getElementById('deactivateBtn').disabled = true;
            document.getElementById('machineInfo').classList.add('hidden');
            
            alert('Machine deactivated successfully!');
            logEvent('Machine Deactivated', {});
        } catch (error) {
            logError('Machine deactivation error', error);
            alert(`Error: ${error}`);
        }
    }
}

// Feature Entitlement Functions
async function checkEntitlement() {
    const featureName = document.getElementById('featureName').value.trim();
    if (!featureName) {
        alert('Please enter a feature name');
        return;
    }

    try {
        const enabled = await window.App.CheckEntitlement(featureName);
        const resultDiv = document.getElementById('entitlementResult');
        const resultText = document.getElementById('entitlementResultText');
        
        resultDiv.classList.remove('hidden');
        resultText.textContent = `Feature "${featureName}" is ${enabled ? 'ENABLED' : 'DISABLED'}`;
        resultText.className = enabled ? 'success' : 'error';
        
        logEvent('Entitlement Checked', { feature: featureName, enabled });
    } catch (error) {
        logError('Entitlement check error', error);
        alert(`Error: ${error}`);
    }
}

// Update Management Functions
async function setUpdateChannel() {
    const channel = document.getElementById('updateChannel').value;
    try {
        await window.App.SetUpdateChannel(channel);
        logEvent('Update Channel Changed', { channel });
    } catch (error) {
        logError('Failed to set update channel', error);
        alert(`Error: ${error}`);
    }
}

async function checkForUpdates() {
    try {
        const updateInfo = await window.App.CheckForUpdates();
        currentUpdateInfo = updateInfo;
        
        if (updateInfo.available) {
            showUpdateAvailable(updateInfo);
        } else {
            alert('No updates available. You are running the latest version!');
        }
        
        logEvent('Update Check', updateInfo);
    } catch (error) {
        logError('Update check error', error);
        alert(`Error: ${error}`);
    }
}

async function downloadUpdate() {
    if (!currentUpdateInfo || !currentUpdateInfo.releaseId) {
        alert('No update information available');
        return;
    }

    if (downloadInProgress) {
        alert('Download already in progress');
        return;
    }

    try {
        downloadInProgress = true;
        document.getElementById('downloadBtn').disabled = true;
        document.getElementById('downloadProgress').classList.remove('hidden');
        
        const progress = await window.App.DownloadUpdate(currentUpdateInfo.releaseId);
        logEvent('Download Started', { releaseId: currentUpdateInfo.releaseId });
    } catch (error) {
        downloadInProgress = false;
        document.getElementById('downloadBtn').disabled = false;
        logError('Download error', error);
        alert(`Error: ${error}`);
    }
}

async function installUpdate() {
    if (confirm('The application will restart to complete the update. Continue?')) {
        try {
            await window.App.InstallUpdate();
            logEvent('Update Installation Started', {});
        } catch (error) {
            logError('Installation error', error);
            alert(`Error: ${error}`);
        }
    }
}

// UI Update Functions
function updateLicenseStatus(data) {
    document.getElementById('statusValue').textContent = data.valid ? 'Valid' : 'Invalid';
    document.getElementById('statusValue').className = `value ${data.valid ? 'success' : 'error'}`;
    document.getElementById('emailValue').textContent = data.email || '-';
    document.getElementById('expiresValue').textContent = data.expiresAt ? new Date(data.expiresAt).toLocaleDateString() : 'Never';
    document.getElementById('lastCheckedValue').textContent = data.lastChecked ? new Date(data.lastChecked).toLocaleString() : '-';
}

function showUpdateAvailable(updateInfo) {
    document.getElementById('updateStatus').classList.remove('hidden');
    document.getElementById('newVersionValue').textContent = updateInfo.latestVersion || updateInfo.version;
    document.getElementById('releaseDateValue').textContent = updateInfo.publishedAt ? 
        new Date(updateInfo.publishedAt).toLocaleDateString() : 
        (updateInfo.releaseDate ? new Date(updateInfo.releaseDate).toLocaleDateString() : '-');
    document.getElementById('updateSizeValue').textContent = formatBytes(updateInfo.size);
    document.getElementById('releaseNotes').textContent = updateInfo.releaseNotes || updateInfo.notes || 'No release notes available';
    
    currentUpdateInfo = updateInfo;
}

function updateDownloadProgress(data) {
    const progressBar = document.getElementById('progressBar');
    const progressPercent = document.getElementById('progressPercent');
    const progressSpeed = document.getElementById('progressSpeed');
    const progressETA = document.getElementById('progressETA');
    
    progressBar.style.width = `${data.progress}%`;
    progressPercent.textContent = `${Math.round(data.progress)}%`;
    progressSpeed.textContent = `${formatBytes(data.speed)}/s`;
    progressETA.textContent = data.timeRemaining > 0 ? `ETA: ${formatTime(data.timeRemaining)}` : 'ETA: --:--';
    
    // Enable install button when download is complete
    if (data.status === 'completed') {
        document.getElementById('installBtn').disabled = false;
        downloadInProgress = false;
        logEvent('Download Completed', { progress: 100 });
    } else if (data.status === 'failed') {
        downloadInProgress = false;
        document.getElementById('downloadBtn').disabled = false;
        alert('Download failed: ' + (data.error || 'Unknown error'));
        logError('Download Failed', data.error);
    }
}

// Event Logging Functions
function logEvent(event, data) {
    const log = document.getElementById('eventLog');
    const entry = document.createElement('div');
    entry.className = 'event-entry';
    entry.innerHTML = `
        <span class="event-time">${new Date().toLocaleTimeString()}</span>
        <span class="event-name">${event}</span>
        <span class="event-data">${JSON.stringify(data)}</span>
    `;
    log.insertBefore(entry, log.firstChild);
    
    // Keep only last 50 events
    while (log.children.length > 50) {
        log.removeChild(log.lastChild);
    }
}

function logError(context, error) {
    logEvent(`ERROR: ${context}`, { error: error.toString() });
}

function clearEventLog() {
    document.getElementById('eventLog').innerHTML = '';
}

// Utility Functions
function formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function formatTime(seconds) {
    if (!seconds || seconds < 0) return '--:--';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
}