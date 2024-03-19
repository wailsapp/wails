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

/**
 * @typedef {Object} Position
 * @property {number} X - The X coordinate.
 * @property {number} Y - The Y coordinate.
 */

/**
 * @typedef {Object} Size
 * @property {number} X - The width.
 * @property {number} Y - The height.
 */


/**
 * @typedef {Object} Rect
 * @property {number} X - The X coordinate of the top-left corner.
 * @property {number} Y - The Y coordinate of the top-left corner.
 * @property {number} Width - The width of the rectangle.
 * @property {number} Height - The height of the rectangle.
 */


/**
 * @typedef {Object} Screen
 * @property {string} Id - Unique identifier for the screen.
 * @property {string} Name - Human readable name of the screen.
 * @property {number} Scale - The resolution scale of the screen. 1 = standard resolution, 2 = high (Retina), etc.
 * @property {Position} Position - Contains the X and Y coordinates of the screen's position.
 * @property {Size} Size - Contains the width and height of the screen.
 * @property {Rect} Bounds - Contains the bounds of the screen in terms of X, Y, Width, and Height.
 * @property {Rect} WorkArea - Contains the area of the screen that is actually usable (excluding taskbar and other system UI).
 * @property {boolean} IsPrimary - True if this is the primary monitor selected by the user in the operating system.
 * @property {number} Rotation - The rotation of the screen.
 */


import {newRuntimeCallerWithID, objectNames} from "./runtime";
const call = newRuntimeCallerWithID(objectNames.Screens, '');

const getAll = 0;
const getPrimary = 1;
const getCurrent = 2;

/**
 * Gets all screens.
 * @returns {Promise<Screen[]>} A promise that resolves to an array of Screen objects.
 */
export function GetAll() {
    return call(getAll);
}
/**
 * Gets the primary screen.
 * @returns {Promise<Screen>} A promise that resolves to the primary screen.
 */
export function GetPrimary() {
    return call(getPrimary);
}
/**
 * Gets the current active screen.
 *
 * @returns {Promise<Screen>} A promise that resolves with the current active screen.
 */
export function GetCurrent() {
    return call(getCurrent);
}
