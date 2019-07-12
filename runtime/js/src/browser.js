/*
 _       __      _ __    
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  ) 
|__/|__/\__,_/_/_/____/  
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/

import { SystemCall } from './calls';

export function OpenURL(url) {
  return SystemCall('Browser.OpenURL', url);
}

export function OpenFile(filename) {
  return SystemCall('Browser.OpenFile', filename);
}
