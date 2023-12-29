import {Greet} from "./bindings/main/GreetService.js";

window.doGreet = () => {
    Greet('test').then((result) => {
        console.log(result);
    }).catch((err) => {
        console.log(err);
    });
}