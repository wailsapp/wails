/* @refresh reload */
import { render } from 'solid-js/web'
import {WML} from '@wailsio/runtime'
import App from './App'

// Enable Wails Markup Language (WML) for data-wml-* attributes
WML.Enable();

const root = document.getElementById('root')

render(() => <App />, root!)
