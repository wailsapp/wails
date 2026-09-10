"""Exercise a native Wails development session in a disposable fixture."""

import subprocess, os, time, signal, json, sys, socket
from pathlib import Path
import argparse, shutil, tempfile
parser = argparse.ArgumentParser(description='Native Wails dev acceptance against a disposable fixture')
parser.add_argument('--cli', required=True, type=Path)
parser.add_argument('--output', required=True, type=Path)
parser.add_argument('--repo', type=Path, default=Path(__file__).resolve().parents[3])
parser.add_argument('--expected-app-args', type=json.loads, help='Expected application argument vector as JSON')
parser.add_argument('dev_args', nargs=argparse.REMAINDER)
options = parser.parse_args()
base = options.output.resolve()
base.mkdir(parents=True, exist_ok=True)
root = Path(tempfile.mkdtemp(prefix='app-', dir=base))
shutil.copytree(Path(__file__).parent / 'app', root, dirs_exist_ok=True)
repo = options.repo.resolve()
(root / 'go.mod').write_text('module devacceptance\n\ngo 1.26.2\n\nrequire github.com/wailsapp/wails/v3 v3.0.0\nreplace github.com/wailsapp/wails/v3 => ' + json.dumps(str(repo)) + '\n')
package = json.loads((root / 'frontend/package.json').read_text())
package['dependencies'] = {'@wailsio/runtime': 'file:' + str(repo / 'internal/runtime/desktop/@wailsio/runtime')}
(root / 'frontend/package.json').write_text(json.dumps(package))
(root / 'frontend/dist').mkdir(parents=True)
(root / 'frontend/dist/index.html').write_text('<html></html>')
subprocess.run(['go', 'mod', 'tidy'], cwd=root, check=True)
name = 'acceptance'
args = options.dev_args
if args and args[0] == '--':
    args = args[1:]
log = base / ('dev-native-' + name + '.log')
original = (root / 'main.go').read_text()
env = os.environ.copy()
env['WAILS_EXP_USE_WAKE'] = ''
r = {}
production = root / 'bin' / 'production-sentinel'
production.parent.mkdir(parents=True, exist_ok=True)
production.write_bytes(b'production must remain untouched')
port = int(env.get('WAILS_VITE_PORT', '9245'))
host = '127.0.0.1'
network_parser = argparse.ArgumentParser(add_help=False)
network_parser.add_argument('--port', '-port', type=int, default=port)
network_parser.add_argument('--host', '-host', default=host)
network_options, _ = network_parser.parse_known_args(args)
port, host = network_options.port, network_options.host
frontend = root / 'frontend/src/main.ts'
original_frontend = frontend.read_text()
with log.open('w') as f:
    p = subprocess.Popen([str(options.cli.resolve()), 'dev', *args], cwd=root, env=env, stdout=f, stderr=subprocess.STDOUT, **{'creationflags': subprocess.CREATE_NEW_PROCESS_GROUP} if os.name == 'nt' else {'start_new_session': True})

    def wait(text, count=1, timeout=180):
        deadline = time.time() + timeout
        while time.time() < deadline:
            if log.read_text().count(text) >= count:
                return True
            if p.poll() is not None:
                return False
            time.sleep(0.3)
        return False
    try:
        r['runtime_ready'] = wait('Backend built and started', timeout=600)
        r['binding_roundtrip'] = wait('HCL_BINDING_ROUNDTRIP', timeout=15) if r['runtime_ready'] else False
        r['stderr_logging'] = wait('HCL_STDERR_LOG_PASS', timeout=15) if r['runtime_ready'] else False
        if options.expected_app_args is not None:
            marker = 'HCL_APP_ARGS ' + json.dumps(options.expected_app_args, separators=(',', ':'), ensure_ascii=False)
            r['initial_application_args'] = wait(marker, timeout=15) if r['runtime_ready'] else False
        if r['runtime_ready']:
            frontend.write_text(original_frontend + "\nCall.ByName('main.AcceptanceService.Confirm', 'HCL_HMR_PASS');\n")
            r['frontend_hmr'] = wait('HCL_BINDING_ROUNDTRIP HCL_HMR_PASS', timeout=30)
            r['hmr_without_backend_restart'] = 'Backend rebuilt and restarted' not in log.read_text()
            baseline = log.read_text().count('HCL_BINDING_ROUNDTRIP')
            (root / 'main.go').write_text(original + '\nfunc invalid syntax\n')
            r['failed_build_keeps_app'] = wait('Current app is still running.', timeout=90)
            (root / 'main.go').write_text(original + '\nvar devAcceptanceGeneration = 2\n')
            r['rebuild_restart'] = wait('Backend rebuilt and restarted', timeout=300)
            r['second_binding_roundtrip'] = wait('HCL_BINDING_ROUNDTRIP', baseline + 1, 15)
            if options.expected_app_args is not None:
                r['restarted_application_args'] = wait(marker, count=2, timeout=15)
    finally:
        if p.poll() is None:
            p.send_signal(signal.CTRL_BREAK_EVENT if os.name == 'nt' else signal.SIGINT)
            try:
                p.wait(timeout=25)
            except subprocess.TimeoutExpired:
                r['graceful_shutdown'] = False
                if os.name == 'nt':
                    subprocess.run(['taskkill', '/T', '/F', '/PID', str(p.pid)], check=False)
                else:
                    os.killpg(p.pid, signal.SIGTERM)
                p.wait(timeout=10)
        (root / 'main.go').write_text(original)
        frontend.write_text(original_frontend)
    r['exit'] = p.returncode
    r.setdefault('graceful_shutdown', p.returncode == 0)
    r['production_unchanged'] = production.read_bytes() == b'production must remain untouched'
    try:
        with socket.socket() as listener:
            if os.name != "nt":
                listener.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            listener.bind((host, port))
        r['port_reusable'] = True
    except OSError:
        r['port_reusable'] = False
(base / ('dev-native-' + name + '.json')).write_text(json.dumps(r, indent=2))
print(json.dumps(r))
if not all((value is True for (key, value) in r.items() if key != 'exit')) or r['exit'] != 0:
    sys.exit(1)
