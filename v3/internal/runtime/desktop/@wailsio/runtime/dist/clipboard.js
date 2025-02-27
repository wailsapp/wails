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
const call = newRuntimeCaller(objectNames.Clipboard);
const ClipboardSetText = 0;
const ClipboardText = 1;
/**
 * Sets the text to the Clipboard.
 *
 * @param text - The text to be set to the Clipboard.
 * @return A Promise that resolves when the operation is successful.
 */
export function SetText(text) {
    return call(ClipboardSetText, { text });
}
/**
 * Get the Clipboard text
 *
 * @returns A promise that resolves with the text from the Clipboard.
 */
export function Text() {
    return call(ClipboardText);
}
