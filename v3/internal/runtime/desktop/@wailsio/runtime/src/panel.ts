/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { newRuntimeCaller, objectNames } from "./runtime.js";

// Panel method constants - must match messageprocessor_panel.go
const SetBoundsMethod    = 0;
const GetBoundsMethod    = 1;
const SetZIndexMethod    = 2;
const SetURLMethod       = 3;
const SetHTMLMethod      = 4;
const ExecJSMethod       = 5;
const ReloadMethod       = 6;
const ForceReloadMethod  = 7;
const ShowMethod         = 8;
const HideMethod         = 9;
const IsVisibleMethod    = 10;
const SetZoomMethod      = 11;
const GetZoomMethod      = 12;
const FocusMethod        = 13;
const IsFocusedMethod    = 14;
const OpenDevToolsMethod = 15;
const DestroyMethod      = 16;
const NameMethod         = 17;

/**
 * A record describing the bounds (position and size) of a panel.
 */
export interface Bounds {
    /** The X position of the panel within the window. */
    x: number;
    /** The Y position of the panel within the window. */
    y: number;
    /** The width of the panel. */
    width: number;
    /** The height of the panel. */
    height: number;
}

// Private field names
const callerSym = Symbol("caller");
const panelNameSym = Symbol("panelName");

/**
 * Panel represents an embedded webview panel within a window.
 * Panels allow embedding multiple webview instances in a single window,
 * similar to Electron's BrowserView/WebContentsView.
 */
export class Panel {
    // Private fields
    private [callerSym]: (method: number, args?: any) => Promise<any>;
    private [panelNameSym]: string;

    /**
     * Creates a new Panel instance.
     *
     * @param panelName - The name of the panel to control.
     * @param windowName - The name of the parent window (optional, defaults to current window).
     */
    constructor(panelName: string, windowName: string = '') {
        this[panelNameSym] = panelName;
        this[callerSym] = newRuntimeCaller(objectNames.Panel, windowName);

        // Bind instance methods for use in event handlers
        for (const method of Object.getOwnPropertyNames(Panel.prototype)) {
            if (method !== "constructor" && typeof (this as any)[method] === "function") {
                (this as any)[method] = (this as any)[method].bind(this);
            }
        }
    }

    /**
     * Gets a reference to the specified panel.
     *
     * @param panelName - The name of the panel to get.
     * @param windowName - The name of the parent window (optional).
     * @returns A new Panel instance.
     */
    static Get(panelName: string, windowName: string = ''): Panel {
        return new Panel(panelName, windowName);
    }

    /**
     * Sets the position and size of the panel.
     *
     * @param bounds - The new bounds for the panel.
     */
    SetBounds(bounds: Bounds): Promise<void> {
        return this[callerSym](SetBoundsMethod, {
            panel: this[panelNameSym],
            x: bounds.x,
            y: bounds.y,
            width: bounds.width,
            height: bounds.height
        });
    }

    /**
     * Gets the current position and size of the panel.
     *
     * @returns The current bounds of the panel.
     */
    GetBounds(): Promise<Bounds> {
        return this[callerSym](GetBoundsMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Sets the z-index (stacking order) of the panel.
     *
     * @param zIndex - The z-index value (higher values appear on top).
     */
    SetZIndex(zIndex: number): Promise<void> {
        return this[callerSym](SetZIndexMethod, {
            panel: this[panelNameSym],
            zIndex
        });
    }

    /**
     * Navigates the panel to the specified URL.
     *
     * @param url - The URL to navigate to.
     */
    SetURL(url: string): Promise<void> {
        return this[callerSym](SetURLMethod, {
            panel: this[panelNameSym],
            url
        });
    }

    /**
     * Sets the HTML content of the panel directly.
     *
     * @param html - The HTML content to load.
     */
    SetHTML(html: string): Promise<void> {
        return this[callerSym](SetHTMLMethod, {
            panel: this[panelNameSym],
            html
        });
    }

    /**
     * Executes JavaScript code in the panel's context.
     *
     * @param js - The JavaScript code to execute.
     */
    ExecJS(js: string): Promise<void> {
        return this[callerSym](ExecJSMethod, {
            panel: this[panelNameSym],
            js
        });
    }

    /**
     * Reloads the current page in the panel.
     */
    Reload(): Promise<void> {
        return this[callerSym](ReloadMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Forces a reload of the page, ignoring cached content.
     */
    ForceReload(): Promise<void> {
        return this[callerSym](ForceReloadMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Shows the panel (makes it visible).
     */
    Show(): Promise<void> {
        return this[callerSym](ShowMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Hides the panel (makes it invisible).
     */
    Hide(): Promise<void> {
        return this[callerSym](HideMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Checks if the panel is currently visible.
     *
     * @returns True if the panel is visible, false otherwise.
     */
    IsVisible(): Promise<boolean> {
        return this[callerSym](IsVisibleMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Sets the zoom level of the panel.
     *
     * @param zoom - The zoom level (1.0 = 100%).
     */
    SetZoom(zoom: number): Promise<void> {
        return this[callerSym](SetZoomMethod, {
            panel: this[panelNameSym],
            zoom
        });
    }

    /**
     * Gets the current zoom level of the panel.
     *
     * @returns The current zoom level.
     */
    GetZoom(): Promise<number> {
        return this[callerSym](GetZoomMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Focuses the panel (gives it keyboard input focus).
     */
    Focus(): Promise<void> {
        return this[callerSym](FocusMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Checks if the panel currently has keyboard focus.
     *
     * @returns True if the panel is focused, false otherwise.
     */
    IsFocused(): Promise<boolean> {
        return this[callerSym](IsFocusedMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Opens the developer tools for the panel.
     */
    OpenDevTools(): Promise<void> {
        return this[callerSym](OpenDevToolsMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Destroys the panel, removing it from the window.
     */
    Destroy(): Promise<void> {
        return this[callerSym](DestroyMethod, {
            panel: this[panelNameSym]
        });
    }

    /**
     * Gets the name of the panel.
     *
     * @returns The name of the panel.
     */
    Name(): Promise<string> {
        return this[callerSym](NameMethod, {
            panel: this[panelNameSym]
        });
    }
}
