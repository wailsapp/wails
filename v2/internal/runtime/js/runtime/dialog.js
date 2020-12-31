/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 6 */

/**
 * @type {Object} OpenDialog
 * @param {string} [DefaultDirectory=""]           
 * @param {string} [DefaultFilename=""]            
 * @param {string} [Title=""]                      
 * @param {string} [Filters=""]                    
 * @param {boolean} [AllowFiles=false]
 * @param {boolean} [AllowDirectories=false]
 * @param {boolean} [AllowMultiple=false]
 * @param {boolean} [ShowHiddenFiles=false]
 * @param {boolean} [CanCreateDirectories=false]
 * @param {boolean} [ResolvesAliases=false] - Mac Only: Resolves aliases (symlinks)
 * @param {boolean} [TreatPackagesAsDirectories=false] - Mac Only: Show packages (EG Applications) as folders
 */

/**
 * Opens a dialog using the given parameters, prompting the user to
 * select files/folders.
 *
 * @export
 * @param {OpenDialogOptions} options
 * @returns {Promise<Array<string>>} - List of files/folders selected
 */
export function Open(options) {
	return window.wails.Dialog.Open(options);
}

/**
 * 
 * @type {Object} SaveDialogOptions 
 * @param {string} [DefaultDirectory=""]           
 * @param {string} [DefaultFilename=""]            
 * @param {string} [Title=""]                      
 * @param {string} [Filters=""]                    
 * @param {boolean} [ShowHiddenFiles=false]
 * @param {boolean} [CanCreateDirectories=false]
 * @param {boolean} [TreatPackagesAsDirectories=false]
 */

/**
 * Opens a dialog using the given parameters, prompting the user to
 * select a single file/folder.
 * 
 * @export
 * @param {SaveDialogOptions} options
 * @returns {Promise<string>} 
 */
export function Save(options) {
	return window.wails.Dialog.Save(options);
}

/**
 *
 * @type {Object} MessageDialogOptions
 * @param {DialogType} [Type=InfoDialog] - The type of the dialog
 * @param {string} [Title=""] - The dialog title
 * @param {string} [Message=""] - The dialog message
 * @param {string[]} [Buttons=[]] - The button titles
 * @param {string} [DefaultButton=""] - The button that should be used as the default button
 * @param {string} [CancelButton=""] - The button that should be used as the cancel button
 * @param {string} [Icon=""] - The name of the icon to use in the dialog
 */

/**
 * Opens a dialog using the given parameters, to display a message
 * or prompt the user to select an option
 *
 * @export
 * @param {MessageDialogOptions} options
 * @returns {Promise<string>} - The button text that was selected
 */
export function Message(options) {
	return window.wails.Dialog.Message(options);
}
