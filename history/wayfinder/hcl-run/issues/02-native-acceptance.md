# Complete native run acceptance on remaining platforms

Type: task
Status: open
Blocked by: none
Labels: ready-for-human

## Question

Run the platform acceptance commands in the
[verification handoff](../verification.md) on Windows, macOS, an Android device
or working emulator, and an iOS simulator and signed physical device. Verify
application arguments, environment handling, shutdown, signing and bundle identity.
Exercise each example's interactive features beyond the automated launch smoke.

## Comments

2026-09-10: Linux launch smoke covers all 105 example entry directories: 100
launches and five explicit unsupported-platform rejections. Android APK assembly
and signature verification pass. The installed Android emulator crashes with exit
139 before stable ADB readiness in three graphics/acceleration configurations;
no Wails application was installed or launched on it. Windows/macOS CLI cross-
compilation passes but does not prove native launch behaviour. GTK3 development
libraries are absent; the build-tag selection test covers both GTK stacks, while
native compilation and launch were verified only with GTK4 / WebKitGTK 6.0.
