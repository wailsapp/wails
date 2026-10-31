"""Generate release-local notes from PRs whose merge commits are in a Git range.

All GitHub operations are reads. No model, merge-time writer, or shared staging
file is involved. The receipt records the exact source, baseline and selection.
"""
import argparse
import hashlib
import html
import json
import os
from pathlib import Path
import re
import subprocess

SECTIONS = ('Added', 'Changed', 'Fixed', 'Deprecated', 'Removed', 'Security')
INTERNAL = {'build', 'chore', 'ci', 'docs', 'style', 'test'}
CONVENTIONAL = re.compile(r'^([a-z]+)(?:\([^)]*\))?!?:\s*(.+)$', re.I)
# Keep release boundaries aligned with release.go's releaseTagShape.
VERSION = re.compile(r'v3\.\d+\.\d+(?:-(?:alpha\d*|beta|rc)\.\d+|-nightly\.\d{8}(?:\.\d+)?)?$')
ARCHIVE = 'docs/mpress/content/changelog.md'


def git(root, *args):
    return subprocess.check_output(['git', '-C', str(root), *args], text=True).strip()


def github_pages(endpoint):
    """gh follows every Link page; failures never become an empty changelog."""
    output = subprocess.check_output([
        'gh', 'api', '--method', 'GET', '--paginate', '--slurp', endpoint,
        '-H', 'Accept: application/vnd.github+json',
    ], text=True, timeout=120)
    pages = json.loads(output)
    if not isinstance(pages, list) or any(not isinstance(page, list) for page in pages):
        raise ValueError(f'Expected paginated lists from {endpoint}')
    return [item for page in pages for item in page]


def resolve(root, ref):
    return git(root, 'rev-parse', '--verify', '--end-of-options', f'{ref}^{{commit}}')


def baseline(root, target, stable):
    """Use the nearest reachable v3 tag; stable releases include all nightlies."""
    tags = [tag for tag in git(root, 'tag', '--list', 'v3.*').splitlines()
            if VERSION.fullmatch(tag) and (not stable or '-' not in tag)]
    if not tags:
        raise ValueError('No eligible v3 release tag; use --since for an explicit bootstrap baseline')
    # Exact tag patterns avoid treating unrelated v3.* tags as release boundaries.
    try:
        return git(root, 'describe', '--tags', '--abbrev=0',
                   *[arg for tag in tags for arg in ('--match', tag)], target)
    except subprocess.CalledProcessError as error:
        raise ValueError('No reachable v3 release tag; use --since to bootstrap explicitly') from error


def section_and_title(pr):
    """Explicit category labels override conventional title classification."""
    labels = {label['name'].lower() for label in pr.get('labels', [])}
    if 'skip-changelog' in labels:
        return None
    match = CONVENTIONAL.match(pr['title'].strip())
    kind, title = (match[1].lower(), match[2]) if match else ('', pr['title'])
    explicit = [section for section in SECTIONS if f'changelog:{section.lower()}' in labels]
    if len(explicit) > 1:
        raise ValueError(f'PR #{pr["number"]} has conflicting changelog categories')
    if not explicit and kind in INTERNAL:
        return None
    section = explicit[0] if explicit else {'feat': 'Added', 'fix': 'Fixed',
                'security': 'Security', 'deprecate': 'Deprecated', 'remove': 'Removed'}.get(kind, 'Changed')
    # A dedicated Release notes paragraph lets maintainers refine a generic title
    # without editing generated files. Other PR-template sections are ignored.
    summary = re.search(r'^## Release notes[ \t]*\n(.*?)(?=^## |\Z)', pr.get('body') or '', re.M | re.S | re.I)
    if summary:
        text = re.sub(r'<!--.*?-->', '', summary[1], flags=re.S).strip()
        if text:
            title = text.removeprefix('- ').removeprefix('* ')
    # Preserve inline Markdown but prevent raw HTML and line/heading injection.
    title = html.escape(' '.join(title.split()), quote=False)
    if not title:
        raise ValueError(f'PR #{pr["number"]} has no release-note title')
    return section, title


def collect(root, repo, start, target, api=github_pages):
    """Require actual merge-commit membership, including squash/rebase merges.

    A closed PR may be unmerged, or merged later than the tested snapshot.
    Require merge-commit membership; never fetch or advance the tested target.
    """
    commits = git(root, 'rev-list', '--reverse', f'{start}..{target}').splitlines()
    included = set(commits)
    selected = {}
    # Listing PRs once avoids an API round trip for every commit in the range.
    for pr in api(f'repos/{repo}/pulls?state=closed&per_page=100'):
        if (not pr.get('merged_at') or pr.get('merge_commit_sha') not in included
                or pr['base']['repo']['full_name'].lower() != repo.lower()):
            continue
        selected[pr['number']] = pr
    entries = []
    skipped = []
    for number, pr in sorted(selected.items()):
        category = section_and_title(pr)
        if category is None:
            skipped.append(number)
            continue
        files = api(f'repos/{repo}/pulls/{number}/files?per_page=100')
        # GitHub caps this API at 3,000 files. Do not silently process a truncation.
        if len(files) >= 3000:
            raise ValueError(f'PR #{number} reaches the GitHub changed-files limit')
        if not any(item.get(key, '').startswith(('v3/', 'docs/mpress/'))
                   for item in files for key in ('filename', 'previous_filename')):
            skipped.append(number)
            continue
        author = pr['user']['login']
        if not re.fullmatch(r'[A-Za-z0-9-]+(?:\[bot\])?', author):
            raise ValueError(f'Invalid author for PR #{number}')
        section, title = category
        entries.append({'number': number, 'section': section, 'title': title,
                        'author': author, 'merge_commit': pr['merge_commit_sha']})
    return entries, skipped, commits


def migrate_pending(root, repo, commits, entries):
    """Retain approved pre-cutover prose only while its merge is in range."""
    path = root / 'v3/tasks/release/changelog-migration.json'
    if not path.exists():
        return entries
    seed = json.loads(path.read_text())
    if seed['repository'] != repo:
        return entries
    pending = {entry['number']: entry for entry in seed['entries']
               if entry['merge_commit'] in commits}
    current = {entry['number']: entry for entry in entries}
    current.update(pending)
    return [entry for _, entry in sorted(current.items())]


def render(repo, entries):
    lines = ['# Unreleased Changes', '']
    for section in SECTIONS:
        matching = [entry for entry in entries if entry['section'] == section]
        if not matching:
            continue
        lines.extend([f'## {section}', ''])
        for entry in matching:
            lines.append(f'- {entry["title"]} in [PR](https://github.com/{repo}/pull/{entry["number"]}) by @{entry["author"]}')
        lines.append('')
    return '\n'.join(lines) + '\n'


def generate(root, repo, target='HEAD', channel='nightly', version='', since=None, api=github_pages):
    # Invalidate outputs from any earlier invocation before doing fallible work.
    for name in ('release_changelog.md', 'release_changelog.json'):
        (root / 'v3' / name).unlink(missing_ok=True)
    if not re.fullmatch(r'[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+', repo):
        raise ValueError('Expected an explicit owner/repository')
    target = resolve(root, target)
    if target != resolve(root, 'HEAD'):
        raise ValueError('Target must equal the checked-out, tested HEAD')
    if version and not VERSION.fullmatch('v' + version.removeprefix('v')):
        raise ValueError('Expected a v3 release version')
    stable = '-' not in version if version else channel == 'weekly'
    start_ref = since or baseline(root, target, stable)
    start = resolve(root, start_ref)
    subprocess.run(['git', '-C', str(root), 'merge-base', '--is-ancestor', start, target], check=True)
    entries, skipped, commits = collect(root, repo, start, target, api)
    entries = migrate_pending(root, repo, set(commits), entries)
    notes = render(repo, entries)
    receipt = {'notes_sha256': hashlib.sha256(notes.encode()).hexdigest(), 'repository': repo, 'target': target, 'baseline': start, 'baseline_ref': start_ref,
               'entries': entries, 'skipped_prs': skipped, 'commit_count': len(commits)}
    # Only write after every API page and validation succeeds. Errors leave no
    # apparently successful empty release input in a fresh runner checkout.
    (root / 'v3/release_changelog.json').write_text(json.dumps(receipt, indent=2) + '\n')
    (root / 'v3/release_changelog.md').write_bytes(notes.encode('utf-8'))
    return receipt


def sync_archive(root, tag):
    """Insert only the tagged release section; preserve newer master history."""
    if not VERSION.fullmatch(tag):
        raise ValueError('Expected a v3 release tag')
    tagged = git(root, 'show', f'refs/tags/{tag}:{ARCHIVE}')
    pattern = rf'^## {re.escape(tag)} - [^\n]+\n.*?(?=^## v3\.\d+\.\d+|\Z)'
    match = re.search(pattern, tagged, re.M | re.S)
    if match is None:
        raise ValueError(f'Tag {tag} has no release changelog section')
    path = root / ARCHIVE
    current = path.read_text()
    existing = re.search(pattern, current, re.M | re.S)
    if existing:
        if existing[0].strip() != match[0].strip():
            raise ValueError(f'Existing changelog section for {tag} differs from the tag')
        return
    if current.count('## [Unreleased]\n') != 1:
        raise ValueError('Expected exactly one Unreleased heading')
    path.write_text(current.replace('## [Unreleased]\n', '## [Unreleased]\n\n' + match[0].rstrip() + '\n', 1))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest='command', required=True)
    gen = sub.add_parser('generate')
    gen.add_argument('--repo', default=os.environ.get('GITHUB_REPOSITORY'))
    gen.add_argument('--target', default='HEAD')
    gen.add_argument('--channel', choices=['nightly', 'weekly'], default='nightly')
    gen.add_argument('--version', default='')
    gen.add_argument('--since', help='Explicit ancestor baseline for bootstrapping')
    sync = sub.add_parser('sync')
    sync.add_argument('--tag', required=True)
    args = parser.parse_args()
    root = Path(git(Path.cwd(), 'rev-parse', '--show-toplevel'))
    if args.command == 'sync':
        sync_archive(root, args.tag)
    else:
        if not args.repo:
            parser.error('--repo or GITHUB_REPOSITORY is required; no public repository default')
        receipt = generate(root, args.repo, args.target, args.channel, args.version, args.since)
        print(f'Generated {len(receipt["entries"])} entries from {receipt["baseline_ref"]} to {receipt["target"]}')


if __name__ == '__main__':
    main()
