# OpenHarmony native bridge device probe

This standalone HAP checks whether the OpenHarmony host can submit work to a Go goroutine and receive results on the ArkTS thread via a NAPI thread-safe function. It neither imports Wails nor opens a network port. Passing this probe is not evidence that the Wails bridge or ArkWeb integration works.

## Build and run (Windows, OpenHarmony x86_64)

1. Use a Go OpenHarmony toolchain (verified with `go1.27.4-hmos-devel`) and DevEco Studio with OpenHarmony SDK API 26. Set `GOOS=openharmony`, `GOARCH=amd64`, `CGO_ENABLED=1` and `CC=<SDK>/native/llvm/bin/clang.exe --target=x86_64-linux-ohos --sysroot=<SDK>/native/sysroot`. Use a short path if the compiler path contains spaces.
2. From `go/`, run `go build -buildmode=c-shared -ldflags=-extldflags=-Wl,--no-undefined -o libhomecast_go.so .`. Copy the generated `.so` into `harmony/entry/src/main/cpp/libs/x86_64/` and the generated `libhomecast_go.h` into `harmony/entry/src/main/cpp/include/`. The library name is retained from the working shell template; this package does not contain HomeCast code.
3. From `harmony/`, configure the local SDK in DevEco Studio (or set `sdk.dir` in an untracked `local.properties` to the DevEco SDK default parent directory) and run `hvigorw.bat assembleHap --no-daemon` using DevEco Studio's Hvigor. The verified project uses compile, compatible and target SDK `26.0.0`.
4. On the unlocked x86_64 emulator, run `hdc install entry/build/default/outputs/default/entry-default-unsigned.hap`, then `hdc shell aa start -b com.wails.nativebridgeprobe -a EntryAbility`. Use `hdc -t <emulator>` if multiple devices are connected. `hdc shell hilog -x -T WailsM2` should show `M2 PASS` from the new process.

`M2 PASS` checks three correct correlated replies, a pending request at the deadline, suppression of both late replies after closing the session, and a responsive ArkTS timer. For cold restart, force-stop only `com.wails.nativebridgeprobe`, launch again, and verify a fresh `M2 PASS` from the new PID. The probe does not test in-process UIAbility recreation or long-running bridge load.
