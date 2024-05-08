// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "/wails/runtime.js";

import * as $models from "./models.ts";

/**
 * Get someone
 */
export function Get(aliasValue: $models.Alias): Promise<$models.Person> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1928502664, aliasValue);
    return $resultPromise as any;
}

/**
 * Get someone quite different
 */
export function GetButDifferent(): Promise<$models.GenericPerson<boolean>> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2240931744);
    return $resultPromise as any;
}

/**
 * Greet a lot of unusual things.
 */
export function Greet($0: $models.EmptyAliasStruct, $1: $models.EmptyStruct): Promise<$models.AliasStruct> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1411160069, $0, $1);
    return $resultPromise as any;
}
