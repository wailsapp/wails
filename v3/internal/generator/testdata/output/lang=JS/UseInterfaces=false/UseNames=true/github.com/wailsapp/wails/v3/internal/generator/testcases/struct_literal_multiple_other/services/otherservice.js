// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * OtherService is a struct
 * that does things
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

/**
 * Yay does this and that
 * @returns {Promise<$models.Address | null> & { cancel(): void }}
 */
export function Yay() {
    let $resultPromise = /** @type {any} */($Call.ByName("github.com/wailsapp/wails/v3/internal/generator/testcases/struct_literal_multiple_other/services.OtherService.Yay"));
    let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
        return $$createType1($result);
    }));
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

// Private type creation functions
const $$createType0 = $models.Address.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
