import '@builder.io/qwik/qwikloader.js'

import { render } from '@builder.io/qwik'
import {WML} from '@wailsio/runtime'
import { App } from './app.jsx.tmpl'

// Enable Wails Markup Language (WML) for data-wml-* attributes
WML.Enable();

render(document.getElementById('app'), <App />)
