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
 * @typedef {import("./api/types").MessageDialogOptions} MessageDialogOptions
 * @typedef {import("./api/types").OpenDialogOptions} OpenDialogOptions
 * @typedef {import("./api/types").SaveDialogOptions} SaveDialogOptions
 */

import {newRuntimeCaller} from "./runtime";

import { nanoid } from 'nanoid/non-secure';

let call = newRuntimeCaller("dialog");

let dialogResponses = new Map();

function generateID() {
    let result;
    do {
        result = nanoid();
    } while (dialogResponses.has(result));
    return result;
}

export function dialogCallback(id, data, isJSON) {
    let p = dialogResponses.get(id);
    if (p) {
        if (isJSON) {
            p.resolve(JSON.parse(data));
        } else {
            p.resolve(data);
        }
        dialogResponses.delete(id);
    }
}
export function dialogErrorCallback(id, message) {
    let p = dialogResponses.get(id);
    if (p) {
        p.reject(message);
        dialogResponses.delete(id);
    }
}

function dialog(type, options) {
    return new Promise((resolve, reject) => {
        let id = generateID();
        options = options || {};
        options["dialog-id"] = id;
        dialogResponses.set(id, {resolve, reject});
        call(type, options).catch((error) => {
            reject(error);
            dialogResponses.delete(id);
        });
    });
}


/**
 * Shows an Info dialog with the given options.
 * @param {MessageDialogOptions} options
 * @returns {Promise<string>} The label of the button pressed
 */
export function Info(options) {
    return dialog("Info", options);
}

/**
 * Shows an Warning dialog with the given options.
 * @param {MessageDialogOptions} options
 * @returns {Promise<string>} The label of the button pressed
 */
export function Warning(options) {
    return dialog("Warning", options);
}

/**
 * Shows an Error dialog with the given options.
 * @param {MessageDialogOptions} options
 * @returns {Promise<string>} The label of the button pressed
 */
export function Error(options) {
    return dialog("Error", options);
}

/**
 * Shows a Question dialog with the given options.
 * @param {MessageDialogOptions} options
 * @returns {Promise<string>} The label of the button pressed
 */
export function Question(options) {
    return dialog("Question", options);
}

/**
 * Shows an Open dialog with the given options.
 * @param {OpenDialogOptions} options
 * @returns {Promise<string[]|string>} Returns the selected file or an array of selected files if AllowsMultipleSelection is true. A blank string is returned if no file was selected.
 */
export function OpenFile(options) {
    return dialog("OpenFile", options);
}

/**
 * Shows a Save dialog with the given options.
 * @param {OpenDialogOptions} options
 * @returns {Promise<string>} Returns the selected file. A blank string is returned if no file was selected.
 */
export function SaveFile(options) {
    return dialog("SaveFile", options);
}

