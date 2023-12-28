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
 * @typedef {import("./types").MessageDialogOptions} MessageDialogOptions
 * @typedef {import("./types").OpenDialogOptions} OpenDialogOptions
 * @typedef {import("./types").SaveDialogOptions} SaveDialogOptions
 */

import {newRuntimeCallerWithID, objectNames} from "./runtime";

import { nanoid } from 'nanoid/non-secure';

// Define constants from the `methods` object in Title Case
const DialogInfo = 0;
const DialogWarning = 1;
const DialogError = 2;
const DialogQuestion = 3;
const DialogOpenFile = 4;
const DialogSaveFile = 5;

const call = newRuntimeCallerWithID(objectNames.Dialog, '');
const dialogResponses = new Map();

/**
 * Generates a unique id that is not present in dialogResponses.
 * @returns {string} unique id
 */
function generateID() {
    let result;
    do {
        result = nanoid();
    } while (dialogResponses.has(result));
    return result;
}

/**
 * Shows a dialog of specified type with the given options.
 * @param {number} type - type of dialog
 * @param {object} options - options for the dialog
 * @returns {Promise} promise that resolves with result of dialog
 */
function dialog(type, options = {}) {
    const id = generateID();
    options["dialog-id"] = id;
    return new Promise((resolve, reject) => {
        dialogResponses.set(id, {resolve, reject});
        call(type, options).catch((error) => {
            reject(error);
            dialogResponses.delete(id);
        });
    });
}

/**
 * Handles the callback from a dialog.
 *
 * @param {string} id - The ID of the dialog response.
 * @param {string} data - The data received from the dialog.
 * @param {boolean} isJSON - Flag indicating whether the data is in JSON format.
 *
 * @return {undefined}
 */
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

/**
 * Callback function for handling errors in dialog.
 *
 * @param {string} id - The id of the dialog response.
 * @param {string} message - The error message.
 *
 * @return {void}
 */
export function dialogErrorCallback(id, message) {
    let p = dialogResponses.get(id);
    if (p) {
        p.reject(message);
        dialogResponses.delete(id);
    }
}


// Replace `methods` with constants in Title Case

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Info = (options) => dialog(DialogInfo, options);

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Warning = (options) => dialog(DialogWarning, options);

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Error = (options) => dialog(DialogError, options);

/**
 * @param {MessageDialogOptions} options - Dialog options
 * @returns {Promise<string>} - The label of the button pressed
 */
export const Question = (options) => dialog(DialogQuestion, options);

/**
 * @param {OpenDialogOptions} options - Dialog options
 * @returns {Promise<string[]|string>} Returns selected file or list of files. Returns blank string if no file is selected.
 */
export const OpenFile = (options) => dialog(DialogOpenFile, options);

/**
 * @param {SaveDialogOptions} options - Dialog options
 * @returns {Promise<string>} Returns the selected file. Returns blank string if no file is selected.
 */
export const SaveFile = (options) => dialog(DialogSaveFile, options);
