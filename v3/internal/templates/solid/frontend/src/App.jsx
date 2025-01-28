import { createSignal, onMount } from 'solid-js'
import {GreetService} from "../bindings/changeme";
import {Events} from "@wailsio/runtime";

function App() {
  const [name, setName] = createSignal('');
  const [result, setResult] = createSignal('Please enter your name below 👇');
  const [time, setTime] = createSignal('Listening for Time event...');

  const doGreet = () => {
    let localName = name();
    if (!localName) {
      localName = 'anonymous';
    }
    GreetService.Greet(localName).then((resultValue) => {
      setResult(resultValue);
    }).catch((err) => {
      console.log(err);
    });
  }

  onMount(() => {
    Events.On('time', (timeValue) => {
      setTime(timeValue.data);
    });
  });

  return (
    <div className="container">
      <div>
        <a data-wml-openURL="https://wails.io">
          <img src="/wails.png" className="logo" alt="Wails logo"/>
        </a>
        <a data-wml-openURL="https://solidjs.com">
          <img src="/solid.svg" className="logo solid" alt="Solid logo"/>
        </a>
      </div>
      <h1>Wails + Solid</h1>
      <div className="result">{result()}</div>
      <div className="card">
        <div className="input-box">
          <input className="input" value={name()} onInput={(e) => setName(e.target.value)} type="text" autocomplete="off"/>
          <button className="btn" onClick={doGreet}>Greet</button>
        </div>
      </div>
      <div className="footer">
        <div><p>Click on the Wails logo to learn more</p></div>
        <div><p>{time()}</p></div>
      </div>
    </div>
  )
}

export default App
