#!/usr/bin/env python3
"""Plan a v2 release from reachable tags and commits touching v2/ only.

Writes GitHub Actions outputs when GITHUB_OUTPUT is set. Never mutates Git.
"""
import os
import re
import subprocess


def git(*args):
    return subprocess.check_output(['git', *args], text=True).strip()


def plan():
    tags = git('tag', '--merged', 'HEAD', '--list', 'v2.*').splitlines()
    tags = [tag for tag in tags if re.fullmatch(r'v2\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?', tag)]
    if not tags:
        raise RuntimeError('No reachable v2 release tag; establish an explicit release baseline first')
    # Nearest release in this history; a tag from another branch cannot hide changes.
    baseline = git('describe', '--tags', '--abbrev=0', *[arg for tag in tags for arg in ('--match', tag)], 'HEAD')
    revision_range = f'{baseline}..HEAD'
    count = int(git('rev-list', '--count', revision_range, '--', 'v2/'))
    stable = [tag for tag in tags if re.fullmatch(r'v2\.\d+\.\d+', tag)]
    if not stable:
        raise RuntimeError('No reachable stable v2 version')
    major, minor, patch = max(tuple(map(int, tag[1:].split('.'))) for tag in stable)
    next_version = f'v{major}.{minor}.{patch + 1}'
    return {'next_version': next_version, 'latest_tag': baseline, 'commit_count': str(count),
            'has_commits': str(count > 0).lower()}


if __name__ == '__main__':
    result = plan()
    output = '\n'.join(f'{key}={value}' for key, value in result.items()) + '\n'
    print(output, end='')
    if os.environ.get('GITHUB_OUTPUT'):
        with open(os.environ['GITHUB_OUTPUT'], 'a') as stream:
            stream.write(output)
