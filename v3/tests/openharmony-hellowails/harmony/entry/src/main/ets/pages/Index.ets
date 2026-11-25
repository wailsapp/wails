import * as native from 'libentry.so';
import { webview } from '@kit.ArkWeb';
import { hilog } from '@kit.PerformanceAnalysisKit';
import { display, promptAction } from '@kit.ArkUI';
import { pasteboard } from '@kit.BasicServicesKit';
import { BusinessError } from '@kit.BasicServicesKit';
import { getContext } from '../common/WailsContext';

const ORIGIN: string = 'https://wails.localhost/';

interface ScriptCommand {
  script: string;
}

interface MainThreadCommand {
  id: number;
}

interface LogArgs {
  level?: string;
  message?: string;
}

interface DialogButton {
  label: string;
}

interface DialogOptions {
  title: string;
  message: string;
  buttons: DialogButton[];
}

interface DialogArgs {
  id: number;
  options: string;
}

interface ThemeConfig {
  colorMode?: number;
}

interface ThemeContext {
  config?: ThemeConfig;
}

interface PromptButton {
  text: string;
  color: string;
}

// The Go app can outlive a single stage, so the bridge belongs to the process,
// not to a page. Pages attach their Web dispatcher and clear it on exit; the
// bridge itself is opened once (AppStorage survives stage recreation).
class BridgeState {
  opened: boolean = false;
  dispatch: native.Dispatcher | null = null;
}

function bridgeState(): BridgeState {
  const state: BridgeState | undefined = AppStorage.get<BridgeState>('wailsBridgeState');
  if (state !== undefined) {
    return state;
  }
  const created: BridgeState = new BridgeState();
  AppStorage.setOrCreate('wailsBridgeState', created);
  return created;
}

function dispatchBridge(id: number, name: string, payload: string): void {
  if (name === 'runOnMainThread') {
    // Main-thread work must be drained even while no page owns a Web view;
    // otherwise Go's dispatchOnMainThread never wakes and its panic fires.
    const command = JSON.parse(payload) as MainThreadCommand;
    native.runOnMainThread(command.id);
    return;
  }
  const callback: native.Dispatcher | null = bridgeState().dispatch;
  if (callback !== null) {
    callback(id, name, payload);
    return;
  }
  if (id > 0) {
    native.result(id, '');
  }
}

class WailsProxy {
  private advance: () => void;
  session: string = '';

  constructor(advance: () => void) {
    this.advance = advance;
  }

  advanceLifecycleProbe(): void {
    this.advance();
  }

  invoke(value: string): void {
    native.message(value);
  }

  ready(session: string): boolean {
    if (this.session !== '' || !/^[0-9a-f-]{36}$/.test(session)) {
      return false;
    }
    this.session = session;
    return true;
  }

  request(session: string, id: string, payload: string): boolean {
    return this.session !== '' && session === this.session && native.request(`${session}:${id}`, payload);
  }

  report(message: string): void {
    hilog.info(0x0, 'WailsM3', '%{public}s', message);
  }
}

@Entry
@Component
struct Index {
  @State private runtimeReady: boolean = false;
  @State private status: string = 'Starting Wails runtime';
  private controller: webview.WebviewController = new webview.WebviewController();
  private proxy: WailsProxy = new WailsProxy(() => this.advanceLifecycleProbe());
  @StorageProp('wailsLifecycle') private probe: string = '';
  @StorageProp('wailsPhase') private phase: string = '';
  private probeAdvanced: boolean = false;
  private rebuildTimer: number = -1;

  private entryURL(): string {
    if (this.probe === 'stage') {
      return `${ORIGIN}?probe=stage&phase=${this.phase === 'new' ? 'new' : 'old'}`;
    }
    return this.probe === '' ? ORIGIN : `${ORIGIN}?probe=${this.probe}&phase=${this.probeAdvanced ? 'new' : 'old'}`;
  }

  private advanceLifecycleProbe(): void {
    if (this.probe === '' || this.probe === 'stage' || this.probeAdvanced || !this.active) {
      return;
    }
    this.probeAdvanced = true;
    this.proxy.session = '';
    if (this.probe === 'refresh') {
      this.controller.loadUrl(this.entryURL());
      return;
    }
    this.runtimeReady = false;
    this.rebuildTimer = setTimeout((): void => {
      this.rebuildTimer = -1;
      this.controller = new webview.WebviewController();
      this.runtimeReady = true;
    }, 200);
  }
  private timer: number = -1;
  private active: boolean = false;
  private dispatchWeb: native.Dispatcher = (id: number, name: string, payload: string): void => {
    this.dispatch(id, name, payload);
  };

  private reply(id: number, value: string): void {
    if (id > 0 && this.active) {
      native.result(id, value);
    }
  }

  private dispatch(id: number, name: string, payload: string): void {
    if (!this.active) {
      return;
    }
    try {
      if (name === 'executeJavaScript') {
        const session: string = this.proxy.session;
        if (session === '') {
          this.reply(id, '');
          return;
        }
        const command = JSON.parse(payload) as ScriptCommand;
        // Recheck inside the document in case ArkWeb runs the command after navigation.
        const script: string = `if (window.wailsOpenHarmonySession === ${JSON.stringify(session)}) {\n${command.script}\n}`;
        this.controller.runJavaScript(script).then((value: string): void => {
          this.reply(id, value);
        }).catch((error: Error): void => {
          hilog.error(0x0, 'WailsM3', 'JavaScript error: %{public}s', error.message);
          this.reply(id, '');
        });
        return;
      }
      if (name === 'getScreenInfoJson') {
        const current = display.getDefaultDisplaySync();
        this.reply(id, JSON.stringify({
          widthPx: current.width,
          heightPx: current.height,
          density: current.densityPixels,
          insetTop: 0,
          insetBottom: 0,
          insetLeft: 0,
          insetRight: 0
        }));
        return;
      }
      if (name === 'isDarkMode') {
        const context = getContext() as ThemeContext;
        const dark: boolean = context.config?.colorMode === 0;
        this.reply(id, JSON.stringify({ dark: dark }));
        return;
      }
      if (name === 'getClipboardText') {
        const text = pasteboard.getSystemPasteboard().getDataSync().getPrimaryText();
        this.reply(id, JSON.stringify({ text: text }));
        return;
      }
      if (name === 'log') {
        const args = JSON.parse(payload) as LogArgs;
        hilog.info(0x0, 'WailsM3', '[%{public}s] %{public}s', args.level ?? 'INFO', args.message ?? '');
        this.reply(id, '');
        return;
      }
      if (name === 'showMessageDialog') {
        const args = JSON.parse(payload) as DialogArgs;
        const options = JSON.parse(args.options) as DialogOptions;
        const buttons = (options.buttons ?? []).slice(0, 3);
        if (buttons.length === 0) {
          native.dialogResult(args.id ?? 0, -1);
          this.reply(id, '');
          return;
        }
        const styled: promptAction.Button[] = buttons.map((button: DialogButton): promptAction.Button => {
          const styledButton: PromptButton = { text: button.label, color: '#000000' };
          return styledButton as promptAction.Button;
        });
        promptAction.showDialog({
          title: options.title,
          message: options.message,
          buttons: styled
        }).then((result: promptAction.ShowDialogSuccessResponse): void => {
          native.dialogResult(args.id ?? 0,
            result.index >= 0 && result.index < buttons.length ? result.index : -1);
          this.reply(id, '');
        }).catch((error: BusinessError): void => {
          hilog.error(0x0, 'WailsM3', 'Dialog failed: %{public}s', error.message);
          native.dialogResult(args.id ?? 0, -1);
          this.reply(id, '');
        });
        return;
      }
      if (name === 'runOnMainThread') {
        const command = JSON.parse(payload) as MainThreadCommand;
        native.runOnMainThread(command.id);
        this.reply(id, '');
        return;
      }
      hilog.error(0x0, 'WailsM3', 'Unsupported dispatcher method: %{public}s', name);
      this.reply(id, '');
    } catch (error) {
      hilog.error(0x0, 'WailsM3', 'Dispatcher error: %{public}s', String(error));
      this.reply(id, '');
    }
  }

  aboutToAppear(): void {
    this.active = true;
    const state: BridgeState = bridgeState();
    if (!state.opened) {
      state.opened = native.open(dispatchBridge);
    }
    if (!state.opened) {
      this.status = 'Unable to initialize the Wails native dispatcher';
      hilog.error(0x0, 'WailsM3', '%{public}s', this.status);
      return;
    }
    state.dispatch = this.dispatchWeb;
    native.start();
    // The Go application starts on a native thread. Do not create ArkWeb until
    // its application, window and embedded asset server have been initialized.
    this.timer = setInterval((): void => {
      if (native.ready()) {
        clearInterval(this.timer);
        this.timer = -1;
        this.runtimeReady = true;
        hilog.info(0x0, 'WailsM3', 'Runtime ready; creating local-origin ArkWeb');
      }
    }, 20);
  }

  aboutToDisappear(): void {
    this.active = false;
    this.proxy.session = '';
    if (this.rebuildTimer >= 0) {
      clearTimeout(this.rebuildTimer);
      this.rebuildTimer = -1;
    }
    if (this.timer >= 0) {
      clearInterval(this.timer);
      this.timer = -1;
    }
    this.runtimeReady = false;
    const state: BridgeState = bridgeState();
    if (state.dispatch === this.dispatchWeb) {
      state.dispatch = null;
    }
  }

  private resource(url: string): WebResourceResponse {
    hilog.info(0x0, 'WailsM3', 'Intercept %{public}s', url);
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
    const asset: native.Asset | null = native.asset(path);
    if (asset === null) {
      response.setResponseCode(404);
      response.setReasonMessage('Not Found');
      response.setResponseMimeType('text/plain');
      response.setResponseData('Not Found');
      return response;
    }
    response.setResponseCode(200);
    response.setReasonMessage('OK');
    response.setResponseMimeType(asset.mime);
    response.setResponseEncoding('UTF-8');
    response.setResponseData(asset.body);
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
            methodList: ['invoke', 'ready', 'request', 'report', 'advanceLifecycleProbe'],
            controller: this.controller
          })
          .onControllerAttached((): void => {
            hilog.info(0x0, 'WailsM3', 'Controller attached; loading local origin');
            this.controller.loadUrl(this.entryURL());
          })
          .onLoadIntercept((event: OnLoadInterceptEvent): boolean => {
            return event === undefined || !event.data.getRequestUrl().startsWith(ORIGIN);
          })
          .onInterceptRequest((event: OnInterceptRequestEvent): WebResourceResponse => {
            return this.resource(event === undefined ? '' : event.request.getRequestUrl());
          })
          .onPageBegin((event: OnPageBeginEvent): void => {
            this.proxy.session = '';
            native.pageStarted(event === undefined ? ORIGIN : event.url);
          })
          .onErrorReceive((event: OnErrorReceiveEvent): void => {
            hilog.error(0x0, 'WailsM3', 'Web load error: %{public}s', event.error.getErrorInfo());
          })
          .onPageEnd((event: OnPageEndEvent): void => {
            hilog.info(0x0, 'WailsM3', 'Local-origin frontend loaded');
            native.pageFinished(event === undefined ? ORIGIN : event.url);
          })
      } else {
        Text(this.status).fontSize(20)
      }
    }
    .width('100%')
    .height('100%')
  }
}
