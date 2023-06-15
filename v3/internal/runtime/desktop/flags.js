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

let flags = new Map();

function convertToMap(obj) {
    const map = new Map();

    for (const [key, value] of Object.entries(obj)) {
        if (typeof value === 'object' && value !== null) {
            map.set(key, convertToMap(value)); // Recursively convert nested object
        } else {
            map.set(key, value);
        }
    }

    return map;
}

fetch("/wails/flags").then((response) => {
    response.json().then((data) => {
        flags = convertToMap(data);
    });
});


function getValueFromMap(keyString) {
    const keys = keyString.split('.');
    let value = flags;

    for (const key of keys) {
        if (value instanceof Map) {
            value = value.get(key);
        } else {
            value = value[key];
        }

        if (value === undefined) {
            break;
        }
    }

    return value;
}

export function GetFlag(keyString) {
    return getValueFromMap(keyString);
}
