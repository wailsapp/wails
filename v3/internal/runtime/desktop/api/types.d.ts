
export interface Button {
    // The label of the button
    Label?: string;
    // True if this button is the cancel button (selected when pressing escape)
    IsCancel?: boolean;
    // True if this button is the default button (selected when pressing enter)
    IsDefault?: boolean;
}

interface MessageDialogOptions {
    // The title for the dialog
    Title?: string;
    // The message to display
    Message?: string;
    // The buttons to use on the dialog
    Buttons?: Button[];
}

export interface OpenFileDialogOptions {
    // Allows the user to be able to select directories
    CanChooseDirectories?: boolean;
    // Allows the user to be able to select files
    CanChooseFiles?: boolean;
    // Provide an option to create directories in the dialog
    CanCreateDirectories?: boolean;
    // Makes the dialog show hidden files
    ShowHiddenFiles?: boolean;
    // Whether the dialog should follow filesystem aliases
    ResolvesAliases?: boolean;
    // Allow the user to select multiple files or directories
    AllowsMultipleSelection?: boolean;
    // Hide the extension when showing the filename
    HideExtension?: boolean;
    // Allow the user to select files where the system hides their extensions
    CanSelectHiddenExtension?: boolean;
    // Treats file packages as directories, e.g. .app on macOS
    TreatsFilePackagesAsDirectories?: boolean;
    // Allows selection of filetypes not specified in the filters
    AllowsOtherFiletypes?: boolean;
    // The file filters to use in the dialog
    Filters?: FileFilter[];
    // The title of the dialog
    Title?: string;
    // The message to display
    Message?: string;
    // The label for the select button
    ButtonText?: string;
    // The default directory to open the dialog in
    Directory?: string;
}
export interface FileFilter {
    // The display name for the filter, e.g. "Text Files"
    DisplayName?: string;
    // The pattern to use for the filter, e.g. "*.txt;*.md"
    Pattern?: string;
}
export interface SaveFileDialogOptions {
    // Provide an option to create directories in the dialog
    CanCreateDirectories?: boolean;
    // Makes the dialog show hidden files
    ShowHiddenFiles?: boolean;
    // Allow the user to select files where the system hides their extensions
    CanSelectHiddenExtension?: boolean;
    // Allows selection of filetypes not specified in the filters
    AllowOtherFiletypes?: boolean;
    // Hide the extension when showing the filename
    HideExtension?: boolean;
    // Treats file packages as directories, e.g. .app on macOS
    TreatsFilePackagesAsDirectories?: boolean;
    // The message to show in the dialog
    Message?: string;
    // The default directory to open the dialog in
    Directory?: string;
    // The default filename to use in the dialog
    Filename?: string;
    // The label for the select button
    ButtonText?: string;
}

export interface Screen {
    // The screen ID
    Id: string;
    // The screen name
    Name: string;
    // The screen scale. 1 = standard resolution, 2: 2x retina, etc.
    Scale: number;
    // The X position of the screen
    X: number;
    // The Y position of the screen
    Y: number;
    // The width and height of the screen
    Size: Size;
    // The bounds of the screen
    Bounds: Rect;
    // The work area of the screen
    WorkArea: Rect;
    // True if this is the primary screen
    IsPrimary: boolean;
    // The rotation of the screen
    Rotation: number;
}
export interface Rect {
    X: number;
    Y: number;
    Width: number;
    Height: number;
}

export interface WailsEvent {
    // The name of the event
    Name: string;
    // The data associated with the event
    Data?: any;
}

export interface Size {
    Width: number;
    Height: number;
}
export interface Position {
    X: number;
    Y: number;
}
