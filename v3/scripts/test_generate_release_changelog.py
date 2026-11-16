"""Real Git histories with fixture GitHub responses; never publish anything."""
import importlib.util
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

SPEC = importlib.util.spec_from_file_location('changelog', Path(__file__).with_name('generate-release-changelog.py'))
MODULE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(MODULE)


class Generation(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        (self.root / 'v3').mkdir()
        self.git('init', '-q', '-b', 'master')
        self.git('config', 'user.name', 'Fixture')
        self.git('config', 'user.email', 'fixture@example.invalid')
        self.start = self.commit('baseline')
        self.git('tag', 'v3.0.0')
        self.repo = 'fixture/private'
        self.responses = {}
        self.requests = []

    def git(self, *args):
        return MODULE.git(self.root, *args)

    def commit(self, message):
        self.git('commit', '-q', '--allow-empty', '-m', message)
        return self.git('rev-parse', 'HEAD')

    def pr(self, number, sha, title='fix: Correct window size', **extra):
        pr = {'number': number, 'merge_commit_sha': sha, 'merged_at': '2026-10-31T00:00:00Z',
              'title': title, 'base': {'repo': {'full_name': self.repo}},
              'user': {'login': 'fixture'}, 'labels': [], **extra}
        self.responses.setdefault(f'repos/{self.repo}/pulls?state=closed&per_page=100', []).append(pr)
        self.responses[f'repos/{self.repo}/pulls/{number}/files?per_page=100'] = [{'filename': 'v3/pkg/application/window.go'}]
        return pr

    def api(self, endpoint):
        self.requests.append(endpoint)
        return self.responses.get(endpoint, [])

    def generate(self, **kwargs):
        return MODULE.generate(self.root, self.repo, api=self.api, **kwargs)

    def notes(self):
        return (self.root / 'v3/release_changelog.md').read_text()

    def test_nightly_and_stable_have_different_boundaries(self):
        first = self.commit('squashed fix')
        self.pr(1, first)
        self.git('tag', 'v3.1.0-nightly.20261030')
        second = self.commit('feature')
        self.pr(2, second, 'feat(window): Add a feature')
        nightly = self.generate()
        self.assertEqual([e['number'] for e in nightly['entries']], [2])
        weekly = self.generate(channel='weekly')
        self.assertEqual([e['number'] for e in weekly['entries']], [1, 2])
        self.assertEqual(weekly['baseline'], self.start)
        explicit = self.generate(version='v3.1.1')
        self.assertEqual(explicit['entries'], weekly['entries'])
        self.assertIn('## Added', self.notes())
        self.assertIn('## Fixed', self.notes())
        self.assertNotIn('fix:', self.notes())

    def test_real_merge_rebase_duplicates_open_late_and_foreign_prs(self):
        self.git('switch', '-q', '-c', 'feature')
        feature = self.commit('branch feature')
        self.git('switch', '-q', 'master')
        self.git('merge', '-q', '--no-ff', 'feature', '-m', 'merge feature')
        merge = self.git('rev-parse', 'HEAD')
        pr = self.pr(2, merge, 'feat: Add feature')
        self.responses[f'repos/{self.repo}/pulls?state=closed&per_page=100'].append(pr)
        rebase_first = self.commit('rebased first')
        rebase_last = self.commit('rebased last')
        rebased = self.pr(3, rebase_last)
        self.responses[f'repos/{self.repo}/pulls?state=closed&per_page=100'].append(rebased)
        self.pr(4, rebase_first, merged_at=None)
        self.pr(5, rebase_first, merge_commit_sha='f' * 40)
        self.pr(6, rebase_first, base={'repo': {'full_name': 'other/repo'}})
        receipt = self.generate()
        self.assertEqual([e['number'] for e in receipt['entries']], [2, 3])
        self.assertEqual(self.notes().count('/pull/2)'), 1)
        self.assertEqual(receipt['target'], rebase_last)

    def test_category_overrides_internal_v2_and_empty_notes(self):
        sha = self.commit('many associations')
        self.pr(1, sha, 'chore: release bookkeeping')
        self.pr(2, sha, 'fix: v2-only bug')
        self.responses[f'repos/{self.repo}/pulls/2/files?per_page=100'] = [{'filename': 'v2/file.go'}]
        self.pr(3, sha, 'fix: hidden', labels=[{'name': 'skip-changelog'}])
        receipt = self.generate()
        self.assertEqual(receipt['entries'], [])
        self.assertEqual(receipt['skipped_prs'], [1, 2, 3])
        self.assertEqual(self.notes(), '# Unreleased Changes\n\n')
        self.pr(4, sha, 'chore: Fix a vulnerability', labels=[{'name': 'changelog:security'}])
        self.assertEqual(self.generate()['entries'][0]['section'], 'Security')

    def test_snapshot_mismatch_nonancestor_and_missing_baseline_fail(self):
        sha = self.commit('fix')
        with self.assertRaisesRegex(ValueError, 'tested HEAD'):
            self.generate(target=self.start)
        self.git('tag', '-d', 'v3.0.0')
        with self.assertRaisesRegex(ValueError, 'bootstrap'):
            self.generate()
        self.assertEqual(self.generate(since=self.start)['target'], sha)
        self.git('switch', '-q', '--orphan', 'unrelated')
        unrelated = self.commit('unrelated')
        with self.assertRaises(subprocess.CalledProcessError):
            self.generate(since=sha)
        self.assertNotEqual(unrelated, sha)

    def test_api_failures_and_invalid_payloads_do_not_write_notes(self):
        sha = self.commit('fix')
        self.pr(1, sha)
        def failure(endpoint):
            raise subprocess.CalledProcessError(1, ['gh', 'api', endpoint])
        with self.assertRaises(subprocess.CalledProcessError):
            MODULE.generate(self.root, self.repo, api=failure)
        self.assertFalse((self.root / 'v3/release_changelog.md').exists())
        self.pr(2, sha, labels=[{'name': 'changelog:fixed'}, {'name': 'changelog:added'}])
        with self.assertRaisesRegex(ValueError, 'conflicting'):
            self.generate()
        self.assertFalse((self.root / 'v3/release_changelog.md').exists())

    def test_changed_files_limit_fails_closed(self):
        sha = self.commit('large PR')
        self.pr(1, sha)
        self.responses[f'repos/{self.repo}/pulls/1/files?per_page=100'] *= 3000
        with self.assertRaisesRegex(ValueError, 'limit'):
            self.generate()

    def test_render_escapes_html_and_newlines_and_is_repeatable(self):
        sha = self.commit('fix')
        self.pr(1, sha, 'fix: Correct <script>\n## injected heading')
        first = self.generate()
        notes = self.notes()
        self.assertIn('&lt;script&gt; ## injected heading', notes)
        self.assertEqual(self.generate(), first)
        self.assertEqual(self.notes(), notes)
        receipt = json.loads((self.root / 'v3/release_changelog.json').read_text())
        self.assertEqual(receipt, first)
        self.assertNotIn('UNRELEASED_CHANGELOG', self.git('status', '--porcelain'))

    def test_paginated_gh_command_reads_all_pages(self):
        with tempfile.TemporaryDirectory() as directory:
            executable = Path(directory) / 'gh'
            executable.write_text('#!/usr/bin/env python3\nimport json, sys\nassert "--paginate" in sys.argv and "--slurp" in sys.argv\nprint(json.dumps([[{"number": 1}], [{"number": 2}]]))\n')
            executable.chmod(0o755)
            before = os.environ['PATH']
            try:
                os.environ['PATH'] = directory + os.pathsep + before
                self.assertEqual(MODULE.github_pages('repos/fixture/private/example'), [{'number': 1}, {'number': 2}])
                executable.write_text('#!/usr/bin/env python3\nprint("[{\\\"message\\\": \\\"failure\\\"}]")\n')
                with self.assertRaises(ValueError):
                    MODULE.github_pages('repos/fixture/private/example')
            finally:
                os.environ['PATH'] = before

    def test_release_notes_description_and_all_categories(self):
        sha = self.commit('categorized PRs')
        for number, kind in enumerate(['feat', 'perf', 'fix', 'deprecate', 'remove', 'security'], 1):
            self.pr(number, sha, f'{kind}: Summary {number}')
        self.pr(7, sha, 'fix: Generic title', body='## Description\nIgnore this\n## Release notes\n<!-- ignored -->\n- Correct the user-visible failure.\n## Tests\nIgnore tests')
        receipt = self.generate()
        self.assertEqual({entry['section'] for entry in receipt['entries']}, set(MODULE.SECTIONS))
        self.assertIn('Correct the user-visible failure.', self.notes())
        self.assertNotIn('Generic title', self.notes())
        self.assertNotIn('Ignore tests', self.notes())

    def test_renames_out_of_v3_still_count_as_v3_changes(self):
        sha = self.commit('move code')
        self.pr(1, sha, 'remove: Move legacy v3 tool out of the module')
        self.responses[f'repos/{self.repo}/pulls/1/files?per_page=100'] = [
            {'filename': 'tools/legacy.go', 'previous_filename': 'v3/legacy.go'}]
        self.assertEqual(self.generate()['entries'][0]['section'], 'Removed')

    def test_failed_regeneration_invalidates_previous_outputs(self):
        sha = self.commit('fix')
        self.pr(1, sha)
        self.generate()
        def failure(endpoint):
            raise RuntimeError('GitHub unavailable')
        with self.assertRaisesRegex(RuntimeError, 'unavailable'):
            MODULE.generate(self.root, self.repo, api=failure)
        self.assertFalse((self.root / 'v3/release_changelog.md').exists())
        self.assertFalse((self.root / 'v3/release_changelog.json').exists())

    def test_unreachable_tags_and_unrelated_tag_names_are_not_boundaries(self):
        sha = self.commit('fix')
        self.pr(1, sha)
        self.git('tag', 'v3.999.0-tui')
        self.git('tag', 'v3.not-a-release')
        self.assertEqual(self.generate(channel='weekly')['baseline_ref'], 'v3.0.0')
        self.assertEqual(self.generate()['baseline_ref'], 'v3.0.0')
        self.git('switch', '-q', '--orphan', 'other-history')
        self.commit('different repository history')
        self.git('tag', 'v3.999.0')
        self.git('switch', '-q', 'master')
        self.assertEqual(self.generate(channel='weekly')['baseline_ref'], 'v3.0.0')

    def test_cutover_preserves_pending_prose_only_in_its_release_range(self):
        sha = self.commit('pending fix')
        self.pr(1, sha, 'fix: A differently worded PR title')
        entry = {'number': 1, 'section': 'Fixed', 'title': 'Original approved pending entry',
                 'author': 'fixture', 'merge_commit': sha}
        seed = self.root / 'v3/tasks/release/changelog-migration.json'
        seed.parent.mkdir(parents=True)
        seed.write_text(json.dumps({'repository': self.repo, 'entries': [entry]}))
        receipt = self.generate()
        self.assertEqual(receipt['entries'], [entry])
        self.assertEqual(self.notes().count('Original approved pending entry'), 1)
        self.assertNotIn('differently worded', self.notes())
        self.assertEqual(MODULE.migrate_pending(self.root, 'other/repo', {sha}, []), [])
        self.git('tag', 'v3.1.0-nightly.20261030')
        self.commit('release bookkeeping')
        self.assertEqual(self.generate()['entries'], [])
        self.assertEqual(self.generate(channel='weekly')['entries'], [entry])

    def test_archive_sync_preserves_newer_sections_is_idempotent_and_rejects_conflicts(self):
        path = self.root / MODULE.ARCHIVE
        path.parent.mkdir(parents=True)
        header = 'Header\n\n## [Unreleased]\n\n'
        old = '## v3.0.0 - 2026-10-01\n\n## Fixed\n- Old fix\n'
        new = '## v3.1.0 - 2026-10-31\n\n## Added\n- Weekly feature\n\n'
        path.write_text(header + new + old)
        self.git('add', MODULE.ARCHIVE)
        self.git('commit', '-q', '-m', 'weekly release')
        self.git('tag', 'v3.1.0')
        concurrent = '## v3.2.0-nightly.20261101 - 2026-11-01\n\n## Fixed\n- Concurrent fix\n\n'
        path.write_text(header + concurrent + old)
        MODULE.sync_archive(self.root, 'v3.1.0')
        expected = path.read_text()
        self.assertIn(new.strip(), expected)
        self.assertIn(concurrent.strip(), expected)
        MODULE.sync_archive(self.root, 'v3.1.0')
        self.assertEqual(path.read_text(), expected)
        path.write_text(expected.replace('Weekly feature', 'Conflicting edit'))
        with self.assertRaisesRegex(ValueError, 'differs'):
            MODULE.sync_archive(self.root, 'v3.1.0')
        self.git('tag', 'v3.2.0')
        with self.assertRaisesRegex(ValueError, 'no release changelog section'):
            MODULE.sync_archive(self.root, 'v3.2.0')


if __name__ == '__main__':
    unittest.main()
