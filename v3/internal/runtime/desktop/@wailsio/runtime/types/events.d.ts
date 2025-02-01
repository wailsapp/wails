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
        APMPowerSettingChange: string;
        APMPowerStatusChange: string;
        APMResumeAutomatic: string;
        APMResumeSuspend: string;
        APMSuspend: string;
        ApplicationStarted: string;
        SystemThemeChanged: string;
        WebViewNavigationCompleted: string;
        WindowActive: string;
        WindowBackgroundErase: string;
        WindowClickActive: string;
        WindowClosing: string;
        WindowDidMove: string;
        WindowDidResize: string;
        WindowDPIChanged: string;
        WindowDragDrop: string;
        WindowDragEnter: string;
        WindowDragLeave: string;
        WindowDragOver: string;
        WindowEndMove: string;
        WindowEndResize: string;
        WindowFullscreen: string;
        WindowHide: string;
        WindowInactive: string;
        WindowKeyDown: string;
        WindowKeyUp: string;
        WindowKillFocus: string;
        WindowNonClientHit: string;
        WindowNonClientMouseDown: string;
        WindowNonClientMouseLeave: string;
        WindowNonClientMouseMove: string;
        WindowNonClientMouseUp: string;
        WindowPaint: string;
        WindowRestore: string;
        WindowSetFocus: string;
        WindowShow: string;
        WindowStartMove: string;
        WindowStartResize: string;
        WindowUnFullscreen: string;
        WindowZOrderChanged: string;
        WindowMinimise: string;
        WindowUnMinimise: string;
        WindowMaximise: string;
        WindowUnMaximise: string;
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
        ApplicationDidChangeTheme: string;
        ApplicationDidFinishLaunching: string;
        ApplicationDidHide: string;
        ApplicationDidResignActive: string;
        ApplicationDidUnhide: string;
        ApplicationDidUpdate: string;
        ApplicationShouldHandleReopen: string;
        ApplicationWillBecomeActive: string;
        ApplicationWillFinishLaunching: string;
        ApplicationWillHide: string;
        ApplicationWillResignActive: string;
        ApplicationWillTerminate: string;
        ApplicationWillUnhide: string;
        ApplicationWillUpdate: string;
        MenuDidAddItem: string;
        MenuDidBeginTracking: string;
        MenuDidClose: string;
        MenuDidDisplayItem: string;
        MenuDidEndTracking: string;
        MenuDidHighlightItem: string;
        MenuDidOpen: string;
        MenuDidPopUp: string;
        MenuDidRemoveItem: string;
        MenuDidSendAction: string;
        MenuDidSendActionToItem: string;
        MenuDidUpdate: string;
        MenuWillAddItem: string;
        MenuWillBeginTracking: string;
        MenuWillDisplayItem: string;
        MenuWillEndTracking: string;
        MenuWillHighlightItem: string;
        MenuWillOpen: string;
        MenuWillPopUp: string;
        MenuWillRemoveItem: string;
        MenuWillSendAction: string;
        MenuWillSendActionToItem: string;
        MenuWillUpdate: string;
        WebViewDidCommitNavigation: string;
        WebViewDidFinishNavigation: string;
        WebViewDidReceiveServerRedirectForProvisionalNavigation: string;
        WebViewDidStartProvisionalNavigation: string;
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
        WindowDidZoom: string;
        WindowFileDraggingEntered: string;
        WindowFileDraggingExited: string;
        WindowFileDraggingPerformed: string;
        WindowHide: string;
        WindowMaximise: string;
        WindowUnMaximise: string;
        WindowMinimise: string;
        WindowUnMinimise: string;
        WindowShouldClose: string;
        WindowShow: string;
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
        WindowZoomIn: string;
        WindowZoomOut: string;
        WindowZoomReset: string;
    };
    Linux: {
        ApplicationStartup: string;
        SystemThemeChanged: string;
        WindowDeleteEvent: string;
        WindowDidMove: string;
        WindowDidResize: string;
        WindowFocusIn: string;
        WindowFocusOut: string;
        WindowLoadChanged: string;
    };
    Common: {
        ApplicationOpenedWithFile: string;
        ApplicationStarted: string;
        ThemeChanged: string;
        WindowClosing: string;
        WindowDidMove: string;
        WindowDidResize: string;
        WindowDPIChanged: string;
        WindowFilesDropped: string;
        WindowFocus: string;
        WindowFullscreen: string;
        WindowHide: string;
        WindowLostFocus: string;
        WindowMaximise: string;
        WindowMinimise: string;
        WindowRestore: string;
        WindowRuntimeReady: string;
        WindowShow: string;
        WindowUnFullscreen: string;
        WindowUnMaximise: string;
        WindowUnMinimise: string;
        WindowZoom: string;
        WindowZoomIn: string;
        WindowZoomOut: string;
        WindowZoomReset: string;
    };
};
export class WailsEvent {
    constructor(name: any, data?: any);
    name: any;
    data: any;
}
