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
export type Screen = {
    /**
     * - Unique identifier for the screen.
     */
    Id: string;
    /**
     * - Human readable name of the screen.
     */
    Name: string;
    /**
     * - The resolution scale of the screen. 1 = standard resolution, 2 = high (Retina), etc.
     */
    Scale: number;
    /**
     * - Contains the X and Y coordinates of the screen's position.
     */
    Position: Position;
    /**
     * - Contains the width and height of the screen.
     */
    Size: Size;
    /**
     * - Contains the bounds of the screen in terms of X, Y, Width, and Height.
     */
    Bounds: Rect;
    /**
     * - Contains the area of the screen that is actually usable (excluding taskbar and other system UI).
     */
    WorkArea: Rect;
    /**
     * - True if this is the primary monitor selected by the user in the operating system.
     */
    IsPrimary: boolean;
    /**
     * - The rotation of the screen.
     */
    Rotation: Rotation;
};
