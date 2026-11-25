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
 */

import * as native from 'libentry.so';
import { hilog } from '@kit.PerformanceAnalysisKit';
import { BusinessError } from '@kit.BasicServicesKit';
import { promptAction } from '@kit.ArkUI';
import { pasteboard } from '@kit.BasicServicesKit';
import { display } from '@kit.ArkUI';
import { picker } from '@kit.CoreFileKit';
import { fileIo } from '@kit.CoreFileKit';
import { webview } from '@kit.ArkWeb';

import { getContext } from '../common/WailsContext';

const DOMAIN: number = 0x0000;
const TAG: string = 'WailsBridge';

/**
 * The dispatcher registered with the Go bridge.
 *
 * The Go side owns the TSFN and calls this function on the ArkTS thread with
 * (callId, name, argsJSON).
 */
type Dispatcher = (id: number, name: string, payload: string) => void;

/** The page that currently owns the ArkWeb view. */
let pageDispatcher: Dispatcher | null = null;

/** The ArkWeb controller, used to run JavaScript in the page. */
let webController: webview.WebviewController | null = null;

export function setPageDispatcher(dispatcher: Dispatcher | null): void {
  pageDispatcher = dispatcher;
}

export function setWebController(controller: webview.WebviewController | null): void {
  webController = controller;
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
// Helpers
// ---------------------------------------------------------------------------

interface BridgeArgs {
  id?: number;
  script?: string;
  text?: string;
  url?: string;
  value?: number;
  level?: string;
  message?: string;
  options?: string;
  orientation?: string;
  duration?: number;
}

/** The dark mode answer. */
interface DarkModeResult {
  dark: boolean;
}

/** The device information answer. */
interface DeviceInfoResult {
  platform: string;
  manufacturer: string;
  brand: string;
  model: string;
  device: string;
  version: string;
  sdkInt: number;
}

/** Screen information in the shape expected by pkg/application/screen_openharmony.go. */
interface ScreenInfoResult {
  widthPx?: number;
  heightPx?: number;
  density?: number;
  insetTop?: number;
  insetBottom?: number;
  insetLeft?: number;
  insetRight?: number;
}

/** The clipboard answer. */
interface ClipboardResult {
  text: string;
}


function parseArgs(payload: string): BridgeArgs {
  if (payload === '') {
    return {};
  }
  try {
    return JSON.parse(payload) as BridgeArgs;
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'unable to parse bridge arguments');
    return {};
  }
}

/**
 * Run JavaScript in the page and report the result.
 *
 * The script is executed against the current document, so the Go side
 * includes its own session guard when the reply matters.
 */
function executeJavaScript(callId: number, script: string): void {
  const controller = webController;
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

/** Reply with a JSON value. */
function replyJSON(callId: number, value: Object): void {
  wailsBridgeResult(callId, JSON.stringify(value));
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
  const args = parseArgs(payload);

  switch (name) {
    // --- window -------------------------------------------------------
    case 'show':
    case 'hide':
      wailsBridgeResult(callId, '');
      return;
    case 'isDarkMode':
      replyJSON(callId, { dark: currentColorMode() } as DarkModeResult);
      return;

    // --- scripting ----------------------------------------------------
    case 'executeJavaScript':
      executeJavaScript(callId, args.script ?? '');
      return;
    case 'runOnMainThread':
      native.runOnMainThread(args.id ?? 0);
      wailsBridgeResult(callId, '');
      return;

    // --- logging ------------------------------------------------------
    case 'log':
      hilog.info(DOMAIN, TAG, '%{public}s', `[${args.level ?? 'INFO'}] ${args.message ?? ''}`);
      wailsBridgeResult(callId, '');
      return;

    // --- system information -------------------------------------------
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

    // --- clipboard ----------------------------------------------------
    case 'setClipboardText':
      setClipboardText(args.text ?? '');
      wailsBridgeResult(callId, '');
      return;
    case 'getClipboardText':
      getClipboardText(callId);
      return;

    // --- screen -------------------------------------------------------
    case 'getScreenInfoJson':
      replyJSON(callId, screenInfo());
      return;

    // --- dialogs ------------------------------------------------------
    case 'showMessageDialog':
      showMessageDialog(args);
      return;
    case 'showFilePicker':
      showFilePicker(args);
      return;

    // --- toasts and sharing -------------------------------------------
    case 'showToast':
      showToast(args.message ?? '');
      wailsBridgeResult(callId, '');
      return;
    case 'share':
      share(payload);
      wailsBridgeResult(callId, '');
      return;
    case 'openURL':
      openURL(args.url ?? '');
      wailsBridgeResult(callId, '');
      return;

    // --- device features ----------------------------------------------
    case 'setKeepAwake':
    case 'setTorch':
    case 'setBrightness':
    case 'setOrientation':
    case 'setStatusBar':
    case 'postNotification':
    case 'haptic':
    case 'authenticate':
    case 'setMotion':
    case 'setProximity':
    case 'setKeyboardWatch':
    case 'setScreenProtect':
    case 'speak':
    case 'stopSpeak':
    case 'showAboutDialog':
      // These need device capabilities that the template does not claim.
      // They are acknowledged so the Go side does not block on them.
      hilog.debug(DOMAIN, TAG, '%{public}s is not implemented by the OpenHarmony shell', name);
      wailsBridgeResult(callId, '');
      return;

    default:
      hilog.warn(DOMAIN, TAG, 'unknown bridge command: %{public}s', name);
      wailsBridgeResult(callId, '');
      return;
  }
}

// ---------------------------------------------------------------------------
// Command implementations
// ---------------------------------------------------------------------------

function currentColorMode(): boolean {
  // ConfigurationConstant.ColorMode.COLOR_MODE_DARK === 0.
  const context: object = getContext();
  if ('config' in context && context.config && typeof context.config === 'object' && 'colorMode' in context.config) {
    return context.config.colorMode === 0;
  }
  return false;
}

function screenInfo(): ScreenInfoResult {
  try {
    const current = display.getDefaultDisplaySync();
    return {
      widthPx: current.width,
      heightPx: current.height,
      density: current.densityPixels,
      insetTop: 0,
      insetBottom: 0,
      insetLeft: 0,
      insetRight: 0
    };
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'screen info unavailable');
    return {};
  }
}

function setClipboardText(text: string): void {
  try {
    const board = pasteboard.getSystemPasteboard();
    const data = pasteboard.createData(pasteboard.MIMETYPE_TEXT_PLAIN, text);
    board.setDataSync(data);
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'setClipboardText failed: %{public}s', JSON.stringify(e));
  }
}

function getClipboardText(callId: number): void {
  try {
    const board = pasteboard.getSystemPasteboard();
    const data = board.getDataSync();
    wailsBridgeResult(callId, JSON.stringify({ text: data.getPrimaryText() } as ClipboardResult));
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'getClipboardText failed: %{public}s', JSON.stringify(e));
    wailsBridgeResult(callId, '');
  }
}

function showToast(message: string): void {
  try {
    promptAction.showToast({ message: message });
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'showToast failed: %{public}s', JSON.stringify(e));
  }
}

/**
 * Show a message dialog and report the pressed button index.
 *
 * The button index is the index into the options array, or -1 when the dialog
 * was dismissed without choosing.
 */
function showMessageDialog(args: BridgeArgs): void {
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

  const callbackID = args.id ?? 0;
  let options: DialogOptions = { title: '', message: '', buttons: [] };
  try {
    options = JSON.parse(args.options ?? '{}') as DialogOptions;
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'showMessageDialog: invalid options ignored');
  }

  const buttons = (options.buttons ?? []).slice(0, 3);
  if (buttons.length === 0) {
    native.dialogResult(callbackID, -1);
    return;
  }

  try {
    promptAction.showDialog({
      title: options.title,
      message: options.message,
      buttons: buttons.map((button: DialogButton) => ({ text: button.label, color: '#000000' }))
    }).then((result: promptAction.ShowDialogSuccessResponse): void => {
      const index = result.index;
      native.dialogResult(callbackID, index >= 0 && index < buttons.length ? index : -1);
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'showMessageDialog failed: %{public}s', error.message);
      native.dialogResult(callbackID, -1);
    });
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'showMessageDialog failed: %{public}s', JSON.stringify(e));
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
  interface PickerOptions {
    multiple: boolean;
  }
  const callbackID = args.id ?? 0;
  let options: PickerOptions = { multiple: false };
  try {
    options = JSON.parse(args.options ?? '{}') as PickerOptions;
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'showFilePicker: invalid options ignored');
  }

  try {
    const documentViewPicker = new picker.DocumentViewPicker();
    const pickerOptions: picker.DocumentSelectOptions = {
      maxSelectNumber: options.multiple ? 100 : 1,
      fileSuffixFilters: []
    };
    documentViewPicker.select(pickerOptions).then((result: string[]): void => {
      if (result.length === 0) {
        native.filePickerDone(callbackID);
        return;
      }
      const uris = options.multiple ? result : [result[0]];
      copyUrisToCache(callbackID, uris);
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'showFilePicker failed: %{public}s', error.message);
      native.filePickerDone(callbackID);
    });
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'showFilePicker failed: %{public}s', JSON.stringify(e));
    native.filePickerDone(callbackID);
  }
}

function copyUrisToCache(callbackID: number, uris: string[]): void {
  const cacheDir = getContext().cacheDir;
  const paths: string[] = [];
  for (const uri of uris) {
    try {
      const source = fileIo.openSync(uri, fileIo.OpenMode.READ_ONLY);
      const target = `${cacheDir}/wails_pick_${Date.now()}_${paths.length}`;
      const destination = fileIo.openSync(target, fileIo.OpenMode.CREATE | fileIo.OpenMode.READ_WRITE);
      fileIo.copyFileSync(source.fd, destination.fd);
      fileIo.closeSync(source.fd);
      fileIo.closeSync(destination.fd);
      paths.push(target);
    } catch (e) {
      hilog.error(DOMAIN, TAG, 'copyUrisToCache failed: %{public}s', JSON.stringify(e));
    }
  }
  for (const path of paths) {
    native.filePickerResult(callbackID, path);
  }
  native.filePickerDone(callbackID);
}

function share(payload: string): void {
  // The template does not claim the share capability; the payload is logged
  // so an application can see that the call reached the shell.
  hilog.info(DOMAIN, TAG, 'share requested: %{public}s', payload);
}

function openURL(url: string): void {
  try {
    const context = getContext();
    context.startAbility({
      action: 'ohos.want.action.viewData',
      uri: url
    }).catch((error: BusinessError): void => {
      hilog.error(DOMAIN, TAG, 'openURL failed: %{public}s', error.message);
    });
  } catch (e) {
    hilog.error(DOMAIN, TAG, 'openURL failed: %{public}s', JSON.stringify(e));
  }
}
