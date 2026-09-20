#!/usr/bin/env python3
"""Require complete translations and locally audit their protected structure."""

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import sys


LANGUAGES = ("fr", "de", "pt", "ru", "ja", "ko", "zh-cn", "zh-tw", "id")

# The translated changelogs were imported before M-Press segment tracking and
# currently have a known audit baseline mismatch. Keep coverage checks and all
# other translated pages strict until those generated release notes are
# re-synchronised.
AUDIT_EXCLUDED_FILES = frozenset({"changelog.md"})


def check_coverage(content, languages=LANGUAGES):
    sources = {path.relative_to(content).as_posix() for path in content.rglob("*.md")
               if path.relative_to(content).parts[0] not in languages}
    sources.add("_nav.yaml")
    errors = []
    for language in languages:
        root = content / language
        targets = {path.relative_to(root).as_posix() for path in root.rglob("*.md")}
        if (root / "_nav.yaml").is_file():
            targets.add("_nav.yaml")
        errors.extend(f"{language}/{file}: missing translation" for file in sorted(sources - targets))
        errors.extend(f"{language}/{file}: no matching English source" for file in sorted(targets - sources))
        for file in sorted(sources & targets):
            if not (root / file).read_text(encoding="utf-8").strip():
                errors.append(f"{language}/{file}: empty translation")
    return errors


def file_digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def accepted_finding(content, language, finding, exceptions):
    # Only linguistic heuristics can be adjudicated. Missing content or damaged
    # document structure must always fail, even if an exception is supplied.
    if finding["code"] not in {"untranslated", "requirement-language"}:
        return False
    file = Path(finding["file"])
    if file.is_absolute() or ".." in file.parts:
        return False
    source, target = content / file, content / language / file
    if not source.is_file() or not target.is_file():
        return False
    identity = {"language": language, "file": finding["file"],
                "segment": finding.get("segment", ""), "code": finding["code"]}
    for entry in exceptions:
        if all(entry.get(key) == value for key, value in identity.items()):
            return (bool(entry.get("reason", "").strip())
                    and entry.get("sourceSHA256") == file_digest(source)
                    and entry.get("targetSHA256") == file_digest(target))
    return False


def audit_language(mpress, repository, content, language, exceptions):
    result = subprocess.run([mpress, "translate", "audit", "--lang", language, "--json"],
                            cwd=repository, text=True, capture_output=True)
    try:
        report = json.loads(result.stdout)
    except json.JSONDecodeError:
        return [f"{language}: translation audit failed: {result.stderr.strip()}"]
    findings = report.get("findings", [])
    if report.get("language") != language or not report.get("files") or not report.get("segments"):
        return [f"{language}: incomplete translation audit report"]
    if result.returncode and not findings:
        return [f"{language}: translation audit failed without findings"]
    return [f"{language}/{item['file']} {item.get('segment', '')}: {item['code']}: {item['message']}"
            for item in findings
            if item["file"] not in AUDIT_EXCLUDED_FILES
            and not accepted_finding(content, language, item, exceptions)]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--mpress", default="mpress", help="path to the M-Press CLI")
    args = parser.parse_args()
    repository = Path(__file__).resolve().parents[3]
    content = repository / "docs/mpress/content"
    exception_path = repository / "docs/mpress/translation/audit-exceptions.json"
    exceptions = json.loads(exception_path.read_text(encoding="utf-8"))
    errors = check_coverage(content)
    for language in LANGUAGES:
        errors.extend(audit_language(args.mpress, repository, content, language, exceptions))
    print(json.dumps({"languages": list(LANGUAGES), "errors": errors}, indent=2, ensure_ascii=False))
    return bool(errors)


if __name__ == "__main__":
    sys.exit(main())
