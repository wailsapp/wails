// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "/wails/runtime.js";

import * as $models from "./models.ts";

/**
 * Greet does XYZ
 * It has a multiline doc comment
 * The comment has even some * / traps!!
 */
export function Greet(str: string, people: $models.Person[] | null, $2: {"AnotherCount": number, "AnotherOne": $models.Person | null}, assoc: { [_: `${number}`]: boolean | null } | null, $4: (number | null)[] | null, ...other: string[]): Promise<[$models.Person, any, number[] | null]> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", str, people, $2, assoc, $4, other) as any;
    return $resultPromise;
}
