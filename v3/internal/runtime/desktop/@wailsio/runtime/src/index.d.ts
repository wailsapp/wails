declare module '@wailsio/runtime';

/**
 * Describes the properties of a Button object.
 */
export interface Button {
    // Text that appears within the button.
    Label?: string;
    // True if the button should cancel an operation when clicked.
    IsCancel?: boolean;
    // True if the button should be the default action when the user presses enter.
    IsDefault?: boolean;
}

/**
 * Describes the properties of a FileFilter object.
 */
export interface FileFilter {
    // Display name for the filter, it could be "Text Files", "Images" etc.
    DisplayName?: string;
    // Pattern to match for the filter, e.g. "*.txt;*.md" for text markdown files.
    Pattern?: string;
}

/**
 * Describes the properties of a Position object.
 */
export interface Position {
    X: number;
    Y: number;
}

/**
 * Describes the properties of a Size object.
 * Extends the properties of Position.
 */
export type Size = Position;

/**
 * Describes the properties of a Rect object.
 * Extends the properties of Size.
 */
export type Rect = Size & Position;

/**
 * Describes the possible rotations a Screen can have.
 */
export enum Rotation {
    Zero = 0,
    Ninety = 90,
    OneEighty = 180,
    TwoSeventy = 270
}

/**
 * Describes the properties of a Screen object.
 */
export interface Screen {
    // Unique identifier for the screen.
    Id: string;
    // Human readable name of the screen.
    Name: string;
    // The resolution scale of the screen. 1 = standard resolution, 2 = high (Retina), etc.
    Scale: number;
    // Contains the X and Y coordinates of the screen's position.
    Position: Position;
    // Contains the width and height of the screen.
    Size: Size;
    // Contains the bounds of the screen in terms of X, Y, Width, and Height.
    Bounds: Rect;
    // Contains the area of the screen that is actually usable (excluding taskbar and other system UI).
    WorkArea: Rect;
    // True if this is the primary monitor selected by the user in the operating system.
    IsPrimary: boolean;
    // The rotation of the screen. Can be one of 0, 90, 180, 270 degrees.
    Rotation: Rotation;
}

/**
 * Describes the properties of a MessageDialogOptions object.
 */
export interface MessageDialogOptions {
    // The title of the dialog window.
    Title?: string;
    // The main message to show in the dialog.
    Message?: string;
    // Array of button options to show in the dialog.
    Buttons?: Button[];
    // True if the dialog should appear detached from the main window (if applicable).
    Detached?: boolean;
}

/**
 * Describes the properties common to OpenFileDialogOptions and SaveFileDialogOptions objects.
 */
export interface DialogOptions {
    CanChooseDirectories?: boolean;
    CanChooseFiles?: boolean;
    CanCreateDirectories?: boolean;
    ShowHiddenFiles?: boolean;
    ResolvesAliases?: boolean;
    AllowsMultipleSelection?: boolean;
    HideExtension?: boolean;
    CanSelectHiddenExtension?: boolean;
    TreatsFilePackagesAsDirectories?: boolean;
    AllowsOtherFiletypes?: boolean;
    Filters?: FileFilter[];
    Title?: string;
    Message?: string;
    ButtonText?: string;
    Directory?: string;
    Detached?: boolean;
}

/**
 * Describes properties unique to OpenFileDialogOptions.
 */
export type OpenFileDialogOptions = DialogOptions;

/**
 * Describes properties unique to SaveFileDialogOptions.
 */
export type SaveFileDialogOptions = DialogOptions & {
    // Default filename to use for the save dialog.
    Filename?: string;
};

/**
 * Describes the properties of a WailsEvent object.
 */
export interface WailsEvent {
    // The name of the event.
    Name: string;
    // Any data associated with the event.
    Data?: unknown;
}
