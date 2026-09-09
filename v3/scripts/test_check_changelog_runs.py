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


if __name__ == '__main__':
    unittest.main()
