// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import { Create as $Create } from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as nobindingshere$0 from "../no_bindings_here/models.js";

/**
 * SomeClass renders as a TS class.
 */
export class SomeClass {
    /**
     * Creates a new SomeClass instance.
     * @param {Partial<SomeClass>} [$$source = {}] - The source object to create the SomeClass.
     */
    constructor($$source = {}) {
        if (!("Field" in $$source)) {
            /**
             * @member
             * @type {string}
             */
            this["Field"] = "";
        }
        if (!("Meadow" in $$source)) {
            /**
             * @member
             * @type {nobindingshere$0.HowDifferent<number>}
             */
            this["Meadow"] = (new nobindingshere$0.HowDifferent());
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new SomeClass instance from a string or object.
     * @param {any} [$$source = {}]
     * @returns {SomeClass}
     */
    static createFrom($$source = {}) {
        const $$createField1_0 = $$createType0;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("Meadow" in $$parsedSource) {
            $$parsedSource["Meadow"] = $$createField1_0($$parsedSource["Meadow"]);
        }
        return new SomeClass(/** @type {Partial<SomeClass>} */($$parsedSource));
    }
}

// Private type creation functions
const $$createType0 = nobindingshere$0.HowDifferent.createFrom($Create.Any);
