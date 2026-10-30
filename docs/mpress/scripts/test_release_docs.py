"""Release policy and failure safety tests; standard library only."""
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch

import release_docs as docs


class ReleaseDocumentationTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        registry = self.root / docs.REGISTRY
        registry.parent.mkdir(parents=True)
        registry.write_text(json.dumps({"schemaVersion": 1, "latestStable": "v3.0.0", "versions": ["v3.0.0"]}))
        (self.root / "mpress.yaml").write_text('site:\n  baseURL: https://wails.io\n')

    def test_stable_registry_uses_numeric_order_and_never_regresses(self):
        for version in ["v3.2.0", "v3.10.0", "v3.1.1", "v3.10.0"]:
            docs.register(self.root, version)
        self.assertEqual(docs.read_registry(self.root)["versions"], ["v3.10.0", "v3.2.0", "v3.1.1", "v3.0.0"])
        for invalid in ["v3.11.0-nightly.20261030", "--help", "v3.01.0", "v2.0.0"]:
            with self.assertRaises(ValueError):
                docs.register(self.root, invalid)

    def test_invalid_registry_fails_closed(self):
        path = self.root / docs.REGISTRY
        data = docs.read_registry(self.root)
        data["versions"].append("v3.1.0")
        path.write_text(json.dumps(data))
        with self.assertRaises(ValueError):
            docs.read_registry(self.root)

    def test_prepare_registers_only_after_successful_stable_validation(self):
        path = self.root / docs.REGISTRY
        original = path.read_bytes()
        with patch.object(docs, "build", side_effect=RuntimeError("broken links")):
            with self.assertRaises(RuntimeError):
                docs.prepare(self.root, "v3.1.0", "mpress")
        self.assertEqual(path.read_bytes(), original)
        with patch.object(docs, "build"):
            docs.prepare(self.root, "v3.1.0", "mpress", dry_run=True)
            docs.prepare(self.root, "v3.1.0-nightly.20261030", "mpress")
            self.assertEqual(path.read_bytes(), original)
            docs.prepare(self.root, "v3.1.0", "mpress")
        self.assertEqual(docs.read_registry(self.root)["latestStable"], "v3.1.0")

    def test_snapshot_source_is_the_tag_not_head_or_working_tree(self):
        def git(*args):
            return subprocess.check_output(["git", *args], cwd=self.root, text=True).strip()
        git("init", "-q")
        git("config", "user.name", "Fixture")
        git("config", "user.email", "fixture@example.invalid")
        page = self.root / "docs/mpress/content/page.md"
        page.parent.mkdir()
        released = "Released [v2](https://wails.io/docs/introduction) and [API](https://wails.io/wails/runtime.js)"
        page.write_text(released)
        git("add", ".")
        git("commit", "-qm", "release")
        commit = git("rev-parse", "HEAD")
        git("tag", "v3.0.0")
        page.write_text("unreleased")
        destination = self.root / "extracted"
        destination.mkdir()
        self.assertEqual(docs.source_tree(self.root, destination, "v3.0.0"), commit)
        self.assertEqual((destination / "docs/mpress/content/page.md").read_text(),
                         "Released [v2](https://v2.wails.io/docs/introduction) and [API](/guides/custom-transport/)")
        self.assertEqual(git("show", "v3.0.0:docs/mpress/content/page.md"), released)

    def test_search_stays_in_the_selected_version_and_language(self):
        path = self.root / "fr/search-index.json"
        path.parent.mkdir()
        path.write_text(json.dumps([{"title": "Événements", "url": "/fr/reference/events/"}]))
        mounted = self.root / "versions/v3.0.0/search-index.json"
        mounted.parent.mkdir(parents=True)
        mounted.write_text('[{"url": "/versions/v3.0.0/"}]')
        docs.version_search(self.root, "nightly")
        docs.version_search(self.root, "nightly")
        self.assertEqual(json.loads(path.read_text())[0]["url"], "/versions/nightly/fr/reference/events/")
        self.assertEqual(json.loads(mounted.read_text())[0]["url"], "/versions/v3.0.0/")

    def test_corrupted_cache_is_rejected(self):
        host_key = hashlib.sha256(b"https://wails.io").hexdigest()[:12]
        cache = self.root / (".mpress/release-docs/v3.0.0-abc-mpress-1.0.25-layout3-" + host_key)
        cache.mkdir(parents=True)
        (cache / "wails-docs-v3.0.0.tar.gz").write_bytes(b"damaged")
        (cache / "wails-docs-v3.0.0.sha256").write_text(hashlib.sha256(b"original").hexdigest())
        with patch.object(docs, "binary_version", return_value="mpress 1.0.25"), patch.object(docs, "resolve_ref", return_value="abc"):
            with self.assertRaisesRegex(ValueError, "checksum mismatch"):
                docs.snapshot(self.root, "v3.0.0", "mpress", self.root / "output")
        self.assertFalse((self.root / "output").exists())

    def test_failed_site_build_preserves_previous_output(self):
        output = self.root / "published"
        output.mkdir()
        (output / "index.html").write_text("last known good")
        with patch.object(docs, "snapshot", side_effect=RuntimeError("invalid capture")):
            with self.assertRaises(RuntimeError):
                docs.site(self.root, "mpress", output)
        self.assertEqual((output / "index.html").read_text(), "last known good")

    def test_cache_headers_keep_nightly_mutable(self):
        headers = self.root / "_headers"
        headers.write_text('/assets/*\n  X-Test: retained\n\n/versions/*\n  Cache-Control: public, max-age=31536000, immutable\n')
        docs.cache_headers(self.root)
        result = headers.read_text()
        self.assertNotIn("/versions/*\n", result)
        self.assertIn('/versions/nightly/*\n  Cache-Control: no-cache', result)
        self.assertIn('/versions/v3.*\n  Cache-Control: public, max-age=31536000, immutable', result)
        self.assertIn('X-Test: retained', result)


if __name__ == "__main__":
    unittest.main()
