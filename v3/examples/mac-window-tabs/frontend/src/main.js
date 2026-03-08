import {Events} from "@wailsio/runtime";
import {WindowService} from "../bindings/changeme";
const timeElement = document.getElementById('time');

window.openTabbedWindow = async () => {
    try {
        await WindowService.OpenTabbedWindow();
    } catch (err) {
        console.error(err);
    }
}

window.openNonTabbedWindow = async () => {
    try {
        await WindowService.OpenNonTabbedWindow();
    } catch (err) {
        console.error(err);
    }
}

Events.On('time', (time) => {
    timeElement.innerText = time.data;
});
