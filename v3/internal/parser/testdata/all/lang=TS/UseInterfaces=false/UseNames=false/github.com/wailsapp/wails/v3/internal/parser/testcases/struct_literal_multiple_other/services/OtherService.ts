// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * OtherService is a struct
 * that does things
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import * as $models from "./models.ts";

/**
 * Yay does this and that
 */
export function Yay(): Promise<$models.Address | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1956770239);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

// Private type creation functions
const $$createType0 = $models.Address.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
