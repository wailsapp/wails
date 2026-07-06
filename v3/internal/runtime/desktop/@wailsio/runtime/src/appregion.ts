/*
 _     __    _ __
| |   / /___(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { invoke } from "./system.js";
import { whenReady } from "./utils.js";
import { hasDOM } from "./environment.js";

type NonClientRegionKind = "caption" | "minimize" | "maximize" | "close";

interface NonClientRegion {
    kind: NonClientRegionKind;
    left: number;
    top: number;
    right: number;
    bottom: number;
}

/*
--wails-non-client-region: caption;  marks an area that can drag the window
--wails-non-client-region: minimize; marks a custom minimize button
--wails-non-client-region: maximize; marks a custom maximize button
--wails-non-client-region: close;    marks a custom close button
*/
const regionProperty = "--wails-non-client-region";
const runtimeConfigReadyEvent = "wails:runtime-config-ready";
const validRegions = new Set<NonClientRegionKind>(["caption", "minimize", "maximize", "close"]);

// Setup
if (hasDOM) {
    window._wails = window._wails || {};
}

let updatePending = false;
let lastPayload = "";
let observedElements = new Set<Element>();
let resizeObserver: ResizeObserver | undefined;
let trackingStarted = false;

function normaliseRegionKind(value: string): NonClientRegionKind | undefined {
    const region = value.trim().toLowerCase();
    if (validRegions.has(region as NonClientRegionKind)) {
        return region as NonClientRegionKind;
    }
    return undefined;
}

function nonClientRegionForElement(element: Element): NonClientRegionKind | undefined {
    if (!(element instanceof HTMLElement)) {
        return undefined;
    }

    const style = window.getComputedStyle(element);
    const region = normaliseRegionKind(style.getPropertyValue(regionProperty));
    if (!region) {
        return undefined;
    }

    const parent = element.parentElement;
    if (parent) {
        const parentStyle = window.getComputedStyle(parent);
        // The CSS property is inherited. Only report the outermost element for
        // each contiguous region so native hit testing sees stable rectangles.
        if (normaliseRegionKind(parentStyle.getPropertyValue(regionProperty)) === region) {
            return undefined;
        }
    }

    return region;
}

function isVisible(element: HTMLElement): boolean {
    const style = window.getComputedStyle(element);
    return style.display !== "none" &&
        style.visibility !== "hidden" &&
        style.contentVisibility !== "hidden";
}

function elementRegion(element: Element): NonClientRegion | undefined {
    if (!(element instanceof HTMLElement)) {
        return undefined;
    }

    const kind = nonClientRegionForElement(element);
    if (!kind || !isVisible(element)) {
        return undefined;
    }

    const rect = element.getBoundingClientRect();
    if (rect.width <= 0 || rect.height <= 0) {
        return undefined;
    }

    // Native hit testing runs in physical pixels, while DOM geometry is in CSS pixels.
    const scale = window.devicePixelRatio || 1;
    const left = Math.floor(rect.left * scale);
    const top = Math.floor(rect.top * scale);
    const right = Math.ceil(rect.right * scale);
    const bottom = Math.ceil(rect.bottom * scale);

    if (right <= left || bottom <= top) {
        return undefined;
    }

    return { kind, left, top, right, bottom };
}

function regionElements(): Element[] {
    const elements: Element[] = [];

    if (document.documentElement) {
        elements.push(document.documentElement);
    }
    if (document.body) {
        elements.push(document.body);
        // Append via a loop: spreading a huge NodeList into push() overflows
        // the engine's argument limit on very large documents.
        for (const element of document.body.querySelectorAll("*")) {
            elements.push(element);
        }
    }

    return elements;
}

function observeRegionElements(elements: Element[]): void {
    if (typeof ResizeObserver === "undefined") {
        return;
    }

    // Track size changes only for active region elements. DOM structure and style
    // changes are covered by MutationObserver in startNonClientRegionTracking().
    resizeObserver ??= new ResizeObserver(scheduleUpdate);
    const nextElements = new Set(elements);

    for (const element of observedElements) {
        if (!nextElements.has(element)) {
            resizeObserver.unobserve(element);
        }
    }

    for (const element of nextElements) {
        if (!observedElements.has(element)) {
            resizeObserver.observe(element);
        }
    }

    observedElements = nextElements;
}

function updateNonClientRegions(): void {
    updatePending = false;

    const elements = regionElements();
    const regions: NonClientRegion[] = [];
    const activeElements: Element[] = [];

    for (const element of elements) {
        const region = elementRegion(element);
        if (region) {
            regions.push(region);
            activeElements.push(element);
        }
    }

    observeRegionElements(activeElements);

    const payload = JSON.stringify({ version: 1, regions });
    if (payload === lastPayload) {
        // Avoid sending duplicate native messages during resize or style churn.
        return;
    }

    lastPayload = payload;
    invoke("wails:non-client-region:" + payload);
}

function scheduleUpdate(): void {
    if (updatePending) {
        return;
    }

    // Batch region updates to animation frames so layout is measured once per frame.
    updatePending = true;
    window.requestAnimationFrame(updateNonClientRegions);
}

function startNonClientRegionTracking(): void {
    if (trackingStarted) {
        return;
    }

    trackingStarted = true;
    // Send an initial empty or populated region list once the DOM is ready.
    scheduleUpdate();

    const mutationObserver = new MutationObserver(scheduleUpdate);
    mutationObserver.observe(document.documentElement, {
        attributes: true,
        childList: true,
        subtree: true,
    });

    window.addEventListener("resize", scheduleUpdate);
    window.addEventListener("scroll", scheduleUpdate, true);
    window.visualViewport?.addEventListener("resize", scheduleUpdate);
    window.visualViewport?.addEventListener("scroll", scheduleUpdate);
}

function tryStartNonClientRegionTracking(): boolean {
    const os = window._wails.environment?.OS;
    if (os === undefined) {
        return false;
    }

    const enabled = window._wails.flags?.nonClientRegionTracking;
    if (os === "windows") {
        if (enabled === true) {
            whenReady(startNonClientRegionTracking);
        }
        return true;
    }

    return true;
}

if (hasDOM && !tryStartNonClientRegionTracking()) {
    window.addEventListener(runtimeConfigReadyEvent, tryStartNonClientRegionTracking, { once: true });
}
