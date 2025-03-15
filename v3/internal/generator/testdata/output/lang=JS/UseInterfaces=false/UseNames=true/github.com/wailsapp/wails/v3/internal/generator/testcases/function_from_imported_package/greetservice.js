// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import { Call as $Call, CancellablePromise as $CancellablePromise, Create as $Create } from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

/**
 * Greet does XYZ
 * @param {string} name
 * @returns {$CancellablePromise<string>}
 */
export function Greet(name) {
    return $Call.ByName("main.GreetService.Greet", name);
}

/**
 * NewPerson creates a new person
 * @param {string} name
 * @returns {$CancellablePromise<$models.Person | null>}
 */
export function NewPerson(name) {
    return $Call.ByName("main.GreetService.NewPerson", name).then(/** @type {($result: any) => any} */(($result) => {
        return $$createType1($result);
    }));
}

// Private type creation functions
const $$createType0 = $models.Person.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
