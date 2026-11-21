import * as native from 'libentry.so';
import { webview } from '@kit.ArkWeb';
import { hilog } from '@kit.PerformanceAnalysisKit';

const ORIGIN: string = 'https://wails.localhost/';

interface ScriptCommand {
  script: string;
}

interface MainThreadCommand {
  id: number;
}

class WailsProxy {
  invoke(value: string): void {
    native.message(value);
  }

  request(id: string, payload: string): boolean {
    return native.request(id, payload);
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
  private proxy: WailsProxy = new WailsProxy();
  private timer: number = -1;
  private active: boolean = false;

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
        const command = JSON.parse(payload) as ScriptCommand;
        this.controller.runJavaScript(command.script).then((value: string): void => {
          this.reply(id, value);
        }).catch((error: Error): void => {
          hilog.error(0x0, 'WailsM3', 'JavaScript error: %{public}s', error.message);
          this.reply(id, '');
        });
        return;
      }
      if (name === 'getScreenInfoJson') {
        this.reply(id, '{}');
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
    if (!native.open((id: number, name: string, payload: string): void => {
      this.dispatch(id, name, payload);
    })) {
      this.status = 'Unable to initialize the Wails native dispatcher';
      hilog.error(0x0, 'WailsM3', '%{public}s', this.status);
      return;
    }
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
    if (this.timer >= 0) {
      clearInterval(this.timer);
      this.timer = -1;
    }
    this.runtimeReady = false;
    native.close();
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
            methodList: ['invoke', 'request', 'report'],
            controller: this.controller
          })
          .onControllerAttached((): void => {
            hilog.info(0x0, 'WailsM3', 'Controller attached; loading local origin');
            this.controller.loadUrl(ORIGIN);
          })
          .onLoadIntercept((event: OnLoadInterceptEvent): boolean => {
            return event === undefined || !event.data.getRequestUrl().startsWith(ORIGIN);
          })
          .onInterceptRequest((event: OnInterceptRequestEvent): WebResourceResponse => {
            return this.resource(event === undefined ? '' : event.request.getRequestUrl());
          })
          .onErrorReceive((event: OnErrorReceiveEvent): void => {
            hilog.error(0x0, 'WailsM3', 'Web load error: %{public}s', event.error.getErrorInfo());
          })
          .onPageEnd((): void => {
            hilog.info(0x0, 'WailsM3', 'Local-origin frontend loaded');
          })
      } else {
        Text(this.status).fontSize(20)
      }
    }
    .width('100%')
    .height('100%')
  }
}
