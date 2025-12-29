/*
 * Wails Linux Media Interceptor
 *
 * On Linux, WebKitGTK uses GStreamer for media playback. GStreamer does not
 * have a URI handler for the "wails://" protocol, causing video and audio
 * elements to fail loading.
 *
 * This script intercepts media elements and converts wails:// URLs to blob URLs
 * by fetching the content through fetch() (which works via WebKit's URI scheme
 * handler) and creating object URLs.
 *
 * See: https://github.com/wailsapp/wails/issues/4412
 * See: https://bugs.webkit.org/show_bug.cgi?id=146351
 */
(function() {
    'use strict';

    // This constant is replaced by the server based on EnableGStreamerCaching option
    const ENABLE_CACHING = true;

    const blobUrlCache = ENABLE_CACHING ? new Map() : null;
    const processingElements = new WeakSet();
    const processingSourceElements = new WeakSet();

    function shouldInterceptUrl(src) {
        if (!src || src.startsWith('blob:') || src.startsWith('data:')) {
            return false;
        }
        if (src.startsWith('wails://')) {
            return true;
        }
        if (src.startsWith('/') || (!src.includes('://') && !src.startsWith('//'))) {
            return true;
        }
        if (src.startsWith(window.location.origin)) {
            return true;
        }
        return false;
    }

    function toAbsoluteUrl(src) {
        if (src.startsWith('wails://') || src.startsWith('http://') || src.startsWith('https://')) {
            return src;
        }
        if (src.startsWith('/')) {
            return window.location.origin + src;
        }
        const base = window.location.href.substring(0, window.location.href.lastIndexOf('/') + 1);
        return base + src;
    }

    async function convertToBlob(url) {
        const absoluteUrl = toAbsoluteUrl(url);
        if (ENABLE_CACHING && blobUrlCache.has(absoluteUrl)) {
            return blobUrlCache.get(absoluteUrl);
        }
        const response = await fetch(absoluteUrl);
        if (!response.ok) {
            throw new Error('Failed to fetch media: ' + response.status + ' ' + response.statusText);
        }
        const blob = await response.blob();
        const blobUrl = URL.createObjectURL(blob);
        if (ENABLE_CACHING) {
            blobUrlCache.set(absoluteUrl, blobUrl);
        }
        return blobUrl;
    }

    async function processSourceElement(source) {
        if (processingSourceElements.has(source)) {
            return;
        }
        const src = source.src || source.getAttribute('src');
        if (!src || !shouldInterceptUrl(src)) {
            return;
        }
        processingSourceElements.add(source);
        try {
            source.dataset.wailsOriginalSrc = src;
            const blobUrl = await convertToBlob(src);
            source.src = blobUrl;
            console.debug('[Wails] Converted source element:', src);
        } catch (err) {
            console.error('[Wails] Failed to convert source element:', src, err);
        }
    }

    async function processMediaElement(element) {
        if (processingElements.has(element)) {
            return;
        }
        const src = element.src || element.getAttribute('src');
        if (src && shouldInterceptUrl(src)) {
            processingElements.add(element);
            try {
                element.dataset.wailsOriginalSrc = src;
                const blobUrl = await convertToBlob(src);
                element.src = blobUrl;
                console.debug('[Wails] Converted media element:', src);
            } catch (err) {
                console.error('[Wails] Failed to convert media element:', src, err);
            }
        }

        const sources = element.querySelectorAll('source');
        for (const source of sources) {
            await processSourceElement(source);
        }
        if (sources.length > 0 && element.dataset.wailsOriginalSrc === undefined) {
            element.load();
        }
    }

    function scanForMediaElements() {
        const mediaElements = document.querySelectorAll('video, audio');
        mediaElements.forEach(function(element) {
            processMediaElement(element);
        });
    }

    function setupMutationObserver() {
        const observer = new MutationObserver(function(mutations) {
            for (const mutation of mutations) {
                for (const node of mutation.addedNodes) {
                    if (node instanceof HTMLMediaElement) {
                        processMediaElement(node);
                    } else if (node instanceof Element) {
                        const mediaElements = node.querySelectorAll('video, audio');
                        mediaElements.forEach(function(el) {
                            processMediaElement(el);
                        });
                    }
                }
                if (mutation.type === 'attributes' &&
                    mutation.attributeName === 'src' &&
                    mutation.target instanceof HTMLMediaElement) {
                    processingElements.delete(mutation.target);
                    processMediaElement(mutation.target);
                }
                if (mutation.type === 'attributes' &&
                    mutation.attributeName === 'src' &&
                    mutation.target instanceof HTMLSourceElement) {
                    processingSourceElements.delete(mutation.target);
                    processSourceElement(mutation.target);
                }
            }
        });

        observer.observe(document.documentElement, {
            childList: true,
            subtree: true,
            attributes: true,
            attributeFilter: ['src']
        });
    }

    console.debug('[Wails] Enabling media interceptor for Linux/WebKitGTK');

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', scanForMediaElements);
    } else {
        scanForMediaElements();
    }

    setupMutationObserver();
})();
