/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/**
 * Any is a dummy creation function for simple or unknown types.
 */
export function Any<T>(source: any): T {
    return source;
}

/**
 * ByteSlice is a creation function that replaces
 * null strings with empty strings.
 */
export function ByteSlice(source: any): string {
    return ((source == null) ? "" : source);
}

/**
 * Array takes a creation function for an arbitrary type
 * and returns an in-place creation function for an array
 * whose elements are of that type.
 */
export function Array<T>(element: (source: any) => T): (source: any) => T[] {
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
 */
export function Map<V>(key: (source: any) => string, value: (source: any) => V): (source: any) => Record<string, V> {
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
 */
export function Nullable<T>(element: (source: any) => T): (source: any) => (T | null) {
    if (element === Any) {
        return Any;
    }

    return (source) => (source === null ? null : element(source));
}

/**
 * Struct takes an object mapping field names to creation functions
 * and returns an in-place creation function for a struct.
 */
export function Struct<
    T extends { [_: string]: ((source: any) => any) },
    U extends { [Key in keyof T]?: ReturnType<T[Key]> }
>(createField: T): (source: any) => U {
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
