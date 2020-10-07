
import { Log } from '@wailsapp/runtime2';

function doSomeOperation() {
    // Do things
    let value = doSomething();
    Log.Debug("I got: " + value);
    /*
    Log.Info("An Info message");
    Log.Warning("A Warning message");
    Log.Error("An Error message");
    */
}

function abort() {
    // Do some things
    Log.Fatal("I accidentally the whole application!");
}