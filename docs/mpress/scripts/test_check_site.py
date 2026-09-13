"""Regression tests for failures a successful static build can miss."""

from pathlib import Path
from tempfile import TemporaryDirectory
import unittest

from check_site import Page, check_links, check_page


class SiteChecks(unittest.TestCase):
    def setUp(self):
        self.temp = TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.repository = Path(self.temp.name).resolve()
        self.site = self.repository / "site"
        self.site.mkdir()

    def page(self, name, body):
        path = self.site / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(body, encoding="utf-8")
        return Page(path)

    def test_absolute_internal_link_is_checked(self):
        page = self.page("index.html", '<a href="https://v3.wails.io/missing/">broken</a>')
        errors = []
        check_links(page, self.site, {page.path: page}, {}, errors)
        self.assertEqual(len(errors), 1)
        self.assertIn("missing link", errors[0])

    def test_redirects_and_encoded_fragments_resolve(self):
        page = self.page("index.html", '<a href="/old/#caf%C3%A9">valid</a>')
        target = self.page("new/index.html", '<h2 id="café">Heading</h2>')
        errors = []
        check_links(page, self.site, {target.path: target}, {"/old/": "/new/"}, errors)
        self.assertEqual(errors, [])

    def test_missing_fragment_is_rejected(self):
        page = self.page("index.html", '<a href="#missing">broken</a>')
        errors = []
        check_links(page, self.site, {page.path: page}, {}, errors)
        self.assertIn("missing anchor", errors[0])

    def test_contribution_path_cannot_escape_content(self):
        (self.repository / "private.txt").write_text("outside content")
        page = self.page("index.html", '<html lang="en">'
                         '<link rel="canonical" href="https://v3.wails.io/">'
                         '<meta name="mpress:source" content="private.txt">')
        errors = []
        check_page(page, self.site, self.repository, errors)
        self.assertEqual(len(errors), 1)
        self.assertIn("invalid contribution source", errors[0])

    def test_missing_language_and_wrong_canonical_are_rejected(self):
        page = self.page("index.html", '<link rel="canonical" href="https://old.example/">')
        errors = []
        check_page(page, self.site, self.repository, errors)
        self.assertEqual(len(errors), 2)

    def test_unrendered_d2_is_rejected(self):
        page = self.page("index.html", '<html lang="en">'
                         '<link rel="canonical" href="https://v3.wails.io/">'
                         '<pre><code class="language-d2">a -&gt; b</code></pre>')
        errors = []
        check_page(page, self.site, self.repository, errors)
        self.assertEqual(errors, ["index.html: D2 diagram was displayed as source code"])


if __name__ == "__main__":
    unittest.main()
