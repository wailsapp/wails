"""Regression tests for failures a successful static build can miss."""

import json
from pathlib import Path
from tempfile import TemporaryDirectory
import unittest

from check_site import Page, check_links, check_page, check_translation_routes, check_knowledge_assets


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

    def test_raw_html_cannot_hide_unrendered_markdown_code(self):
        page = self.page("index.html", '<html lang="en">'
                         '<link rel="canonical" href="https://v3.wails.io/">'
                         '<details><summary>Example</summary>```go\nfunc main() {}\n```</details>')
        errors = []
        check_page(page, self.site, self.repository, errors)
        self.assertEqual(errors, ["index.html: Markdown code fence was displayed as prose"])
        literal = self.page("literal.html", '<pre><code>```go</code></pre><script>"```"</script>')
        self.assertFalse(literal.unrendered_fence)

    def test_every_english_route_requires_all_translations(self):
        page = self.page("guide/index.html", '<html lang="en"><p>English</p>')
        errors = []
        check_translation_routes(self.site, {page.path: page}, errors)
        self.assertEqual(len(errors), 9)
        self.assertIn("fr/guide/index.html: missing translated page", errors)

    def test_knowledge_artifacts_follow_manifest_paths(self):
        directory = self.site / "knowledge"
        directory.mkdir()
        names = {"pages": "pages.json", "chunks": "chunks.json.gz", "index": "index.json.gz"}
        (directory / "manifest.json").write_text(json.dumps({"artifacts": names}))
        for name in names.values():
            (directory / name).write_bytes(b"data")
        errors = []
        check_knowledge_assets(self.site, errors)
        self.assertEqual(errors, [])
        (directory / "index.json.gz").unlink()
        check_knowledge_assets(self.site, errors)
        self.assertEqual(errors, ["Missing or invalid knowledge artifact: index.json.gz"])

    def test_knowledge_manifest_rejects_path_escape(self):
        directory = self.site / "knowledge"
        directory.mkdir()
        (self.site / "outside.json").write_text("[]")
        (directory / "manifest.json").write_text(json.dumps({"artifacts": {
            "pages": "../outside.json", "chunks": "../outside.json", "index": "../outside.json"}}))
        errors = []
        check_knowledge_assets(self.site, errors)
        self.assertEqual(len(errors), 3)

    def test_english_fallback_is_rejected(self):
        page = self.page("fr/guide/index.html", '<html lang="en"><p>English</p>')
        errors = []
        check_translation_routes(self.site, {page.path: page}, errors)
        self.assertEqual(errors, ["fr/guide/index.html: English fallback in a translated route"])


if __name__ == "__main__":
    unittest.main()
