import {Events} from "@wailsio/runtime";
import {SetBadge, RemoveBadge, SetCustomBadge} from "../bindings/github.com/wailsapp/wails/v3/pkg/services/badge/badgeservice";
import { RGBA } from "../bindings/image/color/models";

const setCustomButton = document.getElementById('set-custom')! as HTMLButtonElement;
const setButton = document.getElementById('set')! as HTMLButtonElement;
const removeButton = document.getElementById('remove')! as HTMLButtonElement;
const setButtonUsingGo = document.getElementById('set-go')! as HTMLButtonElement;
const removeButtonUsingGo = document.getElementById('remove-go')! as HTMLButtonElement;
const labelElement : HTMLInputElement = document.getElementById('label')! as HTMLInputElement;
const timeElement = document.getElementById('time')! as HTMLDivElement;

setCustomButton.addEventListener('click', () => {
    console.log("click!")
    let label = (labelElement as HTMLInputElement).value
    SetCustomBadge(label, {
        BackgroundColour: RGBA.createFrom({
            R: 0,
            G: 255,
            B: 255,
            A: 255,
        }),
        FontName: "arialb.ttf", // System font
        FontSize: 16,
        SmallFontSize: 10,
        TextColour: RGBA.createFrom({
            R: 0,
            G: 0,
            B: 0,
            A: 255,
        }),
    });
})

setButton.addEventListener('click', () => {
    let label = (labelElement as HTMLInputElement).value
    SetBadge(label);
});

removeButton.addEventListener('click', () => {
    RemoveBadge();
});

setButtonUsingGo.addEventListener('click', () => {
    let label = (labelElement as HTMLInputElement).value
    void Events.Emit({
        name: "set:badge",
        data: label,
    })
})

removeButtonUsingGo.addEventListener('click', () => {
    void Events.Emit({name:"remove:badge", data: null})
})

Events.On('time', (time: {data: any}) => {
    timeElement.innerText = time.data;
});

