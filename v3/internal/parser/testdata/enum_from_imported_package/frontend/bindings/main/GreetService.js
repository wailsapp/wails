// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from '@wailsio/runtime';
/**
 * @typedef {import('../services/models').Title} servicesTitle
 */

/**
 * Greet does XYZ
 * @function Greet
 * @param name {string}
 * @param title {servicesTitle}
 * @returns {Promise<string>}
 **/
export async function Greet(name, title) {
	return Call.ByID(1411160069, ...Array.prototype.slice.call(arguments, 0));
}
