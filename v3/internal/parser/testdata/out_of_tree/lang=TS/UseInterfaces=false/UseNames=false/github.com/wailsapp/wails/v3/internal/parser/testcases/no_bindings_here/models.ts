// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Create as $Create} from "@wailsio/runtime";

import * as other$0 from "./other/models.ts";

/**
 * HowDifferent is a curious kind of person
 * that lets other people decide how they are different.
 */
export const HowDifferent = other$0.OtherPerson;

/**
 * HowDifferent is a curious kind of person
 * that lets other people decide how they are different.
 */
export type HowDifferent<How> = other$0.OtherPerson<{ [_: string]: How | null }>;

/**
 * Impersonator gets their fields from other people.
 */
export const Impersonator = other$0.OtherPerson;

/**
 * Impersonator gets their fields from other people.
 */
export type Impersonator = other$0.OtherPerson<number>;

/**
 * Person is not a number.
 */
export class Person {
    /**
     * They have a name.
     */
    "Name": string;

    /** Creates a new Person instance. */
    constructor($$source: Partial<Person> = {}) {
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
 * PrivatePerson gets their fields from hidden sources.
 */
export class PrivatePerson {
    /**
     * Nickname conceals a person's identity.
     */
    "Nickname": string;

    /**
     * They have a name.
     */
    "Name": string;

    /** Creates a new PrivatePerson instance. */
    constructor($$source: Partial<PrivatePerson> = {}) {
        if (!("Nickname" in $$source)) {
            this["Nickname"] = "";
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new PrivatePerson instance from a string or object.
     */
    static createFrom($$source: any = {}): PrivatePerson {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new PrivatePerson($$parsedSource as Partial<PrivatePerson>);
    }
}
