/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

export function LogPrint(message) {
    window.runtime.LogPrint(message);
}

export function LogTrace(message) {
    window.runtime.LogTrace(message);
}

export function LogDebug(message) {
    window.runtime.LogDebug(message);
}

export function LogInfo(message) {
    window.runtime.LogInfo(message);
}

export function LogWarning(message) {
    window.runtime.LogWarning(message);
}

export function LogError(message) {
    window.runtime.LogError(message);
}

export function LogFatal(message) {
    window.runtime.LogFatal(message);
}

export function EventsOnMultiple(eventName, callback, maxCallbacks) {
    window.runtime.EventsOnMultiple(eventName, callback, maxCallbacks);
}

export function EventsOn(eventName, callback) {
    EventsOnMultiple(eventName, callback, -1);
}

export function EventsOff(eventName) {
    return window.runtime.EventsOff(eventName);
}

export function EventsOnce(eventName, callback) {
    EventsOnMultiple(eventName, callback, 1);
}

export function EventsEmit(eventName) {
    let args = [eventName].slice.call(arguments);
    return window.runtime.EventsEmit.apply(null, args);
}

export function WindowReload() {
    window.runtime.WindowReload();
}

export function WindowReloadApp() {
    window.runtime.WindowReloadApp();
}

export function WindowSetSystemDefaultTheme() {
    window.runtime.WindowSetSystemDefaultTheme();
}

export function WindowSetLightTheme() {
    window.runtime.WindowSetLightTheme();
}

export function WindowSetDarkTheme() {
    window.runtime.WindowSetDarkTheme();
}

export function WindowCenter() {
    window.runtime.WindowCenter();
}

export function WindowSetTitle(title) {
    window.runtime.WindowSetTitle(title);
}

export function WindowFullscreen() {
    window.runtime.WindowFullscreen();
}

export function WindowUnfullscreen() {
    window.runtime.WindowUnfullscreen();
}

export function WindowGetSize() {
    return window.runtime.WindowGetSize();
}

export function WindowSetSize(width, height) {
    window.runtime.WindowSetSize(width, height);
}

export function WindowSetMaxSize(width, height) {
    window.runtime.WindowSetMaxSize(width, height);
}

export function WindowSetMinSize(width, height) {
    window.runtime.WindowSetMinSize(width, height);
}

export function WindowSetPosition(x, y) {
    window.runtime.WindowSetPosition(x, y);
}

export function WindowGetPosition() {
    return window.runtime.WindowGetPosition();
}

export function WindowHide() {
    window.runtime.WindowHide();
}

export function WindowShow() {
    window.runtime.WindowShow();
}

export function WindowMaximise() {
    window.runtime.WindowMaximise();
}

export function WindowToggleMaximise() {
    window.runtime.WindowToggleMaximise();
}

export function WindowUnmaximise() {
    window.runtime.WindowUnmaximise();
}

export function WindowMinimise() {
    window.runtime.WindowMinimise();
}

export function WindowUnminimise() {
    window.runtime.WindowUnminimise();
}

export function WindowSetRGBA(RGBA) {
    window.runtime.WindowSetRGBA(RGBA);
}

export function BrowserOpenURL(url) {
    window.runtime.BrowserOpenURL(url);
}

export function Environment() {
    return window.runtime.Environment();
}

export function Quit() {
    window.runtime.Quit();
}
