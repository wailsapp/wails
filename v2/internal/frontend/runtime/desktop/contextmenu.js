/*
--default-contextmenu: auto; (default) will show the default context menu if contentEditable is true OR text has been selected OR element is input or textarea
--default-contextmenu: show; will always show the default context menu
--default-contextmenu: hide; will always hide the default context menu

This rule is inherited like normal CSS rules, so nesting works as expected
*/
export function processDefaultContextMenu(event) {
    // Process default context menu
    const element = event.target;
    const computedStyle = window.getComputedStyle(element);
    const defaultContextMenuAction = computedStyle.getPropertyValue("--default-contextmenu").trim();
    switch (defaultContextMenuAction) {
        case "show":
            return;
        case "hide":
            event.preventDefault();
            return;
        default:
            // Check if contentEditable is true
            if (element.isContentEditable) {
                return;
            }

            // Check if text has been selected and action is on the selected elements
            const selection = window.getSelection();
            const hasSelection = (selection.toString().length > 0)
            if (hasSelection) {
                for (let i = 0; i < selection.rangeCount; i++) {
                    const range = selection.getRangeAt(i);
                    const rects = range.getClientRects();
                    for (let j = 0; j < rects.length; j++) {
                        const rect = rects[j];
                        if (document.elementFromPoint(rect.left, rect.top) === element) {
                            return;
                        }
                    }
                }
            }
            // Check if tagname is input or textarea
            if (element.tagName === "INPUT" || element.tagName === "TEXTAREA") {
                if (hasSelection || (!element.readOnly && !element.disabled)) {
                    return;
                }
            }

            // hide default context menu
            event.preventDefault();
    }
}
