# Launch the configured Android application ID

Type: task
Status: resolved
Blocked by: none
Label: ready-for-agent
Priority: P1

## Question

With android.application_id different from project.identifier, android run installs the correct APK but launches the wrong package. Reproduced on the API 36 x86_64 emulator: Activity class com.mycompany.myproduct/com.wails.app.MainActivity does not exist, while the installed APK belongs to com.wails.gareview.badge.

## Acceptance criteria

Use the Android application ID with the project ID as fallback for both build-and-run and existing-APK deployment. Add a real manifest-loading regression and verify installation, launch and package-scoped attachment on the emulator.

## Comments

2026-09-08: Native Android launch and permanent manifest-loading test reproduce the same ID mismatch.

## Answer

2026-09-08: Fixed Android application-ID selection in both launch paths. Manifest regression passes. The API 36 emulator successfully built/installed/launched the APK, launched an existing APK, streamed logs and cancelled cleanly.
