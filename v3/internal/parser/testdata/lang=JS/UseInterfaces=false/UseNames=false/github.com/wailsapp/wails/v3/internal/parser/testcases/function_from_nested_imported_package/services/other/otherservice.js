// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * OtherService is a struct
 * that does things
 * @module
 */

import {Call as $Call, Create as $Create} from "/wails/runtime.js";

import * as $models from "./models.js";

/**
 * Yay does this and that
 * @returns {Promise<$models.Address | null> & { cancel(): void }}
 */
export function Yay() {
    let $resultPromise = /** @type {any} */($Call.ByID(3249920254));
    let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
        return $$createType1($result);
    }));
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

// Private type creation functions
const $$createType0 = $models.Address.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
