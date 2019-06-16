import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';

import Bridge from "./wailsbridge";

Bridge.Start(() => {
  ReactDOM.render(<App />, document.getElementById('app'));
});
