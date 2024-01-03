/**
 * Gets all screens.
 * @returns {Promise<Screen[]>} A promise that resolves to an array of Screen objects.
 */
export function GetAll(): Promise<Screen[]>;
/**
 * Gets the primary screen.
 * @returns {Promise<Screen>} A promise that resolves to the primary screen.
 */
export function GetPrimary(): Promise<Screen>;
/**
 * Gets the current active screen.
 *
 * @returns {Promise<Screen>} A promise that resolves with the current active screen.
 */
export function GetCurrent(): Promise<Screen>;
export type Screen = any;
