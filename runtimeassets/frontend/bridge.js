/*
    Wails Bridge (c) 2019-present Lea Anthony

    This prod version is to get around having to rewrite your code
    for production. When doing a release build, this file will be used 
    instead of the full version.
*/

export default {
  // The main function
  // Passes the main Wails object to the callback if given.
  Start: function (callback) {
    document.addEventListener('DOMContentLoaded', () => {
      try {
        if (callback) {
          window.wails.events.on("wails:ready", callback);
        }
      }
      catch (e) {
        console.log(e);
      }
    });

  }
};
