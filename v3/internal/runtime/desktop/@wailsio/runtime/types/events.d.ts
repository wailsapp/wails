export function setup(): void;
/**
 * Register a callback function to be called multiple times for a specific event.
 *
 * @param {string} eventName - The name of the event to register the callback for.
 * @param {function} callback - The callback function to be called when the event is triggered.
 * @param {number} maxCallbacks - The maximum number of times the callback can be called for the event. Once the maximum number is reached, the callback will no longer be called.
 *
 @return {function} - A function that, when called, will unregister the callback from the event.
 */
export function OnMultiple(eventName: string, callback: Function, maxCallbacks: number): Function;
/**
 * Registers a callback function to be executed when the specified event occurs.
 *
 * @param {string} eventName - The name of the event.
 * @param {function} callback - The callback function to be executed. It takes no parameters.
 * @return {function} - A function that, when called, will unregister the callback from the event. */
export function On(eventName: string, callback: Function): Function;
/**
 * Registers a callback function to be executed only once for the specified event.
 *
 * @param {string} eventName - The name of the event.
 * @param {function} callback - The function to be executed when the event occurs.
 * @return {function} - A function that, when called, will unregister the callback from the event.
 */
export function Once(eventName: string, callback: Function): Function;
/**
 * Removes event listeners for the specified event names.
 *
 * @param {string} eventName - The name of the event to remove listeners for.
 * @param {...string} additionalEventNames - Additional event names to remove listeners for.
 * @return {undefined}
 */
export function Off(eventName: string, ...additionalEventNames: string[]): undefined;
/**
 * Removes all event listeners.
 *
 * @function OffAll
 * @returns {void}
 */
export function OffAll(): void;
/**
 * Emits an event using the given event name.
 *
 * @param {WailsEvent} event - The name of the event to emit.
 * @returns {any} - The result of the emitted event.
 */
export function Emit(event: WailsEvent): any;
export const Types: {
    Windows: {
        SystemThemeChanged: string;
        APMPowerStatusChange: string;
        APMSuspend: string;
        APMResumeAutomatic: string;
        APMResumeSuspend: string;
        APMPowerSettingChange: string;
        ApplicationStarted: string;
        WebViewNavigationCompleted: string;
        WindowInactive: string;
        WindowActive: string;
        WindowClickActive: string;
        WindowMaximise: string;
        WindowUnMaximise: string;
        WindowFullscreen: string;
        WindowUnFullscreen: string;
        WindowRestore: string;
        WindowMinimise: string;
        WindowUnMinimise: string;
        WindowClose: string;
        WindowSetFocus: string;
        WindowKillFocus: string;
        WindowDragDrop: string;
        WindowDragEnter: string;
        WindowDragLeave: string;
        WindowDragOver: string;
    };
    Mac: {
        ApplicationDidBecomeActive: string;
        ApplicationDidChangeBackingProperties: string;
        ApplicationDidChangeEffectiveAppearance: string;
        ApplicationDidChangeIcon: string;
        ApplicationDidChangeOcclusionState: string;
        ApplicationDidChangeScreenParameters: string;
        ApplicationDidChangeStatusBarFrame: string;
        ApplicationDidChangeStatusBarOrientation: string;
        ApplicationDidFinishLaunching: string;
        ApplicationDidHide: string;
        ApplicationDidResignActiveNotification: string;
        ApplicationDidUnhide: string;
        ApplicationDidUpdate: string;
        ApplicationWillBecomeActive: string;
        ApplicationWillFinishLaunching: string;
        ApplicationWillHide: string;
        ApplicationWillResignActive: string;
        ApplicationWillTerminate: string;
        ApplicationWillUnhide: string;
        ApplicationWillUpdate: string;
        ApplicationDidChangeTheme: string;
        ApplicationShouldHandleReopen: string;
        WindowDidBecomeKey: string;
        WindowDidBecomeMain: string;
        WindowDidBeginSheet: string;
        WindowDidChangeAlpha: string;
        WindowDidChangeBackingLocation: string;
        WindowDidChangeBackingProperties: string;
        WindowDidChangeCollectionBehavior: string;
        WindowDidChangeEffectiveAppearance: string;
        WindowDidChangeOcclusionState: string;
        WindowDidChangeOrderingMode: string;
        WindowDidChangeScreen: string;
        WindowDidChangeScreenParameters: string;
        WindowDidChangeScreenProfile: string;
        WindowDidChangeScreenSpace: string;
        WindowDidChangeScreenSpaceProperties: string;
        WindowDidChangeSharingType: string;
        WindowDidChangeSpace: string;
        WindowDidChangeSpaceOrderingMode: string;
        WindowDidChangeTitle: string;
        WindowDidChangeToolbar: string;
        WindowDidChangeVisibility: string;
        WindowDidDeminiaturize: string;
        WindowDidEndSheet: string;
        WindowDidEnterFullScreen: string;
        WindowDidEnterVersionBrowser: string;
        WindowDidExitFullScreen: string;
        WindowDidExitVersionBrowser: string;
        WindowDidExpose: string;
        WindowDidFocus: string;
        WindowDidMiniaturize: string;
        WindowDidMove: string;
        WindowDidOrderOffScreen: string;
        WindowDidOrderOnScreen: string;
        WindowDidResignKey: string;
        WindowDidResignMain: string;
        WindowDidResize: string;
        WindowDidUpdate: string;
        WindowDidUpdateAlpha: string;
        WindowDidUpdateCollectionBehavior: string;
        WindowDidUpdateCollectionProperties: string;
        WindowDidUpdateShadow: string;
        WindowDidUpdateTitle: string;
        WindowDidUpdateToolbar: string;
        WindowDidUpdateVisibility: string;
        WindowShouldClose: string;
        WindowWillBecomeKey: string;
        WindowWillBecomeMain: string;
        WindowWillBeginSheet: string;
        WindowWillChangeOrderingMode: string;
        WindowWillClose: string;
        WindowWillDeminiaturize: string;
        WindowWillEnterFullScreen: string;
        WindowWillEnterVersionBrowser: string;
        WindowWillExitFullScreen: string;
        WindowWillExitVersionBrowser: string;
        WindowWillFocus: string;
        WindowWillMiniaturize: string;
        WindowWillMove: string;
        WindowWillOrderOffScreen: string;
        WindowWillOrderOnScreen: string;
        WindowWillResignMain: string;
        WindowWillResize: string;
        WindowWillUnfocus: string;
        WindowWillUpdate: string;
        WindowWillUpdateAlpha: string;
        WindowWillUpdateCollectionBehavior: string;
        WindowWillUpdateCollectionProperties: string;
        WindowWillUpdateShadow: string;
        WindowWillUpdateTitle: string;
        WindowWillUpdateToolbar: string;
        WindowWillUpdateVisibility: string;
        WindowWillUseStandardFrame: string;
        MenuWillOpen: string;
        MenuDidOpen: string;
        MenuDidClose: string;
        MenuWillSendAction: string;
        MenuDidSendAction: string;
        MenuWillHighlightItem: string;
        MenuDidHighlightItem: string;
        MenuWillDisplayItem: string;
        MenuDidDisplayItem: string;
        MenuWillAddItem: string;
        MenuDidAddItem: string;
        MenuWillRemoveItem: string;
        MenuDidRemoveItem: string;
        MenuWillBeginTracking: string;
        MenuDidBeginTracking: string;
        MenuWillEndTracking: string;
        MenuDidEndTracking: string;
        MenuWillUpdate: string;
        MenuDidUpdate: string;
        MenuWillPopUp: string;
        MenuDidPopUp: string;
        MenuWillSendActionToItem: string;
        MenuDidSendActionToItem: string;
        WebViewDidStartProvisionalNavigation: string;
        WebViewDidReceiveServerRedirectForProvisionalNavigation: string;
        WebViewDidFinishNavigation: string;
        WebViewDidCommitNavigation: string;
        WindowFileDraggingEntered: string;
        WindowFileDraggingPerformed: string;
        WindowFileDraggingExited: string;
    };
    Linux: {
        SystemThemeChanged: string;
        WindowLoadChanged: string;
        WindowDeleteEvent: string;
        WindowFocusIn: string;
        WindowFocusOut: string;
        ApplicationStartup: string;
    };
    Common: {
        ApplicationStarted: string;
        WindowMaximise: string;
        WindowUnMaximise: string;
        WindowFullscreen: string;
        WindowUnFullscreen: string;
        WindowRestore: string;
        WindowMinimise: string;
        WindowUnMinimise: string;
        WindowClosing: string;
        WindowZoom: string;
        WindowZoomIn: string;
        WindowZoomOut: string;
        WindowZoomReset: string;
        WindowFocus: string;
        WindowLostFocus: string;
        WindowShow: string;
        WindowHide: string;
        WindowDPIChanged: string;
        WindowFilesDropped: string;
        WindowRuntimeReady: string;
        ThemeChanged: string;
    };
};
export class WailsEvent {
    constructor(name: any, data?: any);
    name: any;
    data: any;
}
