// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import { Call as $Call, CancellablePromise as $CancellablePromise, Create as $Create } from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function Method(): $CancellablePromise<$models.Data> {
    return $Call.ByName("main.Service.Method").then(($result: any) => {
        return $$createType0($result);
    });
}

// Private type creation functions
const $$createType0 = $models.Data.createFrom;
