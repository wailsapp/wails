"""Fail a nightly release when changelog runs for its commits are incomplete."""
import json
import os
import subprocess


def incomplete_runs(runs, includes):
    # Closing an unmerged PR intentionally skips the workflow's only job.
    """Return unsuccessful runs whose commits belong to the release snapshot."""
    return [run for run in runs
            if includes(run['head_sha']) and run['conclusion'] not in ('success', 'skipped')]


def main():
    """Check recorded changelog runs since the previous reachable v3 release tag."""
    tag = subprocess.check_output(
        ['git', 'describe', '--tags', '--abbrev=0', '--match', 'v3.*'], text=True).strip()
    since = subprocess.check_output(['git', 'show', '-s', '--format=%cI', tag], text=True).strip()
    repo = os.environ['GITHUB_REPOSITORY']
    output = subprocess.check_output([
        'gh', 'api', '--paginate', '--slurp', '--method', 'GET',
        f'repos/{repo}/actions/workflows/auto-changelog-v3.yml/runs',
        '-f', f'created=>={since}', '-f', 'per_page=100',
    ], text=True)
    runs = [run for page in json.loads(output) for run in page['workflow_runs']]

    def includes(sha):
        """Check commit ancestry, failing closed when Git cannot establish it."""
        result = subprocess.run(['git', 'merge-base', '--is-ancestor', sha, 'HEAD'])
        if result.returncode not in (0, 1):
            raise RuntimeError(f'Cannot determine whether changelog run commit {sha} is in this release')
        return result.returncode == 0

    missing = incomplete_runs(runs, includes)
    if missing:
        for run in missing:
            print(f"Changelog run requires recovery: {run['html_url']} ({run['status']}/{run['conclusion']})")
        raise SystemExit('Release blocked: rerun/recover the unsuccessful changelog runs first.')
    print('All observed changelog runs for this release snapshot succeeded.')


if __name__ == '__main__':
    main()
