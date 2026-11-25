/*
 * WailsBridge.ets — OpenHarmony native shell bridge.
 *
 * The Go backend calls named commands over the NAPI threadsafe function, and
 * this dispatcher routes them to the OpenHarmony system APIs.
 *
 * Protocol (see pkg/application/bridge_openharmony.go):
 *   Go -> ArkTS : wailsBridgeDispatch(callId, name, argsJSON)
 *                 callId 0 means "no reply expected"; any other callId must be
 *                 answered exactly once via wailsBridgeResult(callId, result).
 *   ArkTS -> Go : wailsBridgeResult(callId, result)
 *                 result is a JSON string; "" means "not available".
 *
 * Every command below either performs a real system call or replies with an
 * explicit "unsupported:"/"error:" string. Nothing is silently acknowledged.
 */

import * as native from 'libentry.so';
import { hilog } from '@kit.PerformanceAnalysisKit';
import { BusinessError, batteryInfo, commonEventManager } from '@kit.BasicServicesKit';
import { promptAction, display, window } from '@kit.ArkUI';
import { pasteboard } from '@kit.BasicServicesKit';
import { picker, fileIo, statfs } from '@kit.CoreFileKit';
import { webview } from '@kit.ArkWeb';
import { abilityAccessCtrl, wantAgent, WantAgent, Permissions, PermissionRequestResult } from '@kit.AbilityKit';
import { geoLocationManager } from '@kit.LocationKit';
import { camera, cameraPicker } from '@kit.CameraKit';
import { asset } from '@kit.AssetStoreKit';
import { util } from '@kit.ArkTS';
import { connection } from '@kit.NetworkKit';
import { backgroundTaskManager } from '@kit.BackgroundTasksKit';

import { getContext, getMainWindow } from '../common/WailsContext';

const DOMAIN: number = 0x0000;
const TAG: string = 'WailsBridge';

/** The only origin the shell is allowed to navigate or load content for. */
const ORIGIN: string = 'https://wails.localhost/';

/**
 * The dispatcher registered with the Go bridge.
 *
 * The Go side owns the TSFN and calls this function on the ArkTS thread with
 * (callId, name, argsJSON).
 */
type Dispatcher = (id: number, name: string, payload: string) => void;

/** The page that currently owns the ArkWeb view. */
let pageDispatcher: Dispatcher | null = null;

/** The ArkWeb controller, used for navigation, zoom and script execution. */
let webController: webview.WebviewController | null = null;

/** The URL the controlled ArkWeb view is currently displaying. */
let currentURL: string = ORIGIN;

/** ArkWeb has no getZoom(); keep the factor we last applied. */
let currentZoom: number = 1.0;

export function setPageDispatcher(dispatcher: Dispatcher | null): void {
  pageDispatcher = dispatcher;
}

export function setWebController(controller: webview.WebviewController | null): void {
  webController = controller;
}

/** Called by the shell page whenever the controlled origin changes. */
export function setCurrentURL(url: string): void {
  currentURL = url;
}

/**
 * Reply to a Go bridge call.
 *
 * A call whose reply is a JSON string is delivered as-is; failures are
 * reported as an empty string so the Go side can treat them as "not
 * available".
 */
export function wailsBridgeResult(callId: number, result: string): void {
  if (callId === 0) {
    return;
  }
  native.result(callId, result);
}

// ---------------------------------------------------------------------------
// Argument and result shapes
// ---------------------------------------------------------------------------

interface BridgeArgs {
  id?: number;
  script?: string;
  text?: string;
  name?: string;
  description?: string;
  url?: string;
  title?: string;
  html?: string;
  icon?: string;
  key?: string;
  value?: number | boolean;
  width?: number;
  height?: number;
  x?: number;
  y?: number;
  zoom?: number;
  r?: number;
  g?: number;
  b?: number;
  a?: number;
  level?: string;
  message?: string;
  options?: string;
  orientation?: string;
  duration?: number;
}

interface DarkModeResult {
  dark: boolean;
}

interface DeviceInfoResult {
  platform: string;
  manufacturer: string;
  brand: string;
  model: string;
  device: string;
  version: string;
  sdkInt: number;
}

interface ScreenInfoResult {
  widthPx?: number;
  heightPx?: number;
  density?: number;
  insetTop?: number;
  insetBottom?: number;
  insetLeft?: number;
  insetRight?: number;
}

interface ClipboardResult {
  text: string;
}

interface ZoomResult {
  zoom: number;
}

interface ColorResult {
  color: string;
}

interface ErrorEnvelope {
  error: string;
}

interface SavePickerOptions {
  newFileNames?: string[];
  fileSuffixChoices?: string[];
}

interface DialogButton {
  label: string;
  isCancel: boolean;
  isDefault: boolean;
}

interface DialogOptions {
  title: string;
  message: string;
  buttons: DialogButton[];
}

interface PickerOptions {
  multiple: boolean;
}

interface LocationResult {
  lat: number;
  lng: number;
  accuracy: number;
}

interface CameraResult {
  uri: string;
  mediaType: string;
}

function parseArgs(payload: string): BridgeArgs {
  if (payload === '') {
    return {};
  }
  try {
    return JSON.parse(payload) as BridgeArgs;
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'unable to parse bridge arguments');
    return {};
  }
}

function flag(value: number | boolean | undefined): boolean {
  if (typeof value === 'boolean') {
    return value;
  }
  return (value ?? 0) !== 0;
}

function numberArg(value: number | undefined, fallback: number): number {
  return typeof value === 'number' && Number.isFinite(value) ? value : fallback;
}

function clampByte(value: number): number {
  return Math.max(0, Math.min(255, Math.round(value)));
}

function hexByte(value: number): string {
  return clampByte(value).toString(16).padStart(2, '0');
}

function byteArray(value: string): Uint8Array {
  return new util.TextEncoder().encodeInto(value);
}

/** Reply with a JSON value. */
function replyJSON(callId: number, value: Object): void {
  wailsBridgeResult(callId, JSON.stringify(value));
}

/** Explicit, non-silent failure replies. */
function unsupported(callId: number, name: string, detail: string): void {
  const message: string = `unsupported: ${name} ${detail}`;
  hilog.info(DOMAIN, TAG, '%{public}s', message);
  wailsBridgeResult(callId, message);
}

function failed(callId: number, name: string, detail: string): void {
  const message: string = `error: ${name}: ${detail}`;
  hilog.error(DOMAIN, TAG, '%{public}s', message);
  wailsBridgeResult(callId, message);
}

// ---------------------------------------------------------------------------
// Window helpers
// ---------------------------------------------------------------------------

/** Run an asynchronous window operation and answer exactly once. */
function windowPromise(callId: number, name: string, op: (win: window.Window) => Promise<void>): void {
  const win: window.Window | null = getMainWindow();
  if (win === null) {
    failed(callId, name, 'main window is not available');
    return;
  }
  try {
    op(win).then((): void => {
      wailsBridgeResult(callId, '');
    }).catch((error: BusinessError): void => {
      failed(callId, name, error.message);
    });
  } catch (error) {
    failed(callId, name, JSON.stringify(error));
  }
}

/** Run a synchronous window operation and answer exactly once. */
function windowSync(callId: number, name: string, op: (win: window.Window) => void): void {
  const win: window.Window | null = getMainWindow();
  if (win === null) {
    failed(callId, name, 'main window is not available');
    return;
  }
  try {
    op(win);
    wailsBridgeResult(callId, '');
  } catch (error) {
    failed(callId, name, JSON.stringify(error));
  }
}

/** Minimize is the nearest system operation; true hide is not available. */
function hideWindow(callId: number): void {
  const win: window.Window | null = getMainWindow();
  if (win === null) {
    failed(callId, 'hide', 'main window is not available');
    return;
  }
  try {
    win.minimize().then((): void => {
      unsupported(callId, 'hide', 'true hide is not available; the window was minimized');
    }).catch((error: BusinessError): void => {
      failed(callId, 'hide', error.message);
    });
  } catch (error) {
    failed(callId, 'hide', JSON.stringify(error));
  }
}

function orientationOf(value: string): window.Orientation {
  switch (value) {
    case 'portrait':
      return window.Orientation.PORTRAIT;
    case 'landscape':
      return window.Orientation.LANDSCAPE;
    case 'auto':
    case '':
      return window.Orientation.AUTO_ROTATION;
    default:
      return window.Orientation.UNSPECIFIED;
  }
}

// ---------------------------------------------------------------------------
// Web helpers — the local-origin policy is enforced here.
// ---------------------------------------------------------------------------

function isLocalURL(url: string): boolean {
  return url.startsWith(ORIGIN);
}

function webSync(callId: number, name: string, op: (controller: webview.WebviewController) => void): void {
  const controller: webview.WebviewController | null = webController;
  if (controller === null) {
    failed(callId, name, 'WebviewController is not attached');
    return;
  }
  if (!isLocalURL(currentURL)) {
    failed(callId, name, 'only the controlled wails.localhost origin is available');
    return;
  }
  try {
    op(controller);
    wailsBridgeResult(callId, '');
  } catch (error) {
    failed(callId, name, JSON.stringify(error));
  }
}

/** Editing commands are delivered to the controlled document via ArkWeb. */
function webEdit(callId: number, name: string, command: string): void {
  webSync(callId, name, (controller: webview.WebviewController): void => {
    controller.runJavaScript(`document.execCommand(${JSON.stringify(command)})`);
  });
}

/** ArkWeb has no zoomReset(); apply the default factor and say so. */
function resetZoom(callId: number): void {
  const controller: webview.WebviewController | null = webController;
  if (controller === null) {
    failed(callId, 'zoomReset', 'WebviewController is not attached');
    return;
  }
  if (!isLocalURL(currentURL)) {
    failed(callId, 'zoomReset', 'only the controlled wails.localhost origin is available');
    return;
  }
  try {
    controller.zoom(1.0);
    currentZoom = 1.0;
    unsupported(callId, 'zoomReset', 'ArkWeb has no zoomReset; zoom was set to 1.0');
  } catch (error) {
    failed(callId, 'zoomReset', JSON.stringify(error));
  }
}

function executeJavaScript(callId: number, script: string): void {
  const controller: webview.WebviewController | null = webController;
  if (controller === null) {
    wailsBridgeResult(callId, '');
    return;
  }
  controller.runJavaScript(script).then((value: string): void => {
    wailsBridgeResult(callId, value);
  }).catch((error: BusinessError): void => {
    hilog.error(DOMAIN, TAG, 'JavaScript error: %{public}s', error.message);
    wailsBridgeResult(callId, '');
  });
}

// ---------------------------------------------------------------------------
// Runtime permission helper
// ---------------------------------------------------------------------------

async function hasPermission(permission: string): Promise<boolean> {
  const manager: abilityAccessCtrl.AtManager = abilityAccessCtrl.createAtManager();
  try {
    const permissions: Permissions[] = [permission as Permissions];
    const result: PermissionRequestResult =
      await manager.requestPermissionsFromUser(getContext(), permissions);
    return result.authResults.length > 0 &&
      result.authResults[0] === abilityAccessCtrl.GrantStatus.PERMISSION_GRANTED;
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'permission request failed: %{public}s', JSON.stringify(error));
    return false;
  }
}

// ---------------------------------------------------------------------------
// Dispatcher
// ---------------------------------------------------------------------------

/**
 * Handle one Go -> ArkTS command.
 *
 * @param callId the reply id; 0 means no reply is expected
 * @param name the command name
 * @param payload the command arguments as a JSON string
 */
export function dispatchBridge(callId: number, name: string, payload: string): void {
  const args: BridgeArgs = parseArgs(payload);

  switch (name) {
    // --- window -------------------------------------------------------
    case 'setTitle':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.setWindowTitle(args.title ?? ''));
      return;
    case 'setSize':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.resize(Math.round(numberArg(args.width, 0)), Math.round(numberArg(args.height, 0))));
      return;
    case 'setPosition':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.moveWindowTo(Math.round(numberArg(args.x, 0)), Math.round(numberArg(args.y, 0))));
      return;
    case 'setBounds': {
      const x: number = Math.round(numberArg(args.x, 0));
      const y: number = Math.round(numberArg(args.y, 0));
      const width: number = Math.round(numberArg(args.width, 0));
      const height: number = Math.round(numberArg(args.height, 0));
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.moveWindowTo(x, y).then((): Promise<void> => win.resize(width, height)));
      return;
    }
    case 'center': {
      const screen: display.Display = display.getDefaultDisplaySync();
      const width: number = Math.round(numberArg(args.width, 0));
      const height: number = Math.round(numberArg(args.height, 0));
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.moveWindowTo(Math.max(0, Math.floor((screen.width - width) / 2)),
          Math.max(0, Math.floor((screen.height - height) / 2))));
      return;
    }
    case 'fullscreen':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.setWindowLayoutFullScreen(true));
      return;
    case 'unfullscreen':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.setWindowLayoutFullScreen(false));
      return;
    case 'minimise':
      windowPromise(callId, name, (win: window.Window): Promise<void> => win.minimize());
      return;
    case 'unminimise':
      windowPromise(callId, name, (win: window.Window): Promise<void> => win.showWindow());
      return;
    case 'maximise':
      windowPromise(callId, name, (win: window.Window): Promise<void> => win.maximize());
      return;
    case 'unmaximise':
      windowPromise(callId, name, (win: window.Window): Promise<void> => win.recover());
      return;
    case 'close': {
      try {
        getContext().terminateSelf().then((): void => {
          wailsBridgeResult(callId, '');
        }).catch((error: BusinessError): void => {
          failed(callId, name, error.message);
        });
      } catch (error) {
        failed(callId, name, JSON.stringify(error));
      }
      return;
    }
    case 'focus':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.showWindow().then((): Promise<void> => win.setWindowFocusable(true)));
      return;
    case 'show':
      windowPromise(callId, name, (win: window.Window): Promise<void> => win.showWindow());
      return;
    case 'hide':
      hideWindow(callId);
      return;
    case 'setResizable':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.setResizeByDragEnabled(flag(args.value)));
      return;
    case 'setAlwaysOnTop':
      unsupported(callId, name, 'the OpenHarmony window API does not expose always-on-top');
      return;
    case 'setBackgroundColour': {
      const color: string = '#' + hexByte(args.a ?? 255) + hexByte(args.r ?? 255) +
        hexByte(args.g ?? 255) + hexByte(args.b ?? 255);
      windowSync(callId, name, (win: window.Window): void => win.setWindowBackgroundColor(color));
      return;
    }
    case 'setIcon':
      unsupported(callId, name, 'the OpenHarmony window API does not expose an application icon setter');
      return;

    // --- zoom and navigation -----------------------------------------
    case 'setZoom': {
      const zoom: number = numberArg(args.zoom, 1.0);
      webSync(callId, name, (controller: webview.WebviewController): void => {
        controller.zoom(zoom);
        currentZoom = zoom;
      });
      return;
    }
    case 'getZoom':
      replyJSON(callId, { zoom: currentZoom } as ZoomResult);
      return;
    case 'zoomIn':
      webSync(callId, name, (controller: webview.WebviewController): void => {
        controller.zoomIn();
        currentZoom += 0.1;
      });
      return;
    case 'zoomOut':
      webSync(callId, name, (controller: webview.WebviewController): void => {
        controller.zoomOut();
        currentZoom = Math.max(0.1, currentZoom - 0.1);
      });
      return;
    case 'zoomReset':
      resetZoom(callId);
      return;
    case 'setURL': {
      const url: string = args.url ?? '';
      if (!isLocalURL(url)) {
        failed(callId, name, `external URLs are disabled; only ${ORIGIN} is allowed`);
        return;
      }
      webSync(callId, name, (controller: webview.WebviewController): void => {
        controller.loadUrl(url);
        currentURL = url;
      });
      return;
    }
    case 'reload':
      webSync(callId, name, (controller: webview.WebviewController): void => controller.refresh());
      return;
    case 'forceReload':
      webSync(callId, name, (controller: webview.WebviewController): void => controller.refresh(true));
      return;
    case 'setHTML': {
      const html: string = args.html ?? '';
      webSync(callId, name, (controller: webview.WebviewController): void => {
        controller.loadData(html, 'text/html', 'UTF-8');
      });
      return;
    }

    // --- scripting ----------------------------------------------------
    case 'executeJavaScript':
      executeJavaScript(callId, args.script ?? '');
      return;
    case 'runOnMainThread':
      native.runOnMainThread(args.id ?? 0);
      wailsBridgeResult(callId, '');
      return;
    case 'copy':
      webEdit(callId, name, 'copy');
      return;
    case 'cut':
      webEdit(callId, name, 'cut');
      return;
    case 'paste':
      webEdit(callId, name, 'paste');
      return;
    case 'selectAll':
      webEdit(callId, name, 'selectAll');
      return;
    case 'undo':
      webEdit(callId, name, 'undo');
      return;
    case 'redo':
      webEdit(callId, name, 'redo');
      return;
    case 'delete':
      webEdit(callId, name, 'delete');
      return;

    // --- logging and information --------------------------------------
    case 'log':
      hilog.info(DOMAIN, TAG, '%{public}s', `[${args.level ?? 'INFO'}] ${args.message ?? ''}`);
      wailsBridgeResult(callId, '');
      return;
    case 'isDarkMode':
      replyJSON(callId, { dark: currentColorMode() } as DarkModeResult);
      return;
    case 'getAccentColor':
      replyJSON(callId, { color: accentColor() } as ColorResult);
      return;
    case 'getDeviceInfoJson':
      replyJSON(callId, {
        platform: 'openharmony',
        manufacturer: 'OpenHarmony',
        brand: 'OpenHarmony',
        model: 'OpenHarmony',
        device: 'OpenHarmony',
        version: 'OpenHarmony',
        sdkInt: 26
      } as DeviceInfoResult);
      return;
    case 'getScreenInfoJson':
      replyJSON(callId, screenInfo());
      return;
    case 'getAppInfoJson':
      replyJSON(callId, appInfo());
      return;
    case 'getSafeAreaJson':
      replyJSON(callId, safeArea());
      return;
    case 'getBrightnessJson':
      replyJSON(callId, brightnessInfo());
      return;
    case 'getOrientationJson':
      replyJSON(callId, orientationInfo());
      return;
    case 'getStorageJson':
      replyJSON(callId, storageInfo());
      return;
    case 'getStoragePath':
      wailsBridgeResult(callId, getContext().filesDir);
      return;
    case 'getPowerJson':
      replyJSON(callId, powerInfo());
      return;
    case 'getNetworkJson':
      replyJSON(callId, networkInfo());
      return;

    // --- clipboard ----------------------------------------------------
    case 'setClipboardText':
      setClipboardText(args.text ?? '');
      wailsBridgeResult(callId, '');
      return;
    case 'getClipboardText':
      getClipboardText(callId);
      return;

    // --- dialogs and pickers ------------------------------------------
    case 'showMessageDialog':
      showMessageDialog(args);
      return;
    case 'showFilePicker':
      showFilePicker(args);
      return;
    case 'saveFilePicker':
      saveFilePicker(callId, args);
      return;
    case 'showAboutDialog':
      showAboutDialog(args);
      return;
    case 'showToast':
      showToast(args.message ?? '');
      wailsBridgeResult(callId, '');
      return;

    // --- one-way mobile capabilities ----------------------------------
    case 'setKeepAwake':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.setWindowKeepScreenOn(flag(args.value)));
      return;
    case 'setBrightness': {
      const brightness: number = Math.max(0, Math.min(1,
        (typeof args.value === 'number' ? args.value : 0) / 255));
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.setWindowBrightness(brightness));
      return;
    }
    case 'setOrientation':
      windowPromise(callId, name, (win: window.Window): Promise<void> =>
        win.setPreferredOrientation(orientationOf(args.orientation ?? '')));
      return;
    case 'openURL':
      openURL(args.url ?? '');
      wailsBridgeResult(callId, '');
      return;
    case 'getLocation':
      getLocation(callId);
      return;
    case 'capturePhoto':
      capture(callId, name, cameraPicker.PickerMediaType.PHOTO);
      return;
    case 'captureVideo':
      capture(callId, name, cameraPicker.PickerMediaType.VIDEO);
      return;
    case 'startForegroundService':
      startForegroundService(callId);
      return;
    case 'stopForegroundService':
      stopForegroundService(callId);
      return;
    case 'secureSet':
      secureSet(callId, args.key ?? '', String(args.value ?? ''));
      return;
    case 'secureGet':
      secureGet(callId, args.key ?? '');
      return;
    case 'secureDelete':
      secureDelete(callId, args.key ?? '');
      return;
    case 'share':
    case 'setStatusBar':
    case 'postNotification':
    case 'haptic':
    case 'authenticate':
    case 'setMotion':
    case 'setProximity':
    case 'setKeyboardWatch':
    case 'setScreenProtect':
    case 'setTorch':
    case 'speak':
    case 'stopSpeak':
      unsupported(callId, name, 'is not exposed by the OpenHarmony shell');
      return;

    default:
      failed(callId, name, 'unknown bridge command');
      return;
  }
}

// ---------------------------------------------------------------------------
// Command implementations
// ---------------------------------------------------------------------------

function currentColorMode(): boolean {
  interface ThemeConfig {
    colorMode?: number;
  }
  interface ThemeContext {
    config?: ThemeConfig;
  }
  const themed: ThemeContext = getContext() as ThemeContext;
  return themed.config?.colorMode === 0;
}

function accentColor(): string {
  try {
    const value: number = getContext().resourceManager.getColorByNameSync('ohos_id_color_emphasize');
    return '#' + (value & 0x00FFFFFF).toString(16).padStart(6, '0');
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'accent color unavailable: %{public}s', JSON.stringify(error));
    return '#0A59F7';
  }
}

function screenInfo(): ScreenInfoResult {
  try {
    const current: display.Display = display.getDefaultDisplaySync();
    return {
      widthPx: current.width,
      heightPx: current.height,
      density: current.densityPixels,
      insetTop: 0,
      insetBottom: 0,
      insetLeft: 0,
      insetRight: 0
    };
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'screen info unavailable: %{public}s', JSON.stringify(error));
    return {};
  }
}

interface AppInfoResult {
  name: string;
  process: string;
  bundle: string;
}

function appInfo(): AppInfoResult {
  try {
    const info = getContext().applicationInfo;
    const name: string = info?.name ?? 'unknown';
    const process: string = info?.process ?? 'unknown';
    return { name, process, bundle: name } as AppInfoResult;
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'app info unavailable: %{public}s', JSON.stringify(error));
    return { name: 'unknown', process: 'unknown', bundle: 'unknown' } as AppInfoResult;
  }
}

function safeArea(): ScreenInfoResult {
  return screenInfo();
}

interface BrightnessResult {
  brightness: number;
}

function brightnessInfo(): BrightnessResult {
  const win: window.Window | null = getMainWindow();
  if (win === null) {
    return { brightness: -1 };
  }
  return { brightness: win.getWindowProperties().brightness };
}

interface OrientationResult {
  orientation: string;
}

function orientationInfo(): OrientationResult {
  const rotation: number = display.getDefaultDisplaySync().rotation;
  if (rotation === 90 || rotation === 270) {
    return { orientation: 'landscape' };
  }
  return { orientation: 'portrait' };
}

interface StorageResult {
  total: number;
  free: number;
  path: string;
}

function storageInfo(): StorageResult {
  const path: string = getContext().filesDir;
  return {
    total: statfs.getTotalSizeSync(path),
    free: statfs.getFreeSizeSync(path),
    path: path
  };
}

interface PowerResult {
  level: number;
  state: string;
  charging: boolean;
  lowPowerMode: boolean;
}

function powerInfo(): PowerResult {
  const charging: boolean = batteryInfo.chargingStatus === batteryInfo.BatteryChargeState.ENABLE;
  return {
    level: batteryInfo.batterySOC / 100,
    state: charging ? 'charging' : 'discharging',
    charging: charging,
    lowPowerMode: false
  };
}

interface NetworkResult {
  connected: boolean;
  type: string;
  metered: boolean;
}

function networkInfo(): NetworkResult {
  const handle: connection.NetHandle = connection.getDefaultNetSync();
  if (handle.netId === 0) {
    return { connected: false, type: 'none', metered: false };
  }
  const capabilities: connection.NetCapabilities = connection.getNetCapabilitiesSync(handle);
  let type: string = 'unknown';
  const bearers: connection.NetBearType[] = capabilities.bearerTypes;
  for (const bearer of bearers) {
    if (bearer === connection.NetBearType.BEARER_WIFI) {
      type = 'wifi';
      break;
    }
    if (bearer === connection.NetBearType.BEARER_CELLULAR) {
      type = 'cellular';
      break;
    }
    if (bearer === connection.NetBearType.BEARER_ETHERNET) {
      type = 'ethernet';
      break;
    }
  }
  return { connected: true, type: type, metered: type === 'cellular' };
}

function setClipboardText(text: string): void {
  try {
    const board = pasteboard.getSystemPasteboard();
    const data = pasteboard.createData(pasteboard.MIMETYPE_TEXT_PLAIN, text);
    board.setDataSync(data);
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'setClipboardText failed: %{public}s', JSON.stringify(error));
  }
}

function getClipboardText(callId: number): void {
  try {
    const board = pasteboard.getSystemPasteboard();
    const data = board.getDataSync();
    wailsBridgeResult(callId, JSON.stringify({ text: data.getPrimaryText() } as ClipboardResult));
  } catch (error) {
    failed(callId, 'getClipboardText', JSON.stringify(error));
  }
}

function showAboutDialog(args: BridgeArgs): void {
  const callbackID: number = args.id ?? 0;
  const title: string = args.name ?? args.title ?? '';
  const message: string = args.description ?? args.text ?? '';
  promptAction.showDialog({
    title: title,
    message: message,
    buttons: [{ text: 'OK', color: '#000000' } as promptAction.Button]
  }).then((): void => {
    native.dialogResult(callbackID, 0);
  }).catch((error: BusinessError): void => {
    hilog.error(DOMAIN, TAG, 'showAboutDialog failed: %{public}s', error.message);
    native.dialogResult(callbackID, -1);
  });
}

function showMessageDialog(args: BridgeArgs): void {
  const callbackID: number = args.id ?? 0;
  let options: DialogOptions = { title: '', message: '', buttons: [] };
  try {
    options = JSON.parse(args.options ?? '{}') as DialogOptions;
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'showMessageDialog: invalid options ignored');
  }

  const buttons: DialogButton[] = (options.buttons ?? []).slice(0, 3);
  if (buttons.length === 0) {
    native.dialogResult(callbackID, -1);
    return;
  }

  try {
    promptAction.showDialog({
      title: options.title,
      message: options.message,
      buttons: buttons.map((button: DialogButton): promptAction.Button => {
        const styled: promptAction.Button = { text: button.label, color: '#000000' };
        return styled;
      })
    }).then((result: promptAction.ShowDialogSuccessResponse): void => {
      const index: number = result.index;
      native.dialogResult(callbackID, index >= 0 && index < buttons.length ? index : -1);
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'showMessageDialog failed: %{public}s', error.message);
      native.dialogResult(callbackID, -1);
    });
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'showMessageDialog failed: %{public}s', JSON.stringify(error));
    native.dialogResult(callbackID, -1);
  }
}

/**
 * Show the system file picker and report the selected path.
 *
 * OpenHarmony's picker returns content URIs rather than filesystem paths, so
 * the selected URI is copied into the app's cache directory and the copied
 * path is what the Go backend receives.
 */
function showFilePicker(args: BridgeArgs): void {
  const callbackID: number = args.id ?? 0;
  let options: PickerOptions = { multiple: false };
  try {
    options = JSON.parse(args.options ?? '{}') as PickerOptions;
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'showFilePicker: invalid options ignored');
  }

  try {
    const documentViewPicker = new picker.DocumentViewPicker(getContext());
    const pickerOptions: picker.DocumentSelectOptions = {
      maxSelectNumber: options.multiple ? 100 : 1,
      fileSuffixFilters: []
    };
    documentViewPicker.select(pickerOptions).then((result: string[]): void => {
      if (result.length === 0) {
        native.filePickerDone(callbackID);
        return;
      }
      const uris: string[] = options.multiple ? result : [result[0]];
      copyUrisToCache(callbackID, uris);
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'showFilePicker failed: %{public}s', error.message);
      native.filePickerDone(callbackID);
    });
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'showFilePicker failed: %{public}s', JSON.stringify(error));
    native.filePickerDone(callbackID);
  }
}

/** Save dialog: hand the selected content URIs to the Go callbacks. */
function saveFilePicker(callId: number, args: BridgeArgs): void {
  const callbackID: number = args.id ?? callId;
  let options: SavePickerOptions = { newFileNames: [], fileSuffixChoices: [] };
  try {
    options = JSON.parse(args.options ?? '{}') as SavePickerOptions;
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'saveFilePicker: invalid options ignored');
  }
  try {
    const saveOptions: picker.DocumentSaveOptions = new picker.DocumentSaveOptions();
    saveOptions.newFileNames = options.newFileNames ?? [];
    saveOptions.fileSuffixChoices = options.fileSuffixChoices ?? [];
    const documentViewPicker = new picker.DocumentViewPicker(getContext());
    documentViewPicker.save(saveOptions).then((result: string[]): void => {
      for (const uri of result) {
        native.filePickerResult(callbackID, uri);
      }
      native.filePickerDone(callbackID);
      wailsBridgeResult(callId, '');
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'saveFilePicker failed: %{public}s', error.message);
      native.filePickerDone(callbackID);
      failed(callId, 'saveFilePicker', error.message);
    });
  } catch (error) {
    native.filePickerDone(callbackID);
    failed(callId, 'saveFilePicker', JSON.stringify(error));
  }
}

function copyUrisToCache(callbackID: number, uris: string[]): void {
  const cacheDir: string = getContext().cacheDir;
  const paths: string[] = [];
  for (const uri of uris) {
    try {
      const source = fileIo.openSync(uri, fileIo.OpenMode.READ_ONLY);
      const target: string = `${cacheDir}/wails_pick_${Date.now()}_${paths.length}`;
      const destination = fileIo.openSync(target, fileIo.OpenMode.CREATE | fileIo.OpenMode.READ_WRITE);
      fileIo.copyFileSync(source.fd, destination.fd);
      fileIo.closeSync(source.fd);
      fileIo.closeSync(destination.fd);
      paths.push(target);
    } catch (error) {
      hilog.error(DOMAIN, TAG, 'copyUrisToCache failed: %{public}s', JSON.stringify(error));
    }
  }
  for (const path of paths) {
    native.filePickerResult(callbackID, path);
  }
  native.filePickerDone(callbackID);
}

function showToast(message: string): void {
  try {
    promptAction.showToast({ message: message });
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'showToast failed: %{public}s', JSON.stringify(error));
  }
}

function openURL(url: string): void {
  try {
    const abilityContext = getContext();
    abilityContext.startAbility({
      action: 'ohos.want.action.viewData',
      uri: url
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'openURL failed: %{public}s', error.message);
    });
  } catch (error) {
    hilog.error(DOMAIN, TAG, 'openURL failed: %{public}s', JSON.stringify(error));
  }
}

// ---------------------------------------------------------------------------
// Mobile capabilities
// ---------------------------------------------------------------------------

function getLocation(callId: number): void {
  hasPermission('ohos.permission.APPROXIMATELY_LOCATION').then((granted: boolean): void => {
    if (!granted) {
      const denial: ErrorEnvelope = { error: 'location permission denied' };
      native.emitEvent('common:location', JSON.stringify(denial));
      failed(callId, 'getLocation', 'location permission denied');
      return;
    }
    geoLocationManager.getCurrentLocation().then((location: geoLocationManager.Location): void => {
      const result: LocationResult = {
        lat: location.latitude,
        lng: location.longitude,
        accuracy: location.accuracy
      };
      native.emitEvent('common:location', JSON.stringify(result));
      wailsBridgeResult(callId, '');
    }).catch((error: BusinessError): void => {
      const envelope: ErrorEnvelope = { error: error.message };
      native.emitEvent('common:location', JSON.stringify(envelope));
      failed(callId, 'getLocation', error.message);
    });
  });
}

function capture(callId: number, name: string, mediaType: cameraPicker.PickerMediaType): void {
  hasPermission('ohos.permission.CAMERA').then((granted: boolean): void => {
    if (!granted) {
      failed(callId, name, 'camera permission denied');
      return;
    }
    const profile: cameraPicker.PickerProfile = {
      cameraPosition: camera.CameraPosition.CAMERA_POSITION_BACK
    };
    cameraPicker.pick(getContext(), [mediaType], profile).then((result: cameraPicker.PickerResult): void => {
      if (result.resultCode !== 0) {
        failed(callId, name, 'camera picker was cancelled');
        return;
      }
      const event: CameraResult = { uri: result.resultUri, mediaType: result.mediaType };
      native.emitEvent(name === 'capturePhoto' ? 'common:photo' : 'common:video', JSON.stringify(event));
      wailsBridgeResult(callId, '');
    }).catch((error: BusinessError): void => {
      failed(callId, name, error.message);
    });
  });
}

function startForegroundService(callId: number): void {
  const info: wantAgent.WantAgentInfo = {
    wants: [{
      bundleName: getContext().applicationInfo.name,
      abilityName: 'EntryAbility'
    }],
    actionType: wantAgent.OperationType.START_ABILITY,
    requestCode: 0,
    wantAgentFlags: [wantAgent.WantAgentFlags.UPDATE_PRESENT_FLAG]
  };
  wantAgent.getWantAgent(info).then((agent: WantAgent): void => {
    backgroundTaskManager.startBackgroundRunning(getContext(),
      backgroundTaskManager.BackgroundMode.DATA_TRANSFER, agent).then((): void => {
      wailsBridgeResult(callId, '');
    }).catch((error: BusinessError): void => {
      failed(callId, 'startForegroundService', error.message);
    });
  }).catch((error: BusinessError): void => {
    failed(callId, 'startForegroundService', error.message);
  });
}

function stopForegroundService(callId: number): void {
  backgroundTaskManager.stopBackgroundRunning(getContext()).then((): void => {
    wailsBridgeResult(callId, '');
  }).catch((error: BusinessError): void => {
    failed(callId, 'stopForegroundService', error.message);
  });
}

function secureSet(callId: number, key: string, value: string): void {
  if (key === '') {
    failed(callId, 'secureSet', 'key is required');
    return;
  }
  const attributes: asset.AssetMap = new Map();
  attributes.set(asset.Tag.ALIAS, byteArray(key));
  attributes.set(asset.Tag.SECRET, byteArray(value));
  attributes.set(asset.Tag.ACCESSIBILITY, asset.Accessibility.DEVICE_FIRST_UNLOCKED);
  asset.add(attributes).then((): void => {
    wailsBridgeResult(callId, '');
  }).catch((error: BusinessError): void => {
    // An existing alias is updated rather than silently ignored.
    if (error.code === 24000003) {
      const update: asset.AssetMap = new Map();
      update.set(asset.Tag.SECRET, byteArray(value));
      asset.update(attributes, update).then((): void => {
        wailsBridgeResult(callId, '');
      }).catch((updateError: BusinessError): void => {
        failed(callId, 'secureSet', updateError.message);
      });
      return;
    }
    failed(callId, 'secureSet', error.message);
  });
}

function secureGet(callId: number, key: string): void {
  if (key === '') {
    failed(callId, 'secureGet', 'key is required');
    return;
  }
  const query: asset.AssetMap = new Map();
  query.set(asset.Tag.ALIAS, byteArray(key));
  query.set(asset.Tag.RETURN_TYPE, asset.ReturnType.ALL);
  asset.query(query).then((results: asset.AssetMap[]): void => {
    if (results.length === 0) {
      failed(callId, 'secureGet', 'key not found');
      return;
    }
    const secret: asset.Value | undefined = results[0].get(asset.Tag.SECRET);
    if (secret instanceof Uint8Array) {
      const decoder: util.TextDecoder = new util.TextDecoder();
      wailsBridgeResult(callId, decoder.decode(secret));
      return;
    }
    failed(callId, 'secureGet', 'stored value is not a string');
  }).catch((error: BusinessError): void => {
    failed(callId, 'secureGet', error.message);
  });
}

function secureDelete(callId: number, key: string): void {
  if (key === '') {
    failed(callId, 'secureDelete', 'key is required');
    return;
  }
  const query: asset.AssetMap = new Map();
  query.set(asset.Tag.ALIAS, byteArray(key));
  asset.remove(query).then((): void => {
    wailsBridgeResult(callId, '');
  }).catch((error: BusinessError): void => {
    failed(callId, 'secureDelete', error.message);
  });
}
