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
    let $resultPromise = $Call.ByName("main.GreetService.Get", aliasValue) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType0($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

/**
 * Get someone quite different
 */
export function GetButDifferent(): Promise<$models.GenericPerson<boolean>> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.GreetService.GetButDifferent") as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

/**
 * Greet a lot of unusual things.
 */
export function Greet($0: $models.EmptyAliasStruct, $1: $models.EmptyStruct): Promise<$models.AliasStruct> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0, $1) as any;
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType5($result);
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

// Private type creation functions
const $$createType0 = $models.Person.createFrom;
const $$createType1 = $models.GenericPerson.createFrom($Create.Any);
const $$createType2 = $Create.Array($Create.Any);
const $$createType3 = $Create.Array($Create.Any);
const $$createType4 = $Create.Struct({
    "NoMoreIdeas": $$createType3,
});
const $$createType5 = $Create.Struct({
    "Foo": $$createType2,
    "Other": $$createType4,
});
