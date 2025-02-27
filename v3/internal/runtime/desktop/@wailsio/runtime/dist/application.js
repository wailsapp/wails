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
const call = newRuntimeCaller(objectNames.Application);
const HideMethod = 0;
const ShowMethod = 1;
const QuitMethod = 2;
/**
 * Hides a certain method by calling the HideMethod function.
 */
export function Hide() {
    return call(HideMethod);
}
/**
 * Calls the ShowMethod and returns the result.
 */
export function Show() {
    return call(ShowMethod);
}
/**
 * Calls the QuitMethod to terminate the program.
 */
export function Quit() {
    return call(QuitMethod);
}
