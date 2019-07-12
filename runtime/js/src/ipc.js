/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/

// var Invoke = window.external.invoke;

var Invoke;

if (window && window.external && window.external.invoke) {
  Invoke = window.external.invoke;
} else {
  Invoke = console.log;
}

export function SendMessage(type, payload, callbackID) {
  const message = {
    type,
    callbackID,
    payload
  };
  Invoke(JSON.stringify(message));
}
