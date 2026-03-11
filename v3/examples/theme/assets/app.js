import { Call, CancellablePromise, Create} from "/wails/runtime.js";
import { Events, Window} from "/wails/runtime.js";

const resultsApp = document.getElementById("app-theme");
const resultsWin = document.getElementById("win-theme");

// Call Function [Services Functions] by name
async function callBinding(name, ...params) {
    return Call.ByName(name, ...params);
}

async function setAppTheme(theme) {
    await callBinding("main.WindowService.SetAppTheme", theme);
}

async function setWinTheme(theme) {
    await callBinding("main.WindowService.SetWinTheme", theme);
}

// Window Event Listeners
window.addEventListener("DOMContentLoaded", async () => {
    // fetch the current theme from Go when the page loads
    const appTheme = await callBinding("main.WindowService.GetAppTheme");
    resultsApp.innerText = appTheme;

    const winTheme = await callBinding("main.WindowService.GetWinTheme");
    resultsWin.innerText = winTheme;
});

// Button Event Listeners
document.getElementById("app-theme-system").addEventListener("click", () => setAppTheme("system"));
document.getElementById("app-theme-light").addEventListener("click", () => setAppTheme("light"));
document.getElementById("app-theme-dark").addEventListener("click", () => setAppTheme("dark"));

document.getElementById("win-theme-app").addEventListener("click", () => setWinTheme("application"));
document.getElementById("win-theme-system").addEventListener("click", () => setWinTheme("system"));
document.getElementById("win-theme-light").addEventListener("click", () => setWinTheme("light"));
document.getElementById("win-theme-dark").addEventListener("click", () => setWinTheme("dark"));

// Go Event Listeners
Events.On("common:ApplicationThemeChanged", async (ev) => {
    const appTheme = await callBinding("main.WindowService.GetAppTheme");
    resultsApp.innerText = appTheme;
});

Events.On("common:ThemeChanged", async (ev) => {
    const winTheme = await callBinding("main.WindowService.GetWinTheme");
    resultsWin.innerText = winTheme;
});