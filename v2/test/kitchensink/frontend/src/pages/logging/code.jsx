import { Log } from '@wails/runtime';

function doSomeOperation() {
  // Do things
  let value = doSomething();
  Log.Trace("I got: " + value);
  Log.Debug("A debug message");
  Log.Info("An Info message");
  Log.Warning("A Warning message");
  Log.Error("An Error message");
}

function abort() {
  // Do some things
  Log.Fatal("I accidentally the whole application!");
}