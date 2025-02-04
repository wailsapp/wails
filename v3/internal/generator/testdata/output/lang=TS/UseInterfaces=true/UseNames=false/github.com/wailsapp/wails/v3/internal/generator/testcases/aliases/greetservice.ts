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
import * as nobindingshere$0 from "../no_bindings_here/models.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

/**
 * Get someone.
 */
export function Get(aliasValue: $models.Alias): Promise<$models.Person> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1928502664, aliasValue) as any;
    return $resultPromise;
}

/**
 * Apparently, aliases are all the rage right now.
 */
export function GetButAliased(p: $models.AliasedPerson): Promise<$models.StrangelyAliasedPerson> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1896499664, p) as any;
    return $resultPromise;
}

/**
 * Get someone quite different.
 */
export function GetButDifferent(): Promise<$models.GenericPerson<boolean>> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2240931744) as any;
    return $resultPromise;
}

export function GetButForeignPrivateAlias(): Promise<nobindingshere$0.PrivatePerson> & { cancel(): void } {
    let $resultPromise = $Call.ByID(643456960) as any;
    return $resultPromise;
}

export function GetButGenericAliases(): Promise<$models.AliasGroup> & { cancel(): void } {
    let $resultPromise = $Call.ByID(914093800) as any;
    return $resultPromise;
}

/**
 * Greet a lot of unusual things.
 */
export function Greet($0: $models.EmptyAliasStruct, $1: $models.EmptyStruct): Promise<$models.AliasStruct> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1411160069, $0, $1) as any;
    return $resultPromise;
}
