// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * Service represents the notifications service
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import { Call as $Call, CancellablePromise as $CancellablePromise, Create as $Create } from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function CheckNotificationAuthorization(): $CancellablePromise<boolean> {
    return $Call.ByID(2789931702);
}

export function RegisterNotificationCategory(category: $models.NotificationCategory): $CancellablePromise<void> {
    return $Call.ByID(2679064664, category);
}

export function RemoveAllDeliveredNotifications(): $CancellablePromise<void> {
    return $Call.ByID(384520397);
}

export function RemoveAllPendingNotifications(): $CancellablePromise<void> {
    return $Call.ByID(1423986276);
}

export function RemoveDeliveredNotification(identifier: string): $CancellablePromise<void> {
    return $Call.ByID(149440045, identifier);
}

export function RemoveNotification(identifier: string): $CancellablePromise<void> {
    return $Call.ByID(3702062929, identifier);
}

export function RemoveNotificationCategory(categoryID: string): $CancellablePromise<void> {
    return $Call.ByID(229511469, categoryID);
}

export function RemovePendingNotification(identifier: string): $CancellablePromise<void> {
    return $Call.ByID(3872412470, identifier);
}

/**
 * Public methods that delegate to the implementation.
 */
export function RequestNotificationAuthorization(): $CancellablePromise<boolean> {
    return $Call.ByID(729898933);
}

export function SendNotification(options: $models.NotificationOptions): $CancellablePromise<void> {
    return $Call.ByID(2246903123, options);
}

export function SendNotificationWithActions(options: $models.NotificationOptions): $CancellablePromise<void> {
    return $Call.ByID(1615199806, options);
}
