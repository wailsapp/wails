// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * OtherService is a struct
 * that does things
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import {Address} from "./models.js";

/**
 * Yay does this and that
 * @returns {Promise<Address | null>}
 */
export function Yay() {
    let $resultPromise = $Call.ByID(3249920254);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType0($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

// Internal type creation functions
const $$createType0 = $Create.Nullable(Address.createFrom);
