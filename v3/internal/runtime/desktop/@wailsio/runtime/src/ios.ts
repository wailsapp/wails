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

const call = newRuntimeCaller(objectNames.IOS);

// Method IDs
const HapticsImpact = 0;
const DeviceInfo = 1;
const NativeTabsSetEnabled = 9;
const NativeTabsIsEnabled = 10;
const NativeTabsSelect = 11;

export namespace Haptics {
    export type ImpactStyle = "light"|"medium"|"heavy"|"soft"|"rigid";
    export function Impact(style: ImpactStyle = "medium"): Promise<void> {
        return call(HapticsImpact, { style });
    }
}

export namespace Device {
    export interface Info {
        model: string;
        systemName: string;
        systemVersion: string;
        isSimulator: boolean;
    }
    export function Info(): Promise<Info> {
        return call(DeviceInfo);
    }
}

/** Controls the native iOS tab bar configured in the application options. */
export namespace NativeTabs {
    /** Shows or hides the native tab bar. */
    export function SetEnabled(enabled: boolean): Promise<void> {
        return call(NativeTabsSetEnabled, { enabled });
    }

    /** Reports whether the native tab bar is currently enabled. */
    export function IsEnabled(): Promise<boolean> {
        return call(NativeTabsIsEnabled);
    }

    /** Selects the tab at the zero-based index. */
    export function Select(index: number): Promise<void> {
        return call(NativeTabsSelect, { index });
    }
}
