"""Fail a nightly release when changelog runs for its commits are incomplete."""
import json
import os
import subprocess


def incomplete_runs(runs, includes):
    # Closing an unmerged PR intentionally skips the workflow's only job.
    """Return unsuccessful runs whose commits belong to the release snapshot."""
    return [run for run in runs
            if run['conclusion'] not in ('success', 'skipped') and includes(run['head_sha'])]


def make_includes(is_ancestor, merge_commits):
    """Build the snapshot test for a changelog run's head commit.

    Auto Changelog runs on pull_request_target, so head_sha is the PR branch
    head. PRs are squash-merged and their branches deleted, so that commit is
    usually neither in master's history nor fetchable. What belongs to the
    release is the PR's merge commit, so fall back to it.

    is_ancestor(sha) returns True/False, or None when Git does not have sha.
    merge_commits(sha) returns the merge commits of merged PRs headed at sha.
    """
    def includes(sha):
        direct = is_ancestor(sha)
        if direct:
            return True
        merges = merge_commits(sha)
        if merges:
            return any(is_ancestor(merge) for merge in merges)
        if direct is None:
            # Unknown commit and no merged PR: fail closed rather than guess.
            raise RuntimeError(f'Cannot determine whether changelog run commit {sha} is in this release')
        return False
    return includes


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

    def is_ancestor(sha):
        """Check commit ancestry; None when Git cannot establish it."""
        result = subprocess.run(['git', 'merge-base', '--is-ancestor', sha, 'HEAD'],
                                stderr=subprocess.DEVNULL)
        if result.returncode not in (0, 1):
            return None
        return result.returncode == 0

    def merge_commits(sha):
        """Merge commits of the merged PRs whose head is sha."""
        pulls = json.loads(subprocess.check_output(
            ['gh', 'api', f'repos/{repo}/commits/{sha}/pulls'], text=True))
        return [pr['merge_commit_sha'] for pr in pulls
                if pr.get('merged_at') and pr.get('merge_commit_sha')]

    missing = incomplete_runs(runs, make_includes(is_ancestor, merge_commits))
    if missing:
        for run in missing:
            print(f"Changelog run requires recovery: {run['html_url']} ({run['status']}/{run['conclusion']})")
        raise SystemExit('Release blocked: rerun/recover the unsuccessful changelog runs first.')
    print('All observed changelog runs for this release snapshot succeeded.')


if __name__ == '__main__':
    main()
