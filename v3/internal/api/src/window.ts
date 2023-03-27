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

import {newRuntimeCaller} from "./runtime";

export function newWindow(id: number) {
    let call = newRuntimeCaller("window", id);
    return {
        // Reload: () => call('WR'),
        // ReloadApp: () => call('WR'),
        // SetSystemDefaultTheme: () => call('WASDT'),
        // SetLightTheme: () => call('WALT'),
        // SetDarkTheme: () => call('WADT'),
        // IsFullscreen: () => call('WIF'),
        // IsMaximized: () => call('WIM'),
        // IsMinimized: () => call('WIMN'),
        // IsWindowed: () => call('WIF'),
        Center: () => call('Center'),
        SetTitle: (title: string) => call('SetTitle', {title}),
        Fullscreen: () => call('Fullscreen'),
        UnFullscreen: () => call('UnFullscreen'),
        SetSize: (width: number, height:number) => call('SetSize', {width,height}),
        Size: () => { return call('Size') },
        SetMaxSize: (width:number, height:number) => call('SetMaxSize', {width,height}),
        SetMinSize: (width:number, height:number) => call('SetMinSize', {width,height}),
        SetAlwaysOnTop: (b:boolean) => call('SetAlwaysOnTop', {alwaysOnTop:b}),
        SetPosition: (x:number, y:number) => call('SetPosition', {x,y}),
        Position: () => { return call('Position') },
        Screen: () => { return call('Screen') },
        Hide: () => call('Hide'),
        Maximise: () => call('Maximise'),
        Show: () => call('Show'),
        Close: () => call('Close'),
        ToggleMaximise: () => call('ToggleMaximise'),
        UnMaximise: () => call('UnMaximise'),
        Minimise: () => call('Minimise'),
        UnMinimise: () => call('UnMinimise'),
        SetBackgroundColour: (r:number, g:number, b:number, a:number) => call('SetBackgroundColour', {r, g, b, a}),
    }
}
