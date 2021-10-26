// Copyright (C) Microsoft Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef __core_webview2_environment_options_h__
#define __core_webview2_environment_options_h__

#include <objbase.h>
#include <wrl/implements.h>

#include "webview2.h"
#define CORE_WEBVIEW_TARGET_PRODUCT_VERSION L"91.0.864.35"

#define COREWEBVIEW2ENVIRONMENTOPTIONS_STRING_PROPERTY(p)     \
 public:                                                      \
  HRESULT STDMETHODCALLTYPE get_##p(LPWSTR* value) override { \
    if (!value)                                               \
      return E_POINTER;                                       \
    *value = m_##p.Copy();                                    \
    if ((*value == nullptr) && (m_##p.Get() != nullptr))      \
      return HRESULT_FROM_WIN32(GetLastError());              \
    return S_OK;                                              \
  }                                                           \
  HRESULT STDMETHODCALLTYPE put_##p(LPCWSTR value) override { \
    LPCWSTR result = m_##p.Set(value);                        \
    if ((result == nullptr) && (value != nullptr))            \
      return HRESULT_FROM_WIN32(GetLastError());              \
    return S_OK;                                              \
  }                                                           \
                                                              \
 protected:                                                   \
  AutoCoMemString m_##p;

#define COREWEBVIEW2ENVIRONMENTOPTIONS_BOOL_PROPERTY(p)     \
 public:                                                    \
  HRESULT STDMETHODCALLTYPE get_##p(BOOL* value) override { \
    if (!value)                                             \
      return E_POINTER;                                     \
    *value = m_##p;                                         \
    return S_OK;                                            \
  }                                                         \
  HRESULT STDMETHODCALLTYPE put_##p(BOOL value) override {  \
    m_##p = value;                                          \
    return S_OK;                                            \
  }                                                         \
                                                            \
 protected:                                                 \
  BOOL m_##p = FALSE;

// This is a base COM class that implements ICoreWebView2EnvironmentOptions.
template <typename allocate_fn_t,
          allocate_fn_t allocate_fn,
          typename deallocate_fn_t,
          deallocate_fn_t deallocate_fn>
class CoreWebView2EnvironmentOptionsBase
    : public Microsoft::WRL::Implements<
          Microsoft::WRL::RuntimeClassFlags<Microsoft::WRL::ClassicCom>,
          ICoreWebView2EnvironmentOptions> {
 public:
  CoreWebView2EnvironmentOptionsBase() {
    // Initialize the target compatible browser version value to the version of
    // the browser binaries corresponding to this version of the SDK.
    m_TargetCompatibleBrowserVersion.Set(CORE_WEBVIEW_TARGET_PRODUCT_VERSION);
  }

 protected:
  ~CoreWebView2EnvironmentOptionsBase(){};

  class AutoCoMemString {
   public:
    AutoCoMemString() {}
    ~AutoCoMemString() { Release(); }
    void Release() {
      if (m_string) {
        deallocate_fn(m_string);
        m_string = nullptr;
      }
    }

    LPCWSTR Set(LPCWSTR str) {
      Release();
      if (str) {
        m_string = MakeCoMemString(str);
      }
      return m_string;
    }
    LPCWSTR Get() { return m_string; }
    LPWSTR Copy() {
      if (m_string)
        return MakeCoMemString(m_string);
      return nullptr;
    }

   protected:
    LPWSTR MakeCoMemString(LPCWSTR source) {
      const size_t length = wcslen(source);
      const size_t bytes = (length + 1) * sizeof(*source);
      // Ensure we didn't overflow during our size calculation.
      if (bytes <= length) {
        return nullptr;
      }

      wchar_t* result = reinterpret_cast<wchar_t*>(allocate_fn(bytes));
      if (result)
        memcpy(result, source, bytes);

      return result;
    }

    LPWSTR m_string = nullptr;
  };

  COREWEBVIEW2ENVIRONMENTOPTIONS_STRING_PROPERTY(AdditionalBrowserArguments)
  COREWEBVIEW2ENVIRONMENTOPTIONS_STRING_PROPERTY(Language)
  COREWEBVIEW2ENVIRONMENTOPTIONS_STRING_PROPERTY(TargetCompatibleBrowserVersion)
  COREWEBVIEW2ENVIRONMENTOPTIONS_BOOL_PROPERTY(
      AllowSingleSignOnUsingOSPrimaryAccount)
};

template <typename allocate_fn_t,
          allocate_fn_t allocate_fn,
          typename deallocate_fn_t,
          deallocate_fn_t deallocate_fn>
class CoreWebView2EnvironmentOptionsBaseClass
    : public Microsoft::WRL::RuntimeClass<
          Microsoft::WRL::RuntimeClassFlags<Microsoft::WRL::ClassicCom>,
          CoreWebView2EnvironmentOptionsBase<allocate_fn_t,
                                             allocate_fn,
                                             deallocate_fn_t,
                                             deallocate_fn>> {
 public:
  CoreWebView2EnvironmentOptionsBaseClass() {}

 protected:
  ~CoreWebView2EnvironmentOptionsBaseClass() override{};
};

typedef CoreWebView2EnvironmentOptionsBaseClass<decltype(&::CoTaskMemAlloc),
                                                ::CoTaskMemAlloc,
                                                decltype(&::CoTaskMemFree),
                                                ::CoTaskMemFree>
    CoreWebView2EnvironmentOptions;

#endif  // __core_webview2_environment_options_h__
