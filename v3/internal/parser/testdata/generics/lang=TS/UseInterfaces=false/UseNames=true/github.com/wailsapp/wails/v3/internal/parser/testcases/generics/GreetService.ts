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
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType0($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

/**
 * Get someone quite different
 */
export function GetButDifferent(): Promise<$models.Person<boolean>> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.GreetService.GetButDifferent");
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

// Private type creation functions
const $$createType0 = $models.Person.createFrom($Create.Any);
const $$createType1 = $models.Person.createFrom($Create.Any);