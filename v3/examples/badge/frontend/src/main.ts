import {Events} from "@wailsio/runtime";
import {SetBadge, RemoveBadge} from "../bindings/github.com/wailsapp/wails/v3/pkg/services/badge/service"; 

const setButton = document.getElementById('set')! as HTMLButtonElement;
const removeButton = document.getElementById('remove')! as HTMLButtonElement;
const labelElement : HTMLInputElement = document.getElementById('label')! as HTMLInputElement;
const timeElement = document.getElementById('time')! as HTMLDivElement;

setButton.addEventListener('click', () => {
    let label = (labelElement as HTMLInputElement).value
    SetBadge(label);
});

removeButton.addEventListener('click', () => {
    RemoveBadge();
});

Events.On('time', (time: {data: any}) => {
    timeElement.innerText = time.data;
});
