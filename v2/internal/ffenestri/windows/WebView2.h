

/* this ALWAYS GENERATED file contains the definitions for the interfaces */


 /* File created by MIDL compiler version 8.xx.xxxx */
/* at a redacted point in time
 */
/* Compiler settings for ../../edge_embedded_browser/client/win/current/webview2.idl:
    Oicf, W1, Zp8, env=Win64 (32b run), target_arch=AMD64 8.xx.xxxx 
    protocol : dce , ms_ext, c_ext, robust
    error checks: allocation ref bounds_check enum stub_data 
    VC __declspec() decoration level: 
         __declspec(uuid()), __declspec(selectany), __declspec(novtable)
         DECLSPEC_UUID(), MIDL_INTERFACE()
*/
/* @@MIDL_FILE_HEADING(  ) */

#pragma warning( disable: 4049 )  /* more than 64k source lines */


/* verify that the <rpcndr.h> version is high enough to compile this file*/
#ifndef __REQUIRED_RPCNDR_H_VERSION__
#define __REQUIRED_RPCNDR_H_VERSION__ 475
#endif

#include "rpc.h"
#include "rpcndr.h"

#ifndef __RPCNDR_H_VERSION__
#error this stub requires an updated version of <rpcndr.h>
#endif /* __RPCNDR_H_VERSION__ */


#ifndef __webview2_h__
#define __webview2_h__

#if defined(_MSC_VER) && (_MSC_VER >= 1020)
#pragma once
#endif

/* Forward Declarations */ 

#ifndef __ICoreWebView2AcceleratorKeyPressedEventArgs_FWD_DEFINED__
#define __ICoreWebView2AcceleratorKeyPressedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2AcceleratorKeyPressedEventArgs ICoreWebView2AcceleratorKeyPressedEventArgs;

#endif 	/* __ICoreWebView2AcceleratorKeyPressedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2AcceleratorKeyPressedEventHandler_FWD_DEFINED__
#define __ICoreWebView2AcceleratorKeyPressedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2AcceleratorKeyPressedEventHandler ICoreWebView2AcceleratorKeyPressedEventHandler;

#endif 	/* __ICoreWebView2AcceleratorKeyPressedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler;

#endif 	/* __ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2CallDevToolsProtocolMethodCompletedHandler ICoreWebView2CallDevToolsProtocolMethodCompletedHandler;

#endif 	/* __ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CapturePreviewCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2CapturePreviewCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2CapturePreviewCompletedHandler ICoreWebView2CapturePreviewCompletedHandler;

#endif 	/* __ICoreWebView2CapturePreviewCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2_FWD_DEFINED__
#define __ICoreWebView2_FWD_DEFINED__
typedef interface ICoreWebView2 ICoreWebView2;

#endif 	/* __ICoreWebView2_FWD_DEFINED__ */


#ifndef __ICoreWebView2_2_FWD_DEFINED__
#define __ICoreWebView2_2_FWD_DEFINED__
typedef interface ICoreWebView2_2 ICoreWebView2_2;

#endif 	/* __ICoreWebView2_2_FWD_DEFINED__ */


#ifndef __ICoreWebView2_3_FWD_DEFINED__
#define __ICoreWebView2_3_FWD_DEFINED__
typedef interface ICoreWebView2_3 ICoreWebView2_3;

#endif 	/* __ICoreWebView2_3_FWD_DEFINED__ */


#ifndef __ICoreWebView2CompositionController_FWD_DEFINED__
#define __ICoreWebView2CompositionController_FWD_DEFINED__
typedef interface ICoreWebView2CompositionController ICoreWebView2CompositionController;

#endif 	/* __ICoreWebView2CompositionController_FWD_DEFINED__ */


#ifndef __ICoreWebView2CompositionController2_FWD_DEFINED__
#define __ICoreWebView2CompositionController2_FWD_DEFINED__
typedef interface ICoreWebView2CompositionController2 ICoreWebView2CompositionController2;

#endif 	/* __ICoreWebView2CompositionController2_FWD_DEFINED__ */


#ifndef __ICoreWebView2Controller_FWD_DEFINED__
#define __ICoreWebView2Controller_FWD_DEFINED__
typedef interface ICoreWebView2Controller ICoreWebView2Controller;

#endif 	/* __ICoreWebView2Controller_FWD_DEFINED__ */


#ifndef __ICoreWebView2Controller2_FWD_DEFINED__
#define __ICoreWebView2Controller2_FWD_DEFINED__
typedef interface ICoreWebView2Controller2 ICoreWebView2Controller2;

#endif 	/* __ICoreWebView2Controller2_FWD_DEFINED__ */


#ifndef __ICoreWebView2Controller3_FWD_DEFINED__
#define __ICoreWebView2Controller3_FWD_DEFINED__
typedef interface ICoreWebView2Controller3 ICoreWebView2Controller3;

#endif 	/* __ICoreWebView2Controller3_FWD_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventArgs_FWD_DEFINED__
#define __ICoreWebView2ContentLoadingEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2ContentLoadingEventArgs ICoreWebView2ContentLoadingEventArgs;

#endif 	/* __ICoreWebView2ContentLoadingEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventHandler_FWD_DEFINED__
#define __ICoreWebView2ContentLoadingEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ContentLoadingEventHandler ICoreWebView2ContentLoadingEventHandler;

#endif 	/* __ICoreWebView2ContentLoadingEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2Cookie_FWD_DEFINED__
#define __ICoreWebView2Cookie_FWD_DEFINED__
typedef interface ICoreWebView2Cookie ICoreWebView2Cookie;

#endif 	/* __ICoreWebView2Cookie_FWD_DEFINED__ */


#ifndef __ICoreWebView2CookieList_FWD_DEFINED__
#define __ICoreWebView2CookieList_FWD_DEFINED__
typedef interface ICoreWebView2CookieList ICoreWebView2CookieList;

#endif 	/* __ICoreWebView2CookieList_FWD_DEFINED__ */


#ifndef __ICoreWebView2CookieManager_FWD_DEFINED__
#define __ICoreWebView2CookieManager_FWD_DEFINED__
typedef interface ICoreWebView2CookieManager ICoreWebView2CookieManager;

#endif 	/* __ICoreWebView2CookieManager_FWD_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler;

#endif 	/* __ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2CreateCoreWebView2ControllerCompletedHandler ICoreWebView2CreateCoreWebView2ControllerCompletedHandler;

#endif 	/* __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler;

#endif 	/* __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2ContainsFullScreenElementChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2ContainsFullScreenElementChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ContainsFullScreenElementChangedEventHandler ICoreWebView2ContainsFullScreenElementChangedEventHandler;

#endif 	/* __ICoreWebView2ContainsFullScreenElementChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CursorChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2CursorChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2CursorChangedEventHandler ICoreWebView2CursorChangedEventHandler;

#endif 	/* __ICoreWebView2CursorChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2DocumentTitleChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2DocumentTitleChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2DocumentTitleChangedEventHandler ICoreWebView2DocumentTitleChangedEventHandler;

#endif 	/* __ICoreWebView2DocumentTitleChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2DOMContentLoadedEventArgs_FWD_DEFINED__
#define __ICoreWebView2DOMContentLoadedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2DOMContentLoadedEventArgs ICoreWebView2DOMContentLoadedEventArgs;

#endif 	/* __ICoreWebView2DOMContentLoadedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2DOMContentLoadedEventHandler_FWD_DEFINED__
#define __ICoreWebView2DOMContentLoadedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2DOMContentLoadedEventHandler ICoreWebView2DOMContentLoadedEventHandler;

#endif 	/* __ICoreWebView2DOMContentLoadedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2Deferral_FWD_DEFINED__
#define __ICoreWebView2Deferral_FWD_DEFINED__
typedef interface ICoreWebView2Deferral ICoreWebView2Deferral;

#endif 	/* __ICoreWebView2Deferral_FWD_DEFINED__ */


#ifndef __ICoreWebView2DevToolsProtocolEventReceivedEventArgs_FWD_DEFINED__
#define __ICoreWebView2DevToolsProtocolEventReceivedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2DevToolsProtocolEventReceivedEventArgs ICoreWebView2DevToolsProtocolEventReceivedEventArgs;

#endif 	/* __ICoreWebView2DevToolsProtocolEventReceivedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2DevToolsProtocolEventReceivedEventHandler_FWD_DEFINED__
#define __ICoreWebView2DevToolsProtocolEventReceivedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2DevToolsProtocolEventReceivedEventHandler ICoreWebView2DevToolsProtocolEventReceivedEventHandler;

#endif 	/* __ICoreWebView2DevToolsProtocolEventReceivedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2DevToolsProtocolEventReceiver_FWD_DEFINED__
#define __ICoreWebView2DevToolsProtocolEventReceiver_FWD_DEFINED__
typedef interface ICoreWebView2DevToolsProtocolEventReceiver ICoreWebView2DevToolsProtocolEventReceiver;

#endif 	/* __ICoreWebView2DevToolsProtocolEventReceiver_FWD_DEFINED__ */


#ifndef __ICoreWebView2Environment_FWD_DEFINED__
#define __ICoreWebView2Environment_FWD_DEFINED__
typedef interface ICoreWebView2Environment ICoreWebView2Environment;

#endif 	/* __ICoreWebView2Environment_FWD_DEFINED__ */


#ifndef __ICoreWebView2Environment2_FWD_DEFINED__
#define __ICoreWebView2Environment2_FWD_DEFINED__
typedef interface ICoreWebView2Environment2 ICoreWebView2Environment2;

#endif 	/* __ICoreWebView2Environment2_FWD_DEFINED__ */


#ifndef __ICoreWebView2Environment3_FWD_DEFINED__
#define __ICoreWebView2Environment3_FWD_DEFINED__
typedef interface ICoreWebView2Environment3 ICoreWebView2Environment3;

#endif 	/* __ICoreWebView2Environment3_FWD_DEFINED__ */


#ifndef __ICoreWebView2Environment4_FWD_DEFINED__
#define __ICoreWebView2Environment4_FWD_DEFINED__
typedef interface ICoreWebView2Environment4 ICoreWebView2Environment4;

#endif 	/* __ICoreWebView2Environment4_FWD_DEFINED__ */


#ifndef __ICoreWebView2EnvironmentOptions_FWD_DEFINED__
#define __ICoreWebView2EnvironmentOptions_FWD_DEFINED__
typedef interface ICoreWebView2EnvironmentOptions ICoreWebView2EnvironmentOptions;

#endif 	/* __ICoreWebView2EnvironmentOptions_FWD_DEFINED__ */


#ifndef __ICoreWebView2ExecuteScriptCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2ExecuteScriptCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2ExecuteScriptCompletedHandler ICoreWebView2ExecuteScriptCompletedHandler;

#endif 	/* __ICoreWebView2ExecuteScriptCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2FrameInfo_FWD_DEFINED__
#define __ICoreWebView2FrameInfo_FWD_DEFINED__
typedef interface ICoreWebView2FrameInfo ICoreWebView2FrameInfo;

#endif 	/* __ICoreWebView2FrameInfo_FWD_DEFINED__ */


#ifndef __ICoreWebView2FrameInfoCollection_FWD_DEFINED__
#define __ICoreWebView2FrameInfoCollection_FWD_DEFINED__
typedef interface ICoreWebView2FrameInfoCollection ICoreWebView2FrameInfoCollection;

#endif 	/* __ICoreWebView2FrameInfoCollection_FWD_DEFINED__ */


#ifndef __ICoreWebView2FrameInfoCollectionIterator_FWD_DEFINED__
#define __ICoreWebView2FrameInfoCollectionIterator_FWD_DEFINED__
typedef interface ICoreWebView2FrameInfoCollectionIterator ICoreWebView2FrameInfoCollectionIterator;

#endif 	/* __ICoreWebView2FrameInfoCollectionIterator_FWD_DEFINED__ */


#ifndef __ICoreWebView2FocusChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2FocusChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2FocusChangedEventHandler ICoreWebView2FocusChangedEventHandler;

#endif 	/* __ICoreWebView2FocusChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2GetCookiesCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2GetCookiesCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2GetCookiesCompletedHandler ICoreWebView2GetCookiesCompletedHandler;

#endif 	/* __ICoreWebView2GetCookiesCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2HistoryChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2HistoryChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2HistoryChangedEventHandler ICoreWebView2HistoryChangedEventHandler;

#endif 	/* __ICoreWebView2HistoryChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2HttpHeadersCollectionIterator_FWD_DEFINED__
#define __ICoreWebView2HttpHeadersCollectionIterator_FWD_DEFINED__
typedef interface ICoreWebView2HttpHeadersCollectionIterator ICoreWebView2HttpHeadersCollectionIterator;

#endif 	/* __ICoreWebView2HttpHeadersCollectionIterator_FWD_DEFINED__ */


#ifndef __ICoreWebView2HttpRequestHeaders_FWD_DEFINED__
#define __ICoreWebView2HttpRequestHeaders_FWD_DEFINED__
typedef interface ICoreWebView2HttpRequestHeaders ICoreWebView2HttpRequestHeaders;

#endif 	/* __ICoreWebView2HttpRequestHeaders_FWD_DEFINED__ */


#ifndef __ICoreWebView2HttpResponseHeaders_FWD_DEFINED__
#define __ICoreWebView2HttpResponseHeaders_FWD_DEFINED__
typedef interface ICoreWebView2HttpResponseHeaders ICoreWebView2HttpResponseHeaders;

#endif 	/* __ICoreWebView2HttpResponseHeaders_FWD_DEFINED__ */


#ifndef __ICoreWebView2Interop_FWD_DEFINED__
#define __ICoreWebView2Interop_FWD_DEFINED__
typedef interface ICoreWebView2Interop ICoreWebView2Interop;

#endif 	/* __ICoreWebView2Interop_FWD_DEFINED__ */


#ifndef __ICoreWebView2MoveFocusRequestedEventArgs_FWD_DEFINED__
#define __ICoreWebView2MoveFocusRequestedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2MoveFocusRequestedEventArgs ICoreWebView2MoveFocusRequestedEventArgs;

#endif 	/* __ICoreWebView2MoveFocusRequestedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2MoveFocusRequestedEventHandler_FWD_DEFINED__
#define __ICoreWebView2MoveFocusRequestedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2MoveFocusRequestedEventHandler ICoreWebView2MoveFocusRequestedEventHandler;

#endif 	/* __ICoreWebView2MoveFocusRequestedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2NavigationCompletedEventArgs_FWD_DEFINED__
#define __ICoreWebView2NavigationCompletedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2NavigationCompletedEventArgs ICoreWebView2NavigationCompletedEventArgs;

#endif 	/* __ICoreWebView2NavigationCompletedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2NavigationCompletedEventHandler_FWD_DEFINED__
#define __ICoreWebView2NavigationCompletedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2NavigationCompletedEventHandler ICoreWebView2NavigationCompletedEventHandler;

#endif 	/* __ICoreWebView2NavigationCompletedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2NavigationStartingEventArgs_FWD_DEFINED__
#define __ICoreWebView2NavigationStartingEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2NavigationStartingEventArgs ICoreWebView2NavigationStartingEventArgs;

#endif 	/* __ICoreWebView2NavigationStartingEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2NavigationStartingEventHandler_FWD_DEFINED__
#define __ICoreWebView2NavigationStartingEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2NavigationStartingEventHandler ICoreWebView2NavigationStartingEventHandler;

#endif 	/* __ICoreWebView2NavigationStartingEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2NewBrowserVersionAvailableEventHandler_FWD_DEFINED__
#define __ICoreWebView2NewBrowserVersionAvailableEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2NewBrowserVersionAvailableEventHandler ICoreWebView2NewBrowserVersionAvailableEventHandler;

#endif 	/* __ICoreWebView2NewBrowserVersionAvailableEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2NewWindowRequestedEventArgs_FWD_DEFINED__
#define __ICoreWebView2NewWindowRequestedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2NewWindowRequestedEventArgs ICoreWebView2NewWindowRequestedEventArgs;

#endif 	/* __ICoreWebView2NewWindowRequestedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2NewWindowRequestedEventHandler_FWD_DEFINED__
#define __ICoreWebView2NewWindowRequestedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2NewWindowRequestedEventHandler ICoreWebView2NewWindowRequestedEventHandler;

#endif 	/* __ICoreWebView2NewWindowRequestedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2PermissionRequestedEventArgs_FWD_DEFINED__
#define __ICoreWebView2PermissionRequestedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2PermissionRequestedEventArgs ICoreWebView2PermissionRequestedEventArgs;

#endif 	/* __ICoreWebView2PermissionRequestedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2PermissionRequestedEventHandler_FWD_DEFINED__
#define __ICoreWebView2PermissionRequestedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2PermissionRequestedEventHandler ICoreWebView2PermissionRequestedEventHandler;

#endif 	/* __ICoreWebView2PermissionRequestedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2PointerInfo_FWD_DEFINED__
#define __ICoreWebView2PointerInfo_FWD_DEFINED__
typedef interface ICoreWebView2PointerInfo ICoreWebView2PointerInfo;

#endif 	/* __ICoreWebView2PointerInfo_FWD_DEFINED__ */


#ifndef __ICoreWebView2ProcessFailedEventArgs_FWD_DEFINED__
#define __ICoreWebView2ProcessFailedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2ProcessFailedEventArgs ICoreWebView2ProcessFailedEventArgs;

#endif 	/* __ICoreWebView2ProcessFailedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2ProcessFailedEventArgs2_FWD_DEFINED__
#define __ICoreWebView2ProcessFailedEventArgs2_FWD_DEFINED__
typedef interface ICoreWebView2ProcessFailedEventArgs2 ICoreWebView2ProcessFailedEventArgs2;

#endif 	/* __ICoreWebView2ProcessFailedEventArgs2_FWD_DEFINED__ */


#ifndef __ICoreWebView2ProcessFailedEventHandler_FWD_DEFINED__
#define __ICoreWebView2ProcessFailedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ProcessFailedEventHandler ICoreWebView2ProcessFailedEventHandler;

#endif 	/* __ICoreWebView2ProcessFailedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2RasterizationScaleChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2RasterizationScaleChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2RasterizationScaleChangedEventHandler ICoreWebView2RasterizationScaleChangedEventHandler;

#endif 	/* __ICoreWebView2RasterizationScaleChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2ScriptDialogOpeningEventArgs_FWD_DEFINED__
#define __ICoreWebView2ScriptDialogOpeningEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2ScriptDialogOpeningEventArgs ICoreWebView2ScriptDialogOpeningEventArgs;

#endif 	/* __ICoreWebView2ScriptDialogOpeningEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2ScriptDialogOpeningEventHandler_FWD_DEFINED__
#define __ICoreWebView2ScriptDialogOpeningEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ScriptDialogOpeningEventHandler ICoreWebView2ScriptDialogOpeningEventHandler;

#endif 	/* __ICoreWebView2ScriptDialogOpeningEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2Settings_FWD_DEFINED__
#define __ICoreWebView2Settings_FWD_DEFINED__
typedef interface ICoreWebView2Settings ICoreWebView2Settings;

#endif 	/* __ICoreWebView2Settings_FWD_DEFINED__ */


#ifndef __ICoreWebView2Settings2_FWD_DEFINED__
#define __ICoreWebView2Settings2_FWD_DEFINED__
typedef interface ICoreWebView2Settings2 ICoreWebView2Settings2;

#endif 	/* __ICoreWebView2Settings2_FWD_DEFINED__ */


#ifndef __ICoreWebView2Settings3_FWD_DEFINED__
#define __ICoreWebView2Settings3_FWD_DEFINED__
typedef interface ICoreWebView2Settings3 ICoreWebView2Settings3;

#endif 	/* __ICoreWebView2Settings3_FWD_DEFINED__ */


#ifndef __ICoreWebView2SourceChangedEventArgs_FWD_DEFINED__
#define __ICoreWebView2SourceChangedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2SourceChangedEventArgs ICoreWebView2SourceChangedEventArgs;

#endif 	/* __ICoreWebView2SourceChangedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2SourceChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2SourceChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2SourceChangedEventHandler ICoreWebView2SourceChangedEventHandler;

#endif 	/* __ICoreWebView2SourceChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2TrySuspendCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2TrySuspendCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2TrySuspendCompletedHandler ICoreWebView2TrySuspendCompletedHandler;

#endif 	/* __ICoreWebView2TrySuspendCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebMessageReceivedEventArgs_FWD_DEFINED__
#define __ICoreWebView2WebMessageReceivedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2WebMessageReceivedEventArgs ICoreWebView2WebMessageReceivedEventArgs;

#endif 	/* __ICoreWebView2WebMessageReceivedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebMessageReceivedEventHandler_FWD_DEFINED__
#define __ICoreWebView2WebMessageReceivedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2WebMessageReceivedEventHandler ICoreWebView2WebMessageReceivedEventHandler;

#endif 	/* __ICoreWebView2WebMessageReceivedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceRequest_FWD_DEFINED__
#define __ICoreWebView2WebResourceRequest_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceRequest ICoreWebView2WebResourceRequest;

#endif 	/* __ICoreWebView2WebResourceRequest_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceRequestedEventArgs_FWD_DEFINED__
#define __ICoreWebView2WebResourceRequestedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceRequestedEventArgs ICoreWebView2WebResourceRequestedEventArgs;

#endif 	/* __ICoreWebView2WebResourceRequestedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceRequestedEventHandler_FWD_DEFINED__
#define __ICoreWebView2WebResourceRequestedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceRequestedEventHandler ICoreWebView2WebResourceRequestedEventHandler;

#endif 	/* __ICoreWebView2WebResourceRequestedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponse_FWD_DEFINED__
#define __ICoreWebView2WebResourceResponse_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceResponse ICoreWebView2WebResourceResponse;

#endif 	/* __ICoreWebView2WebResourceResponse_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseReceivedEventHandler_FWD_DEFINED__
#define __ICoreWebView2WebResourceResponseReceivedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceResponseReceivedEventHandler ICoreWebView2WebResourceResponseReceivedEventHandler;

#endif 	/* __ICoreWebView2WebResourceResponseReceivedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseReceivedEventArgs_FWD_DEFINED__
#define __ICoreWebView2WebResourceResponseReceivedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceResponseReceivedEventArgs ICoreWebView2WebResourceResponseReceivedEventArgs;

#endif 	/* __ICoreWebView2WebResourceResponseReceivedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseView_FWD_DEFINED__
#define __ICoreWebView2WebResourceResponseView_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceResponseView ICoreWebView2WebResourceResponseView;

#endif 	/* __ICoreWebView2WebResourceResponseView_FWD_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2WebResourceResponseViewGetContentCompletedHandler ICoreWebView2WebResourceResponseViewGetContentCompletedHandler;

#endif 	/* __ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2WindowCloseRequestedEventHandler_FWD_DEFINED__
#define __ICoreWebView2WindowCloseRequestedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2WindowCloseRequestedEventHandler ICoreWebView2WindowCloseRequestedEventHandler;

#endif 	/* __ICoreWebView2WindowCloseRequestedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2WindowFeatures_FWD_DEFINED__
#define __ICoreWebView2WindowFeatures_FWD_DEFINED__
typedef interface ICoreWebView2WindowFeatures ICoreWebView2WindowFeatures;

#endif 	/* __ICoreWebView2WindowFeatures_FWD_DEFINED__ */


#ifndef __ICoreWebView2ZoomFactorChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2ZoomFactorChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ZoomFactorChangedEventHandler ICoreWebView2ZoomFactorChangedEventHandler;

#endif 	/* __ICoreWebView2ZoomFactorChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CompositionControllerInterop_FWD_DEFINED__
#define __ICoreWebView2CompositionControllerInterop_FWD_DEFINED__
typedef interface ICoreWebView2CompositionControllerInterop ICoreWebView2CompositionControllerInterop;

#endif 	/* __ICoreWebView2CompositionControllerInterop_FWD_DEFINED__ */


#ifndef __ICoreWebView2EnvironmentInterop_FWD_DEFINED__
#define __ICoreWebView2EnvironmentInterop_FWD_DEFINED__
typedef interface ICoreWebView2EnvironmentInterop ICoreWebView2EnvironmentInterop;

#endif 	/* __ICoreWebView2EnvironmentInterop_FWD_DEFINED__ */


/* header files for imported files */
#include "objidl.h"
#include "oaidl.h"
#include "EventToken.h"

#ifdef __cplusplus
extern "C"{
#endif 



#ifndef __WebView2_LIBRARY_DEFINED__
#define __WebView2_LIBRARY_DEFINED__

/* library WebView2 */
/* [version][uuid] */ 




















































































typedef /* [v1_enum] */ 
enum COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT
    {
        COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT_PNG	= 0,
        COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT_JPEG	= ( COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT_PNG + 1 ) 
    } 	COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_COOKIE_SAME_SITE_KIND
    {
        COREWEBVIEW2_COOKIE_SAME_SITE_KIND_NONE	= 0,
        COREWEBVIEW2_COOKIE_SAME_SITE_KIND_LAX	= ( COREWEBVIEW2_COOKIE_SAME_SITE_KIND_NONE + 1 ) ,
        COREWEBVIEW2_COOKIE_SAME_SITE_KIND_STRICT	= ( COREWEBVIEW2_COOKIE_SAME_SITE_KIND_LAX + 1 ) 
    } 	COREWEBVIEW2_COOKIE_SAME_SITE_KIND;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND
    {
        COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_DENY	= 0,
        COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_ALLOW	= ( COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_DENY + 1 ) ,
        COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_DENY_CORS	= ( COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_ALLOW + 1 ) 
    } 	COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_SCRIPT_DIALOG_KIND
    {
        COREWEBVIEW2_SCRIPT_DIALOG_KIND_ALERT	= 0,
        COREWEBVIEW2_SCRIPT_DIALOG_KIND_CONFIRM	= ( COREWEBVIEW2_SCRIPT_DIALOG_KIND_ALERT + 1 ) ,
        COREWEBVIEW2_SCRIPT_DIALOG_KIND_PROMPT	= ( COREWEBVIEW2_SCRIPT_DIALOG_KIND_CONFIRM + 1 ) ,
        COREWEBVIEW2_SCRIPT_DIALOG_KIND_BEFOREUNLOAD	= ( COREWEBVIEW2_SCRIPT_DIALOG_KIND_PROMPT + 1 ) 
    } 	COREWEBVIEW2_SCRIPT_DIALOG_KIND;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_PROCESS_FAILED_KIND
    {
        COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED	= 0,
        COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_BROWSER_PROCESS_EXITED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_UNRESPONSIVE	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_FRAME_RENDER_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_UNRESPONSIVE + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_UTILITY_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_FRAME_RENDER_PROCESS_EXITED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_SANDBOX_HELPER_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_UTILITY_PROCESS_EXITED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_GPU_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_SANDBOX_HELPER_PROCESS_EXITED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_PPAPI_PLUGIN_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_GPU_PROCESS_EXITED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_PPAPI_BROKER_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_PPAPI_PLUGIN_PROCESS_EXITED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_KIND_UNKNOWN_PROCESS_EXITED	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_PPAPI_BROKER_PROCESS_EXITED + 1 ) 
    } 	COREWEBVIEW2_PROCESS_FAILED_KIND;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_PROCESS_FAILED_REASON
    {
        COREWEBVIEW2_PROCESS_FAILED_REASON_UNEXPECTED	= 0,
        COREWEBVIEW2_PROCESS_FAILED_REASON_UNRESPONSIVE	= ( COREWEBVIEW2_PROCESS_FAILED_REASON_UNEXPECTED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_REASON_TERMINATED	= ( COREWEBVIEW2_PROCESS_FAILED_REASON_UNRESPONSIVE + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_REASON_CRASHED	= ( COREWEBVIEW2_PROCESS_FAILED_REASON_TERMINATED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_REASON_LAUNCH_FAILED	= ( COREWEBVIEW2_PROCESS_FAILED_REASON_CRASHED + 1 ) ,
        COREWEBVIEW2_PROCESS_FAILED_REASON_OUT_OF_MEMORY	= ( COREWEBVIEW2_PROCESS_FAILED_REASON_LAUNCH_FAILED + 1 ) 
    } 	COREWEBVIEW2_PROCESS_FAILED_REASON;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_PERMISSION_KIND
    {
        COREWEBVIEW2_PERMISSION_KIND_UNKNOWN_PERMISSION	= 0,
        COREWEBVIEW2_PERMISSION_KIND_MICROPHONE	= ( COREWEBVIEW2_PERMISSION_KIND_UNKNOWN_PERMISSION + 1 ) ,
        COREWEBVIEW2_PERMISSION_KIND_CAMERA	= ( COREWEBVIEW2_PERMISSION_KIND_MICROPHONE + 1 ) ,
        COREWEBVIEW2_PERMISSION_KIND_GEOLOCATION	= ( COREWEBVIEW2_PERMISSION_KIND_CAMERA + 1 ) ,
        COREWEBVIEW2_PERMISSION_KIND_NOTIFICATIONS	= ( COREWEBVIEW2_PERMISSION_KIND_GEOLOCATION + 1 ) ,
        COREWEBVIEW2_PERMISSION_KIND_OTHER_SENSORS	= ( COREWEBVIEW2_PERMISSION_KIND_NOTIFICATIONS + 1 ) ,
        COREWEBVIEW2_PERMISSION_KIND_CLIPBOARD_READ	= ( COREWEBVIEW2_PERMISSION_KIND_OTHER_SENSORS + 1 ) 
    } 	COREWEBVIEW2_PERMISSION_KIND;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_PERMISSION_STATE
    {
        COREWEBVIEW2_PERMISSION_STATE_DEFAULT	= 0,
        COREWEBVIEW2_PERMISSION_STATE_ALLOW	= ( COREWEBVIEW2_PERMISSION_STATE_DEFAULT + 1 ) ,
        COREWEBVIEW2_PERMISSION_STATE_DENY	= ( COREWEBVIEW2_PERMISSION_STATE_ALLOW + 1 ) 
    } 	COREWEBVIEW2_PERMISSION_STATE;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_WEB_ERROR_STATUS
    {
        COREWEBVIEW2_WEB_ERROR_STATUS_UNKNOWN	= 0,
        COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_COMMON_NAME_IS_INCORRECT	= ( COREWEBVIEW2_WEB_ERROR_STATUS_UNKNOWN + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_EXPIRED	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_COMMON_NAME_IS_INCORRECT + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_CLIENT_CERTIFICATE_CONTAINS_ERRORS	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_EXPIRED + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_REVOKED	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CLIENT_CERTIFICATE_CONTAINS_ERRORS + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_IS_INVALID	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_REVOKED + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_SERVER_UNREACHABLE	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CERTIFICATE_IS_INVALID + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_TIMEOUT	= ( COREWEBVIEW2_WEB_ERROR_STATUS_SERVER_UNREACHABLE + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_ERROR_HTTP_INVALID_SERVER_RESPONSE	= ( COREWEBVIEW2_WEB_ERROR_STATUS_TIMEOUT + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_CONNECTION_ABORTED	= ( COREWEBVIEW2_WEB_ERROR_STATUS_ERROR_HTTP_INVALID_SERVER_RESPONSE + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_CONNECTION_RESET	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CONNECTION_ABORTED + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_DISCONNECTED	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CONNECTION_RESET + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_CANNOT_CONNECT	= ( COREWEBVIEW2_WEB_ERROR_STATUS_DISCONNECTED + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_HOST_NAME_NOT_RESOLVED	= ( COREWEBVIEW2_WEB_ERROR_STATUS_CANNOT_CONNECT + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_OPERATION_CANCELED	= ( COREWEBVIEW2_WEB_ERROR_STATUS_HOST_NAME_NOT_RESOLVED + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_REDIRECT_FAILED	= ( COREWEBVIEW2_WEB_ERROR_STATUS_OPERATION_CANCELED + 1 ) ,
        COREWEBVIEW2_WEB_ERROR_STATUS_UNEXPECTED_ERROR	= ( COREWEBVIEW2_WEB_ERROR_STATUS_REDIRECT_FAILED + 1 ) 
    } 	COREWEBVIEW2_WEB_ERROR_STATUS;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_WEB_RESOURCE_CONTEXT
    {
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL	= 0,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_DOCUMENT	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_STYLESHEET	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_DOCUMENT + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_IMAGE	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_STYLESHEET + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_MEDIA	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_IMAGE + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_FONT	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_MEDIA + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_SCRIPT	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_FONT + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_XML_HTTP_REQUEST	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_SCRIPT + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_FETCH	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_XML_HTTP_REQUEST + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_TEXT_TRACK	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_FETCH + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_EVENT_SOURCE	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_TEXT_TRACK + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_WEBSOCKET	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_EVENT_SOURCE + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_MANIFEST	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_WEBSOCKET + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_SIGNED_EXCHANGE	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_MANIFEST + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_PING	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_SIGNED_EXCHANGE + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_CSP_VIOLATION_REPORT	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_PING + 1 ) ,
        COREWEBVIEW2_WEB_RESOURCE_CONTEXT_OTHER	= ( COREWEBVIEW2_WEB_RESOURCE_CONTEXT_CSP_VIOLATION_REPORT + 1 ) 
    } 	COREWEBVIEW2_WEB_RESOURCE_CONTEXT;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_MOVE_FOCUS_REASON
    {
        COREWEBVIEW2_MOVE_FOCUS_REASON_PROGRAMMATIC	= 0,
        COREWEBVIEW2_MOVE_FOCUS_REASON_NEXT	= ( COREWEBVIEW2_MOVE_FOCUS_REASON_PROGRAMMATIC + 1 ) ,
        COREWEBVIEW2_MOVE_FOCUS_REASON_PREVIOUS	= ( COREWEBVIEW2_MOVE_FOCUS_REASON_NEXT + 1 ) 
    } 	COREWEBVIEW2_MOVE_FOCUS_REASON;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_KEY_EVENT_KIND
    {
        COREWEBVIEW2_KEY_EVENT_KIND_KEY_DOWN	= 0,
        COREWEBVIEW2_KEY_EVENT_KIND_KEY_UP	= ( COREWEBVIEW2_KEY_EVENT_KIND_KEY_DOWN + 1 ) ,
        COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN	= ( COREWEBVIEW2_KEY_EVENT_KIND_KEY_UP + 1 ) ,
        COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_UP	= ( COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN + 1 ) 
    } 	COREWEBVIEW2_KEY_EVENT_KIND;

typedef struct COREWEBVIEW2_PHYSICAL_KEY_STATUS
    {
    UINT32 RepeatCount;
    UINT32 ScanCode;
    BOOL IsExtendedKey;
    BOOL IsMenuKeyDown;
    BOOL WasKeyDown;
    BOOL IsKeyReleased;
    } 	COREWEBVIEW2_PHYSICAL_KEY_STATUS;

typedef struct COREWEBVIEW2_COLOR
    {
    BYTE A;
    BYTE R;
    BYTE G;
    BYTE B;
    } 	COREWEBVIEW2_COLOR;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_MOUSE_EVENT_KIND
    {
        COREWEBVIEW2_MOUSE_EVENT_KIND_HORIZONTAL_WHEEL	= 0x20e,
        COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_DOUBLE_CLICK	= 0x203,
        COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_DOWN	= 0x201,
        COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_UP	= 0x202,
        COREWEBVIEW2_MOUSE_EVENT_KIND_LEAVE	= 0x2a3,
        COREWEBVIEW2_MOUSE_EVENT_KIND_MIDDLE_BUTTON_DOUBLE_CLICK	= 0x209,
        COREWEBVIEW2_MOUSE_EVENT_KIND_MIDDLE_BUTTON_DOWN	= 0x207,
        COREWEBVIEW2_MOUSE_EVENT_KIND_MIDDLE_BUTTON_UP	= 0x208,
        COREWEBVIEW2_MOUSE_EVENT_KIND_MOVE	= 0x200,
        COREWEBVIEW2_MOUSE_EVENT_KIND_RIGHT_BUTTON_DOUBLE_CLICK	= 0x206,
        COREWEBVIEW2_MOUSE_EVENT_KIND_RIGHT_BUTTON_DOWN	= 0x204,
        COREWEBVIEW2_MOUSE_EVENT_KIND_RIGHT_BUTTON_UP	= 0x205,
        COREWEBVIEW2_MOUSE_EVENT_KIND_WHEEL	= 0x20a,
        COREWEBVIEW2_MOUSE_EVENT_KIND_X_BUTTON_DOUBLE_CLICK	= 0x20d,
        COREWEBVIEW2_MOUSE_EVENT_KIND_X_BUTTON_DOWN	= 0x20b,
        COREWEBVIEW2_MOUSE_EVENT_KIND_X_BUTTON_UP	= 0x20c
    } 	COREWEBVIEW2_MOUSE_EVENT_KIND;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS
    {
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_NONE	= 0,
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_LEFT_BUTTON	= 0x1,
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_RIGHT_BUTTON	= 0x2,
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_SHIFT	= 0x4,
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_CONTROL	= 0x8,
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_MIDDLE_BUTTON	= 0x10,
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_X_BUTTON1	= 0x20,
        COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_X_BUTTON2	= 0x40
    } 	COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS;

DEFINE_ENUM_FLAG_OPERATORS(COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS);
typedef /* [v1_enum] */ 
enum COREWEBVIEW2_POINTER_EVENT_KIND
    {
        COREWEBVIEW2_POINTER_EVENT_KIND_ACTIVATE	= 0x24b,
        COREWEBVIEW2_POINTER_EVENT_KIND_DOWN	= 0x246,
        COREWEBVIEW2_POINTER_EVENT_KIND_ENTER	= 0x249,
        COREWEBVIEW2_POINTER_EVENT_KIND_LEAVE	= 0x24a,
        COREWEBVIEW2_POINTER_EVENT_KIND_UP	= 0x247,
        COREWEBVIEW2_POINTER_EVENT_KIND_UPDATE	= 0x245
    } 	COREWEBVIEW2_POINTER_EVENT_KIND;

typedef /* [v1_enum] */ 
enum COREWEBVIEW2_BOUNDS_MODE
    {
        COREWEBVIEW2_BOUNDS_MODE_USE_RAW_PIXELS	= 0,
        COREWEBVIEW2_BOUNDS_MODE_USE_RASTERIZATION_SCALE	= ( COREWEBVIEW2_BOUNDS_MODE_USE_RAW_PIXELS + 1 ) 
    } 	COREWEBVIEW2_BOUNDS_MODE;

STDAPI CreateCoreWebView2EnvironmentWithOptions(PCWSTR browserExecutableFolder, PCWSTR userDataFolder, ICoreWebView2EnvironmentOptions* environmentOptions, ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler* environmentCreatedHandler);
STDAPI CreateCoreWebView2Environment(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler* environmentCreatedHandler);
STDAPI GetAvailableCoreWebView2BrowserVersionString(PCWSTR browserExecutableFolder, LPWSTR* versionInfo);
STDAPI CompareBrowserVersions(PCWSTR version1, PCWSTR version2, int* result);

EXTERN_C const IID LIBID_WebView2;

#ifndef __ICoreWebView2AcceleratorKeyPressedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2AcceleratorKeyPressedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2AcceleratorKeyPressedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2AcceleratorKeyPressedEventArgs = {0x9f760f8a,0xfb79,0x42be,{0x99,0x90,0x7b,0x56,0x90,0x0f,0xa9,0xc7}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("9f760f8a-fb79-42be-9990-7b56900fa9c7")
    ICoreWebView2AcceleratorKeyPressedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_KeyEventKind( 
            /* [retval][out] */ COREWEBVIEW2_KEY_EVENT_KIND *keyEventKind) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_VirtualKey( 
            /* [retval][out] */ UINT *virtualKey) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_KeyEventLParam( 
            /* [retval][out] */ INT *lParam) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PhysicalKeyStatus( 
            /* [retval][out] */ COREWEBVIEW2_PHYSICAL_KEY_STATUS *physicalKeyStatus) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Handled( 
            /* [retval][out] */ BOOL *handled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Handled( 
            /* [in] */ BOOL handled) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2AcceleratorKeyPressedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_KeyEventKind )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_KEY_EVENT_KIND *keyEventKind);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_VirtualKey )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This,
            /* [retval][out] */ UINT *virtualKey);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_KeyEventLParam )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This,
            /* [retval][out] */ INT *lParam);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PhysicalKeyStatus )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_PHYSICAL_KEY_STATUS *physicalKeyStatus);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Handled )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This,
            /* [retval][out] */ BOOL *handled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Handled )( 
            ICoreWebView2AcceleratorKeyPressedEventArgs * This,
            /* [in] */ BOOL handled);
        
        END_INTERFACE
    } ICoreWebView2AcceleratorKeyPressedEventArgsVtbl;

    interface ICoreWebView2AcceleratorKeyPressedEventArgs
    {
        CONST_VTBL struct ICoreWebView2AcceleratorKeyPressedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2AcceleratorKeyPressedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2AcceleratorKeyPressedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2AcceleratorKeyPressedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2AcceleratorKeyPressedEventArgs_get_KeyEventKind(This,keyEventKind)	\
    ( (This)->lpVtbl -> get_KeyEventKind(This,keyEventKind) ) 

#define ICoreWebView2AcceleratorKeyPressedEventArgs_get_VirtualKey(This,virtualKey)	\
    ( (This)->lpVtbl -> get_VirtualKey(This,virtualKey) ) 

#define ICoreWebView2AcceleratorKeyPressedEventArgs_get_KeyEventLParam(This,lParam)	\
    ( (This)->lpVtbl -> get_KeyEventLParam(This,lParam) ) 

#define ICoreWebView2AcceleratorKeyPressedEventArgs_get_PhysicalKeyStatus(This,physicalKeyStatus)	\
    ( (This)->lpVtbl -> get_PhysicalKeyStatus(This,physicalKeyStatus) ) 

#define ICoreWebView2AcceleratorKeyPressedEventArgs_get_Handled(This,handled)	\
    ( (This)->lpVtbl -> get_Handled(This,handled) ) 

#define ICoreWebView2AcceleratorKeyPressedEventArgs_put_Handled(This,handled)	\
    ( (This)->lpVtbl -> put_Handled(This,handled) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2AcceleratorKeyPressedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2AcceleratorKeyPressedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2AcceleratorKeyPressedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2AcceleratorKeyPressedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2AcceleratorKeyPressedEventHandler = {0xb29c7e28,0xfa79,0x41a8,{0x8e,0x44,0x65,0x81,0x1c,0x76,0xdc,0xb2}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("b29c7e28-fa79-41a8-8e44-65811c76dcb2")
    ICoreWebView2AcceleratorKeyPressedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ ICoreWebView2AcceleratorKeyPressedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2AcceleratorKeyPressedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2AcceleratorKeyPressedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2AcceleratorKeyPressedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2AcceleratorKeyPressedEventHandler * This,
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ ICoreWebView2AcceleratorKeyPressedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl;

    interface ICoreWebView2AcceleratorKeyPressedEventHandler
    {
        CONST_VTBL struct ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2AcceleratorKeyPressedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2AcceleratorKeyPressedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2AcceleratorKeyPressedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2AcceleratorKeyPressedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2AcceleratorKeyPressedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler = {0xb99369f3,0x9b11,0x47b5,{0xbc,0x6f,0x8e,0x78,0x95,0xfc,0xea,0x17}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("b99369f3-9b11-47b5-bc6f-8e7895fcea17")
    ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ HRESULT errorCode,
            /* [in] */ LPCWSTR id) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler * This,
            /* [in] */ HRESULT errorCode,
            /* [in] */ LPCWSTR id);
        
        END_INTERFACE
    } ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl;

    interface ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_Invoke(This,errorCode,id)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,id) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CallDevToolsProtocolMethodCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CallDevToolsProtocolMethodCompletedHandler = {0x5c4889f0,0x5ef6,0x4c5a,{0x95,0x2c,0xd8,0xf1,0xb9,0x2d,0x05,0x74}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("5c4889f0-5ef6-4c5a-952c-d8f1b92d0574")
    ICoreWebView2CallDevToolsProtocolMethodCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ HRESULT errorCode,
            /* [in] */ LPCWSTR returnObjectAsJson) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CallDevToolsProtocolMethodCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CallDevToolsProtocolMethodCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CallDevToolsProtocolMethodCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2CallDevToolsProtocolMethodCompletedHandler * This,
            /* [in] */ HRESULT errorCode,
            /* [in] */ LPCWSTR returnObjectAsJson);
        
        END_INTERFACE
    } ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl;

    interface ICoreWebView2CallDevToolsProtocolMethodCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_Invoke(This,errorCode,returnObjectAsJson)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,returnObjectAsJson) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CallDevToolsProtocolMethodCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CapturePreviewCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CapturePreviewCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CapturePreviewCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CapturePreviewCompletedHandler = {0x697e05e9,0x3d8f,0x45fa,{0x96,0xf4,0x8f,0xfe,0x1e,0xde,0xda,0xf5}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("697e05e9-3d8f-45fa-96f4-8ffe1ededaf5")
    ICoreWebView2CapturePreviewCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ HRESULT errorCode) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CapturePreviewCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CapturePreviewCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CapturePreviewCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CapturePreviewCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2CapturePreviewCompletedHandler * This,
            /* [in] */ HRESULT errorCode);
        
        END_INTERFACE
    } ICoreWebView2CapturePreviewCompletedHandlerVtbl;

    interface ICoreWebView2CapturePreviewCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2CapturePreviewCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CapturePreviewCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CapturePreviewCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CapturePreviewCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CapturePreviewCompletedHandler_Invoke(This,errorCode)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CapturePreviewCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2_INTERFACE_DEFINED__
#define __ICoreWebView2_INTERFACE_DEFINED__

/* interface ICoreWebView2 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2 = {0x76eceacb,0x0462,0x4d94,{0xac,0x83,0x42,0x3a,0x67,0x93,0x77,0x5e}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("76eceacb-0462-4d94-ac83-423a6793775e")
    ICoreWebView2 : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Settings( 
            /* [retval][out] */ ICoreWebView2Settings **settings) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Source( 
            /* [retval][out] */ LPWSTR *uri) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Navigate( 
            /* [in] */ LPCWSTR uri) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE NavigateToString( 
            /* [in] */ LPCWSTR htmlContent) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_NavigationStarting( 
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_NavigationStarting( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_ContentLoading( 
            /* [in] */ ICoreWebView2ContentLoadingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_ContentLoading( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_SourceChanged( 
            /* [in] */ ICoreWebView2SourceChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_SourceChanged( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_HistoryChanged( 
            /* [in] */ ICoreWebView2HistoryChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_HistoryChanged( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_NavigationCompleted( 
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_NavigationCompleted( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_FrameNavigationStarting( 
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_FrameNavigationStarting( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_FrameNavigationCompleted( 
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_FrameNavigationCompleted( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_ScriptDialogOpening( 
            /* [in] */ ICoreWebView2ScriptDialogOpeningEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_ScriptDialogOpening( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_PermissionRequested( 
            /* [in] */ ICoreWebView2PermissionRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_PermissionRequested( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_ProcessFailed( 
            /* [in] */ ICoreWebView2ProcessFailedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_ProcessFailed( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE AddScriptToExecuteOnDocumentCreated( 
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler *handler) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE RemoveScriptToExecuteOnDocumentCreated( 
            /* [in] */ LPCWSTR id) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE ExecuteScript( 
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2ExecuteScriptCompletedHandler *handler) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE CapturePreview( 
            /* [in] */ COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT imageFormat,
            /* [in] */ IStream *imageStream,
            /* [in] */ ICoreWebView2CapturePreviewCompletedHandler *handler) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Reload( void) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE PostWebMessageAsJson( 
            /* [in] */ LPCWSTR webMessageAsJson) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE PostWebMessageAsString( 
            /* [in] */ LPCWSTR webMessageAsString) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_WebMessageReceived( 
            /* [in] */ ICoreWebView2WebMessageReceivedEventHandler *handler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_WebMessageReceived( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE CallDevToolsProtocolMethod( 
            /* [in] */ LPCWSTR methodName,
            /* [in] */ LPCWSTR parametersAsJson,
            /* [in] */ ICoreWebView2CallDevToolsProtocolMethodCompletedHandler *handler) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_BrowserProcessId( 
            /* [retval][out] */ UINT32 *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_CanGoBack( 
            /* [retval][out] */ BOOL *canGoBack) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_CanGoForward( 
            /* [retval][out] */ BOOL *canGoForward) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GoBack( void) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GoForward( void) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetDevToolsProtocolEventReceiver( 
            /* [in] */ LPCWSTR eventName,
            /* [retval][out] */ ICoreWebView2DevToolsProtocolEventReceiver **receiver) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Stop( void) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_NewWindowRequested( 
            /* [in] */ ICoreWebView2NewWindowRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_NewWindowRequested( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_DocumentTitleChanged( 
            /* [in] */ ICoreWebView2DocumentTitleChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_DocumentTitleChanged( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_DocumentTitle( 
            /* [retval][out] */ LPWSTR *title) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE AddHostObjectToScript( 
            /* [in] */ LPCWSTR name,
            /* [in] */ VARIANT *object) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE RemoveHostObjectFromScript( 
            /* [in] */ LPCWSTR name) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE OpenDevToolsWindow( void) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_ContainsFullScreenElementChanged( 
            /* [in] */ ICoreWebView2ContainsFullScreenElementChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_ContainsFullScreenElementChanged( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ContainsFullScreenElement( 
            /* [retval][out] */ BOOL *containsFullScreenElement) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_WebResourceRequested( 
            /* [in] */ ICoreWebView2WebResourceRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_WebResourceRequested( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE AddWebResourceRequestedFilter( 
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE RemoveWebResourceRequestedFilter( 
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_WindowCloseRequested( 
            /* [in] */ ICoreWebView2WindowCloseRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_WindowCloseRequested( 
            /* [in] */ EventRegistrationToken token) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Settings )( 
            ICoreWebView2 * This,
            /* [retval][out] */ ICoreWebView2Settings **settings);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Source )( 
            ICoreWebView2 * This,
            /* [retval][out] */ LPWSTR *uri);
        
        HRESULT ( STDMETHODCALLTYPE *Navigate )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR uri);
        
        HRESULT ( STDMETHODCALLTYPE *NavigateToString )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR htmlContent);
        
        HRESULT ( STDMETHODCALLTYPE *add_NavigationStarting )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NavigationStarting )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ContentLoading )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2ContentLoadingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ContentLoading )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_SourceChanged )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2SourceChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_SourceChanged )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_HistoryChanged )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2HistoryChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_HistoryChanged )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_NavigationCompleted )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NavigationCompleted )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_FrameNavigationStarting )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_FrameNavigationStarting )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_FrameNavigationCompleted )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_FrameNavigationCompleted )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ScriptDialogOpening )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2ScriptDialogOpeningEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ScriptDialogOpening )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_PermissionRequested )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2PermissionRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_PermissionRequested )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ProcessFailed )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2ProcessFailedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ProcessFailed )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *AddScriptToExecuteOnDocumentCreated )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveScriptToExecuteOnDocumentCreated )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR id);
        
        HRESULT ( STDMETHODCALLTYPE *ExecuteScript )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2ExecuteScriptCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CapturePreview )( 
            ICoreWebView2 * This,
            /* [in] */ COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT imageFormat,
            /* [in] */ IStream *imageStream,
            /* [in] */ ICoreWebView2CapturePreviewCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *Reload )( 
            ICoreWebView2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *PostWebMessageAsJson )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR webMessageAsJson);
        
        HRESULT ( STDMETHODCALLTYPE *PostWebMessageAsString )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR webMessageAsString);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebMessageReceived )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2WebMessageReceivedEventHandler *handler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebMessageReceived )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *CallDevToolsProtocolMethod )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR methodName,
            /* [in] */ LPCWSTR parametersAsJson,
            /* [in] */ ICoreWebView2CallDevToolsProtocolMethodCompletedHandler *handler);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BrowserProcessId )( 
            ICoreWebView2 * This,
            /* [retval][out] */ UINT32 *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CanGoBack )( 
            ICoreWebView2 * This,
            /* [retval][out] */ BOOL *canGoBack);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CanGoForward )( 
            ICoreWebView2 * This,
            /* [retval][out] */ BOOL *canGoForward);
        
        HRESULT ( STDMETHODCALLTYPE *GoBack )( 
            ICoreWebView2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *GoForward )( 
            ICoreWebView2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *GetDevToolsProtocolEventReceiver )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR eventName,
            /* [retval][out] */ ICoreWebView2DevToolsProtocolEventReceiver **receiver);
        
        HRESULT ( STDMETHODCALLTYPE *Stop )( 
            ICoreWebView2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *add_NewWindowRequested )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2NewWindowRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NewWindowRequested )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_DocumentTitleChanged )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2DocumentTitleChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_DocumentTitleChanged )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_DocumentTitle )( 
            ICoreWebView2 * This,
            /* [retval][out] */ LPWSTR *title);
        
        HRESULT ( STDMETHODCALLTYPE *AddHostObjectToScript )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ VARIANT *object);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveHostObjectFromScript )( 
            ICoreWebView2 * This,
            /* [in] */ LPCWSTR name);
        
        HRESULT ( STDMETHODCALLTYPE *OpenDevToolsWindow )( 
            ICoreWebView2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *add_ContainsFullScreenElementChanged )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2ContainsFullScreenElementChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ContainsFullScreenElementChanged )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ContainsFullScreenElement )( 
            ICoreWebView2 * This,
            /* [retval][out] */ BOOL *containsFullScreenElement);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebResourceRequested )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2WebResourceRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebResourceRequested )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *AddWebResourceRequestedFilter )( 
            ICoreWebView2 * This,
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveWebResourceRequestedFilter )( 
            ICoreWebView2 * This,
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext);
        
        HRESULT ( STDMETHODCALLTYPE *add_WindowCloseRequested )( 
            ICoreWebView2 * This,
            /* [in] */ ICoreWebView2WindowCloseRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WindowCloseRequested )( 
            ICoreWebView2 * This,
            /* [in] */ EventRegistrationToken token);
        
        END_INTERFACE
    } ICoreWebView2Vtbl;

    interface ICoreWebView2
    {
        CONST_VTBL struct ICoreWebView2Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2_get_Settings(This,settings)	\
    ( (This)->lpVtbl -> get_Settings(This,settings) ) 

#define ICoreWebView2_get_Source(This,uri)	\
    ( (This)->lpVtbl -> get_Source(This,uri) ) 

#define ICoreWebView2_Navigate(This,uri)	\
    ( (This)->lpVtbl -> Navigate(This,uri) ) 

#define ICoreWebView2_NavigateToString(This,htmlContent)	\
    ( (This)->lpVtbl -> NavigateToString(This,htmlContent) ) 

#define ICoreWebView2_add_NavigationStarting(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NavigationStarting(This,eventHandler,token) ) 

#define ICoreWebView2_remove_NavigationStarting(This,token)	\
    ( (This)->lpVtbl -> remove_NavigationStarting(This,token) ) 

#define ICoreWebView2_add_ContentLoading(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ContentLoading(This,eventHandler,token) ) 

#define ICoreWebView2_remove_ContentLoading(This,token)	\
    ( (This)->lpVtbl -> remove_ContentLoading(This,token) ) 

#define ICoreWebView2_add_SourceChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_SourceChanged(This,eventHandler,token) ) 

#define ICoreWebView2_remove_SourceChanged(This,token)	\
    ( (This)->lpVtbl -> remove_SourceChanged(This,token) ) 

#define ICoreWebView2_add_HistoryChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_HistoryChanged(This,eventHandler,token) ) 

#define ICoreWebView2_remove_HistoryChanged(This,token)	\
    ( (This)->lpVtbl -> remove_HistoryChanged(This,token) ) 

#define ICoreWebView2_add_NavigationCompleted(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NavigationCompleted(This,eventHandler,token) ) 

#define ICoreWebView2_remove_NavigationCompleted(This,token)	\
    ( (This)->lpVtbl -> remove_NavigationCompleted(This,token) ) 

#define ICoreWebView2_add_FrameNavigationStarting(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_FrameNavigationStarting(This,eventHandler,token) ) 

#define ICoreWebView2_remove_FrameNavigationStarting(This,token)	\
    ( (This)->lpVtbl -> remove_FrameNavigationStarting(This,token) ) 

#define ICoreWebView2_add_FrameNavigationCompleted(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_FrameNavigationCompleted(This,eventHandler,token) ) 

#define ICoreWebView2_remove_FrameNavigationCompleted(This,token)	\
    ( (This)->lpVtbl -> remove_FrameNavigationCompleted(This,token) ) 

#define ICoreWebView2_add_ScriptDialogOpening(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ScriptDialogOpening(This,eventHandler,token) ) 

#define ICoreWebView2_remove_ScriptDialogOpening(This,token)	\
    ( (This)->lpVtbl -> remove_ScriptDialogOpening(This,token) ) 

#define ICoreWebView2_add_PermissionRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_PermissionRequested(This,eventHandler,token) ) 

#define ICoreWebView2_remove_PermissionRequested(This,token)	\
    ( (This)->lpVtbl -> remove_PermissionRequested(This,token) ) 

#define ICoreWebView2_add_ProcessFailed(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ProcessFailed(This,eventHandler,token) ) 

#define ICoreWebView2_remove_ProcessFailed(This,token)	\
    ( (This)->lpVtbl -> remove_ProcessFailed(This,token) ) 

#define ICoreWebView2_AddScriptToExecuteOnDocumentCreated(This,javaScript,handler)	\
    ( (This)->lpVtbl -> AddScriptToExecuteOnDocumentCreated(This,javaScript,handler) ) 

#define ICoreWebView2_RemoveScriptToExecuteOnDocumentCreated(This,id)	\
    ( (This)->lpVtbl -> RemoveScriptToExecuteOnDocumentCreated(This,id) ) 

#define ICoreWebView2_ExecuteScript(This,javaScript,handler)	\
    ( (This)->lpVtbl -> ExecuteScript(This,javaScript,handler) ) 

#define ICoreWebView2_CapturePreview(This,imageFormat,imageStream,handler)	\
    ( (This)->lpVtbl -> CapturePreview(This,imageFormat,imageStream,handler) ) 

#define ICoreWebView2_Reload(This)	\
    ( (This)->lpVtbl -> Reload(This) ) 

#define ICoreWebView2_PostWebMessageAsJson(This,webMessageAsJson)	\
    ( (This)->lpVtbl -> PostWebMessageAsJson(This,webMessageAsJson) ) 

#define ICoreWebView2_PostWebMessageAsString(This,webMessageAsString)	\
    ( (This)->lpVtbl -> PostWebMessageAsString(This,webMessageAsString) ) 

#define ICoreWebView2_add_WebMessageReceived(This,handler,token)	\
    ( (This)->lpVtbl -> add_WebMessageReceived(This,handler,token) ) 

#define ICoreWebView2_remove_WebMessageReceived(This,token)	\
    ( (This)->lpVtbl -> remove_WebMessageReceived(This,token) ) 

#define ICoreWebView2_CallDevToolsProtocolMethod(This,methodName,parametersAsJson,handler)	\
    ( (This)->lpVtbl -> CallDevToolsProtocolMethod(This,methodName,parametersAsJson,handler) ) 

#define ICoreWebView2_get_BrowserProcessId(This,value)	\
    ( (This)->lpVtbl -> get_BrowserProcessId(This,value) ) 

#define ICoreWebView2_get_CanGoBack(This,canGoBack)	\
    ( (This)->lpVtbl -> get_CanGoBack(This,canGoBack) ) 

#define ICoreWebView2_get_CanGoForward(This,canGoForward)	\
    ( (This)->lpVtbl -> get_CanGoForward(This,canGoForward) ) 

#define ICoreWebView2_GoBack(This)	\
    ( (This)->lpVtbl -> GoBack(This) ) 

#define ICoreWebView2_GoForward(This)	\
    ( (This)->lpVtbl -> GoForward(This) ) 

#define ICoreWebView2_GetDevToolsProtocolEventReceiver(This,eventName,receiver)	\
    ( (This)->lpVtbl -> GetDevToolsProtocolEventReceiver(This,eventName,receiver) ) 

#define ICoreWebView2_Stop(This)	\
    ( (This)->lpVtbl -> Stop(This) ) 

#define ICoreWebView2_add_NewWindowRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NewWindowRequested(This,eventHandler,token) ) 

#define ICoreWebView2_remove_NewWindowRequested(This,token)	\
    ( (This)->lpVtbl -> remove_NewWindowRequested(This,token) ) 

#define ICoreWebView2_add_DocumentTitleChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_DocumentTitleChanged(This,eventHandler,token) ) 

#define ICoreWebView2_remove_DocumentTitleChanged(This,token)	\
    ( (This)->lpVtbl -> remove_DocumentTitleChanged(This,token) ) 

#define ICoreWebView2_get_DocumentTitle(This,title)	\
    ( (This)->lpVtbl -> get_DocumentTitle(This,title) ) 

#define ICoreWebView2_AddHostObjectToScript(This,name,object)	\
    ( (This)->lpVtbl -> AddHostObjectToScript(This,name,object) ) 

#define ICoreWebView2_RemoveHostObjectFromScript(This,name)	\
    ( (This)->lpVtbl -> RemoveHostObjectFromScript(This,name) ) 

#define ICoreWebView2_OpenDevToolsWindow(This)	\
    ( (This)->lpVtbl -> OpenDevToolsWindow(This) ) 

#define ICoreWebView2_add_ContainsFullScreenElementChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ContainsFullScreenElementChanged(This,eventHandler,token) ) 

#define ICoreWebView2_remove_ContainsFullScreenElementChanged(This,token)	\
    ( (This)->lpVtbl -> remove_ContainsFullScreenElementChanged(This,token) ) 

#define ICoreWebView2_get_ContainsFullScreenElement(This,containsFullScreenElement)	\
    ( (This)->lpVtbl -> get_ContainsFullScreenElement(This,containsFullScreenElement) ) 

#define ICoreWebView2_add_WebResourceRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WebResourceRequested(This,eventHandler,token) ) 

#define ICoreWebView2_remove_WebResourceRequested(This,token)	\
    ( (This)->lpVtbl -> remove_WebResourceRequested(This,token) ) 

#define ICoreWebView2_AddWebResourceRequestedFilter(This,uri,resourceContext)	\
    ( (This)->lpVtbl -> AddWebResourceRequestedFilter(This,uri,resourceContext) ) 

#define ICoreWebView2_RemoveWebResourceRequestedFilter(This,uri,resourceContext)	\
    ( (This)->lpVtbl -> RemoveWebResourceRequestedFilter(This,uri,resourceContext) ) 

#define ICoreWebView2_add_WindowCloseRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WindowCloseRequested(This,eventHandler,token) ) 

#define ICoreWebView2_remove_WindowCloseRequested(This,token)	\
    ( (This)->lpVtbl -> remove_WindowCloseRequested(This,token) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2_2_INTERFACE_DEFINED__
#define __ICoreWebView2_2_INTERFACE_DEFINED__

/* interface ICoreWebView2_2 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2_2 = {0x9E8F0CF8,0xE670,0x4B5E,{0xB2,0xBC,0x73,0xE0,0x61,0xE3,0x18,0x4C}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("9E8F0CF8-E670-4B5E-B2BC-73E061E3184C")
    ICoreWebView2_2 : public ICoreWebView2
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE add_WebResourceResponseReceived( 
            /* [in] */ ICoreWebView2WebResourceResponseReceivedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_WebResourceResponseReceived( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE NavigateWithWebResourceRequest( 
            /* [in] */ ICoreWebView2WebResourceRequest *request) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_DOMContentLoaded( 
            /* [in] */ ICoreWebView2DOMContentLoadedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_DOMContentLoaded( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_CookieManager( 
            /* [retval][out] */ ICoreWebView2CookieManager **cookieManager) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Environment( 
            /* [retval][out] */ ICoreWebView2Environment **environment) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2_2Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2_2 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2_2 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2_2 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Settings )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ ICoreWebView2Settings **settings);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Source )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ LPWSTR *uri);
        
        HRESULT ( STDMETHODCALLTYPE *Navigate )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR uri);
        
        HRESULT ( STDMETHODCALLTYPE *NavigateToString )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR htmlContent);
        
        HRESULT ( STDMETHODCALLTYPE *add_NavigationStarting )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NavigationStarting )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ContentLoading )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2ContentLoadingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ContentLoading )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_SourceChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2SourceChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_SourceChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_HistoryChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2HistoryChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_HistoryChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_NavigationCompleted )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NavigationCompleted )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_FrameNavigationStarting )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_FrameNavigationStarting )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_FrameNavigationCompleted )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_FrameNavigationCompleted )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ScriptDialogOpening )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2ScriptDialogOpeningEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ScriptDialogOpening )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_PermissionRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2PermissionRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_PermissionRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ProcessFailed )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2ProcessFailedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ProcessFailed )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *AddScriptToExecuteOnDocumentCreated )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveScriptToExecuteOnDocumentCreated )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR id);
        
        HRESULT ( STDMETHODCALLTYPE *ExecuteScript )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2ExecuteScriptCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CapturePreview )( 
            ICoreWebView2_2 * This,
            /* [in] */ COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT imageFormat,
            /* [in] */ IStream *imageStream,
            /* [in] */ ICoreWebView2CapturePreviewCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *Reload )( 
            ICoreWebView2_2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *PostWebMessageAsJson )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR webMessageAsJson);
        
        HRESULT ( STDMETHODCALLTYPE *PostWebMessageAsString )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR webMessageAsString);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebMessageReceived )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2WebMessageReceivedEventHandler *handler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebMessageReceived )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *CallDevToolsProtocolMethod )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR methodName,
            /* [in] */ LPCWSTR parametersAsJson,
            /* [in] */ ICoreWebView2CallDevToolsProtocolMethodCompletedHandler *handler);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BrowserProcessId )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ UINT32 *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CanGoBack )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ BOOL *canGoBack);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CanGoForward )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ BOOL *canGoForward);
        
        HRESULT ( STDMETHODCALLTYPE *GoBack )( 
            ICoreWebView2_2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *GoForward )( 
            ICoreWebView2_2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *GetDevToolsProtocolEventReceiver )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR eventName,
            /* [retval][out] */ ICoreWebView2DevToolsProtocolEventReceiver **receiver);
        
        HRESULT ( STDMETHODCALLTYPE *Stop )( 
            ICoreWebView2_2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *add_NewWindowRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2NewWindowRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NewWindowRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_DocumentTitleChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2DocumentTitleChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_DocumentTitleChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_DocumentTitle )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ LPWSTR *title);
        
        HRESULT ( STDMETHODCALLTYPE *AddHostObjectToScript )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ VARIANT *object);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveHostObjectFromScript )( 
            ICoreWebView2_2 * This,
            /* [in] */ LPCWSTR name);
        
        HRESULT ( STDMETHODCALLTYPE *OpenDevToolsWindow )( 
            ICoreWebView2_2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *add_ContainsFullScreenElementChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2ContainsFullScreenElementChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ContainsFullScreenElementChanged )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ContainsFullScreenElement )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ BOOL *containsFullScreenElement);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebResourceRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2WebResourceRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebResourceRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *AddWebResourceRequestedFilter )( 
            ICoreWebView2_2 * This,
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveWebResourceRequestedFilter )( 
            ICoreWebView2_2 * This,
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext);
        
        HRESULT ( STDMETHODCALLTYPE *add_WindowCloseRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2WindowCloseRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WindowCloseRequested )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebResourceResponseReceived )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2WebResourceResponseReceivedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebResourceResponseReceived )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *NavigateWithWebResourceRequest )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2WebResourceRequest *request);
        
        HRESULT ( STDMETHODCALLTYPE *add_DOMContentLoaded )( 
            ICoreWebView2_2 * This,
            /* [in] */ ICoreWebView2DOMContentLoadedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_DOMContentLoaded )( 
            ICoreWebView2_2 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CookieManager )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ ICoreWebView2CookieManager **cookieManager);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Environment )( 
            ICoreWebView2_2 * This,
            /* [retval][out] */ ICoreWebView2Environment **environment);
        
        END_INTERFACE
    } ICoreWebView2_2Vtbl;

    interface ICoreWebView2_2
    {
        CONST_VTBL struct ICoreWebView2_2Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2_2_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2_2_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2_2_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2_2_get_Settings(This,settings)	\
    ( (This)->lpVtbl -> get_Settings(This,settings) ) 

#define ICoreWebView2_2_get_Source(This,uri)	\
    ( (This)->lpVtbl -> get_Source(This,uri) ) 

#define ICoreWebView2_2_Navigate(This,uri)	\
    ( (This)->lpVtbl -> Navigate(This,uri) ) 

#define ICoreWebView2_2_NavigateToString(This,htmlContent)	\
    ( (This)->lpVtbl -> NavigateToString(This,htmlContent) ) 

#define ICoreWebView2_2_add_NavigationStarting(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NavigationStarting(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_NavigationStarting(This,token)	\
    ( (This)->lpVtbl -> remove_NavigationStarting(This,token) ) 

#define ICoreWebView2_2_add_ContentLoading(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ContentLoading(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_ContentLoading(This,token)	\
    ( (This)->lpVtbl -> remove_ContentLoading(This,token) ) 

#define ICoreWebView2_2_add_SourceChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_SourceChanged(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_SourceChanged(This,token)	\
    ( (This)->lpVtbl -> remove_SourceChanged(This,token) ) 

#define ICoreWebView2_2_add_HistoryChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_HistoryChanged(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_HistoryChanged(This,token)	\
    ( (This)->lpVtbl -> remove_HistoryChanged(This,token) ) 

#define ICoreWebView2_2_add_NavigationCompleted(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NavigationCompleted(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_NavigationCompleted(This,token)	\
    ( (This)->lpVtbl -> remove_NavigationCompleted(This,token) ) 

#define ICoreWebView2_2_add_FrameNavigationStarting(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_FrameNavigationStarting(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_FrameNavigationStarting(This,token)	\
    ( (This)->lpVtbl -> remove_FrameNavigationStarting(This,token) ) 

#define ICoreWebView2_2_add_FrameNavigationCompleted(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_FrameNavigationCompleted(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_FrameNavigationCompleted(This,token)	\
    ( (This)->lpVtbl -> remove_FrameNavigationCompleted(This,token) ) 

#define ICoreWebView2_2_add_ScriptDialogOpening(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ScriptDialogOpening(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_ScriptDialogOpening(This,token)	\
    ( (This)->lpVtbl -> remove_ScriptDialogOpening(This,token) ) 

#define ICoreWebView2_2_add_PermissionRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_PermissionRequested(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_PermissionRequested(This,token)	\
    ( (This)->lpVtbl -> remove_PermissionRequested(This,token) ) 

#define ICoreWebView2_2_add_ProcessFailed(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ProcessFailed(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_ProcessFailed(This,token)	\
    ( (This)->lpVtbl -> remove_ProcessFailed(This,token) ) 

#define ICoreWebView2_2_AddScriptToExecuteOnDocumentCreated(This,javaScript,handler)	\
    ( (This)->lpVtbl -> AddScriptToExecuteOnDocumentCreated(This,javaScript,handler) ) 

#define ICoreWebView2_2_RemoveScriptToExecuteOnDocumentCreated(This,id)	\
    ( (This)->lpVtbl -> RemoveScriptToExecuteOnDocumentCreated(This,id) ) 

#define ICoreWebView2_2_ExecuteScript(This,javaScript,handler)	\
    ( (This)->lpVtbl -> ExecuteScript(This,javaScript,handler) ) 

#define ICoreWebView2_2_CapturePreview(This,imageFormat,imageStream,handler)	\
    ( (This)->lpVtbl -> CapturePreview(This,imageFormat,imageStream,handler) ) 

#define ICoreWebView2_2_Reload(This)	\
    ( (This)->lpVtbl -> Reload(This) ) 

#define ICoreWebView2_2_PostWebMessageAsJson(This,webMessageAsJson)	\
    ( (This)->lpVtbl -> PostWebMessageAsJson(This,webMessageAsJson) ) 

#define ICoreWebView2_2_PostWebMessageAsString(This,webMessageAsString)	\
    ( (This)->lpVtbl -> PostWebMessageAsString(This,webMessageAsString) ) 

#define ICoreWebView2_2_add_WebMessageReceived(This,handler,token)	\
    ( (This)->lpVtbl -> add_WebMessageReceived(This,handler,token) ) 

#define ICoreWebView2_2_remove_WebMessageReceived(This,token)	\
    ( (This)->lpVtbl -> remove_WebMessageReceived(This,token) ) 

#define ICoreWebView2_2_CallDevToolsProtocolMethod(This,methodName,parametersAsJson,handler)	\
    ( (This)->lpVtbl -> CallDevToolsProtocolMethod(This,methodName,parametersAsJson,handler) ) 

#define ICoreWebView2_2_get_BrowserProcessId(This,value)	\
    ( (This)->lpVtbl -> get_BrowserProcessId(This,value) ) 

#define ICoreWebView2_2_get_CanGoBack(This,canGoBack)	\
    ( (This)->lpVtbl -> get_CanGoBack(This,canGoBack) ) 

#define ICoreWebView2_2_get_CanGoForward(This,canGoForward)	\
    ( (This)->lpVtbl -> get_CanGoForward(This,canGoForward) ) 

#define ICoreWebView2_2_GoBack(This)	\
    ( (This)->lpVtbl -> GoBack(This) ) 

#define ICoreWebView2_2_GoForward(This)	\
    ( (This)->lpVtbl -> GoForward(This) ) 

#define ICoreWebView2_2_GetDevToolsProtocolEventReceiver(This,eventName,receiver)	\
    ( (This)->lpVtbl -> GetDevToolsProtocolEventReceiver(This,eventName,receiver) ) 

#define ICoreWebView2_2_Stop(This)	\
    ( (This)->lpVtbl -> Stop(This) ) 

#define ICoreWebView2_2_add_NewWindowRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NewWindowRequested(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_NewWindowRequested(This,token)	\
    ( (This)->lpVtbl -> remove_NewWindowRequested(This,token) ) 

#define ICoreWebView2_2_add_DocumentTitleChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_DocumentTitleChanged(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_DocumentTitleChanged(This,token)	\
    ( (This)->lpVtbl -> remove_DocumentTitleChanged(This,token) ) 

#define ICoreWebView2_2_get_DocumentTitle(This,title)	\
    ( (This)->lpVtbl -> get_DocumentTitle(This,title) ) 

#define ICoreWebView2_2_AddHostObjectToScript(This,name,object)	\
    ( (This)->lpVtbl -> AddHostObjectToScript(This,name,object) ) 

#define ICoreWebView2_2_RemoveHostObjectFromScript(This,name)	\
    ( (This)->lpVtbl -> RemoveHostObjectFromScript(This,name) ) 

#define ICoreWebView2_2_OpenDevToolsWindow(This)	\
    ( (This)->lpVtbl -> OpenDevToolsWindow(This) ) 

#define ICoreWebView2_2_add_ContainsFullScreenElementChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ContainsFullScreenElementChanged(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_ContainsFullScreenElementChanged(This,token)	\
    ( (This)->lpVtbl -> remove_ContainsFullScreenElementChanged(This,token) ) 

#define ICoreWebView2_2_get_ContainsFullScreenElement(This,containsFullScreenElement)	\
    ( (This)->lpVtbl -> get_ContainsFullScreenElement(This,containsFullScreenElement) ) 

#define ICoreWebView2_2_add_WebResourceRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WebResourceRequested(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_WebResourceRequested(This,token)	\
    ( (This)->lpVtbl -> remove_WebResourceRequested(This,token) ) 

#define ICoreWebView2_2_AddWebResourceRequestedFilter(This,uri,resourceContext)	\
    ( (This)->lpVtbl -> AddWebResourceRequestedFilter(This,uri,resourceContext) ) 

#define ICoreWebView2_2_RemoveWebResourceRequestedFilter(This,uri,resourceContext)	\
    ( (This)->lpVtbl -> RemoveWebResourceRequestedFilter(This,uri,resourceContext) ) 

#define ICoreWebView2_2_add_WindowCloseRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WindowCloseRequested(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_WindowCloseRequested(This,token)	\
    ( (This)->lpVtbl -> remove_WindowCloseRequested(This,token) ) 


#define ICoreWebView2_2_add_WebResourceResponseReceived(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WebResourceResponseReceived(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_WebResourceResponseReceived(This,token)	\
    ( (This)->lpVtbl -> remove_WebResourceResponseReceived(This,token) ) 

#define ICoreWebView2_2_NavigateWithWebResourceRequest(This,request)	\
    ( (This)->lpVtbl -> NavigateWithWebResourceRequest(This,request) ) 

#define ICoreWebView2_2_add_DOMContentLoaded(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_DOMContentLoaded(This,eventHandler,token) ) 

#define ICoreWebView2_2_remove_DOMContentLoaded(This,token)	\
    ( (This)->lpVtbl -> remove_DOMContentLoaded(This,token) ) 

#define ICoreWebView2_2_get_CookieManager(This,cookieManager)	\
    ( (This)->lpVtbl -> get_CookieManager(This,cookieManager) ) 

#define ICoreWebView2_2_get_Environment(This,environment)	\
    ( (This)->lpVtbl -> get_Environment(This,environment) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2_2_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2_3_INTERFACE_DEFINED__
#define __ICoreWebView2_3_INTERFACE_DEFINED__

/* interface ICoreWebView2_3 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2_3 = {0xA0D6DF20,0x3B92,0x416D,{0xAA,0x0C,0x43,0x7A,0x9C,0x72,0x78,0x57}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("A0D6DF20-3B92-416D-AA0C-437A9C727857")
    ICoreWebView2_3 : public ICoreWebView2_2
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE TrySuspend( 
            /* [in] */ ICoreWebView2TrySuspendCompletedHandler *handler) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Resume( void) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsSuspended( 
            /* [retval][out] */ BOOL *isSuspended) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE SetVirtualHostNameToFolderMapping( 
            /* [in] */ LPCWSTR hostName,
            /* [in] */ LPCWSTR folderPath,
            /* [in] */ COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND accessKind) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE ClearVirtualHostNameToFolderMapping( 
            /* [in] */ LPCWSTR hostName) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2_3Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2_3 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2_3 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2_3 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Settings )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ ICoreWebView2Settings **settings);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Source )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ LPWSTR *uri);
        
        HRESULT ( STDMETHODCALLTYPE *Navigate )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR uri);
        
        HRESULT ( STDMETHODCALLTYPE *NavigateToString )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR htmlContent);
        
        HRESULT ( STDMETHODCALLTYPE *add_NavigationStarting )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NavigationStarting )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ContentLoading )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2ContentLoadingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ContentLoading )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_SourceChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2SourceChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_SourceChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_HistoryChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2HistoryChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_HistoryChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_NavigationCompleted )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NavigationCompleted )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_FrameNavigationStarting )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2NavigationStartingEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_FrameNavigationStarting )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_FrameNavigationCompleted )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2NavigationCompletedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_FrameNavigationCompleted )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ScriptDialogOpening )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2ScriptDialogOpeningEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ScriptDialogOpening )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_PermissionRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2PermissionRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_PermissionRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_ProcessFailed )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2ProcessFailedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ProcessFailed )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *AddScriptToExecuteOnDocumentCreated )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveScriptToExecuteOnDocumentCreated )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR id);
        
        HRESULT ( STDMETHODCALLTYPE *ExecuteScript )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR javaScript,
            /* [in] */ ICoreWebView2ExecuteScriptCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CapturePreview )( 
            ICoreWebView2_3 * This,
            /* [in] */ COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT imageFormat,
            /* [in] */ IStream *imageStream,
            /* [in] */ ICoreWebView2CapturePreviewCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *Reload )( 
            ICoreWebView2_3 * This);
        
        HRESULT ( STDMETHODCALLTYPE *PostWebMessageAsJson )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR webMessageAsJson);
        
        HRESULT ( STDMETHODCALLTYPE *PostWebMessageAsString )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR webMessageAsString);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebMessageReceived )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2WebMessageReceivedEventHandler *handler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebMessageReceived )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *CallDevToolsProtocolMethod )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR methodName,
            /* [in] */ LPCWSTR parametersAsJson,
            /* [in] */ ICoreWebView2CallDevToolsProtocolMethodCompletedHandler *handler);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BrowserProcessId )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ UINT32 *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CanGoBack )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ BOOL *canGoBack);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CanGoForward )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ BOOL *canGoForward);
        
        HRESULT ( STDMETHODCALLTYPE *GoBack )( 
            ICoreWebView2_3 * This);
        
        HRESULT ( STDMETHODCALLTYPE *GoForward )( 
            ICoreWebView2_3 * This);
        
        HRESULT ( STDMETHODCALLTYPE *GetDevToolsProtocolEventReceiver )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR eventName,
            /* [retval][out] */ ICoreWebView2DevToolsProtocolEventReceiver **receiver);
        
        HRESULT ( STDMETHODCALLTYPE *Stop )( 
            ICoreWebView2_3 * This);
        
        HRESULT ( STDMETHODCALLTYPE *add_NewWindowRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2NewWindowRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NewWindowRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_DocumentTitleChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2DocumentTitleChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_DocumentTitleChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_DocumentTitle )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ LPWSTR *title);
        
        HRESULT ( STDMETHODCALLTYPE *AddHostObjectToScript )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ VARIANT *object);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveHostObjectFromScript )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR name);
        
        HRESULT ( STDMETHODCALLTYPE *OpenDevToolsWindow )( 
            ICoreWebView2_3 * This);
        
        HRESULT ( STDMETHODCALLTYPE *add_ContainsFullScreenElementChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2ContainsFullScreenElementChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ContainsFullScreenElementChanged )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ContainsFullScreenElement )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ BOOL *containsFullScreenElement);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebResourceRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2WebResourceRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebResourceRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *AddWebResourceRequestedFilter )( 
            ICoreWebView2_3 * This,
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveWebResourceRequestedFilter )( 
            ICoreWebView2_3 * This,
            /* [in] */ const LPCWSTR uri,
            /* [in] */ const COREWEBVIEW2_WEB_RESOURCE_CONTEXT resourceContext);
        
        HRESULT ( STDMETHODCALLTYPE *add_WindowCloseRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2WindowCloseRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WindowCloseRequested )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_WebResourceResponseReceived )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2WebResourceResponseReceivedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_WebResourceResponseReceived )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *NavigateWithWebResourceRequest )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2WebResourceRequest *request);
        
        HRESULT ( STDMETHODCALLTYPE *add_DOMContentLoaded )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2DOMContentLoadedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_DOMContentLoaded )( 
            ICoreWebView2_3 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CookieManager )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ ICoreWebView2CookieManager **cookieManager);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Environment )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ ICoreWebView2Environment **environment);
        
        HRESULT ( STDMETHODCALLTYPE *TrySuspend )( 
            ICoreWebView2_3 * This,
            /* [in] */ ICoreWebView2TrySuspendCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *Resume )( 
            ICoreWebView2_3 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsSuspended )( 
            ICoreWebView2_3 * This,
            /* [retval][out] */ BOOL *isSuspended);
        
        HRESULT ( STDMETHODCALLTYPE *SetVirtualHostNameToFolderMapping )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR hostName,
            /* [in] */ LPCWSTR folderPath,
            /* [in] */ COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND accessKind);
        
        HRESULT ( STDMETHODCALLTYPE *ClearVirtualHostNameToFolderMapping )( 
            ICoreWebView2_3 * This,
            /* [in] */ LPCWSTR hostName);
        
        END_INTERFACE
    } ICoreWebView2_3Vtbl;

    interface ICoreWebView2_3
    {
        CONST_VTBL struct ICoreWebView2_3Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2_3_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2_3_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2_3_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2_3_get_Settings(This,settings)	\
    ( (This)->lpVtbl -> get_Settings(This,settings) ) 

#define ICoreWebView2_3_get_Source(This,uri)	\
    ( (This)->lpVtbl -> get_Source(This,uri) ) 

#define ICoreWebView2_3_Navigate(This,uri)	\
    ( (This)->lpVtbl -> Navigate(This,uri) ) 

#define ICoreWebView2_3_NavigateToString(This,htmlContent)	\
    ( (This)->lpVtbl -> NavigateToString(This,htmlContent) ) 

#define ICoreWebView2_3_add_NavigationStarting(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NavigationStarting(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_NavigationStarting(This,token)	\
    ( (This)->lpVtbl -> remove_NavigationStarting(This,token) ) 

#define ICoreWebView2_3_add_ContentLoading(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ContentLoading(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_ContentLoading(This,token)	\
    ( (This)->lpVtbl -> remove_ContentLoading(This,token) ) 

#define ICoreWebView2_3_add_SourceChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_SourceChanged(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_SourceChanged(This,token)	\
    ( (This)->lpVtbl -> remove_SourceChanged(This,token) ) 

#define ICoreWebView2_3_add_HistoryChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_HistoryChanged(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_HistoryChanged(This,token)	\
    ( (This)->lpVtbl -> remove_HistoryChanged(This,token) ) 

#define ICoreWebView2_3_add_NavigationCompleted(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NavigationCompleted(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_NavigationCompleted(This,token)	\
    ( (This)->lpVtbl -> remove_NavigationCompleted(This,token) ) 

#define ICoreWebView2_3_add_FrameNavigationStarting(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_FrameNavigationStarting(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_FrameNavigationStarting(This,token)	\
    ( (This)->lpVtbl -> remove_FrameNavigationStarting(This,token) ) 

#define ICoreWebView2_3_add_FrameNavigationCompleted(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_FrameNavigationCompleted(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_FrameNavigationCompleted(This,token)	\
    ( (This)->lpVtbl -> remove_FrameNavigationCompleted(This,token) ) 

#define ICoreWebView2_3_add_ScriptDialogOpening(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ScriptDialogOpening(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_ScriptDialogOpening(This,token)	\
    ( (This)->lpVtbl -> remove_ScriptDialogOpening(This,token) ) 

#define ICoreWebView2_3_add_PermissionRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_PermissionRequested(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_PermissionRequested(This,token)	\
    ( (This)->lpVtbl -> remove_PermissionRequested(This,token) ) 

#define ICoreWebView2_3_add_ProcessFailed(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ProcessFailed(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_ProcessFailed(This,token)	\
    ( (This)->lpVtbl -> remove_ProcessFailed(This,token) ) 

#define ICoreWebView2_3_AddScriptToExecuteOnDocumentCreated(This,javaScript,handler)	\
    ( (This)->lpVtbl -> AddScriptToExecuteOnDocumentCreated(This,javaScript,handler) ) 

#define ICoreWebView2_3_RemoveScriptToExecuteOnDocumentCreated(This,id)	\
    ( (This)->lpVtbl -> RemoveScriptToExecuteOnDocumentCreated(This,id) ) 

#define ICoreWebView2_3_ExecuteScript(This,javaScript,handler)	\
    ( (This)->lpVtbl -> ExecuteScript(This,javaScript,handler) ) 

#define ICoreWebView2_3_CapturePreview(This,imageFormat,imageStream,handler)	\
    ( (This)->lpVtbl -> CapturePreview(This,imageFormat,imageStream,handler) ) 

#define ICoreWebView2_3_Reload(This)	\
    ( (This)->lpVtbl -> Reload(This) ) 

#define ICoreWebView2_3_PostWebMessageAsJson(This,webMessageAsJson)	\
    ( (This)->lpVtbl -> PostWebMessageAsJson(This,webMessageAsJson) ) 

#define ICoreWebView2_3_PostWebMessageAsString(This,webMessageAsString)	\
    ( (This)->lpVtbl -> PostWebMessageAsString(This,webMessageAsString) ) 

#define ICoreWebView2_3_add_WebMessageReceived(This,handler,token)	\
    ( (This)->lpVtbl -> add_WebMessageReceived(This,handler,token) ) 

#define ICoreWebView2_3_remove_WebMessageReceived(This,token)	\
    ( (This)->lpVtbl -> remove_WebMessageReceived(This,token) ) 

#define ICoreWebView2_3_CallDevToolsProtocolMethod(This,methodName,parametersAsJson,handler)	\
    ( (This)->lpVtbl -> CallDevToolsProtocolMethod(This,methodName,parametersAsJson,handler) ) 

#define ICoreWebView2_3_get_BrowserProcessId(This,value)	\
    ( (This)->lpVtbl -> get_BrowserProcessId(This,value) ) 

#define ICoreWebView2_3_get_CanGoBack(This,canGoBack)	\
    ( (This)->lpVtbl -> get_CanGoBack(This,canGoBack) ) 

#define ICoreWebView2_3_get_CanGoForward(This,canGoForward)	\
    ( (This)->lpVtbl -> get_CanGoForward(This,canGoForward) ) 

#define ICoreWebView2_3_GoBack(This)	\
    ( (This)->lpVtbl -> GoBack(This) ) 

#define ICoreWebView2_3_GoForward(This)	\
    ( (This)->lpVtbl -> GoForward(This) ) 

#define ICoreWebView2_3_GetDevToolsProtocolEventReceiver(This,eventName,receiver)	\
    ( (This)->lpVtbl -> GetDevToolsProtocolEventReceiver(This,eventName,receiver) ) 

#define ICoreWebView2_3_Stop(This)	\
    ( (This)->lpVtbl -> Stop(This) ) 

#define ICoreWebView2_3_add_NewWindowRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NewWindowRequested(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_NewWindowRequested(This,token)	\
    ( (This)->lpVtbl -> remove_NewWindowRequested(This,token) ) 

#define ICoreWebView2_3_add_DocumentTitleChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_DocumentTitleChanged(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_DocumentTitleChanged(This,token)	\
    ( (This)->lpVtbl -> remove_DocumentTitleChanged(This,token) ) 

#define ICoreWebView2_3_get_DocumentTitle(This,title)	\
    ( (This)->lpVtbl -> get_DocumentTitle(This,title) ) 

#define ICoreWebView2_3_AddHostObjectToScript(This,name,object)	\
    ( (This)->lpVtbl -> AddHostObjectToScript(This,name,object) ) 

#define ICoreWebView2_3_RemoveHostObjectFromScript(This,name)	\
    ( (This)->lpVtbl -> RemoveHostObjectFromScript(This,name) ) 

#define ICoreWebView2_3_OpenDevToolsWindow(This)	\
    ( (This)->lpVtbl -> OpenDevToolsWindow(This) ) 

#define ICoreWebView2_3_add_ContainsFullScreenElementChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ContainsFullScreenElementChanged(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_ContainsFullScreenElementChanged(This,token)	\
    ( (This)->lpVtbl -> remove_ContainsFullScreenElementChanged(This,token) ) 

#define ICoreWebView2_3_get_ContainsFullScreenElement(This,containsFullScreenElement)	\
    ( (This)->lpVtbl -> get_ContainsFullScreenElement(This,containsFullScreenElement) ) 

#define ICoreWebView2_3_add_WebResourceRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WebResourceRequested(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_WebResourceRequested(This,token)	\
    ( (This)->lpVtbl -> remove_WebResourceRequested(This,token) ) 

#define ICoreWebView2_3_AddWebResourceRequestedFilter(This,uri,resourceContext)	\
    ( (This)->lpVtbl -> AddWebResourceRequestedFilter(This,uri,resourceContext) ) 

#define ICoreWebView2_3_RemoveWebResourceRequestedFilter(This,uri,resourceContext)	\
    ( (This)->lpVtbl -> RemoveWebResourceRequestedFilter(This,uri,resourceContext) ) 

#define ICoreWebView2_3_add_WindowCloseRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WindowCloseRequested(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_WindowCloseRequested(This,token)	\
    ( (This)->lpVtbl -> remove_WindowCloseRequested(This,token) ) 


#define ICoreWebView2_3_add_WebResourceResponseReceived(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_WebResourceResponseReceived(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_WebResourceResponseReceived(This,token)	\
    ( (This)->lpVtbl -> remove_WebResourceResponseReceived(This,token) ) 

#define ICoreWebView2_3_NavigateWithWebResourceRequest(This,request)	\
    ( (This)->lpVtbl -> NavigateWithWebResourceRequest(This,request) ) 

#define ICoreWebView2_3_add_DOMContentLoaded(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_DOMContentLoaded(This,eventHandler,token) ) 

#define ICoreWebView2_3_remove_DOMContentLoaded(This,token)	\
    ( (This)->lpVtbl -> remove_DOMContentLoaded(This,token) ) 

#define ICoreWebView2_3_get_CookieManager(This,cookieManager)	\
    ( (This)->lpVtbl -> get_CookieManager(This,cookieManager) ) 

#define ICoreWebView2_3_get_Environment(This,environment)	\
    ( (This)->lpVtbl -> get_Environment(This,environment) ) 


#define ICoreWebView2_3_TrySuspend(This,handler)	\
    ( (This)->lpVtbl -> TrySuspend(This,handler) ) 

#define ICoreWebView2_3_Resume(This)	\
    ( (This)->lpVtbl -> Resume(This) ) 

#define ICoreWebView2_3_get_IsSuspended(This,isSuspended)	\
    ( (This)->lpVtbl -> get_IsSuspended(This,isSuspended) ) 

#define ICoreWebView2_3_SetVirtualHostNameToFolderMapping(This,hostName,folderPath,accessKind)	\
    ( (This)->lpVtbl -> SetVirtualHostNameToFolderMapping(This,hostName,folderPath,accessKind) ) 

#define ICoreWebView2_3_ClearVirtualHostNameToFolderMapping(This,hostName)	\
    ( (This)->lpVtbl -> ClearVirtualHostNameToFolderMapping(This,hostName) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2_3_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CompositionController_INTERFACE_DEFINED__
#define __ICoreWebView2CompositionController_INTERFACE_DEFINED__

/* interface ICoreWebView2CompositionController */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CompositionController = {0x3df9b733,0xb9ae,0x4a15,{0x86,0xb4,0xeb,0x9e,0xe9,0x82,0x64,0x69}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("3df9b733-b9ae-4a15-86b4-eb9ee9826469")
    ICoreWebView2CompositionController : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_RootVisualTarget( 
            /* [retval][out] */ IUnknown **target) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_RootVisualTarget( 
            /* [in] */ IUnknown *target) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE SendMouseInput( 
            /* [in] */ COREWEBVIEW2_MOUSE_EVENT_KIND eventKind,
            /* [in] */ COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS virtualKeys,
            /* [in] */ UINT32 mouseData,
            /* [in] */ POINT point) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE SendPointerInput( 
            /* [in] */ COREWEBVIEW2_POINTER_EVENT_KIND eventKind,
            /* [in] */ ICoreWebView2PointerInfo *pointerInfo) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Cursor( 
            /* [retval][out] */ HCURSOR *cursor) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_SystemCursorId( 
            /* [retval][out] */ UINT32 *systemCursorId) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_CursorChanged( 
            /* [in] */ ICoreWebView2CursorChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_CursorChanged( 
            /* [in] */ EventRegistrationToken token) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CompositionControllerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CompositionController * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CompositionController * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CompositionController * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_RootVisualTarget )( 
            ICoreWebView2CompositionController * This,
            /* [retval][out] */ IUnknown **target);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_RootVisualTarget )( 
            ICoreWebView2CompositionController * This,
            /* [in] */ IUnknown *target);
        
        HRESULT ( STDMETHODCALLTYPE *SendMouseInput )( 
            ICoreWebView2CompositionController * This,
            /* [in] */ COREWEBVIEW2_MOUSE_EVENT_KIND eventKind,
            /* [in] */ COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS virtualKeys,
            /* [in] */ UINT32 mouseData,
            /* [in] */ POINT point);
        
        HRESULT ( STDMETHODCALLTYPE *SendPointerInput )( 
            ICoreWebView2CompositionController * This,
            /* [in] */ COREWEBVIEW2_POINTER_EVENT_KIND eventKind,
            /* [in] */ ICoreWebView2PointerInfo *pointerInfo);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Cursor )( 
            ICoreWebView2CompositionController * This,
            /* [retval][out] */ HCURSOR *cursor);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_SystemCursorId )( 
            ICoreWebView2CompositionController * This,
            /* [retval][out] */ UINT32 *systemCursorId);
        
        HRESULT ( STDMETHODCALLTYPE *add_CursorChanged )( 
            ICoreWebView2CompositionController * This,
            /* [in] */ ICoreWebView2CursorChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_CursorChanged )( 
            ICoreWebView2CompositionController * This,
            /* [in] */ EventRegistrationToken token);
        
        END_INTERFACE
    } ICoreWebView2CompositionControllerVtbl;

    interface ICoreWebView2CompositionController
    {
        CONST_VTBL struct ICoreWebView2CompositionControllerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CompositionController_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CompositionController_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CompositionController_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CompositionController_get_RootVisualTarget(This,target)	\
    ( (This)->lpVtbl -> get_RootVisualTarget(This,target) ) 

#define ICoreWebView2CompositionController_put_RootVisualTarget(This,target)	\
    ( (This)->lpVtbl -> put_RootVisualTarget(This,target) ) 

#define ICoreWebView2CompositionController_SendMouseInput(This,eventKind,virtualKeys,mouseData,point)	\
    ( (This)->lpVtbl -> SendMouseInput(This,eventKind,virtualKeys,mouseData,point) ) 

#define ICoreWebView2CompositionController_SendPointerInput(This,eventKind,pointerInfo)	\
    ( (This)->lpVtbl -> SendPointerInput(This,eventKind,pointerInfo) ) 

#define ICoreWebView2CompositionController_get_Cursor(This,cursor)	\
    ( (This)->lpVtbl -> get_Cursor(This,cursor) ) 

#define ICoreWebView2CompositionController_get_SystemCursorId(This,systemCursorId)	\
    ( (This)->lpVtbl -> get_SystemCursorId(This,systemCursorId) ) 

#define ICoreWebView2CompositionController_add_CursorChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_CursorChanged(This,eventHandler,token) ) 

#define ICoreWebView2CompositionController_remove_CursorChanged(This,token)	\
    ( (This)->lpVtbl -> remove_CursorChanged(This,token) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CompositionController_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CompositionController2_INTERFACE_DEFINED__
#define __ICoreWebView2CompositionController2_INTERFACE_DEFINED__

/* interface ICoreWebView2CompositionController2 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CompositionController2 = {0x0b6a3d24,0x49cb,0x4806,{0xba,0x20,0xb5,0xe0,0x73,0x4a,0x7b,0x26}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("0b6a3d24-49cb-4806-ba20-b5e0734a7b26")
    ICoreWebView2CompositionController2 : public ICoreWebView2CompositionController
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_UIAProvider( 
            /* [retval][out] */ IUnknown **provider) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CompositionController2Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CompositionController2 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CompositionController2 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CompositionController2 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_RootVisualTarget )( 
            ICoreWebView2CompositionController2 * This,
            /* [retval][out] */ IUnknown **target);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_RootVisualTarget )( 
            ICoreWebView2CompositionController2 * This,
            /* [in] */ IUnknown *target);
        
        HRESULT ( STDMETHODCALLTYPE *SendMouseInput )( 
            ICoreWebView2CompositionController2 * This,
            /* [in] */ COREWEBVIEW2_MOUSE_EVENT_KIND eventKind,
            /* [in] */ COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS virtualKeys,
            /* [in] */ UINT32 mouseData,
            /* [in] */ POINT point);
        
        HRESULT ( STDMETHODCALLTYPE *SendPointerInput )( 
            ICoreWebView2CompositionController2 * This,
            /* [in] */ COREWEBVIEW2_POINTER_EVENT_KIND eventKind,
            /* [in] */ ICoreWebView2PointerInfo *pointerInfo);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Cursor )( 
            ICoreWebView2CompositionController2 * This,
            /* [retval][out] */ HCURSOR *cursor);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_SystemCursorId )( 
            ICoreWebView2CompositionController2 * This,
            /* [retval][out] */ UINT32 *systemCursorId);
        
        HRESULT ( STDMETHODCALLTYPE *add_CursorChanged )( 
            ICoreWebView2CompositionController2 * This,
            /* [in] */ ICoreWebView2CursorChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_CursorChanged )( 
            ICoreWebView2CompositionController2 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_UIAProvider )( 
            ICoreWebView2CompositionController2 * This,
            /* [retval][out] */ IUnknown **provider);
        
        END_INTERFACE
    } ICoreWebView2CompositionController2Vtbl;

    interface ICoreWebView2CompositionController2
    {
        CONST_VTBL struct ICoreWebView2CompositionController2Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CompositionController2_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CompositionController2_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CompositionController2_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CompositionController2_get_RootVisualTarget(This,target)	\
    ( (This)->lpVtbl -> get_RootVisualTarget(This,target) ) 

#define ICoreWebView2CompositionController2_put_RootVisualTarget(This,target)	\
    ( (This)->lpVtbl -> put_RootVisualTarget(This,target) ) 

#define ICoreWebView2CompositionController2_SendMouseInput(This,eventKind,virtualKeys,mouseData,point)	\
    ( (This)->lpVtbl -> SendMouseInput(This,eventKind,virtualKeys,mouseData,point) ) 

#define ICoreWebView2CompositionController2_SendPointerInput(This,eventKind,pointerInfo)	\
    ( (This)->lpVtbl -> SendPointerInput(This,eventKind,pointerInfo) ) 

#define ICoreWebView2CompositionController2_get_Cursor(This,cursor)	\
    ( (This)->lpVtbl -> get_Cursor(This,cursor) ) 

#define ICoreWebView2CompositionController2_get_SystemCursorId(This,systemCursorId)	\
    ( (This)->lpVtbl -> get_SystemCursorId(This,systemCursorId) ) 

#define ICoreWebView2CompositionController2_add_CursorChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_CursorChanged(This,eventHandler,token) ) 

#define ICoreWebView2CompositionController2_remove_CursorChanged(This,token)	\
    ( (This)->lpVtbl -> remove_CursorChanged(This,token) ) 


#define ICoreWebView2CompositionController2_get_UIAProvider(This,provider)	\
    ( (This)->lpVtbl -> get_UIAProvider(This,provider) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CompositionController2_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Controller_INTERFACE_DEFINED__
#define __ICoreWebView2Controller_INTERFACE_DEFINED__

/* interface ICoreWebView2Controller */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Controller = {0x4d00c0d1,0x9434,0x4eb6,{0x80,0x78,0x86,0x97,0xa5,0x60,0x33,0x4f}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("4d00c0d1-9434-4eb6-8078-8697a560334f")
    ICoreWebView2Controller : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsVisible( 
            /* [retval][out] */ BOOL *isVisible) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsVisible( 
            /* [in] */ BOOL isVisible) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Bounds( 
            /* [retval][out] */ RECT *bounds) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Bounds( 
            /* [in] */ RECT bounds) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ZoomFactor( 
            /* [retval][out] */ double *zoomFactor) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_ZoomFactor( 
            /* [in] */ double zoomFactor) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_ZoomFactorChanged( 
            /* [in] */ ICoreWebView2ZoomFactorChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_ZoomFactorChanged( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE SetBoundsAndZoomFactor( 
            /* [in] */ RECT bounds,
            /* [in] */ double zoomFactor) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE MoveFocus( 
            /* [in] */ COREWEBVIEW2_MOVE_FOCUS_REASON reason) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_MoveFocusRequested( 
            /* [in] */ ICoreWebView2MoveFocusRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_MoveFocusRequested( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_GotFocus( 
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_GotFocus( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_LostFocus( 
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_LostFocus( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_AcceleratorKeyPressed( 
            /* [in] */ ICoreWebView2AcceleratorKeyPressedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_AcceleratorKeyPressed( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ParentWindow( 
            /* [retval][out] */ HWND *parentWindow) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_ParentWindow( 
            /* [in] */ HWND parentWindow) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE NotifyParentWindowPositionChanged( void) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Close( void) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_CoreWebView2( 
            /* [retval][out] */ ICoreWebView2 **coreWebView2) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ControllerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Controller * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Controller * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Controller * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsVisible )( 
            ICoreWebView2Controller * This,
            /* [retval][out] */ BOOL *isVisible);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsVisible )( 
            ICoreWebView2Controller * This,
            /* [in] */ BOOL isVisible);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Bounds )( 
            ICoreWebView2Controller * This,
            /* [retval][out] */ RECT *bounds);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Bounds )( 
            ICoreWebView2Controller * This,
            /* [in] */ RECT bounds);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ZoomFactor )( 
            ICoreWebView2Controller * This,
            /* [retval][out] */ double *zoomFactor);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ZoomFactor )( 
            ICoreWebView2Controller * This,
            /* [in] */ double zoomFactor);
        
        HRESULT ( STDMETHODCALLTYPE *add_ZoomFactorChanged )( 
            ICoreWebView2Controller * This,
            /* [in] */ ICoreWebView2ZoomFactorChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ZoomFactorChanged )( 
            ICoreWebView2Controller * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *SetBoundsAndZoomFactor )( 
            ICoreWebView2Controller * This,
            /* [in] */ RECT bounds,
            /* [in] */ double zoomFactor);
        
        HRESULT ( STDMETHODCALLTYPE *MoveFocus )( 
            ICoreWebView2Controller * This,
            /* [in] */ COREWEBVIEW2_MOVE_FOCUS_REASON reason);
        
        HRESULT ( STDMETHODCALLTYPE *add_MoveFocusRequested )( 
            ICoreWebView2Controller * This,
            /* [in] */ ICoreWebView2MoveFocusRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_MoveFocusRequested )( 
            ICoreWebView2Controller * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_GotFocus )( 
            ICoreWebView2Controller * This,
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_GotFocus )( 
            ICoreWebView2Controller * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_LostFocus )( 
            ICoreWebView2Controller * This,
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_LostFocus )( 
            ICoreWebView2Controller * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_AcceleratorKeyPressed )( 
            ICoreWebView2Controller * This,
            /* [in] */ ICoreWebView2AcceleratorKeyPressedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_AcceleratorKeyPressed )( 
            ICoreWebView2Controller * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ParentWindow )( 
            ICoreWebView2Controller * This,
            /* [retval][out] */ HWND *parentWindow);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ParentWindow )( 
            ICoreWebView2Controller * This,
            /* [in] */ HWND parentWindow);
        
        HRESULT ( STDMETHODCALLTYPE *NotifyParentWindowPositionChanged )( 
            ICoreWebView2Controller * This);
        
        HRESULT ( STDMETHODCALLTYPE *Close )( 
            ICoreWebView2Controller * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CoreWebView2 )( 
            ICoreWebView2Controller * This,
            /* [retval][out] */ ICoreWebView2 **coreWebView2);
        
        END_INTERFACE
    } ICoreWebView2ControllerVtbl;

    interface ICoreWebView2Controller
    {
        CONST_VTBL struct ICoreWebView2ControllerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Controller_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Controller_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Controller_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Controller_get_IsVisible(This,isVisible)	\
    ( (This)->lpVtbl -> get_IsVisible(This,isVisible) ) 

#define ICoreWebView2Controller_put_IsVisible(This,isVisible)	\
    ( (This)->lpVtbl -> put_IsVisible(This,isVisible) ) 

#define ICoreWebView2Controller_get_Bounds(This,bounds)	\
    ( (This)->lpVtbl -> get_Bounds(This,bounds) ) 

#define ICoreWebView2Controller_put_Bounds(This,bounds)	\
    ( (This)->lpVtbl -> put_Bounds(This,bounds) ) 

#define ICoreWebView2Controller_get_ZoomFactor(This,zoomFactor)	\
    ( (This)->lpVtbl -> get_ZoomFactor(This,zoomFactor) ) 

#define ICoreWebView2Controller_put_ZoomFactor(This,zoomFactor)	\
    ( (This)->lpVtbl -> put_ZoomFactor(This,zoomFactor) ) 

#define ICoreWebView2Controller_add_ZoomFactorChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ZoomFactorChanged(This,eventHandler,token) ) 

#define ICoreWebView2Controller_remove_ZoomFactorChanged(This,token)	\
    ( (This)->lpVtbl -> remove_ZoomFactorChanged(This,token) ) 

#define ICoreWebView2Controller_SetBoundsAndZoomFactor(This,bounds,zoomFactor)	\
    ( (This)->lpVtbl -> SetBoundsAndZoomFactor(This,bounds,zoomFactor) ) 

#define ICoreWebView2Controller_MoveFocus(This,reason)	\
    ( (This)->lpVtbl -> MoveFocus(This,reason) ) 

#define ICoreWebView2Controller_add_MoveFocusRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_MoveFocusRequested(This,eventHandler,token) ) 

#define ICoreWebView2Controller_remove_MoveFocusRequested(This,token)	\
    ( (This)->lpVtbl -> remove_MoveFocusRequested(This,token) ) 

#define ICoreWebView2Controller_add_GotFocus(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_GotFocus(This,eventHandler,token) ) 

#define ICoreWebView2Controller_remove_GotFocus(This,token)	\
    ( (This)->lpVtbl -> remove_GotFocus(This,token) ) 

#define ICoreWebView2Controller_add_LostFocus(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_LostFocus(This,eventHandler,token) ) 

#define ICoreWebView2Controller_remove_LostFocus(This,token)	\
    ( (This)->lpVtbl -> remove_LostFocus(This,token) ) 

#define ICoreWebView2Controller_add_AcceleratorKeyPressed(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_AcceleratorKeyPressed(This,eventHandler,token) ) 

#define ICoreWebView2Controller_remove_AcceleratorKeyPressed(This,token)	\
    ( (This)->lpVtbl -> remove_AcceleratorKeyPressed(This,token) ) 

#define ICoreWebView2Controller_get_ParentWindow(This,parentWindow)	\
    ( (This)->lpVtbl -> get_ParentWindow(This,parentWindow) ) 

#define ICoreWebView2Controller_put_ParentWindow(This,parentWindow)	\
    ( (This)->lpVtbl -> put_ParentWindow(This,parentWindow) ) 

#define ICoreWebView2Controller_NotifyParentWindowPositionChanged(This)	\
    ( (This)->lpVtbl -> NotifyParentWindowPositionChanged(This) ) 

#define ICoreWebView2Controller_Close(This)	\
    ( (This)->lpVtbl -> Close(This) ) 

#define ICoreWebView2Controller_get_CoreWebView2(This,coreWebView2)	\
    ( (This)->lpVtbl -> get_CoreWebView2(This,coreWebView2) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Controller_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Controller2_INTERFACE_DEFINED__
#define __ICoreWebView2Controller2_INTERFACE_DEFINED__

/* interface ICoreWebView2Controller2 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Controller2 = {0xc979903e,0xd4ca,0x4228,{0x92,0xeb,0x47,0xee,0x3f,0xa9,0x6e,0xab}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("c979903e-d4ca-4228-92eb-47ee3fa96eab")
    ICoreWebView2Controller2 : public ICoreWebView2Controller
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_DefaultBackgroundColor( 
            /* [retval][out] */ COREWEBVIEW2_COLOR *backgroundColor) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_DefaultBackgroundColor( 
            /* [in] */ COREWEBVIEW2_COLOR backgroundColor) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Controller2Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Controller2 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Controller2 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsVisible )( 
            ICoreWebView2Controller2 * This,
            /* [retval][out] */ BOOL *isVisible);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsVisible )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ BOOL isVisible);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Bounds )( 
            ICoreWebView2Controller2 * This,
            /* [retval][out] */ RECT *bounds);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Bounds )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ RECT bounds);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ZoomFactor )( 
            ICoreWebView2Controller2 * This,
            /* [retval][out] */ double *zoomFactor);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ZoomFactor )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ double zoomFactor);
        
        HRESULT ( STDMETHODCALLTYPE *add_ZoomFactorChanged )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ ICoreWebView2ZoomFactorChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ZoomFactorChanged )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *SetBoundsAndZoomFactor )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ RECT bounds,
            /* [in] */ double zoomFactor);
        
        HRESULT ( STDMETHODCALLTYPE *MoveFocus )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ COREWEBVIEW2_MOVE_FOCUS_REASON reason);
        
        HRESULT ( STDMETHODCALLTYPE *add_MoveFocusRequested )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ ICoreWebView2MoveFocusRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_MoveFocusRequested )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_GotFocus )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_GotFocus )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_LostFocus )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_LostFocus )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_AcceleratorKeyPressed )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ ICoreWebView2AcceleratorKeyPressedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_AcceleratorKeyPressed )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ParentWindow )( 
            ICoreWebView2Controller2 * This,
            /* [retval][out] */ HWND *parentWindow);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ParentWindow )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ HWND parentWindow);
        
        HRESULT ( STDMETHODCALLTYPE *NotifyParentWindowPositionChanged )( 
            ICoreWebView2Controller2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *Close )( 
            ICoreWebView2Controller2 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CoreWebView2 )( 
            ICoreWebView2Controller2 * This,
            /* [retval][out] */ ICoreWebView2 **coreWebView2);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_DefaultBackgroundColor )( 
            ICoreWebView2Controller2 * This,
            /* [retval][out] */ COREWEBVIEW2_COLOR *backgroundColor);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_DefaultBackgroundColor )( 
            ICoreWebView2Controller2 * This,
            /* [in] */ COREWEBVIEW2_COLOR backgroundColor);
        
        END_INTERFACE
    } ICoreWebView2Controller2Vtbl;

    interface ICoreWebView2Controller2
    {
        CONST_VTBL struct ICoreWebView2Controller2Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Controller2_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Controller2_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Controller2_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Controller2_get_IsVisible(This,isVisible)	\
    ( (This)->lpVtbl -> get_IsVisible(This,isVisible) ) 

#define ICoreWebView2Controller2_put_IsVisible(This,isVisible)	\
    ( (This)->lpVtbl -> put_IsVisible(This,isVisible) ) 

#define ICoreWebView2Controller2_get_Bounds(This,bounds)	\
    ( (This)->lpVtbl -> get_Bounds(This,bounds) ) 

#define ICoreWebView2Controller2_put_Bounds(This,bounds)	\
    ( (This)->lpVtbl -> put_Bounds(This,bounds) ) 

#define ICoreWebView2Controller2_get_ZoomFactor(This,zoomFactor)	\
    ( (This)->lpVtbl -> get_ZoomFactor(This,zoomFactor) ) 

#define ICoreWebView2Controller2_put_ZoomFactor(This,zoomFactor)	\
    ( (This)->lpVtbl -> put_ZoomFactor(This,zoomFactor) ) 

#define ICoreWebView2Controller2_add_ZoomFactorChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ZoomFactorChanged(This,eventHandler,token) ) 

#define ICoreWebView2Controller2_remove_ZoomFactorChanged(This,token)	\
    ( (This)->lpVtbl -> remove_ZoomFactorChanged(This,token) ) 

#define ICoreWebView2Controller2_SetBoundsAndZoomFactor(This,bounds,zoomFactor)	\
    ( (This)->lpVtbl -> SetBoundsAndZoomFactor(This,bounds,zoomFactor) ) 

#define ICoreWebView2Controller2_MoveFocus(This,reason)	\
    ( (This)->lpVtbl -> MoveFocus(This,reason) ) 

#define ICoreWebView2Controller2_add_MoveFocusRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_MoveFocusRequested(This,eventHandler,token) ) 

#define ICoreWebView2Controller2_remove_MoveFocusRequested(This,token)	\
    ( (This)->lpVtbl -> remove_MoveFocusRequested(This,token) ) 

#define ICoreWebView2Controller2_add_GotFocus(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_GotFocus(This,eventHandler,token) ) 

#define ICoreWebView2Controller2_remove_GotFocus(This,token)	\
    ( (This)->lpVtbl -> remove_GotFocus(This,token) ) 

#define ICoreWebView2Controller2_add_LostFocus(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_LostFocus(This,eventHandler,token) ) 

#define ICoreWebView2Controller2_remove_LostFocus(This,token)	\
    ( (This)->lpVtbl -> remove_LostFocus(This,token) ) 

#define ICoreWebView2Controller2_add_AcceleratorKeyPressed(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_AcceleratorKeyPressed(This,eventHandler,token) ) 

#define ICoreWebView2Controller2_remove_AcceleratorKeyPressed(This,token)	\
    ( (This)->lpVtbl -> remove_AcceleratorKeyPressed(This,token) ) 

#define ICoreWebView2Controller2_get_ParentWindow(This,parentWindow)	\
    ( (This)->lpVtbl -> get_ParentWindow(This,parentWindow) ) 

#define ICoreWebView2Controller2_put_ParentWindow(This,parentWindow)	\
    ( (This)->lpVtbl -> put_ParentWindow(This,parentWindow) ) 

#define ICoreWebView2Controller2_NotifyParentWindowPositionChanged(This)	\
    ( (This)->lpVtbl -> NotifyParentWindowPositionChanged(This) ) 

#define ICoreWebView2Controller2_Close(This)	\
    ( (This)->lpVtbl -> Close(This) ) 

#define ICoreWebView2Controller2_get_CoreWebView2(This,coreWebView2)	\
    ( (This)->lpVtbl -> get_CoreWebView2(This,coreWebView2) ) 


#define ICoreWebView2Controller2_get_DefaultBackgroundColor(This,backgroundColor)	\
    ( (This)->lpVtbl -> get_DefaultBackgroundColor(This,backgroundColor) ) 

#define ICoreWebView2Controller2_put_DefaultBackgroundColor(This,backgroundColor)	\
    ( (This)->lpVtbl -> put_DefaultBackgroundColor(This,backgroundColor) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Controller2_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Controller3_INTERFACE_DEFINED__
#define __ICoreWebView2Controller3_INTERFACE_DEFINED__

/* interface ICoreWebView2Controller3 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Controller3 = {0xf9614724,0x5d2b,0x41dc,{0xae,0xf7,0x73,0xd6,0x2b,0x51,0x54,0x3b}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("f9614724-5d2b-41dc-aef7-73d62b51543b")
    ICoreWebView2Controller3 : public ICoreWebView2Controller2
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_RasterizationScale( 
            /* [retval][out] */ double *scale) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_RasterizationScale( 
            /* [in] */ double scale) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ShouldDetectMonitorScaleChanges( 
            /* [retval][out] */ BOOL *value) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_ShouldDetectMonitorScaleChanges( 
            /* [in] */ BOOL value) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_RasterizationScaleChanged( 
            /* [in] */ ICoreWebView2RasterizationScaleChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_RasterizationScaleChanged( 
            /* [in] */ EventRegistrationToken token) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_BoundsMode( 
            /* [retval][out] */ COREWEBVIEW2_BOUNDS_MODE *boundsMode) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_BoundsMode( 
            /* [in] */ COREWEBVIEW2_BOUNDS_MODE boundsMode) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Controller3Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Controller3 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Controller3 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsVisible )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ BOOL *isVisible);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsVisible )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ BOOL isVisible);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Bounds )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ RECT *bounds);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Bounds )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ RECT bounds);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ZoomFactor )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ double *zoomFactor);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ZoomFactor )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ double zoomFactor);
        
        HRESULT ( STDMETHODCALLTYPE *add_ZoomFactorChanged )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ ICoreWebView2ZoomFactorChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_ZoomFactorChanged )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *SetBoundsAndZoomFactor )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ RECT bounds,
            /* [in] */ double zoomFactor);
        
        HRESULT ( STDMETHODCALLTYPE *MoveFocus )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ COREWEBVIEW2_MOVE_FOCUS_REASON reason);
        
        HRESULT ( STDMETHODCALLTYPE *add_MoveFocusRequested )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ ICoreWebView2MoveFocusRequestedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_MoveFocusRequested )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_GotFocus )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_GotFocus )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_LostFocus )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ ICoreWebView2FocusChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_LostFocus )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *add_AcceleratorKeyPressed )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ ICoreWebView2AcceleratorKeyPressedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_AcceleratorKeyPressed )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ParentWindow )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ HWND *parentWindow);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ParentWindow )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ HWND parentWindow);
        
        HRESULT ( STDMETHODCALLTYPE *NotifyParentWindowPositionChanged )( 
            ICoreWebView2Controller3 * This);
        
        HRESULT ( STDMETHODCALLTYPE *Close )( 
            ICoreWebView2Controller3 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_CoreWebView2 )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ ICoreWebView2 **coreWebView2);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_DefaultBackgroundColor )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ COREWEBVIEW2_COLOR *backgroundColor);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_DefaultBackgroundColor )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ COREWEBVIEW2_COLOR backgroundColor);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_RasterizationScale )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ double *scale);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_RasterizationScale )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ double scale);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ShouldDetectMonitorScaleChanges )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ BOOL *value);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ShouldDetectMonitorScaleChanges )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ BOOL value);
        
        HRESULT ( STDMETHODCALLTYPE *add_RasterizationScaleChanged )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ ICoreWebView2RasterizationScaleChangedEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_RasterizationScaleChanged )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ EventRegistrationToken token);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BoundsMode )( 
            ICoreWebView2Controller3 * This,
            /* [retval][out] */ COREWEBVIEW2_BOUNDS_MODE *boundsMode);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_BoundsMode )( 
            ICoreWebView2Controller3 * This,
            /* [in] */ COREWEBVIEW2_BOUNDS_MODE boundsMode);
        
        END_INTERFACE
    } ICoreWebView2Controller3Vtbl;

    interface ICoreWebView2Controller3
    {
        CONST_VTBL struct ICoreWebView2Controller3Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Controller3_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Controller3_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Controller3_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Controller3_get_IsVisible(This,isVisible)	\
    ( (This)->lpVtbl -> get_IsVisible(This,isVisible) ) 

#define ICoreWebView2Controller3_put_IsVisible(This,isVisible)	\
    ( (This)->lpVtbl -> put_IsVisible(This,isVisible) ) 

#define ICoreWebView2Controller3_get_Bounds(This,bounds)	\
    ( (This)->lpVtbl -> get_Bounds(This,bounds) ) 

#define ICoreWebView2Controller3_put_Bounds(This,bounds)	\
    ( (This)->lpVtbl -> put_Bounds(This,bounds) ) 

#define ICoreWebView2Controller3_get_ZoomFactor(This,zoomFactor)	\
    ( (This)->lpVtbl -> get_ZoomFactor(This,zoomFactor) ) 

#define ICoreWebView2Controller3_put_ZoomFactor(This,zoomFactor)	\
    ( (This)->lpVtbl -> put_ZoomFactor(This,zoomFactor) ) 

#define ICoreWebView2Controller3_add_ZoomFactorChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_ZoomFactorChanged(This,eventHandler,token) ) 

#define ICoreWebView2Controller3_remove_ZoomFactorChanged(This,token)	\
    ( (This)->lpVtbl -> remove_ZoomFactorChanged(This,token) ) 

#define ICoreWebView2Controller3_SetBoundsAndZoomFactor(This,bounds,zoomFactor)	\
    ( (This)->lpVtbl -> SetBoundsAndZoomFactor(This,bounds,zoomFactor) ) 

#define ICoreWebView2Controller3_MoveFocus(This,reason)	\
    ( (This)->lpVtbl -> MoveFocus(This,reason) ) 

#define ICoreWebView2Controller3_add_MoveFocusRequested(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_MoveFocusRequested(This,eventHandler,token) ) 

#define ICoreWebView2Controller3_remove_MoveFocusRequested(This,token)	\
    ( (This)->lpVtbl -> remove_MoveFocusRequested(This,token) ) 

#define ICoreWebView2Controller3_add_GotFocus(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_GotFocus(This,eventHandler,token) ) 

#define ICoreWebView2Controller3_remove_GotFocus(This,token)	\
    ( (This)->lpVtbl -> remove_GotFocus(This,token) ) 

#define ICoreWebView2Controller3_add_LostFocus(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_LostFocus(This,eventHandler,token) ) 

#define ICoreWebView2Controller3_remove_LostFocus(This,token)	\
    ( (This)->lpVtbl -> remove_LostFocus(This,token) ) 

#define ICoreWebView2Controller3_add_AcceleratorKeyPressed(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_AcceleratorKeyPressed(This,eventHandler,token) ) 

#define ICoreWebView2Controller3_remove_AcceleratorKeyPressed(This,token)	\
    ( (This)->lpVtbl -> remove_AcceleratorKeyPressed(This,token) ) 

#define ICoreWebView2Controller3_get_ParentWindow(This,parentWindow)	\
    ( (This)->lpVtbl -> get_ParentWindow(This,parentWindow) ) 

#define ICoreWebView2Controller3_put_ParentWindow(This,parentWindow)	\
    ( (This)->lpVtbl -> put_ParentWindow(This,parentWindow) ) 

#define ICoreWebView2Controller3_NotifyParentWindowPositionChanged(This)	\
    ( (This)->lpVtbl -> NotifyParentWindowPositionChanged(This) ) 

#define ICoreWebView2Controller3_Close(This)	\
    ( (This)->lpVtbl -> Close(This) ) 

#define ICoreWebView2Controller3_get_CoreWebView2(This,coreWebView2)	\
    ( (This)->lpVtbl -> get_CoreWebView2(This,coreWebView2) ) 


#define ICoreWebView2Controller3_get_DefaultBackgroundColor(This,backgroundColor)	\
    ( (This)->lpVtbl -> get_DefaultBackgroundColor(This,backgroundColor) ) 

#define ICoreWebView2Controller3_put_DefaultBackgroundColor(This,backgroundColor)	\
    ( (This)->lpVtbl -> put_DefaultBackgroundColor(This,backgroundColor) ) 


#define ICoreWebView2Controller3_get_RasterizationScale(This,scale)	\
    ( (This)->lpVtbl -> get_RasterizationScale(This,scale) ) 

#define ICoreWebView2Controller3_put_RasterizationScale(This,scale)	\
    ( (This)->lpVtbl -> put_RasterizationScale(This,scale) ) 

#define ICoreWebView2Controller3_get_ShouldDetectMonitorScaleChanges(This,value)	\
    ( (This)->lpVtbl -> get_ShouldDetectMonitorScaleChanges(This,value) ) 

#define ICoreWebView2Controller3_put_ShouldDetectMonitorScaleChanges(This,value)	\
    ( (This)->lpVtbl -> put_ShouldDetectMonitorScaleChanges(This,value) ) 

#define ICoreWebView2Controller3_add_RasterizationScaleChanged(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_RasterizationScaleChanged(This,eventHandler,token) ) 

#define ICoreWebView2Controller3_remove_RasterizationScaleChanged(This,token)	\
    ( (This)->lpVtbl -> remove_RasterizationScaleChanged(This,token) ) 

#define ICoreWebView2Controller3_get_BoundsMode(This,boundsMode)	\
    ( (This)->lpVtbl -> get_BoundsMode(This,boundsMode) ) 

#define ICoreWebView2Controller3_put_BoundsMode(This,boundsMode)	\
    ( (This)->lpVtbl -> put_BoundsMode(This,boundsMode) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Controller3_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2ContentLoadingEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2ContentLoadingEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ContentLoadingEventArgs = {0x0c8a1275,0x9b6b,0x4901,{0x87,0xad,0x70,0xdf,0x25,0xba,0xfa,0x6e}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("0c8a1275-9b6b-4901-87ad-70df25bafa6e")
    ICoreWebView2ContentLoadingEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsErrorPage( 
            /* [retval][out] */ BOOL *isErrorPage) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_NavigationId( 
            /* [retval][out] */ UINT64 *navigationId) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ContentLoadingEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ContentLoadingEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ContentLoadingEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ContentLoadingEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsErrorPage )( 
            ICoreWebView2ContentLoadingEventArgs * This,
            /* [retval][out] */ BOOL *isErrorPage);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_NavigationId )( 
            ICoreWebView2ContentLoadingEventArgs * This,
            /* [retval][out] */ UINT64 *navigationId);
        
        END_INTERFACE
    } ICoreWebView2ContentLoadingEventArgsVtbl;

    interface ICoreWebView2ContentLoadingEventArgs
    {
        CONST_VTBL struct ICoreWebView2ContentLoadingEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ContentLoadingEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ContentLoadingEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ContentLoadingEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ContentLoadingEventArgs_get_IsErrorPage(This,isErrorPage)	\
    ( (This)->lpVtbl -> get_IsErrorPage(This,isErrorPage) ) 

#define ICoreWebView2ContentLoadingEventArgs_get_NavigationId(This,navigationId)	\
    ( (This)->lpVtbl -> get_NavigationId(This,navigationId) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ContentLoadingEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ContentLoadingEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ContentLoadingEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ContentLoadingEventHandler = {0x364471e7,0xf2be,0x4910,{0xbd,0xba,0xd7,0x20,0x77,0xd5,0x1c,0x4b}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("364471e7-f2be-4910-bdba-d72077d51c4b")
    ICoreWebView2ContentLoadingEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2ContentLoadingEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ContentLoadingEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ContentLoadingEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ContentLoadingEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ContentLoadingEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2ContentLoadingEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2ContentLoadingEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2ContentLoadingEventHandlerVtbl;

    interface ICoreWebView2ContentLoadingEventHandler
    {
        CONST_VTBL struct ICoreWebView2ContentLoadingEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ContentLoadingEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ContentLoadingEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ContentLoadingEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ContentLoadingEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ContentLoadingEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Cookie_INTERFACE_DEFINED__
#define __ICoreWebView2Cookie_INTERFACE_DEFINED__

/* interface ICoreWebView2Cookie */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Cookie = {0xAD26D6BE,0x1486,0x43E6,{0xBF,0x87,0xA2,0x03,0x40,0x06,0xCA,0x21}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("AD26D6BE-1486-43E6-BF87-A2034006CA21")
    ICoreWebView2Cookie : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Name( 
            /* [retval][out] */ LPWSTR *name) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Value( 
            /* [retval][out] */ LPWSTR *value) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Value( 
            /* [in] */ LPCWSTR value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Domain( 
            /* [retval][out] */ LPWSTR *domain) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Path( 
            /* [retval][out] */ LPWSTR *path) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Expires( 
            /* [retval][out] */ double *expires) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Expires( 
            /* [in] */ double expires) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsHttpOnly( 
            /* [retval][out] */ BOOL *isHttpOnly) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsHttpOnly( 
            /* [in] */ BOOL isHttpOnly) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_SameSite( 
            /* [retval][out] */ COREWEBVIEW2_COOKIE_SAME_SITE_KIND *sameSite) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_SameSite( 
            /* [in] */ COREWEBVIEW2_COOKIE_SAME_SITE_KIND sameSite) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsSecure( 
            /* [retval][out] */ BOOL *isSecure) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsSecure( 
            /* [in] */ BOOL isSecure) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsSession( 
            /* [retval][out] */ BOOL *isSession) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CookieVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Cookie * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Cookie * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Cookie * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Name )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ LPWSTR *name);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Value )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ LPWSTR *value);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Value )( 
            ICoreWebView2Cookie * This,
            /* [in] */ LPCWSTR value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Domain )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ LPWSTR *domain);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Path )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ LPWSTR *path);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Expires )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ double *expires);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Expires )( 
            ICoreWebView2Cookie * This,
            /* [in] */ double expires);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsHttpOnly )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ BOOL *isHttpOnly);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsHttpOnly )( 
            ICoreWebView2Cookie * This,
            /* [in] */ BOOL isHttpOnly);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_SameSite )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ COREWEBVIEW2_COOKIE_SAME_SITE_KIND *sameSite);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_SameSite )( 
            ICoreWebView2Cookie * This,
            /* [in] */ COREWEBVIEW2_COOKIE_SAME_SITE_KIND sameSite);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsSecure )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ BOOL *isSecure);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsSecure )( 
            ICoreWebView2Cookie * This,
            /* [in] */ BOOL isSecure);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsSession )( 
            ICoreWebView2Cookie * This,
            /* [retval][out] */ BOOL *isSession);
        
        END_INTERFACE
    } ICoreWebView2CookieVtbl;

    interface ICoreWebView2Cookie
    {
        CONST_VTBL struct ICoreWebView2CookieVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Cookie_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Cookie_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Cookie_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Cookie_get_Name(This,name)	\
    ( (This)->lpVtbl -> get_Name(This,name) ) 

#define ICoreWebView2Cookie_get_Value(This,value)	\
    ( (This)->lpVtbl -> get_Value(This,value) ) 

#define ICoreWebView2Cookie_put_Value(This,value)	\
    ( (This)->lpVtbl -> put_Value(This,value) ) 

#define ICoreWebView2Cookie_get_Domain(This,domain)	\
    ( (This)->lpVtbl -> get_Domain(This,domain) ) 

#define ICoreWebView2Cookie_get_Path(This,path)	\
    ( (This)->lpVtbl -> get_Path(This,path) ) 

#define ICoreWebView2Cookie_get_Expires(This,expires)	\
    ( (This)->lpVtbl -> get_Expires(This,expires) ) 

#define ICoreWebView2Cookie_put_Expires(This,expires)	\
    ( (This)->lpVtbl -> put_Expires(This,expires) ) 

#define ICoreWebView2Cookie_get_IsHttpOnly(This,isHttpOnly)	\
    ( (This)->lpVtbl -> get_IsHttpOnly(This,isHttpOnly) ) 

#define ICoreWebView2Cookie_put_IsHttpOnly(This,isHttpOnly)	\
    ( (This)->lpVtbl -> put_IsHttpOnly(This,isHttpOnly) ) 

#define ICoreWebView2Cookie_get_SameSite(This,sameSite)	\
    ( (This)->lpVtbl -> get_SameSite(This,sameSite) ) 

#define ICoreWebView2Cookie_put_SameSite(This,sameSite)	\
    ( (This)->lpVtbl -> put_SameSite(This,sameSite) ) 

#define ICoreWebView2Cookie_get_IsSecure(This,isSecure)	\
    ( (This)->lpVtbl -> get_IsSecure(This,isSecure) ) 

#define ICoreWebView2Cookie_put_IsSecure(This,isSecure)	\
    ( (This)->lpVtbl -> put_IsSecure(This,isSecure) ) 

#define ICoreWebView2Cookie_get_IsSession(This,isSession)	\
    ( (This)->lpVtbl -> get_IsSession(This,isSession) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Cookie_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CookieList_INTERFACE_DEFINED__
#define __ICoreWebView2CookieList_INTERFACE_DEFINED__

/* interface ICoreWebView2CookieList */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CookieList = {0xF7F6F714,0x5D2A,0x43C6,{0x95,0x03,0x34,0x6E,0xCE,0x02,0xD1,0x86}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("F7F6F714-5D2A-43C6-9503-346ECE02D186")
    ICoreWebView2CookieList : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Count( 
            /* [retval][out] */ UINT *count) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetValueAtIndex( 
            /* [in] */ UINT index,
            /* [retval][out] */ ICoreWebView2Cookie **cookie) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CookieListVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CookieList * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CookieList * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CookieList * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Count )( 
            ICoreWebView2CookieList * This,
            /* [retval][out] */ UINT *count);
        
        HRESULT ( STDMETHODCALLTYPE *GetValueAtIndex )( 
            ICoreWebView2CookieList * This,
            /* [in] */ UINT index,
            /* [retval][out] */ ICoreWebView2Cookie **cookie);
        
        END_INTERFACE
    } ICoreWebView2CookieListVtbl;

    interface ICoreWebView2CookieList
    {
        CONST_VTBL struct ICoreWebView2CookieListVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CookieList_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CookieList_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CookieList_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CookieList_get_Count(This,count)	\
    ( (This)->lpVtbl -> get_Count(This,count) ) 

#define ICoreWebView2CookieList_GetValueAtIndex(This,index,cookie)	\
    ( (This)->lpVtbl -> GetValueAtIndex(This,index,cookie) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CookieList_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CookieManager_INTERFACE_DEFINED__
#define __ICoreWebView2CookieManager_INTERFACE_DEFINED__

/* interface ICoreWebView2CookieManager */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CookieManager = {0x177CD9E7,0xB6F5,0x451A,{0x94,0xA0,0x5D,0x7A,0x3A,0x4C,0x41,0x41}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("177CD9E7-B6F5-451A-94A0-5D7A3A4C4141")
    ICoreWebView2CookieManager : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE CreateCookie( 
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR value,
            /* [in] */ LPCWSTR domain,
            /* [in] */ LPCWSTR path,
            /* [retval][out] */ ICoreWebView2Cookie **cookie) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE CopyCookie( 
            /* [in] */ ICoreWebView2Cookie *cookieParam,
            /* [retval][out] */ ICoreWebView2Cookie **cookie) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetCookies( 
            /* [in] */ LPCWSTR uri,
            /* [in] */ ICoreWebView2GetCookiesCompletedHandler *handler) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE AddOrUpdateCookie( 
            /* [in] */ ICoreWebView2Cookie *cookie) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE DeleteCookie( 
            /* [in] */ ICoreWebView2Cookie *cookie) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE DeleteCookies( 
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR uri) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE DeleteCookiesWithDomainAndPath( 
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR domain,
            /* [in] */ LPCWSTR path) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE DeleteAllCookies( void) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CookieManagerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CookieManager * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CookieManager * This);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCookie )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR value,
            /* [in] */ LPCWSTR domain,
            /* [in] */ LPCWSTR path,
            /* [retval][out] */ ICoreWebView2Cookie **cookie);
        
        HRESULT ( STDMETHODCALLTYPE *CopyCookie )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ ICoreWebView2Cookie *cookieParam,
            /* [retval][out] */ ICoreWebView2Cookie **cookie);
        
        HRESULT ( STDMETHODCALLTYPE *GetCookies )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ LPCWSTR uri,
            /* [in] */ ICoreWebView2GetCookiesCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *AddOrUpdateCookie )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ ICoreWebView2Cookie *cookie);
        
        HRESULT ( STDMETHODCALLTYPE *DeleteCookie )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ ICoreWebView2Cookie *cookie);
        
        HRESULT ( STDMETHODCALLTYPE *DeleteCookies )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR uri);
        
        HRESULT ( STDMETHODCALLTYPE *DeleteCookiesWithDomainAndPath )( 
            ICoreWebView2CookieManager * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR domain,
            /* [in] */ LPCWSTR path);
        
        HRESULT ( STDMETHODCALLTYPE *DeleteAllCookies )( 
            ICoreWebView2CookieManager * This);
        
        END_INTERFACE
    } ICoreWebView2CookieManagerVtbl;

    interface ICoreWebView2CookieManager
    {
        CONST_VTBL struct ICoreWebView2CookieManagerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CookieManager_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CookieManager_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CookieManager_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CookieManager_CreateCookie(This,name,value,domain,path,cookie)	\
    ( (This)->lpVtbl -> CreateCookie(This,name,value,domain,path,cookie) ) 

#define ICoreWebView2CookieManager_CopyCookie(This,cookieParam,cookie)	\
    ( (This)->lpVtbl -> CopyCookie(This,cookieParam,cookie) ) 

#define ICoreWebView2CookieManager_GetCookies(This,uri,handler)	\
    ( (This)->lpVtbl -> GetCookies(This,uri,handler) ) 

#define ICoreWebView2CookieManager_AddOrUpdateCookie(This,cookie)	\
    ( (This)->lpVtbl -> AddOrUpdateCookie(This,cookie) ) 

#define ICoreWebView2CookieManager_DeleteCookie(This,cookie)	\
    ( (This)->lpVtbl -> DeleteCookie(This,cookie) ) 

#define ICoreWebView2CookieManager_DeleteCookies(This,name,uri)	\
    ( (This)->lpVtbl -> DeleteCookies(This,name,uri) ) 

#define ICoreWebView2CookieManager_DeleteCookiesWithDomainAndPath(This,name,domain,path)	\
    ( (This)->lpVtbl -> DeleteCookiesWithDomainAndPath(This,name,domain,path) ) 

#define ICoreWebView2CookieManager_DeleteAllCookies(This)	\
    ( (This)->lpVtbl -> DeleteAllCookies(This) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CookieManager_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler = {0x02fab84b,0x1428,0x4fb7,{0xad,0x45,0x1b,0x2e,0x64,0x73,0x61,0x84}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("02fab84b-1428-4fb7-ad45-1b2e64736184")
    ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            HRESULT errorCode,
            ICoreWebView2CompositionController *webView) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler * This,
            HRESULT errorCode,
            ICoreWebView2CompositionController *webView);
        
        END_INTERFACE
    } ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl;

    interface ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_Invoke(This,errorCode,webView)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,webView) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CreateCoreWebView2ControllerCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler = {0x6c4819f3,0xc9b7,0x4260,{0x81,0x27,0xc9,0xf5,0xbd,0xe7,0xf6,0x8c}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("6c4819f3-c9b7-4260-8127-c9f5bde7f68c")
    ICoreWebView2CreateCoreWebView2ControllerCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            HRESULT errorCode,
            ICoreWebView2Controller *createdController) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler * This,
            HRESULT errorCode,
            ICoreWebView2Controller *createdController);
        
        END_INTERFACE
    } ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl;

    interface ICoreWebView2CreateCoreWebView2ControllerCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_Invoke(This,errorCode,createdController)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,createdController) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler = {0x4e8a3389,0xc9d8,0x4bd2,{0xb6,0xb5,0x12,0x4f,0xee,0x6c,0xc1,0x4d}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("4e8a3389-c9d8-4bd2-b6b5-124fee6cc14d")
    ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            HRESULT errorCode,
            ICoreWebView2Environment *createdEnvironment) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler * This,
            HRESULT errorCode,
            ICoreWebView2Environment *createdEnvironment);
        
        END_INTERFACE
    } ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl;

    interface ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_Invoke(This,errorCode,createdEnvironment)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,createdEnvironment) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ContainsFullScreenElementChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ContainsFullScreenElementChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ContainsFullScreenElementChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ContainsFullScreenElementChangedEventHandler = {0xe45d98b1,0xafef,0x45be,{0x8b,0xaf,0x6c,0x77,0x28,0x86,0x7f,0x73}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("e45d98b1-afef-45be-8baf-6c7728867f73")
    ICoreWebView2ContainsFullScreenElementChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ContainsFullScreenElementChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ContainsFullScreenElementChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ContainsFullScreenElementChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ContainsFullScreenElementChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2ContainsFullScreenElementChangedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2ContainsFullScreenElementChangedEventHandlerVtbl;

    interface ICoreWebView2ContainsFullScreenElementChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2ContainsFullScreenElementChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ContainsFullScreenElementChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ContainsFullScreenElementChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ContainsFullScreenElementChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ContainsFullScreenElementChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ContainsFullScreenElementChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CursorChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CursorChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CursorChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CursorChangedEventHandler = {0x9da43ccc,0x26e1,0x4dad,{0xb5,0x6c,0xd8,0x96,0x1c,0x94,0xc5,0x71}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("9da43ccc-26e1-4dad-b56c-d8961c94c571")
    ICoreWebView2CursorChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2CompositionController *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CursorChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CursorChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CursorChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CursorChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2CursorChangedEventHandler * This,
            /* [in] */ ICoreWebView2CompositionController *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2CursorChangedEventHandlerVtbl;

    interface ICoreWebView2CursorChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2CursorChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CursorChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CursorChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CursorChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CursorChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CursorChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2DocumentTitleChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2DocumentTitleChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2DocumentTitleChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2DocumentTitleChangedEventHandler = {0xf5f2b923,0x953e,0x4042,{0x9f,0x95,0xf3,0xa1,0x18,0xe1,0xaf,0xd4}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("f5f2b923-953e-4042-9f95-f3a118e1afd4")
    ICoreWebView2DocumentTitleChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2DocumentTitleChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2DocumentTitleChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2DocumentTitleChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2DocumentTitleChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2DocumentTitleChangedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2DocumentTitleChangedEventHandlerVtbl;

    interface ICoreWebView2DocumentTitleChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2DocumentTitleChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2DocumentTitleChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2DocumentTitleChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2DocumentTitleChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2DocumentTitleChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2DocumentTitleChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2DOMContentLoadedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2DOMContentLoadedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2DOMContentLoadedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2DOMContentLoadedEventArgs = {0x16B1E21A,0xC503,0x44F2,{0x84,0xC9,0x70,0xAB,0xA5,0x03,0x12,0x83}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("16B1E21A-C503-44F2-84C9-70ABA5031283")
    ICoreWebView2DOMContentLoadedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_NavigationId( 
            /* [retval][out] */ UINT64 *navigationId) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2DOMContentLoadedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2DOMContentLoadedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2DOMContentLoadedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2DOMContentLoadedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_NavigationId )( 
            ICoreWebView2DOMContentLoadedEventArgs * This,
            /* [retval][out] */ UINT64 *navigationId);
        
        END_INTERFACE
    } ICoreWebView2DOMContentLoadedEventArgsVtbl;

    interface ICoreWebView2DOMContentLoadedEventArgs
    {
        CONST_VTBL struct ICoreWebView2DOMContentLoadedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2DOMContentLoadedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2DOMContentLoadedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2DOMContentLoadedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2DOMContentLoadedEventArgs_get_NavigationId(This,navigationId)	\
    ( (This)->lpVtbl -> get_NavigationId(This,navigationId) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2DOMContentLoadedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2DOMContentLoadedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2DOMContentLoadedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2DOMContentLoadedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2DOMContentLoadedEventHandler = {0x4BAC7E9C,0x199E,0x49ED,{0x87,0xED,0x24,0x93,0x03,0xAC,0xF0,0x19}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("4BAC7E9C-199E-49ED-87ED-249303ACF019")
    ICoreWebView2DOMContentLoadedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2DOMContentLoadedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2DOMContentLoadedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2DOMContentLoadedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2DOMContentLoadedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2DOMContentLoadedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2DOMContentLoadedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2DOMContentLoadedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2DOMContentLoadedEventHandlerVtbl;

    interface ICoreWebView2DOMContentLoadedEventHandler
    {
        CONST_VTBL struct ICoreWebView2DOMContentLoadedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2DOMContentLoadedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2DOMContentLoadedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2DOMContentLoadedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2DOMContentLoadedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2DOMContentLoadedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Deferral_INTERFACE_DEFINED__
#define __ICoreWebView2Deferral_INTERFACE_DEFINED__

/* interface ICoreWebView2Deferral */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Deferral = {0xc10e7f7b,0xb585,0x46f0,{0xa6,0x23,0x8b,0xef,0xbf,0x3e,0x4e,0xe0}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("c10e7f7b-b585-46f0-a623-8befbf3e4ee0")
    ICoreWebView2Deferral : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Complete( void) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2DeferralVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Deferral * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Deferral * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Deferral * This);
        
        HRESULT ( STDMETHODCALLTYPE *Complete )( 
            ICoreWebView2Deferral * This);
        
        END_INTERFACE
    } ICoreWebView2DeferralVtbl;

    interface ICoreWebView2Deferral
    {
        CONST_VTBL struct ICoreWebView2DeferralVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Deferral_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Deferral_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Deferral_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Deferral_Complete(This)	\
    ( (This)->lpVtbl -> Complete(This) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Deferral_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2DevToolsProtocolEventReceivedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2DevToolsProtocolEventReceivedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2DevToolsProtocolEventReceivedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2DevToolsProtocolEventReceivedEventArgs = {0x653c2959,0xbb3a,0x4377,{0x86,0x32,0xb5,0x8a,0xda,0x4e,0x66,0xc4}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("653c2959-bb3a-4377-8632-b58ada4e66c4")
    ICoreWebView2DevToolsProtocolEventReceivedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ParameterObjectAsJson( 
            /* [retval][out] */ LPWSTR *parameterObjectAsJson) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2DevToolsProtocolEventReceivedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ParameterObjectAsJson )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventArgs * This,
            /* [retval][out] */ LPWSTR *parameterObjectAsJson);
        
        END_INTERFACE
    } ICoreWebView2DevToolsProtocolEventReceivedEventArgsVtbl;

    interface ICoreWebView2DevToolsProtocolEventReceivedEventArgs
    {
        CONST_VTBL struct ICoreWebView2DevToolsProtocolEventReceivedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2DevToolsProtocolEventReceivedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2DevToolsProtocolEventReceivedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2DevToolsProtocolEventReceivedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2DevToolsProtocolEventReceivedEventArgs_get_ParameterObjectAsJson(This,parameterObjectAsJson)	\
    ( (This)->lpVtbl -> get_ParameterObjectAsJson(This,parameterObjectAsJson) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2DevToolsProtocolEventReceivedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2DevToolsProtocolEventReceivedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2DevToolsProtocolEventReceivedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2DevToolsProtocolEventReceivedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2DevToolsProtocolEventReceivedEventHandler = {0xe2fda4be,0x5456,0x406c,{0xa2,0x61,0x3d,0x45,0x21,0x38,0x36,0x2c}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("e2fda4be-5456-406c-a261-3d452138362c")
    ICoreWebView2DevToolsProtocolEventReceivedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2DevToolsProtocolEventReceivedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2DevToolsProtocolEventReceivedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2DevToolsProtocolEventReceivedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl;

    interface ICoreWebView2DevToolsProtocolEventReceivedEventHandler
    {
        CONST_VTBL struct ICoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2DevToolsProtocolEventReceivedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2DevToolsProtocolEventReceivedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2DevToolsProtocolEventReceivedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2DevToolsProtocolEventReceivedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2DevToolsProtocolEventReceivedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2DevToolsProtocolEventReceiver_INTERFACE_DEFINED__
#define __ICoreWebView2DevToolsProtocolEventReceiver_INTERFACE_DEFINED__

/* interface ICoreWebView2DevToolsProtocolEventReceiver */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2DevToolsProtocolEventReceiver = {0xb32ca51a,0x8371,0x45e9,{0x93,0x17,0xaf,0x02,0x1d,0x08,0x03,0x67}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("b32ca51a-8371-45e9-9317-af021d080367")
    ICoreWebView2DevToolsProtocolEventReceiver : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE add_DevToolsProtocolEventReceived( 
            /* [in] */ ICoreWebView2DevToolsProtocolEventReceivedEventHandler *handler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_DevToolsProtocolEventReceived( 
            /* [in] */ EventRegistrationToken token) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2DevToolsProtocolEventReceiverVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2DevToolsProtocolEventReceiver * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2DevToolsProtocolEventReceiver * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2DevToolsProtocolEventReceiver * This);
        
        HRESULT ( STDMETHODCALLTYPE *add_DevToolsProtocolEventReceived )( 
            ICoreWebView2DevToolsProtocolEventReceiver * This,
            /* [in] */ ICoreWebView2DevToolsProtocolEventReceivedEventHandler *handler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_DevToolsProtocolEventReceived )( 
            ICoreWebView2DevToolsProtocolEventReceiver * This,
            /* [in] */ EventRegistrationToken token);
        
        END_INTERFACE
    } ICoreWebView2DevToolsProtocolEventReceiverVtbl;

    interface ICoreWebView2DevToolsProtocolEventReceiver
    {
        CONST_VTBL struct ICoreWebView2DevToolsProtocolEventReceiverVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2DevToolsProtocolEventReceiver_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2DevToolsProtocolEventReceiver_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2DevToolsProtocolEventReceiver_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2DevToolsProtocolEventReceiver_add_DevToolsProtocolEventReceived(This,handler,token)	\
    ( (This)->lpVtbl -> add_DevToolsProtocolEventReceived(This,handler,token) ) 

#define ICoreWebView2DevToolsProtocolEventReceiver_remove_DevToolsProtocolEventReceived(This,token)	\
    ( (This)->lpVtbl -> remove_DevToolsProtocolEventReceived(This,token) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2DevToolsProtocolEventReceiver_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Environment_INTERFACE_DEFINED__
#define __ICoreWebView2Environment_INTERFACE_DEFINED__

/* interface ICoreWebView2Environment */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Environment = {0xb96d755e,0x0319,0x4e92,{0xa2,0x96,0x23,0x43,0x6f,0x46,0xa1,0xfc}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("b96d755e-0319-4e92-a296-23436f46a1fc")
    ICoreWebView2Environment : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE CreateCoreWebView2Controller( 
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler *handler) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE CreateWebResourceResponse( 
            /* [in] */ IStream *content,
            /* [in] */ int statusCode,
            /* [in] */ LPCWSTR reasonPhrase,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceResponse **response) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_BrowserVersionString( 
            /* [retval][out] */ LPWSTR *versionInfo) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE add_NewBrowserVersionAvailable( 
            /* [in] */ ICoreWebView2NewBrowserVersionAvailableEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE remove_NewBrowserVersionAvailable( 
            /* [in] */ EventRegistrationToken token) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2EnvironmentVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Environment * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Environment * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Environment * This);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2Controller )( 
            ICoreWebView2Environment * This,
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CreateWebResourceResponse )( 
            ICoreWebView2Environment * This,
            /* [in] */ IStream *content,
            /* [in] */ int statusCode,
            /* [in] */ LPCWSTR reasonPhrase,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceResponse **response);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BrowserVersionString )( 
            ICoreWebView2Environment * This,
            /* [retval][out] */ LPWSTR *versionInfo);
        
        HRESULT ( STDMETHODCALLTYPE *add_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment * This,
            /* [in] */ ICoreWebView2NewBrowserVersionAvailableEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment * This,
            /* [in] */ EventRegistrationToken token);
        
        END_INTERFACE
    } ICoreWebView2EnvironmentVtbl;

    interface ICoreWebView2Environment
    {
        CONST_VTBL struct ICoreWebView2EnvironmentVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Environment_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Environment_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Environment_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Environment_CreateCoreWebView2Controller(This,parentWindow,handler)	\
    ( (This)->lpVtbl -> CreateCoreWebView2Controller(This,parentWindow,handler) ) 

#define ICoreWebView2Environment_CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response)	\
    ( (This)->lpVtbl -> CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response) ) 

#define ICoreWebView2Environment_get_BrowserVersionString(This,versionInfo)	\
    ( (This)->lpVtbl -> get_BrowserVersionString(This,versionInfo) ) 

#define ICoreWebView2Environment_add_NewBrowserVersionAvailable(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NewBrowserVersionAvailable(This,eventHandler,token) ) 

#define ICoreWebView2Environment_remove_NewBrowserVersionAvailable(This,token)	\
    ( (This)->lpVtbl -> remove_NewBrowserVersionAvailable(This,token) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Environment_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Environment2_INTERFACE_DEFINED__
#define __ICoreWebView2Environment2_INTERFACE_DEFINED__

/* interface ICoreWebView2Environment2 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Environment2 = {0x41F3632B,0x5EF4,0x404F,{0xAD,0x82,0x2D,0x60,0x6C,0x5A,0x9A,0x21}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("41F3632B-5EF4-404F-AD82-2D606C5A9A21")
    ICoreWebView2Environment2 : public ICoreWebView2Environment
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE CreateWebResourceRequest( 
            /* [in] */ LPCWSTR uri,
            /* [in] */ LPCWSTR method,
            /* [in] */ IStream *postData,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Environment2Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Environment2 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Environment2 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Environment2 * This);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2Controller )( 
            ICoreWebView2Environment2 * This,
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CreateWebResourceResponse )( 
            ICoreWebView2Environment2 * This,
            /* [in] */ IStream *content,
            /* [in] */ int statusCode,
            /* [in] */ LPCWSTR reasonPhrase,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceResponse **response);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BrowserVersionString )( 
            ICoreWebView2Environment2 * This,
            /* [retval][out] */ LPWSTR *versionInfo);
        
        HRESULT ( STDMETHODCALLTYPE *add_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment2 * This,
            /* [in] */ ICoreWebView2NewBrowserVersionAvailableEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment2 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *CreateWebResourceRequest )( 
            ICoreWebView2Environment2 * This,
            /* [in] */ LPCWSTR uri,
            /* [in] */ LPCWSTR method,
            /* [in] */ IStream *postData,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request);
        
        END_INTERFACE
    } ICoreWebView2Environment2Vtbl;

    interface ICoreWebView2Environment2
    {
        CONST_VTBL struct ICoreWebView2Environment2Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Environment2_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Environment2_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Environment2_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Environment2_CreateCoreWebView2Controller(This,parentWindow,handler)	\
    ( (This)->lpVtbl -> CreateCoreWebView2Controller(This,parentWindow,handler) ) 

#define ICoreWebView2Environment2_CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response)	\
    ( (This)->lpVtbl -> CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response) ) 

#define ICoreWebView2Environment2_get_BrowserVersionString(This,versionInfo)	\
    ( (This)->lpVtbl -> get_BrowserVersionString(This,versionInfo) ) 

#define ICoreWebView2Environment2_add_NewBrowserVersionAvailable(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NewBrowserVersionAvailable(This,eventHandler,token) ) 

#define ICoreWebView2Environment2_remove_NewBrowserVersionAvailable(This,token)	\
    ( (This)->lpVtbl -> remove_NewBrowserVersionAvailable(This,token) ) 


#define ICoreWebView2Environment2_CreateWebResourceRequest(This,uri,method,postData,headers,request)	\
    ( (This)->lpVtbl -> CreateWebResourceRequest(This,uri,method,postData,headers,request) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Environment2_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Environment3_INTERFACE_DEFINED__
#define __ICoreWebView2Environment3_INTERFACE_DEFINED__

/* interface ICoreWebView2Environment3 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Environment3 = {0x80a22ae3,0xbe7c,0x4ce2,{0xaf,0xe1,0x5a,0x50,0x05,0x6c,0xde,0xeb}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("80a22ae3-be7c-4ce2-afe1-5a50056cdeeb")
    ICoreWebView2Environment3 : public ICoreWebView2Environment2
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE CreateCoreWebView2CompositionController( 
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler *handler) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE CreateCoreWebView2PointerInfo( 
            /* [retval][out] */ ICoreWebView2PointerInfo **pointerInfo) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Environment3Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Environment3 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Environment3 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Environment3 * This);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2Controller )( 
            ICoreWebView2Environment3 * This,
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CreateWebResourceResponse )( 
            ICoreWebView2Environment3 * This,
            /* [in] */ IStream *content,
            /* [in] */ int statusCode,
            /* [in] */ LPCWSTR reasonPhrase,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceResponse **response);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BrowserVersionString )( 
            ICoreWebView2Environment3 * This,
            /* [retval][out] */ LPWSTR *versionInfo);
        
        HRESULT ( STDMETHODCALLTYPE *add_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment3 * This,
            /* [in] */ ICoreWebView2NewBrowserVersionAvailableEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment3 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *CreateWebResourceRequest )( 
            ICoreWebView2Environment3 * This,
            /* [in] */ LPCWSTR uri,
            /* [in] */ LPCWSTR method,
            /* [in] */ IStream *postData,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2CompositionController )( 
            ICoreWebView2Environment3 * This,
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2PointerInfo )( 
            ICoreWebView2Environment3 * This,
            /* [retval][out] */ ICoreWebView2PointerInfo **pointerInfo);
        
        END_INTERFACE
    } ICoreWebView2Environment3Vtbl;

    interface ICoreWebView2Environment3
    {
        CONST_VTBL struct ICoreWebView2Environment3Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Environment3_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Environment3_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Environment3_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Environment3_CreateCoreWebView2Controller(This,parentWindow,handler)	\
    ( (This)->lpVtbl -> CreateCoreWebView2Controller(This,parentWindow,handler) ) 

#define ICoreWebView2Environment3_CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response)	\
    ( (This)->lpVtbl -> CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response) ) 

#define ICoreWebView2Environment3_get_BrowserVersionString(This,versionInfo)	\
    ( (This)->lpVtbl -> get_BrowserVersionString(This,versionInfo) ) 

#define ICoreWebView2Environment3_add_NewBrowserVersionAvailable(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NewBrowserVersionAvailable(This,eventHandler,token) ) 

#define ICoreWebView2Environment3_remove_NewBrowserVersionAvailable(This,token)	\
    ( (This)->lpVtbl -> remove_NewBrowserVersionAvailable(This,token) ) 


#define ICoreWebView2Environment3_CreateWebResourceRequest(This,uri,method,postData,headers,request)	\
    ( (This)->lpVtbl -> CreateWebResourceRequest(This,uri,method,postData,headers,request) ) 


#define ICoreWebView2Environment3_CreateCoreWebView2CompositionController(This,parentWindow,handler)	\
    ( (This)->lpVtbl -> CreateCoreWebView2CompositionController(This,parentWindow,handler) ) 

#define ICoreWebView2Environment3_CreateCoreWebView2PointerInfo(This,pointerInfo)	\
    ( (This)->lpVtbl -> CreateCoreWebView2PointerInfo(This,pointerInfo) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Environment3_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Environment4_INTERFACE_DEFINED__
#define __ICoreWebView2Environment4_INTERFACE_DEFINED__

/* interface ICoreWebView2Environment4 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Environment4 = {0x20944379,0x6dcf,0x41d6,{0xa0,0xa0,0xab,0xc0,0xfc,0x50,0xde,0x0d}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("20944379-6dcf-41d6-a0a0-abc0fc50de0d")
    ICoreWebView2Environment4 : public ICoreWebView2Environment3
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE GetProviderForHwnd( 
            /* [in] */ HWND hwnd,
            /* [retval][out] */ IUnknown **provider) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Environment4Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Environment4 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Environment4 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Environment4 * This);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2Controller )( 
            ICoreWebView2Environment4 * This,
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2ControllerCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CreateWebResourceResponse )( 
            ICoreWebView2Environment4 * This,
            /* [in] */ IStream *content,
            /* [in] */ int statusCode,
            /* [in] */ LPCWSTR reasonPhrase,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceResponse **response);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_BrowserVersionString )( 
            ICoreWebView2Environment4 * This,
            /* [retval][out] */ LPWSTR *versionInfo);
        
        HRESULT ( STDMETHODCALLTYPE *add_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment4 * This,
            /* [in] */ ICoreWebView2NewBrowserVersionAvailableEventHandler *eventHandler,
            /* [out] */ EventRegistrationToken *token);
        
        HRESULT ( STDMETHODCALLTYPE *remove_NewBrowserVersionAvailable )( 
            ICoreWebView2Environment4 * This,
            /* [in] */ EventRegistrationToken token);
        
        HRESULT ( STDMETHODCALLTYPE *CreateWebResourceRequest )( 
            ICoreWebView2Environment4 * This,
            /* [in] */ LPCWSTR uri,
            /* [in] */ LPCWSTR method,
            /* [in] */ IStream *postData,
            /* [in] */ LPCWSTR headers,
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2CompositionController )( 
            ICoreWebView2Environment4 * This,
            HWND parentWindow,
            ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler *handler);
        
        HRESULT ( STDMETHODCALLTYPE *CreateCoreWebView2PointerInfo )( 
            ICoreWebView2Environment4 * This,
            /* [retval][out] */ ICoreWebView2PointerInfo **pointerInfo);
        
        HRESULT ( STDMETHODCALLTYPE *GetProviderForHwnd )( 
            ICoreWebView2Environment4 * This,
            /* [in] */ HWND hwnd,
            /* [retval][out] */ IUnknown **provider);
        
        END_INTERFACE
    } ICoreWebView2Environment4Vtbl;

    interface ICoreWebView2Environment4
    {
        CONST_VTBL struct ICoreWebView2Environment4Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Environment4_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Environment4_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Environment4_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Environment4_CreateCoreWebView2Controller(This,parentWindow,handler)	\
    ( (This)->lpVtbl -> CreateCoreWebView2Controller(This,parentWindow,handler) ) 

#define ICoreWebView2Environment4_CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response)	\
    ( (This)->lpVtbl -> CreateWebResourceResponse(This,content,statusCode,reasonPhrase,headers,response) ) 

#define ICoreWebView2Environment4_get_BrowserVersionString(This,versionInfo)	\
    ( (This)->lpVtbl -> get_BrowserVersionString(This,versionInfo) ) 

#define ICoreWebView2Environment4_add_NewBrowserVersionAvailable(This,eventHandler,token)	\
    ( (This)->lpVtbl -> add_NewBrowserVersionAvailable(This,eventHandler,token) ) 

#define ICoreWebView2Environment4_remove_NewBrowserVersionAvailable(This,token)	\
    ( (This)->lpVtbl -> remove_NewBrowserVersionAvailable(This,token) ) 


#define ICoreWebView2Environment4_CreateWebResourceRequest(This,uri,method,postData,headers,request)	\
    ( (This)->lpVtbl -> CreateWebResourceRequest(This,uri,method,postData,headers,request) ) 


#define ICoreWebView2Environment4_CreateCoreWebView2CompositionController(This,parentWindow,handler)	\
    ( (This)->lpVtbl -> CreateCoreWebView2CompositionController(This,parentWindow,handler) ) 

#define ICoreWebView2Environment4_CreateCoreWebView2PointerInfo(This,pointerInfo)	\
    ( (This)->lpVtbl -> CreateCoreWebView2PointerInfo(This,pointerInfo) ) 


#define ICoreWebView2Environment4_GetProviderForHwnd(This,hwnd,provider)	\
    ( (This)->lpVtbl -> GetProviderForHwnd(This,hwnd,provider) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Environment4_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2EnvironmentOptions_INTERFACE_DEFINED__
#define __ICoreWebView2EnvironmentOptions_INTERFACE_DEFINED__

/* interface ICoreWebView2EnvironmentOptions */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2EnvironmentOptions = {0x2fde08a8,0x1e9a,0x4766,{0x8c,0x05,0x95,0xa9,0xce,0xb9,0xd1,0xc5}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("2fde08a8-1e9a-4766-8c05-95a9ceb9d1c5")
    ICoreWebView2EnvironmentOptions : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_AdditionalBrowserArguments( 
            /* [retval][out] */ LPWSTR *value) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_AdditionalBrowserArguments( 
            /* [in] */ LPCWSTR value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Language( 
            /* [retval][out] */ LPWSTR *value) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Language( 
            /* [in] */ LPCWSTR value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_TargetCompatibleBrowserVersion( 
            /* [retval][out] */ LPWSTR *value) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_TargetCompatibleBrowserVersion( 
            /* [in] */ LPCWSTR value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_AllowSingleSignOnUsingOSPrimaryAccount( 
            /* [retval][out] */ BOOL *allow) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_AllowSingleSignOnUsingOSPrimaryAccount( 
            /* [in] */ BOOL allow) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2EnvironmentOptionsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2EnvironmentOptions * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2EnvironmentOptions * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AdditionalBrowserArguments )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [retval][out] */ LPWSTR *value);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AdditionalBrowserArguments )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [in] */ LPCWSTR value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Language )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [retval][out] */ LPWSTR *value);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Language )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [in] */ LPCWSTR value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_TargetCompatibleBrowserVersion )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [retval][out] */ LPWSTR *value);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_TargetCompatibleBrowserVersion )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [in] */ LPCWSTR value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AllowSingleSignOnUsingOSPrimaryAccount )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [retval][out] */ BOOL *allow);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AllowSingleSignOnUsingOSPrimaryAccount )( 
            ICoreWebView2EnvironmentOptions * This,
            /* [in] */ BOOL allow);
        
        END_INTERFACE
    } ICoreWebView2EnvironmentOptionsVtbl;

    interface ICoreWebView2EnvironmentOptions
    {
        CONST_VTBL struct ICoreWebView2EnvironmentOptionsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2EnvironmentOptions_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2EnvironmentOptions_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2EnvironmentOptions_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2EnvironmentOptions_get_AdditionalBrowserArguments(This,value)	\
    ( (This)->lpVtbl -> get_AdditionalBrowserArguments(This,value) ) 

#define ICoreWebView2EnvironmentOptions_put_AdditionalBrowserArguments(This,value)	\
    ( (This)->lpVtbl -> put_AdditionalBrowserArguments(This,value) ) 

#define ICoreWebView2EnvironmentOptions_get_Language(This,value)	\
    ( (This)->lpVtbl -> get_Language(This,value) ) 

#define ICoreWebView2EnvironmentOptions_put_Language(This,value)	\
    ( (This)->lpVtbl -> put_Language(This,value) ) 

#define ICoreWebView2EnvironmentOptions_get_TargetCompatibleBrowserVersion(This,value)	\
    ( (This)->lpVtbl -> get_TargetCompatibleBrowserVersion(This,value) ) 

#define ICoreWebView2EnvironmentOptions_put_TargetCompatibleBrowserVersion(This,value)	\
    ( (This)->lpVtbl -> put_TargetCompatibleBrowserVersion(This,value) ) 

#define ICoreWebView2EnvironmentOptions_get_AllowSingleSignOnUsingOSPrimaryAccount(This,allow)	\
    ( (This)->lpVtbl -> get_AllowSingleSignOnUsingOSPrimaryAccount(This,allow) ) 

#define ICoreWebView2EnvironmentOptions_put_AllowSingleSignOnUsingOSPrimaryAccount(This,allow)	\
    ( (This)->lpVtbl -> put_AllowSingleSignOnUsingOSPrimaryAccount(This,allow) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2EnvironmentOptions_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ExecuteScriptCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ExecuteScriptCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ExecuteScriptCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ExecuteScriptCompletedHandler = {0x49511172,0xcc67,0x4bca,{0x99,0x23,0x13,0x71,0x12,0xf4,0xc4,0xcc}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("49511172-cc67-4bca-9923-137112f4c4cc")
    ICoreWebView2ExecuteScriptCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ HRESULT errorCode,
            /* [in] */ LPCWSTR resultObjectAsJson) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ExecuteScriptCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ExecuteScriptCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ExecuteScriptCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ExecuteScriptCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2ExecuteScriptCompletedHandler * This,
            /* [in] */ HRESULT errorCode,
            /* [in] */ LPCWSTR resultObjectAsJson);
        
        END_INTERFACE
    } ICoreWebView2ExecuteScriptCompletedHandlerVtbl;

    interface ICoreWebView2ExecuteScriptCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2ExecuteScriptCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ExecuteScriptCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ExecuteScriptCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ExecuteScriptCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ExecuteScriptCompletedHandler_Invoke(This,errorCode,resultObjectAsJson)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,resultObjectAsJson) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ExecuteScriptCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2FrameInfo_INTERFACE_DEFINED__
#define __ICoreWebView2FrameInfo_INTERFACE_DEFINED__

/* interface ICoreWebView2FrameInfo */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2FrameInfo = {0xda86b8a1,0xbdf3,0x4f11,{0x99,0x55,0x52,0x8c,0xef,0xa5,0x97,0x27}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("da86b8a1-bdf3-4f11-9955-528cefa59727")
    ICoreWebView2FrameInfo : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Name( 
            /* [retval][out] */ LPWSTR *name) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Source( 
            /* [retval][out] */ LPWSTR *source) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2FrameInfoVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2FrameInfo * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2FrameInfo * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2FrameInfo * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Name )( 
            ICoreWebView2FrameInfo * This,
            /* [retval][out] */ LPWSTR *name);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Source )( 
            ICoreWebView2FrameInfo * This,
            /* [retval][out] */ LPWSTR *source);
        
        END_INTERFACE
    } ICoreWebView2FrameInfoVtbl;

    interface ICoreWebView2FrameInfo
    {
        CONST_VTBL struct ICoreWebView2FrameInfoVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2FrameInfo_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2FrameInfo_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2FrameInfo_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2FrameInfo_get_Name(This,name)	\
    ( (This)->lpVtbl -> get_Name(This,name) ) 

#define ICoreWebView2FrameInfo_get_Source(This,source)	\
    ( (This)->lpVtbl -> get_Source(This,source) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2FrameInfo_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2FrameInfoCollection_INTERFACE_DEFINED__
#define __ICoreWebView2FrameInfoCollection_INTERFACE_DEFINED__

/* interface ICoreWebView2FrameInfoCollection */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2FrameInfoCollection = {0x8f834154,0xd38e,0x4d90,{0xaf,0xfb,0x68,0x00,0xa7,0x27,0x28,0x39}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("8f834154-d38e-4d90-affb-6800a7272839")
    ICoreWebView2FrameInfoCollection : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE GetIterator( 
            /* [retval][out] */ ICoreWebView2FrameInfoCollectionIterator **iterator) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2FrameInfoCollectionVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2FrameInfoCollection * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2FrameInfoCollection * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2FrameInfoCollection * This);
        
        HRESULT ( STDMETHODCALLTYPE *GetIterator )( 
            ICoreWebView2FrameInfoCollection * This,
            /* [retval][out] */ ICoreWebView2FrameInfoCollectionIterator **iterator);
        
        END_INTERFACE
    } ICoreWebView2FrameInfoCollectionVtbl;

    interface ICoreWebView2FrameInfoCollection
    {
        CONST_VTBL struct ICoreWebView2FrameInfoCollectionVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2FrameInfoCollection_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2FrameInfoCollection_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2FrameInfoCollection_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2FrameInfoCollection_GetIterator(This,iterator)	\
    ( (This)->lpVtbl -> GetIterator(This,iterator) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2FrameInfoCollection_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2FrameInfoCollectionIterator_INTERFACE_DEFINED__
#define __ICoreWebView2FrameInfoCollectionIterator_INTERFACE_DEFINED__

/* interface ICoreWebView2FrameInfoCollectionIterator */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2FrameInfoCollectionIterator = {0x1bf89e2d,0x1b2b,0x4629,{0xb2,0x8f,0x05,0x09,0x9b,0x41,0xbb,0x03}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("1bf89e2d-1b2b-4629-b28f-05099b41bb03")
    ICoreWebView2FrameInfoCollectionIterator : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HasCurrent( 
            /* [retval][out] */ BOOL *hasCurrent) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetCurrent( 
            /* [retval][out] */ ICoreWebView2FrameInfo **frameInfo) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE MoveNext( 
            /* [retval][out] */ BOOL *hasNext) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2FrameInfoCollectionIteratorVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2FrameInfoCollectionIterator * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2FrameInfoCollectionIterator * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2FrameInfoCollectionIterator * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_HasCurrent )( 
            ICoreWebView2FrameInfoCollectionIterator * This,
            /* [retval][out] */ BOOL *hasCurrent);
        
        HRESULT ( STDMETHODCALLTYPE *GetCurrent )( 
            ICoreWebView2FrameInfoCollectionIterator * This,
            /* [retval][out] */ ICoreWebView2FrameInfo **frameInfo);
        
        HRESULT ( STDMETHODCALLTYPE *MoveNext )( 
            ICoreWebView2FrameInfoCollectionIterator * This,
            /* [retval][out] */ BOOL *hasNext);
        
        END_INTERFACE
    } ICoreWebView2FrameInfoCollectionIteratorVtbl;

    interface ICoreWebView2FrameInfoCollectionIterator
    {
        CONST_VTBL struct ICoreWebView2FrameInfoCollectionIteratorVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2FrameInfoCollectionIterator_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2FrameInfoCollectionIterator_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2FrameInfoCollectionIterator_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2FrameInfoCollectionIterator_get_HasCurrent(This,hasCurrent)	\
    ( (This)->lpVtbl -> get_HasCurrent(This,hasCurrent) ) 

#define ICoreWebView2FrameInfoCollectionIterator_GetCurrent(This,frameInfo)	\
    ( (This)->lpVtbl -> GetCurrent(This,frameInfo) ) 

#define ICoreWebView2FrameInfoCollectionIterator_MoveNext(This,hasNext)	\
    ( (This)->lpVtbl -> MoveNext(This,hasNext) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2FrameInfoCollectionIterator_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2FocusChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2FocusChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2FocusChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2FocusChangedEventHandler = {0x05ea24bd,0x6452,0x4926,{0x90,0x14,0x4b,0x82,0xb4,0x98,0x13,0x5d}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("05ea24bd-6452-4926-9014-4b82b498135d")
    ICoreWebView2FocusChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2FocusChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2FocusChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2FocusChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2FocusChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2FocusChangedEventHandler * This,
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2FocusChangedEventHandlerVtbl;

    interface ICoreWebView2FocusChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2FocusChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2FocusChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2FocusChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2FocusChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2FocusChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2FocusChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2GetCookiesCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2GetCookiesCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2GetCookiesCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2GetCookiesCompletedHandler = {0x5A4F5069,0x5C15,0x47C3,{0x86,0x46,0xF4,0xDE,0x1C,0x11,0x66,0x70}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("5A4F5069-5C15-47C3-8646-F4DE1C116670")
    ICoreWebView2GetCookiesCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            HRESULT result,
            ICoreWebView2CookieList *cookieList) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2GetCookiesCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2GetCookiesCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2GetCookiesCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2GetCookiesCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2GetCookiesCompletedHandler * This,
            HRESULT result,
            ICoreWebView2CookieList *cookieList);
        
        END_INTERFACE
    } ICoreWebView2GetCookiesCompletedHandlerVtbl;

    interface ICoreWebView2GetCookiesCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2GetCookiesCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2GetCookiesCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2GetCookiesCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2GetCookiesCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2GetCookiesCompletedHandler_Invoke(This,result,cookieList)	\
    ( (This)->lpVtbl -> Invoke(This,result,cookieList) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2GetCookiesCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2HistoryChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2HistoryChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2HistoryChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2HistoryChangedEventHandler = {0xc79a420c,0xefd9,0x4058,{0x92,0x95,0x3e,0x8b,0x4b,0xca,0xb6,0x45}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("c79a420c-efd9-4058-9295-3e8b4bcab645")
    ICoreWebView2HistoryChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2HistoryChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2HistoryChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2HistoryChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2HistoryChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2HistoryChangedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2HistoryChangedEventHandlerVtbl;

    interface ICoreWebView2HistoryChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2HistoryChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2HistoryChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2HistoryChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2HistoryChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2HistoryChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2HistoryChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2HttpHeadersCollectionIterator_INTERFACE_DEFINED__
#define __ICoreWebView2HttpHeadersCollectionIterator_INTERFACE_DEFINED__

/* interface ICoreWebView2HttpHeadersCollectionIterator */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2HttpHeadersCollectionIterator = {0x0702fc30,0xf43b,0x47bb,{0xab,0x52,0xa4,0x2c,0xb5,0x52,0xad,0x9f}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("0702fc30-f43b-47bb-ab52-a42cb552ad9f")
    ICoreWebView2HttpHeadersCollectionIterator : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE GetCurrentHeader( 
            /* [out] */ LPWSTR *name,
            /* [out] */ LPWSTR *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HasCurrentHeader( 
            /* [retval][out] */ BOOL *hasCurrent) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE MoveNext( 
            /* [retval][out] */ BOOL *hasNext) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2HttpHeadersCollectionIteratorVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2HttpHeadersCollectionIterator * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2HttpHeadersCollectionIterator * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2HttpHeadersCollectionIterator * This);
        
        HRESULT ( STDMETHODCALLTYPE *GetCurrentHeader )( 
            ICoreWebView2HttpHeadersCollectionIterator * This,
            /* [out] */ LPWSTR *name,
            /* [out] */ LPWSTR *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_HasCurrentHeader )( 
            ICoreWebView2HttpHeadersCollectionIterator * This,
            /* [retval][out] */ BOOL *hasCurrent);
        
        HRESULT ( STDMETHODCALLTYPE *MoveNext )( 
            ICoreWebView2HttpHeadersCollectionIterator * This,
            /* [retval][out] */ BOOL *hasNext);
        
        END_INTERFACE
    } ICoreWebView2HttpHeadersCollectionIteratorVtbl;

    interface ICoreWebView2HttpHeadersCollectionIterator
    {
        CONST_VTBL struct ICoreWebView2HttpHeadersCollectionIteratorVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2HttpHeadersCollectionIterator_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2HttpHeadersCollectionIterator_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2HttpHeadersCollectionIterator_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2HttpHeadersCollectionIterator_GetCurrentHeader(This,name,value)	\
    ( (This)->lpVtbl -> GetCurrentHeader(This,name,value) ) 

#define ICoreWebView2HttpHeadersCollectionIterator_get_HasCurrentHeader(This,hasCurrent)	\
    ( (This)->lpVtbl -> get_HasCurrentHeader(This,hasCurrent) ) 

#define ICoreWebView2HttpHeadersCollectionIterator_MoveNext(This,hasNext)	\
    ( (This)->lpVtbl -> MoveNext(This,hasNext) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2HttpHeadersCollectionIterator_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2HttpRequestHeaders_INTERFACE_DEFINED__
#define __ICoreWebView2HttpRequestHeaders_INTERFACE_DEFINED__

/* interface ICoreWebView2HttpRequestHeaders */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2HttpRequestHeaders = {0xe86cac0e,0x5523,0x465c,{0xb5,0x36,0x8f,0xb9,0xfc,0x8c,0x8c,0x60}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("e86cac0e-5523-465c-b536-8fb9fc8c8c60")
    ICoreWebView2HttpRequestHeaders : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE GetHeader( 
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ LPWSTR *value) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetHeaders( 
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Contains( 
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ BOOL *contains) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE SetHeader( 
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR value) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE RemoveHeader( 
            /* [in] */ LPCWSTR name) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetIterator( 
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2HttpRequestHeadersVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2HttpRequestHeaders * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2HttpRequestHeaders * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2HttpRequestHeaders * This);
        
        HRESULT ( STDMETHODCALLTYPE *GetHeader )( 
            ICoreWebView2HttpRequestHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ LPWSTR *value);
        
        HRESULT ( STDMETHODCALLTYPE *GetHeaders )( 
            ICoreWebView2HttpRequestHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator);
        
        HRESULT ( STDMETHODCALLTYPE *Contains )( 
            ICoreWebView2HttpRequestHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ BOOL *contains);
        
        HRESULT ( STDMETHODCALLTYPE *SetHeader )( 
            ICoreWebView2HttpRequestHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR value);
        
        HRESULT ( STDMETHODCALLTYPE *RemoveHeader )( 
            ICoreWebView2HttpRequestHeaders * This,
            /* [in] */ LPCWSTR name);
        
        HRESULT ( STDMETHODCALLTYPE *GetIterator )( 
            ICoreWebView2HttpRequestHeaders * This,
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator);
        
        END_INTERFACE
    } ICoreWebView2HttpRequestHeadersVtbl;

    interface ICoreWebView2HttpRequestHeaders
    {
        CONST_VTBL struct ICoreWebView2HttpRequestHeadersVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2HttpRequestHeaders_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2HttpRequestHeaders_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2HttpRequestHeaders_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2HttpRequestHeaders_GetHeader(This,name,value)	\
    ( (This)->lpVtbl -> GetHeader(This,name,value) ) 

#define ICoreWebView2HttpRequestHeaders_GetHeaders(This,name,iterator)	\
    ( (This)->lpVtbl -> GetHeaders(This,name,iterator) ) 

#define ICoreWebView2HttpRequestHeaders_Contains(This,name,contains)	\
    ( (This)->lpVtbl -> Contains(This,name,contains) ) 

#define ICoreWebView2HttpRequestHeaders_SetHeader(This,name,value)	\
    ( (This)->lpVtbl -> SetHeader(This,name,value) ) 

#define ICoreWebView2HttpRequestHeaders_RemoveHeader(This,name)	\
    ( (This)->lpVtbl -> RemoveHeader(This,name) ) 

#define ICoreWebView2HttpRequestHeaders_GetIterator(This,iterator)	\
    ( (This)->lpVtbl -> GetIterator(This,iterator) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2HttpRequestHeaders_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2HttpResponseHeaders_INTERFACE_DEFINED__
#define __ICoreWebView2HttpResponseHeaders_INTERFACE_DEFINED__

/* interface ICoreWebView2HttpResponseHeaders */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2HttpResponseHeaders = {0x03c5ff5a,0x9b45,0x4a88,{0x88,0x1c,0x89,0xa9,0xf3,0x28,0x61,0x9c}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("03c5ff5a-9b45-4a88-881c-89a9f328619c")
    ICoreWebView2HttpResponseHeaders : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE AppendHeader( 
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR value) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Contains( 
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ BOOL *contains) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetHeader( 
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ LPWSTR *value) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetHeaders( 
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetIterator( 
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2HttpResponseHeadersVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2HttpResponseHeaders * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2HttpResponseHeaders * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2HttpResponseHeaders * This);
        
        HRESULT ( STDMETHODCALLTYPE *AppendHeader )( 
            ICoreWebView2HttpResponseHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ LPCWSTR value);
        
        HRESULT ( STDMETHODCALLTYPE *Contains )( 
            ICoreWebView2HttpResponseHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ BOOL *contains);
        
        HRESULT ( STDMETHODCALLTYPE *GetHeader )( 
            ICoreWebView2HttpResponseHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ LPWSTR *value);
        
        HRESULT ( STDMETHODCALLTYPE *GetHeaders )( 
            ICoreWebView2HttpResponseHeaders * This,
            /* [in] */ LPCWSTR name,
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator);
        
        HRESULT ( STDMETHODCALLTYPE *GetIterator )( 
            ICoreWebView2HttpResponseHeaders * This,
            /* [retval][out] */ ICoreWebView2HttpHeadersCollectionIterator **iterator);
        
        END_INTERFACE
    } ICoreWebView2HttpResponseHeadersVtbl;

    interface ICoreWebView2HttpResponseHeaders
    {
        CONST_VTBL struct ICoreWebView2HttpResponseHeadersVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2HttpResponseHeaders_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2HttpResponseHeaders_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2HttpResponseHeaders_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2HttpResponseHeaders_AppendHeader(This,name,value)	\
    ( (This)->lpVtbl -> AppendHeader(This,name,value) ) 

#define ICoreWebView2HttpResponseHeaders_Contains(This,name,contains)	\
    ( (This)->lpVtbl -> Contains(This,name,contains) ) 

#define ICoreWebView2HttpResponseHeaders_GetHeader(This,name,value)	\
    ( (This)->lpVtbl -> GetHeader(This,name,value) ) 

#define ICoreWebView2HttpResponseHeaders_GetHeaders(This,name,iterator)	\
    ( (This)->lpVtbl -> GetHeaders(This,name,iterator) ) 

#define ICoreWebView2HttpResponseHeaders_GetIterator(This,iterator)	\
    ( (This)->lpVtbl -> GetIterator(This,iterator) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2HttpResponseHeaders_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Interop_INTERFACE_DEFINED__
#define __ICoreWebView2Interop_INTERFACE_DEFINED__

/* interface ICoreWebView2Interop */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Interop = {0x912b34a7,0xd10b,0x49c4,{0xaf,0x18,0x7c,0xb7,0xe6,0x04,0xe0,0x1a}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("912b34a7-d10b-49c4-af18-7cb7e604e01a")
    ICoreWebView2Interop : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE AddHostObjectToScript( 
            /* [in] */ LPCWSTR name,
            /* [in] */ VARIANT *object) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2InteropVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Interop * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Interop * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Interop * This);
        
        HRESULT ( STDMETHODCALLTYPE *AddHostObjectToScript )( 
            ICoreWebView2Interop * This,
            /* [in] */ LPCWSTR name,
            /* [in] */ VARIANT *object);
        
        END_INTERFACE
    } ICoreWebView2InteropVtbl;

    interface ICoreWebView2Interop
    {
        CONST_VTBL struct ICoreWebView2InteropVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Interop_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Interop_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Interop_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Interop_AddHostObjectToScript(This,name,object)	\
    ( (This)->lpVtbl -> AddHostObjectToScript(This,name,object) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Interop_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2MoveFocusRequestedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2MoveFocusRequestedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2MoveFocusRequestedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2MoveFocusRequestedEventArgs = {0x2d6aa13b,0x3839,0x4a15,{0x92,0xfc,0xd8,0x8b,0x3c,0x0d,0x9c,0x9d}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("2d6aa13b-3839-4a15-92fc-d88b3c0d9c9d")
    ICoreWebView2MoveFocusRequestedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Reason( 
            /* [retval][out] */ COREWEBVIEW2_MOVE_FOCUS_REASON *reason) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Handled( 
            /* [retval][out] */ BOOL *value) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Handled( 
            /* [in] */ BOOL value) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2MoveFocusRequestedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2MoveFocusRequestedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2MoveFocusRequestedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2MoveFocusRequestedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Reason )( 
            ICoreWebView2MoveFocusRequestedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_MOVE_FOCUS_REASON *reason);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Handled )( 
            ICoreWebView2MoveFocusRequestedEventArgs * This,
            /* [retval][out] */ BOOL *value);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Handled )( 
            ICoreWebView2MoveFocusRequestedEventArgs * This,
            /* [in] */ BOOL value);
        
        END_INTERFACE
    } ICoreWebView2MoveFocusRequestedEventArgsVtbl;

    interface ICoreWebView2MoveFocusRequestedEventArgs
    {
        CONST_VTBL struct ICoreWebView2MoveFocusRequestedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2MoveFocusRequestedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2MoveFocusRequestedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2MoveFocusRequestedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2MoveFocusRequestedEventArgs_get_Reason(This,reason)	\
    ( (This)->lpVtbl -> get_Reason(This,reason) ) 

#define ICoreWebView2MoveFocusRequestedEventArgs_get_Handled(This,value)	\
    ( (This)->lpVtbl -> get_Handled(This,value) ) 

#define ICoreWebView2MoveFocusRequestedEventArgs_put_Handled(This,value)	\
    ( (This)->lpVtbl -> put_Handled(This,value) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2MoveFocusRequestedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2MoveFocusRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2MoveFocusRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2MoveFocusRequestedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2MoveFocusRequestedEventHandler = {0x69035451,0x6dc7,0x4cb8,{0x9b,0xce,0xb2,0xbd,0x70,0xad,0x28,0x9f}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("69035451-6dc7-4cb8-9bce-b2bd70ad289f")
    ICoreWebView2MoveFocusRequestedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ ICoreWebView2MoveFocusRequestedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2MoveFocusRequestedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2MoveFocusRequestedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2MoveFocusRequestedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2MoveFocusRequestedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2MoveFocusRequestedEventHandler * This,
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ ICoreWebView2MoveFocusRequestedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2MoveFocusRequestedEventHandlerVtbl;

    interface ICoreWebView2MoveFocusRequestedEventHandler
    {
        CONST_VTBL struct ICoreWebView2MoveFocusRequestedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2MoveFocusRequestedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2MoveFocusRequestedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2MoveFocusRequestedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2MoveFocusRequestedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2MoveFocusRequestedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NavigationCompletedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2NavigationCompletedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2NavigationCompletedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2NavigationCompletedEventArgs = {0x30d68b7d,0x20d9,0x4752,{0xa9,0xca,0xec,0x84,0x48,0xfb,0xb5,0xc1}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("30d68b7d-20d9-4752-a9ca-ec8448fbb5c1")
    ICoreWebView2NavigationCompletedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsSuccess( 
            /* [retval][out] */ BOOL *isSuccess) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_WebErrorStatus( 
            /* [retval][out] */ COREWEBVIEW2_WEB_ERROR_STATUS *webErrorStatus) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_NavigationId( 
            /* [retval][out] */ UINT64 *navigationId) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2NavigationCompletedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2NavigationCompletedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2NavigationCompletedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2NavigationCompletedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsSuccess )( 
            ICoreWebView2NavigationCompletedEventArgs * This,
            /* [retval][out] */ BOOL *isSuccess);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_WebErrorStatus )( 
            ICoreWebView2NavigationCompletedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_WEB_ERROR_STATUS *webErrorStatus);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_NavigationId )( 
            ICoreWebView2NavigationCompletedEventArgs * This,
            /* [retval][out] */ UINT64 *navigationId);
        
        END_INTERFACE
    } ICoreWebView2NavigationCompletedEventArgsVtbl;

    interface ICoreWebView2NavigationCompletedEventArgs
    {
        CONST_VTBL struct ICoreWebView2NavigationCompletedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2NavigationCompletedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2NavigationCompletedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2NavigationCompletedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2NavigationCompletedEventArgs_get_IsSuccess(This,isSuccess)	\
    ( (This)->lpVtbl -> get_IsSuccess(This,isSuccess) ) 

#define ICoreWebView2NavigationCompletedEventArgs_get_WebErrorStatus(This,webErrorStatus)	\
    ( (This)->lpVtbl -> get_WebErrorStatus(This,webErrorStatus) ) 

#define ICoreWebView2NavigationCompletedEventArgs_get_NavigationId(This,navigationId)	\
    ( (This)->lpVtbl -> get_NavigationId(This,navigationId) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NavigationCompletedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NavigationCompletedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2NavigationCompletedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2NavigationCompletedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2NavigationCompletedEventHandler = {0xd33a35bf,0x1c49,0x4f98,{0x93,0xab,0x00,0x6e,0x05,0x33,0xfe,0x1c}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("d33a35bf-1c49-4f98-93ab-006e0533fe1c")
    ICoreWebView2NavigationCompletedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2NavigationCompletedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2NavigationCompletedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2NavigationCompletedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2NavigationCompletedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2NavigationCompletedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2NavigationCompletedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2NavigationCompletedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2NavigationCompletedEventHandlerVtbl;

    interface ICoreWebView2NavigationCompletedEventHandler
    {
        CONST_VTBL struct ICoreWebView2NavigationCompletedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2NavigationCompletedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2NavigationCompletedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2NavigationCompletedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2NavigationCompletedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NavigationCompletedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NavigationStartingEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2NavigationStartingEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2NavigationStartingEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2NavigationStartingEventArgs = {0x5b495469,0xe119,0x438a,{0x9b,0x18,0x76,0x04,0xf2,0x5f,0x2e,0x49}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("5b495469-e119-438a-9b18-7604f25f2e49")
    ICoreWebView2NavigationStartingEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Uri( 
            /* [retval][out] */ LPWSTR *uri) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsUserInitiated( 
            /* [retval][out] */ BOOL *isUserInitiated) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsRedirected( 
            /* [retval][out] */ BOOL *isRedirected) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_RequestHeaders( 
            /* [retval][out] */ ICoreWebView2HttpRequestHeaders **requestHeaders) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Cancel( 
            /* [retval][out] */ BOOL *cancel) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Cancel( 
            /* [in] */ BOOL cancel) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_NavigationId( 
            /* [retval][out] */ UINT64 *navigationId) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2NavigationStartingEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2NavigationStartingEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2NavigationStartingEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Uri )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [retval][out] */ LPWSTR *uri);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsUserInitiated )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [retval][out] */ BOOL *isUserInitiated);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsRedirected )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [retval][out] */ BOOL *isRedirected);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_RequestHeaders )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [retval][out] */ ICoreWebView2HttpRequestHeaders **requestHeaders);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Cancel )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [retval][out] */ BOOL *cancel);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Cancel )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [in] */ BOOL cancel);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_NavigationId )( 
            ICoreWebView2NavigationStartingEventArgs * This,
            /* [retval][out] */ UINT64 *navigationId);
        
        END_INTERFACE
    } ICoreWebView2NavigationStartingEventArgsVtbl;

    interface ICoreWebView2NavigationStartingEventArgs
    {
        CONST_VTBL struct ICoreWebView2NavigationStartingEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2NavigationStartingEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2NavigationStartingEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2NavigationStartingEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2NavigationStartingEventArgs_get_Uri(This,uri)	\
    ( (This)->lpVtbl -> get_Uri(This,uri) ) 

#define ICoreWebView2NavigationStartingEventArgs_get_IsUserInitiated(This,isUserInitiated)	\
    ( (This)->lpVtbl -> get_IsUserInitiated(This,isUserInitiated) ) 

#define ICoreWebView2NavigationStartingEventArgs_get_IsRedirected(This,isRedirected)	\
    ( (This)->lpVtbl -> get_IsRedirected(This,isRedirected) ) 

#define ICoreWebView2NavigationStartingEventArgs_get_RequestHeaders(This,requestHeaders)	\
    ( (This)->lpVtbl -> get_RequestHeaders(This,requestHeaders) ) 

#define ICoreWebView2NavigationStartingEventArgs_get_Cancel(This,cancel)	\
    ( (This)->lpVtbl -> get_Cancel(This,cancel) ) 

#define ICoreWebView2NavigationStartingEventArgs_put_Cancel(This,cancel)	\
    ( (This)->lpVtbl -> put_Cancel(This,cancel) ) 

#define ICoreWebView2NavigationStartingEventArgs_get_NavigationId(This,navigationId)	\
    ( (This)->lpVtbl -> get_NavigationId(This,navigationId) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NavigationStartingEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NavigationStartingEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2NavigationStartingEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2NavigationStartingEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2NavigationStartingEventHandler = {0x9adbe429,0xf36d,0x432b,{0x9d,0xdc,0xf8,0x88,0x1f,0xbd,0x76,0xe3}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("9adbe429-f36d-432b-9ddc-f8881fbd76e3")
    ICoreWebView2NavigationStartingEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2NavigationStartingEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2NavigationStartingEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2NavigationStartingEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2NavigationStartingEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2NavigationStartingEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2NavigationStartingEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2NavigationStartingEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2NavigationStartingEventHandlerVtbl;

    interface ICoreWebView2NavigationStartingEventHandler
    {
        CONST_VTBL struct ICoreWebView2NavigationStartingEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2NavigationStartingEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2NavigationStartingEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2NavigationStartingEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2NavigationStartingEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NavigationStartingEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NewBrowserVersionAvailableEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2NewBrowserVersionAvailableEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2NewBrowserVersionAvailableEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2NewBrowserVersionAvailableEventHandler = {0xf9a2976e,0xd34e,0x44fc,{0xad,0xee,0x81,0xb6,0xb5,0x7c,0xa9,0x14}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("f9a2976e-d34e-44fc-adee-81b6b57ca914")
    ICoreWebView2NewBrowserVersionAvailableEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2Environment *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2NewBrowserVersionAvailableEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2NewBrowserVersionAvailableEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2NewBrowserVersionAvailableEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2NewBrowserVersionAvailableEventHandler * This,
            /* [in] */ ICoreWebView2Environment *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl;

    interface ICoreWebView2NewBrowserVersionAvailableEventHandler
    {
        CONST_VTBL struct ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2NewBrowserVersionAvailableEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2NewBrowserVersionAvailableEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2NewBrowserVersionAvailableEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2NewBrowserVersionAvailableEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NewBrowserVersionAvailableEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NewWindowRequestedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2NewWindowRequestedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2NewWindowRequestedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2NewWindowRequestedEventArgs = {0x34acb11c,0xfc37,0x4418,{0x91,0x32,0xf9,0xc2,0x1d,0x1e,0xaf,0xb9}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("34acb11c-fc37-4418-9132-f9c21d1eafb9")
    ICoreWebView2NewWindowRequestedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Uri( 
            /* [retval][out] */ LPWSTR *uri) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_NewWindow( 
            /* [in] */ ICoreWebView2 *newWindow) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_NewWindow( 
            /* [retval][out] */ ICoreWebView2 **newWindow) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Handled( 
            /* [in] */ BOOL handled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Handled( 
            /* [retval][out] */ BOOL *handled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsUserInitiated( 
            /* [retval][out] */ BOOL *isUserInitiated) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetDeferral( 
            /* [retval][out] */ ICoreWebView2Deferral **deferral) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_WindowFeatures( 
            /* [retval][out] */ ICoreWebView2WindowFeatures **value) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2NewWindowRequestedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2NewWindowRequestedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2NewWindowRequestedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Uri )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [retval][out] */ LPWSTR *uri);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_NewWindow )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [in] */ ICoreWebView2 *newWindow);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_NewWindow )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [retval][out] */ ICoreWebView2 **newWindow);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Handled )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [in] */ BOOL handled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Handled )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [retval][out] */ BOOL *handled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsUserInitiated )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [retval][out] */ BOOL *isUserInitiated);
        
        HRESULT ( STDMETHODCALLTYPE *GetDeferral )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [retval][out] */ ICoreWebView2Deferral **deferral);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_WindowFeatures )( 
            ICoreWebView2NewWindowRequestedEventArgs * This,
            /* [retval][out] */ ICoreWebView2WindowFeatures **value);
        
        END_INTERFACE
    } ICoreWebView2NewWindowRequestedEventArgsVtbl;

    interface ICoreWebView2NewWindowRequestedEventArgs
    {
        CONST_VTBL struct ICoreWebView2NewWindowRequestedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2NewWindowRequestedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2NewWindowRequestedEventArgs_get_Uri(This,uri)	\
    ( (This)->lpVtbl -> get_Uri(This,uri) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_put_NewWindow(This,newWindow)	\
    ( (This)->lpVtbl -> put_NewWindow(This,newWindow) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_get_NewWindow(This,newWindow)	\
    ( (This)->lpVtbl -> get_NewWindow(This,newWindow) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_put_Handled(This,handled)	\
    ( (This)->lpVtbl -> put_Handled(This,handled) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_get_Handled(This,handled)	\
    ( (This)->lpVtbl -> get_Handled(This,handled) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_get_IsUserInitiated(This,isUserInitiated)	\
    ( (This)->lpVtbl -> get_IsUserInitiated(This,isUserInitiated) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_GetDeferral(This,deferral)	\
    ( (This)->lpVtbl -> GetDeferral(This,deferral) ) 

#define ICoreWebView2NewWindowRequestedEventArgs_get_WindowFeatures(This,value)	\
    ( (This)->lpVtbl -> get_WindowFeatures(This,value) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NewWindowRequestedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NewWindowRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2NewWindowRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2NewWindowRequestedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2NewWindowRequestedEventHandler = {0xd4c185fe,0xc81c,0x4989,{0x97,0xaf,0x2d,0x3f,0xa7,0xab,0x56,0x51}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("d4c185fe-c81c-4989-97af-2d3fa7ab5651")
    ICoreWebView2NewWindowRequestedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2NewWindowRequestedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2NewWindowRequestedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2NewWindowRequestedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2NewWindowRequestedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2NewWindowRequestedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2NewWindowRequestedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2NewWindowRequestedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2NewWindowRequestedEventHandlerVtbl;

    interface ICoreWebView2NewWindowRequestedEventHandler
    {
        CONST_VTBL struct ICoreWebView2NewWindowRequestedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2NewWindowRequestedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2NewWindowRequestedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2NewWindowRequestedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2NewWindowRequestedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NewWindowRequestedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2PermissionRequestedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2PermissionRequestedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2PermissionRequestedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2PermissionRequestedEventArgs = {0x973ae2ef,0xff18,0x4894,{0x8f,0xb2,0x3c,0x75,0x8f,0x04,0x68,0x10}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("973ae2ef-ff18-4894-8fb2-3c758f046810")
    ICoreWebView2PermissionRequestedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Uri( 
            /* [retval][out] */ LPWSTR *uri) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PermissionKind( 
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_KIND *permissionKind) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsUserInitiated( 
            /* [retval][out] */ BOOL *isUserInitiated) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_State( 
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_STATE *state) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_State( 
            /* [in] */ COREWEBVIEW2_PERMISSION_STATE state) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetDeferral( 
            /* [retval][out] */ ICoreWebView2Deferral **deferral) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2PermissionRequestedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2PermissionRequestedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2PermissionRequestedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Uri )( 
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [retval][out] */ LPWSTR *uri);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PermissionKind )( 
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_KIND *permissionKind);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsUserInitiated )( 
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [retval][out] */ BOOL *isUserInitiated);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_State )( 
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_STATE *state);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_State )( 
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [in] */ COREWEBVIEW2_PERMISSION_STATE state);
        
        HRESULT ( STDMETHODCALLTYPE *GetDeferral )( 
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [retval][out] */ ICoreWebView2Deferral **deferral);
        
        END_INTERFACE
    } ICoreWebView2PermissionRequestedEventArgsVtbl;

    interface ICoreWebView2PermissionRequestedEventArgs
    {
        CONST_VTBL struct ICoreWebView2PermissionRequestedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2PermissionRequestedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2PermissionRequestedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2PermissionRequestedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2PermissionRequestedEventArgs_get_Uri(This,uri)	\
    ( (This)->lpVtbl -> get_Uri(This,uri) ) 

#define ICoreWebView2PermissionRequestedEventArgs_get_PermissionKind(This,permissionKind)	\
    ( (This)->lpVtbl -> get_PermissionKind(This,permissionKind) ) 

#define ICoreWebView2PermissionRequestedEventArgs_get_IsUserInitiated(This,isUserInitiated)	\
    ( (This)->lpVtbl -> get_IsUserInitiated(This,isUserInitiated) ) 

#define ICoreWebView2PermissionRequestedEventArgs_get_State(This,state)	\
    ( (This)->lpVtbl -> get_State(This,state) ) 

#define ICoreWebView2PermissionRequestedEventArgs_put_State(This,state)	\
    ( (This)->lpVtbl -> put_State(This,state) ) 

#define ICoreWebView2PermissionRequestedEventArgs_GetDeferral(This,deferral)	\
    ( (This)->lpVtbl -> GetDeferral(This,deferral) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2PermissionRequestedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2PermissionRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2PermissionRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2PermissionRequestedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2PermissionRequestedEventHandler = {0x15e1c6a3,0xc72a,0x4df3,{0x91,0xd7,0xd0,0x97,0xfb,0xec,0x6b,0xfd}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("15e1c6a3-c72a-4df3-91d7-d097fbec6bfd")
    ICoreWebView2PermissionRequestedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2PermissionRequestedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2PermissionRequestedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2PermissionRequestedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2PermissionRequestedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2PermissionRequestedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2PermissionRequestedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2PermissionRequestedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2PermissionRequestedEventHandlerVtbl;

    interface ICoreWebView2PermissionRequestedEventHandler
    {
        CONST_VTBL struct ICoreWebView2PermissionRequestedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2PermissionRequestedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2PermissionRequestedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2PermissionRequestedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2PermissionRequestedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2PermissionRequestedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2PointerInfo_INTERFACE_DEFINED__
#define __ICoreWebView2PointerInfo_INTERFACE_DEFINED__

/* interface ICoreWebView2PointerInfo */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2PointerInfo = {0xe6995887,0xd10d,0x4f5d,{0x93,0x59,0x4c,0xe4,0x6e,0x4f,0x96,0xb9}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("e6995887-d10d-4f5d-9359-4ce46e4f96b9")
    ICoreWebView2PointerInfo : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PointerKind( 
            /* [retval][out] */ DWORD *pointerKind) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PointerKind( 
            /* [in] */ DWORD pointerKind) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PointerId( 
            /* [retval][out] */ UINT32 *pointerId) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PointerId( 
            /* [in] */ UINT32 pointerId) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_FrameId( 
            /* [retval][out] */ UINT32 *frameId) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_FrameId( 
            /* [in] */ UINT32 frameId) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PointerFlags( 
            /* [retval][out] */ UINT32 *pointerFlags) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PointerFlags( 
            /* [in] */ UINT32 pointerFlags) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PointerDeviceRect( 
            /* [retval][out] */ RECT *pointerDeviceRect) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PointerDeviceRect( 
            /* [in] */ RECT pointerDeviceRect) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_DisplayRect( 
            /* [retval][out] */ RECT *displayRect) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_DisplayRect( 
            /* [in] */ RECT displayRect) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PixelLocation( 
            /* [retval][out] */ POINT *pixelLocation) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PixelLocation( 
            /* [in] */ POINT pixelLocation) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HimetricLocation( 
            /* [retval][out] */ POINT *himetricLocation) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_HimetricLocation( 
            /* [in] */ POINT himetricLocation) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PixelLocationRaw( 
            /* [retval][out] */ POINT *pixelLocationRaw) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PixelLocationRaw( 
            /* [in] */ POINT pixelLocationRaw) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HimetricLocationRaw( 
            /* [retval][out] */ POINT *himetricLocationRaw) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_HimetricLocationRaw( 
            /* [in] */ POINT himetricLocationRaw) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Time( 
            /* [retval][out] */ DWORD *time) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Time( 
            /* [in] */ DWORD time) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HistoryCount( 
            /* [retval][out] */ UINT32 *historyCount) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_HistoryCount( 
            /* [in] */ UINT32 historyCount) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_InputData( 
            /* [retval][out] */ INT32 *inputData) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_InputData( 
            /* [in] */ INT32 inputData) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_KeyStates( 
            /* [retval][out] */ DWORD *keyStates) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_KeyStates( 
            /* [in] */ DWORD keyStates) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PerformanceCount( 
            /* [retval][out] */ UINT64 *performanceCount) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PerformanceCount( 
            /* [in] */ UINT64 performanceCount) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ButtonChangeKind( 
            /* [retval][out] */ INT32 *buttonChangeKind) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_ButtonChangeKind( 
            /* [in] */ INT32 buttonChangeKind) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PenFlags( 
            /* [retval][out] */ UINT32 *penFLags) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PenFlags( 
            /* [in] */ UINT32 penFLags) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PenMask( 
            /* [retval][out] */ UINT32 *penMask) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PenMask( 
            /* [in] */ UINT32 penMask) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PenPressure( 
            /* [retval][out] */ UINT32 *penPressure) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PenPressure( 
            /* [in] */ UINT32 penPressure) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PenRotation( 
            /* [retval][out] */ UINT32 *penRotation) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PenRotation( 
            /* [in] */ UINT32 penRotation) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PenTiltX( 
            /* [retval][out] */ INT32 *penTiltX) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PenTiltX( 
            /* [in] */ INT32 penTiltX) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PenTiltY( 
            /* [retval][out] */ INT32 *penTiltY) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_PenTiltY( 
            /* [in] */ INT32 penTiltY) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_TouchFlags( 
            /* [retval][out] */ UINT32 *touchFlags) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_TouchFlags( 
            /* [in] */ UINT32 touchFlags) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_TouchMask( 
            /* [retval][out] */ UINT32 *touchMask) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_TouchMask( 
            /* [in] */ UINT32 touchMask) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_TouchContact( 
            /* [retval][out] */ RECT *touchContact) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_TouchContact( 
            /* [in] */ RECT touchContact) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_TouchContactRaw( 
            /* [retval][out] */ RECT *touchContactRaw) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_TouchContactRaw( 
            /* [in] */ RECT touchContactRaw) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_TouchOrientation( 
            /* [retval][out] */ UINT32 *touchOrientation) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_TouchOrientation( 
            /* [in] */ UINT32 touchOrientation) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_TouchPressure( 
            /* [retval][out] */ UINT32 *touchPressure) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_TouchPressure( 
            /* [in] */ UINT32 touchPressure) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2PointerInfoVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2PointerInfo * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2PointerInfo * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PointerKind )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ DWORD *pointerKind);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PointerKind )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ DWORD pointerKind);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PointerId )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *pointerId);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PointerId )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 pointerId);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_FrameId )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *frameId);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_FrameId )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 frameId);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PointerFlags )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *pointerFlags);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PointerFlags )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 pointerFlags);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PointerDeviceRect )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ RECT *pointerDeviceRect);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PointerDeviceRect )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ RECT pointerDeviceRect);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_DisplayRect )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ RECT *displayRect);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_DisplayRect )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ RECT displayRect);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PixelLocation )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ POINT *pixelLocation);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PixelLocation )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ POINT pixelLocation);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_HimetricLocation )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ POINT *himetricLocation);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_HimetricLocation )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ POINT himetricLocation);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PixelLocationRaw )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ POINT *pixelLocationRaw);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PixelLocationRaw )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ POINT pixelLocationRaw);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_HimetricLocationRaw )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ POINT *himetricLocationRaw);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_HimetricLocationRaw )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ POINT himetricLocationRaw);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Time )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ DWORD *time);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Time )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ DWORD time);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_HistoryCount )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *historyCount);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_HistoryCount )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 historyCount);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_InputData )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ INT32 *inputData);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_InputData )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ INT32 inputData);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_KeyStates )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ DWORD *keyStates);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_KeyStates )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ DWORD keyStates);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PerformanceCount )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT64 *performanceCount);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PerformanceCount )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT64 performanceCount);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ButtonChangeKind )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ INT32 *buttonChangeKind);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ButtonChangeKind )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ INT32 buttonChangeKind);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PenFlags )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *penFLags);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PenFlags )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 penFLags);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PenMask )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *penMask);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PenMask )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 penMask);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PenPressure )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *penPressure);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PenPressure )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 penPressure);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PenRotation )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *penRotation);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PenRotation )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 penRotation);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PenTiltX )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ INT32 *penTiltX);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PenTiltX )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ INT32 penTiltX);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_PenTiltY )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ INT32 *penTiltY);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_PenTiltY )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ INT32 penTiltY);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_TouchFlags )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *touchFlags);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_TouchFlags )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 touchFlags);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_TouchMask )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *touchMask);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_TouchMask )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 touchMask);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_TouchContact )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ RECT *touchContact);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_TouchContact )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ RECT touchContact);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_TouchContactRaw )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ RECT *touchContactRaw);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_TouchContactRaw )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ RECT touchContactRaw);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_TouchOrientation )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *touchOrientation);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_TouchOrientation )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 touchOrientation);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_TouchPressure )( 
            ICoreWebView2PointerInfo * This,
            /* [retval][out] */ UINT32 *touchPressure);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_TouchPressure )( 
            ICoreWebView2PointerInfo * This,
            /* [in] */ UINT32 touchPressure);
        
        END_INTERFACE
    } ICoreWebView2PointerInfoVtbl;

    interface ICoreWebView2PointerInfo
    {
        CONST_VTBL struct ICoreWebView2PointerInfoVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2PointerInfo_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2PointerInfo_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2PointerInfo_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2PointerInfo_get_PointerKind(This,pointerKind)	\
    ( (This)->lpVtbl -> get_PointerKind(This,pointerKind) ) 

#define ICoreWebView2PointerInfo_put_PointerKind(This,pointerKind)	\
    ( (This)->lpVtbl -> put_PointerKind(This,pointerKind) ) 

#define ICoreWebView2PointerInfo_get_PointerId(This,pointerId)	\
    ( (This)->lpVtbl -> get_PointerId(This,pointerId) ) 

#define ICoreWebView2PointerInfo_put_PointerId(This,pointerId)	\
    ( (This)->lpVtbl -> put_PointerId(This,pointerId) ) 

#define ICoreWebView2PointerInfo_get_FrameId(This,frameId)	\
    ( (This)->lpVtbl -> get_FrameId(This,frameId) ) 

#define ICoreWebView2PointerInfo_put_FrameId(This,frameId)	\
    ( (This)->lpVtbl -> put_FrameId(This,frameId) ) 

#define ICoreWebView2PointerInfo_get_PointerFlags(This,pointerFlags)	\
    ( (This)->lpVtbl -> get_PointerFlags(This,pointerFlags) ) 

#define ICoreWebView2PointerInfo_put_PointerFlags(This,pointerFlags)	\
    ( (This)->lpVtbl -> put_PointerFlags(This,pointerFlags) ) 

#define ICoreWebView2PointerInfo_get_PointerDeviceRect(This,pointerDeviceRect)	\
    ( (This)->lpVtbl -> get_PointerDeviceRect(This,pointerDeviceRect) ) 

#define ICoreWebView2PointerInfo_put_PointerDeviceRect(This,pointerDeviceRect)	\
    ( (This)->lpVtbl -> put_PointerDeviceRect(This,pointerDeviceRect) ) 

#define ICoreWebView2PointerInfo_get_DisplayRect(This,displayRect)	\
    ( (This)->lpVtbl -> get_DisplayRect(This,displayRect) ) 

#define ICoreWebView2PointerInfo_put_DisplayRect(This,displayRect)	\
    ( (This)->lpVtbl -> put_DisplayRect(This,displayRect) ) 

#define ICoreWebView2PointerInfo_get_PixelLocation(This,pixelLocation)	\
    ( (This)->lpVtbl -> get_PixelLocation(This,pixelLocation) ) 

#define ICoreWebView2PointerInfo_put_PixelLocation(This,pixelLocation)	\
    ( (This)->lpVtbl -> put_PixelLocation(This,pixelLocation) ) 

#define ICoreWebView2PointerInfo_get_HimetricLocation(This,himetricLocation)	\
    ( (This)->lpVtbl -> get_HimetricLocation(This,himetricLocation) ) 

#define ICoreWebView2PointerInfo_put_HimetricLocation(This,himetricLocation)	\
    ( (This)->lpVtbl -> put_HimetricLocation(This,himetricLocation) ) 

#define ICoreWebView2PointerInfo_get_PixelLocationRaw(This,pixelLocationRaw)	\
    ( (This)->lpVtbl -> get_PixelLocationRaw(This,pixelLocationRaw) ) 

#define ICoreWebView2PointerInfo_put_PixelLocationRaw(This,pixelLocationRaw)	\
    ( (This)->lpVtbl -> put_PixelLocationRaw(This,pixelLocationRaw) ) 

#define ICoreWebView2PointerInfo_get_HimetricLocationRaw(This,himetricLocationRaw)	\
    ( (This)->lpVtbl -> get_HimetricLocationRaw(This,himetricLocationRaw) ) 

#define ICoreWebView2PointerInfo_put_HimetricLocationRaw(This,himetricLocationRaw)	\
    ( (This)->lpVtbl -> put_HimetricLocationRaw(This,himetricLocationRaw) ) 

#define ICoreWebView2PointerInfo_get_Time(This,time)	\
    ( (This)->lpVtbl -> get_Time(This,time) ) 

#define ICoreWebView2PointerInfo_put_Time(This,time)	\
    ( (This)->lpVtbl -> put_Time(This,time) ) 

#define ICoreWebView2PointerInfo_get_HistoryCount(This,historyCount)	\
    ( (This)->lpVtbl -> get_HistoryCount(This,historyCount) ) 

#define ICoreWebView2PointerInfo_put_HistoryCount(This,historyCount)	\
    ( (This)->lpVtbl -> put_HistoryCount(This,historyCount) ) 

#define ICoreWebView2PointerInfo_get_InputData(This,inputData)	\
    ( (This)->lpVtbl -> get_InputData(This,inputData) ) 

#define ICoreWebView2PointerInfo_put_InputData(This,inputData)	\
    ( (This)->lpVtbl -> put_InputData(This,inputData) ) 

#define ICoreWebView2PointerInfo_get_KeyStates(This,keyStates)	\
    ( (This)->lpVtbl -> get_KeyStates(This,keyStates) ) 

#define ICoreWebView2PointerInfo_put_KeyStates(This,keyStates)	\
    ( (This)->lpVtbl -> put_KeyStates(This,keyStates) ) 

#define ICoreWebView2PointerInfo_get_PerformanceCount(This,performanceCount)	\
    ( (This)->lpVtbl -> get_PerformanceCount(This,performanceCount) ) 

#define ICoreWebView2PointerInfo_put_PerformanceCount(This,performanceCount)	\
    ( (This)->lpVtbl -> put_PerformanceCount(This,performanceCount) ) 

#define ICoreWebView2PointerInfo_get_ButtonChangeKind(This,buttonChangeKind)	\
    ( (This)->lpVtbl -> get_ButtonChangeKind(This,buttonChangeKind) ) 

#define ICoreWebView2PointerInfo_put_ButtonChangeKind(This,buttonChangeKind)	\
    ( (This)->lpVtbl -> put_ButtonChangeKind(This,buttonChangeKind) ) 

#define ICoreWebView2PointerInfo_get_PenFlags(This,penFLags)	\
    ( (This)->lpVtbl -> get_PenFlags(This,penFLags) ) 

#define ICoreWebView2PointerInfo_put_PenFlags(This,penFLags)	\
    ( (This)->lpVtbl -> put_PenFlags(This,penFLags) ) 

#define ICoreWebView2PointerInfo_get_PenMask(This,penMask)	\
    ( (This)->lpVtbl -> get_PenMask(This,penMask) ) 

#define ICoreWebView2PointerInfo_put_PenMask(This,penMask)	\
    ( (This)->lpVtbl -> put_PenMask(This,penMask) ) 

#define ICoreWebView2PointerInfo_get_PenPressure(This,penPressure)	\
    ( (This)->lpVtbl -> get_PenPressure(This,penPressure) ) 

#define ICoreWebView2PointerInfo_put_PenPressure(This,penPressure)	\
    ( (This)->lpVtbl -> put_PenPressure(This,penPressure) ) 

#define ICoreWebView2PointerInfo_get_PenRotation(This,penRotation)	\
    ( (This)->lpVtbl -> get_PenRotation(This,penRotation) ) 

#define ICoreWebView2PointerInfo_put_PenRotation(This,penRotation)	\
    ( (This)->lpVtbl -> put_PenRotation(This,penRotation) ) 

#define ICoreWebView2PointerInfo_get_PenTiltX(This,penTiltX)	\
    ( (This)->lpVtbl -> get_PenTiltX(This,penTiltX) ) 

#define ICoreWebView2PointerInfo_put_PenTiltX(This,penTiltX)	\
    ( (This)->lpVtbl -> put_PenTiltX(This,penTiltX) ) 

#define ICoreWebView2PointerInfo_get_PenTiltY(This,penTiltY)	\
    ( (This)->lpVtbl -> get_PenTiltY(This,penTiltY) ) 

#define ICoreWebView2PointerInfo_put_PenTiltY(This,penTiltY)	\
    ( (This)->lpVtbl -> put_PenTiltY(This,penTiltY) ) 

#define ICoreWebView2PointerInfo_get_TouchFlags(This,touchFlags)	\
    ( (This)->lpVtbl -> get_TouchFlags(This,touchFlags) ) 

#define ICoreWebView2PointerInfo_put_TouchFlags(This,touchFlags)	\
    ( (This)->lpVtbl -> put_TouchFlags(This,touchFlags) ) 

#define ICoreWebView2PointerInfo_get_TouchMask(This,touchMask)	\
    ( (This)->lpVtbl -> get_TouchMask(This,touchMask) ) 

#define ICoreWebView2PointerInfo_put_TouchMask(This,touchMask)	\
    ( (This)->lpVtbl -> put_TouchMask(This,touchMask) ) 

#define ICoreWebView2PointerInfo_get_TouchContact(This,touchContact)	\
    ( (This)->lpVtbl -> get_TouchContact(This,touchContact) ) 

#define ICoreWebView2PointerInfo_put_TouchContact(This,touchContact)	\
    ( (This)->lpVtbl -> put_TouchContact(This,touchContact) ) 

#define ICoreWebView2PointerInfo_get_TouchContactRaw(This,touchContactRaw)	\
    ( (This)->lpVtbl -> get_TouchContactRaw(This,touchContactRaw) ) 

#define ICoreWebView2PointerInfo_put_TouchContactRaw(This,touchContactRaw)	\
    ( (This)->lpVtbl -> put_TouchContactRaw(This,touchContactRaw) ) 

#define ICoreWebView2PointerInfo_get_TouchOrientation(This,touchOrientation)	\
    ( (This)->lpVtbl -> get_TouchOrientation(This,touchOrientation) ) 

#define ICoreWebView2PointerInfo_put_TouchOrientation(This,touchOrientation)	\
    ( (This)->lpVtbl -> put_TouchOrientation(This,touchOrientation) ) 

#define ICoreWebView2PointerInfo_get_TouchPressure(This,touchPressure)	\
    ( (This)->lpVtbl -> get_TouchPressure(This,touchPressure) ) 

#define ICoreWebView2PointerInfo_put_TouchPressure(This,touchPressure)	\
    ( (This)->lpVtbl -> put_TouchPressure(This,touchPressure) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2PointerInfo_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ProcessFailedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2ProcessFailedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2ProcessFailedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ProcessFailedEventArgs = {0x8155a9a4,0x1474,0x4a86,{0x8c,0xae,0x15,0x1b,0x0f,0xa6,0xb8,0xca}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("8155a9a4-1474-4a86-8cae-151b0fa6b8ca")
    ICoreWebView2ProcessFailedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ProcessFailedKind( 
            /* [retval][out] */ COREWEBVIEW2_PROCESS_FAILED_KIND *processFailedKind) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ProcessFailedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ProcessFailedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ProcessFailedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ProcessFailedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ProcessFailedKind )( 
            ICoreWebView2ProcessFailedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_PROCESS_FAILED_KIND *processFailedKind);
        
        END_INTERFACE
    } ICoreWebView2ProcessFailedEventArgsVtbl;

    interface ICoreWebView2ProcessFailedEventArgs
    {
        CONST_VTBL struct ICoreWebView2ProcessFailedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ProcessFailedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ProcessFailedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ProcessFailedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ProcessFailedEventArgs_get_ProcessFailedKind(This,processFailedKind)	\
    ( (This)->lpVtbl -> get_ProcessFailedKind(This,processFailedKind) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ProcessFailedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ProcessFailedEventArgs2_INTERFACE_DEFINED__
#define __ICoreWebView2ProcessFailedEventArgs2_INTERFACE_DEFINED__

/* interface ICoreWebView2ProcessFailedEventArgs2 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ProcessFailedEventArgs2 = {0x4dab9422,0x46fa,0x4c3e,{0xa5,0xd2,0x41,0xd2,0x07,0x1d,0x36,0x80}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("4dab9422-46fa-4c3e-a5d2-41d2071d3680")
    ICoreWebView2ProcessFailedEventArgs2 : public ICoreWebView2ProcessFailedEventArgs
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Reason( 
            /* [retval][out] */ COREWEBVIEW2_PROCESS_FAILED_REASON *reason) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ExitCode( 
            /* [retval][out] */ int *exitCode) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ProcessDescription( 
            /* [retval][out] */ LPWSTR *processDescription) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_FrameInfosForFailedProcess( 
            /* [retval][out] */ ICoreWebView2FrameInfoCollection **frames) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ProcessFailedEventArgs2Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ProcessFailedEventArgs2 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ProcessFailedEventArgs2 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ProcessFailedEventArgs2 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ProcessFailedKind )( 
            ICoreWebView2ProcessFailedEventArgs2 * This,
            /* [retval][out] */ COREWEBVIEW2_PROCESS_FAILED_KIND *processFailedKind);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Reason )( 
            ICoreWebView2ProcessFailedEventArgs2 * This,
            /* [retval][out] */ COREWEBVIEW2_PROCESS_FAILED_REASON *reason);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ExitCode )( 
            ICoreWebView2ProcessFailedEventArgs2 * This,
            /* [retval][out] */ int *exitCode);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ProcessDescription )( 
            ICoreWebView2ProcessFailedEventArgs2 * This,
            /* [retval][out] */ LPWSTR *processDescription);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_FrameInfosForFailedProcess )( 
            ICoreWebView2ProcessFailedEventArgs2 * This,
            /* [retval][out] */ ICoreWebView2FrameInfoCollection **frames);
        
        END_INTERFACE
    } ICoreWebView2ProcessFailedEventArgs2Vtbl;

    interface ICoreWebView2ProcessFailedEventArgs2
    {
        CONST_VTBL struct ICoreWebView2ProcessFailedEventArgs2Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ProcessFailedEventArgs2_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ProcessFailedEventArgs2_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ProcessFailedEventArgs2_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ProcessFailedEventArgs2_get_ProcessFailedKind(This,processFailedKind)	\
    ( (This)->lpVtbl -> get_ProcessFailedKind(This,processFailedKind) ) 


#define ICoreWebView2ProcessFailedEventArgs2_get_Reason(This,reason)	\
    ( (This)->lpVtbl -> get_Reason(This,reason) ) 

#define ICoreWebView2ProcessFailedEventArgs2_get_ExitCode(This,exitCode)	\
    ( (This)->lpVtbl -> get_ExitCode(This,exitCode) ) 

#define ICoreWebView2ProcessFailedEventArgs2_get_ProcessDescription(This,processDescription)	\
    ( (This)->lpVtbl -> get_ProcessDescription(This,processDescription) ) 

#define ICoreWebView2ProcessFailedEventArgs2_get_FrameInfosForFailedProcess(This,frames)	\
    ( (This)->lpVtbl -> get_FrameInfosForFailedProcess(This,frames) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ProcessFailedEventArgs2_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ProcessFailedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ProcessFailedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ProcessFailedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ProcessFailedEventHandler = {0x79e0aea4,0x990b,0x42d9,{0xaa,0x1d,0x0f,0xcc,0x2e,0x5b,0xc7,0xf1}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("79e0aea4-990b-42d9-aa1d-0fcc2e5bc7f1")
    ICoreWebView2ProcessFailedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2ProcessFailedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ProcessFailedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ProcessFailedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ProcessFailedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ProcessFailedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2ProcessFailedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2ProcessFailedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2ProcessFailedEventHandlerVtbl;

    interface ICoreWebView2ProcessFailedEventHandler
    {
        CONST_VTBL struct ICoreWebView2ProcessFailedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ProcessFailedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ProcessFailedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ProcessFailedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ProcessFailedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ProcessFailedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2RasterizationScaleChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2RasterizationScaleChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2RasterizationScaleChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2RasterizationScaleChangedEventHandler = {0x9c98c8b1,0xac53,0x427e,{0xa3,0x45,0x30,0x49,0xb5,0x52,0x4b,0xbe}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("9c98c8b1-ac53-427e-a345-3049b5524bbe")
    ICoreWebView2RasterizationScaleChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2RasterizationScaleChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2RasterizationScaleChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2RasterizationScaleChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2RasterizationScaleChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2RasterizationScaleChangedEventHandler * This,
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2RasterizationScaleChangedEventHandlerVtbl;

    interface ICoreWebView2RasterizationScaleChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2RasterizationScaleChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2RasterizationScaleChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2RasterizationScaleChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2RasterizationScaleChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2RasterizationScaleChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2RasterizationScaleChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ScriptDialogOpeningEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2ScriptDialogOpeningEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2ScriptDialogOpeningEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ScriptDialogOpeningEventArgs = {0x7390bb70,0xabe0,0x4843,{0x95,0x29,0xf1,0x43,0xb3,0x1b,0x03,0xd6}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("7390bb70-abe0-4843-9529-f143b31b03d6")
    ICoreWebView2ScriptDialogOpeningEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Uri( 
            /* [retval][out] */ LPWSTR *uri) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Kind( 
            /* [retval][out] */ COREWEBVIEW2_SCRIPT_DIALOG_KIND *kind) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Message( 
            /* [retval][out] */ LPWSTR *message) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE Accept( void) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_DefaultText( 
            /* [retval][out] */ LPWSTR *defaultText) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ResultText( 
            /* [retval][out] */ LPWSTR *resultText) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_ResultText( 
            /* [in] */ LPCWSTR resultText) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetDeferral( 
            /* [retval][out] */ ICoreWebView2Deferral **deferral) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ScriptDialogOpeningEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Uri )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [retval][out] */ LPWSTR *uri);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Kind )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_SCRIPT_DIALOG_KIND *kind);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Message )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [retval][out] */ LPWSTR *message);
        
        HRESULT ( STDMETHODCALLTYPE *Accept )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_DefaultText )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [retval][out] */ LPWSTR *defaultText);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ResultText )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [retval][out] */ LPWSTR *resultText);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ResultText )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [in] */ LPCWSTR resultText);
        
        HRESULT ( STDMETHODCALLTYPE *GetDeferral )( 
            ICoreWebView2ScriptDialogOpeningEventArgs * This,
            /* [retval][out] */ ICoreWebView2Deferral **deferral);
        
        END_INTERFACE
    } ICoreWebView2ScriptDialogOpeningEventArgsVtbl;

    interface ICoreWebView2ScriptDialogOpeningEventArgs
    {
        CONST_VTBL struct ICoreWebView2ScriptDialogOpeningEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ScriptDialogOpeningEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ScriptDialogOpeningEventArgs_get_Uri(This,uri)	\
    ( (This)->lpVtbl -> get_Uri(This,uri) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_get_Kind(This,kind)	\
    ( (This)->lpVtbl -> get_Kind(This,kind) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_get_Message(This,message)	\
    ( (This)->lpVtbl -> get_Message(This,message) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_Accept(This)	\
    ( (This)->lpVtbl -> Accept(This) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_get_DefaultText(This,defaultText)	\
    ( (This)->lpVtbl -> get_DefaultText(This,defaultText) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_get_ResultText(This,resultText)	\
    ( (This)->lpVtbl -> get_ResultText(This,resultText) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_put_ResultText(This,resultText)	\
    ( (This)->lpVtbl -> put_ResultText(This,resultText) ) 

#define ICoreWebView2ScriptDialogOpeningEventArgs_GetDeferral(This,deferral)	\
    ( (This)->lpVtbl -> GetDeferral(This,deferral) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ScriptDialogOpeningEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ScriptDialogOpeningEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ScriptDialogOpeningEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ScriptDialogOpeningEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ScriptDialogOpeningEventHandler = {0xef381bf9,0xafa8,0x4e37,{0x91,0xc4,0x8a,0xc4,0x85,0x24,0xbd,0xfb}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("ef381bf9-afa8-4e37-91c4-8ac48524bdfb")
    ICoreWebView2ScriptDialogOpeningEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2ScriptDialogOpeningEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ScriptDialogOpeningEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ScriptDialogOpeningEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ScriptDialogOpeningEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ScriptDialogOpeningEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2ScriptDialogOpeningEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2ScriptDialogOpeningEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2ScriptDialogOpeningEventHandlerVtbl;

    interface ICoreWebView2ScriptDialogOpeningEventHandler
    {
        CONST_VTBL struct ICoreWebView2ScriptDialogOpeningEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ScriptDialogOpeningEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ScriptDialogOpeningEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ScriptDialogOpeningEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ScriptDialogOpeningEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ScriptDialogOpeningEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Settings_INTERFACE_DEFINED__
#define __ICoreWebView2Settings_INTERFACE_DEFINED__

/* interface ICoreWebView2Settings */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Settings = {0xe562e4f0,0xd7fa,0x43ac,{0x8d,0x71,0xc0,0x51,0x50,0x49,0x9f,0x00}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("e562e4f0-d7fa-43ac-8d71-c05150499f00")
    ICoreWebView2Settings : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsScriptEnabled( 
            /* [retval][out] */ BOOL *isScriptEnabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsScriptEnabled( 
            /* [in] */ BOOL isScriptEnabled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsWebMessageEnabled( 
            /* [retval][out] */ BOOL *isWebMessageEnabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsWebMessageEnabled( 
            /* [in] */ BOOL isWebMessageEnabled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_AreDefaultScriptDialogsEnabled( 
            /* [retval][out] */ BOOL *areDefaultScriptDialogsEnabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_AreDefaultScriptDialogsEnabled( 
            /* [in] */ BOOL areDefaultScriptDialogsEnabled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsStatusBarEnabled( 
            /* [retval][out] */ BOOL *isStatusBarEnabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsStatusBarEnabled( 
            /* [in] */ BOOL isStatusBarEnabled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_AreDevToolsEnabled( 
            /* [retval][out] */ BOOL *areDevToolsEnabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_AreDevToolsEnabled( 
            /* [in] */ BOOL areDevToolsEnabled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_AreDefaultContextMenusEnabled( 
            /* [retval][out] */ BOOL *enabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_AreDefaultContextMenusEnabled( 
            /* [in] */ BOOL enabled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_AreHostObjectsAllowed( 
            /* [retval][out] */ BOOL *allowed) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_AreHostObjectsAllowed( 
            /* [in] */ BOOL allowed) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsZoomControlEnabled( 
            /* [retval][out] */ BOOL *enabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsZoomControlEnabled( 
            /* [in] */ BOOL enabled) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsBuiltInErrorPageEnabled( 
            /* [retval][out] */ BOOL *enabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_IsBuiltInErrorPageEnabled( 
            /* [in] */ BOOL enabled) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2SettingsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Settings * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Settings * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Settings * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsScriptEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *isScriptEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsScriptEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL isScriptEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsWebMessageEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *isWebMessageEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsWebMessageEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL isWebMessageEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDefaultScriptDialogsEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *areDefaultScriptDialogsEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDefaultScriptDialogsEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL areDefaultScriptDialogsEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsStatusBarEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *isStatusBarEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsStatusBarEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL isStatusBarEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDevToolsEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *areDevToolsEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDevToolsEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL areDevToolsEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDefaultContextMenusEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDefaultContextMenusEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreHostObjectsAllowed )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *allowed);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreHostObjectsAllowed )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL allowed);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsZoomControlEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsZoomControlEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsBuiltInErrorPageEnabled )( 
            ICoreWebView2Settings * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsBuiltInErrorPageEnabled )( 
            ICoreWebView2Settings * This,
            /* [in] */ BOOL enabled);
        
        END_INTERFACE
    } ICoreWebView2SettingsVtbl;

    interface ICoreWebView2Settings
    {
        CONST_VTBL struct ICoreWebView2SettingsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Settings_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Settings_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Settings_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Settings_get_IsScriptEnabled(This,isScriptEnabled)	\
    ( (This)->lpVtbl -> get_IsScriptEnabled(This,isScriptEnabled) ) 

#define ICoreWebView2Settings_put_IsScriptEnabled(This,isScriptEnabled)	\
    ( (This)->lpVtbl -> put_IsScriptEnabled(This,isScriptEnabled) ) 

#define ICoreWebView2Settings_get_IsWebMessageEnabled(This,isWebMessageEnabled)	\
    ( (This)->lpVtbl -> get_IsWebMessageEnabled(This,isWebMessageEnabled) ) 

#define ICoreWebView2Settings_put_IsWebMessageEnabled(This,isWebMessageEnabled)	\
    ( (This)->lpVtbl -> put_IsWebMessageEnabled(This,isWebMessageEnabled) ) 

#define ICoreWebView2Settings_get_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled)	\
    ( (This)->lpVtbl -> get_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled) ) 

#define ICoreWebView2Settings_put_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled)	\
    ( (This)->lpVtbl -> put_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled) ) 

#define ICoreWebView2Settings_get_IsStatusBarEnabled(This,isStatusBarEnabled)	\
    ( (This)->lpVtbl -> get_IsStatusBarEnabled(This,isStatusBarEnabled) ) 

#define ICoreWebView2Settings_put_IsStatusBarEnabled(This,isStatusBarEnabled)	\
    ( (This)->lpVtbl -> put_IsStatusBarEnabled(This,isStatusBarEnabled) ) 

#define ICoreWebView2Settings_get_AreDevToolsEnabled(This,areDevToolsEnabled)	\
    ( (This)->lpVtbl -> get_AreDevToolsEnabled(This,areDevToolsEnabled) ) 

#define ICoreWebView2Settings_put_AreDevToolsEnabled(This,areDevToolsEnabled)	\
    ( (This)->lpVtbl -> put_AreDevToolsEnabled(This,areDevToolsEnabled) ) 

#define ICoreWebView2Settings_get_AreDefaultContextMenusEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_AreDefaultContextMenusEnabled(This,enabled) ) 

#define ICoreWebView2Settings_put_AreDefaultContextMenusEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_AreDefaultContextMenusEnabled(This,enabled) ) 

#define ICoreWebView2Settings_get_AreHostObjectsAllowed(This,allowed)	\
    ( (This)->lpVtbl -> get_AreHostObjectsAllowed(This,allowed) ) 

#define ICoreWebView2Settings_put_AreHostObjectsAllowed(This,allowed)	\
    ( (This)->lpVtbl -> put_AreHostObjectsAllowed(This,allowed) ) 

#define ICoreWebView2Settings_get_IsZoomControlEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_IsZoomControlEnabled(This,enabled) ) 

#define ICoreWebView2Settings_put_IsZoomControlEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_IsZoomControlEnabled(This,enabled) ) 

#define ICoreWebView2Settings_get_IsBuiltInErrorPageEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_IsBuiltInErrorPageEnabled(This,enabled) ) 

#define ICoreWebView2Settings_put_IsBuiltInErrorPageEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_IsBuiltInErrorPageEnabled(This,enabled) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Settings_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Settings2_INTERFACE_DEFINED__
#define __ICoreWebView2Settings2_INTERFACE_DEFINED__

/* interface ICoreWebView2Settings2 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Settings2 = {0xee9a0f68,0xf46c,0x4e32,{0xac,0x23,0xef,0x8c,0xac,0x22,0x4d,0x2a}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("ee9a0f68-f46c-4e32-ac23-ef8cac224d2a")
    ICoreWebView2Settings2 : public ICoreWebView2Settings
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_UserAgent( 
            /* [retval][out] */ LPWSTR *userAgent) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_UserAgent( 
            /* [in] */ LPCWSTR userAgent) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Settings2Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Settings2 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Settings2 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsScriptEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *isScriptEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsScriptEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL isScriptEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsWebMessageEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *isWebMessageEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsWebMessageEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL isWebMessageEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDefaultScriptDialogsEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *areDefaultScriptDialogsEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDefaultScriptDialogsEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL areDefaultScriptDialogsEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsStatusBarEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *isStatusBarEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsStatusBarEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL isStatusBarEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDevToolsEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *areDevToolsEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDevToolsEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL areDevToolsEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDefaultContextMenusEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDefaultContextMenusEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreHostObjectsAllowed )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *allowed);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreHostObjectsAllowed )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL allowed);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsZoomControlEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsZoomControlEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsBuiltInErrorPageEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsBuiltInErrorPageEnabled )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_UserAgent )( 
            ICoreWebView2Settings2 * This,
            /* [retval][out] */ LPWSTR *userAgent);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_UserAgent )( 
            ICoreWebView2Settings2 * This,
            /* [in] */ LPCWSTR userAgent);
        
        END_INTERFACE
    } ICoreWebView2Settings2Vtbl;

    interface ICoreWebView2Settings2
    {
        CONST_VTBL struct ICoreWebView2Settings2Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Settings2_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Settings2_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Settings2_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Settings2_get_IsScriptEnabled(This,isScriptEnabled)	\
    ( (This)->lpVtbl -> get_IsScriptEnabled(This,isScriptEnabled) ) 

#define ICoreWebView2Settings2_put_IsScriptEnabled(This,isScriptEnabled)	\
    ( (This)->lpVtbl -> put_IsScriptEnabled(This,isScriptEnabled) ) 

#define ICoreWebView2Settings2_get_IsWebMessageEnabled(This,isWebMessageEnabled)	\
    ( (This)->lpVtbl -> get_IsWebMessageEnabled(This,isWebMessageEnabled) ) 

#define ICoreWebView2Settings2_put_IsWebMessageEnabled(This,isWebMessageEnabled)	\
    ( (This)->lpVtbl -> put_IsWebMessageEnabled(This,isWebMessageEnabled) ) 

#define ICoreWebView2Settings2_get_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled)	\
    ( (This)->lpVtbl -> get_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled) ) 

#define ICoreWebView2Settings2_put_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled)	\
    ( (This)->lpVtbl -> put_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled) ) 

#define ICoreWebView2Settings2_get_IsStatusBarEnabled(This,isStatusBarEnabled)	\
    ( (This)->lpVtbl -> get_IsStatusBarEnabled(This,isStatusBarEnabled) ) 

#define ICoreWebView2Settings2_put_IsStatusBarEnabled(This,isStatusBarEnabled)	\
    ( (This)->lpVtbl -> put_IsStatusBarEnabled(This,isStatusBarEnabled) ) 

#define ICoreWebView2Settings2_get_AreDevToolsEnabled(This,areDevToolsEnabled)	\
    ( (This)->lpVtbl -> get_AreDevToolsEnabled(This,areDevToolsEnabled) ) 

#define ICoreWebView2Settings2_put_AreDevToolsEnabled(This,areDevToolsEnabled)	\
    ( (This)->lpVtbl -> put_AreDevToolsEnabled(This,areDevToolsEnabled) ) 

#define ICoreWebView2Settings2_get_AreDefaultContextMenusEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_AreDefaultContextMenusEnabled(This,enabled) ) 

#define ICoreWebView2Settings2_put_AreDefaultContextMenusEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_AreDefaultContextMenusEnabled(This,enabled) ) 

#define ICoreWebView2Settings2_get_AreHostObjectsAllowed(This,allowed)	\
    ( (This)->lpVtbl -> get_AreHostObjectsAllowed(This,allowed) ) 

#define ICoreWebView2Settings2_put_AreHostObjectsAllowed(This,allowed)	\
    ( (This)->lpVtbl -> put_AreHostObjectsAllowed(This,allowed) ) 

#define ICoreWebView2Settings2_get_IsZoomControlEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_IsZoomControlEnabled(This,enabled) ) 

#define ICoreWebView2Settings2_put_IsZoomControlEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_IsZoomControlEnabled(This,enabled) ) 

#define ICoreWebView2Settings2_get_IsBuiltInErrorPageEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_IsBuiltInErrorPageEnabled(This,enabled) ) 

#define ICoreWebView2Settings2_put_IsBuiltInErrorPageEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_IsBuiltInErrorPageEnabled(This,enabled) ) 


#define ICoreWebView2Settings2_get_UserAgent(This,userAgent)	\
    ( (This)->lpVtbl -> get_UserAgent(This,userAgent) ) 

#define ICoreWebView2Settings2_put_UserAgent(This,userAgent)	\
    ( (This)->lpVtbl -> put_UserAgent(This,userAgent) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Settings2_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Settings3_INTERFACE_DEFINED__
#define __ICoreWebView2Settings3_INTERFACE_DEFINED__

/* interface ICoreWebView2Settings3 */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2Settings3 = {0xfdb5ab74,0xaf33,0x4854,{0x84,0xf0,0x0a,0x63,0x1d,0xeb,0x5e,0xba}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("fdb5ab74-af33-4854-84f0-0a631deb5eba")
    ICoreWebView2Settings3 : public ICoreWebView2Settings2
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_AreBrowserAcceleratorKeysEnabled( 
            /* [retval][out] */ BOOL *areBrowserAcceleratorKeysEnabled) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_AreBrowserAcceleratorKeysEnabled( 
            /* [in] */ BOOL areBrowserAcceleratorKeysEnabled) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2Settings3Vtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2Settings3 * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2Settings3 * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsScriptEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *isScriptEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsScriptEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL isScriptEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsWebMessageEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *isWebMessageEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsWebMessageEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL isWebMessageEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDefaultScriptDialogsEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *areDefaultScriptDialogsEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDefaultScriptDialogsEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL areDefaultScriptDialogsEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsStatusBarEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *isStatusBarEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsStatusBarEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL isStatusBarEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDevToolsEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *areDevToolsEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDevToolsEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL areDevToolsEnabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreDefaultContextMenusEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreDefaultContextMenusEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreHostObjectsAllowed )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *allowed);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreHostObjectsAllowed )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL allowed);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsZoomControlEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsZoomControlEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsBuiltInErrorPageEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *enabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_IsBuiltInErrorPageEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL enabled);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_UserAgent )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ LPWSTR *userAgent);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_UserAgent )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ LPCWSTR userAgent);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_AreBrowserAcceleratorKeysEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [retval][out] */ BOOL *areBrowserAcceleratorKeysEnabled);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_AreBrowserAcceleratorKeysEnabled )( 
            ICoreWebView2Settings3 * This,
            /* [in] */ BOOL areBrowserAcceleratorKeysEnabled);
        
        END_INTERFACE
    } ICoreWebView2Settings3Vtbl;

    interface ICoreWebView2Settings3
    {
        CONST_VTBL struct ICoreWebView2Settings3Vtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2Settings3_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2Settings3_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2Settings3_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2Settings3_get_IsScriptEnabled(This,isScriptEnabled)	\
    ( (This)->lpVtbl -> get_IsScriptEnabled(This,isScriptEnabled) ) 

#define ICoreWebView2Settings3_put_IsScriptEnabled(This,isScriptEnabled)	\
    ( (This)->lpVtbl -> put_IsScriptEnabled(This,isScriptEnabled) ) 

#define ICoreWebView2Settings3_get_IsWebMessageEnabled(This,isWebMessageEnabled)	\
    ( (This)->lpVtbl -> get_IsWebMessageEnabled(This,isWebMessageEnabled) ) 

#define ICoreWebView2Settings3_put_IsWebMessageEnabled(This,isWebMessageEnabled)	\
    ( (This)->lpVtbl -> put_IsWebMessageEnabled(This,isWebMessageEnabled) ) 

#define ICoreWebView2Settings3_get_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled)	\
    ( (This)->lpVtbl -> get_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled) ) 

#define ICoreWebView2Settings3_put_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled)	\
    ( (This)->lpVtbl -> put_AreDefaultScriptDialogsEnabled(This,areDefaultScriptDialogsEnabled) ) 

#define ICoreWebView2Settings3_get_IsStatusBarEnabled(This,isStatusBarEnabled)	\
    ( (This)->lpVtbl -> get_IsStatusBarEnabled(This,isStatusBarEnabled) ) 

#define ICoreWebView2Settings3_put_IsStatusBarEnabled(This,isStatusBarEnabled)	\
    ( (This)->lpVtbl -> put_IsStatusBarEnabled(This,isStatusBarEnabled) ) 

#define ICoreWebView2Settings3_get_AreDevToolsEnabled(This,areDevToolsEnabled)	\
    ( (This)->lpVtbl -> get_AreDevToolsEnabled(This,areDevToolsEnabled) ) 

#define ICoreWebView2Settings3_put_AreDevToolsEnabled(This,areDevToolsEnabled)	\
    ( (This)->lpVtbl -> put_AreDevToolsEnabled(This,areDevToolsEnabled) ) 

#define ICoreWebView2Settings3_get_AreDefaultContextMenusEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_AreDefaultContextMenusEnabled(This,enabled) ) 

#define ICoreWebView2Settings3_put_AreDefaultContextMenusEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_AreDefaultContextMenusEnabled(This,enabled) ) 

#define ICoreWebView2Settings3_get_AreHostObjectsAllowed(This,allowed)	\
    ( (This)->lpVtbl -> get_AreHostObjectsAllowed(This,allowed) ) 

#define ICoreWebView2Settings3_put_AreHostObjectsAllowed(This,allowed)	\
    ( (This)->lpVtbl -> put_AreHostObjectsAllowed(This,allowed) ) 

#define ICoreWebView2Settings3_get_IsZoomControlEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_IsZoomControlEnabled(This,enabled) ) 

#define ICoreWebView2Settings3_put_IsZoomControlEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_IsZoomControlEnabled(This,enabled) ) 

#define ICoreWebView2Settings3_get_IsBuiltInErrorPageEnabled(This,enabled)	\
    ( (This)->lpVtbl -> get_IsBuiltInErrorPageEnabled(This,enabled) ) 

#define ICoreWebView2Settings3_put_IsBuiltInErrorPageEnabled(This,enabled)	\
    ( (This)->lpVtbl -> put_IsBuiltInErrorPageEnabled(This,enabled) ) 


#define ICoreWebView2Settings3_get_UserAgent(This,userAgent)	\
    ( (This)->lpVtbl -> get_UserAgent(This,userAgent) ) 

#define ICoreWebView2Settings3_put_UserAgent(This,userAgent)	\
    ( (This)->lpVtbl -> put_UserAgent(This,userAgent) ) 


#define ICoreWebView2Settings3_get_AreBrowserAcceleratorKeysEnabled(This,areBrowserAcceleratorKeysEnabled)	\
    ( (This)->lpVtbl -> get_AreBrowserAcceleratorKeysEnabled(This,areBrowserAcceleratorKeysEnabled) ) 

#define ICoreWebView2Settings3_put_AreBrowserAcceleratorKeysEnabled(This,areBrowserAcceleratorKeysEnabled)	\
    ( (This)->lpVtbl -> put_AreBrowserAcceleratorKeysEnabled(This,areBrowserAcceleratorKeysEnabled) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Settings3_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2SourceChangedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2SourceChangedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2SourceChangedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2SourceChangedEventArgs = {0x31e0e545,0x1dba,0x4266,{0x89,0x14,0xf6,0x38,0x48,0xa1,0xf7,0xd7}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("31e0e545-1dba-4266-8914-f63848a1f7d7")
    ICoreWebView2SourceChangedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsNewDocument( 
            /* [retval][out] */ BOOL *isNewDocument) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2SourceChangedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2SourceChangedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2SourceChangedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2SourceChangedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsNewDocument )( 
            ICoreWebView2SourceChangedEventArgs * This,
            /* [retval][out] */ BOOL *isNewDocument);
        
        END_INTERFACE
    } ICoreWebView2SourceChangedEventArgsVtbl;

    interface ICoreWebView2SourceChangedEventArgs
    {
        CONST_VTBL struct ICoreWebView2SourceChangedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2SourceChangedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2SourceChangedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2SourceChangedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2SourceChangedEventArgs_get_IsNewDocument(This,isNewDocument)	\
    ( (This)->lpVtbl -> get_IsNewDocument(This,isNewDocument) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2SourceChangedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2SourceChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2SourceChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2SourceChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2SourceChangedEventHandler = {0x3c067f9f,0x5388,0x4772,{0x8b,0x48,0x79,0xf7,0xef,0x1a,0xb3,0x7c}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("3c067f9f-5388-4772-8b48-79f7ef1ab37c")
    ICoreWebView2SourceChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2SourceChangedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2SourceChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2SourceChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2SourceChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2SourceChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2SourceChangedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2SourceChangedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2SourceChangedEventHandlerVtbl;

    interface ICoreWebView2SourceChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2SourceChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2SourceChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2SourceChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2SourceChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2SourceChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2SourceChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2TrySuspendCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2TrySuspendCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2TrySuspendCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2TrySuspendCompletedHandler = {0x00F206A7,0x9D17,0x4605,{0x91,0xF6,0x4E,0x8E,0x4D,0xE1,0x92,0xE3}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("00F206A7-9D17-4605-91F6-4E8E4DE192E3")
    ICoreWebView2TrySuspendCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ HRESULT errorCode,
            /* [in] */ BOOL isSuccessful) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2TrySuspendCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2TrySuspendCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2TrySuspendCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2TrySuspendCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2TrySuspendCompletedHandler * This,
            /* [in] */ HRESULT errorCode,
            /* [in] */ BOOL isSuccessful);
        
        END_INTERFACE
    } ICoreWebView2TrySuspendCompletedHandlerVtbl;

    interface ICoreWebView2TrySuspendCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2TrySuspendCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2TrySuspendCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2TrySuspendCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2TrySuspendCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2TrySuspendCompletedHandler_Invoke(This,errorCode,isSuccessful)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,isSuccessful) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2TrySuspendCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebMessageReceivedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2WebMessageReceivedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2WebMessageReceivedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebMessageReceivedEventArgs = {0x0f99a40c,0xe962,0x4207,{0x9e,0x92,0xe3,0xd5,0x42,0xef,0xf8,0x49}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("0f99a40c-e962-4207-9e92-e3d542eff849")
    ICoreWebView2WebMessageReceivedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Source( 
            /* [retval][out] */ LPWSTR *source) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_WebMessageAsJson( 
            /* [retval][out] */ LPWSTR *webMessageAsJson) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE TryGetWebMessageAsString( 
            /* [retval][out] */ LPWSTR *webMessageAsString) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebMessageReceivedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebMessageReceivedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebMessageReceivedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebMessageReceivedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Source )( 
            ICoreWebView2WebMessageReceivedEventArgs * This,
            /* [retval][out] */ LPWSTR *source);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_WebMessageAsJson )( 
            ICoreWebView2WebMessageReceivedEventArgs * This,
            /* [retval][out] */ LPWSTR *webMessageAsJson);
        
        HRESULT ( STDMETHODCALLTYPE *TryGetWebMessageAsString )( 
            ICoreWebView2WebMessageReceivedEventArgs * This,
            /* [retval][out] */ LPWSTR *webMessageAsString);
        
        END_INTERFACE
    } ICoreWebView2WebMessageReceivedEventArgsVtbl;

    interface ICoreWebView2WebMessageReceivedEventArgs
    {
        CONST_VTBL struct ICoreWebView2WebMessageReceivedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebMessageReceivedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebMessageReceivedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebMessageReceivedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebMessageReceivedEventArgs_get_Source(This,source)	\
    ( (This)->lpVtbl -> get_Source(This,source) ) 

#define ICoreWebView2WebMessageReceivedEventArgs_get_WebMessageAsJson(This,webMessageAsJson)	\
    ( (This)->lpVtbl -> get_WebMessageAsJson(This,webMessageAsJson) ) 

#define ICoreWebView2WebMessageReceivedEventArgs_TryGetWebMessageAsString(This,webMessageAsString)	\
    ( (This)->lpVtbl -> TryGetWebMessageAsString(This,webMessageAsString) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebMessageReceivedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebMessageReceivedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2WebMessageReceivedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2WebMessageReceivedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebMessageReceivedEventHandler = {0x57213f19,0x00e6,0x49fa,{0x8e,0x07,0x89,0x8e,0xa0,0x1e,0xcb,0xd2}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("57213f19-00e6-49fa-8e07-898ea01ecbd2")
    ICoreWebView2WebMessageReceivedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2WebMessageReceivedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebMessageReceivedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebMessageReceivedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebMessageReceivedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebMessageReceivedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2WebMessageReceivedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2WebMessageReceivedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2WebMessageReceivedEventHandlerVtbl;

    interface ICoreWebView2WebMessageReceivedEventHandler
    {
        CONST_VTBL struct ICoreWebView2WebMessageReceivedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebMessageReceivedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebMessageReceivedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebMessageReceivedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebMessageReceivedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebMessageReceivedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceRequest_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceRequest_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceRequest */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceRequest = {0x97055cd4,0x512c,0x4264,{0x8b,0x5f,0xe3,0xf4,0x46,0xce,0xa6,0xa5}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("97055cd4-512c-4264-8b5f-e3f446cea6a5")
    ICoreWebView2WebResourceRequest : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Uri( 
            /* [retval][out] */ LPWSTR *uri) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Uri( 
            /* [in] */ LPCWSTR uri) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Method( 
            /* [retval][out] */ LPWSTR *method) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Method( 
            /* [in] */ LPCWSTR method) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Content( 
            /* [retval][out] */ IStream **content) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Content( 
            /* [in] */ IStream *content) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Headers( 
            /* [retval][out] */ ICoreWebView2HttpRequestHeaders **headers) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceRequestVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceRequest * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceRequest * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceRequest * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Uri )( 
            ICoreWebView2WebResourceRequest * This,
            /* [retval][out] */ LPWSTR *uri);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Uri )( 
            ICoreWebView2WebResourceRequest * This,
            /* [in] */ LPCWSTR uri);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Method )( 
            ICoreWebView2WebResourceRequest * This,
            /* [retval][out] */ LPWSTR *method);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Method )( 
            ICoreWebView2WebResourceRequest * This,
            /* [in] */ LPCWSTR method);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Content )( 
            ICoreWebView2WebResourceRequest * This,
            /* [retval][out] */ IStream **content);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Content )( 
            ICoreWebView2WebResourceRequest * This,
            /* [in] */ IStream *content);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Headers )( 
            ICoreWebView2WebResourceRequest * This,
            /* [retval][out] */ ICoreWebView2HttpRequestHeaders **headers);
        
        END_INTERFACE
    } ICoreWebView2WebResourceRequestVtbl;

    interface ICoreWebView2WebResourceRequest
    {
        CONST_VTBL struct ICoreWebView2WebResourceRequestVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceRequest_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceRequest_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceRequest_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceRequest_get_Uri(This,uri)	\
    ( (This)->lpVtbl -> get_Uri(This,uri) ) 

#define ICoreWebView2WebResourceRequest_put_Uri(This,uri)	\
    ( (This)->lpVtbl -> put_Uri(This,uri) ) 

#define ICoreWebView2WebResourceRequest_get_Method(This,method)	\
    ( (This)->lpVtbl -> get_Method(This,method) ) 

#define ICoreWebView2WebResourceRequest_put_Method(This,method)	\
    ( (This)->lpVtbl -> put_Method(This,method) ) 

#define ICoreWebView2WebResourceRequest_get_Content(This,content)	\
    ( (This)->lpVtbl -> get_Content(This,content) ) 

#define ICoreWebView2WebResourceRequest_put_Content(This,content)	\
    ( (This)->lpVtbl -> put_Content(This,content) ) 

#define ICoreWebView2WebResourceRequest_get_Headers(This,headers)	\
    ( (This)->lpVtbl -> get_Headers(This,headers) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceRequest_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceRequestedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceRequestedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceRequestedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceRequestedEventArgs = {0x453e667f,0x12c7,0x49d4,{0xbe,0x6d,0xdd,0xbe,0x79,0x56,0xf5,0x7a}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("453e667f-12c7-49d4-be6d-ddbe7956f57a")
    ICoreWebView2WebResourceRequestedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Request( 
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Response( 
            /* [retval][out] */ ICoreWebView2WebResourceResponse **response) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Response( 
            /* [in] */ ICoreWebView2WebResourceResponse *response) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetDeferral( 
            /* [retval][out] */ ICoreWebView2Deferral **deferral) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ResourceContext( 
            /* [retval][out] */ COREWEBVIEW2_WEB_RESOURCE_CONTEXT *context) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceRequestedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceRequestedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceRequestedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceRequestedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Request )( 
            ICoreWebView2WebResourceRequestedEventArgs * This,
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Response )( 
            ICoreWebView2WebResourceRequestedEventArgs * This,
            /* [retval][out] */ ICoreWebView2WebResourceResponse **response);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Response )( 
            ICoreWebView2WebResourceRequestedEventArgs * This,
            /* [in] */ ICoreWebView2WebResourceResponse *response);
        
        HRESULT ( STDMETHODCALLTYPE *GetDeferral )( 
            ICoreWebView2WebResourceRequestedEventArgs * This,
            /* [retval][out] */ ICoreWebView2Deferral **deferral);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ResourceContext )( 
            ICoreWebView2WebResourceRequestedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_WEB_RESOURCE_CONTEXT *context);
        
        END_INTERFACE
    } ICoreWebView2WebResourceRequestedEventArgsVtbl;

    interface ICoreWebView2WebResourceRequestedEventArgs
    {
        CONST_VTBL struct ICoreWebView2WebResourceRequestedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceRequestedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceRequestedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceRequestedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceRequestedEventArgs_get_Request(This,request)	\
    ( (This)->lpVtbl -> get_Request(This,request) ) 

#define ICoreWebView2WebResourceRequestedEventArgs_get_Response(This,response)	\
    ( (This)->lpVtbl -> get_Response(This,response) ) 

#define ICoreWebView2WebResourceRequestedEventArgs_put_Response(This,response)	\
    ( (This)->lpVtbl -> put_Response(This,response) ) 

#define ICoreWebView2WebResourceRequestedEventArgs_GetDeferral(This,deferral)	\
    ( (This)->lpVtbl -> GetDeferral(This,deferral) ) 

#define ICoreWebView2WebResourceRequestedEventArgs_get_ResourceContext(This,context)	\
    ( (This)->lpVtbl -> get_ResourceContext(This,context) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceRequestedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceRequestedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceRequestedEventHandler = {0xab00b74c,0x15f1,0x4646,{0x80,0xe8,0xe7,0x63,0x41,0xd2,0x5d,0x71}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("ab00b74c-15f1-4646-80e8-e76341d25d71")
    ICoreWebView2WebResourceRequestedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2WebResourceRequestedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceRequestedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceRequestedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceRequestedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceRequestedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2WebResourceRequestedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2WebResourceRequestedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2WebResourceRequestedEventHandlerVtbl;

    interface ICoreWebView2WebResourceRequestedEventHandler
    {
        CONST_VTBL struct ICoreWebView2WebResourceRequestedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceRequestedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceRequestedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceRequestedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceRequestedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceRequestedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponse_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceResponse_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceResponse */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceResponse = {0xaafcc94f,0xfa27,0x48fd,{0x97,0xdf,0x83,0x0e,0xf7,0x5a,0xae,0xc9}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("aafcc94f-fa27-48fd-97df-830ef75aaec9")
    ICoreWebView2WebResourceResponse : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Content( 
            /* [retval][out] */ IStream **content) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_Content( 
            /* [in] */ IStream *content) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Headers( 
            /* [retval][out] */ ICoreWebView2HttpResponseHeaders **headers) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_StatusCode( 
            /* [retval][out] */ int *statusCode) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_StatusCode( 
            /* [in] */ int statusCode) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ReasonPhrase( 
            /* [retval][out] */ LPWSTR *reasonPhrase) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_ReasonPhrase( 
            /* [in] */ LPCWSTR reasonPhrase) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceResponseVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceResponse * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceResponse * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceResponse * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Content )( 
            ICoreWebView2WebResourceResponse * This,
            /* [retval][out] */ IStream **content);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_Content )( 
            ICoreWebView2WebResourceResponse * This,
            /* [in] */ IStream *content);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Headers )( 
            ICoreWebView2WebResourceResponse * This,
            /* [retval][out] */ ICoreWebView2HttpResponseHeaders **headers);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_StatusCode )( 
            ICoreWebView2WebResourceResponse * This,
            /* [retval][out] */ int *statusCode);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_StatusCode )( 
            ICoreWebView2WebResourceResponse * This,
            /* [in] */ int statusCode);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ReasonPhrase )( 
            ICoreWebView2WebResourceResponse * This,
            /* [retval][out] */ LPWSTR *reasonPhrase);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ReasonPhrase )( 
            ICoreWebView2WebResourceResponse * This,
            /* [in] */ LPCWSTR reasonPhrase);
        
        END_INTERFACE
    } ICoreWebView2WebResourceResponseVtbl;

    interface ICoreWebView2WebResourceResponse
    {
        CONST_VTBL struct ICoreWebView2WebResourceResponseVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceResponse_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceResponse_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceResponse_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceResponse_get_Content(This,content)	\
    ( (This)->lpVtbl -> get_Content(This,content) ) 

#define ICoreWebView2WebResourceResponse_put_Content(This,content)	\
    ( (This)->lpVtbl -> put_Content(This,content) ) 

#define ICoreWebView2WebResourceResponse_get_Headers(This,headers)	\
    ( (This)->lpVtbl -> get_Headers(This,headers) ) 

#define ICoreWebView2WebResourceResponse_get_StatusCode(This,statusCode)	\
    ( (This)->lpVtbl -> get_StatusCode(This,statusCode) ) 

#define ICoreWebView2WebResourceResponse_put_StatusCode(This,statusCode)	\
    ( (This)->lpVtbl -> put_StatusCode(This,statusCode) ) 

#define ICoreWebView2WebResourceResponse_get_ReasonPhrase(This,reasonPhrase)	\
    ( (This)->lpVtbl -> get_ReasonPhrase(This,reasonPhrase) ) 

#define ICoreWebView2WebResourceResponse_put_ReasonPhrase(This,reasonPhrase)	\
    ( (This)->lpVtbl -> put_ReasonPhrase(This,reasonPhrase) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceResponse_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseReceivedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceResponseReceivedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceResponseReceivedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceResponseReceivedEventHandler = {0x7DE9898A,0x24F5,0x40C3,{0xA2,0xDE,0xD4,0xF4,0x58,0xE6,0x98,0x28}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("7DE9898A-24F5-40C3-A2DE-D4F458E69828")
    ICoreWebView2WebResourceResponseReceivedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2WebResourceResponseReceivedEventArgs *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceResponseReceivedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceResponseReceivedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceResponseReceivedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2WebResourceResponseReceivedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ ICoreWebView2WebResourceResponseReceivedEventArgs *args);
        
        END_INTERFACE
    } ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl;

    interface ICoreWebView2WebResourceResponseReceivedEventHandler
    {
        CONST_VTBL struct ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceResponseReceivedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceResponseReceivedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceResponseReceivedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceResponseReceivedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceResponseReceivedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseReceivedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceResponseReceivedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceResponseReceivedEventArgs */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceResponseReceivedEventArgs = {0xD1DB483D,0x6796,0x4B8B,{0x80,0xFC,0x13,0x71,0x2B,0xB7,0x16,0xF4}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("D1DB483D-6796-4B8B-80FC-13712BB716F4")
    ICoreWebView2WebResourceResponseReceivedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Request( 
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Response( 
            /* [retval][out] */ ICoreWebView2WebResourceResponseView **response) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceResponseReceivedEventArgsVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceResponseReceivedEventArgs * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceResponseReceivedEventArgs * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceResponseReceivedEventArgs * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Request )( 
            ICoreWebView2WebResourceResponseReceivedEventArgs * This,
            /* [retval][out] */ ICoreWebView2WebResourceRequest **request);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Response )( 
            ICoreWebView2WebResourceResponseReceivedEventArgs * This,
            /* [retval][out] */ ICoreWebView2WebResourceResponseView **response);
        
        END_INTERFACE
    } ICoreWebView2WebResourceResponseReceivedEventArgsVtbl;

    interface ICoreWebView2WebResourceResponseReceivedEventArgs
    {
        CONST_VTBL struct ICoreWebView2WebResourceResponseReceivedEventArgsVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceResponseReceivedEventArgs_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceResponseReceivedEventArgs_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceResponseReceivedEventArgs_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceResponseReceivedEventArgs_get_Request(This,request)	\
    ( (This)->lpVtbl -> get_Request(This,request) ) 

#define ICoreWebView2WebResourceResponseReceivedEventArgs_get_Response(This,response)	\
    ( (This)->lpVtbl -> get_Response(This,response) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceResponseReceivedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseView_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceResponseView_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceResponseView */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceResponseView = {0x79701053,0x7759,0x4162,{0x8F,0x7D,0xF1,0xB3,0xF0,0x84,0x92,0x8D}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("79701053-7759-4162-8F7D-F1B3F084928D")
    ICoreWebView2WebResourceResponseView : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Headers( 
            /* [retval][out] */ ICoreWebView2HttpResponseHeaders **headers) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_StatusCode( 
            /* [retval][out] */ int *statusCode) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ReasonPhrase( 
            /* [retval][out] */ LPWSTR *reasonPhrase) = 0;
        
        virtual HRESULT STDMETHODCALLTYPE GetContent( 
            /* [in] */ ICoreWebView2WebResourceResponseViewGetContentCompletedHandler *handler) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceResponseViewVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceResponseView * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceResponseView * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceResponseView * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Headers )( 
            ICoreWebView2WebResourceResponseView * This,
            /* [retval][out] */ ICoreWebView2HttpResponseHeaders **headers);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_StatusCode )( 
            ICoreWebView2WebResourceResponseView * This,
            /* [retval][out] */ int *statusCode);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ReasonPhrase )( 
            ICoreWebView2WebResourceResponseView * This,
            /* [retval][out] */ LPWSTR *reasonPhrase);
        
        HRESULT ( STDMETHODCALLTYPE *GetContent )( 
            ICoreWebView2WebResourceResponseView * This,
            /* [in] */ ICoreWebView2WebResourceResponseViewGetContentCompletedHandler *handler);
        
        END_INTERFACE
    } ICoreWebView2WebResourceResponseViewVtbl;

    interface ICoreWebView2WebResourceResponseView
    {
        CONST_VTBL struct ICoreWebView2WebResourceResponseViewVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceResponseView_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceResponseView_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceResponseView_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceResponseView_get_Headers(This,headers)	\
    ( (This)->lpVtbl -> get_Headers(This,headers) ) 

#define ICoreWebView2WebResourceResponseView_get_StatusCode(This,statusCode)	\
    ( (This)->lpVtbl -> get_StatusCode(This,statusCode) ) 

#define ICoreWebView2WebResourceResponseView_get_ReasonPhrase(This,reasonPhrase)	\
    ( (This)->lpVtbl -> get_ReasonPhrase(This,reasonPhrase) ) 

#define ICoreWebView2WebResourceResponseView_GetContent(This,handler)	\
    ( (This)->lpVtbl -> GetContent(This,handler) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceResponseView_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2WebResourceResponseViewGetContentCompletedHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WebResourceResponseViewGetContentCompletedHandler = {0x875738E1,0x9FA2,0x40E3,{0x8B,0x74,0x2E,0x89,0x72,0xDD,0x6F,0xE7}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("875738E1-9FA2-40E3-8B74-2E8972DD6FE7")
    ICoreWebView2WebResourceResponseViewGetContentCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ HRESULT errorCode,
            /* [in] */ IStream *content) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WebResourceResponseViewGetContentCompletedHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WebResourceResponseViewGetContentCompletedHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WebResourceResponseViewGetContentCompletedHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2WebResourceResponseViewGetContentCompletedHandler * This,
            /* [in] */ HRESULT errorCode,
            /* [in] */ IStream *content);
        
        END_INTERFACE
    } ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl;

    interface ICoreWebView2WebResourceResponseViewGetContentCompletedHandler
    {
        CONST_VTBL struct ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_Invoke(This,errorCode,content)	\
    ( (This)->lpVtbl -> Invoke(This,errorCode,content) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WebResourceResponseViewGetContentCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WindowCloseRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2WindowCloseRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2WindowCloseRequestedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WindowCloseRequestedEventHandler = {0x5c19e9e0,0x092f,0x486b,{0xaf,0xfa,0xca,0x82,0x31,0x91,0x30,0x39}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("5c19e9e0-092f-486b-affa-ca8231913039")
    ICoreWebView2WindowCloseRequestedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WindowCloseRequestedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WindowCloseRequestedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WindowCloseRequestedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WindowCloseRequestedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2WindowCloseRequestedEventHandler * This,
            /* [in] */ ICoreWebView2 *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2WindowCloseRequestedEventHandlerVtbl;

    interface ICoreWebView2WindowCloseRequestedEventHandler
    {
        CONST_VTBL struct ICoreWebView2WindowCloseRequestedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WindowCloseRequestedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WindowCloseRequestedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WindowCloseRequestedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WindowCloseRequestedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WindowCloseRequestedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WindowFeatures_INTERFACE_DEFINED__
#define __ICoreWebView2WindowFeatures_INTERFACE_DEFINED__

/* interface ICoreWebView2WindowFeatures */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2WindowFeatures = {0x5eaf559f,0xb46e,0x4397,{0x88,0x60,0xe4,0x22,0xf2,0x87,0xff,0x1e}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("5eaf559f-b46e-4397-8860-e422f287ff1e")
    ICoreWebView2WindowFeatures : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HasPosition( 
            /* [retval][out] */ BOOL *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HasSize( 
            /* [retval][out] */ BOOL *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Left( 
            /* [retval][out] */ UINT32 *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Top( 
            /* [retval][out] */ UINT32 *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Height( 
            /* [retval][out] */ UINT32 *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Width( 
            /* [retval][out] */ UINT32 *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ShouldDisplayMenuBar( 
            /* [retval][out] */ BOOL *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ShouldDisplayStatus( 
            /* [retval][out] */ BOOL *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ShouldDisplayToolbar( 
            /* [retval][out] */ BOOL *value) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_ShouldDisplayScrollBars( 
            /* [retval][out] */ BOOL *value) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2WindowFeaturesVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2WindowFeatures * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2WindowFeatures * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2WindowFeatures * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_HasPosition )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ BOOL *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_HasSize )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ BOOL *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Left )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ UINT32 *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Top )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ UINT32 *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Height )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ UINT32 *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_Width )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ UINT32 *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ShouldDisplayMenuBar )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ BOOL *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ShouldDisplayStatus )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ BOOL *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ShouldDisplayToolbar )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ BOOL *value);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_ShouldDisplayScrollBars )( 
            ICoreWebView2WindowFeatures * This,
            /* [retval][out] */ BOOL *value);
        
        END_INTERFACE
    } ICoreWebView2WindowFeaturesVtbl;

    interface ICoreWebView2WindowFeatures
    {
        CONST_VTBL struct ICoreWebView2WindowFeaturesVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2WindowFeatures_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2WindowFeatures_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2WindowFeatures_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2WindowFeatures_get_HasPosition(This,value)	\
    ( (This)->lpVtbl -> get_HasPosition(This,value) ) 

#define ICoreWebView2WindowFeatures_get_HasSize(This,value)	\
    ( (This)->lpVtbl -> get_HasSize(This,value) ) 

#define ICoreWebView2WindowFeatures_get_Left(This,value)	\
    ( (This)->lpVtbl -> get_Left(This,value) ) 

#define ICoreWebView2WindowFeatures_get_Top(This,value)	\
    ( (This)->lpVtbl -> get_Top(This,value) ) 

#define ICoreWebView2WindowFeatures_get_Height(This,value)	\
    ( (This)->lpVtbl -> get_Height(This,value) ) 

#define ICoreWebView2WindowFeatures_get_Width(This,value)	\
    ( (This)->lpVtbl -> get_Width(This,value) ) 

#define ICoreWebView2WindowFeatures_get_ShouldDisplayMenuBar(This,value)	\
    ( (This)->lpVtbl -> get_ShouldDisplayMenuBar(This,value) ) 

#define ICoreWebView2WindowFeatures_get_ShouldDisplayStatus(This,value)	\
    ( (This)->lpVtbl -> get_ShouldDisplayStatus(This,value) ) 

#define ICoreWebView2WindowFeatures_get_ShouldDisplayToolbar(This,value)	\
    ( (This)->lpVtbl -> get_ShouldDisplayToolbar(This,value) ) 

#define ICoreWebView2WindowFeatures_get_ShouldDisplayScrollBars(This,value)	\
    ( (This)->lpVtbl -> get_ShouldDisplayScrollBars(This,value) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2WindowFeatures_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ZoomFactorChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ZoomFactorChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ZoomFactorChangedEventHandler */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2ZoomFactorChangedEventHandler = {0xb52d71d6,0xc4df,0x4543,{0xa9,0x0c,0x64,0xa3,0xe6,0x0f,0x38,0xcb}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("b52d71d6-c4df-4543-a90c-64a3e60f38cb")
    ICoreWebView2ZoomFactorChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke( 
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ IUnknown *args) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2ZoomFactorChangedEventHandlerVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2ZoomFactorChangedEventHandler * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2ZoomFactorChangedEventHandler * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2ZoomFactorChangedEventHandler * This);
        
        HRESULT ( STDMETHODCALLTYPE *Invoke )( 
            ICoreWebView2ZoomFactorChangedEventHandler * This,
            /* [in] */ ICoreWebView2Controller *sender,
            /* [in] */ IUnknown *args);
        
        END_INTERFACE
    } ICoreWebView2ZoomFactorChangedEventHandlerVtbl;

    interface ICoreWebView2ZoomFactorChangedEventHandler
    {
        CONST_VTBL struct ICoreWebView2ZoomFactorChangedEventHandlerVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2ZoomFactorChangedEventHandler_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2ZoomFactorChangedEventHandler_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2ZoomFactorChangedEventHandler_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2ZoomFactorChangedEventHandler_Invoke(This,sender,args)	\
    ( (This)->lpVtbl -> Invoke(This,sender,args) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ZoomFactorChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CompositionControllerInterop_INTERFACE_DEFINED__
#define __ICoreWebView2CompositionControllerInterop_INTERFACE_DEFINED__

/* interface ICoreWebView2CompositionControllerInterop */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2CompositionControllerInterop = {0x8e9922ce,0x9c80,0x42e6,{0xba,0xd7,0xfc,0xeb,0xf2,0x91,0xa4,0x95}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("8e9922ce-9c80-42e6-bad7-fcebf291a495")
    ICoreWebView2CompositionControllerInterop : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_UIAProvider( 
            /* [retval][out] */ IUnknown **provider) = 0;
        
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_RootVisualTarget( 
            /* [retval][out] */ IUnknown **target) = 0;
        
        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_RootVisualTarget( 
            /* [in] */ IUnknown *target) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2CompositionControllerInteropVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2CompositionControllerInterop * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2CompositionControllerInterop * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2CompositionControllerInterop * This);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_UIAProvider )( 
            ICoreWebView2CompositionControllerInterop * This,
            /* [retval][out] */ IUnknown **provider);
        
        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_RootVisualTarget )( 
            ICoreWebView2CompositionControllerInterop * This,
            /* [retval][out] */ IUnknown **target);
        
        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_RootVisualTarget )( 
            ICoreWebView2CompositionControllerInterop * This,
            /* [in] */ IUnknown *target);
        
        END_INTERFACE
    } ICoreWebView2CompositionControllerInteropVtbl;

    interface ICoreWebView2CompositionControllerInterop
    {
        CONST_VTBL struct ICoreWebView2CompositionControllerInteropVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2CompositionControllerInterop_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2CompositionControllerInterop_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2CompositionControllerInterop_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2CompositionControllerInterop_get_UIAProvider(This,provider)	\
    ( (This)->lpVtbl -> get_UIAProvider(This,provider) ) 

#define ICoreWebView2CompositionControllerInterop_get_RootVisualTarget(This,target)	\
    ( (This)->lpVtbl -> get_RootVisualTarget(This,target) ) 

#define ICoreWebView2CompositionControllerInterop_put_RootVisualTarget(This,target)	\
    ( (This)->lpVtbl -> put_RootVisualTarget(This,target) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CompositionControllerInterop_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2EnvironmentInterop_INTERFACE_DEFINED__
#define __ICoreWebView2EnvironmentInterop_INTERFACE_DEFINED__

/* interface ICoreWebView2EnvironmentInterop */
/* [unique][object][uuid] */ 


EXTERN_C __declspec(selectany) const IID IID_ICoreWebView2EnvironmentInterop = {0xee503a63,0xc1e2,0x4fbf,{0x8a,0x4d,0x82,0x4e,0x95,0xf8,0xbb,0x13}};

#if defined(__cplusplus) && !defined(CINTERFACE)
    
    MIDL_INTERFACE("ee503a63-c1e2-4fbf-8a4d-824e95f8bb13")
    ICoreWebView2EnvironmentInterop : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE GetProviderForHwnd( 
            /* [in] */ HWND hwnd,
            /* [retval][out] */ IUnknown **provider) = 0;
        
    };
    
    
#else 	/* C style interface */

    typedef struct ICoreWebView2EnvironmentInteropVtbl
    {
        BEGIN_INTERFACE
        
        HRESULT ( STDMETHODCALLTYPE *QueryInterface )( 
            ICoreWebView2EnvironmentInterop * This,
            /* [in] */ REFIID riid,
            /* [annotation][iid_is][out] */ 
            _COM_Outptr_  void **ppvObject);
        
        ULONG ( STDMETHODCALLTYPE *AddRef )( 
            ICoreWebView2EnvironmentInterop * This);
        
        ULONG ( STDMETHODCALLTYPE *Release )( 
            ICoreWebView2EnvironmentInterop * This);
        
        HRESULT ( STDMETHODCALLTYPE *GetProviderForHwnd )( 
            ICoreWebView2EnvironmentInterop * This,
            /* [in] */ HWND hwnd,
            /* [retval][out] */ IUnknown **provider);
        
        END_INTERFACE
    } ICoreWebView2EnvironmentInteropVtbl;

    interface ICoreWebView2EnvironmentInterop
    {
        CONST_VTBL struct ICoreWebView2EnvironmentInteropVtbl *lpVtbl;
    };

    

#ifdef COBJMACROS


#define ICoreWebView2EnvironmentInterop_QueryInterface(This,riid,ppvObject)	\
    ( (This)->lpVtbl -> QueryInterface(This,riid,ppvObject) ) 

#define ICoreWebView2EnvironmentInterop_AddRef(This)	\
    ( (This)->lpVtbl -> AddRef(This) ) 

#define ICoreWebView2EnvironmentInterop_Release(This)	\
    ( (This)->lpVtbl -> Release(This) ) 


#define ICoreWebView2EnvironmentInterop_GetProviderForHwnd(This,hwnd,provider)	\
    ( (This)->lpVtbl -> GetProviderForHwnd(This,hwnd,provider) ) 

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2EnvironmentInterop_INTERFACE_DEFINED__ */

#endif /* __WebView2_LIBRARY_DEFINED__ */

/* Additional Prototypes for ALL interfaces */

/* end of Additional Prototypes */

#ifdef __cplusplus
}
#endif

#endif


