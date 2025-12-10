import {GreetService} from "./bindings/changeme";
import {Events} from "@wailsio/runtime";

const resultElement = document.getElementById('result');
const timeElement = document.getElementById('time');

window.doGreet = () => {
    let name = document.getElementById('name').value;
    if (!name) {
        name = 'anonymous';
    }
    GreetService.Greet(name).then((result) => {
        resultElement.innerText = result;
    }).catch((err) => {
        console.log(err);
    });
}

Events.On('time', (payload) => {
    // payload may be a plain value or an object with a `data` field depending on emitter/runtime
    const value = (payload && typeof payload === 'object' && 'data' in payload) ? payload.data : payload;
    console.log('[frontend] time event:', payload, '->', value);
    timeElement.innerText = value;
});
