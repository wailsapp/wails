/**
 * Handles the callback from a dialog.
 *
 * @param {string} id - The ID of the dialog response.
 * @param {string} data - The data received from the dialog.
 * @param {boolean} isJSON - Flag indicating whether the data is in JSON format.
 *
 * @return {undefined}
 */
export function dialogCallback(id: string, data: string, isJSON: boolean): undefined;
/**
 * Callback function for handling errors in dialog.
 *
 * @param {string} id - The id of the dialog response.
 * @param {string} message - The error message.
 *
 * @return {void}
 */
export function dialogErrorCallback(id: string, message: string): void;
export function Info(options: any): Promise<string>;
export function Warning(options: any): Promise<string>;
export function Error(options: any): Promise<string>;
export function Question(options: any): Promise<string>;
export function OpenFile(options: any): Promise<string[] | string>;
export function SaveFile(options: any): Promise<string>;
export type MessageDialogOptions = any;
export type OpenDialogOptions = any;
export type SaveDialogOptions = any;
