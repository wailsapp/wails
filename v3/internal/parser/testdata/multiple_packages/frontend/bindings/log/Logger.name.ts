// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from "/wails/runtime.js";

/**
 * Fatal is equivalent to l.Print() followed by a call to [os.Exit](1).
 */
export function Fatal(...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Fatal", v);
}

/**
 * Fatalf is equivalent to l.Printf() followed by a call to [os.Exit](1).
 */
export function Fatalf(format: string, ...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Fatalf", format, v);
}

/**
 * Fatalln is equivalent to l.Println() followed by a call to [os.Exit](1).
 */
export function Fatalln(...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Fatalln", v);
}

/**
 * Flags returns the output flags for the logger.
 * The flag bits are [Ldate], [Ltime], and so on.
 */
export function Flags(): Promise<number> {
    return Call.ByName("log.Logger.Flags");
}

/**
 * Output writes the output for a logging event. The string s contains
 * the text to print after the prefix specified by the flags of the
 * Logger. A newline is appended if the last character of s is not
 * already a newline. Calldepth is used to recover the PC and is
 * provided for generality, although at the moment on all pre-defined
 * paths it will be 2.
 */
export function Output(calldepth: number, s: string): Promise<void> {
    return Call.ByName("log.Logger.Output", calldepth, s);
}

/**
 * Panic is equivalent to l.Print() followed by a call to panic().
 */
export function Panic(...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Panic", v);
}

/**
 * Panicf is equivalent to l.Printf() followed by a call to panic().
 */
export function Panicf(format: string, ...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Panicf", format, v);
}

/**
 * Panicln is equivalent to l.Println() followed by a call to panic().
 */
export function Panicln(...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Panicln", v);
}

/**
 * Prefix returns the output prefix for the logger.
 */
export function Prefix(): Promise<string> {
    return Call.ByName("log.Logger.Prefix");
}

/**
 * Print calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Print].
 */
export function Print(...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Print", v);
}

/**
 * Printf calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Printf].
 */
export function Printf(format: string, ...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Printf", format, v);
}

/**
 * Println calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Println].
 */
export function Println(...v: any[]): Promise<void> {
    return Call.ByName("log.Logger.Println", v);
}

/**
 * SetFlags sets the output flags for the logger.
 * The flag bits are [Ldate], [Ltime], and so on.
 */
export function SetFlags(flag: number): Promise<void> {
    return Call.ByName("log.Logger.SetFlags", flag);
}

/**
 * SetPrefix sets the output prefix for the logger.
 */
export function SetPrefix(prefix: string): Promise<void> {
    return Call.ByName("log.Logger.SetPrefix", prefix);
}
