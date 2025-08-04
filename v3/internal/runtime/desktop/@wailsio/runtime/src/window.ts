/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import {newRuntimeCaller, objectNames} from "./runtime.js";
import type { Screen } from "./screens.js";

// NEW: Dropzone constants
const DROPZONE_ATTRIBUTE = 'data-wails-dropzone';
const DROPZONE_HOVER_CLASS = 'wails-dropzone-hover'; // User can style this class
let currentHoveredDropzone: Element | null = null;

const PositionMethod                    = 0;
const CenterMethod                      = 1;
const CloseMethod                       = 2;
const DisableSizeConstraintsMethod      = 3;
const EnableSizeConstraintsMethod       = 4;
const FocusMethod                       = 5;
const ForceReloadMethod                 = 6;
const FullscreenMethod                  = 7;
const GetScreenMethod                   = 8;
const GetZoomMethod                     = 9;
const HeightMethod                      = 10;
const HideMethod                        = 11;
const IsFocusedMethod                   = 12;
const IsFullscreenMethod                = 13;
const IsMaximisedMethod                 = 14;
const IsMinimisedMethod                 = 15;
const MaximiseMethod                    = 16;
const MinimiseMethod                    = 17;
const NameMethod                        = 18;
const OpenDevToolsMethod                = 19;
const RelativePositionMethod            = 20;
const ReloadMethod                      = 21;
const ResizableMethod                   = 22;
const RestoreMethod                     = 23;
const SetPositionMethod                 = 24;
const SetAlwaysOnTopMethod              = 25;
const SetBackgroundColourMethod         = 26;
const SetFramelessMethod                = 27;
const SetFullscreenButtonEnabledMethod  = 28;
const SetMaxSizeMethod                  = 29;
const SetMinSizeMethod                  = 30;
const SetRelativePositionMethod         = 31;
const SetResizableMethod                = 32;
const SetSizeMethod                     = 33;
const SetTitleMethod                    = 34;
const SetZoomMethod                     = 35;
const ShowMethod                        = 36;
const SizeMethod                        = 37;
const ToggleFullscreenMethod            = 38;
const ToggleMaximiseMethod              = 39;
const ToggleFramelessMethod             = 40; 
const UnFullscreenMethod                = 41;
const UnMaximiseMethod                  = 42;
const UnMinimiseMethod                  = 43;
const WidthMethod                       = 44;
const ZoomMethod                        = 45;
const ZoomInMethod                      = 46;
const ZoomOutMethod                     = 47;
const ZoomResetMethod                   = 48;
const SnapAssistMethod                  = 49;
const WindowDropZoneDropped             = 50;

function getDropzoneElement(element: Element | null): Element | null {
    if (!element) {
        return null;
    }
    // Allow dropzone attribute to be on the element itself or any parent
    return element.closest(`[${DROPZONE_ATTRIBUTE}]`);
}

/**
 * A record describing the position of a window.
 */
interface Position {
    /** The horizontal position of the window. */
    x: number;
    /** The vertical position of the window. */
    y: number;
}

/**
 * A record describing the size of a window.
 */
interface Size {
    /** The width of the window. */
    width: number;
    /** The height of the window. */
    height: number;
}

// Private field names.
const callerSym = Symbol("caller");

class Window {
    // Private fields.
    private [callerSym]: (message: number, args?: any) => Promise<any>;

    /**
     * Initialises a window object with the specified name.
     *
     * @private
     * @param name - The name of the target window.
     */
    constructor(name: string = '') {
        this[callerSym] = newRuntimeCaller(objectNames.Window, name)

        // bind instance method to make them easily usable in event handlers
        for (const method of Object.getOwnPropertyNames(Window.prototype)) {
            if (
                method !== "constructor"
                && typeof (this as any)[method] === "function"
            ) {
                (this as any)[method] = (this as any)[method].bind(this);
            }
        }
    }

    /**
     * Gets the specified window.
     *
     * @param name - The name of the window to get.
     * @returns The corresponding window object.
     */
    Get(name: string): Window {
        return new Window(name);
    }

    /**
     * Returns the absolute position of the window.
     *
     * @returns The current absolute position of the window.
     */
    Position(): Promise<Position> {
        return this[callerSym](PositionMethod);
    }

    /**
     * Centers the window on the screen.
     */
    Center(): Promise<void> {
        return this[callerSym](CenterMethod);
    }

    /**
     * Closes the window.
     */
    Close(): Promise<void> {
        return this[callerSym](CloseMethod);
    }

    /**
     * Disables min/max size constraints.
     */
    DisableSizeConstraints(): Promise<void> {
        return this[callerSym](DisableSizeConstraintsMethod);
    }

    /**
     * Enables min/max size constraints.
     */
    EnableSizeConstraints(): Promise<void> {
        return this[callerSym](EnableSizeConstraintsMethod);
    }

    /**
     * Focuses the window.
     */
    Focus(): Promise<void> {
        return this[callerSym](FocusMethod);
    }

    /**
     * Forces the window to reload the page assets.
     */
    ForceReload(): Promise<void> {
        return this[callerSym](ForceReloadMethod);
    }

    /**
     * Switches the window to fullscreen mode.
     */
    Fullscreen(): Promise<void> {
        return this[callerSym](FullscreenMethod);
    }

    /**
     * Returns the screen that the window is on.
     *
     * @returns The screen the window is currently on.
     */
    GetScreen(): Promise<Screen> {
        return this[callerSym](GetScreenMethod);
    }

    /**
     * Returns the current zoom level of the window.
     *
     * @returns The current zoom level.
     */
    GetZoom(): Promise<number> {
        return this[callerSym](GetZoomMethod);
    }

    /**
     * Returns the height of the window.
     *
     * @returns The current height of the window.
     */
    Height(): Promise<number> {
        return this[callerSym](HeightMethod);
    }

    /**
     * Hides the window.
     */
    Hide(): Promise<void> {
        return this[callerSym](HideMethod);
    }

    /**
     * Returns true if the window is focused.
     *
     * @returns Whether the window is currently focused.
     */
    IsFocused(): Promise<boolean> {
        return this[callerSym](IsFocusedMethod);
    }

    /**
     * Returns true if the window is fullscreen.
     *
     * @returns Whether the window is currently fullscreen.
     */
    IsFullscreen(): Promise<boolean> {
        return this[callerSym](IsFullscreenMethod);
    }

    /**
     * Returns true if the window is maximised.
     *
     * @returns Whether the window is currently maximised.
     */
    IsMaximised(): Promise<boolean> {
        return this[callerSym](IsMaximisedMethod);
    }

    /**
     * Returns true if the window is minimised.
     *
     * @returns Whether the window is currently minimised.
     */
    IsMinimised(): Promise<boolean> {
        return this[callerSym](IsMinimisedMethod);
    }

    /**
     * Maximises the window.
     */
    Maximise(): Promise<void> {
        return this[callerSym](MaximiseMethod);
    }

    /**
     * Minimises the window.
     */
    Minimise(): Promise<void> {
        return this[callerSym](MinimiseMethod);
    }

    /**
     * Returns the name of the window.
     *
     * @returns The name of the window.
     */
    Name(): Promise<string> {
        return this[callerSym](NameMethod);
    }

    /**
     * Opens the development tools pane.
     */
    OpenDevTools(): Promise<void> {
        return this[callerSym](OpenDevToolsMethod);
    }

    /**
     * Returns the relative position of the window to the screen.
     *
     * @returns The current relative position of the window.
     */
    RelativePosition(): Promise<Position> {
        return this[callerSym](RelativePositionMethod);
    }

    /**
     * Reloads the page assets.
     */
    Reload(): Promise<void> {
        return this[callerSym](ReloadMethod);
    }

    /**
     * Returns true if the window is resizable.
     *
     * @returns Whether the window is currently resizable.
     */
    Resizable(): Promise<boolean> {
        return this[callerSym](ResizableMethod);
    }

    /**
     * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
     */
    Restore(): Promise<void> {
        return this[callerSym](RestoreMethod);
    }

    /**
     * Sets the absolute position of the window.
     *
     * @param x - The desired horizontal absolute position of the window.
     * @param y - The desired vertical absolute position of the window.
     */
    SetPosition(x: number, y: number): Promise<void> {
        return this[callerSym](SetPositionMethod, { x, y });
    }

    /**
     * Sets the window to be always on top.
     *
     * @param alwaysOnTop - Whether the window should stay on top.
     */
    SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void> {
        return this[callerSym](SetAlwaysOnTopMethod, { alwaysOnTop });
    }

    /**
     * Sets the background colour of the window.
     *
     * @param r - The desired red component of the window background.
     * @param g - The desired green component of the window background.
     * @param b - The desired blue component of the window background.
     * @param a - The desired alpha component of the window background.
     */
    SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void> {
        return this[callerSym](SetBackgroundColourMethod, { r, g, b, a });
    }

    /**
     * Removes the window frame and title bar.
     *
     * @param frameless - Whether the window should be frameless.
     */
    SetFrameless(frameless: boolean): Promise<void> {
        return this[callerSym](SetFramelessMethod, { frameless });
    }

    /**
     * Disables the system fullscreen button.
     *
     * @param enabled - Whether the fullscreen button should be enabled.
     */
    SetFullscreenButtonEnabled(enabled: boolean): Promise<void> {
        return this[callerSym](SetFullscreenButtonEnabledMethod, { enabled });
    }

    /**
     * Sets the maximum size of the window.
     *
     * @param width - The desired maximum width of the window.
     * @param height - The desired maximum height of the window.
     */
    SetMaxSize(width: number, height: number): Promise<void> {
        return this[callerSym](SetMaxSizeMethod, { width, height });
    }

    /**
     * Sets the minimum size of the window.
     *
     * @param width - The desired minimum width of the window.
     * @param height - The desired minimum height of the window.
     */
    SetMinSize(width: number, height: number): Promise<void> {
        return this[callerSym](SetMinSizeMethod, { width, height });
    }

    /**
     * Sets the relative position of the window to the screen.
     *
     * @param x - The desired horizontal relative position of the window.
     * @param y - The desired vertical relative position of the window.
     */
    SetRelativePosition(x: number, y: number): Promise<void> {
        return this[callerSym](SetRelativePositionMethod, { x, y });
    }

    /**
     * Sets whether the window is resizable.
     *
     * @param resizable - Whether the window should be resizable.
     */
    SetResizable(resizable: boolean): Promise<void> {
        return this[callerSym](SetResizableMethod, { resizable });
    }

    /**
     * Sets the size of the window.
     *
     * @param width - The desired width of the window.
     * @param height - The desired height of the window.
     */
    SetSize(width: number, height: number): Promise<void> {
        return this[callerSym](SetSizeMethod, { width, height });
    }

    /**
     * Sets the title of the window.
     *
     * @param title - The desired title of the window.
     */
    SetTitle(title: string): Promise<void> {
        return this[callerSym](SetTitleMethod, { title });
    }

    /**
     * Sets the zoom level of the window.
     *
     * @param zoom - The desired zoom level.
     */
    SetZoom(zoom: number): Promise<void> {
        return this[callerSym](SetZoomMethod, { zoom });
    }

    /**
     * Shows the window.
     */
    Show(): Promise<void> {
        return this[callerSym](ShowMethod);
    }

    /**
     * Returns the size of the window.
     *
     * @returns The current size of the window.
     */
    Size(): Promise<Size> {
        return this[callerSym](SizeMethod);
    }

    /**
     * Toggles the window between fullscreen and normal.
     */
    ToggleFullscreen(): Promise<void> {
        return this[callerSym](ToggleFullscreenMethod);
    }

    /**
     * Toggles the window between maximised and normal.
     */
    ToggleMaximise(): Promise<void> {
        return this[callerSym](ToggleMaximiseMethod);
    }

    /**
     * Toggles the window between frameless and normal.
     */
    ToggleFrameless(): Promise<void> {
        return this[callerSym](ToggleFramelessMethod);
    }

    /**
     * Un-fullscreens the window.
     */
    UnFullscreen(): Promise<void> {
        return this[callerSym](UnFullscreenMethod);
    }

    /**
     * Un-maximises the window.
     */
    UnMaximise(): Promise<void> {
        return this[callerSym](UnMaximiseMethod);
    }

    /**
     * Un-minimises the window.
     */
    UnMinimise(): Promise<void> {
        return this[callerSym](UnMinimiseMethod);
    }

    /**
     * Returns the width of the window.
     *
     * @returns The current width of the window.
     */
    Width(): Promise<number> {
        return this[callerSym](WidthMethod);
    }

    /**
     * Zooms the window.
     */
    Zoom(): Promise<void> {
        return this[callerSym](ZoomMethod);
    }

    /**
     * Increases the zoom level of the webview content.
     */
    ZoomIn(): Promise<void> {
        return this[callerSym](ZoomInMethod);
    }

    /**
     * Decreases the zoom level of the webview content.
     */
    ZoomOut(): Promise<void> {
        return this[callerSym](ZoomOutMethod);
    }

    /**
     * Resets the zoom level of the webview content.
     */
    ZoomReset(): Promise<void> {
        return this[callerSym](ZoomResetMethod);
    }

    /**
     * Handles file drops originating from platform-specific code (e.g., macOS native drag-and-drop).
     * Gathers information about the drop target element and sends it back to the Go backend.
     *
     * @param filenames - An array of file paths (strings) that were dropped.
     * @param x - The x-coordinate of the drop event.
     * @param y - The y-coordinate of the drop event.
     */
    HandlePlatformFileDrop(filenames: string[], x: number, y: number): void {
        const element = document.elementFromPoint(x, y);

        // NEW: Check if the drop target is a valid dropzone
        const dropzoneTarget = getDropzoneElement(element);

        if (!dropzoneTarget) {
            console.log(`Wails Runtime: Drop on element (or no element) at ${x},${y} which is not a designated dropzone. Ignoring. Element:`, element);
            // No need to call backend if not a valid dropzone target
            return;
        }

        console.log(`Wails Runtime: Drop on designated dropzone. Element at (${x}, ${y}):`, element, 'Effective dropzone:', dropzoneTarget);
        const elementDetails = {
            id: dropzoneTarget.id,
            classList: Array.from(dropzoneTarget.classList),
            attributes: {} as { [key: string]: string },
        };
        for (let i = 0; i < dropzoneTarget.attributes.length; i++) {
            const attr = dropzoneTarget.attributes[i];
            elementDetails.attributes[attr.name] = attr.value;
        }

        const payload = {
            filenames,
            x,
            y,
            elementDetails,
        };

        this[callerSym](WindowDropZoneDropped, payload);
    }
  
    /* Triggers Windows 11 Snap Assist feature (Windows only).
     * This is equivalent to pressing Win+Z and shows snap layout options.
     */
    SnapAssist(): Promise<void> {
        return this[callerSym](SnapAssistMethod);
    }
}

/**
 * The window within which the script is running.
 */
const thisWindow = new Window('');

// NEW: Global Drag Event Listeners
function setupGlobalDropzoneListeners() {
    const docElement = document.documentElement;
    let dragEnterCounter = 0; // To handle dragenter/dragleave on child elements

    docElement.addEventListener('dragenter', (event) => {
        event.preventDefault();
        if (event.dataTransfer && event.dataTransfer.types.includes('Files')) {
            dragEnterCounter++;
            const targetElement = document.elementFromPoint(event.clientX, event.clientY);
            const dropzone = getDropzoneElement(targetElement);

            // Clear previous hover regardless, then apply new if valid
            if (currentHoveredDropzone && currentHoveredDropzone !== dropzone) {
                currentHoveredDropzone.classList.remove(DROPZONE_HOVER_CLASS);
            }

            if (dropzone) {
                dropzone.classList.add(DROPZONE_HOVER_CLASS);
                event.dataTransfer.dropEffect = 'copy';
                currentHoveredDropzone = dropzone;
            } else {
                event.dataTransfer.dropEffect = 'none';
                currentHoveredDropzone = null; // Ensure it's cleared if no dropzone found
            }
        }
    }, false);

    docElement.addEventListener('dragover', (event) => {
        event.preventDefault(); // Necessary to allow drop
        if (event.dataTransfer && event.dataTransfer.types.includes('Files')) {
            // No need to query elementFromPoint again if already handled by dragenter correctly
            // Just ensure dropEffect is continuously set based on currentHoveredDropzone
            if (currentHoveredDropzone) {
                 // Re-apply class just in case it was removed by some other JS
                if(!currentHoveredDropzone.classList.contains(DROPZONE_HOVER_CLASS)) {
                    currentHoveredDropzone.classList.add(DROPZONE_HOVER_CLASS);
                }
                event.dataTransfer.dropEffect = 'copy';
            } else {
                event.dataTransfer.dropEffect = 'none';
            }
        }
    }, false);

    docElement.addEventListener('dragleave', (event) => {
        event.preventDefault();
        if (event.dataTransfer && event.dataTransfer.types.includes('Files')) {
            dragEnterCounter--;
            // Only remove hover if drag truly left the window or the last dropzone
            if (dragEnterCounter === 0 || event.relatedTarget === null || (currentHoveredDropzone && !currentHoveredDropzone.contains(event.relatedTarget as Node))) {
                if (currentHoveredDropzone) {
                    currentHoveredDropzone.classList.remove(DROPZONE_HOVER_CLASS);
                    currentHoveredDropzone = null;
                }
                dragEnterCounter = 0; // Reset counter if it went negative or left window
            }
        }
    }, false);

    docElement.addEventListener('drop', (event) => {
        event.preventDefault(); // Prevent default browser file handling
        dragEnterCounter = 0; // Reset counter
        if (currentHoveredDropzone) {
            currentHoveredDropzone.classList.remove(DROPZONE_HOVER_CLASS);
            currentHoveredDropzone = null;
        }
        // The actual drop processing is initiated by the native side calling HandlePlatformFileDrop
        // HandlePlatformFileDrop will then check if the drop was on a valid zone.
    }, false);
}

// Initialize listeners when the script loads
if (typeof window !== "undefined" && typeof document !== "undefined") {
    setupGlobalDropzoneListeners();
}

export default thisWindow;
