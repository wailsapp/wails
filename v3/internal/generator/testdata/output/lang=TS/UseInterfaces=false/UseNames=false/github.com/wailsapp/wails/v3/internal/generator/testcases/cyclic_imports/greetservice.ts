// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

/**
 * Make a cycle.
 */
export function MakeCycles(): Promise<[$models.StructA, $models.StructC]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(440020721) as any;
    let $typingPromise = $resultPromise.then(($result: any) => {
        $result[0] = $$createType0($result[0]);
        $result[1] = $$createType1($result[1]);
        return $result;
    }) as any;
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

// Private type creation functions
const $$createType0 = $models.StructA.createFrom;
const $$createType1 = $models.StructC.createFrom;
