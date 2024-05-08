// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Create as $Create} from "/wails/runtime.js";

/**
 * A generic struct
 * @template T
 */
export class Person {
    /**
     * Creates a new Person instance.
     * @param {Partial<Person<T>>} [$$source = {}] - The source object to create the Person.
     */
    constructor($$source = {}) {
        if (!("Name" in $$source)) {
            /**
             * @member
             * @type {T | null}
             */
            this["Name"] = null;
        }
        if (!("AliasedField" in $$source)) {
            /**
             * @member
             * @type {number}
             */
            this["AliasedField"] = 0;
        }

        Object.assign(this, $$source);
    }

    /**
     * Given creation functions for each type parameter,
     * returns a creation function for a concrete instance
     * of the generic class Person.
     * @template T
     * @param {(any) => T} $$createParamT
     * @returns {($$source?: any) => Person<T>}
     */
    static createFrom($$createParamT) {
        const $$createField0_0 = $$createParamT;
        return ($$source = {}) => {
            let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
            if ("Name" in $$parsedSource) {
                $$parsedSource["Name"] = $$createField0_0($$parsedSource["Name"]);
            }
            return new Person(/** @type {Partial<Person<T>>} */($$parsedSource));
        };
    }
}