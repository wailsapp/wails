import {WML} from '@wailsio/runtime'
import App from './App.svelte'

// Enable Wails Markup Language (WML) for data-wml-* attributes
WML.Enable();

const app = new App({
  target: document.getElementById('app'),
})

export default app
