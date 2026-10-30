#!/usr/bin/env python3
"""Release-tagged M-Press snapshots and a stable-default, HEAD/nightly site."""
import argparse
import hashlib
import io
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tarfile
import tempfile

REGISTRY = Path("docs/mpress/releases.json")
STABLE = re.compile(r"v3\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\Z")
RELEASE = re.compile(r"v3\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z][0-9A-Za-z.-]*)?\Z")
DOC_PATHS = ("mpress.yaml", "docs/mpress")
SNAPSHOT_LAYOUT = 3


def run(args, root, **kwargs):
    subprocess.run([str(arg) for arg in args], cwd=root, check=True, **kwargs)


def stable_key(version):
    if not STABLE.fullmatch(version):
        raise ValueError(f"not a stable Wails v3 version: {version}")
    return tuple(int(part) for part in version[1:].split("."))


def read_registry(root):
    data = json.loads((root / REGISTRY).read_text())
    versions = data["versions"]
    if data.get("schemaVersion") != 1 or not versions or len(set(versions)) != len(versions):
        raise ValueError("invalid documentation release registry")
    if data["latestStable"] != max(versions, key=stable_key):
        raise ValueError("latestStable must name the highest registered stable version")
    return data


def register(root, version):
    stable_key(version)
    data = read_registry(root)
    data["versions"] = sorted(set(data["versions"] + [version]), key=stable_key, reverse=True)
    data["latestStable"] = data["versions"][0]
    destination = root / REGISTRY
    temporary = destination.with_suffix(".json.tmp")
    temporary.write_text(json.dumps(data, indent=2) + "\n")
    temporary.replace(destination)


def replace_section(source, key, section):
    # Only replace a complete top-level stanza; M-Press still parses/validates YAML.
    pattern = re.compile(r"^" + re.escape(key) + r":[^\n]*\n(?:[ \t][^\n]*\n|\n)*", re.M)
    if pattern.search(source):
        return pattern.sub(lambda _: section + "\n", source, count=1)
    return source.rstrip() + "\n\n" + section + "\n"


def site_base_url(root):
    source = (root / "mpress.yaml").read_text()
    match = re.search(r"^  baseURL:\s*[\"']?([^\"'\n]+)", source, re.M)
    if not match:
        raise ValueError("site.baseURL is missing")
    return match.group(1).strip().rstrip("/")


def configure(root, label, base_url=None):
    path = root / "mpress.yaml"
    source = path.read_text()
    source = replace_section(source, "versioning", "versioning:\n  enabled: true\n  current: " + label + "\n  artifactsDir: .mpress/versions")
    if base_url:
        source = re.sub(r"^(  baseURL:)[^\n]*", lambda match: match[1] + " " + json.dumps(base_url), source, count=1, flags=re.M)
    path.write_text(source)


def resolve_ref(root, ref):
    stable_key(ref)
    resolved = subprocess.run(["git", "rev-parse", "--verify", "--end-of-options", "refs/tags/" + ref + "^{commit}"], cwd=root, text=True, capture_output=True)
    if resolved.returncode:
        run(["git", "fetch", "--no-tags", "https://github.com/wailsapp/wails.git", f"refs/tags/{ref}:refs/tags/{ref}"], root)
        resolved = subprocess.run(["git", "rev-parse", "--verify", "refs/tags/" + ref + "^{commit}"], cwd=root, text=True, capture_output=True, check=True)
    return resolved.stdout.strip()


def source_tree(root, destination, ref=None):
    if ref is None:
        shutil.copy2(root / "mpress.yaml", destination / "mpress.yaml")
        shutil.copytree(root / "docs/mpress", destination / "docs/mpress", ignore=shutil.ignore_patterns("site", "site.new", "site.previous", "releases", ".mpress", "__pycache__"))
        return None
    commit = resolve_ref(root, ref)
    data = subprocess.check_output(["git", "archive", commit, *DOC_PATHS], cwd=root)
    with tarfile.open(fileobj=io.BytesIO(data)) as archive:
        archive.extractall(destination, filter="data")
    if ref == "v3.0.0":
        # Routing corrections only; preserve the tagged prose and examples.
        for page in (destination / "docs/mpress/content").rglob("*.md"):
            source = page.read_text()
            corrected = source.replace("](https://wails.io/docs/", "](https://v2.wails.io/docs/")
            corrected = corrected.replace("](https://wails.io/wails/runtime.js)", "](/guides/custom-transport/)")
            if corrected != source:
                page.write_text(corrected)
    return commit


def build(project, mpress):
    run([mpress, "build", "--strict", "--no-purge-css"], project)
    run([mpress, "check"], project)
    return project / "docs/mpress/site"


def version_search(output, label):
    # M-Press rebases HTML when mounting, but search result URLs live in JSON.
    for path in output.rglob("search-index.json"):
        if path.relative_to(output).parts[0] == "versions":
            continue
        entries = json.loads(path.read_text())
        for entry in entries:
            url = entry["url"]
            if url.startswith("/") and not url.startswith("/versions/"):
                entry["url"] = "/versions/" + label + url
        path.write_text(json.dumps(entries, ensure_ascii=False) + "\n")


def binary_version(mpress):
    return subprocess.check_output([str(mpress), "version"], text=True).strip()


def snapshot(root, version, mpress, destination):
    stable_key(version)
    renderer = binary_version(mpress)
    commit = resolve_ref(root, version)
    base_url = site_base_url(root)
    host_key = hashlib.sha256(base_url.encode()).hexdigest()[:12]
    cache = root / ".mpress/release-docs" / (version + "-" + commit + "-" + renderer.replace(" ", "-") + f"-layout{SNAPSHOT_LAYOUT}-" + host_key)
    cached_bundle = cache / f"wails-docs-{version}.tar.gz"
    if cached_bundle.exists():
        checksum = (cache / f"wails-docs-{version}.sha256").read_text().split()[0]
        if hashlib.sha256(cached_bundle.read_bytes()).hexdigest() != checksum:
            raise ValueError(f"cached documentation snapshot checksum mismatch: {version}")
        destination.mkdir(parents=True, exist_ok=True)
        for item in cache.iterdir():
            if item.is_file():
                shutil.copy2(item, destination / item.name)
        return destination / cached_bundle.name
    with tempfile.TemporaryDirectory(prefix="wails-docs-release-") as temporary:
        project = Path(temporary)
        commit = source_tree(root, project, version)
        configure(project, version, base_url + "/versions/" + version)
        version_search(build(project, mpress), version)
        run([mpress, "versions", "capture", version], project)
        run([mpress, "versions", "verify", version], project)
        captured = project / ".mpress/versions" / version
        metadata = {"schemaVersion": 1, "version": version, "sourceTag": version, "sourceCommit": commit, "renderer": binary_version(mpress), "snapshotLayout": SNAPSHOT_LAYOUT, "baseURL": base_url, "routingCorrections": version == "v3.0.0"}
        destination.mkdir(parents=True, exist_ok=True)
        bundle = destination / f"wails-docs-{version}.tar.gz"
        with tarfile.open(bundle, "w:gz") as archive:
            archive.add(captured, arcname=version)
        (destination / f"wails-docs-{version}.json").write_text(json.dumps(metadata, indent=2) + "\n")
        digest = hashlib.sha256(bundle.read_bytes()).hexdigest()
        (destination / f"wails-docs-{version}.sha256").write_text(f"{digest}  {bundle.name}\n")
        cache.parent.mkdir(parents=True, exist_ok=True)
        with tempfile.TemporaryDirectory(prefix="snapshot-", dir=cache.parent) as pending:
            for item in [bundle, destination / f"wails-docs-{version}.json", destination / f"wails-docs-{version}.sha256"]:
                shutil.copy2(item, Path(pending) / item.name)
            if not cache.exists():
                Path(pending).rename(cache)
        print(json.dumps(metadata))
        return bundle


def prepare(root, version, mpress, dry_run=False):
    if not RELEASE.fullmatch(version):
        raise ValueError(f"invalid Wails v3 release version: {version}")
    with tempfile.TemporaryDirectory(prefix="wails-docs-prepare-") as temporary:
        project = Path(temporary)
        source_tree(root, project)
        configure(project, version if STABLE.fullmatch(version) else "nightly")
        build(project, mpress)
    if STABLE.fullmatch(version) and not dry_run:
        register(root, version)


def install_snapshot(bundle, store):
    with tarfile.open(bundle) as archive:
        archive.extractall(store, filter="data")


def cache_headers(output):
    path = output / "_headers"
    source = path.read_text()
    source = re.sub(r"(?m)^/versions/\*\n(?:[ \t][^\n]*\n)*", "", source)
    source += "\n/versions/v3.*\n  Cache-Control: public, max-age=31536000, immutable\n\n/versions/nightly/*\n  Cache-Control: no-cache\n\n/versions/versions.json\n  Cache-Control: no-cache\n"
    path.write_text(source)


def site(root, mpress, output):
    registry = read_registry(root)
    latest = registry["latestStable"]
    with tempfile.TemporaryDirectory(prefix="wails-docs-site-") as temporary:
        stage = Path(temporary)
        releases = stage / "releases"
        releases.mkdir()
        for version in registry["versions"]:
            snapshot(root, version, mpress, releases)
        head = stage / "head"
        head.mkdir()
        source_tree(root, head)
        configure(head, "nightly", site_base_url(root) + "/versions/nightly")
        head_store = head / ".mpress/versions"
        head_store.mkdir(parents=True)
        for version in registry["versions"]:
            install_snapshot(releases / f"wails-docs-{version}.tar.gz", head_store)
            run([mpress, "versions", "verify", version], head)
        version_search(build(head, mpress), "nightly")
        run([mpress, "versions", "capture", "nightly"], head)
        stable = stage / "stable"
        stable.mkdir()
        source_tree(root, stable, latest)
        # Use the live site's hostname while keeping release-tagged content.
        configure(stable, latest, site_base_url(root))
        store = stable / ".mpress/versions"
        store.mkdir(parents=True)
        shutil.copytree(head_store / "nightly", store / "nightly")
        for version in registry["versions"]:
            install_snapshot(releases / f"wails-docs-{version}.tar.gz", store)
        generated = build(stable, mpress)
        cache_headers(generated)
        (generated / "documentation-release.json").write_text(json.dumps({"stable": latest, "nightly": "HEAD", "sourceCommit": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=root, text=True).strip()}, indent=2) + "\n")
        run([sys.executable, root / "docs/mpress/scripts/check_site.py", generated], root)
        # Replace the output only after both views and mounted snapshots pass.
        output.parent.mkdir(parents=True, exist_ok=True)
        staged = output.with_name(output.name + ".new")
        if staged.exists():
            shutil.rmtree(staged)
        shutil.copytree(generated, staged)
        backup = output.with_name(output.name + ".previous")
        if backup.exists():
            shutil.rmtree(backup)
        if output.exists():
            output.rename(backup)
        try:
            staged.rename(output)
        except OSError:
            if backup.exists():
                backup.rename(output)
            raise
        if backup.exists():
            shutil.rmtree(backup)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=["prepare", "register", "snapshot", "build"])
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parents[3])
    parser.add_argument("--mpress", default=os.environ.get("MPRESS_BIN", "mpress"))
    parser.add_argument("--version")
    parser.add_argument("--output", type=Path)
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()
    root = args.root.resolve()
    mpress = shutil.which(args.mpress)
    if args.command != "register" and not mpress:
        parser.error("M-Press executable not found; set MPRESS_BIN or install it")
    if args.command != "build" and not args.version:
        parser.error("--version is required")
    if args.command == "prepare":
        prepare(root, args.version, mpress, args.dry_run)
    elif args.command == "register":
        register(root, args.version)
    elif args.command == "snapshot":
        snapshot(root, args.version, mpress, (args.output or root / "docs/mpress/releases").resolve())
    else:
        site(root, mpress, (args.output or root / "docs/mpress/site").resolve())


if __name__ == "__main__":
    main()
