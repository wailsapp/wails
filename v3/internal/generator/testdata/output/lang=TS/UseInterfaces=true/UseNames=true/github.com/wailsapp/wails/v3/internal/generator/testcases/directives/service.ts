// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as otherpackage$0 from "./otherpackage/models.js";

function InternalMethod($0: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.Service.InternalMethod", $0) as any;
    return $resultPromise;
}

export function VisibleMethod($0: otherpackage$0.Dummy): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByName("main.Service.VisibleMethod", $0) as any;
    return $resultPromise;
}

export async function CustomMethod(arg: string): Promise<void> {
    await InternalMethod("Hello " + arg + "!");
}
