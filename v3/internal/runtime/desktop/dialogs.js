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


export function Info(options) {
    return dialog("Info", options);
}

export function Warning(options) {
    return dialog("Warning", options);
}

export function Error(options) {
    return dialog("Error", options);
}

export function Question(options) {
    return dialog("Question", options);
}

export function OpenFile(options) {
    return dialog("OpenFile", options);
}

export function SaveFile(options) {
    return dialog("SaveFile", options);
}

