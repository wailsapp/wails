
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
            public ICoreWebView2PermissionRequestedEventHandler,
            public ICoreWebView2AcceleratorKeyPressedEventHandler
{

    struct Application *app;
    HWND window;
    messageCallback mcb;
    comHandlerCallback cb;

    public:
        wv2ComHandler(struct Application *app, HWND window, messageCallback mcb, comHandlerCallback cb) {
            this->app = app;
            this->window = window;
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
          env->CreateCoreWebView2Controller(window, this);
          return S_OK;
        }
        HRESULT STDMETHODCALLTYPE Invoke(HRESULT res,
                                         ICoreWebView2Controller *controller) {
          controller->AddRef();

          ICoreWebView2 *webview;
          ::EventRegistrationToken token;
          controller->get_CoreWebView2(&webview);
          controller->add_AcceleratorKeyPressed(this, &token);
          webview->add_WebMessageReceived(this, &token);
          webview->add_PermissionRequested(this, &token);

          cb(controller);
          return S_OK;
        }

        // This is our keyboard callback method
        HRESULT STDMETHODCALLTYPE Invoke(ICoreWebView2Controller *controller, ICoreWebView2AcceleratorKeyPressedEventArgs * args) {
            // Prevent WebView2 from processing the key
            args->put_Handled(TRUE);

            COREWEBVIEW2_KEY_EVENT_KIND kind;
            args->get_KeyEventKind(&kind);
            if (kind == COREWEBVIEW2_KEY_EVENT_KIND_KEY_DOWN ||
                kind == COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN)
            {
                UINT key;
                args->get_VirtualKey(&key);
                COREWEBVIEW2_PHYSICAL_KEY_STATUS status;
                args->get_PhysicalKeyStatus(&status);
                if (!status.WasKeyDown)
                {
                    processKeyPress(key);
                }
            }
            return S_OK;
        }

        // This is called when JS posts a message back to webkit
        HRESULT STDMETHODCALLTYPE Invoke(
            ICoreWebView2 *sender, ICoreWebView2WebMessageReceivedEventArgs *args) {
          LPWSTR message;
          args->TryGetWebMessageAsString(&message);
          if ( message == nullptr ) {
            return S_OK;
          }
          const char *m = LPWSTRToCstr(message);

          // check for internal messages
          if (strcmp(m, "completed") == 0) {
            completed(app);
            return S_OK;
          }
          else if (strcmp(m, "initialised") == 0) {
            loadAssets(app);
            return S_OK;
          }
          else if (strcmp(m, "wails-drag") == 0) {
            // We don't drag in fullscreen mode
            if (!app->isFullscreen) {
                ReleaseCapture();
                SendMessage(this->window, WM_NCLBUTTONDOWN, HTCAPTION, 0);
            }
            return S_OK;
          }
          else {
            messageFromWindowCallback(m);
          }
          delete[] m;
          return S_OK;
        }
        HRESULT STDMETHODCALLTYPE
        Invoke(ICoreWebView2 *sender,
               ICoreWebView2PermissionRequestedEventArgs *args) {
                                                        printf("DDDDDDDDDDDD\n");

          COREWEBVIEW2_PERMISSION_KIND kind;
          args->get_PermissionKind(&kind);
          if (kind == COREWEBVIEW2_PERMISSION_KIND_CLIPBOARD_READ) {
            args->put_State(COREWEBVIEW2_PERMISSION_STATE_ALLOW);
          }
          return S_OK;
        }

};

#endif