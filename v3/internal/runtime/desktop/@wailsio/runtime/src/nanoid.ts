// Source: https://github.com/ai/nanoid

// The MIT License (MIT)
//
// Copyright 2017 Andrey Sitnik <andrey@sitnik.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
//     subject to the following conditions:
//
//     The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
//     THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// This alphabet uses `A-Za-z0-9_-` symbols.
// The order of characters is optimized for better gzip and brotli compression.
// References to the same file (works both for gzip and brotli):
// `'use`, `andom`, and `rict'`
// References to the brotli default dictionary:
// `-26T`, `1983`, `40px`, `75px`, `bush`, `jack`, `mind`, `very`, and `wolf`
const urlAlphabet =
    'useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict'

// MODIFIED FOR WAILS: the upstream implementation draws from Math.random(),
// which is not a cryptographic source. These ids are used for the runtime
// client id, binding call ids, chunk ids, and desktop stream session ids. They
// are correlation identifiers rather than authentication capabilities, but a
// stronger source makes accidental collisions negligible even across many
// concurrently loaded runtimes. The alphabet is 64 characters, so masking a
// random byte with 63 is uniform — no rejection sampling needed.
export function nanoid(size: number = 21): string {
    const source = typeof crypto !== "undefined" && typeof crypto.getRandomValues === "function"
        ? crypto
        : null;

    if (source) {
        const bytes = new Uint8Array(size);
        source.getRandomValues(bytes);
        let id = '';
        for (let i = 0; i < size; i++) {
            id += urlAlphabet[bytes[i] & 63];
        }
        return id;
    }

    // Fallback for environments without Web Crypto (very old engines, some
    // server-side rendering contexts). Ids remain unique but not unguessable.
    let id = '';
    let i = size | 0;
    while (i--) {
        id += urlAlphabet[(Math.random() * 64) | 0];
    }
    return id;
}
