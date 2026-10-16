#!/usr/bin/env python3
"""Validate the rendered documentation and Cloudflare Pages upload contract."""

from collections import Counter
from html.parser import HTMLParser
from pathlib import Path
import json
import sys
from urllib.parse import unquote, urljoin, urlsplit

from check_translations import LANGUAGES


class Page(HTMLParser):
    def __init__(self, path):
        super().__init__(convert_charrefs=True)
        self.path = path
        self.lang = ""
        self.canonical = ""
        self.ids = set()
        self.links = []
        self.source = ""
        self.unrendered_d2 = False
        self.unrendered_fence = False
        self.literal_contexts = Counter()
        self.feed(path.read_text(encoding="utf-8"))

    def handle_starttag(self, tag, attrs):
        if tag in {"pre", "code", "script", "style", "svg"}:
            self.literal_contexts[tag] += 1
        attrs = dict(attrs)
        if tag == "code" and "language-d2" in attrs.get("class", "").split():
            self.unrendered_d2 = True
        if "id" in attrs:
            self.ids.add(attrs["id"])
        if tag == "html":
            self.lang = attrs.get("lang", "")
        if tag == "link" and attrs.get("rel") == "canonical":
            self.canonical = attrs.get("href", "")
        if tag == "meta" and attrs.get("name") == "mpress:source":
            self.source = attrs.get("content", "")
        if tag == "a" and attrs.get("href"):
            self.links.append(attrs["href"])

    def handle_endtag(self, tag):
        if self.literal_contexts[tag]:
            self.literal_contexts[tag] -= 1

    def handle_data(self, data):
        if not any(self.literal_contexts.values()) and "```" in data:
            self.unrendered_fence = True


def route(path, site):
    relative = path.relative_to(site).as_posix()
    return "/" + relative.removesuffix("index.html")


def redirects(site):
    result = {}
    for line in (site / "_redirects").read_text().splitlines():
        if line.strip() and not line.startswith("#"):
            source, target, status = line.split()
            if status not in {"301", "302"}:
                raise ValueError("Unexpected redirect status: " + line)
            result[source] = target
    return result


def check_page(page, site, repository, errors):
    relative = page.path.relative_to(site).as_posix()
    if not page.lang:
        errors.append(relative + ": missing language")
    if page.unrendered_d2:
        errors.append(relative + ": D2 diagram was displayed as source code")
    if page.unrendered_fence:
        errors.append(relative + ": Markdown code fence was displayed as prose")
    if relative == "404.html":
        return
    expected = "https://v3.wails.io" + route(page.path, site)
    if page.canonical != expected:
        errors.append(relative + ": incorrect canonical " + page.canonical)
    if page.source:
        source = (repository / page.source).resolve()
        if not source.is_relative_to(repository / "docs/mpress/content") or not source.is_file():
            errors.append(relative + ": invalid contribution source " + page.source)


def check_links(page, site, pages, aliases, errors):
    current = "https://v3.wails.io" + route(page.path, site)
    for link in page.links:
        target = urlsplit(urljoin(current, link))
        if target.scheme not in {"http", "https"} or target.netloc != "v3.wails.io":
            continue
        path = unquote(target.path)
        path = aliases.get(path, path)
        candidate = (site / path.lstrip("/")).resolve()
        if candidate.is_dir():
            candidate /= "index.html"
        if not candidate.is_file() and not candidate.suffix:
            candidate = candidate / "index.html"
        if not candidate.is_relative_to(site) or not candidate.is_file():
            errors.append(str(page.path.relative_to(site)) + ": missing link " + link)
        elif target.fragment and candidate in pages:
            fragment = unquote(target.fragment)
            if fragment not in pages[candidate].ids:
                errors.append(str(page.path.relative_to(site)) + ": missing anchor " + link)


def check_translation_routes(site, pages, errors):
    for path, page in pages.items():
        if page.lang != "en" or path.name == "404.html":
            continue
        relative = path.relative_to(site)
        if relative.parts[0] in LANGUAGES:
            errors.append(relative.as_posix() + ": English fallback in a translated route")
            continue
        for language in LANGUAGES:
            translated = site / language / relative
            target = pages.get(translated)
            if target is None or target.lang != language:
                errors.append(translated.relative_to(site).as_posix() + ": missing translated page")


def check_knowledge_assets(site, errors):
    directory = site / "knowledge"
    try:
        manifest = json.loads((directory / "manifest.json").read_text())
        for kind in ["pages", "chunks", "index"]:
            name = manifest["artifacts"][kind]
            if not isinstance(name, str) or not name:
                raise ValueError("invalid artifact filename")
            path = (directory / name).resolve()
            if not path.is_relative_to(directory.resolve()) or not path.is_file():
                errors.append("Missing or invalid knowledge artifact: " + name)
    except (OSError, ValueError, KeyError, TypeError) as exc:
        errors.append("Invalid knowledge manifest: " + str(exc))


def main():
    site = Path(sys.argv[1] if len(sys.argv) > 1 else "docs/mpress/site").resolve()
    repository = Path(__file__).resolve().parents[3]
    files = [path for path in site.rglob("*") if path.is_file()]
    errors = []
    if not files or len(files) > 20000:
        errors.append("Pages upload must contain 1–20,000 files")
    for path in files:
        if path.stat().st_size > 25 * 1024 * 1024:
            errors.append("File exceeds the Pages 25 MiB limit: " + str(path))
    pages = {path: Page(path) for path in files if path.suffix == ".html"}
    aliases = redirects(site)
    for page in pages.values():
        check_page(page, site, repository, errors)
        check_links(page, site, pages, aliases, errors)
    check_translation_routes(site, pages, errors)
    if site / "index.html" not in pages or pages[site / "index.html"].lang != "en":
        errors.append("The root page must default to English")
    for name in ["search-index.json", "sitemap.xml", "robots.txt", "llms.txt"]:
        if not (site / name).is_file():
            errors.append("Missing generated asset: " + name)
    check_knowledge_assets(site, errors)
    for source, target in aliases.items():
        if not (site / target.lstrip("/") / "index.html").is_file():
            errors.append("Redirect target is missing: " + source + " -> " + target)
    print(json.dumps({"files": len(files), "html_pages": len(pages),
                      "languages": dict(Counter(page.lang for page in pages.values())),
                      "redirects": len(aliases), "errors": errors}, indent=2))
    return bool(errors)


if __name__ == "__main__":
    sys.exit(main())
