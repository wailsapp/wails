import { component$, useSignal, useVisibleTask$ } from '@builder.io/qwik'
import {GreetService} from "../bindings/changeme";
import {Events, WML} from "@wailsio/runtime";

export const App = component$(() => {
  const name = useSignal('');
  const result = useSignal('Please enter your name below 👇');
  const time = useSignal('Listening for Time event...');

  const doGreet = () => {
    let localName = name.value;
    if (!localName) {
      localName = 'anonymous';
    }
    GreetService.Greet(localName).then((resultValue) => {
      result.value = resultValue;
    }).catch((err) => {
      console.log(err);
    });
  }

  useVisibleTask$(() => {
    Events.On('time', (timeValue) => {
      time.value = timeValue.data;
    });
    // Reload WML so it picks up the wml tags
    WML.Reload();
  });

  return (
    <div class="container">
      <div>
        <a wml-openURL="https://wails.io">
          <img src="/wails.png" class="logo" alt="Wails logo"/>
        </a>
        <a wml-openURL="https://qwik.builder.io">
          <img src="/qwik.svg" class="logo qwik" alt="Qwik logo"/>
        </a>
      </div>
      <h1>Wails + Qwik</h1>
      <div class="result">{result.value}</div>
      <div class="card">
        <div class="input-box">
          <input class="input" value={name.value} onInput$={(e) => name.value = e.target.value} type="text" autocomplete="off"/>
          <button class="btn" onClick$={doGreet}>Greet</button>
        </div>
      </div>
      <div class="footer">
        <div><p>Click on the Wails logo to learn more</p></div>
        <div><p>{time.value}</p></div>
      </div>
    </div>
  )
})
