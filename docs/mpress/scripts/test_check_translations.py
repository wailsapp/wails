"""The publication gate must reject fallback files and stale audit exceptions."""

import json
from pathlib import Path
from tempfile import TemporaryDirectory
import unittest
from unittest.mock import patch
from types import SimpleNamespace

from check_translations import (AUDIT_EXCLUDED_FILES, accepted_finding,
                                audit_language, check_coverage, file_digest)


class TranslationChecks(unittest.TestCase):
    def setUp(self):
        temp = TemporaryDirectory()
        self.addCleanup(temp.cleanup)
        self.content = Path(temp.name)
        (self.content / "fr").mkdir()
        (self.content / "guide.md").write_text("# Guide\n\nInstall the application.\n")
        (self.content / "_nav.yaml").write_text("- label: Guide\n  link: /guide/\n")

    def test_missing_page_and_navigation_are_rejected(self):
        errors = check_coverage(self.content, ("fr",))
        self.assertEqual(len(errors), 2)
        self.assertTrue(all("missing translation" in error for error in errors))

    def test_empty_and_orphaned_translations_are_rejected(self):
        (self.content / "fr/guide.md").write_text(" \n")
        (self.content / "fr/old.md").write_text("# Ancien\n")
        (self.content / "fr/_nav.yaml").write_text("- label: Guide\n  link: /guide/\n")
        errors = check_coverage(self.content, ("fr",))
        self.assertEqual(len(errors), 2)
        self.assertTrue(any("empty translation" in error for error in errors))
        self.assertTrue(any("no matching English source" in error for error in errors))

    def test_exception_requires_current_source_and_translation(self):
        target = self.content / "fr/guide.md"
        target.write_text("# Guide\n\nInstallez l’application.\n")
        finding = {"file": "guide.md", "segment": "title", "code": "untranslated"}
        exception = {**finding, "language": "fr", "reason": "Guide is also the French title.",
                     "sourceSHA256": file_digest(self.content / "guide.md"),
                     "targetSHA256": file_digest(target)}
        self.assertTrue(accepted_finding(self.content, "fr", finding, [exception]))
        target.write_text("# Guide\n\nInstall the application.\n")
        self.assertFalse(accepted_finding(self.content, "fr", finding, [exception]))

    def test_structure_errors_cannot_be_waived(self):
        finding = {"file": "guide.md", "segment": "body", "code": "structure"}
        self.assertFalse(accepted_finding(self.content, "fr", finding, [finding]))

    def test_historical_changelog_audit_findings_are_excluded(self):
        (self.content / "changelog.md").write_text("# Changelog\n")
        (self.content / "fr/changelog.md").write_text("# Journal\n")
        report = {
            "language": "fr",
            "files": ["changelog.md", "guide.md"],
            "segments": ["changelog-title", "guide-title"],
            "findings": [
                {"file": "changelog.md", "segment": "body", "code": "structure",
                 "message": "legacy baseline"},
                {"file": "guide.md", "segment": "title", "code": "structure",
                 "message": "current issue"},
            ],
        }
        with patch("check_translations.subprocess.run") as run:
            run.return_value = SimpleNamespace(stdout=json.dumps(report), stderr="", returncode=1)
            errors = audit_language("mpress", self.content, self.content, "fr", [])
        self.assertIn("changelog.md", AUDIT_EXCLUDED_FILES)
        self.assertEqual(errors, ["fr/guide.md title: structure: current issue"])


if __name__ == "__main__":
    unittest.main()
