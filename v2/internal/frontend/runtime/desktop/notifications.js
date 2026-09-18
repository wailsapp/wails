/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 9 */

import {Call} from "./calls";

/**
 * Initialize the notification service for the application.
 * This must be called before sending any notifications.
 * On macOS, this also ensures the notification delegate is properly initialized.
 *
 * @export
 * @return {Promise<void>}
 */
export function InitializeNotifications() {
    return Call(":wails:InitializeNotifications");
}

/**
 * Clean up notification resources and release any held connections.
 * This should be called when shutting down the application to properly release resources
 * (primarily needed on Linux to close D-Bus connections).
 *
 * @export
 * @return {Promise<void>}
 */
export function CleanupNotifications() {
    return Call(":wails:CleanupNotifications");
}

/**
 * Check if notifications are available on the current platform.
 *
 * @export
 * @return {Promise<boolean>} True if notifications are available, false otherwise
 */
export function IsNotificationAvailable() {
    return Call(":wails:IsNotificationAvailable");
}

/**
 * Request notification authorization from the user.
 * On macOS, this prompts the user to allow notifications.
 * On other platforms, this always returns true.
 *
 * @export
 * @return {Promise<boolean>} True if authorization was granted, false otherwise
 */
export function RequestNotificationAuthorization() {
    return Call(":wails:RequestNotificationAuthorization");
}

/**
 * Check the current notification authorization status.
 * On macOS, this checks if the app has notification permissions.
 * On other platforms, this always returns true.
 *
 * @export
 * @return {Promise<boolean>} True if authorized, false otherwise
 */
export function CheckNotificationAuthorization() {
    return Call(":wails:CheckNotificationAuthorization");
}

/**
 * Send a basic notification with the given options.
 * The notification will display with the provided title, subtitle (if supported), and body text.
 *
 * @export
 * @param {Object} options - Notification options
 * @param {string} options.id - Unique identifier for the notification
 * @param {string} options.title - Notification title
 * @param {string} [options.subtitle] - Notification subtitle (macOS and Linux only)
 * @param {string} [options.body] - Notification body text
 * @param {string} [options.categoryId] - Category ID for action buttons (requires SendNotificationWithActions)
 * @param {Object<string, any>} [options.data] - Additional user data to attach to the notification
 * @return {Promise<void>}
 */
export function SendNotification(options) {
    return Call(":wails:SendNotification", [options]);
}

/**
 * Send a notification with action buttons.
 * A NotificationCategory must be registered first using RegisterNotificationCategory.
 * The options.categoryId must match a previously registered category ID.
 * If the category is not found, a basic notification will be sent instead.
 *
 * @export
 * @param {Object} options - Notification options
 * @param {string} options.id - Unique identifier for the notification
 * @param {string} options.title - Notification title
 * @param {string} [options.subtitle] - Notification subtitle (macOS and Linux only)
 * @param {string} [options.body] - Notification body text
 * @param {string} options.categoryId - Category ID that matches a registered category
 * @param {Object<string, any>} [options.data] - Additional user data to attach to the notification
 * @return {Promise<void>}
 */
export function SendNotificationWithActions(options) {
    return Call(":wails:SendNotificationWithActions", [options]);
}

/**
 * Register a notification category that can be used with SendNotificationWithActions.
 * Categories define the action buttons and optional reply fields that will appear on notifications.
 * Registering a category with the same ID as a previously registered category will override it.
 *
 * @export
 * @param {Object} category - Notification category definition
 * @param {string} category.id - Unique identifier for the category
 * @param {Array<Object>} [category.actions] - Array of action buttons
 * @param {string} category.actions[].id - Unique identifier for the action
 * @param {string} category.actions[].title - Display title for the action button
 * @param {boolean} [category.actions[].destructive] - Whether the action is destructive (macOS-specific)
 * @param {boolean} [category.hasReplyField] - Whether to include a text input field for replies
 * @param {string} [category.replyPlaceholder] - Placeholder text for the reply field (required if hasReplyField is true)
 * @param {string} [category.replyButtonTitle] - Title for the reply button (required if hasReplyField is true)
 * @return {Promise<void>}
 */
export function RegisterNotificationCategory(category) {
    return Call(":wails:RegisterNotificationCategory", [category]);
}

/**
 * Remove a previously registered notification category.
 *
 * @export
 * @param {string} categoryId - The ID of the category to remove
 * @return {Promise<void>}
 */
export function RemoveNotificationCategory(categoryId) {
    return Call(":wails:RemoveNotificationCategory", [categoryId]);
}

/**
 * Remove all pending notifications from the notification center.
 * On Windows, this is a no-op as the platform manages notification lifecycle automatically.
 *
 * @export
 * @return {Promise<void>}
 */
export function RemoveAllPendingNotifications() {
    return Call(":wails:RemoveAllPendingNotifications");
}

/**
 * Remove a specific pending notification by its identifier.
 * On Windows, this is a no-op as the platform manages notification lifecycle automatically.
 *
 * @export
 * @param {string} identifier - The ID of the notification to remove
 * @return {Promise<void>}
 */
export function RemovePendingNotification(identifier) {
    return Call(":wails:RemovePendingNotification", [identifier]);
}

/**
 * Remove all delivered notifications from the notification center.
 * On Windows, this is a no-op as the platform manages notification lifecycle automatically.
 *
 * @export
 * @return {Promise<void>}
 */
export function RemoveAllDeliveredNotifications() {
    return Call(":wails:RemoveAllDeliveredNotifications");
}

/**
 * Remove a specific delivered notification by its identifier.
 * On Windows, this is a no-op as the platform manages notification lifecycle automatically.
 *
 * @export
 * @param {string} identifier - The ID of the notification to remove
 * @return {Promise<void>}
 */
export function RemoveDeliveredNotification(identifier) {
    return Call(":wails:RemoveDeliveredNotification", [identifier]);
}

/**
 * Remove a notification by its identifier.
 * This is a convenience function that works across platforms.
 * On macOS, use the more specific RemovePendingNotification or RemoveDeliveredNotification functions.
 *
 * @export
 * @param {string} identifier - The ID of the notification to remove
 * @return {Promise<void>}
 */
export function RemoveNotification(identifier) {
    return Call(":wails:RemoveNotification", [identifier]);
}

