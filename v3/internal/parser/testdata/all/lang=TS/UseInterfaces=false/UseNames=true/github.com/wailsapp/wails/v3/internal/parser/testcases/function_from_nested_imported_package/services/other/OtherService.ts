// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * OtherService is a struct
 * that does things
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import {Address} from "./models.ts";

/**
 * Yay does this and that
 */
export function Yay(): Promise<Address | null> {
    let $resultPromise = $Call.ByName("github.com/wailsapp/wails/v3/internal/parser/testcases/function_from_nested_imported_package/services/other.OtherService.Yay");
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType0($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

// Internal type creation functions
const $$createType0 = $Create.Nullable(Address.createFrom);
