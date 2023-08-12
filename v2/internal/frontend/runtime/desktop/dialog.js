import {Call} from "./calls";

/**
 * Returns the selected directory.
 *
 * @export
 * @return {Promise<string>} The selected directory
 */
export function OpenDirectoryDialog(dialogOptions) {
  return Call(":wails:OpenDirectoryDialog", [dialogOptions])
}

export function OpenMultipleDirectoriesDialog(dialogOptions) {
  return Call(":wails:OpenMultipleDirectoriesDialog", [dialogOptions])
}

export function OpenFileDialog(dialogOptions) {
  return Call(":wails:OpenFileDialog", [dialogOptions])
}

export function OpenMultipleFilesDialog(dialogOptions) {
  return Call(":wails:OpenMultipleFilesDialog", [dialogOptions])
}

export function SaveFileDialog(dialogOptions) {
  return Call(":wails:SaveFileDialog", [dialogOptions])
}

export function MessageDialog(dialogOptions) {
  return Call(":wails:MessageDialog", [dialogOptions])
}
