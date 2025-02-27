/**
 * Any is a dummy creation function for simple or unknown types.
 */
export declare function Any<T = any>(source: any): T;
/**
 * ByteSlice is a creation function that replaces
 * null strings with empty strings.
 */
export declare function ByteSlice(source: any): string;
/**
 * Array takes a creation function for an arbitrary type
 * and returns an in-place creation function for an array
 * whose elements are of that type.
 */
export declare function Array<T = any>(element: (source: any) => T): (source: any) => T[];
/**
 * Map takes creation functions for two arbitrary types
 * and returns an in-place creation function for an object
 * whose keys and values are of those types.
 */
export declare function Map<V = any>(key: (source: any) => string, value: (source: any) => V): (source: any) => Record<string, V>;
/**
 * Nullable takes a creation function for an arbitrary type
 * and returns a creation function for a nullable value of that type.
 */
export declare function Nullable<T = any>(element: (source: any) => T): (source: any) => (T | null);
/**
 * Struct takes an object mapping field names to creation functions
 * and returns an in-place creation function for a struct.
 */
export declare function Struct(createField: Record<string, (source: any) => any>): <U extends Record<string, any> = any>(source: any) => U;
