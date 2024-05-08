// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

export enum Embedded3 {
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
}

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
    "NamingThingsIsHard": boolean;

    /** Creates a new Embedded1 instance. */
    constructor(source: Partial<Embedded1> = {}) {
        if (!("Friends" in source)) {
            this["Friends"] = 0;
        }
        if (!("Vanish" in source)) {
            this["Vanish"] = 0;
        }
        if (!("StillThere" in source)) {
            this["StillThere"] = "";
        }
        if (!("NamingThingsIsHard" in source)) {
            this["NamingThingsIsHard"] = "false";
        }

        Object.assign(this, source);
    }

    /** Creates a new Embedded1 instance from a string or object. */
    static createFrom(source: string | object = {}): Embedded1 {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Embedded1(parsedSource as Partial<Embedded1>);
    }
}

export class Embedded2 {
    /**
     * Vanish should be omitted from Person because there is another field with same depth and no tag
     */
    "Vanish": boolean;

    /**
     * StillThereButRenamed should shadow in Person the other field with same depth and no json tag
     */
    "StillThere": Embedded3 | null;

    /** Creates a new Embedded2 instance. */
    constructor(source: Partial<Embedded2> = {}) {
        if (!("Vanish" in source)) {
            this["Vanish"] = "false";
        }
        if (!("StillThere" in source)) {
            this["StillThere"] = null;
        }

        Object.assign(this, source);
    }

    /** Creates a new Embedded2 instance from a string or object. */
    static createFrom(source: string | object = {}): Embedded2 {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Embedded2(parsedSource as Partial<Embedded2>);
    }
}

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
    "the person\'s partner ❤️": Person | null;
    "Friends": (Person | null)[];

    /**
     * NamingThingsIsHard is a law of programming
     */
    "NamingThingsIsHard": boolean;

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
    "emb4"?: embedded4;

    /** Creates a new Person instance. */
    constructor(source: Partial<Person> = {}) {
        if (!("Names" in source)) {
            this["Names"] = [];
        }
        if (!("the person\'s partner ❤️" in source)) {
            this["the person\'s partner ❤️"] = null;
        }
        if (!("Friends" in source)) {
            this["Friends"] = [];
        }
        if (!("NamingThingsIsHard" in source)) {
            this["NamingThingsIsHard"] = "false";
        }
        if (!("StillThere" in source)) {
            this["StillThere"] = null;
        }
        if (!("-" in source)) {
            this["-"] = 0;
        }
        if (!("Embedded3" in source)) {
            this["Embedded3"] = null;
        }
        if (!("StrangerNumber" in source)) {
            this["StrangerNumber"] = "0";
        }

        Object.assign(this, source);
    }

    /** Creates a new Person instance from a string or object. */
    static createFrom(source: string | object = {}): Person {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Person(parsedSource as Partial<Person>);
    }
}

export class embedded4 {
    /**
     * NamingThingsIsHard is a law of programming
     */
    "NamingThingsIsHard": boolean;

    /**
     * Friends should not be shadowed in Person as embedded4 is not embedded
     * from encoding/json's point of view;
     * however, it should be shadowed in Embedded1
     */
    "Friends": boolean;

    /** Creates a new embedded4 instance. */
    constructor(source: Partial<embedded4> = {}) {
        if (!("NamingThingsIsHard" in source)) {
            this["NamingThingsIsHard"] = "false";
        }
        if (!("Friends" in source)) {
            this["Friends"] = "false";
        }

        Object.assign(this, source);
    }

    /** Creates a new embedded4 instance from a string or object. */
    static createFrom(source: string | object = {}): embedded4 {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new embedded4(parsedSource as Partial<embedded4>);
    }
}