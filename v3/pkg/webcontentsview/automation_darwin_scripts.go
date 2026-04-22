//go:build darwin

package webcontentsview

const darwinPageAutomationScript = `
(() => {
  if (window.__wailsAutomationPageInstrumentationInstalled) {
    return;
  }

  window.__wailsAutomationPageInstrumentationInstalled = true;

  const post = (type, payload) => {
    try {
      window.webkit.messageHandlers.wailsAutomationPage.postMessage({ type, payload });
    } catch (_) {
    }
  };

  const encode = (value) => {
    if (value instanceof Error) {
      return value.message || value.toString();
    }
    if (typeof value === 'string') {
      return value;
    }
    try {
      return JSON.stringify(value);
    } catch (_) {
      return String(value);
    }
  };

  for (const level of ['log', 'info', 'warn', 'error', 'debug']) {
    const original = console[level] && console[level].bind(console);
    if (!original) {
      continue;
    }

    console[level] = (...args) => {
      post('console', {
        level,
        text: args.map(encode).join(' '),
        args: args.map(encode),
        timestamp: Date.now(),
      });
      return original(...args);
    };
  }

  window.addEventListener('error', (event) => {
    post('exception', {
      message: event.message || '',
      stack: event.error && event.error.stack ? event.error.stack : '',
      source: event.filename || '',
      line: event.lineno || 0,
      column: event.colno || 0,
      timestamp: Date.now(),
    });
  });

  window.addEventListener('unhandledrejection', (event) => {
    const reason = event.reason;
    post('exception', {
      message: reason && reason.message ? reason.message : encode(reason),
      stack: reason && reason.stack ? reason.stack : '',
      timestamp: Date.now(),
      unhandledRejection: true,
    });
  });

  document.addEventListener('DOMContentLoaded', () => {
    post('domcontentloaded', {
      url: location.href,
      title: document.title || '',
      timestamp: Date.now(),
    });
  }, { once: true });

  window.addEventListener('load', () => {
    post('load', {
      url: location.href,
      title: document.title || '',
      timestamp: Date.now(),
    });
  }, { once: true });

  let requestSequence = 0;
  const capturedResponseBodies = new Map();
  const extraHTTPHeaders = {};
  const nextRequestID = () => 'net-' + Date.now() + '-' + (++requestSequence);
  const trimCapturedBodies = () => {
    while (capturedResponseBodies.size > 50) {
      const oldestKey = capturedResponseBodies.keys().next().value;
      capturedResponseBodies.delete(oldestKey);
    }
  };
  const headerObject = (headers) => {
    const result = {};
    if (!headers) {
      return result;
    }

    if (headers instanceof Headers) {
      headers.forEach((value, key) => {
        result[key] = value;
      });
      return result;
    }

    if (Array.isArray(headers)) {
      for (const [key, value] of headers) {
        result[String(key)] = String(value);
      }
      return result;
    }

    if (typeof headers === 'object') {
      for (const [key, value] of Object.entries(headers)) {
        result[String(key)] = String(value);
      }
    }
    return result;
  };

  const mergeExtraHeaders = (headers) => {
    const merged = { ...extraHTTPHeaders };
    for (const [key, value] of Object.entries(headerObject(headers))) {
      merged[key] = value;
    }
    return merged;
  };

  if (typeof window.fetch === 'function') {
    const originalFetch = window.fetch.bind(window);
    window.fetch = async (...args) => {
      const init = args[0] instanceof Request ? {
        method: args[0].method,
        headers: mergeExtraHeaders(args[0].headers),
        body: args[1] && args[1].body !== undefined ? args[1].body : undefined,
        credentials: args[1] && args[1].credentials ? args[1].credentials : args[0].credentials,
        mode: args[1] && args[1].mode ? args[1].mode : args[0].mode,
        cache: args[1] && args[1].cache ? args[1].cache : args[0].cache,
        redirect: args[1] && args[1].redirect ? args[1].redirect : args[0].redirect,
        referrer: args[1] && args[1].referrer ? args[1].referrer : args[0].referrer,
        referrerPolicy: args[1] && args[1].referrerPolicy ? args[1].referrerPolicy : args[0].referrerPolicy,
      } : {
        ...(args[1] || {}),
        headers: mergeExtraHeaders(args[1] && args[1].headers),
      };
      const request = args[0] instanceof Request ? new Request(args[0], init) : new Request(args[0], init);
      const requestId = nextRequestID();
      const startTime = Date.now();

      post('networkRequest', {
        requestId,
        url: request.url,
        method: request.method || 'GET',
        requestHeaders: headerObject(request.headers),
        type: 'fetch',
        startTime,
      });

      try {
        const response = await originalFetch(request);
        const endTime = Date.now();
        let responseBody = '';
        try {
          responseBody = await response.clone().text();
        } catch (_) {
        }
        capturedResponseBodies.set(requestId, {
          body: responseBody,
          base64Encoded: false,
          timestamp: endTime,
        });
        trimCapturedBodies();
        const payload = {
          requestId,
          url: response.url || request.url,
          method: request.method || 'GET',
          status: response.status,
          statusText: response.statusText || '',
          headers: headerObject(response.headers),
          requestHeaders: headerObject(request.headers),
          type: 'fetch',
          startTime,
          endTime,
          duration: endTime - startTime,
        };
        post('networkResponse', payload);
        post('networkFinished', payload);
        return response;
      } catch (error) {
        const endTime = Date.now();
        post('networkFailed', {
          requestId,
          url: request.url,
          method: request.method || 'GET',
          requestHeaders: headerObject(request.headers),
          errorText: error && error.message ? error.message : String(error),
          type: 'fetch',
          startTime,
          endTime,
          duration: endTime - startTime,
        });
        throw error;
      }
    };
  }

  if (typeof window.XMLHttpRequest === 'function') {
    const OriginalXMLHttpRequest = window.XMLHttpRequest;
    window.XMLHttpRequest = function WailsAutomationXMLHttpRequest() {
      const xhr = new OriginalXMLHttpRequest();
      const requestId = nextRequestID();
      const requestHeaders = {};
      let url = '';
      let method = 'GET';
      let startTime = 0;

      const originalOpen = xhr.open;
      xhr.open = function patchedOpen(nextMethod, nextURL, ...rest) {
        method = nextMethod ? String(nextMethod) : 'GET';
        url = nextURL ? String(nextURL) : '';
        return originalOpen.call(this, nextMethod, nextURL, ...rest);
      };

      const originalSetRequestHeader = xhr.setRequestHeader;
      xhr.setRequestHeader = function patchedSetRequestHeader(name, value) {
        requestHeaders[String(name)] = String(value);
        return originalSetRequestHeader.call(this, name, value);
      };

      xhr.addEventListener('loadstart', () => {
        startTime = Date.now();
        for (const [name, value] of Object.entries(extraHTTPHeaders)) {
          if (!(name in requestHeaders)) {
            requestHeaders[name] = value;
            try {
              originalSetRequestHeader.call(xhr, name, value);
            } catch (_) {
            }
          }
        }
        post('networkRequest', {
          requestId,
          url,
          method,
          requestHeaders,
          type: 'xhr',
          startTime,
        });
      });

      xhr.addEventListener('loadend', () => {
        const endTime = Date.now();
        const rawHeaders = xhr.getAllResponseHeaders ? xhr.getAllResponseHeaders() : '';
        const headers = {};
        for (const line of rawHeaders.split(/\r?\n/)) {
          if (!line) {
            continue;
          }
          const index = line.indexOf(':');
          if (index <= 0) {
            continue;
          }
          const key = line.slice(0, index).trim();
          const value = line.slice(index + 1).trim();
          headers[key] = value;
        }

        const payload = {
          requestId,
          url: xhr.responseURL || url,
          method,
          status: xhr.status,
          statusText: xhr.statusText || '',
          headers,
          requestHeaders,
          type: 'xhr',
          startTime,
          endTime,
          duration: startTime ? endTime - startTime : 0,
        };

        try {
          if (xhr.responseType === '' || xhr.responseType === 'text') {
            capturedResponseBodies.set(requestId, {
              body: xhr.responseText || '',
              base64Encoded: false,
              timestamp: endTime,
            });
            trimCapturedBodies();
          }
        } catch (_) {
        }

        if (xhr.status === 0) {
          payload.errorText = 'XMLHttpRequest failed';
          post('networkFailed', payload);
          return;
        }

        post('networkResponse', payload);
        post('networkFinished', payload);
      });

      xhr.addEventListener('error', () => {
        const endTime = Date.now();
        post('networkFailed', {
          requestId,
          url,
          method,
          requestHeaders,
          errorText: 'XMLHttpRequest error',
          type: 'xhr',
          startTime,
          endTime,
          duration: startTime ? endTime - startTime : 0,
        });
      });

      return xhr;
    };
    window.XMLHttpRequest.prototype = OriginalXMLHttpRequest.prototype;
  }
})();
`

const darwinRuntimeAutomationScript = `
(() => {
  if (globalThis.__wailsAutomation) {
    return;
  }

  const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms));
  const normalize = (value) => String(value ?? '').replace(/\s+/g, ' ').trim().toLowerCase();
  const stringValue = (value) => String(value ?? '');

  const attributesOf = (element) => {
    const result = {};
    if (!element || !element.attributes) {
      return result;
    }

    for (const attribute of Array.from(element.attributes)) {
      result[attribute.name] = attribute.value;
    }

    return result;
  };

  const rectOf = (element) => {
    if (!element || !element.getBoundingClientRect) {
      return null;
    }

    const rect = element.getBoundingClientRect();
    return {
      x: rect.x,
      y: rect.y,
      width: rect.width,
      height: rect.height,
      top: rect.top,
      left: rect.left,
      right: rect.right,
      bottom: rect.bottom,
    };
  };

  const textOf = (element) => {
    if (!element) {
      return '';
    }
    return element.innerText || element.textContent || '';
  };

  const implicitRole = (element) => {
    if (!element || !element.tagName) {
      return '';
    }

    const tagName = element.tagName.toLowerCase();
    switch (tagName) {
      case 'a':
        return element.href ? 'link' : '';
      case 'button':
        return 'button';
      case 'input': {
        const type = (element.getAttribute('type') || 'text').toLowerCase();
        if (type === 'checkbox') {
          return 'checkbox';
        }
        if (type === 'radio') {
          return 'radio';
        }
        if (type === 'range') {
          return 'slider';
        }
        if (type === 'button' || type === 'submit' || type === 'reset') {
          return 'button';
        }
        return 'textbox';
      }
      case 'select':
        return 'combobox';
      case 'textarea':
        return 'textbox';
      case 'img':
        return 'img';
      case 'table':
        return 'table';
      case 'th':
        return 'columnheader';
      case 'tr':
        return 'row';
      case 'td':
        return 'cell';
      case 'ul':
      case 'ol':
        return 'list';
      case 'li':
        return 'listitem';
      case 'nav':
        return 'navigation';
      case 'main':
        return 'main';
      case 'header':
        return 'banner';
      case 'footer':
        return 'contentinfo';
      case 'form':
        return 'form';
      default:
        return '';
    }
  };

  const accessibleName = (element) => {
    if (!element) {
      return '';
    }

    const ariaLabel = element.getAttribute && element.getAttribute('aria-label');
    if (ariaLabel) {
      return ariaLabel.trim();
    }

    const labelledBy = element.getAttribute && element.getAttribute('aria-labelledby');
    if (labelledBy) {
      const label = labelledBy.split(/\s+/)
        .map((id) => document.getElementById(id))
        .filter(Boolean)
        .map((node) => textOf(node).trim())
        .join(' ')
        .trim();
      if (label) {
        return label;
      }
    }

    if (element.labels && element.labels.length > 0) {
      return Array.from(element.labels)
        .map((label) => textOf(label).trim())
        .join(' ')
        .trim();
    }

    if (element.alt) {
      return stringValue(element.alt).trim();
    }
    if (element.title) {
      return stringValue(element.title).trim();
    }
    if (element.placeholder) {
      return stringValue(element.placeholder).trim();
    }
    if (element.value && element.tagName && element.tagName.toLowerCase() === 'input') {
      return stringValue(element.value).trim();
    }
    return textOf(element).trim();
  };

  const roleOf = (element) => {
    if (!element || !element.getAttribute) {
      return '';
    }
    return element.getAttribute('role') || implicitRole(element);
  };

  const describeElement = (element) => {
    if (!element) {
      return null;
    }

    return {
      tagName: element.tagName ? element.tagName.toLowerCase() : '',
      id: element.id || '',
      name: element.getAttribute ? (element.getAttribute('name') || '') : '',
      role: roleOf(element),
      text: textOf(element),
      value: 'value' in element ? element.value : null,
      accessibleName: accessibleName(element),
      disabled: !!element.disabled,
      checked: !!element.checked,
      attributes: attributesOf(element),
      rect: rectOf(element),
      outerHTML: element.outerHTML || '',
    };
  };

  const serializeValue = (value) => {
    if (value === undefined) {
      return { type: 'undefined' };
    }
    if (value === null) {
      return { type: 'object', subtype: 'null', value: null };
    }

    const type = typeof value;
    if (type === 'string' || type === 'number' || type === 'boolean') {
      return { type, value };
    }
    if (type === 'bigint') {
      return { type: 'bigint', description: value.toString() };
    }
    if (type === 'function') {
      return { type: 'function', description: value.name || 'anonymous' };
    }
    if (Array.isArray(value)) {
      try {
        return { type: 'object', subtype: 'array', value: JSON.parse(JSON.stringify(value)) };
      } catch (_) {
        return { type: 'object', subtype: 'array', description: 'Array(' + value.length + ')' };
      }
    }
    if (value instanceof Element) {
      return { type: 'object', subtype: 'node', value: describeElement(value), description: value.tagName ? value.tagName.toLowerCase() : 'element' };
    }

    try {
      return { type: 'object', value: JSON.parse(JSON.stringify(value)) };
    } catch (_) {
      return { type: 'object', description: Object.prototype.toString.call(value) };
    }
  };

  const allElements = () => Array.from(document.querySelectorAll('*'));

  const resolveLabelElement = (labelText) => {
    const wanted = normalize(labelText);
    if (!wanted) {
      return null;
    }

    for (const label of Array.from(document.querySelectorAll('label'))) {
      const labelName = normalize(textOf(label));
      if (labelName !== wanted && !labelName.includes(wanted)) {
        continue;
      }

      if (label.control) {
        return label.control;
      }

      const targetID = label.getAttribute('for');
      if (targetID) {
        const target = document.getElementById(targetID);
        if (target) {
          return target;
        }
      }

      const nestedControl = label.querySelector('input, textarea, select, button');
      if (nestedControl) {
        return nestedControl;
      }
    }

    return null;
  };

  const resolveElement = (params = {}) => {
    if (params.selector) {
      return document.querySelector(params.selector);
    }
    if (params.label) {
      return resolveLabelElement(params.label);
    }
    if (params.role) {
      const wantedRole = normalize(params.role);
      const wantedName = normalize(params.name || params.text || '');
      return allElements().find((element) => {
        const role = normalize(roleOf(element));
        if (role !== wantedRole) {
          return false;
        }
        if (!wantedName) {
          return true;
        }
        const name = normalize(accessibleName(element));
        return name === wantedName || name.includes(wantedName);
      }) || null;
    }
    if (params.text) {
      const wantedText = normalize(params.text);
      return allElements().find((element) => {
        const text = normalize(textOf(element));
        return text === wantedText || text.includes(wantedText);
      }) || null;
    }
    return null;
  };

  const resolveElements = (params = {}) => {
    if (params.selector) {
      return Array.from(document.querySelectorAll(params.selector));
    }
    if (params.role) {
      const found = resolveElement(params);
      return found ? [found] : [];
    }
    if (params.text) {
      const wantedText = normalize(params.text);
      return allElements().filter((element) => {
        const text = normalize(textOf(element));
        return text === wantedText || text.includes(wantedText);
      });
    }
    if (params.label) {
      const found = resolveLabelElement(params.label);
      return found ? [found] : [];
    }
    return [];
  };

  const getDocumentNode = () => ({
    url: location.href,
    title: document.title || '',
    node: describeElement(document.documentElement),
  });

  const getAttributes = (element) => ({
    attributes: attributesOf(element),
  });

  const boundingRect = (element) => ({
    rect: rectOf(element),
  });

  const focusElement = (element) => {
    if (!element) {
      return { focused: false };
    }
    element.focus();
    return { focused: document.activeElement === element, node: describeElement(element) };
  };

  const clickElement = (element) => {
    if (!element) {
      return { clicked: false };
    }
    element.scrollIntoView({ block: 'center', inline: 'center' });
    if (typeof element.click === 'function') {
      element.click();
      return { clicked: true, node: describeElement(element) };
    }
    return { clicked: false };
  };

  const fillElement = (element, value) => {
    if (!element) {
      return { changed: false };
    }
    element.focus();
    if ('value' in element) {
      element.value = value == null ? '' : String(value);
      element.dispatchEvent(new Event('input', { bubbles: true }));
      element.dispatchEvent(new Event('change', { bubbles: true }));
      return { changed: true, node: describeElement(element) };
    }
    return { changed: false };
  };

  const selectOption = (element, value) => {
    if (!(element instanceof HTMLSelectElement)) {
      return { selected: false };
    }

    const wanted = String(value ?? '');
    let matched = false;
    for (const option of Array.from(element.options)) {
      const shouldSelect = option.value === wanted || option.text === wanted;
      option.selected = shouldSelect;
      matched = matched || shouldSelect;
      if (shouldSelect && !element.multiple) {
        break;
      }
    }

    if (matched) {
      element.dispatchEvent(new Event('input', { bubbles: true }));
      element.dispatchEvent(new Event('change', { bubbles: true }));
    }

    return { selected: matched, node: describeElement(element) };
  };

  const scrollIntoView = (element) => {
    if (!element) {
      return { scrolled: false };
    }
    element.scrollIntoView({ block: 'center', inline: 'center', behavior: 'instant' });
    return { scrolled: true, node: describeElement(element) };
  };

  const waitForSelector = async (params = {}) => {
    const timeout = Number(params.timeout || 30000);
    const pollInterval = Number(params.pollInterval || 50);
    const deadline = Date.now() + timeout;

    while (Date.now() <= deadline) {
      const element = resolveElement(params);
      if (element) {
        return { timedOut: false, node: describeElement(element) };
      }
      await wait(pollInterval);
    }

    return { timedOut: true, node: null };
  };

  const waitForCondition = async (params = {}) => {
    const timeout = Number(params.timeout || 30000);
    const pollInterval = Number(params.pollInterval || 50);
    const expression = params.expression || 'false';
    const deadline = Date.now() + timeout;

    while (Date.now() <= deadline) {
      try {
        if (Boolean((0, eval)(expression))) {
          return { matched: true, timedOut: false };
        }
      } catch (_) {
      }
      await wait(pollInterval);
    }

    return { matched: false, timedOut: true };
  };

  const storageObject = (storage) => {
    const result = {};
    for (let index = 0; index < storage.length; index += 1) {
      const key = storage.key(index);
      if (key !== null) {
        result[key] = storage.getItem(key);
      }
    }
    return result;
  };

  const snapshotNode = (element, depth = 0) => {
    if (!element || depth > 5) {
      return null;
    }

    const children = [];
    for (const child of Array.from(element.children || [])) {
      const childSnapshot = snapshotNode(child, depth + 1);
      if (childSnapshot) {
        children.push(childSnapshot);
      }
    }

    const role = roleOf(element) || (element.tagName ? element.tagName.toLowerCase() : '');
    const name = accessibleName(element);
    if (!role && !name && children.length === 0 && depth > 0) {
      return null;
    }

    return {
      role,
      name,
      tagName: element.tagName ? element.tagName.toLowerCase() : '',
      children,
    };
  };

  const evaluateExpression = async (expression, awaitPromise) => {
    let value = (0, eval)(expression);
    if (awaitPromise) {
      value = await Promise.resolve(value);
    }
    return serializeValue(value);
  };

  const dispatch = async (method, params = {}) => {
    const element = resolveElement(params);
    switch (method) {
      case 'DOM.getDocument':
        return getDocumentNode();
      case 'DOM.querySelector':
      case 'DOM.queryByRole':
      case 'DOM.queryByText':
      case 'DOM.queryByLabel':
      case 'Accessibility.queryByRole':
      case 'Accessibility.queryByLabel':
        return { node: describeElement(element) };
      case 'DOM.querySelectorAll':
        return { nodes: resolveElements(params).map(describeElement) };
      case 'DOM.getOuterHTML':
        return { outerHTML: element ? element.outerHTML || '' : '' };
      case 'DOM.getInnerText':
        return { innerText: textOf(element) };
      case 'DOM.getAttributes':
        return getAttributes(element);
      case 'DOM.getBoundingClientRect':
        return boundingRect(element);
      case 'DOM.scrollIntoView':
        return scrollIntoView(element);
      case 'DOM.focus':
        return focusElement(element);
      case 'DOM.click':
        return clickElement(element);
      case 'DOM.fill':
        return fillElement(element, params.value);
      case 'DOM.selectOption':
        return selectOption(element, params.value);
      case 'DOM.waitForSelector':
        return waitForSelector(params);
      case 'DOM.waitForCondition':
        return waitForCondition(params);
      case 'Storage.getLocalStorage':
        return { items: storageObject(localStorage) };
      case 'Storage.setLocalStorageItem':
        localStorage.setItem(String(params.key || ''), params.value == null ? '' : String(params.value));
        return { items: storageObject(localStorage) };
      case 'Storage.removeLocalStorageItem':
        localStorage.removeItem(String(params.key || ''));
        return { items: storageObject(localStorage) };
      case 'Storage.getSessionStorage':
        return { items: storageObject(sessionStorage) };
      case 'Storage.setSessionStorageItem':
        sessionStorage.setItem(String(params.key || ''), params.value == null ? '' : String(params.value));
        return { items: storageObject(sessionStorage) };
      case 'Storage.removeSessionStorageItem':
        sessionStorage.removeItem(String(params.key || ''));
        return { items: storageObject(sessionStorage) };
      case 'Accessibility.getSnapshot':
        return { snapshot: snapshotNode(document.body || document.documentElement) };
      default:
        throw new Error('unsupported automation method: ' + method);
    }
  };

  globalThis.__wailsAutomation = {
    evaluate: evaluateExpression,
    dispatch,
    serializeValue,
  };

  const originalDispatch = globalThis.__wailsAutomation.dispatch;
  globalThis.__wailsAutomation.dispatch = async (method, params = {}) => {
    if (method === 'Network.getResponseBody') {
      const entry = capturedResponseBodies.get(String(params.requestId || ''));
      if (!entry) {
        return { body: '', base64Encoded: false, found: false };
      }
      return { body: entry.body || '', base64Encoded: !!entry.base64Encoded, found: true };
    }

    if (method === 'Network.setExtraHTTPHeaders') {
      for (const key of Object.keys(extraHTTPHeaders)) {
        delete extraHTTPHeaders[key];
      }
      const headers = params.headers || {};
      for (const [key, value] of Object.entries(headers)) {
        extraHTTPHeaders[String(key)] = String(value);
      }
      return { headers: { ...extraHTTPHeaders } };
    }

    return originalDispatch(method, params);
  };
})();
`
