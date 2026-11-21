"""Build, install, and run both OpenHarmony x86_64 lifecycle checks."""

import argparse
import os
from pathlib import Path
import re
import shutil
import subprocess
import time


ROOT = Path(__file__).resolve().parent
HARMONY = ROOT / "harmony"
CPP = HARMONY / "entry/src/main/cpp"
BUNDLE = "com.wails.nativehello"


def run(args, *, capture=False, timeout=120, **kwargs):
    return subprocess.run(
        args, check=True, timeout=timeout, text=True, encoding="utf-8",
        errors="replace", stdout=subprocess.PIPE if capture else None, **kwargs,
    ).stdout


def build(go, hvigor):
    if not os.environ.get("CC"):
        raise RuntimeError("Set CC to OHOS Clang with x86_64 target and SDK sysroot")
    env = dict(os.environ, GOOS="openharmony", GOARCH="amd64", CGO_ENABLED="1")
    lib = CPP / "libs/x86_64/libhomecast_go.so"
    lib.parent.mkdir(parents=True, exist_ok=True)
    header = CPP / "include/libhomecast_go.h"
    header.parent.mkdir(parents=True, exist_ok=True)
    run([go, "build", "-mod=mod", "-buildmode=c-shared",
         "-ldflags=-extldflags=-Wl,--no-undefined", "-o", str(lib), "."],
        cwd=ROOT / "go", env=env)
    shutil.copyfile(lib.with_suffix(".h"), header)
    run([hvigor, "assembleHap", "--no-daemon"], cwd=HARMONY)
    return HARMONY / "entry/build/default/outputs/default/entry-default-unsigned.hap"


def check(shell, mode, timeout):
    previous = shell("pidof", BUNDLE).split()
    if previous:
        stopped = shell("aa", "force-stop", BUNDLE)
        if "force stop process successfully" not in stopped:
            raise RuntimeError(stopped.strip())
    started = shell("aa", "start", "-b", BUNDLE, "-a", "EntryAbility",
                    "--ps", "wailsLifecycle", mode)
    if "start ability successfully" not in started:
        raise RuntimeError(started.strip())

    deadline = time.monotonic() + timeout
    pid = None
    seen = set()
    phase = {}
    ready = False
    while time.monotonic() < deadline:
        pids = shell("pidof", BUNDLE).split()
        if pid is None:
            if len(pids) != 1 or not pids[0].isdecimal() or pids[0] in previous:
                time.sleep(0.5)
                continue
            pid = pids[0]
        elif pids != [pid]:
            raise RuntimeError(f"{mode}: application process {pid} exited or changed: {pids}")
        logs = shell("hilog", "-x", "-P", pid, "-T", "WailsM3")
        for line in logs.splitlines():
            if line in seen or "WailsM3: " not in line:
                continue
            seen.add(line)
            print(line, flush=True)
            msg = line.split("WailsM3: ", 1)[1]
            ready |= msg == "Local-origin frontend loaded"
            match = re.fullmatch(rf"M4 {mode} (old|new) session=([0-9a-f-]+) id=(\S+)", msg)
            if match:
                phase[match[1]] = match[2], match[3]
            if msg.startswith(f"M4 FAIL {mode}:"):
                raise RuntimeError(msg)
            if msg.startswith(f"M4 PASS {mode} session="):
                old, new = phase.get("old"), phase.get("new")
                if not old or not new or old[0] == new[0] or old[1] != new[1] or msg.rsplit("=", 1)[1] != new[0]:
                    raise RuntimeError(f"{mode}: invalid session evidence: {phase}")
                return
        time.sleep(0.5)
    stage = "process start" if pid is None else "ArkWeb load" if not ready else "lifecycle probe"
    raise RuntimeError(f"{mode}: timed out waiting for {stage} (PID {pid}, phases {list(phase)})")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--go", default="go", help="OpenHarmony-enabled go executable")
    parser.add_argument("--hvigor", default="hvigorw.bat", help="DevEco Hvigor wrapper")
    parser.add_argument("--target", default="127.0.0.1:5555", help="HDC emulator target")
    parser.add_argument("--timeout", type=float, default=45, help="seconds per runtime check")
    args = parser.parse_args()
    hdc = ["hdc", "-t", args.target]

    def shell(*command):
        return run(hdc + ["shell", *command], capture=True, timeout=12)

    try:
        hap = build(args.go, args.hvigor)
        installed = run(hdc + ["install", "-r", str(hap)], capture=True, timeout=60)
        print(installed, end="", flush=True)
        if "install bundle successfully" not in installed:
            raise RuntimeError("HAP installation failed")
        for mode in ("refresh", "rebuild"):
            check(shell, mode, args.timeout)
    except (OSError, subprocess.SubprocessError, RuntimeError) as exc:
        parser.exit(1, f"M4: {exc}\n")


if __name__ == "__main__":
    main()
