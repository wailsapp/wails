/**
 * Returns the selected directory.
 *
 * @export
 * @return {Promise<string>} The selected directory
 */
export function OpenDirectoryDialog(dialogOptions) {
  return window.WailsInvoke(":wails:OpenDirectoryDialog:" + JSON.stringify(dialogOptions))
}

export function OpenMultipleDirectoriesDialog(dialogOptions) {
  return window.WailsInvoke(":wails:OpenMultipleDirectoriesDialog:" + JSON.stringify(dialogOptions))
}

export function OpenFileDialog(dialogOptions) {
  return window.WailsInvoke(":wails:OpenFileDialog:" + JSON.stringify(dialogOptions))
}

export function OpenMultipleFilesDialog(dialogOptions) {
  return window.WailsInvoke(":wails:OpenMultipleFilesDialog:" + JSON.stringify(dialogOptions))
}

export function SaveFileDialog(dialogOptions) {
  return window.WailsInvoke(":wails:SaveFileDialog:" + JSON.stringify(dialogOptions))
}

export function MessageDialog(dialogOptions) {
  return window.WailsInvoke(":wails:MessageDialog:" + JSON.stringify(dialogOptions))
}
