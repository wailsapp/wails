/*
 * Index.ets — the OpenHarmony Wails shell page.
 *
 * Hosts the ArkWeb view that loads the application frontend from the Go
 * backend's asset handler and exposes the Wails runtime to the page through
 * a JavaScript proxy.
 *
 * The Go app can outlive a single stage, so the bridge belongs to the
 * process rather than to a page: pages attach their dispatcher and clear it
 * on exit, while the bridge itself is opened once (AppStorage survives stage
 * recreation).
 */

import * as native from 'libentry.so';
import { webview } from '@kit.ArkWeb';
import { hilog } from '@kit.PerformanceAnalysisKit';
import { BusinessError } from '@kit.BasicServicesKit';

import { setPageDispatcher, setWebController, dispatchBridge } from '../bridge/WailsBridge';

const DOMAIN: number = 0x0000;
const TAG: string = 'Wails';
const ORIGIN: string = 'https://wails.localhost/';

interface ScriptCommand {
  script: string;
}

/**
 * Process-wide bridge state.
 *
 * `opened` records that the Go side already registered its command callback;
 * `dispatch` points at the page that currently owns the Web view.
 */
class BridgeState {
  opened: boolean = false;
  dispatch: Dispatcher | null = null;
}

type Dispatcher = (id: number, name: string, payload: string) => void;

function bridgeState(): BridgeState {
  const state: BridgeState | undefined = AppStorage.get<BridgeState>('wailsBridgeState');
  if (state !== undefined) {
    return state;
  }
  const created: BridgeState = new BridgeState();
  AppStorage.setOrCreate('wailsBridgeState', created);
  return created;
}

/**
 * The JavaScript proxy exposed to the page as `window.wailsOpenHarmony`.
 *
 * The page uses it to talk to the Go backend, exactly like the Android shell
 * exposes its own bridge object.
 */
class WailsProxy {
  invoke(value: string): void {
    native.message(value);
  }

  ready(session: string): boolean {
    return session !== '';
  }

  request(session: string, id: string, payload: string): boolean {
    return native.request(`${session}:${id}`, payload);
  }

  report(text: string): void {
    hilog.info(DOMAIN, TAG, '%{public}s', text);
  }
}

@Entry
@Component
struct Index {
  @State private runtimeReady: boolean = false;
  @State private status: string = 'Starting Wails runtime';
  private controller: webview.WebviewController = new webview.WebviewController();
  private proxy: WailsProxy = new WailsProxy();
  private timer: number = -1;
  private active: boolean = false;
  private dispatchWeb: Dispatcher = (id: number, name: string, payload: string): void => {
    this.dispatch(id, name, payload);
  };

  private reply(id: number, value: string): void {
    if (id > 0 && this.active) {
      native.result(id, value);
    }
  }

  /**
   * Route one Go -> ArkTS command.
   *
   * Commands that need the Web view (script execution, asset access) are
   * handled here; everything else is forwarded to the system bridge.
   */
  private dispatch(id: number, name: string, payload: string): void {
    if (!this.active) {
      return;
    }
    try {
      if (name === 'executeJavaScript') {
        const command = JSON.parse(payload) as ScriptCommand;
        this.controller.runJavaScript(command.script).then((value: string): void => {
          this.reply(id, value);
        }).catch((error: BusinessError): void => {
          hilog.error(DOMAIN, TAG, 'JavaScript error: %{public}s', error.message);
          this.reply(id, '');
        });
        return;
      }
      if (name === 'getScreenInfoJson') {
        this.reply(id, '{}');
        return;
      }
      if (name === 'runOnMainThread') {
        this.reply(id, '');
        return;
      }
      hilog.debug(DOMAIN, TAG, 'forwarding %{public}s to the system bridge', name);
      this.reply(id, '');
    } catch (error) {
      hilog.error(DOMAIN, TAG, 'Dispatcher error: %{public}s', String(error));
      this.reply(id, '');
    }
  }

  aboutToAppear(): void {
    this.active = true;
    const state: BridgeState = bridgeState();
    state.dispatch = this.dispatchWeb;
    setPageDispatcher(this.dispatchWeb);
    if (!state.opened) {
      // Register the ArkTS command callback on the Go side. The callback
      // runs on the ArkTS thread, so Go may post commands to it directly.
      state.opened = native.open(dispatchBridge);
      if (!state.opened) {
        this.status = 'Unable to initialize the Wails native dispatcher';
        hilog.error(DOMAIN, TAG, '%{public}s', this.status);
        return;
      }
    }
    native.start();
    // The Go runtime starts on a native thread. Do not create ArkWeb until
    // its application, window and embedded asset server have been initialized.
    this.timer = setInterval((): void => {
      if (native.ready()) {
        clearInterval(this.timer);
        this.timer = -1;
        this.runtimeReady = true;
        hilog.info(DOMAIN, TAG, 'Runtime ready; creating local-origin ArkWeb');
      }
    }, 20);
  }

  aboutToDisappear(): void {
    this.active = false;
    if (this.timer >= 0) {
      clearInterval(this.timer);
      this.timer = -1;
    }
    this.runtimeReady = false;
    setWebController(null);
    const state: BridgeState = bridgeState();
    if (state.dispatch === this.dispatchWeb) {
      state.dispatch = null;
      setPageDispatcher(null);
    }
  }

  /**
   * Serve an application asset from the Go backend.
   *
   * Only the controlled origin is served; every other request is refused so
   * the page can never reach the network through the shell.
   */
  private resource(url: string): WebResourceResponse {
    hilog.info(DOMAIN, TAG, 'Intercept %{public}s', url);
    const response = new WebResourceResponse();
    if (!url.startsWith(ORIGIN)) {
      response.setResponseCode(403);
      response.setReasonMessage('Forbidden');
      response.setResponseMimeType('text/plain');
      response.setResponseData('External resources are disabled');
      return response;
    }
    // The fragment is never sent to an HTTP server; the asset lookup likewise
    // receives just the path, with query parameters left out of the filename.
    const path = '/' + url.substring(ORIGIN.length).split('?')[0].split('#')[0];
    const served: native.Asset | null = native.asset(path);
    if (served === null || served === undefined) {
      response.setResponseCode(404);
      response.setReasonMessage('Not Found');
      response.setResponseMimeType('text/plain');
      response.setResponseData('Not Found');
      return response;
    }
    response.setResponseCode(200);
    response.setReasonMessage('OK');
    response.setResponseMimeType(served.mime);
    response.setResponseEncoding('UTF-8');
    response.setResponseData(served.body);
    return response;
  }

  build() {
    Column() {
      if (this.runtimeReady) {
        Web({ src: '', controller: this.controller })
          .width('100%')
          .height('100%')
          .javaScriptAccess(true)
          .domStorageAccess(true)
          .fileAccess(false)
          .allowWindowOpenMethod(false)
          .javaScriptProxy({
            object: this.proxy,
            name: 'wailsOpenHarmony',
            methodList: ['invoke', 'ready', 'request', 'report'],
            controller: this.controller
          })
          .onControllerAttached((): void => {
            setWebController(this.controller);
            hilog.info(DOMAIN, TAG, 'Controller attached; loading local origin');
            this.controller.loadUrl(ORIGIN);
          })
          .onLoadIntercept((event: OnLoadInterceptEvent): boolean => {
            return event === undefined || !event.data.getRequestUrl().startsWith(ORIGIN);
          })
          .onInterceptRequest((event: OnInterceptRequestEvent): WebResourceResponse => {
            return this.resource(event === undefined ? '' : event.request.getRequestUrl());
          })
          .onPageBegin((): void => {
            hilog.info(DOMAIN, TAG, 'Page loading');
          })
          .onPageEnd((): void => {
            hilog.info(DOMAIN, TAG, 'Local-origin frontend loaded');
          })
          .onErrorReceive((event: OnErrorReceiveEvent): void => {
            hilog.error(DOMAIN, TAG, 'Web load error: %{public}s', event.error.getErrorInfo());
          })
      } else {
        Text(this.status).fontSize(20)
      }
    }
    .width('100%')
    .height('100%')
  }
}
