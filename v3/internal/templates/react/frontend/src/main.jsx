import React from 'react'
import ReactDOM from 'react-dom/client'
import {WML} from '@wailsio/runtime'
import App from './App'

// Enable Wails Markup Language (WML) for data-wml-* attributes
WML.Enable();

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
