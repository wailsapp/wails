// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from "@wailsio/runtime";

/**
 * @returns {Promise<void>}
 */
export function Hello() {
    return Call.ByName("main.OtherService.Hello");
}
