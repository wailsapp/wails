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

let call: (method: string, args?: any) => Promise<any> = newRuntimeCaller("screens");

export interface Screen
{
    ID: string           // A unique identifier for the display
    Name: string       // The name of the display
    Scale: number     // The scale factor of the display
    X: number             // The x-coordinate of the top-left corner of the rectangle
    Y: number           // The y-coordinate of the top-left corner of the rectangle
    Size: Size               // The size of the display
    Bounds: Rect              // The bounds of the display
    WorkArea: Rect               // The work area of the display
    IsPrimary: boolean    // Whether this is the primary display
    Rotation: number   // The rotation of the display
}

export interface Rect
{
    X:      number
    Y:      number
    Width:  number
    Height: number
}

export interface Size
{
    Width:  number
    Height: number
}

// Get all screens
export function GetAll(): Promise<Screen[]> {
    return call("GetPrimary");
}

// Get the primary screen
export function GetPrimary(): Promise<Screen> {
    return call("GetPrimary");
}

// Get the screen the current window is on
export function GetCurrent(): Promise<Screen> {
    return call("GetCurrent");
}