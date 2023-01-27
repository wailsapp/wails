/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

const runtimeURL = window.location.origin + "/wails/runtime";

function runtimeCall(method, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("method", method);
    if (args) {
        for (let key in args) {
            url.searchParams.append(key, args[key]);
        }
    }
    return new Promise((resolve, reject) => {
        fetch(url)
            .then(response => {
                if (response.ok) {
                    return response.json();
                }
                reject(Error(response.statusText));
            })
            .then(data => resolve(data))
            .catch(error => reject(error));
    });
}

export function newWindow(id) {
    let call = function(method, args) {
        if (id !== -1) {
            args["windowID"] = id;
        }
        return runtimeCall("window." + method, args);
    }
    return {
        // Reload: () => call('WR'),
        // ReloadApp: () => call('WR'),
        // SetSystemDefaultTheme: () => call('WASDT'),
        // SetLightTheme: () => call('WALT'),
        // SetDarkTheme: () => call('WADT'),
        Center: () => call('Center'),
        SetTitle: (title) => call('SetTitle', {title}),
        Fullscreen: () => call('Fullscreen'),
        UnFullscreen: () => call('UnFullscreen'),
        SetSize: (width, height) => call('SetSize', {width,height}),
        GetSize: () => { return call('GetSize') },
        SetMaxSize: (width, height) => call('SetMaxSize', {width,height}),
        SetMinSize: (width, height) => call('SetMinSize', {width,height}),
        SetAlwaysOnTop: (b) => call('SetAlwaysOnTop', {alwaysOnTop:b}),
        SetPosition: (x, y) => call('SetPosition', {x,y}),
        GetPosition: () => { return call('GetPosition') },
        Hide: () => call('Hide'),
        Maximise: () => call('Maximise'),
        Show: () => call('Show'),
        ToggleMaximise: () => call('ToggleMaximise'),
        UnMaximise: () => call('UnMaximise'),
        Minimise: () => call('Minimise'),
        UnMinimise: () => call('UnMinimise'),
        SetBackgroundColour: (r, g, b, a) => call('SetBackgroundColour', {R, G, B, A}),
    }
}

// export function IsFullscreen: ()=> //     return Call(":wails:WindowIsFullscreen"),
//

// export function IsMaximised: ()=> //     return Call(":wails:WindowIsMaximised"),
//

// export function IsMinimised: ()=> //     return Call(":wails:WindowIsMinimised"),
//

// export function IsNormal: ()=> //     return Call(":wails:WindowIsNormal"),
//

