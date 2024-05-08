// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Create as $Create} from "@wailsio/runtime";

/**
 * Person represents a person
 */
export class Person {
    "Title": Title;
    "Name": string;

    /** Creates a new Person instance. */
    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = "";
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Person instance from a string or object.
     */
    static createFrom($$source: any = {}): Person {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new Person($$parsedSource as Partial<Person>);
    }
}

/**
 * Title is a title
 */
export enum Title {
    /**
     * Mister is a title
     */
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
    Mrs = "Mrs",
    Dr = "Dr",
};
