import importlib.util
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('publisher', Path(__file__).with_name('publish-changelog.py'))
publisher = importlib.util.module_from_spec(spec)
spec.loader.exec_module(publisher)


class PublicationTest(unittest.TestCase):
    def setUp(self):
        """Create an isolated bare remote and two competing publisher checkouts."""
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.git(self.root, 'init', '--bare', '--initial-branch=master', 'remote')
        self.git(self.root, 'clone', str(self.root / 'remote'), 'a')
        self.a = self.root / 'a'
        self.configure(self.a)
        (self.a / 'v3').mkdir()
        (self.a / publisher.CHANGELOG).write_text('## Added\n\n## Fixed\n\n---\n')
        self.commit(self.a, 'initial')
        self.git(self.a, 'push', 'origin', 'master')
        self.git(self.root, 'clone', str(self.root / 'remote'), 'b')
        self.b = self.root / 'b'
        self.configure(self.b)
        self.entry = '- Fix A in [PR](https://github.com/wailsapp/wails/pull/123) by @a'
        self.other = '- Fix B in [PR](https://github.com/wailsapp/wails/pull/456) by @b'
        self.add_entry(self.a, self.entry)

    def git(self, cwd, *args):
        """Run Git and return its output, retaining diagnostics on failure."""
        try:
            return subprocess.check_output(['git', *args], cwd=cwd, stderr=subprocess.STDOUT, text=True)
        except subprocess.CalledProcessError as error:
            self.fail(f'Git command failed in {cwd}: {error}\n{error.output}')

    def configure(self, cwd):
        """Set a local test identity without changing the runner configuration."""
        self.git(cwd, 'config', 'user.name', 'Test')
        self.git(cwd, 'config', 'user.email', 'test@example.com')

    def commit(self, cwd, message):
        """Commit all changes in the disposable test checkout."""
        self.git(cwd, 'add', '.')
        self.git(cwd, 'commit', '-m', message)

    def add_entry(self, cwd, entry):
        """Insert a fixture entry into the shared Fixed section."""
        path = cwd / publisher.CHANGELOG
        path.write_text(path.read_text().replace('## Fixed\n', '## Fixed\n' + entry + '\n'))

    def remote(self):
        """Read the changelog actually published to the bare remote."""
        return self.git(self.root, '--git-dir=remote', 'show', 'master:' + publisher.CHANGELOG)

    def publish(self):
        """Publish the generated PR entry with bounded retries and duplicate detection."""
        publisher.publish(self.a, '123', 'fix: A', attempts=3, delay=0)

    def test_competing_push_and_replay(self):
        """Force a push race and verify both entries land exactly once."""
        self.add_entry(self.b, self.other)
        self.commit(self.b, 'competing entry')
        # Force the competing push AFTER the publisher fetches, immediately
        # before its first real push. Both commits edit the same section.
        hook = self.a / '.git/hooks/pre-push'
        hook.write_text('#!/bin/sh\nrm -- "$0"\nunset $(git rev-parse --local-env-vars)\ngit -C ' + shlex.quote(str(self.b)) + ' push origin master\n')
        hook.chmod(0o755)
        self.publish()
        self.assertEqual(self.remote().count(self.entry), 1)
        self.assertEqual(self.remote().count(self.other), 1)
        head = self.git(self.root, '--git-dir=remote', 'rev-parse', 'master')
        self.publish()
        self.assertEqual(head, self.git(self.root, '--git-dir=remote', 'rev-parse', 'master'))

    def test_regenerated_identical_entry_is_idempotent(self):
        """Verify regenerating an existing entry does not duplicate it."""
        self.publish()
        self.git(self.a, 'restore', publisher.CHANGELOG)
        self.git(self.a, 'pull', '--ff-only')
        self.add_entry(self.a, self.entry)
        self.publish()
        self.assertEqual(self.remote().count(self.entry), 1)

    def test_release_reset_preserves_archived_entries(self):
        """Publish after a release reset without resurrecting archived entries."""
        archive = self.b / publisher.ARCHIVE
        archive.parent.mkdir(parents=True)
        archive.write_text('## beta.9\n' + self.other + '\n')
        self.commit(self.b, 'release resets unreleased changelog')
        self.git(self.b, 'push', 'origin', 'master')
        self.publish()
        self.assertEqual(self.remote().count(self.entry), 1)
        self.assertNotIn(self.other, self.remote())
        self.assertIn(self.other, self.git(self.root, '--git-dir=remote', 'show', 'master:' + publisher.ARCHIVE))

    def test_already_archived_does_not_repeat(self):
        """Avoid republishing an entry already present in released notes."""
        archive = self.b / publisher.ARCHIVE
        archive.parent.mkdir(parents=True)
        archive.write_text('## beta.9\n' + self.entry + '\n')
        self.commit(self.b, 'release already contains entry')
        self.git(self.b, 'push', 'origin', 'master')
        self.publish()
        self.assertNotIn(self.entry, self.remote())

    def assert_refresh_retry(self, operation):
        """Inject one refresh failure and require the entry to reach the remote."""
        real_git = publisher.git
        failures = []

        def flaky_git(cwd, *args):
            """Fail the selected Git operation once, then execute real Git commands."""
            if args[0] == operation and not failures:
                failures.append(operation)
                raise subprocess.CalledProcessError(128, ['git', *args])
            return real_git(cwd, *args)

        with patch.object(publisher, 'git', side_effect=flaky_git):
            self.publish()
        self.assertEqual(failures, [operation])
        self.assertEqual(self.remote().count(self.entry), 1)

    def test_transient_fetch_failure_is_retried(self):
        """Recover from a failed first fetch and publish the entry."""
        self.assert_refresh_retry('fetch')

    def test_transient_checkout_failure_is_retried(self):
        """Recover from a failed checkout and publish the entry."""
        self.assert_refresh_retry('checkout')

    def test_commit_failure_preserves_error_and_removes_dirty_worktree(self):
        """Preserve the commit error and generated input while removing scratch state."""
        hook = self.a / '.git/hooks/pre-commit'
        hook.write_text('#!/bin/sh\necho "commit deliberately rejected" >&2\nexit 1\n')
        hook.chmod(0o755)
        worktrees_before = self.git(self.a, 'worktree', 'list', '--porcelain')
        with self.assertRaises(subprocess.CalledProcessError) as caught:
            self.publish()
        self.assertEqual(caught.exception.cmd[1], 'commit')
        self.assertEqual(self.git(self.a, 'worktree', 'list', '--porcelain'), worktrees_before)
        self.assertIn(self.entry, (self.a / publisher.CHANGELOG).read_text())
        self.assertNotIn(self.entry, self.remote())

    def test_exhausted_retries_fail_and_preserve_entry(self):
        """Fail after bounded retries while retaining the unpublished input."""
        hook = self.root / 'remote/hooks/pre-receive'
        hook.write_text('#!/bin/sh\nexit 1\n')
        hook.chmod(0o755)
        with self.assertRaises(RuntimeError):
            self.publish()
        self.assertIn(self.entry, (self.a / publisher.CHANGELOG).read_text())
        self.assertNotIn(self.entry, self.remote())


if __name__ == '__main__':
    unittest.main()
