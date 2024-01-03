/**
 * Sets the text to the Clipboard.
 *
 * @param {string} text - The text to be set to the Clipboard.
 * @return {Promise} - A Promise that resolves when the operation is successful.
 */
export function SetText(text: string): Promise<any>;
/**
 * Get the Clipboard text
 * @returns {Promise<string>} A promise that resolves with the text from the Clipboard.
 */
export function Text(): Promise<string>;
