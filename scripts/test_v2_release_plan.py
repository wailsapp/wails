import importlib.util
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('v2plan', Path(__file__).with_name('v2-release-plan.py'))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class ReleasePlanTests(unittest.TestCase):
    def test_source_gate(self):
        with tempfile.TemporaryDirectory() as root:
            def git(*args):
                return subprocess.check_output(['git', '-C', root, *args], text=True).strip()
            git('init', '-q')
            git('config', 'user.name', 'Test')
            git('config', 'user.email', 'test@example.com')
            def commit(path, text):
                file = Path(root, path)
                file.parent.mkdir(parents=True, exist_ok=True)
                file.write_text(text)
                git('add', '.')
                git('commit', '-qm', path)
            commit('v2/app.go', 'first')
            with patch.object(module, 'git', git):
                with self.assertRaises(RuntimeError):
                    module.plan()
                git('tag', 'v2.12.0')
                self.assertEqual(module.plan()['has_commits'], 'false')
                self.assertEqual(module.plan()['next_version'], 'v2.12.1')
                commit('website/docs/page.md', 'docs')
                commit('v3/app.go', 'v3')
                self.assertEqual(module.plan()['has_commits'], 'false')
                commit('v2/app.go', 'fix')
                self.assertEqual(module.plan()['commit_count'], '1')
                self.assertEqual(module.plan()['has_commits'], 'true')
                commit('v2/app.go', 'first')
                self.assertEqual(module.plan()['has_commits'], 'true')
                git('tag', 'v2.12.1')
                self.assertEqual(module.plan()['has_commits'], 'false')
                git('checkout', '-qb', 'other', 'v2.12.0')
                commit('v2/other.go', 'other')
                git('tag', 'v2.99.0')
                git('checkout', '-')
                self.assertEqual(module.plan()['latest_tag'], 'v2.12.1')
                git('rm', 'v2/app.go')
                git('commit', '-qm', 'remove source')
                self.assertEqual(module.plan()['has_commits'], 'true')


if __name__ == '__main__':
    unittest.main()
