// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import * as $models from "./models.ts";

/**
 * Get someone
 */
export function Get(): Promise<$models.Person<string>> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.GreetService.Get");
    return $resultPromise as any;
}

/**
 * Get someone quite different
 */
export function GetButDifferent(): Promise<$models.Person<boolean>> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.GreetService.GetButDifferent");
    return $resultPromise as any;
}
