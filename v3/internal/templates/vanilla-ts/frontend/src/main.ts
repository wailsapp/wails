import {GreetService} from "../bindings/changeme";
import {Events} from "@wailsio/runtime";

const greetButton = document.getElementById('greet')! as HTMLButtonElement;
const resultElement = document.getElementById('result')! as HTMLDivElement;
const nameElement : HTMLInputElement = document.getElementById('name')! as HTMLInputElement;
const timeElement = document.getElementById('time')! as HTMLDivElement;

greetButton.addEventListener('click', () => {
    let name = (nameElement as HTMLInputElement).value
    if (!name) {
        name = 'anonymous';
    }
    GreetService.Greet(name).then((result: string) => {
        resultElement.innerText = result;
    }).catch((err: Error) => {
        console.log(err);
    });
});

Events.On('time', (time: {data: any}) => {
    timeElement.innerText = time.data;
});
