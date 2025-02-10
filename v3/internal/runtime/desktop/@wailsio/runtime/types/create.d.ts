/**
 * Any is a dummy creation function for simple or unknown types.
 * @template T
 * @param {any} source
 * @returns {T}
 */
export function Any<T>(source: any): T;
/**
 * ByteSlice is a creation function that replaces
 * null strings with empty strings.
 * @param {any} source
 * @returns {string}
 */
export function ByteSlice(source: any): string;
/**
 * Array takes a creation function for an arbitrary type
 * and returns an in-place creation function for an array
 * whose elements are of that type.
 * @template T
 * @param {(source: any) => T} element
 * @returns {(source: any) => T[]}
 */
export function Array<T>(element: (source: any) => T): (source: any) => T[];
/**
 * Map takes creation functions for two arbitrary types
 * and returns an in-place creation function for an object
 * whose keys and values are of those types.
 * @template K, V
 * @param {(source: any) => K} key
 * @param {(source: any) => V} value
 * @returns {(source: any) => { [_: K]: V }}
 */
export function Map<K, V>(key: (source: any) => K, value: (source: any) => V): (source: any) => {};
/**
 * Nullable takes a creation function for an arbitrary type
 * and returns a creation function for a nullable value of that type.
 * @template T
 * @param {(source: any) => T} element
 * @returns {(source: any) => (T | null)}
 */
export function Nullable<T>(element: (source: any) => T): (source: any) => T;
/**
 * Struct takes an object mapping field names to creation functions
 * and returns an in-place creation function for a struct.
 * @template {{ [_: string]: ((source: any) => any) }} T
 * @template {{ [Key in keyof T]?: ReturnType<T[Key]> }} U
 * @param {T} createField
 * @returns {(source: any) => U}
 */
export function Struct<T extends {
    [_: string]: (source: any) => any;
}, U extends { [Key in keyof T]?: ReturnType<T[Key]>; }>(createField: T): (source: any) => U;
