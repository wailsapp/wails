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

// Drop target constants
const DROP_TARGET_ATTRIBUTE = 'data-file-drop-target';
const DROP_TARGET_ACTIVE_CLASS = 'file-drop-target-active';
let currentDropTarget: Element | null = null;

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
const FilesDropped                      = 50;
const PrintMethod                       = 51;

/**
 * Finds the nearest drop target element by walking up the DOM tree.
 */
function getDropTargetElement(element: Element | null): Element | null {
    if (!element) {
        return null;
    }
    return element.closest(`[${DROP_TARGET_ATTRIBUTE}]`);
}

/**
 * Check if we can use WebView2's postMessageWithAdditionalObjects (Windows)
 * Also checks that EnableFileDrop is true for this window.
 */
function canResolveFilePaths(): boolean {
    // Must have WebView2's postMessageWithAdditionalObjects API (Windows only)
    if ((window as any).chrome?.webview?.postMessageWithAdditionalObjects == null) {
        return false;
    }
    // Must have EnableFileDrop set to true for this window
    // This flag is set by the Go backend during runtime initialization
    return (window as any)._wails?.flags?.enableFileDrop === true;
}

/**
 * Send file drop to backend via WebView2 (Windows only)
 */
function resolveFilePaths(x: number, y: number, files: File[]): void {
    if ((window as any).chrome?.webview?.postMessageWithAdditionalObjects) {
        (window as any).chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
    }
}

// Native drag state (Linux/macOS intercept DOM drag events)
let nativeDragActive = false;

/**
 * Cleans up native drag state and hover effects.
 * Called on drop or when drag leaves the window.
 */
function cleanupNativeDrag(): void {
    nativeDragActive = false;
    if (currentDropTarget) {
        currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
        currentDropTarget = null;
    }
}

/**
 * Called from Go when a file drag enters the window on Linux/macOS.
 */
function handleDragEnter(): void {
    // Check if file drops are enabled for this window
    if ((window as any)._wails?.flags?.enableFileDrop === false) {
        return; // File drops disabled, don't activate drag state
    }
    nativeDragActive = true;
}

/**
 * Called from Go when a file drag leaves the window on Linux/macOS.
 */
function handleDragLeave(): void {
    cleanupNativeDrag();
}

/**
 * Called from Go during file drag to update hover state on Linux/macOS.
 * @param x - X coordinate in CSS pixels
 * @param y - Y coordinate in CSS pixels
 */
function handleDragOver(x: number, y: number): void {
    if (!nativeDragActive) return;
    
    // Check if file drops are enabled for this window
    if ((window as any)._wails?.flags?.enableFileDrop === false) {
        return; // File drops disabled, don't show hover effects
    }
    
    const targetElement = document.elementFromPoint(x, y);
    const dropTarget = getDropTargetElement(targetElement);
    
    if (currentDropTarget && currentDropTarget !== dropTarget) {
        currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
    }
    
    if (dropTarget) {
        dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
        currentDropTarget = dropTarget;
    } else {
        currentDropTarget = null;
    }
}



// Export the handlers for use by Go via index.ts
export { handleDragEnter, handleDragLeave, handleDragOver };

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
     * Handles file drops originating from platform-specific code (e.g., macOS/Linux native drag-and-drop).
     * Gathers information about the drop target element and sends it back to the Go backend.
     *
     * @param filenames - An array of file paths (strings) that were dropped.
     * @param x - The x-coordinate of the drop event (CSS pixels).
     * @param y - The y-coordinate of the drop event (CSS pixels).
     */
    HandlePlatformFileDrop(filenames: string[], x: number, y: number): void {
        // Check if file drops are enabled for this window
        if ((window as any)._wails?.flags?.enableFileDrop === false) {
            return; // File drops disabled, ignore the drop
        }
        
        const element = document.elementFromPoint(x, y);
        const dropTarget = getDropTargetElement(element);

        if (!dropTarget) {
            // Drop was not on a designated drop target - ignore
            return;
        }

        const elementDetails = {
            id: dropTarget.id,
            classList: Array.from(dropTarget.classList),
            attributes: {} as { [key: string]: string },
        };
        for (let i = 0; i < dropTarget.attributes.length; i++) {
            const attr = dropTarget.attributes[i];
            elementDetails.attributes[attr.name] = attr.value;
        }

        const payload = {
            filenames,
            x,
            y,
            elementDetails,
        };

        this[callerSym](FilesDropped, payload);
        
        // Clean up native drag state after drop
        cleanupNativeDrag();
    }
  
    /* Triggers Windows 11 Snap Assist feature (Windows only).
     * This is equivalent to pressing Win+Z and shows snap layout options.
     */
    SnapAssist(): Promise<void> {
        return this[callerSym](SnapAssistMethod);
    }

    /**
     * Opens the print dialog for the window.
     */
    Print(): Promise<void> {
        return this[callerSym](PrintMethod);
    }
}

/**
 * The window within which the script is running.
 */
const thisWindow = new Window('');

/**
 * Sets up global drag and drop event listeners for file drops.
 * Handles visual feedback (hover state) and file drop processing.
 */
function setupDropTargetListeners() {
    const docElement = document.documentElement;
    let dragEnterCounter = 0;

    docElement.addEventListener('dragenter', (event) => {
        if (!event.dataTransfer?.types.includes('Files')) {
            return; // Only handle file drags, let other drags pass through
        }
        event.preventDefault(); // Always prevent default to stop browser navigation
        // On Windows, check if file drops are enabled for this window
        if ((window as any)._wails?.flags?.enableFileDrop === false) {
            event.dataTransfer.dropEffect = 'none'; // Show "no drop" cursor
            return; // File drops disabled, don't show hover effects
        }
        dragEnterCounter++;
        
        const targetElement = document.elementFromPoint(event.clientX, event.clientY);
        const dropTarget = getDropTargetElement(targetElement);

        // Update hover state
        if (currentDropTarget && currentDropTarget !== dropTarget) {
            currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
        }

        if (dropTarget) {
            dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
            event.dataTransfer.dropEffect = 'copy';
            currentDropTarget = dropTarget;
        } else {
            event.dataTransfer.dropEffect = 'none';
            currentDropTarget = null;
        }
    }, false);

    docElement.addEventListener('dragover', (event) => {
        if (!event.dataTransfer?.types.includes('Files')) {
            return; // Only handle file drags
        }
        event.preventDefault(); // Always prevent default to stop browser navigation
        // On Windows, check if file drops are enabled for this window
        if ((window as any)._wails?.flags?.enableFileDrop === false) {
            event.dataTransfer.dropEffect = 'none'; // Show "no drop" cursor
            return; // File drops disabled, don't show hover effects
        }
        
        // Update drop target as cursor moves
        const targetElement = document.elementFromPoint(event.clientX, event.clientY);
        const dropTarget = getDropTargetElement(targetElement);
        
        if (currentDropTarget && currentDropTarget !== dropTarget) {
            currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
        }
        
        if (dropTarget) {
            if (!dropTarget.classList.contains(DROP_TARGET_ACTIVE_CLASS)) {
                dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
            }
            event.dataTransfer.dropEffect = 'copy';
            currentDropTarget = dropTarget;
        } else {
            event.dataTransfer.dropEffect = 'none';
            currentDropTarget = null;
        }
    }, false);

    docElement.addEventListener('dragleave', (event) => {
        if (!event.dataTransfer?.types.includes('Files')) {
            return;
        }
        event.preventDefault(); // Always prevent default to stop browser navigation
        // On Windows, check if file drops are enabled for this window
        if ((window as any)._wails?.flags?.enableFileDrop === false) {
            return;
        }
        
        // On Linux/WebKitGTK and macOS, dragleave fires immediately with relatedTarget=null when native
        // drag handling is involved. Ignore these spurious events - we'll clean up on drop instead.
        if (event.relatedTarget === null) {
            return;
        }
        
        dragEnterCounter--;
        
        if (dragEnterCounter === 0 || 
            (currentDropTarget && !currentDropTarget.contains(event.relatedTarget as Node))) {
            if (currentDropTarget) {
                currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
                currentDropTarget = null;
            }
            dragEnterCounter = 0;
        }
    }, false);

    docElement.addEventListener('drop', (event) => {
        if (!event.dataTransfer?.types.includes('Files')) {
            return; // Only handle file drops
        }
        event.preventDefault(); // Always prevent default to stop browser navigation
        // On Windows, check if file drops are enabled for this window
        if ((window as any)._wails?.flags?.enableFileDrop === false) {
            return;
        }
        dragEnterCounter = 0;
        
        if (currentDropTarget) {
            currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
            currentDropTarget = null;
        }

        // On Windows, handle file drops via JavaScript
        // On macOS/Linux, native code will call HandlePlatformFileDrop
        if (canResolveFilePaths()) {
            const files: File[] = [];
            if (event.dataTransfer.items) {
                for (const item of event.dataTransfer.items) {
                    if (item.kind === 'file') {
                        const file = item.getAsFile();
                        if (file) files.push(file);
                    }
                }
            } else if (event.dataTransfer.files) {
                for (const file of event.dataTransfer.files) {
                    files.push(file);
                }
            }
            
            if (files.length > 0) {
                resolveFilePaths(event.clientX, event.clientY, files);
            }
        }
    }, false);
}

// Initialize listeners when the script loads
if (typeof window !== "undefined" && typeof document !== "undefined") {
    setupDropTargetListeners();
}

export default thisWindow;
