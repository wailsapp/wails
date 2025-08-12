import { useState, useEffect } from 'preact/hooks'
import {GreetService} from "../bindings/changeme";
import {Events} from "@wailsio/runtime";

export function App() {
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
  }, []);

  return (
    <div className="container">
      <div>
        <a data-wml-openURL="https://wails.io">
          <img src="/wails.png" className="logo" alt="Wails logo"/>
        </a>
        <a data-wml-openURL="https://preactjs.com">
          <img src="/preact.svg" className="logo preact" alt="Preact logo"/>
        </a>
      </div>
      <h1>Wails + Preact</h1>
      <div className="result">{result}</div>
      <div className="card">
        <div className="input-box">
          <input className="input" value={name} onInput={(e) => setName(e.target.value)} type="text" autoComplete="off"/>
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
