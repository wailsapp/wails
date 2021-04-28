
#ifndef WV2COMHANDLER_H
#define WV2COMHANDLER_H

#include "ffenestri_windows.h"
#include "windows/WebView2.h"

#include <locale>
#include <codecvt>

class wv2ComHandler
        :   public ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler,
            public ICoreWebView2CreateCoreWebView2ControllerCompletedHandler,
            public ICoreWebView2WebMessageReceivedEventHandler,
            public ICoreWebView2PermissionRequestedEventHandler {

    struct Application *app;
    messageCallback mcb;
    comHandlerCallback cb;

    public:
        wv2ComHandler(struct Application *app, messageCallback mcb, comHandlerCallback cb) {
            this->app = app;
            this->mcb = mcb;
            this->cb = cb;
        }
        ULONG STDMETHODCALLTYPE AddRef() { return 1; }
        ULONG STDMETHODCALLTYPE Release() { return 1; }
        HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, LPVOID *ppv) {
          return S_OK;
        }
        HRESULT STDMETHODCALLTYPE Invoke(HRESULT res,
                                         ICoreWebView2Environment *env) {
          env->CreateCoreWebView2Controller(app->window, this);
          return S_OK;
        }
        HRESULT STDMETHODCALLTYPE Invoke(HRESULT res,
                                         ICoreWebView2Controller *controller) {
          controller->AddRef();

          ICoreWebView2 *webview;
          ::EventRegistrationToken token;
          controller->get_CoreWebView2(&webview);
          webview->add_WebMessageReceived(this, &token);
          webview->add_PermissionRequested(this, &token);

          cb(controller);
          return S_OK;
        }
        HRESULT STDMETHODCALLTYPE Invoke(
            ICoreWebView2 *sender, ICoreWebView2WebMessageReceivedEventArgs *args) {
          LPWSTR message;
          args->TryGetWebMessageAsString(&message);

          std::wstring_convert<std::codecvt_utf8_utf16<wchar_t>> wideCharConverter;
          mcb(wideCharConverter.to_bytes(message));
          sender->PostWebMessageAsString(message);

          CoTaskMemFree(message);
          return S_OK;
        }
        HRESULT STDMETHODCALLTYPE
        Invoke(ICoreWebView2 *sender,
               ICoreWebView2PermissionRequestedEventArgs *args) {
          COREWEBVIEW2_PERMISSION_KIND kind;
          args->get_PermissionKind(&kind);
          if (kind == COREWEBVIEW2_PERMISSION_KIND_CLIPBOARD_READ) {
            args->put_State(COREWEBVIEW2_PERMISSION_STATE_ALLOW);
          }
          return S_OK;
        }

};

#endif