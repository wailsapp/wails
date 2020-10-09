import { writable } from 'svelte/store';
import runtime from '@wails/runtime';

export let selectedPage = writable();

export let darkMode = writable(runtime.System.DarkModeEnabled());

// Handle Dark/Light themes automatically
runtime.System.OnThemeChange( (isDarkMode) => {
    darkMode.set(isDarkMode);
});
