import React from 'react';
import ReactDOM from 'react-dom';
import 'core-js/stable';
import './index.css';
import App from './App';

import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
  ReactDOM.render(<App />, document.getElementById('app'));
});
