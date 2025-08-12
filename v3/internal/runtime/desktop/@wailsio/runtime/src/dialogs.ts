/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

import {newRuntimeCaller, objectNames} from "./runtime.js";
import { nanoid } from './nanoid.js';

// setup
window._wails = window._wails || {};
window._wails.dialogErrorCallback = dialogErrorCallback;
window._wails.dialogResultCallback = dialogResultCallback;

type PromiseResolvers = Omit<PromiseWithResolvers<any>, "promise">;

const call = newRuntimeCaller(objectNames.Dialog);
const dialogResponses = new Map<string, PromiseResolvers>();

// Define constants from the `methods` object in Title Case
const DialogInfo = 0;
const DialogWarning = 1;
const DialogError = 2;
const DialogQuestion = 3;
const DialogOpenFile = 4;
const DialogSaveFile = 5;

export interface OpenFileDialogOptions {
    /** Indicates if directories can be chosen. */
    CanChooseDirectories?: boolean;
    /** Indicates if files can be chosen. */
    CanChooseFiles?: boolean;
    /** Indicates if directories can be created. */
    CanCreateDirectories?: boolean;
    /** Indicates if hidden files should be shown. */
    ShowHiddenFiles?: boolean;
    /** Indicates if aliases should be resolved. */
    ResolvesAliases?: boolean;
    /** Indicates if multiple selection is allowed. */
    AllowsMultipleSelection?: boolean;
    /** Indicates if the extension should be hidden. */
    HideExtension?: boolean;
    /** Indicates if hidden extensions can be selected. */
    CanSelectHiddenExtension?: boolean;
    /** Indicates if file packages should be treated as directories. */
    TreatsFilePackagesAsDirectories?: boolean;
    /** Indicates if other file types are allowed. */
    AllowsOtherFiletypes?: boolean;
    /** Array of file filters. */
    Filters?: FileFilter[];
    /** Title of the dialog. */
    Title?: string;
    /** Message to show in the dialog. */
    Message?: string;
    /** Text to display on the button. */
    ButtonText?: string;
    /** Directory to open in the dialog. */
    Directory?: string;
    /** Indicates if the dialog should appear detached from the main window. */
    Detached?: boolean;
}

export interface SaveFileDialogOptions {
    /** Default filename to use in the dialog. */
    Filename?: string;
    /** Indicates if directories can be chosen. */
    CanChooseDirectories?: boolean;
    /** Indicates if files can be chosen. */
    CanChooseFiles?: boolean;
    /** Indicates if directories can be created. */
    CanCreateDirectories?: boolean;
    /** Indicates if hidden files should be shown. */
    ShowHiddenFiles?: boolean;
    /** Indicates if aliases should be resolved. */
    ResolvesAliases?: boolean;
    /** Indicates if the extension should be hidden. */
    HideExtension?: boolean;
    /** Indicates if hidden extensions can be selected. */
    CanSelectHiddenExtension?: boolean;
    /** Indicates if file packages should be treated as directories. */
    TreatsFilePackagesAsDirectories?: boolean;
    /** Indicates if other file types are allowed. */
    AllowsOtherFiletypes?: boolean;
    /** Array of file filters. */
    Filters?: FileFilter[];
    /** Title of the dialog. */
    Title?: string;
    /** Message to show in the dialog. */
    Message?: string;
    /** Text to display on the button. */
    ButtonText?: string;
    /** Directory to open in the dialog. */
    Directory?: string;
    /** Indicates if the dialog should appear detached from the main window. */
    Detached?: boolean;
}

export interface MessageDialogOptions {
    /** The title of the dialog window. */
    Title?: string;
    /** The main message to show in the dialog. */
    Message?: string;
    /** Array of button options to show in the dialog. */
    Buttons?: Button[];
    /** True if the dialog should appear detached from the main window (if applicable). */
    Detached?: boolean;
}

export interface Button {
    /** Text that appears within the button. */
    Label?: string;
    /** True if the button should cancel an operation when clicked. */
    IsCancel?: boolean;
    /** True if the button should be the default action when the user presses enter. */
    IsDefault?: boolean;
}

export interface FileFilter {
    /** Display name for the filter, it could be "Text Files", "Images" etc. */
    DisplayName?: string;
    /** Pattern to match for the filter, e.g. "*.txt;*.md" for text markdown files. */
    Pattern?: string;
}

/**
 * Handles the result of a dialog request.
 *
 * @param id - The id of the request to handle the result for.
 * @param data - The result data of the request.
 * @param isJSON - Indicates whether the data is JSON or not.
 */
function dialogResultCallback(id: string, data: string, isJSON: boolean): void {
    let resolvers = getAndDeleteResponse(id);
    if (!resolvers) {
        return;
    }

    if (isJSON) {
        try {
            resolvers.resolve(JSON.parse(data));
        } catch (err: any) {
            resolvers.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
        }
    } else {
        resolvers.resolve(data);
    }
}

/**
 * Handles the error from a dialog request.
 *
 * @param id - The id of the promise handler.
 * @param message - An error message.
 */
function dialogErrorCallback(id: string, message: string): void {
    getAndDeleteResponse(id)?.reject(new window.Error(message));
}

/**
 * Retrieves and removes the response associated with the given ID from the dialogResponses map.
 *
 * @param id - The ID of the response to be retrieved and removed.
 * @returns The response object associated with the given ID, if any.
 */
function getAndDeleteResponse(id: string): PromiseResolvers | undefined {
    const response = dialogResponses.get(id);
    dialogResponses.delete(id);
    return response;
}

/**
 * Generates a unique ID using the nanoid library.
 *
 * @returns A unique ID that does not exist in the dialogResponses set.
 */
function generateID(): string {
    let result;
    do {
        result = nanoid();
    } while (dialogResponses.has(result));
    return result;
}

/**
 * Presents a dialog of specified type with the given options.
 *
 * @param type - Dialog type.
 * @param options - Options for the dialog.
 * @returns A promise that resolves with result of dialog.
 */
function dialog(type: number, options: MessageDialogOptions | OpenFileDialogOptions | SaveFileDialogOptions = {}): Promise<any> {
    const id = generateID();
    return new Promise((resolve, reject) => {
        dialogResponses.set(id, { resolve, reject });
        call(type, Object.assign({ "dialog-id": id }, options)).catch((err: any) => {
            dialogResponses.delete(id);
            reject(err);
        });
    });
}

/**
 * Presents an info dialog.
 *
 * @param options - Dialog options
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Info(options: MessageDialogOptions): Promise<string> { return dialog(DialogInfo, options); }

/**
 * Presents a warning dialog.
 *
 * @param options - Dialog options.
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Warning(options: MessageDialogOptions): Promise<string> { return dialog(DialogWarning, options); }

/**
 * Presents an error dialog.
 *
 * @param options - Dialog options.
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Error(options: MessageDialogOptions): Promise<string> { return dialog(DialogError, options); }

/**
 * Presents a question dialog.
 *
 * @param options - Dialog options.
 * @returns A promise that resolves with the label of the chosen button.
 */
export function Question(options: MessageDialogOptions): Promise<string> { return dialog(DialogQuestion, options); }

/**
 * Presents a file selection dialog to pick one or more files to open.
 *
 * @param options - Dialog options.
 * @returns Selected file or list of files, or a blank string/empty list if no file has been selected.
 */
export function OpenFile(options: OpenFileDialogOptions & { AllowsMultipleSelection: true }): Promise<string[]>;
export function OpenFile(options: OpenFileDialogOptions & { AllowsMultipleSelection?: false | undefined }): Promise<string>;
export function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>;
export function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]> { return dialog(DialogOpenFile, options) ?? []; }

/**
 * Presents a file selection dialog to pick a file to save.
 *
 * @param options - Dialog options.
 * @returns Selected file, or a blank string if no file has been selected.
 */
export function SaveFile(options: SaveFileDialogOptions): Promise<string> { return dialog(DialogSaveFile, options); }
