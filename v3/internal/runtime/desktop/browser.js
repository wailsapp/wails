/**
 * @description: Use the system default browser to open the url
 * @param {string} url 
 * @return {void}
 */
export function BrowserOpenURL(url) {
  window.WailsInvoke('BO:' + url);
}