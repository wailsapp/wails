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

function runtimeCall(method: string, args?: any) : Promise<any> {
    let url = new URL(runtimeURL);
    url.searchParams.append("method", method);
    if (args) {
        url.searchParams.append("args", JSON.stringify(args));
    }
    return new Promise((resolve, reject) => {
        fetch(url)
            .then(response => {
                if (response && response.ok) {
                    // check content type
                    let content_type = response.headers.get("Content-Type");
                    if ( content_type && content_type.indexOf("application/json") !== -1) {
                        return response.json();
                    } else {
                        return response.text();
                    }
                }
                reject(Error(response.statusText));
            })
            .then(data => resolve(data))
            .catch(error => reject(error));
    });
}

export function newRuntimeCaller(object: string, id?: number) {
    if (!id || id === -1) {
        return function (method: string, args?: any) {
            return runtimeCall(object + "." + method, args);
        };
    }
    return function (method: string, args?: any) {
        args = args || {};
        args["windowID"] = id;
        return runtimeCall(object + "." + method, args);
    }
}