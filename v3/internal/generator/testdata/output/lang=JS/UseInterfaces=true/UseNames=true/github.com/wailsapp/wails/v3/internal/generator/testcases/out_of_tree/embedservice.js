// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * EmbedService is tricky.
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import { Call as $Call, CancellablePromise as $CancellablePromise } from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as nobindingshere$0 from "../no_bindings_here/models.js";

/**
 * LikeThisOne is an example method that does nothing.
 * @returns {$CancellablePromise<[nobindingshere$0.Person, nobindingshere$0.HowDifferent<boolean>, nobindingshere$0.PrivatePerson]>}
 */
export function LikeThisOne() {
    return $Call.ByName("main.EmbedService.LikeThisOne");
}

/**
 * LikeThisOtherOne does nothing as well, but is different.
 * @returns {$CancellablePromise<void>}
 */
export function LikeThisOtherOne() {
    return $Call.ByName("main.EmbedService.LikeThisOtherOne");
}
