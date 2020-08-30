import React from 'react';
import ReactDOM from 'react-dom';
import 'core-js/stable';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';

import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
  ReactDOM.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
    document.getElementById("app")
  );
});

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
