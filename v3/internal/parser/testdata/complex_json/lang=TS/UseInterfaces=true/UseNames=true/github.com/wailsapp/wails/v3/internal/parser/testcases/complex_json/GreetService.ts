// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import {type Embedded1, type Person} from "./models.ts";

/**
 * Greet does XYZ
 */
export function Greet(person: Person, emb: Embedded1): Promise<string> {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", person, emb);
    return $resultPromise as any;
}
