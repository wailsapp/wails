// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * SomeMethods exports some methods.
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import * as $models from "./models.js";

/**
 * LikeThisOne is an example method that does nothing.
 * @returns {Promise<[$models.Person, $models.Impersonator, $models.HowDifferent<boolean>, $models.PrivatePerson]> & { cancel(): void }}
 */
export function LikeThisOne() {
    let $resultPromise = $Call.ByName("github.com/wailsapp/wails/v3/internal/parser/testcases/no_bindings_here.SomeMethods.LikeThisOne");
    let $typingPromise = $resultPromise.then(($result) => {
        $result[0] = $$createType0($result[0]);
        $result[1] = $$createType1($result[1]);
        $result[2] = $$createType2($result[2]);
        $result[3] = $$createType3($result[3]);
        return $result;
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

/**
 * LikeThisOtherOne does nothing as well, but is different.
 * @returns {Promise<void> & { cancel(): void }}
 */
export function LikeThisOtherOne() {
    let $resultPromise = $Call.ByName("github.com/wailsapp/wails/v3/internal/parser/testcases/no_bindings_here.SomeMethods.LikeThisOtherOne");
    return /** @type {any} */($resultPromise);
}

// Private type creation functions
const $$createType0 = $models.Person.createFrom;
const $$createType1 = $models.Impersonator.createFrom;
const $$createType2 = $models.HowDifferent.createFrom($Create.Any);
const $$createType3 = $models.PrivatePerson.createFrom;
