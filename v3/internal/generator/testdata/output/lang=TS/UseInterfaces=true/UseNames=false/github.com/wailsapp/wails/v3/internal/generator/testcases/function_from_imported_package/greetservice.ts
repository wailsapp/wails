// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./internal.js";

/**
 * Greet does XYZ
 */
export function Greet(name: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1411160069, name) as any;
    return $resultPromise;
}

/**
 * NewPerson creates a new person
 */
export function NewPerson(name: string): Promise<$models.Person | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1661412647, name) as any;
    return $resultPromise;
}
