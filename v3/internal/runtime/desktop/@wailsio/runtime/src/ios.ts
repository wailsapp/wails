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
