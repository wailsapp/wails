import {Greet} from "./bindings/main/GreetService.js";
import {Events} from "@wailsio/runtime";

window.doGreet = () => {
    let name = document.getElementById('name').value;
    if (!name) {
        name = 'from Go';
    }
    Greet(name).then((result) => {
        let element = document.getElementById('greeting');
        element.innerText = result;
    }).catch((err) => {
        console.log(err);
    });
}

Events.On('time', (time) => {
    let element = document.getElementById('time');
    element.innerText = time.data;
});