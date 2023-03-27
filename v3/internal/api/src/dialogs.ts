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

import {newRuntimeCaller} from "./runtime";

import {nanoid} from 'nanoid/non-secure';

let call = newRuntimeCaller("dialog");

let dialogResponses = new Map();

function generateID() {
    let result;
    do {
        result = nanoid();
    } while (dialogResponses.has(result));
    return result;
}

export function dialogCallback(id: string, data: any, isJSON: boolean) {
    let p = dialogResponses.get(id);
    if (p) {
        if (isJSON) {
            p.resolve(JSON.parse(data));
        } else {
            p.resolve(data);
        }
        dialogResponses.delete(id);
    }
}

export function dialogErrorCallback(id: string, message: string) {
    let p = dialogResponses.get(id);
    if (p) {
        p.reject(message);
        dialogResponses.delete(id);
    }
}

function dialog(type: string, options: any) {
    return new Promise((resolve, reject) => {
        let id = generateID();
        options = options || {};
        options["dialog-id"] = id;
        dialogResponses.set(id, {resolve, reject});
        call(type, options).catch((error) => {
            reject(error);
            dialogResponses.delete(id);
        });
    });
}

// DialogType is the type of dialog
export enum DialogType {
    InfoDialog,
    QuestionDialog,
    WarningDialog,
    ErrorDialog,
    OpenDirectoryDialog
}

// Button is a button in a dialog
export interface Button {
    Label:     string
    IsCancel:  boolean
    IsDefault: boolean
}

// MessageDialogOptions is the options for a message dialog
export interface MessageDialogOptions {
    DialogType: DialogType
    Title: string
    Message: string
    Buttons: Button[]
}

// OpenFileDialogOptions is the options for an open file dialog
export interface OpenFileDialogOptions {
    CanChooseDirectories:            boolean
    CanChooseFiles:                  boolean
    CanCreateDirectories:            boolean
    ShowHiddenFiles:                 boolean
    ResolvesAliases:                 boolean
    AllowsMultipleSelection:         boolean
    HideExtension:                   boolean
    CanSelectHiddenExtension:        boolean
    TreatsFilePackagesAsDirectories: boolean
    AllowsOtherFileTypes:            boolean
    Filters:                         FileFilter[]

    Title:      string
    Message:    string
    ButtonText: string
    Directory:  string
}

export interface FileFilter {
    DisplayName: string // Filter information EG: "Image Files (*.jpg, *.png)"
    Pattern:     string // semicolon separated list of extensions, EG: "*.jpg;*.png"
}

export interface SaveFileDialogOptions {
    CanCreateDirectories:            boolean
    ShowHiddenFiles:                 boolean
    CanSelectHiddenExtension:        boolean
    AllowOtherFileTypes:             boolean
    HideExtension:                   boolean
    TreatsFilePackagesAsDirectories: boolean
    Message:                         string
    Directory:                       string
    Filename:                        string
    ButtonText:                      string
}

export function Info(options: MessageDialogOptions) {
    return dialog("Info", options);
}

export function Warning(options: MessageDialogOptions) {
    return dialog("Warning", options);
}

export function Error(options: MessageDialogOptions) {
    return dialog("Error", options);
}

export function Question(options: MessageDialogOptions) {
    return dialog("Question", options);
}

export function OpenFile(options: OpenFileDialogOptions) {
    return dialog("OpenFile", options);
}

export function SaveFile(options: SaveFileDialogOptions) {
    return dialog("SaveFile", options);
}

