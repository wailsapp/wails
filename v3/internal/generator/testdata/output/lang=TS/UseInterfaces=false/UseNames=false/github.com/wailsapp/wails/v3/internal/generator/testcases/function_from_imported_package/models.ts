// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as services$0 from "./services/models.js";

/**
 * Person is a person
 */
export class Person {
    "Name": string;
    "Address": services$0.Address | null;

    /** Creates a new Person instance. */
    constructor($$source: Partial<Person> = {}) {
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Address" in $$source)) {
            this["Address"] = null;
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Person instance from a string or object.
     */
    static createFrom($$source: any = {}): Person {
        const $$createField1_0 = $$createType1;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("Address" in $$parsedSource) {
            $$parsedSource["Address"] = $$createField1_0($$parsedSource["Address"]);
        }
        return new Person($$parsedSource as Partial<Person>);
    }
}

// Private type creation functions
const $$createType0 = services$0.Address.createFrom;
const $$createType1 = $Create.Nullable($$createType0);
