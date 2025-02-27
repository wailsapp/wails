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
const PositionMethod = 0;
const CenterMethod = 1;
const CloseMethod = 2;
const DisableSizeConstraintsMethod = 3;
const EnableSizeConstraintsMethod = 4;
const FocusMethod = 5;
const ForceReloadMethod = 6;
const FullscreenMethod = 7;
const GetScreenMethod = 8;
const GetZoomMethod = 9;
const HeightMethod = 10;
const HideMethod = 11;
const IsFocusedMethod = 12;
const IsFullscreenMethod = 13;
const IsMaximisedMethod = 14;
const IsMinimisedMethod = 15;
const MaximiseMethod = 16;
const MinimiseMethod = 17;
const NameMethod = 18;
const OpenDevToolsMethod = 19;
const RelativePositionMethod = 20;
const ReloadMethod = 21;
const ResizableMethod = 22;
const RestoreMethod = 23;
const SetPositionMethod = 24;
const SetAlwaysOnTopMethod = 25;
const SetBackgroundColourMethod = 26;
const SetFramelessMethod = 27;
const SetFullscreenButtonEnabledMethod = 28;
const SetMaxSizeMethod = 29;
const SetMinSizeMethod = 30;
const SetRelativePositionMethod = 31;
const SetResizableMethod = 32;
const SetSizeMethod = 33;
const SetTitleMethod = 34;
const SetZoomMethod = 35;
const ShowMethod = 36;
const SizeMethod = 37;
const ToggleFullscreenMethod = 38;
const ToggleMaximiseMethod = 39;
const UnFullscreenMethod = 40;
const UnMaximiseMethod = 41;
const UnMinimiseMethod = 42;
const WidthMethod = 43;
const ZoomMethod = 44;
const ZoomInMethod = 45;
const ZoomOutMethod = 46;
const ZoomResetMethod = 47;
// Private field names.
const callerSym = Symbol("caller");
class Window {
    /**
     * Initialises a window object with the specified name.
     *
     * @private
     * @param name - The name of the target window.
     */
    constructor(name = '') {
        this[callerSym] = newRuntimeCaller(objectNames.Window, name);
        // bind instance method to make them easily usable in event handlers
        for (const method of Object.getOwnPropertyNames(Window.prototype)) {
            if (method !== "constructor"
                && typeof this[method] === "function") {
                this[method] = this[method].bind(this);
            }
        }
    }
    /**
     * Gets the specified window.
     *
     * @param name - The name of the window to get.
     * @returns The corresponding window object.
     */
    Get(name) {
        return new Window(name);
    }
    /**
     * Returns the absolute position of the window.
     *
     * @returns The current absolute position of the window.
     */
    Position() {
        return this[callerSym](PositionMethod);
    }
    /**
     * Centers the window on the screen.
     */
    Center() {
        return this[callerSym](CenterMethod);
    }
    /**
     * Closes the window.
     */
    Close() {
        return this[callerSym](CloseMethod);
    }
    /**
     * Disables min/max size constraints.
     */
    DisableSizeConstraints() {
        return this[callerSym](DisableSizeConstraintsMethod);
    }
    /**
     * Enables min/max size constraints.
     */
    EnableSizeConstraints() {
        return this[callerSym](EnableSizeConstraintsMethod);
    }
    /**
     * Focuses the window.
     */
    Focus() {
        return this[callerSym](FocusMethod);
    }
    /**
     * Forces the window to reload the page assets.
     */
    ForceReload() {
        return this[callerSym](ForceReloadMethod);
    }
    /**
     * Switches the window to fullscreen mode.
     */
    Fullscreen() {
        return this[callerSym](FullscreenMethod);
    }
    /**
     * Returns the screen that the window is on.
     *
     * @returns The screen the window is currently on.
     */
    GetScreen() {
        return this[callerSym](GetScreenMethod);
    }
    /**
     * Returns the current zoom level of the window.
     *
     * @returns The current zoom level.
     */
    GetZoom() {
        return this[callerSym](GetZoomMethod);
    }
    /**
     * Returns the height of the window.
     *
     * @returns The current height of the window.
     */
    Height() {
        return this[callerSym](HeightMethod);
    }
    /**
     * Hides the window.
     */
    Hide() {
        return this[callerSym](HideMethod);
    }
    /**
     * Returns true if the window is focused.
     *
     * @returns Whether the window is currently focused.
     */
    IsFocused() {
        return this[callerSym](IsFocusedMethod);
    }
    /**
     * Returns true if the window is fullscreen.
     *
     * @returns Whether the window is currently fullscreen.
     */
    IsFullscreen() {
        return this[callerSym](IsFullscreenMethod);
    }
    /**
     * Returns true if the window is maximised.
     *
     * @returns Whether the window is currently maximised.
     */
    IsMaximised() {
        return this[callerSym](IsMaximisedMethod);
    }
    /**
     * Returns true if the window is minimised.
     *
     * @returns Whether the window is currently minimised.
     */
    IsMinimised() {
        return this[callerSym](IsMinimisedMethod);
    }
    /**
     * Maximises the window.
     */
    Maximise() {
        return this[callerSym](MaximiseMethod);
    }
    /**
     * Minimises the window.
     */
    Minimise() {
        return this[callerSym](MinimiseMethod);
    }
    /**
     * Returns the name of the window.
     *
     * @returns The name of the window.
     */
    Name() {
        return this[callerSym](NameMethod);
    }
    /**
     * Opens the development tools pane.
     */
    OpenDevTools() {
        return this[callerSym](OpenDevToolsMethod);
    }
    /**
     * Returns the relative position of the window to the screen.
     *
     * @returns The current relative position of the window.
     */
    RelativePosition() {
        return this[callerSym](RelativePositionMethod);
    }
    /**
     * Reloads the page assets.
     */
    Reload() {
        return this[callerSym](ReloadMethod);
    }
    /**
     * Returns true if the window is resizable.
     *
     * @returns Whether the window is currently resizable.
     */
    Resizable() {
        return this[callerSym](ResizableMethod);
    }
    /**
     * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
     */
    Restore() {
        return this[callerSym](RestoreMethod);
    }
    /**
     * Sets the absolute position of the window.
     *
     * @param x - The desired horizontal absolute position of the window.
     * @param y - The desired vertical absolute position of the window.
     */
    SetPosition(x, y) {
        return this[callerSym](SetPositionMethod, { x, y });
    }
    /**
     * Sets the window to be always on top.
     *
     * @param alwaysOnTop - Whether the window should stay on top.
     */
    SetAlwaysOnTop(alwaysOnTop) {
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
    SetBackgroundColour(r, g, b, a) {
        return this[callerSym](SetBackgroundColourMethod, { r, g, b, a });
    }
    /**
     * Removes the window frame and title bar.
     *
     * @param frameless - Whether the window should be frameless.
     */
    SetFrameless(frameless) {
        return this[callerSym](SetFramelessMethod, { frameless });
    }
    /**
     * Disables the system fullscreen button.
     *
     * @param enabled - Whether the fullscreen button should be enabled.
     */
    SetFullscreenButtonEnabled(enabled) {
        return this[callerSym](SetFullscreenButtonEnabledMethod, { enabled });
    }
    /**
     * Sets the maximum size of the window.
     *
     * @param width - The desired maximum width of the window.
     * @param height - The desired maximum height of the window.
     */
    SetMaxSize(width, height) {
        return this[callerSym](SetMaxSizeMethod, { width, height });
    }
    /**
     * Sets the minimum size of the window.
     *
     * @param width - The desired minimum width of the window.
     * @param height - The desired minimum height of the window.
     */
    SetMinSize(width, height) {
        return this[callerSym](SetMinSizeMethod, { width, height });
    }
    /**
     * Sets the relative position of the window to the screen.
     *
     * @param x - The desired horizontal relative position of the window.
     * @param y - The desired vertical relative position of the window.
     */
    SetRelativePosition(x, y) {
        return this[callerSym](SetRelativePositionMethod, { x, y });
    }
    /**
     * Sets whether the window is resizable.
     *
     * @param resizable - Whether the window should be resizable.
     */
    SetResizable(resizable) {
        return this[callerSym](SetResizableMethod, { resizable });
    }
    /**
     * Sets the size of the window.
     *
     * @param width - The desired width of the window.
     * @param height - The desired height of the window.
     */
    SetSize(width, height) {
        return this[callerSym](SetSizeMethod, { width, height });
    }
    /**
     * Sets the title of the window.
     *
     * @param title - The desired title of the window.
     */
    SetTitle(title) {
        return this[callerSym](SetTitleMethod, { title });
    }
    /**
     * Sets the zoom level of the window.
     *
     * @param zoom - The desired zoom level.
     */
    SetZoom(zoom) {
        return this[callerSym](SetZoomMethod, { zoom });
    }
    /**
     * Shows the window.
     */
    Show() {
        return this[callerSym](ShowMethod);
    }
    /**
     * Returns the size of the window.
     *
     * @returns The current size of the window.
     */
    Size() {
        return this[callerSym](SizeMethod);
    }
    /**
     * Toggles the window between fullscreen and normal.
     */
    ToggleFullscreen() {
        return this[callerSym](ToggleFullscreenMethod);
    }
    /**
     * Toggles the window between maximised and normal.
     */
    ToggleMaximise() {
        return this[callerSym](ToggleMaximiseMethod);
    }
    /**
     * Un-fullscreens the window.
     */
    UnFullscreen() {
        return this[callerSym](UnFullscreenMethod);
    }
    /**
     * Un-maximises the window.
     */
    UnMaximise() {
        return this[callerSym](UnMaximiseMethod);
    }
    /**
     * Un-minimises the window.
     */
    UnMinimise() {
        return this[callerSym](UnMinimiseMethod);
    }
    /**
     * Returns the width of the window.
     *
     * @returns The current width of the window.
     */
    Width() {
        return this[callerSym](WidthMethod);
    }
    /**
     * Zooms the window.
     */
    Zoom() {
        return this[callerSym](ZoomMethod);
    }
    /**
     * Increases the zoom level of the webview content.
     */
    ZoomIn() {
        return this[callerSym](ZoomInMethod);
    }
    /**
     * Decreases the zoom level of the webview content.
     */
    ZoomOut() {
        return this[callerSym](ZoomOutMethod);
    }
    /**
     * Resets the zoom level of the webview content.
     */
    ZoomReset() {
        return this[callerSym](ZoomResetMethod);
    }
}
/**
 * The window within which the script is running.
 */
const thisWindow = new Window('');
export default thisWindow;
