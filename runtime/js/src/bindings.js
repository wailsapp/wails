/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/

import { Call } from './calls';

var bindingsBasePath = window.backend;

// Determines if the given identifier is valid Javascript
function isValidIdentifier(name) {
  // Don't xss yourself :-)
  try {
    new Function('var ' + name);
    return true;
  } catch (e) {
    return false;
  }
}

// Creates the path given in the bindings path
function addBindingPath(pathSections) {
  // Start at the base path
  var currentPath = bindingsBasePath;
  // for each section of the given path
  for (var sectionIndex in pathSections) {

    var section = pathSections[sectionIndex];

    // Is section a valid javascript identifier?
    if (!isValidIdentifier(section)) {
      var errMessage = section + ' is not a valid javascript identifier.';
      var err = new Error(errMessage);
      return [null, err];
    }

    // Add if doesn't exist
    if (!currentPath[section]) {
      currentPath[section] = {};
    }
    // update current path to new path
    currentPath = currentPath[section];
  }
  return [currentPath, null];
}

export function NewBinding(bindingName) {

  // Get all the sections of the binding
  var bindingSections = bindingName.split('.').splice(1);

  // Get the actual function/method call name
  var callName = bindingSections.pop();

  // Add path to binding
  var bs = addBindingPath(bindingSections);
  var pathToBinding = bs[0];
  var err = bs[1];

  if (err != null) {
    // We need to return an error
    return err;
  }

  // Add binding call
  pathToBinding[callName] = function () {

    // No timeout by default
    var timeout = 0;

    // Actual function
    function dynamic() {
      var args = [].slice.call(arguments);
      return Call(bindingName, args, timeout);
    }

    // Allow setting timeout to function
    dynamic.setTimeout = function (newTimeout) {
      timeout = newTimeout;
    };

    // Allow getting timeout to function
    dynamic.getTimeout = function () {
      return timeout;
    };

    return dynamic;
  }();
}
