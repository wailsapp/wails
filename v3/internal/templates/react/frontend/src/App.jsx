import { useState, useEffect } from 'react'
import {GreetService} from "../bindings/changeme";
import {Events, WML} from "@wailsio/runtime";

function App() {
  const [name, setName] = useState('');
  const [result, setResult] = useState('Please enter your name below 👇');
  const [time, setTime] = useState('Listening for Time event...');

  const doGreet = () => {
    let localName = name;
    if (!localName) {
      localName = 'anonymous';
    }
    GreetService.Greet(localName).then((resultValue) => {
      setResult(resultValue);
    }).catch((err) => {
      console.log(err);
    });
  }

  useEffect(() => {
    Events.On('time', (timeValue) => {
      setTime(timeValue.data);
    });
    // Reload WML so it picks up the wml tags
    WML.Reload();
  }, []);

  return (
    <div className="container">
      <div>
        <a wml-openURL="https://wails.io">
          <img src="/wails.png" className="logo" alt="Wails logo"/>
        </a>
        <a wml-openURL="https://reactjs.org">
          <img src='/react.svg' className="logo react" alt="React logo"/>
        </a>
      </div>
      <h1>Wails + React</h1>
      <div className="result">{result}</div>
      <div className="card">
        <div className="input-box">
          <input className="input" value={name} onChange={(e) => setName(e.target.value)} type="text" autoComplete="off"/>
          <button className="btn" onClick={doGreet}>Greet</button>
        </div>
      </div>
      <div className="footer">
        <div><p>Click on the Wails logo to learn more</p></div>
        <div><p>{time}</p></div>
      </div>
    </div>
  )
}

export default App
