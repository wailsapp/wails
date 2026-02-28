/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { newRuntimeCaller, objectNames } from "./runtime.js";

const call = newRuntimeCaller(objectNames.Browser);

const BrowserOpenURL = 0;

/**
 * Open a browser window to the given URL.
 *
 * @param url - The URL to open
 */
export function OpenURL(url: string | URL): Promise<void> {
    const urlString = url.toString();
    const androidOpenURL = (window as any)?.wails?.openURL;
    if (typeof androidOpenURL === "function") {
        try {
            androidOpenURL.call((window as any).wails, urlString);
            return Promise.resolve();
        } catch (e) {
            return Promise.reject(e);
        }
    }
    return call(BrowserOpenURL, {url: urlString});
}
