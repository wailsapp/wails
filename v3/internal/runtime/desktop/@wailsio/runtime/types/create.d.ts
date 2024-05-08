/**
 * Any is a dummy creation function for simple or unknown types.
 * @template T
 * @param {any} source
 * @returns {T}
 */
export function Any<T>(source: any): T;
/**
 * Array takes a creation function for an arbitrary type
 * and returns an in-place creation function for an array
 * whose elements are of that type.
 * @template T
 * @param {(any) => T} element
 * @returns {(any) => T[]}
 */
export function Array<T>(element: (any: any) => T): (any: any) => T[];
/**
 * Map takes creation functions for two arbitrary types
 * and returns an in-place creation function for an object
 * whose keys and values are of those types.
 * @template K, V
 * @param {(any) => K} key
 * @param {(any) => V} value
 * @returns {(any) => { [_: K]: V }}
 */
export function Map<K, V>(key: (any: any) => K, value: (any: any) => V): (any: any) => {};
/**
 * Nullable takes a creation function for an arbitrary type
 * and returns a creation function for a nullable value of that type.
 * @template T
 * @param {(any) => T} element
 * @returns {(any) => (T | null)}
 */
export function Nullable<T>(element: (any: any) => T): (any: any) => T;
/**
 * Struct takes an object mapping field names to creation functions
 * and returns an in-place creation function for a struct.
 * @template {{ [_: string]: ((any) => any) }} T
 * @template {{ [Key in keyof T]?: ReturnType<T[Key]> }} U
 * @param {T} createField
 * @returns {(any) => U}
 */
export function Struct<T extends {
    [_: string]: (any: any) => any;
}, U extends { [Key in keyof T]?: ReturnType<T[Key]>; }>(createField: T): (any: any) => U;
