import React from 'react';
import ReactDOM from 'react-dom';
import "core-js/stable";
import './index.css';
import App from './App';

import Bridge from "./wailsbridge";

Bridge.Start(() => {
  ReactDOM.render(<App />, document.getElementById('app'));
});
