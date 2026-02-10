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

// Method IDs
const HapticsVibrate = 0;
const DeviceInfo = 1;
const ToastShow = 2;
const ScrollSetEnabled = 3;
const ScrollSetBounceEnabled = 4;
const ScrollSetIndicatorsEnabled = 5;
const NavigationSetBackForwardGestures = 6;
const LinksSetPreviewEnabled = 7;
const UserAgentSet = 8;

export namespace Haptics {
    export function Vibrate(duration: number = 100): Promise<void> {
        return call(HapticsVibrate, { duration });
    }
}

export namespace Device {
    export interface Info {
        platform: string;
        model: string;
        version: string;
    }

    export function Info(): Promise<Info> {
        return call(DeviceInfo);
    }
}

export namespace Toast {
    export function Show(message: string): Promise<void> {
        return call(ToastShow, { message });
    }
}

export namespace Scroll {
    export function SetEnabled(enabled: boolean): Promise<void> {
        return call(ScrollSetEnabled, { enabled });
    }

    export function SetBounceEnabled(enabled: boolean): Promise<void> {
        return call(ScrollSetBounceEnabled, { enabled });
    }

    export function SetIndicatorsEnabled(enabled: boolean): Promise<void> {
        return call(ScrollSetIndicatorsEnabled, { enabled });
    }
}

export namespace Navigation {
    export function SetBackForwardGesturesEnabled(enabled: boolean): Promise<void> {
        return call(NavigationSetBackForwardGestures, { enabled });
    }
}

export namespace Links {
    export function SetPreviewEnabled(enabled: boolean): Promise<void> {
        return call(LinksSetPreviewEnabled, { enabled });
    }
}

export namespace UserAgent {
    export function Set(ua: string): Promise<void> {
        return call(UserAgentSet, { ua });
    }
}
