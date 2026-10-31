#!/usr/bin/env python3
"""Render an M-Press external-link report for the Actions job summary."""

import html
import json
from pathlib import Path
import sys


def summarize(path):
    if not path.is_file() or not path.read_text(encoding="utf-8").strip():
        return "External audit did not produce a report. Inspect the build/audit step log.\n"
    report = json.loads(path.read_text(encoding="utf-8"))
    links = report["links"]
    broken = sum(link["status"] == "broken" for link in links)
    inconclusive = sum(link["status"] == "inconclusive" for link in links)
    lines = [f"External audit: {len(links)} destinations; {broken} broken; "
             f"{inconclusive} inconclusive.", "",
             "404/410 confirms a missing destination. Other failures require a retry or investigation.",
             "See the docs-external-links artifact for every referencing page and redirect destination.", ""]
    for link in links:
        if link["status"] == "ok":
            continue
        # Treat source-controlled link/report text as data, including in Markdown.
        lines.append("<details><summary>" + html.escape(link["status"] + ": " + link["url"])
                     + "</summary><pre>" + html.escape(link.get("detail", "") + "\n"
                     + "\n".join(link["pages"])) + "</pre></details>")
    return "\n".join(lines) + "\n"


if __name__ == "__main__":
    print(summarize(Path(sys.argv[1])), end="")
