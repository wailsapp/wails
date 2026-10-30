import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location('gate', Path(__file__).with_name('check-changelog-runs.py'))
gate = importlib.util.module_from_spec(spec)
spec.loader.exec_module(gate)


class EvidenceTest(unittest.TestCase):
    def test_failed_cancelled_pending_and_successful_runs(self):
        """Distinguish successful or intentional skips from incomplete evidence."""
        for conclusion in ['failure', 'cancelled', None, 'timed_out', 'success', 'skipped']:
            run = {'head_sha': 'included', 'conclusion': conclusion}
            self.assertEqual(gate.incomplete_runs([run], lambda sha: True),
                             [] if conclusion in ('success', 'skipped') else [run])

    def test_later_merge_does_not_block_fixed_snapshot(self):
        """Ignore runs belonging to commits outside the release snapshot."""
        run = {'head_sha': 'newer', 'conclusion': None}
        self.assertEqual(gate.incomplete_runs([run], lambda sha: False), [])

    def test_successful_runs_need_no_ancestry_lookup(self):
        """A successful run never blocks, even when its commit is unknown."""
        def unknown(sha):
            raise RuntimeError('should not be called')
        run = {'head_sha': 'deleted-branch-head', 'conclusion': 'success'}
        self.assertEqual(gate.incomplete_runs([run], unknown), [])


class IncludesTest(unittest.TestCase):
    """Map a changelog run's PR head commit to the release snapshot."""

    def includes(self, ancestry, merges):
        return gate.make_includes(lambda sha: ancestry.get(sha), lambda sha: merges.get(sha, []))

    def test_commit_in_history(self):
        self.assertTrue(self.includes({'head': True}, {})('head'))

    def test_squash_merged_pr_with_deleted_branch(self):
        """Unknown PR head, but its merge commit is in the snapshot."""
        self.assertTrue(self.includes({'head': None, 'merge': True}, {'head': ['merge']})('head'))

    def test_squash_merged_after_snapshot(self):
        """PR merged after the snapshot commit: not part of this release."""
        self.assertFalse(self.includes({'head': None, 'merge': False}, {'head': ['merge']})('head'))

    def test_known_commit_outside_snapshot(self):
        self.assertFalse(self.includes({'head': False}, {})('head'))

    def test_unknown_commit_without_merged_pr_fails_closed(self):
        with self.assertRaises(RuntimeError):
            self.includes({'head': None}, {})('head')


if __name__ == '__main__':
    unittest.main()
