// Import the Wails runtime from npm
import '@wailsio/runtime';

// Setup media status listeners
const video = document.getElementById('video');
const videoStatus = document.getElementById('video-status');

video.addEventListener('loadeddata', () => {
    videoStatus.className = 'status success';
    videoStatus.textContent = 'Loaded (' + video.duration.toFixed(1) + 's)';
});

video.addEventListener('error', () => {
    videoStatus.className = 'status error';
    videoStatus.textContent = 'Failed to load';
});

const audio = document.getElementById('audio');
const audioStatus = document.getElementById('audio-status');

audio.addEventListener('loadeddata', () => {
    audioStatus.className = 'status success';
    audioStatus.textContent = 'Loaded (' + audio.duration.toFixed(1) + 's)';
});

audio.addEventListener('error', () => {
    audioStatus.className = 'status error';
    audioStatus.textContent = 'Failed to load';
});

setTimeout(() => {
    if (videoStatus.classList.contains('pending')) {
        videoStatus.className = 'status error';
        videoStatus.textContent = 'Timeout';
    }
    if (audioStatus.classList.contains('pending')) {
        audioStatus.className = 'status error';
        audioStatus.textContent = 'Timeout';
    }
}, 5000);
