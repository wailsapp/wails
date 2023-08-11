import {Call} from "./calls";

/**
 * Returns the selected directory.
 *
 * @export
 * @return {Promise<string>} The selected directory
 */
export function OpenDirectoryDialog() {
  return Call(":wails:OpenDirectoryDialog")
}

export function OpenMultipleDirectoriesDialog() {
  return Call(":wails:OpenMultipleDirectoriesDialog")
}

export function OpenFileDialog() {
  return Call(":wails:OpenFileDialog")
}

export function OpenMultipleFilesDialog() {
  return Call(":wails:OpenMultipleFilesDialog")
}

export function SaveFileDialog() {
  return Call(":wails:SaveFileDialog")
}

export function MessageDialog() {
  return Call(":wails:MessageDialog")
}
