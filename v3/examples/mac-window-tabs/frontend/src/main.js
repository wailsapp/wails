import {Events} from "@wailsio/runtime";
import {WindowService} from "../bindings/mac-window-tabs";
const timeElement = document.getElementById('time');
const summaryElement = document.getElementById('summary');

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

window.addTab = async () => {
    try {
        await WindowService.AddTab();
    } catch (err) {
        summaryElement.innerText = String(err);
        console.error(err);
    }
}

window.detachTab = async () => {
    try {
        await WindowService.DetachTab();
    } catch (err) {
        console.error(err);
    }
}

window.tabSummary = async () => {
    try {
        summaryElement.innerText = await WindowService.TabSummary();
    } catch (err) {
        summaryElement.innerText = String(err);
        console.error(err);
    }
}

Events.On('time', (time) => {
    timeElement.innerText = time.data;
});
