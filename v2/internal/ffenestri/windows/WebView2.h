

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


#ifndef __ICoreWebView2Controller_FWD_DEFINED__
#define __ICoreWebView2Controller_FWD_DEFINED__
typedef interface ICoreWebView2Controller ICoreWebView2Controller;

#endif 	/* __ICoreWebView2Controller_FWD_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventArgs_FWD_DEFINED__
#define __ICoreWebView2ContentLoadingEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2ContentLoadingEventArgs ICoreWebView2ContentLoadingEventArgs;

#endif 	/* __ICoreWebView2ContentLoadingEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventHandler_FWD_DEFINED__
#define __ICoreWebView2ContentLoadingEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ContentLoadingEventHandler ICoreWebView2ContentLoadingEventHandler;

#endif 	/* __ICoreWebView2ContentLoadingEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2DocumentTitleChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2DocumentTitleChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2DocumentTitleChangedEventHandler ICoreWebView2DocumentTitleChangedEventHandler;

#endif 	/* __ICoreWebView2DocumentTitleChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2ContainsFullScreenElementChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2ContainsFullScreenElementChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ContainsFullScreenElementChangedEventHandler ICoreWebView2ContainsFullScreenElementChangedEventHandler;

#endif 	/* __ICoreWebView2ContainsFullScreenElementChangedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2CreateCoreWebView2ControllerCompletedHandler ICoreWebView2CreateCoreWebView2ControllerCompletedHandler;

#endif 	/* __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler;

#endif 	/* __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_FWD_DEFINED__ */


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


#ifndef __ICoreWebView2EnvironmentOptions_FWD_DEFINED__
#define __ICoreWebView2EnvironmentOptions_FWD_DEFINED__
typedef interface ICoreWebView2EnvironmentOptions ICoreWebView2EnvironmentOptions;

#endif 	/* __ICoreWebView2EnvironmentOptions_FWD_DEFINED__ */


#ifndef __ICoreWebView2ExecuteScriptCompletedHandler_FWD_DEFINED__
#define __ICoreWebView2ExecuteScriptCompletedHandler_FWD_DEFINED__
typedef interface ICoreWebView2ExecuteScriptCompletedHandler ICoreWebView2ExecuteScriptCompletedHandler;

#endif 	/* __ICoreWebView2ExecuteScriptCompletedHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2FocusChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2FocusChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2FocusChangedEventHandler ICoreWebView2FocusChangedEventHandler;

#endif 	/* __ICoreWebView2FocusChangedEventHandler_FWD_DEFINED__ */


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


#ifndef __ICoreWebView2ProcessFailedEventArgs_FWD_DEFINED__
#define __ICoreWebView2ProcessFailedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2ProcessFailedEventArgs ICoreWebView2ProcessFailedEventArgs;

#endif 	/* __ICoreWebView2ProcessFailedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2ProcessFailedEventHandler_FWD_DEFINED__
#define __ICoreWebView2ProcessFailedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ProcessFailedEventHandler ICoreWebView2ProcessFailedEventHandler;

#endif 	/* __ICoreWebView2ProcessFailedEventHandler_FWD_DEFINED__ */


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


#ifndef __ICoreWebView2SourceChangedEventArgs_FWD_DEFINED__
#define __ICoreWebView2SourceChangedEventArgs_FWD_DEFINED__
typedef interface ICoreWebView2SourceChangedEventArgs ICoreWebView2SourceChangedEventArgs;

#endif 	/* __ICoreWebView2SourceChangedEventArgs_FWD_DEFINED__ */


#ifndef __ICoreWebView2SourceChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2SourceChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2SourceChangedEventHandler ICoreWebView2SourceChangedEventHandler;

#endif 	/* __ICoreWebView2SourceChangedEventHandler_FWD_DEFINED__ */


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


#ifndef __ICoreWebView2WindowCloseRequestedEventHandler_FWD_DEFINED__
#define __ICoreWebView2WindowCloseRequestedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2WindowCloseRequestedEventHandler ICoreWebView2WindowCloseRequestedEventHandler;

#endif 	/* __ICoreWebView2WindowCloseRequestedEventHandler_FWD_DEFINED__ */


#ifndef __ICoreWebView2ZoomFactorChangedEventHandler_FWD_DEFINED__
#define __ICoreWebView2ZoomFactorChangedEventHandler_FWD_DEFINED__
typedef interface ICoreWebView2ZoomFactorChangedEventHandler ICoreWebView2ZoomFactorChangedEventHandler;

#endif 	/* __ICoreWebView2ZoomFactorChangedEventHandler_FWD_DEFINED__ */


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
        COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_UNRESPONSIVE	= ( COREWEBVIEW2_PROCESS_FAILED_KIND_RENDER_PROCESS_EXITED + 1 )
    } 	COREWEBVIEW2_PROCESS_FAILED_KIND;

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

STDAPI CreateCoreWebView2EnvironmentWithOptions(PCWSTR browserExecutableFolder, PCWSTR userDataFolder, ICoreWebView2EnvironmentOptions* environmentOptions, ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler* environment_created_handler);
STDAPI CreateCoreWebView2Environment(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler* environment_created_handler);
STDAPI GetAvailableCoreWebView2BrowserVersionString(PCWSTR browserExecutableFolder, LPWSTR* versionInfo);
STDAPI CompareBrowserVersions(PCWSTR version1, PCWSTR version2, int* result);

EXTERN_C const IID LIBID_WebView2;

#ifndef __ICoreWebView2AcceleratorKeyPressedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2AcceleratorKeyPressedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2AcceleratorKeyPressedEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2AcceleratorKeyPressedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("9224476E-D8C3-4EB7-BB65-2FD7792B27CE")
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


EXTERN_C const IID IID_ICoreWebView2AcceleratorKeyPressedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("A7D303F9-503C-4B7E-BC40-5C7CE6CABAAA")
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


EXTERN_C const IID IID_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("7082ABED-0591-428F-A722-60C2F814546B")
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


EXTERN_C const IID IID_ICoreWebView2CallDevToolsProtocolMethodCompletedHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("C20CF895-BA7C-493B-AB2E-8A6E3A3602A2")
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


EXTERN_C const IID IID_ICoreWebView2CapturePreviewCompletedHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("DCED64F8-D9C7-4A3C-B9FD-FBBCA0B43496")
    ICoreWebView2CapturePreviewCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke(
            /* [in] */ HRESULT result) = 0;

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
            /* [in] */ HRESULT result);

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


#define ICoreWebView2CapturePreviewCompletedHandler_Invoke(This,result)	\
    ( (This)->lpVtbl -> Invoke(This,result) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CapturePreviewCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2_INTERFACE_DEFINED__
#define __ICoreWebView2_INTERFACE_DEFINED__

/* interface ICoreWebView2 */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("189B8AAF-0426-4748-B9AD-243F537EB46B")
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


#ifndef __ICoreWebView2Controller_INTERFACE_DEFINED__
#define __ICoreWebView2Controller_INTERFACE_DEFINED__

/* interface ICoreWebView2Controller */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2Controller;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("7CCC5C7F-8351-4572-9077-9C1C80913835")
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
            /* [retval][out] */ HWND *topLevelWindow) = 0;

        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_ParentWindow(
            /* [in] */ HWND topLevelWindow) = 0;

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
            /* [retval][out] */ HWND *topLevelWindow);

        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_ParentWindow )(
            ICoreWebView2Controller * This,
            /* [in] */ HWND topLevelWindow);

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

#define ICoreWebView2Controller_get_ParentWindow(This,topLevelWindow)	\
    ( (This)->lpVtbl -> get_ParentWindow(This,topLevelWindow) )

#define ICoreWebView2Controller_put_ParentWindow(This,topLevelWindow)	\
    ( (This)->lpVtbl -> put_ParentWindow(This,topLevelWindow) )

#define ICoreWebView2Controller_NotifyParentWindowPositionChanged(This)	\
    ( (This)->lpVtbl -> NotifyParentWindowPositionChanged(This) )

#define ICoreWebView2Controller_Close(This)	\
    ( (This)->lpVtbl -> Close(This) )

#define ICoreWebView2Controller_get_CoreWebView2(This,coreWebView2)	\
    ( (This)->lpVtbl -> get_CoreWebView2(This,coreWebView2) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2Controller_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2ContentLoadingEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2ContentLoadingEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ContentLoadingEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("2A800835-2179-45D6-A745-6657E9A546B9")
    ICoreWebView2ContentLoadingEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsErrorPage(
            /* [retval][out] */ BOOL *isErrorPage) = 0;

        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_NavigationId(
            /* [retval][out] */ UINT64 *navigation_id) = 0;

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
            /* [retval][out] */ UINT64 *navigation_id);

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

#define ICoreWebView2ContentLoadingEventArgs_get_NavigationId(This,navigation_id)	\
    ( (This)->lpVtbl -> get_NavigationId(This,navigation_id) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ContentLoadingEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ContentLoadingEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ContentLoadingEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ContentLoadingEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ContentLoadingEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("7AF5CC82-AE19-4964-BD71-B9BC5F03E85D")
    ICoreWebView2ContentLoadingEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke(
            /* [in] */ ICoreWebView2 *webview,
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
            /* [in] */ ICoreWebView2 *webview,
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


#define ICoreWebView2ContentLoadingEventHandler_Invoke(This,webview,args)	\
    ( (This)->lpVtbl -> Invoke(This,webview,args) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2ContentLoadingEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2DocumentTitleChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2DocumentTitleChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2DocumentTitleChangedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2DocumentTitleChangedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("6423D6B1-5A57-46C5-BA46-DBB3735EE7C9")
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


#ifndef __ICoreWebView2ContainsFullScreenElementChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ContainsFullScreenElementChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ContainsFullScreenElementChangedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ContainsFullScreenElementChangedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("120888E3-4CAD-4EC2-B627-B2016D05612D")
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


#ifndef __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CreateCoreWebView2ControllerCompletedHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("86EF6808-3C3F-4C6F-975E-8CE0B98F70BA")
    ICoreWebView2CreateCoreWebView2ControllerCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke(
            HRESULT result,
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
            HRESULT result,
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


#define ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_Invoke(This,result,createdController)	\
    ( (This)->lpVtbl -> Invoke(This,result,createdController) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CreateCoreWebView2ControllerCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("8B4F98CE-DB0D-4E71-85FD-C4C4EF1F2630")
    ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke(
            HRESULT result,
            ICoreWebView2Environment *created_environment) = 0;

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
            HRESULT result,
            ICoreWebView2Environment *created_environment);

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


#define ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_Invoke(This,result,created_environment)	\
    ( (This)->lpVtbl -> Invoke(This,result,created_environment) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2Deferral_INTERFACE_DEFINED__
#define __ICoreWebView2Deferral_INTERFACE_DEFINED__

/* interface ICoreWebView2Deferral */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2Deferral;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("A7ED8BF0-3EC9-4E39-8427-3D6F157BD285")
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


EXTERN_C const IID IID_ICoreWebView2DevToolsProtocolEventReceivedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("F661B1C2-5FF5-4700-B723-C439034539B4")
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


EXTERN_C const IID IID_ICoreWebView2DevToolsProtocolEventReceivedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("8E1DED79-A40B-4271-8BE6-57640C167F4A")
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


EXTERN_C const IID IID_ICoreWebView2DevToolsProtocolEventReceiver;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("FE59C48C-540C-4A3C-8898-8E1602E0055D")
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


EXTERN_C const IID IID_ICoreWebView2Environment;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("DA66D884-6DA8-410E-9630-8C48F8B3A40E")
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


#ifndef __ICoreWebView2EnvironmentOptions_INTERFACE_DEFINED__
#define __ICoreWebView2EnvironmentOptions_INTERFACE_DEFINED__

/* interface ICoreWebView2EnvironmentOptions */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2EnvironmentOptions;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("97E9FBD9-646A-4B75-8682-149B71DACE59")
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

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2EnvironmentOptions_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2ExecuteScriptCompletedHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ExecuteScriptCompletedHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ExecuteScriptCompletedHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ExecuteScriptCompletedHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("3B717C93-3ED5-4450-9B13-7F56AA367AC7")
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


#ifndef __ICoreWebView2FocusChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2FocusChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2FocusChangedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2FocusChangedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("76E67C71-663F-4C17-B71A-9381CCF3B94B")
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


#ifndef __ICoreWebView2HistoryChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2HistoryChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2HistoryChangedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2HistoryChangedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("54C9B7D7-D9E9-4158-861F-F97E1C3C6631")
    ICoreWebView2HistoryChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke(
            /* [in] */ ICoreWebView2 *webview,
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
            /* [in] */ ICoreWebView2 *webview,
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


#define ICoreWebView2HistoryChangedEventHandler_Invoke(This,webview,args)	\
    ( (This)->lpVtbl -> Invoke(This,webview,args) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2HistoryChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2HttpHeadersCollectionIterator_INTERFACE_DEFINED__
#define __ICoreWebView2HttpHeadersCollectionIterator_INTERFACE_DEFINED__

/* interface ICoreWebView2HttpHeadersCollectionIterator */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2HttpHeadersCollectionIterator;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("4212F3A7-0FBC-4C9C-8118-17ED6370C1B3")
    ICoreWebView2HttpHeadersCollectionIterator : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE GetCurrentHeader(
            /* [out] */ LPWSTR *name,
            /* [out] */ LPWSTR *value) = 0;

        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_HasCurrentHeader(
            /* [retval][out] */ BOOL *hasCurrent) = 0;

        virtual HRESULT STDMETHODCALLTYPE MoveNext(
            /* [out] */ BOOL *hasNext) = 0;

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
            /* [out] */ BOOL *hasNext);

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


EXTERN_C const IID IID_ICoreWebView2HttpRequestHeaders;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("2C1F04DF-C90E-49E4-BD25-4A659300337B")
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


EXTERN_C const IID IID_ICoreWebView2HttpResponseHeaders;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("B5F6D4D5-1BFF-4869-85B8-158153017B04")
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


#ifndef __ICoreWebView2MoveFocusRequestedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2MoveFocusRequestedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2MoveFocusRequestedEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2MoveFocusRequestedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("71922903-B180-49D0-AED2-C9F9D10064B1")
    ICoreWebView2MoveFocusRequestedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Reason(
            /* [retval][out] */ COREWEBVIEW2_MOVE_FOCUS_REASON *value) = 0;

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
            /* [retval][out] */ COREWEBVIEW2_MOVE_FOCUS_REASON *value);

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


#define ICoreWebView2MoveFocusRequestedEventArgs_get_Reason(This,value)	\
    ( (This)->lpVtbl -> get_Reason(This,value) )

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


EXTERN_C const IID IID_ICoreWebView2MoveFocusRequestedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("4B21D6DD-3DE7-47B0-8019-7D3ACE6E3631")
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


EXTERN_C const IID IID_ICoreWebView2NavigationCompletedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("361F5621-EA7F-4C55-95EC-3C5E6992EA4A")
    ICoreWebView2NavigationCompletedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsSuccess(
            /* [retval][out] */ BOOL *isSuccess) = 0;

        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_WebErrorStatus(
            /* [retval][out] */ COREWEBVIEW2_WEB_ERROR_STATUS *COREWEBVIEW2_WEB_ERROR_STATUS) = 0;

        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_NavigationId(
            /* [retval][out] */ UINT64 *navigation_id) = 0;

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
            /* [retval][out] */ COREWEBVIEW2_WEB_ERROR_STATUS *COREWEBVIEW2_WEB_ERROR_STATUS);

        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_NavigationId )(
            ICoreWebView2NavigationCompletedEventArgs * This,
            /* [retval][out] */ UINT64 *navigation_id);

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

#define ICoreWebView2NavigationCompletedEventArgs_get_WebErrorStatus(This,COREWEBVIEW2_WEB_ERROR_STATUS)	\
    ( (This)->lpVtbl -> get_WebErrorStatus(This,COREWEBVIEW2_WEB_ERROR_STATUS) )

#define ICoreWebView2NavigationCompletedEventArgs_get_NavigationId(This,navigation_id)	\
    ( (This)->lpVtbl -> get_NavigationId(This,navigation_id) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NavigationCompletedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NavigationCompletedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2NavigationCompletedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2NavigationCompletedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2NavigationCompletedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("9F921239-20C4-455F-9E3F-6047A50E248B")
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


EXTERN_C const IID IID_ICoreWebView2NavigationStartingEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("EE1938CE-D385-4CB0-854B-F498F78C3D88")
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
            /* [retval][out] */ UINT64 *navigation_id) = 0;

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
            /* [retval][out] */ UINT64 *navigation_id);

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

#define ICoreWebView2NavigationStartingEventArgs_get_NavigationId(This,navigation_id)	\
    ( (This)->lpVtbl -> get_NavigationId(This,navigation_id) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NavigationStartingEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NavigationStartingEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2NavigationStartingEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2NavigationStartingEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2NavigationStartingEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("073337A4-64D2-4C7E-AC9F-987F0F613497")
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


EXTERN_C const IID IID_ICoreWebView2NewBrowserVersionAvailableEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("E82E8242-EE39-4A57-A065-E13256D60342")
    ICoreWebView2NewBrowserVersionAvailableEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke(
            /* [in] */ ICoreWebView2Environment *webviewEnvironment,
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
            /* [in] */ ICoreWebView2Environment *webviewEnvironment,
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


#define ICoreWebView2NewBrowserVersionAvailableEventHandler_Invoke(This,webviewEnvironment,args)	\
    ( (This)->lpVtbl -> Invoke(This,webviewEnvironment,args) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NewBrowserVersionAvailableEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NewWindowRequestedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2NewWindowRequestedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2NewWindowRequestedEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2NewWindowRequestedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("9EDC7F5F-C6EA-4F3C-827B-A8880794C0A9")
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

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2NewWindowRequestedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2NewWindowRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2NewWindowRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2NewWindowRequestedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2NewWindowRequestedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("ACAA30EF-A40C-47BD-9CB9-D9C2AADC9FCB")
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


EXTERN_C const IID IID_ICoreWebView2PermissionRequestedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("774B5EA1-3FAD-435C-B1FC-A77D1ACD5EAF")
    ICoreWebView2PermissionRequestedEventArgs : public IUnknown
    {
    public:
        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_Uri(
            /* [retval][out] */ LPWSTR *uri) = 0;

        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_PermissionKind(
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_KIND *value) = 0;

        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_IsUserInitiated(
            /* [retval][out] */ BOOL *isUserInitiated) = 0;

        virtual /* [propget] */ HRESULT STDMETHODCALLTYPE get_State(
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_STATE *value) = 0;

        virtual /* [propput] */ HRESULT STDMETHODCALLTYPE put_State(
            /* [in] */ COREWEBVIEW2_PERMISSION_STATE value) = 0;

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
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_KIND *value);

        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_IsUserInitiated )(
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [retval][out] */ BOOL *isUserInitiated);

        /* [propget] */ HRESULT ( STDMETHODCALLTYPE *get_State )(
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [retval][out] */ COREWEBVIEW2_PERMISSION_STATE *value);

        /* [propput] */ HRESULT ( STDMETHODCALLTYPE *put_State )(
            ICoreWebView2PermissionRequestedEventArgs * This,
            /* [in] */ COREWEBVIEW2_PERMISSION_STATE value);

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

#define ICoreWebView2PermissionRequestedEventArgs_get_PermissionKind(This,value)	\
    ( (This)->lpVtbl -> get_PermissionKind(This,value) )

#define ICoreWebView2PermissionRequestedEventArgs_get_IsUserInitiated(This,isUserInitiated)	\
    ( (This)->lpVtbl -> get_IsUserInitiated(This,isUserInitiated) )

#define ICoreWebView2PermissionRequestedEventArgs_get_State(This,value)	\
    ( (This)->lpVtbl -> get_State(This,value) )

#define ICoreWebView2PermissionRequestedEventArgs_put_State(This,value)	\
    ( (This)->lpVtbl -> put_State(This,value) )

#define ICoreWebView2PermissionRequestedEventArgs_GetDeferral(This,deferral)	\
    ( (This)->lpVtbl -> GetDeferral(This,deferral) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2PermissionRequestedEventArgs_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2PermissionRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2PermissionRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2PermissionRequestedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2PermissionRequestedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("543B4ADE-9B0B-4748-9AB7-D76481B223AA")
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


#ifndef __ICoreWebView2ProcessFailedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2ProcessFailedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2ProcessFailedEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ProcessFailedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("EA45D1F4-75C0-471F-A6E9-803FBFF8FEF2")
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


#ifndef __ICoreWebView2ProcessFailedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ProcessFailedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ProcessFailedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ProcessFailedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("7D2183F9-CCA8-40F2-91A9-EAFAD32C8A9B")
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


#ifndef __ICoreWebView2ScriptDialogOpeningEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2ScriptDialogOpeningEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2ScriptDialogOpeningEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ScriptDialogOpeningEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("B8F6356E-24DC-4D74-90FE-AD071E11CB91")
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


EXTERN_C const IID IID_ICoreWebView2ScriptDialogOpeningEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("72D93789-2727-4A9B-A4FC-1B2609CBCBE3")
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


EXTERN_C const IID IID_ICoreWebView2Settings;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("203FBA37-6850-4DCC-A25A-58A351AC625D")
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


#ifndef __ICoreWebView2SourceChangedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2SourceChangedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2SourceChangedEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2SourceChangedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("BD9A4BFB-BE19-40BD-968B-EBCF0D727EF3")
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


EXTERN_C const IID IID_ICoreWebView2SourceChangedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("8FEDD1A7-3A33-416F-AF81-881EEB001433")
    ICoreWebView2SourceChangedEventHandler : public IUnknown
    {
    public:
        virtual HRESULT STDMETHODCALLTYPE Invoke(
            /* [in] */ ICoreWebView2 *webview,
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
            /* [in] */ ICoreWebView2 *webview,
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


#define ICoreWebView2SourceChangedEventHandler_Invoke(This,webview,args)	\
    ( (This)->lpVtbl -> Invoke(This,webview,args) )

#endif /* COBJMACROS */


#endif 	/* C style interface */




#endif 	/* __ICoreWebView2SourceChangedEventHandler_INTERFACE_DEFINED__ */


#ifndef __ICoreWebView2WebMessageReceivedEventArgs_INTERFACE_DEFINED__
#define __ICoreWebView2WebMessageReceivedEventArgs_INTERFACE_DEFINED__

/* interface ICoreWebView2WebMessageReceivedEventArgs */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2WebMessageReceivedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("B263B5AE-9C54-4B75-B632-40AE1A0B6912")
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


EXTERN_C const IID IID_ICoreWebView2WebMessageReceivedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("199328C8-9964-4F5F-84E6-E875B1B763D6")
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


EXTERN_C const IID IID_ICoreWebView2WebResourceRequest;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("11B02254-B827-49F6-8974-30F6E6C55AF6")
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


EXTERN_C const IID IID_ICoreWebView2WebResourceRequestedEventArgs;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("2D7B3282-83B1-41CA-8BBF-FF18F6BFE320")
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


EXTERN_C const IID IID_ICoreWebView2WebResourceRequestedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("F6DC79F2-E1FA-4534-8968-4AFF10BBAA32")
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


EXTERN_C const IID IID_ICoreWebView2WebResourceResponse;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("5953D1FC-B08F-46DD-AFD3-66B172419CD0")
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


#ifndef __ICoreWebView2WindowCloseRequestedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2WindowCloseRequestedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2WindowCloseRequestedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2WindowCloseRequestedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("63C89928-AD32-4421-A0E4-EC99B34AA97E")
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


#ifndef __ICoreWebView2ZoomFactorChangedEventHandler_INTERFACE_DEFINED__
#define __ICoreWebView2ZoomFactorChangedEventHandler_INTERFACE_DEFINED__

/* interface ICoreWebView2ZoomFactorChangedEventHandler */
/* [unique][object][uuid] */


EXTERN_C const IID IID_ICoreWebView2ZoomFactorChangedEventHandler;

#if defined(__cplusplus) && !defined(CINTERFACE)

    MIDL_INTERFACE("F1828246-8B98-4274-B708-ECDB6BF3843A")
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

#endif /* __WebView2_LIBRARY_DEFINED__ */

/* Additional Prototypes for ALL interfaces */

/* end of Additional Prototypes */

#ifdef __cplusplus
}
#endif

#endif

