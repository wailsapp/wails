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

type AppRegionKind = "caption" | "minimize" | "maximize" | "close";

interface AppRegion {
    kind: AppRegionKind;
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
const validRegions = new Set<AppRegionKind>(["caption", "minimize", "maximize", "close"]);

// Setup
window._wails = window._wails || {};

let updatePending = false;
let lastPayload = "";
let observedElements = new Set<Element>();
let resizeObserver: ResizeObserver | undefined;
let trackingStarted = false;

function normaliseRegionKind(value: string): AppRegionKind | undefined {
    const region = value.trim().toLowerCase();
    if (validRegions.has(region as AppRegionKind)) {
        return region as AppRegionKind;
    }
    return undefined;
}

function appRegionForElement(element: Element): AppRegionKind | undefined {
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

function elementRegion(element: Element): AppRegion | undefined {
    if (!(element instanceof HTMLElement)) {
        return undefined;
    }

    const kind = appRegionForElement(element);
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
        elements.push(...document.body.querySelectorAll("*"));
    }

    return elements;
}

function observeRegionElements(elements: Element[]): void {
    if (typeof ResizeObserver === "undefined") {
        return;
    }

    // Track size changes only for active region elements. DOM structure and style
    // changes are covered by MutationObserver in startAppRegionTracking().
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

function updateAppRegions(): void {
    updatePending = false;

    const elements = regionElements();
    const regions: AppRegion[] = [];
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
    window.requestAnimationFrame(updateAppRegions);
}

function startAppRegionTracking(): void {
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

let environmentPolls = 0;
function tryStartAppRegionTracking(): void {
    const os = window._wails.environment?.OS;
    if (os === "windows") {
        whenReady(startAppRegionTracking);
        return;
    }

    if (os === undefined && environmentPolls++ < 100) {
        // The runtime environment can arrive after this side-effect module loads.
        window.setTimeout(tryStartAppRegionTracking, 50);
    }
}

tryStartAppRegionTracking();
