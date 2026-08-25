# Wails v3 Android Example

This example runs on Android (emulator and device) as well as desktop. See
[`ANDROID.md`](../../ANDROID.md) for the full Android guide.

```bash
wails3 task android:run        # build + launch in the Android Emulator
wails3 task android:package    # production release APK
wails3 task android:logs       # stream logcat output
```

It demonstrates service bindings, Go->JS events, native AlertDialog message
dialogs, clipboard, device info, and screen metrics.
