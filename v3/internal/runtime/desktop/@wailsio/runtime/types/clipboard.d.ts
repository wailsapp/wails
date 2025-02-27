/**
 * Sets the text to the Clipboard.
 *
 * @param text - The text to be set to the Clipboard.
 * @return A Promise that resolves when the operation is successful.
 */
export declare function SetText(text: string): Promise<void>;
/**
 * Get the Clipboard text
 *
 * @returns A promise that resolves with the text from the Clipboard.
 */
export declare function Text(): Promise<string>;
