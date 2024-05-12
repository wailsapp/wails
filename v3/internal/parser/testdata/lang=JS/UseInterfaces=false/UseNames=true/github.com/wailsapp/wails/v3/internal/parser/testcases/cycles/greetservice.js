// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "/wails/runtime.js";

import * as $models from "./models.js";

/**
 * Make a cycle.
 * @returns {Promise<[$models.Cyclic, $models.GenericCyclic<$models.GenericCyclic<number>>]> & { cancel(): void }}
 */
export function MakeCycles() {
    let $resultPromise = $Call.ByName("main.GreetService.MakeCycles");
    let $typingPromise = $resultPromise.then(($result) => {
        $result[0] = $$createType0($result[0]);
        $result[1] = $$createType9($result[1]);
        return $result;
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

// Private type creation functions
var $$createType0 = function $$initCreateType0(...args) {
    if ($$createType0 === $$initCreateType0) {
        $$createType0 = $$createType3;
    }
	return $$createType0(...args);
};
const $$createType1 = $Create.Nullable($$createType0);
const $$createType2 = $Create.Map($Create.Any, $$createType1);
const $$createType3 = $Create.Array($$createType2);
var $$createType4 = function $$initCreateType4(...args) {
    if ($$createType4 === $$initCreateType4) {
        $$createType4 = $$createType8;
    }
	return $$createType4(...args);
};
const $$createType5 = $Create.Nullable($$createType4);
const $$createType6 = $Create.Array($Create.Any);
const $$createType7 = $Create.Struct({
    "X": $$createType5,
    "Y": $$createType6,
});
const $$createType8 = $Create.Array($$createType7);
var $$createType9 = function $$initCreateType9(...args) {
    if ($$createType9 === $$initCreateType9) {
        $$createType9 = $$createType13;
    }
	return $$createType9(...args);
};
const $$createType10 = $Create.Nullable($$createType9);
const $$createType11 = $Create.Array($$createType4);
const $$createType12 = $Create.Struct({
    "X": $$createType10,
    "Y": $$createType11,
});
const $$createType13 = $Create.Array($$createType12);
