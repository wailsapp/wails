export function Info(options: MessageDialogOptions): Promise<string>;
export function Warning(options: MessageDialogOptions): Promise<string>;
export function Error(options: MessageDialogOptions): Promise<string>;
export function Question(options: MessageDialogOptions): Promise<string>;
export function OpenFile(options: OpenFileDialogOptions): Promise<string[] | string>;
export function SaveFile(options: SaveFileDialogOptions): Promise<string>;
export type OpenFileDialogOptions = {
    /**
     * - Indicates if directories can be chosen.
     */
    CanChooseDirectories?: boolean;
    /**
     * - Indicates if files can be chosen.
     */
    CanChooseFiles?: boolean;
    /**
     * - Indicates if directories can be created.
     */
    CanCreateDirectories?: boolean;
    /**
     * - Indicates if hidden files should be shown.
     */
    ShowHiddenFiles?: boolean;
    /**
     * - Indicates if aliases should be resolved.
     */
    ResolvesAliases?: boolean;
    /**
     * - Indicates if multiple selection is allowed.
     */
    AllowsMultipleSelection?: boolean;
    /**
     * - Indicates if the extension should be hidden.
     */
    HideExtension?: boolean;
    /**
     * - Indicates if hidden extensions can be selected.
     */
    CanSelectHiddenExtension?: boolean;
    /**
     * - Indicates if file packages should be treated as directories.
     */
    TreatsFilePackagesAsDirectories?: boolean;
    /**
     * - Indicates if other file types are allowed.
     */
    AllowsOtherFiletypes?: boolean;
    /**
     * - Array of file filters.
     */
    Filters?: FileFilter[];
    /**
     * - Title of the dialog.
     */
    Title?: string;
    /**
     * - Message to show in the dialog.
     */
    Message?: string;
    /**
     * - Text to display on the button.
     */
    ButtonText?: string;
    /**
     * - Directory to open in the dialog.
     */
    Directory?: string;
    /**
     * - Indicates if the dialog should appear detached from the main window.
     */
    Detached?: boolean;
};
export type SaveFileDialogOptions = {
    /**
     * - Default filename to use in the dialog.
     */
    Filename?: string;
    /**
     * - Indicates if directories can be chosen.
     */
    CanChooseDirectories?: boolean;
    /**
     * - Indicates if files can be chosen.
     */
    CanChooseFiles?: boolean;
    /**
     * - Indicates if directories can be created.
     */
    CanCreateDirectories?: boolean;
    /**
     * - Indicates if hidden files should be shown.
     */
    ShowHiddenFiles?: boolean;
    /**
     * - Indicates if aliases should be resolved.
     */
    ResolvesAliases?: boolean;
    /**
     * - Indicates if multiple selection is allowed.
     */
    AllowsMultipleSelection?: boolean;
    /**
     * - Indicates if the extension should be hidden.
     */
    HideExtension?: boolean;
    /**
     * - Indicates if hidden extensions can be selected.
     */
    CanSelectHiddenExtension?: boolean;
    /**
     * - Indicates if file packages should be treated as directories.
     */
    TreatsFilePackagesAsDirectories?: boolean;
    /**
     * - Indicates if other file types are allowed.
     */
    AllowsOtherFiletypes?: boolean;
    /**
     * - Array of file filters.
     */
    Filters?: FileFilter[];
    /**
     * - Title of the dialog.
     */
    Title?: string;
    /**
     * - Message to show in the dialog.
     */
    Message?: string;
    /**
     * - Text to display on the button.
     */
    ButtonText?: string;
    /**
     * - Directory to open in the dialog.
     */
    Directory?: string;
    /**
     * - Indicates if the dialog should appear detached from the main window.
     */
    Detached?: boolean;
};
export type MessageDialogOptions = {
    /**
     * - The title of the dialog window.
     */
    Title?: string;
    /**
     * - The main message to show in the dialog.
     */
    Message?: string;
    /**
     * - Array of button options to show in the dialog.
     */
    Buttons?: Button[];
    /**
     * - True if the dialog should appear detached from the main window (if applicable).
     */
    Detached?: boolean;
};
export type Button = {
    /**
     * - Text that appears within the button.
     */
    Label?: string;
    /**
     * - True if the button should cancel an operation when clicked.
     */
    IsCancel?: boolean;
    /**
     * - True if the button should be the default action when the user presses enter.
     */
    IsDefault?: boolean;
};
export type FileFilter = {
    /**
     * - Display name for the filter, it could be "Text Files", "Images" etc.
     */
    DisplayName?: string;
    /**
     * - Pattern to match for the filter, e.g. "*.txt;*.md" for text markdown files.
     */
    Pattern?: string;
};
