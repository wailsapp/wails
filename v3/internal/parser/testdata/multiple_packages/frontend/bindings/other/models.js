// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import * as uuid from "../github.com-google-uuid/models.js";

export class Person {
    /**
     * Creates a new Person instance.
     * @param {Partial<Person>} [source = {}] - The source object to create the Person.
     */
    constructor(source = {}) {
        if (!("UUID" in source)) {
            /**
             * @member
             * @type {uuid.UUID}
             */
            this["UUID"] = "";
        }
        if (!("Name" in source)) {
            /**
             * @member
             * @type {string}
             */
            this["Name"] = "";
        }

        Object.assign(this, source);
    }

    /**
     * Creates a new Person instance from a string or object.
     * @param {string|object} source - The source data to create a Person instance from.
     * @returns {Person} A new Person instance.
     */
    static createFrom(source) {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Person(parsedSource);
    }
};
