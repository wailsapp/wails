"""Publish one generated entry against current master, without merging stale files."""
from collections import Counter
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile
import time

CHANGELOG = 'v3/UNRELEASED_CHANGELOG.md'
ARCHIVE = 'docs/src/content/docs/changelog.mdx'


def git(cwd, *args):
    """Run Git and return its output, retaining diagnostics on failure."""
    return subprocess.check_output(['git', *args], cwd=cwd, text=True).strip()


def publish(root, number, title, attempts=5, delay=2):
    """Publish the generated PR entry with bounded retries and duplicate detection."""
    if not re.fullmatch(r'[1-9][0-9]*', number):
        raise ValueError('PR_NUMBER must be a positive integer')
    repo = os.environ.get('GITHUB_REPOSITORY', 'wailsapp/wails')
    marker = f'https://github.com/{repo}/pull/{number})'
    before = git(root, 'show', f'HEAD:{CHANGELOG}')
    after = (root / CHANGELOG).read_text()
    if before.strip() == after.strip():
        print('No generated entry; nothing to publish.')
        return
    section = None
    entries = []
    remaining = Counter(before.splitlines())
    for line in after.splitlines():
        if line.startswith('## '):
            section = line[3:]
        if line.startswith('- ') and marker in line and remaining[line] <= 0:
            entries.append((section, line))
        remaining[line] -= 1
    if len(entries) != 1 or entries[0][0] not in {
        'Added', 'Changed', 'Fixed', 'Deprecated', 'Removed', 'Security'
    }:
        raise ValueError('Expected exactly one generated PR entry in a changelog section')
    section, entry = entries[0]
    # A detached temporary worktree keeps the caller's generated file and any
    # other local changes intact, including on exhausted retries.
    with tempfile.TemporaryDirectory(prefix='wails-changelog-') as directory:
        work = Path(directory) / 'checkout'
        git(root, 'worktree', 'add', '--detach', str(work), 'HEAD')
        try:
            for attempt in range(attempts):
                try:
                    git(work, 'fetch', 'origin', 'master')
                    git(work, 'checkout', '--detach', 'origin/master')
                except subprocess.CalledProcessError as error:
                    print(f'Git refresh failed ({attempt + 1}/{attempts}): {error}', flush=True)
                    if attempt + 1 < attempts:
                        time.sleep(delay)
                    continue
                current = (work / CHANGELOG).read_text()
                archive = work / ARCHIVE
                if marker in current or (archive.exists() and marker in archive.read_text()):
                    print(f'PR #{number} is already recorded; nothing to publish.')
                    return
                heading = re.compile(r'(^## ' + re.escape(section) + r'\n(?:<!--[^\n]*-->\n)*)', re.M)
                if heading.search(current):
                    current = heading.sub(lambda m: m[0] + entry + '\n', current, count=1)
                elif '\n---\n' in current:
                    current = current.replace('\n---\n', f'\n## {section}\n{entry}\n\n---\n', 1)
                else:
                    raise ValueError('Changelog has neither target section nor footer')
                (work / CHANGELOG).write_text(current)
                git(work, 'add', CHANGELOG)
                git(work, 'commit', '-m', f'chore(changelog): auto-add entry for PR #{number} — {title}')
                result = subprocess.run(['git', 'push', 'origin', 'HEAD:refs/heads/master'], cwd=work)
                if result.returncode == 0:
                    return
                print(f'Push failed ({attempt + 1}/{attempts}); reapplying preserved entry to master.', flush=True)
                if attempt + 1 < attempts:
                    time.sleep(delay)
            raise RuntimeError('Could not publish changelog; generated entry remains in the original checkout')
        finally:
            original_error = sys.exc_info()[0]
            try:
                git(root, 'worktree', 'remove', '--force', str(work))
            except subprocess.CalledProcessError:
                if original_error is None:
                    raise
                print(f'Could not remove temporary worktree {work}; preserving original error.', file=sys.stderr)


if __name__ == '__main__':
    publish(Path.cwd(), os.environ['PR_NUMBER'], os.environ.get('PR_TITLE', ''))
