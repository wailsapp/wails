/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import { newRuntimeCaller, objectNames } from "./runtime.js";

const call = newRuntimeCaller(objectNames.Android);

// Method IDs (must match messageprocessor_android.go)
const HapticsVibrate = 0;
const DeviceInfo = 1;
const ToastShow = 2;

export namespace Haptics {
    /** Vibrate the device for the given duration in milliseconds. */
    export function Vibrate(durationMs: number = 100): Promise<void> {
        return call(HapticsVibrate, { duration: durationMs });
    }
}

export namespace Device {
    export interface Info {
        platform: string;
        manufacturer: string;
        brand: string;
        model: string;
        device: string;
        version: string;
        sdkInt: number;
    }
    /** Return information about the Android device. */
    export function Info(): Promise<Info> {
        return call(DeviceInfo);
    }
}

export namespace Toast {
    /** Show a short Android toast message. */
    export function Show(message: string): Promise<void> {
        return call(ToastShow, { message });
    }
}
