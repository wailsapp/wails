import json
from pathlib import Path
import tempfile
import unittest

from summarize_external_links import summarize


class ExternalSummaryTests(unittest.TestCase):
    def test_report_lists_all_affected_pages_and_separates_uncertainty(self):
        with tempfile.TemporaryDirectory() as temp:
            path = Path(temp) / "report.json"
            path.write_text(json.dumps({"links": [
                {"status": "ok", "url": "https://github.com/good", "pages": ["ok/index.html"]},
                {"status": "broken", "url": "https://github.com/<missing>",
                 "detail": "404 Not Found", "pages": ["en/index.html", "ko/index.html"]},
                {"status": "inconclusive", "url": "https://github.com/limited",
                 "detail": "429 Too Many Requests", "pages": ["fr/index.html"]}]}))
            result = summarize(path)
            self.assertIn("3 destinations; 1 broken; 1 inconclusive", result)
            for page in ("en/index.html", "ko/index.html", "fr/index.html"):
                self.assertIn(page, result)
            self.assertIn("&lt;missing&gt;", result)
            self.assertNotIn("https://github.com/good", result)

    def test_missing_or_empty_report_is_visible(self):
        with tempfile.TemporaryDirectory() as temp:
            path = Path(temp) / "report.json"
            self.assertIn("did not produce", summarize(path))
            path.touch()
            self.assertIn("did not produce", summarize(path))
