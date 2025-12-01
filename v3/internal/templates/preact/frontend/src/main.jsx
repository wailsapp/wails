import { render } from 'preact'
import {WML} from '@wailsio/runtime'
import { App } from './app'

// Enable Wails Markup Language (WML) for data-wml-* attributes
WML.Enable();

render(<App />, document.getElementById('app'))
