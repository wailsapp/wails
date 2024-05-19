/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

/**
 * Any is a dummy creation function for simple or unknown types.
 * @template T
 * @param {any} source
 * @returns {T}
 */
export function Any(source) {
    return /** @type {T} */(source);
}

/**
 * ByteSlice is a creation function that replaces
 * null strings with empty strings.
 * @param {any} source
 * @returns {string}
 */
export function ByteSlice(source) {
    return /** @type {any} */((source == null) ? "" : source);
}

/**
 * Array takes a creation function for an arbitrary type
 * and returns an in-place creation function for an array
 * whose elements are of that type.
 * @template T
 * @param {(source: any) => T} element
 * @returns {(source: any) => T[]}
 */
export function Array(element) {
    if (element === Any) {
        return (source) => (source === null ? [] : source);
    }

    return (source) => {
        if (source === null) {
            return [];
        }
        for (let i = 0; i < source.length; i++) {
            source[i] = element(source[i]);
        }
        return source;
    };
}

/**
 * Map takes creation functions for two arbitrary types
 * and returns an in-place creation function for an object
 * whose keys and values are of those types.
 * @template K, V
 * @param {(source: any) => K} key
 * @param {(source: any) => V} value
 * @returns {(source: any) => { [_: K]: V }}
 */
export function Map(key, value) {
    if (value === Any) {
        return (source) => (source === null ? {} : source);
    }

    return (source) => {
        if (source === null) {
            return {};
        }
        for (const key in source) {
            source[key] = value(source[key]);
        }
        return source;
    };
}

/**
 * Nullable takes a creation function for an arbitrary type
 * and returns a creation function for a nullable value of that type.
 * @template T
 * @param {(source: any) => T} element
 * @returns {(source: any) => (T | null)}
 */
export function Nullable(element) {
    if (element === Any) {
        return Any;
    }

    return (source) => (source === null ? null : element(source));
}

/**
 * Struct takes an object mapping field names to creation functions
 * and returns an in-place creation function for a struct.
 * @template {{ [_: string]: ((source: any) => any) }} T
 * @template {{ [Key in keyof T]?: ReturnType<T[Key]> }} U
 * @param {T} createField
 * @returns {(source: any) => U}
 */
export function Struct(createField) {
    let allAny = true;
    for (const name in createField) {
        if (createField[name] !== Any) {
            allAny = false;
            break;
        }
    }
    if (allAny) {
        return Any;
    }

    return (source) => {
        for (const name in createField) {
            if (name in source) {
                source[name] = createField[name](source[name]);
            }
        }
        return source;
    };
}
