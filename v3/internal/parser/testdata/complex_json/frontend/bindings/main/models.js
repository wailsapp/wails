// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * @readonly
 * @enum
 */
export const Embedded3 = {
};

/**
 * Title is a title
 * @readonly
 * @enum
 */
export const Title = {
    /**
     * Mister is a title
     */
    Mister: "Mr",
    Miss: "Miss",
    Ms: "Ms",
    Mrs: "Mrs",
    Dr: "Dr",
};

export class Embedded1 {
    /**
     * Creates a new Embedded1 instance.
     * @param {Partial<Embedded1>} [source = {}] - The source object to create the Embedded1.
     */
    constructor(source = {}) {
        if (!("Friends" in source)) {
            /**
             * Friends should be shadowed in Person by a field of lesser depth
             * @member
             * @type {number}
             */
            this["Friends"] = 0;
        }
        if (!("Vanish" in source)) {
            /**
             * Vanish should be omitted from Person because there is another field with same depth and no tag
             * @member
             * @type {number}
             */
            this["Vanish"] = 0;
        }
        if (!("StillThere" in source)) {
            /**
             * StillThere should be shadowed in Person by other field with same depth and a json tag
             * @member
             * @type {string}
             */
            this["StillThere"] = "";
        }
        if (!("NamingThingsIsHard" in source)) {
            /**
             * NamingThingsIsHard is a law of programming
             * @member
             * @type {`${boolean}`}
             */
            this["NamingThingsIsHard"] = "false";
        }

        Object.assign(this, source);
    }

    /**
     * Creates a new Embedded1 instance from a string or object.
     * @param {string|object} source - The source data to create a Embedded1 instance from.
     * @returns {Embedded1} A new Embedded1 instance.
     */
    static createFrom(source) {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Embedded1(parsedSource);
    }
};

export class Embedded2 {
    /**
     * Creates a new Embedded2 instance.
     * @param {Partial<Embedded2>} [source = {}] - The source object to create the Embedded2.
     */
    constructor(source = {}) {
        if (!("Vanish" in source)) {
            /**
             * Vanish should be omitted from Person because there is another field with same depth and no tag
             * @member
             * @type {boolean}
             */
            this["Vanish"] = false;
        }
        if (!("StillThere" in source)) {
            /**
             * StillThereButRenamed should shadow in Person the other field with same depth and no json tag
             * @member
             * @type {Embedded3 | null}
             */
            this["StillThere"] = null;
        }

        Object.assign(this, source);
    }

    /**
     * Creates a new Embedded2 instance from a string or object.
     * @param {string|object} source - The source data to create a Embedded2 instance from.
     * @returns {Embedded2} A new Embedded2 instance.
     */
    static createFrom(source) {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Embedded2(parsedSource);
    }
};

/**
 * Person represents a person
 */
export class Person {
    /**
     * Creates a new Person instance.
     * @param {Partial<Person>} [source = {}] - The source object to create the Person.
     */
    constructor(source = {}) {
        if (/** @type {any} */(false)) {
            /**
             * Titles is optional in JSON
             * @member
             * @type {Title[] | undefined}
             */
            this["Titles"] = [];
        }
        if (!("Names" in source)) {
            /**
             * Names has a
             * multiline comment
             * @member
             * @type {string[]}
             */
            this["Names"] = [];
        }
        if (!("the person\'s partner ❤️" in source)) {
            /**
             * Partner has a custom and complex JSON key
             * @member
             * @type {Person | null}
             */
            this["the person\'s partner ❤️"] = null;
        }
        if (!("Friends" in source)) {
            /**
             * @member
             * @type {(Person | null)[]}
             */
            this["Friends"] = [];
        }
        if (!("NamingThingsIsHard" in source)) {
            /**
             * NamingThingsIsHard is a law of programming
             * @member
             * @type {`${boolean}`}
             */
            this["NamingThingsIsHard"] = "false";
        }
        if (!("StillThere" in source)) {
            /**
             * StillThereButRenamed should shadow in Person the other field with same depth and no json tag
             * @member
             * @type {Embedded3 | null}
             */
            this["StillThere"] = null;
        }
        if (!("-" in source)) {
            /**
             * StrangeNumber maps to "-"
             * @member
             * @type {number}
             */
            this["-"] = 0;
        }
        if (!("Embedded3" in source)) {
            /**
             * Embedded3 should appear with key "Embedded3"
             * @member
             * @type {Embedded3}
             */
            this["Embedded3"] = null;
        }
        if (!("StrangerNumber" in source)) {
            /**
             * StrangerNumber is serialized as a string
             * @member
             * @type {`${number}`}
             */
            this["StrangerNumber"] = "0";
        }
        if (/** @type {any} */(false)) {
            /**
             * StrangestString is optional and serialized as a JSON string
             * @member
             * @type {`"${string}"` | undefined}
             */
            this["StrangestString"] = "\"\"";
        }
        if (/** @type {any} */(false)) {
            /**
             * StringStrangest is serialized as a JSON string and optional
             * @member
             * @type {`"${string}"` | undefined}
             */
            this["StringStrangest"] = "\"\"";
        }
        if (/** @type {any} */(false)) {
            /**
             * embedded4 should be optional and appear with key "emb4"
             * @member
             * @type {embedded4 | undefined}
             */
            this["emb4"] = (new embedded4());
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

export class embedded4 {
    /**
     * Creates a new embedded4 instance.
     * @param {Partial<embedded4>} [source = {}] - The source object to create the embedded4.
     */
    constructor(source = {}) {
        if (!("NamingThingsIsHard" in source)) {
            /**
             * NamingThingsIsHard is a law of programming
             * @member
             * @type {`${boolean}`}
             */
            this["NamingThingsIsHard"] = "false";
        }
        if (!("Friends" in source)) {
            /**
             * Friends should not be shadowed in Person as embedded4 is not embedded
             * from encoding/json's point of view;
             * however, it should be shadowed in Embedded1
             * @member
             * @type {boolean}
             */
            this["Friends"] = false;
        }

        Object.assign(this, source);
    }

    /**
     * Creates a new embedded4 instance from a string or object.
     * @param {string|object} source - The source data to create a embedded4 instance from.
     * @returns {embedded4} A new embedded4 instance.
     */
    static createFrom(source) {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new embedded4(parsedSource);
    }
};
