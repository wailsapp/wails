import { Log } from '@wails/runtime';

function setLogLevel() {
  Log.SetLogLevel(Log.Level.TRACE);
  // Log.SetLogLevel(Log.Level.DEBUG);
  // Log.SetLogLevel(Log.Level.INFO);
  // Log.SetLogLevel(Log.Level.WARNING);
  // Log.SetLogLevel(Log.Level.ERROR);
}