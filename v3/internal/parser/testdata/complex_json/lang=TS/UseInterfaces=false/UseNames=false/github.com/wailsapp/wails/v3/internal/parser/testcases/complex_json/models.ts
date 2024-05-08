// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Create as $Create} from "@wailsio/runtime";

import * as $internal from "./internal.ts";

export class Embedded1 {
    /**
     * Friends should be shadowed in Person by a field of lesser depth
     */
    "Friends": number;

    /**
     * Vanish should be omitted from Person because there is another field with same depth and no tag
     */
    "Vanish": number;

    /**
     * StillThere should be shadowed in Person by other field with same depth and a json tag
     */
    "StillThere": string;

    /**
     * NamingThingsIsHard is a law of programming
     */
    "NamingThingsIsHard": `${boolean}`;

    /** Creates a new Embedded1 instance. */
    constructor($$source: Partial<Embedded1> = {}) {
        if (!("Friends" in $$source)) {
            this["Friends"] = 0;
        }
        if (!("Vanish" in $$source)) {
            this["Vanish"] = 0;
        }
        if (!("StillThere" in $$source)) {
            this["StillThere"] = "";
        }
        if (!("NamingThingsIsHard" in $$source)) {
            this["NamingThingsIsHard"] = "false";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Embedded1 instance from a string or object.
     */
    static createFrom($$source: any = {}): Embedded1 {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        return new Embedded1($$parsedSource as Partial<Embedded1>);
    }
}

export type Embedded3 = string;

/**
 * Person represents a person
 */
export class Person {
    /**
     * Titles is optional in JSON
     */
    "Titles"?: Title[];

    /**
     * Names has a
     * multiline comment
     */
    "Names": string[];

    /**
     * Partner has a custom and complex JSON key
     */
    "Partner": Person | null;
    "Friends": (Person | null)[];

    /**
     * NamingThingsIsHard is a law of programming
     */
    "NamingThingsIsHard": `${boolean}`;

    /**
     * StillThereButRenamed should shadow in Person the other field with same depth and no json tag
     */
    "StillThere": Embedded3 | null;

    /**
     * StrangeNumber maps to "-"
     */
    "-": number;

    /**
     * Embedded3 should appear with key "Embedded3"
     */
    "Embedded3": Embedded3;

    /**
     * StrangerNumber is serialized as a string
     */
    "StrangerNumber": `${number}`;

    /**
     * StrangestString is optional and serialized as a JSON string
     */
    "StrangestString"?: `"${string}"`;

    /**
     * StringStrangest is serialized as a JSON string and optional
     */
    "StringStrangest"?: `"${string}"`;

    /**
     * embedded4 should be optional and appear with key "emb4"
     */
    "emb4"?: $internal.embedded4;

    /** Creates a new Person instance. */
    constructor($$source: Partial<Person> = {}) {
        if (!("Names" in $$source)) {
            this["Names"] = [];
        }
        if (!("Partner" in $$source)) {
            this["Partner"] = null;
        }
        if (!("Friends" in $$source)) {
            this["Friends"] = [];
        }
        if (!("NamingThingsIsHard" in $$source)) {
            this["NamingThingsIsHard"] = "false";
        }
        if (!("StillThere" in $$source)) {
            this["StillThere"] = null;
        }
        if (!("-" in $$source)) {
            this["-"] = 0;
        }
        if (!("Embedded3" in $$source)) {
            this["Embedded3"] = "";
        }
        if (!("StrangerNumber" in $$source)) {
            this["StrangerNumber"] = "0";
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Person instance from a string or object.
     */
    static createFrom($$source: any = {}): Person {
        const $$createField0_0 = $$createType0;
        const $$createField1_0 = $$createType1;
        const $$createField2_0 = $$createType3;
        const $$createField3_0 = $$createType4;
        const $$createField11_0 = $$createType5;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("Titles" in $$parsedSource) {
            $$parsedSource["Titles"] = $$createField0_0($$parsedSource["Titles"]);
        }
        if ("Names" in $$parsedSource) {
            $$parsedSource["Names"] = $$createField1_0($$parsedSource["Names"]);
        }
        if ("Partner" in $$parsedSource) {
            $$parsedSource["Partner"] = $$createField2_0($$parsedSource["Partner"]);
        }
        if ("Friends" in $$parsedSource) {
            $$parsedSource["Friends"] = $$createField3_0($$parsedSource["Friends"]);
        }
        if ("emb4" in $$parsedSource) {
            $$parsedSource["emb4"] = $$createField11_0($$parsedSource["emb4"]);
        }
        return new Person($$parsedSource as Partial<Person>);
    }
}

/**
 * Title is a title
 */
export enum Title {
    Dr = "Dr",
    Miss = "Miss",

    /**
     * Mister is a title
     */
    Mister = "Mr",
    Mrs = "Mrs",
    Ms = "Ms",
};

// Private type creation functions
const $$createType0 = $Create.Array($Create.Any);
const $$createType1 = $Create.Array($Create.Any);
const $$createType2 = Person.createFrom;
const $$createType3 = $Create.Nullable($$createType2);
const $$createType4 = $Create.Array($$createType3);
const $$createType5 = $internal.embedded4.createFrom;
