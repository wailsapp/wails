import importlib.util
from pathlib import Path
import shlex
import subprocess
import tempfile
import unittest

spec = importlib.util.spec_from_file_location('publisher', Path(__file__).with_name('publish-changelog.py'))
publisher = importlib.util.module_from_spec(spec)
spec.loader.exec_module(publisher)


class PublicationTest(unittest.TestCase):
    def setUp(self):
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
        return subprocess.check_output(['git', *args], cwd=cwd, stderr=subprocess.DEVNULL, text=True)

    def configure(self, cwd):
        self.git(cwd, 'config', 'user.name', 'Test')
        self.git(cwd, 'config', 'user.email', 'test@example.com')

    def commit(self, cwd, message):
        self.git(cwd, 'add', '.')
        self.git(cwd, 'commit', '-m', message)

    def add_entry(self, cwd, entry):
        path = cwd / publisher.CHANGELOG
        path.write_text(path.read_text().replace('## Fixed\n', '## Fixed\n' + entry + '\n'))

    def remote(self):
        return self.git(self.root, '--git-dir=remote', 'show', 'master:' + publisher.CHANGELOG)

    def publish(self):
        publisher.publish(self.a, '123', 'fix: A', attempts=3, delay=0)

    def test_competing_push_and_replay(self):
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
        self.publish()
        self.git(self.a, 'restore', publisher.CHANGELOG)
        self.git(self.a, 'pull', '--ff-only')
        self.add_entry(self.a, self.entry)
        self.publish()
        self.assertEqual(self.remote().count(self.entry), 1)

    def test_release_reset_preserves_archived_entries(self):
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
        archive = self.b / publisher.ARCHIVE
        archive.parent.mkdir(parents=True)
        archive.write_text('## beta.9\n' + self.entry + '\n')
        self.commit(self.b, 'release already contains entry')
        self.git(self.b, 'push', 'origin', 'master')
        self.publish()
        self.assertNotIn(self.entry, self.remote())

    def test_exhausted_retries_fail_and_preserve_entry(self):
        hook = self.root / 'remote/hooks/pre-receive'
        hook.write_text('#!/bin/sh\nexit 1\n')
        hook.chmod(0o755)
        with self.assertRaises(RuntimeError):
            self.publish()
        self.assertIn(self.entry, (self.a / publisher.CHANGELOG).read_text())
        self.assertNotIn(self.entry, self.remote())


if __name__ == '__main__':
    unittest.main()
