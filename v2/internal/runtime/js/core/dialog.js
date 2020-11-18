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

import { SystemCall } from './calls';

/**
 * @type {Object} OpenDialog
 * @param {string} [DefaultDirectory=""]           
 * @param {string} [DefaultFilename=""]            
 * @param {string} [Title=""]                      
 * @param {string} [Filters=""]                    
 * @param {bool} [AllowFiles=false]                 
 * @param {bool} [AllowDirectories=false]           
 * @param {bool} [AllowMultiple=false]              
 * @param {bool} [ShowHiddenFiles=false]            
 * @param {bool} [CanCreateDirectories=false]       
 * @param {bool} [ResolvesAliases=false] - Mac Only: Resolves aliases (symlinks)            
 * @param {bool} [TreatPackagesAsDirectories=false] - Mac Only: Show packages (EG Applications) as folders
 */



/**
 * Opens a dialog using the given paramaters, prompting the user to 
 * select files/folders.
 *
 * @export
 * @param {OpenDialogOptions} options
 * @returns {Promise<Array<string>>} - List of files/folders selected
 */
export function Open(options) {
	return SystemCall('Dialog.Open', options);
}

/**
 * 
 * @type {Object} SaveDialogOptions 
 * @param {string} [DefaultDirectory=""]           
 * @param {string} [DefaultFilename=""]            
 * @param {string} [Title=""]                      
 * @param {string} [Filters=""]                    
 * @param {bool} [ShowHiddenFiles=false]            
 * @param {bool} [CanCreateDirectories=false]       
 * @param {bool} [TreatPackagesAsDirectories=false] 
 */

/**
 * Opens a dialog using the given paramaters, prompting the user to 
 * select a single file/folder.
 * 
 * @export
 * @param {SaveDialogOptions} options
 * @returns {Promise<string>} 
 */
export function Save(options) {
	return SystemCall('Dialog.Save', options);
}
