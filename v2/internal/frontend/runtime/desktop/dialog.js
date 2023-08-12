/**
 * Returns the selected directory.
 *
 * @export
 * @return {Promise<string>} The selected directory
 */
export function OpenDirectoryDialog(dialogOptions) {
  return window.WailsInvoke("DOD:" + JSON.stringify(dialogOptions))
}

export function OpenMultipleDirectoriesDialog(dialogOptions) {
  return window.WailsInvoke("DOMD:" + JSON.stringify(dialogOptions))
}

export function OpenFileDialog(dialogOptions) {
  return window.WailsInvoke("DOF:" + JSON.stringify(dialogOptions))
}

export function OpenMultipleFilesDialog(dialogOptions) {
  return window.WailsInvoke("DOMF:" + JSON.stringify(dialogOptions))
}

export function SaveFileDialog(dialogOptions) {
  return window.WailsInvoke("DSF:" + JSON.stringify(dialogOptions))
}

export function MessageDialog(dialogOptions) {
  return window.WailsInvoke("DM:" + JSON.stringify(dialogOptions))
}
