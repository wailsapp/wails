// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * Service represents the notifications service
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

/**
 * CheckNotificationAuthorization checks current notification permission status.
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function CheckNotificationAuthorization() {
    let $resultPromise = /** @type {any} */($Call.ByID(2789931702));
    return $resultPromise;
}

/**
 * RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
 * Registering a category with the same name as a previously registered NotificationCategory will override it.
 * @param {$models.NotificationCategory} category
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RegisterNotificationCategory(category) {
    let $resultPromise = /** @type {any} */($Call.ByID(2679064664, category));
    return $resultPromise;
}

/**
 * RemoveAllDeliveredNotifications removes all delivered notifications.
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RemoveAllDeliveredNotifications() {
    let $resultPromise = /** @type {any} */($Call.ByID(384520397));
    return $resultPromise;
}

/**
 * RemoveAllPendingNotifications removes all pending notifications.
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RemoveAllPendingNotifications() {
    let $resultPromise = /** @type {any} */($Call.ByID(1423986276));
    return $resultPromise;
}

/**
 * RemoveDeliveredNotification removes a delivered notification matching the unique identifier.
 * @param {string} identifier
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RemoveDeliveredNotification(identifier) {
    let $resultPromise = /** @type {any} */($Call.ByID(149440045, identifier));
    return $resultPromise;
}

/**
 * RemoveNotification is a macOS stub that always returns nil.
 * Use one of the following instead:
 * RemoveAllPendingNotifications
 * RemovePendingNotification
 * RemoveAllDeliveredNotifications
 * RemoveDeliveredNotification
 * (Linux-specific)
 * @param {string} identifier
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RemoveNotification(identifier) {
    let $resultPromise = /** @type {any} */($Call.ByID(3702062929, identifier));
    return $resultPromise;
}

/**
 * RemoveNotificationCategory remove a previously registered NotificationCategory.
 * @param {string} categoryId
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RemoveNotificationCategory(categoryId) {
    let $resultPromise = /** @type {any} */($Call.ByID(229511469, categoryId));
    return $resultPromise;
}

/**
 * RemovePendingNotification removes a pending notification matching the unique identifier.
 * @param {string} identifier
 * @returns {Promise<void> & { cancel(): void }}
 */
export function RemovePendingNotification(identifier) {
    let $resultPromise = /** @type {any} */($Call.ByID(3872412470, identifier));
    return $resultPromise;
}

/**
 * RequestNotificationAuthorization requests permission for notifications.
 * @returns {Promise<boolean> & { cancel(): void }}
 */
export function RequestNotificationAuthorization() {
    let $resultPromise = /** @type {any} */($Call.ByID(729898933));
    return $resultPromise;
}

/**
 * SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
 * @param {$models.NotificationOptions} options
 * @returns {Promise<void> & { cancel(): void }}
 */
export function SendNotification(options) {
    let $resultPromise = /** @type {any} */($Call.ByID(2246903123, options));
    return $resultPromise;
}

/**
 * SendNotificationWithActions sends a notification with additional actions and inputs.
 * A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
 * If a NotificationCategory is not registered a basic notification will be sent.
 * @param {$models.NotificationOptions} options
 * @returns {Promise<void> & { cancel(): void }}
 */
export function SendNotificationWithActions(options) {
    let $resultPromise = /** @type {any} */($Call.ByID(1615199806, options));
    return $resultPromise;
}
