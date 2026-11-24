"""Verify the OpenHarmony shell template on a device.

The template is the productised shell: it has no lifecycle probes and no
keeper ability, so this checks the three things a Wails application needs:
the local assets load, a bound method returns through the real Wails runtime,
and Go events reach the page.

    python verify_template.py --hap <entry-default-unsigned.hap>
"""

import argparse
import os
import subprocess
import sys
import time


BUNDLE = "com.wails.nativehello"
TAG = "Wails"


def shell(*args, timeout=60):
    return subprocess.run(
        ["hdc", "-t", "127.0.0.1:5555", "shell", *args],
        check=True, timeout=timeout, text=True, encoding="utf-8",
        errors="replace", stdout=subprocess.PIPE,
    ).stdout


def install(hap):
    # hdc on Windows mis-handles a forward-slash absolute path (C:/...) by
    # appending it to the cwd. Change into the HAP's directory and pass only
    # the file name, which works regardless of the drive letter.
    hap = os.path.abspath(hap)
    subprocess.run(
        ["hdc", "-t", "127.0.0.1:5555", "install", "-r", os.path.basename(hap)],
        check=True, timeout=600, text=True, encoding="utf-8", errors="replace",
        cwd=os.path.dirname(hap),
    )


def start_once():
    started = shell("aa", "start", "-b", BUNDLE, "-a", "EntryAbility")
    if "start ability successfully" not in started:
        raise RuntimeError(f"unable to start the application: {started.strip()}")


def wait_for_pass(previous, timeout):
    """Wait for one process to load the page and report success.

    Returns (pid, ready, passed). A process that stays up but never emits any
    Wails log is the known cold-start freeze: the caller retries it.
    """
    deadline = time.monotonic() + timeout
    pid = None
    ready = False
    passed = False
    while time.monotonic() < deadline:
        pids = shell("pidof", BUNDLE).split()
        if pid is None:
            if len(pids) == 1 and pids[0].isdecimal() and pids[0] not in previous:
                pid = pids[0]
        elif pids != [pid]:
            return pid, ready, passed
        if pid:
            logs = shell("hilog", "-x", "-P", pid, "-T", TAG)
            for line in logs.splitlines():
                if f"{TAG}: " not in line:
                    continue
                print(line, flush=True)
                msg = line.split(f"{TAG}: ", 1)[1]
                if msg == "Local-origin frontend loaded":
                    ready = True
                if msg.startswith("M3 PASS"):
                    passed = True
            if passed:
                return pid, ready, passed
        time.sleep(0.5)
    return pid, ready, passed


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--hap", required=True)
    parser.add_argument("--timeout", type=int, default=60)
    parser.add_argument("--attempts", type=int, default=3)
    args = parser.parse_args()

    install(args.hap)

    for attempt in range(1, args.attempts + 1):
        previous = shell("pidof", BUNDLE).split()
        if previous:
            shell("aa", "force-stop", BUNDLE)
            # The ability needs a moment to fully tear down; starting too
            # quickly leaves a live process whose ability never runs onCreate.
            time.sleep(3)
        start_once()
        pid, ready, passed = wait_for_pass(previous, args.timeout)
        if passed:
            print(f"TEMPLATE PASS pid={pid} attempt={attempt}")
            return
        if pid:
            shell("aa", "force-stop", BUNDLE)
            time.sleep(3)
        print(f"attempt {attempt} did not pass (pid={pid}, ready={ready}); retrying",
              file=sys.stderr, flush=True)

    raise RuntimeError(f"no attempt passed after {args.attempts} tries")


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:  # noqa: BLE001
        print(f"TEMPLATE FAIL: {exc}", file=sys.stderr)
        sys.exit(1)
