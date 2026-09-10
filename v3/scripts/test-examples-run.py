#!/usr/bin/env python3
"""Build and smoke-launch every example manifest on the current host.

The report distinguishes launch smoke checks from interactive feature acceptance.
Logs and reports go to the caller's output directory, outside the examples.
"""
import argparse
import json
import os
from pathlib import Path
import signal
import socket
import subprocess
import threading
import time


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--cli", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--example", action="append", help="Relative example path; repeat to narrow a rerun")
    parser.add_argument("--launch-seconds", type=float, default=4)
    parser.add_argument("--plan-timeout", type=float, default=30)
    parser.add_argument("--build-timeout", type=float, default=240)
    options = parser.parse_args()
    if options.launch_seconds <= 0 or options.build_timeout <= 0 or options.plan_timeout <= 0:
        parser.error("timeouts must be positive")
    cli = options.cli.resolve()
    examples = Path(__file__).resolve().parents[1] / "examples"
    output = options.output.resolve()
    output.mkdir(parents=True, exist_ok=True)
    base_env = dict(os.environ, WAILS_EXP_USE_WAKE="1", NO_COLOR="1")
    results = []
    for manifest in sorted(examples.rglob("wails.hcl")):
        if any(part in ("node_modules", ".wails") for part in manifest.relative_to(examples).parts):
            continue
        name = manifest.parent.relative_to(examples).as_posix()
        if options.example and name not in options.example:
            continue
        row = {"example": name}
        env = dict(base_env)
        if name == "server":
            # Avoid colliding with an unrelated local service on the example's default port.
            with socket.socket() as listener:
                listener.bind(("127.0.0.1", 0))
                env["WAILS_SERVER_PORT"] = str(listener.getsockname()[1])
        try:
            plan = subprocess.run([str(cli), "run", "--plan"], cwd=manifest.parent,
                                  env=env, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                                  text=True, timeout=options.plan_timeout)
        except subprocess.TimeoutExpired as error:
            log = error.stdout or b""
            if isinstance(log, bytes):
                log = log.decode(errors="replace")
            (output / (name.replace("/", "_") + ".plan.log")).write_text(log)
            row["status"] = "plan-timeout"
            plan = None
        if plan is not None:
            (output / (name.replace("/", "_") + ".plan.log")).write_text(plan.stdout)
        if plan is None:
            pass  # Record this failure below, then continue the remaining examples.
        elif plan.returncode:
            row["status"] = "unsupported-host" if "does not support" in plan.stdout else "plan-failed"
            row["exit"] = plan.returncode
        else:
            lines = []
            built = threading.Event()
            process = subprocess.Popen([str(cli), "run"], cwd=manifest.parent, env=env,
                                       stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                                       text=True,
                                       creationflags=subprocess.CREATE_NEW_PROCESS_GROUP if os.name == "nt" else 0)
            def read_output():
                for line in process.stdout:
                    lines.append(line)
                    if "build succeeded" in line or "run succeeded" in line:
                        built.set()
            reader = threading.Thread(target=read_output, daemon=True)
            reader.start()
            deadline = time.monotonic() + options.build_timeout
            while process.poll() is None and not built.wait(0.1) and time.monotonic() < deadline:
                pass
            if built.is_set():
                try:
                    process.wait(timeout=options.launch_seconds)
                    row["status"] = "exited-cleanly" if process.returncode == 0 else "launch-failed"
                except subprocess.TimeoutExpired:
                    row["status"] = "launch-smoke"
            else:
                row["status"] = "build-failed" if process.poll() is not None else "build-timeout"
            if process.poll() is None:
                process.send_signal(signal.CTRL_BREAK_EVENT if os.name == "nt" else signal.SIGINT)
                try:
                    process.wait(timeout=10)
                except subprocess.TimeoutExpired:
                    process.kill()
                    process.wait()
                    row["cleanup"] = "forced-cli-kill"
                    row["status"] = "cleanup-failed"
            reader.join(timeout=5)
            if reader.is_alive():
                row["status"] = "cleanup-failed"
            log = "".join(lines)
            (output / (name.replace("/", "_") + ".run.log")).write_text(log)
            row["exit"] = process.returncode
            if row["status"] == "launch-smoke" and any(token in log for token in ("panic:", "fatal error:", "SIGSEGV")):
                row["status"] = "runtime-failed"
        results.append(row)
        (output / "results.json").write_text(json.dumps(results, indent=2) + "\n")
        print(name, row["status"], flush=True)
    missing = set(options.example or []) - {row["example"] for row in results}
    if missing:
        parser.error("unknown examples: " + ", ".join(sorted(missing)))
    if not results:
        parser.error("no example manifests found")
    failures = [row for row in results if row["status"] not in ("unsupported-host", "launch-smoke", "exited-cleanly")]
    print(f"{len(results)} examples checked; {len(failures)} failures", flush=True)
    return bool(failures)


if __name__ == "__main__":
    raise SystemExit(main())
