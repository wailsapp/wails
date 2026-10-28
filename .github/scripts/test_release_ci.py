"""Exercise release CI decisions using the scripts shipped in the workflows."""
import json
import os
from pathlib import Path
import re
import subprocess
import tempfile
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
WORKFLOWS = ROOT / '.github/workflows'


def workflow(name):
    return yaml.safe_load((WORKFLOWS / f'{name}.yml').read_text())


def step(job, name):
    return next(s for s in job['steps'] if s.get('name') == name)


class ReleaseCI(unittest.TestCase):
    def setUp(self):
        self.build = workflow('build-and-test-v3')['jobs']
        self.release = workflow('nightly-release-v3')['jobs']

    def test_pr_scope_uses_author_permissions_and_fails_closed(self):
        script = step(self.build['changes'], 'Choose CI scope from PR author permissions')['with']['script']
        for permission in ['admin', 'maintain', 'write', 'read', 'triage', 'none', 'error']:
            with self.subTest(permission=permission):
                harness = '''
const assert = require('node:assert/strict');
const context = {eventName: 'pull_request', actor: 'maintainer-rerunning-ci',
  repo: {owner: 'wailsapp', repo: 'wails'},
  payload: {pull_request: {user: {login: 'pr-author'}}}};
const core = {setOutput: (k, v) => console.log(JSON.stringify(v)), warning: () => {}, info: () => {}};
const github = {rest: {repos: {getCollaboratorPermissionLevel: async (args) => {
  assert.equal(args.username, 'pr-author');
  if (process.env.TEST_PERMISSION === 'error') throw new Error('API unavailable');
  return {data: {permission: process.env.TEST_PERMISSION}};
}}}};
'''
                result = subprocess.run(['node', '-e', harness + '(async () => {\n' + script + '\n})()'],
                                        env={**os.environ, 'TEST_PERMISSION': permission},
                                        capture_output=True, text=True, check=True)
                self.assertEqual(json.loads(result.stdout), permission not in ['admin', 'maintain', 'write'])

    def test_release_runs_do_not_use_pr_path_filters(self):
        scripts = [
            (step(self.build['changes'], 'Check for changes under v3/')['run'], 'v3=true'),
            (step(workflow('verify-runtime-assets')['jobs']['verify'],
                  'Detect runtime asset or source changes')['run'], 'relevant=true'),
        ]
        for script, expected in scripts:
            for event in ['schedule', 'workflow_dispatch']:
                with self.subTest(expected=expected, event=event), tempfile.TemporaryDirectory() as directory:
                    output = Path(directory) / 'output'
                    result = subprocess.run(['bash', '-e', '-c', script], cwd=directory,
                                            env={**os.environ, 'EVENT_NAME': event, 'GITHUB_OUTPUT': str(output)},
                                            capture_output=True)
                    self.assertEqual(result.returncode, 0)
                    self.assertEqual(output.read_text().strip(), expected)

    def test_results_reject_failed_cancelled_and_skipped_full_checks(self):
        script = self.build['build_results']['steps'][0]['run']
        checks = ['test_go', 'test_js', 'test_templates', 'test_mac_private_apis', 'cross-compile', 'gtk3-menus', 'msix']
        values = {'needs.changes.result': 'success', 'needs.changes.outputs.v3': 'true',
                  'needs.changes.outputs.full': 'true', **{f'needs.{job}.result': 'success' for job in checks}}

        def run(overrides):
            context = {**values, **overrides}
            rendered = re.sub(r'\$\{\{\s*(.*?)\s*\}\}', lambda m: context[m[1]], script)
            return subprocess.run(['bash', '-e', '-c', rendered], capture_output=True).returncode

        self.assertEqual(run({}), 0)
        for job in ['changes', *checks]:
            for status in ['failure', 'cancelled', 'skipped']:
                with self.subTest(job=job, status=status):
                    self.assertNotEqual(run({f'needs.{job}.result': status}), 0)
        quick = {'needs.changes.outputs.full': 'false',
                 **{f'needs.{job}.result': 'skipped' for job in checks[2:]}}
        self.assertEqual(run(quick), 0)
        self.assertNotEqual(run({**quick, 'needs.test_go.result': 'skipped'}), 0)

    def test_all_release_checks_share_one_snapshot(self):
        snapshot = '${{ needs.cadence.outputs.sha }}'
        self.assertEqual(set(self.release['nightly-release']['needs']),
                         {'cadence', 'full-ci', 'runtime-assets', 'changelog-tests', 'security-scan'})
        self.assertTrue(self.release['full-ci']['with']['full'])
        for job in ['full-ci', 'runtime-assets', 'changelog-tests', 'security-scan']:
            self.assertEqual(self.release[job]['with']['ref'], snapshot)
        checkout = step(self.release['nightly-release'], 'Checkout code')
        self.assertEqual(checkout['with']['ref'], snapshot)
        for name in ['build-and-test-v3', 'cross-compile-test-v3', 'gtk3-menu-test-v3',
                     'msix-smoke-v3', 'verify-runtime-assets', 'changelog-publication-tests', 'semgrep']:
            for job in workflow(name)['jobs'].values():
                for item in job.get('steps', []):
                    if item.get('uses', '').startswith('actions/checkout@'):
                        self.assertEqual(item['with']['ref'], '${{ inputs.ref || github.sha }}', name)
        for job in ['cross-compile', 'gtk3-menus', 'msix']:
            self.assertIn(job, self.build['build_results']['needs'])
            self.assertEqual(self.build[job]['with']['ref'], '${{ inputs.ref || github.sha }}')

    def test_branch_movement_cannot_silently_change_release_source(self):
        script = step(self.release['nightly-release'], 'Refuse a stale release snapshot')['run']
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            remote, source, release = [root / name for name in ['remote.git', 'source', 'release']]

            def git(cwd, *args, check=True):
                return subprocess.run(['git', '-C', str(cwd), *args], check=check,
                                      capture_output=True, text=True)

            git(root, 'init', '--bare', str(remote))
            git(root, 'init', '-b', 'master', str(source))
            git(source, 'config', 'user.name', 'CI test')
            git(source, 'config', 'user.email', 'ci@example.invalid')
            (source / 'code').write_text('tested\n')
            git(source, 'add', '.')
            git(source, 'commit', '-m', 'tested snapshot')
            git(source, 'remote', 'add', 'origin', str(remote))
            git(source, 'push', 'origin', 'master')
            git(root, 'clone', '--branch', 'master', str(remote), str(release))
            sha = git(release, 'rev-parse', 'HEAD').stdout.strip()
            git(release, 'checkout', '--detach', sha)
            env = {**os.environ, 'TESTED_SHA': sha, 'RELEASE_REF': 'master'}
            self.assertEqual(subprocess.run(['bash', '-e', '-c', script], cwd=release, env=env,
                                            capture_output=True).returncode, 0)
            (source / 'code').write_text('untested change\n')
            git(source, 'commit', '-am', 'concurrent merge')
            git(source, 'push', 'origin', 'master')
            self.assertNotEqual(subprocess.run(['bash', '-e', '-c', script], cwd=release, env=env,
                                               capture_output=True).returncode, 0)
            self.assertEqual(git(release, 'rev-parse', 'HEAD').stdout.strip(), sha)
            # Even if the branch moves after the guard, the release tool's
            # normal push rejects it without incorporating the new code.
            git(release, 'config', 'user.name', 'CI test')
            git(release, 'config', 'user.email', 'ci@example.invalid')
            (release / 'version').write_text('release metadata\n')
            git(release, 'add', '.')
            git(release, 'commit', '-m', 'release metadata')
            self.assertNotEqual(git(release, 'push', 'origin', 'HEAD:master', check=False).returncode, 0)


if __name__ == '__main__':
    unittest.main()
