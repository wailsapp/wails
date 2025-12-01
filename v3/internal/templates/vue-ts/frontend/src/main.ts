import { createApp } from 'vue'
import {WML} from '@wailsio/runtime'
import App from './App.vue'

// Enable Wails Markup Language (WML) for data-wml-* attributes
WML.Enable();

createApp(App).mount('#app')
