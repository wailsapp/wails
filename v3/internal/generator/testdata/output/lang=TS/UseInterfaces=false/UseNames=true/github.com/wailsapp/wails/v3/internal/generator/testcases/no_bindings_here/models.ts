// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as other$0 from "./other/models.js";

/**
 * HowDifferent is a curious kind of person
 * that lets other people decide how they are different.
 */
export class HowDifferent<How> {
    /**
     * They have a name as well.
     */
    "Name": string;

    /**
     * But they may have many differences.
     */
    "Differences": { [_: string]: How }[];

    /** Creates a new HowDifferent instance. */
    constructor($$source: Partial<HowDifferent<How>> = {}) {
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Differences" in $$source)) {
            this["Differences"] = [];
        }

        Object.assign(this, $$source);
    }

    /**
     * Given creation functions for each type parameter,
     * returns a creation function for a concrete instance
     * of the generic class HowDifferent.
     */
    static createFrom<How>($$createParamHow: (source: any) => How): ($$source?: any) => HowDifferent<How> {
        const $$createField1_0 = $$createType1($$createParamHow);
        return ($$source: any = {}) => {
            let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
            if ("Differences" in $$parsedSource) {
                $$parsedSource["Differences"] = $$createField1_0($$parsedSource["Differences"]);
            }
            return new HowDifferent<How>($$parsedSource as Partial<HowDifferent<How>>);
        };
    }
}

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

    /**
     * Exactly 4 sketchy friends.
     */
    "Friends": Impersonator[];

    /** Creates a new Person instance. */
    constructor($$source: Partial<Person> = {}) {
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Friends" in $$source)) {
            this["Friends"] = Array.from({ length: 4 }, () => (new Impersonator()));
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Person instance from a string or object.
     */
    static createFrom($$source: any = {}): Person {
        const $$createField1_0 = $$createType3;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("Friends" in $$parsedSource) {
            $$parsedSource["Friends"] = $$createField1_0($$parsedSource["Friends"]);
        }
        return new Person($$parsedSource as Partial<Person>);
    }
}

export class personImpl {
    /**
     * Nickname conceals a person's identity.
     */
    "Nickname": string;

    /**
     * They have a name.
     */
    "Name": string;

    /**
     * Exactly 4 sketchy friends.
     */
    "Friends": Impersonator[];

    /** Creates a new personImpl instance. */
    constructor($$source: Partial<personImpl> = {}) {
        if (!("Nickname" in $$source)) {
            this["Nickname"] = "";
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Friends" in $$source)) {
            this["Friends"] = Array.from({ length: 4 }, () => (new Impersonator()));
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new personImpl instance from a string or object.
     */
    static createFrom($$source: any = {}): personImpl {
        const $$createField2_0 = $$createType3;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("Friends" in $$parsedSource) {
            $$parsedSource["Friends"] = $$createField2_0($$parsedSource["Friends"]);
        }
        return new personImpl($$parsedSource as Partial<personImpl>);
    }
}

/**
 * PrivatePerson gets their fields from hidden sources.
 */
export const PrivatePerson = personImpl;

/**
 * PrivatePerson gets their fields from hidden sources.
 */
export type PrivatePerson = personImpl;

// Private type creation functions
const $$createType0 = ($$createParamHow) => $Create.Map($Create.Any, $$createParamHow);
const $$createType1 = ($$createParamHow) => $Create.Array($$createType0($$createParamHow));
const $$createType2 = other$0.OtherPerson.createFrom($Create.Any);
const $$createType3 = $Create.Array($$createType2);
