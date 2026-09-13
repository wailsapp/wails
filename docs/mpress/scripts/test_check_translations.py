"""The publication gate must reject fallback files and stale audit exceptions."""

from pathlib import Path
from tempfile import TemporaryDirectory
import unittest

from check_translations import accepted_finding, check_coverage, file_digest


class TranslationChecks(unittest.TestCase):
    def setUp(self):
        temp = TemporaryDirectory()
        self.addCleanup(temp.cleanup)
        self.content = Path(temp.name)
        (self.content / "fr").mkdir()
        (self.content / "guide.mpd").write_text("# Guide\n\nInstall the application.\n")
        (self.content / "_nav.yaml").write_text("- label: Guide\n  link: /guide/\n")

    def test_missing_page_and_navigation_are_rejected(self):
        errors = check_coverage(self.content, ("fr",))
        self.assertEqual(len(errors), 2)
        self.assertTrue(all("missing translation" in error for error in errors))

    def test_empty_and_orphaned_translations_are_rejected(self):
        (self.content / "fr/guide.mpd").write_text(" \n")
        (self.content / "fr/old.mpd").write_text("# Ancien\n")
        (self.content / "fr/_nav.yaml").write_text("- label: Guide\n  link: /guide/\n")
        errors = check_coverage(self.content, ("fr",))
        self.assertEqual(len(errors), 2)
        self.assertTrue(any("empty translation" in error for error in errors))
        self.assertTrue(any("no matching English source" in error for error in errors))

    def test_exception_requires_current_source_and_translation(self):
        target = self.content / "fr/guide.mpd"
        target.write_text("# Guide\n\nInstallez l’application.\n")
        finding = {"file": "guide.mpd", "segment": "title", "code": "untranslated"}
        exception = {**finding, "language": "fr", "reason": "Guide is also the French title.",
                     "sourceSHA256": file_digest(self.content / "guide.mpd"),
                     "targetSHA256": file_digest(target)}
        self.assertTrue(accepted_finding(self.content, "fr", finding, [exception]))
        target.write_text("# Guide\n\nInstall the application.\n")
        self.assertFalse(accepted_finding(self.content, "fr", finding, [exception]))

    def test_structure_errors_cannot_be_waived(self):
        finding = {"file": "guide.mpd", "segment": "body", "code": "structure"}
        self.assertFalse(accepted_finding(self.content, "fr", finding, [finding]))


if __name__ == "__main__":
    unittest.main()
