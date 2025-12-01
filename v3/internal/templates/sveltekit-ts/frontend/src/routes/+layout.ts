import {WML} from '@wailsio/runtime';

// Enable Wails Markup Language (WML) for data-wml-* attributes
if (typeof window !== 'undefined') {
    WML.Enable();
}

export const prerender = true;
export const ssr = false;
