<script>
import {GreetService} from "../../bindings/changeme";
import {Events} from "@wailsio/runtime";

let name = '';
let result = 'Please enter your name below 👇';
let time = 'Listening for Time event...';

const doGreet = () => {
  let localName = name;
  if (!localName) {
    localName = 'anonymous';
  }
  GreetService.Greet(localName).then((resultValue) => {
    result = resultValue;
  }).catch((err) => {
    console.log(err);
  });
}

Events.On('time', (timeValue) => {
  time = timeValue.data;
});
</script>

<div class="container">
  <div>
    <span wml-openURL="https://wails.io">
      <img src="/wails.png" class="logo" alt="Wails logo"/>
    </span>
    <span wml-openURL="https://svelte.dev">
      <img src="/svelte.svg" class="logo svelte" alt="Svelte logo"/>
    </span>
  </div>
  <h1>Wails + Svelte</h1>
  <div class="result">{result}</div>
  <div class="card">
    <div class="input-box">
      <input class="input" bind:value={name} type="text" autocomplete="off"/>
      <button class="btn" on:click={doGreet}>Greet</button>
    </div>
  </div>
  <div class="footer">
    <div><p>Click on the Wails logo to learn more</p></div>
    <div><p>{time}</p></div>
  </div>
</div>

<style>
</style>
