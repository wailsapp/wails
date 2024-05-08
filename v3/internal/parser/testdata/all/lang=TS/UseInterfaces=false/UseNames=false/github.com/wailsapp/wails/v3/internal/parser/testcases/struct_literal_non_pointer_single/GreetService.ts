// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import * as $models from "./models.ts";

export function ArrayInt($in: number[]): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3862002418, $in);
    return $resultPromise as any;
}

export function BoolInBoolOut($in: boolean): Promise<boolean> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2424639793, $in);
    return $resultPromise as any;
}

export function Float32InFloat32Out($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3132595881, $in);
    return $resultPromise as any;
}

export function Float64InFloat64Out($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2182412247, $in);
    return $resultPromise as any;
}

/**
 * Greet someone
 */
export function Greet(name: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1411160069, name);
    return $resultPromise as any;
}

export function Int16InIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3306292566, $in);
    return $resultPromise as any;
}

export function Int16PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1754277916, $in);
    return $resultPromise as any;
}

export function Int32InIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1909469092, $in);
    return $resultPromise as any;
}

export function Int32PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4251088558, $in);
    return $resultPromise as any;
}

export function Int64InIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1343888303, $in);
    return $resultPromise as any;
}

export function Int64PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2205561041, $in);
    return $resultPromise as any;
}

export function Int8InIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(572240879, $in);
    return $resultPromise as any;
}

export function Int8PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2189402897, $in);
    return $resultPromise as any;
}

export function IntInIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(642881729, $in);
    return $resultPromise as any;
}

export function IntPointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1066151743, $in);
    return $resultPromise as any;
}

export function IntPointerInputNamedOutputs($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2718999663, $in);
    return $resultPromise as any;
}

export function MapIntInt($in: { [_: `${number}`]: number }): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2386486356, $in);
    return $resultPromise as any;
}

export function MapIntPointerInt($in: { [_: string]: number }): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(550413585, $in);
    return $resultPromise as any;
}

export function MapIntSliceInt($in: { [_: `${number}`]: number[] }): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2900172572, $in);
    return $resultPromise as any;
}

export function MapIntSliceIntInMapIntSliceIntOut($in: { [_: `${number}`]: number[] }): Promise<{ [_: `${number}`]: number[] }> & { cancel(): void } {
    let $resultPromise = $Call.ByID(881980169, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

export function NoInputsStringOut(): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1075577233);
    return $resultPromise as any;
}

export function PointerBoolInBoolOut($in: boolean | null): Promise<boolean | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3589606958, $in);
    return $resultPromise as any;
}

export function PointerFloat32InFloat32Out($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(224675106, $in);
    return $resultPromise as any;
}

export function PointerFloat64InFloat64Out($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2124953624, $in);
    return $resultPromise as any;
}

export function PointerMapIntInt($in: { [_: `${number}`]: number } | null): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3516977899, $in);
    return $resultPromise as any;
}

export function PointerStringInStringOut($in: string | null): Promise<string | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(229603958, $in);
    return $resultPromise as any;
}

export function StringArrayInputNamedOutput($in: string[]): Promise<string[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3678582682, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType2($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

export function StringArrayInputNamedOutputs($in: string[]): Promise<string[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(319259595, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType2($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

export function StringArrayInputStringArrayOut($in: string[]): Promise<string[]> & { cancel(): void } {
    let $resultPromise = $Call.ByID(383995060, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType2($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

export function StringArrayInputStringOut($in: string[]): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1091960237, $in);
    return $resultPromise as any;
}

export function StructInputStructOutput($in: $models.Person): Promise<$models.Person> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3835643147, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType3($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

export function StructPointerInputErrorOutput($in: $models.Person | null): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2447692557, $in);
    return $resultPromise as any;
}

export function StructPointerInputStructPointerOutput($in: $models.Person | null): Promise<$models.Person | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2943477349, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType4($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise as any;
}

export function UInt16InUIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3401034892, $in);
    return $resultPromise as any;
}

export function UInt16PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1236957573, $in);
    return $resultPromise as any;
}

export function UInt32InUIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1160383782, $in);
    return $resultPromise as any;
}

export function UInt32PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1739300671, $in);
    return $resultPromise as any;
}

export function UInt64InUIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(793803239, $in);
    return $resultPromise as any;
}

export function UInt64PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1403757716, $in);
    return $resultPromise as any;
}

export function UInt8InUIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2988345717, $in);
    return $resultPromise as any;
}

export function UInt8PointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(518250834, $in);
    return $resultPromise as any;
}

export function UIntInUIntOut($in: number): Promise<number> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2836661285, $in);
    return $resultPromise as any;
}

export function UIntPointerInAndOutput($in: number | null): Promise<number | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1367187362, $in);
    return $resultPromise as any;
}

// Private type creation functions
const $$createType0 = $Create.Array($Create.Any);
const $$createType1 = $Create.Map($Create.Any, $$createType0);
const $$createType2 = $Create.Array($Create.Any);
const $$createType3 = $models.Person.createFrom;
const $$createType4 = $Create.Nullable($$createType3);
