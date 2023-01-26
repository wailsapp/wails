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


import {Call} from "./calls";
import {invoke} from "./ipc";

export function newWindow(id) {
    return {
        // Reload: () => invoke('WR', id),
        // ReloadApp: () => invoke('WR', id),
        // SetSystemDefaultTheme: () => invoke('WASDT', id),
        // SetLightTheme: () => invoke('WALT', id),
        // SetDarkTheme: () => invoke('WADT', id),
        Center: () => invoke('Wc', id),
        SetTitle: (title) => invoke('WT' + title, id),
        Fullscreen: () => invoke('WF', id),
        UnFullscreen: () => invoke('Wf', id),
        SetSize: (width, height) => invoke('WS' + width + ',' + height, id),
        GetSize: () => {
            return Call(":wails:WindowGetSize")
        },
        SetMaxSize: (width, height) => invoke('WZ:' + width + ':' + height, id),
        SetMinSize: (width, height) => invoke('Wz:' + width + ':' + height, id),
        SetAlwaysOnTop: (b) => invoke('WATP:' + (b ? '1' : '0'), id),
        SetPosition: (x, y) => invoke('Wp:' + x + ':' + y, id),
        GetPosition: () => {
            return Call(":wails:WindowGetPos")
        },
        Hide: () => invoke('WH', id),
        Maximise: () => invoke('WM', id),
        Show: () => invoke('WS', id),
        ToggleMaximise: () => invoke('Wt', id),
        UnMaximise: () => invoke('WU', id),
        Minimise: () => invoke('Wm', id),
        UnMinimise: () => invoke('Wu', id),
        SetBackgroundColour: (R, G, B, A) =>
            invoke('Wr:' + JSON.stringify({
                r: R || 0,
                g: G || 0,
                b: B || 0,
                a: A || 255}, id)
            ),
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

