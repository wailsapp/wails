#!/usr/bin/env python3
"""Audit reproduction: expected to fail until compiler cache inputs include .cxx."""
import json
import pathlib
import subprocess
import tempfile

here = pathlib.Path(__file__).resolve().parent
repo = here.parent.parent
with tempfile.TemporaryDirectory(prefix="wails-ripwire-cxx-") as directory:
    overlay = pathlib.Path(directory) / "overlay.json"
    virtual = repo / "v3/internal/wake/pipeline/ripwire_cxx_regression_test.go"
    if virtual.exists():
        raise SystemExit(f"Refusing to shadow {virtual}")
    overlay.write_text(json.dumps({"Replace": {str(virtual): str(here / "cxx-regression_test.go")}}))
    result = subprocess.run(["go", "test", "-overlay=" + str(overlay), "./internal/wake/pipeline", "-run", "TestReviewCompileInputsTrackCXX", "-count=1"], cwd=repo / "v3")
    raise SystemExit(result.returncode)
