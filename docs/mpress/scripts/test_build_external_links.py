"""Exercise the build/audit boundary offline with controlled tool responses."""

import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest


class BuildExternalAuditTests(unittest.TestCase):
    def run_build(self, audit, status=0):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            scripts = root / "docs/mpress/scripts"
            scripts.mkdir(parents=True)
            shutil.copy(Path(__file__).with_name("build.sh"), scripts / "build.sh")
            tools = root / "tools"
            tools.mkdir()
            log = root / "calls.jsonl"
            candidate = tools / "mpress"
            candidate.write_text(f"#!{sys.executable}\n" + '''
import json, os, sys
args = sys.argv[1:]
with open(os.environ["AUDIT_TEST_CALLS"], "a") as stream:
    stream.write(json.dumps(args) + "\\n")
if args[:2] == ["check", "external"]:
    status = int(os.environ["AUDIT_TEST_STATUS"])
    print(json.dumps({"links": [{"url": "https://github.com/example/missing",
                                "status": "broken" if status == 1 else "inconclusive",
                                "pages": ["index.html", "ko/index.html"]}]}))
    sys.exit(status)
''')
            candidate.chmod(0o755)
            # Other validations are tested separately with the actual sources.
            # This test isolates opt-in behavior and status/report preservation.
            python = tools / "python3"
            python.write_text("#!/bin/sh\nexit 0\n")
            python.chmod(0o755)
            env = dict(os.environ, MPRESS_BIN=str(candidate),
                       PATH=str(tools) + os.pathsep + os.environ["PATH"],
                       AUDIT_TEST_CALLS=str(log), AUDIT_TEST_STATUS=str(status))
            report = root / "report.json"
            args = ["bash", str(scripts / "build.sh")]
            if audit:
                args.extend(["--external-audit", str(report)])
            result = subprocess.run(args, env=env, capture_output=True, text=True)
            calls = [json.loads(line) for line in log.read_text().splitlines()]
            self.assertEqual(calls[:2], [["build", "--strict", "--no-purge-css"], ["check"]])
            if audit:
                self.assertEqual(calls[2], ["check", "external", "--host", "github.com",
                    "--url-prefix", "https://github.com/wailsapp/wails/tree/",
                    "--url-prefix", "https://github.com/wailsapp/wails/blob/", "--json"])
                self.assertEqual(json.loads(report.read_text())["links"][0]["pages"],
                                 ["index.html", "ko/index.html"])
            else:
                self.assertEqual(len(calls), 2)
                self.assertFalse(report.exists())
            self.assertEqual(result.returncode, status, result.stderr)

    def test_routine_build_does_not_audit_external_destinations(self):
        self.run_build(False)

    def test_audit_preserves_confirmed_failure_and_report(self):
        self.run_build(True, 1)

    def test_audit_preserves_inconclusive_status_and_report(self):
        self.run_build(True, 2)
